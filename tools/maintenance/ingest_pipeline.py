#!/usr/bin/env python3
"""
ingest_pipeline.py — single off-queue entry point that keeps a growing stash
library's metadata current.

It orchestrates the existing off-queue identifier scripts (in ../identify/) in the
ONLY correct order:

  1. external_identify.py     — fingerprint-match scenes against StashDB/TPDB,
                                stamping stash_ids and filling performers/studio/
                                tags/date/title for everything that has a match.
  2. external_filename_parse.py — for the LEFTOVERS that no stash-box knew, link
                                  studio/performers from the path against names you
                                  already have (link-only, never creates records).

These two steps write the SAME scene fields (studio_id, performer_ids, date, title),
so they must NOT run simultaneously — identify is the higher-fidelity source and
runs FIRST so filename-parse only ever touches what identify could not resolve. This
orchestrator runs them strictly sequentially in one process; never wire them to run
concurrently over the same library.

Both steps default to DRY-RUN. This orchestrator is dry-run too unless you pass
--apply (which is forwarded to both steps). Each step's own summary line is echoed,
and a combined summary plus before/after "gap" snapshot (scenes missing stash_id /
missing performers) is printed at the end.

Stdlib only. Safe to run as a NAS cron — see the cron line at the bottom of this
docstring.

Usage:
  # dry-run, both steps, with a starting/ending gap snapshot:
  python3 ingest_pipeline.py --stash-url http://overwatch-stash:9999
  # actually write, both steps:
  python3 ingest_pipeline.py --stash-url http://overwatch-stash:9999 --apply
  # only one step:
  python3 ingest_pipeline.py --stash-url http://overwatch-stash:9999 --steps identify
  # stop the whole chain if a step errors (default: report and continue):
  python3 ingest_pipeline.py --stash-url http://overwatch-stash:9999 --apply --stop-on-error
  # skip the GraphQL gap snapshot (e.g. NAS busy / auth quirks):
  python3 ingest_pipeline.py --stash-url http://overwatch-stash:9999 --no-snapshot

Recommended NAS cron (DSM Task Scheduler or crontab) — nightly, off-peak, applying:
  # 4:15am daily: keep new/unidentified scenes' metadata current.
  15 4 * * *  /usr/local/bin/python3 /volume1/docker/stash-llm/tools/maintenance/ingest_pipeline.py \
      --stash-url http://localhost:9999 --apply --stop-on-error >> /volume1/docker/stash-llm/logs/ingest_pipeline.log 2>&1
  # (Adjust the repo path and python to the container/host. On the NAS, run it from
  #  *inside* the box that hosts stash so --stash-url can be http://localhost:9999.
  #  Do NOT schedule it to overlap a native stash Identify/Scan over the same scenes.)
"""

import argparse
import json
import os
import subprocess
import sys
import urllib.error
import urllib.request

# Resolve sibling scripts relative to THIS file's location in the repo — never
# hardcode absolute paths, so the tool moves with the checkout (laptop / NAS / etc.).
HERE = os.path.dirname(os.path.abspath(__file__))
IDENTIFY_DIR = os.path.normpath(os.path.join(HERE, "..", "identify"))

# step name -> (script path, human label). Order of this dict is the run order
# default; --steps may narrow it but must NOT reorder it (identify before filename).
STEPS = {
    "identify": (os.path.join(IDENTIFY_DIR, "external_identify.py"),
                 "fingerprint identify (StashDB/TPDB)"),
    "filename": (os.path.join(IDENTIFY_DIR, "external_filename_parse.py"),
                 "filename link-only (leftovers)"),
}
STEP_ORDER = ["identify", "filename"]


# ── light GraphQL gap snapshot (single query, count-only) ───────────────────────


def gql(url, query, variables=None, headers=None, timeout=30):
    body = json.dumps({"query": query, "variables": variables or {}}).encode()
    req = urllib.request.Request(url, data=body, method="POST")
    req.add_header("Content-Type", "application/json")
    # stash's stack 403s the default Python-urllib UA — present a browser UA (same
    # pattern as ../identify/ and the other maintenance scripts).
    req.add_header("User-Agent",
                   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "
                   "(KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
    for k, v in (headers or {}).items():
        req.add_header(k, v)
    try:
        with urllib.request.urlopen(req, timeout=timeout) as resp:
            payload = json.loads(resp.read().decode())
    except urllib.error.HTTPError as e:
        raise RuntimeError(f"HTTP {e.code} from {url}: {e.read().decode()[:300]}")
    except urllib.error.URLError as e:
        raise RuntimeError(f"connection error to {url}: {e}")
    if payload.get("errors"):
        raise RuntimeError(f"GraphQL errors: {json.dumps(payload['errors'])[:400]}")
    return payload.get("data") or {}


# One round-trip, count-only (per_page:0 returns the count without the rows), with
# two aliased findScenes selections so the snapshot is a single light query.
GAP_Q = """
query($total: FindFilterType, $no_sid: FindFilterType, $no_perf: FindFilterType) {
  total:        findScenes(filter: $total) { count }
  missing_sid:  findScenes(filter: $no_sid,  scene_filter: { stash_id: { modifier: IS_NULL,  value: "" } }) { count }
  missing_perf: findScenes(filter: $no_perf, scene_filter: { performer_count: { modifier: EQUALS, value: 0 } }) { count }
}
"""


def gap_snapshot(url, api_key):
    """Return {total, missing_stash_id, missing_performers} or raise RuntimeError."""
    endpoint = url.rstrip("/") + "/graphql"
    headers = {"ApiKey": api_key} if api_key else {}
    f = {"per_page": 0}
    data = gql(endpoint, GAP_Q, {"total": f, "no_sid": f, "no_perf": f}, headers)
    return {
        "total": data["total"]["count"],
        "missing_stash_id": data["missing_sid"]["count"],
        "missing_performers": data["missing_perf"]["count"],
    }


def print_snapshot(label, snap):
    print(f"  [{label}] scenes={snap['total']}  "
          f"missing stash_id={snap['missing_stash_id']}  "
          f"missing performers={snap['missing_performers']}")


# ── step runner ─────────────────────────────────────────────────────────────────


def run_step(name, stash_url, stash_api_key, apply_writes, dry_run_subprocess=False):
    """Run one off-queue step as a subprocess with the SAME python interpreter.

    Returns (rc, summary_line). rc is the subprocess return code (or a synthetic
    non-zero on a launch failure). summary_line is the step's own 'summary:' line
    if it emitted one, else a short fallback. Output is streamed-then-captured: we
    capture so we can extract the summary, and echo the full output through.
    """
    script, label = STEPS[name]
    cmd = [sys.executable, script, "--stash-url", stash_url]
    if stash_api_key:
        cmd += ["--stash-api-key", stash_api_key]
    if apply_writes:
        cmd += ["--apply"]

    print(f"\n=== step: {name} — {label} ===")
    print(f"    $ {' '.join(_shellish(c) for c in cmd)}")
    if dry_run_subprocess:
        # --dry-run on the orchestrator itself: show the command, don't launch it.
        print("    (orchestrator --dry-run: not executing this step)")
        return 0, "(not run — orchestrator dry-run)"

    if not os.path.isfile(script):
        msg = f"step script not found: {script}"
        print(f"    ! {msg}")
        return 127, msg

    try:
        proc = subprocess.run(
            cmd,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            encoding="utf-8",
            errors="replace",
        )
    except OSError as e:
        msg = f"failed to launch step '{name}': {e}"
        print(f"    ! {msg}")
        return 126, msg

    out = proc.stdout or ""
    # echo the child's full output, indented so it's clearly nested under the step.
    for line in out.splitlines():
        print(f"    | {line}")

    summary = _extract_summary(out)
    if proc.returncode != 0:
        print(f"    ! step '{name}' exited {proc.returncode}")
    return proc.returncode, summary


def _extract_summary(text):
    """Pull the last line beginning with 'summary:' (both steps print one)."""
    found = None
    for line in text.splitlines():
        if line.strip().lower().startswith("summary:"):
            found = line.strip()
    return found or "(no summary line emitted)"


def _shellish(arg):
    """Best-effort quoting just for the echoed command line (cosmetic only)."""
    return f'"{arg}"' if (" " in arg or not arg) else arg


# ── main ────────────────────────────────────────────────────────────────────────


def main(argv=None):
    ap = argparse.ArgumentParser(
        description="Off-queue ingest orchestrator: identify FIRST, then filename "
                    "link-only for the leftovers. Default dry-run.")
    ap.add_argument("--stash-url", default="http://localhost:9999",
                    help="stash GraphQL base URL (default: http://localhost:9999)")
    ap.add_argument("--stash-api-key", default=None,
                    help="stash API key, passed through to every step + the snapshot")
    ap.add_argument("--apply", action="store_true",
                    help="forward --apply to BOTH steps (default: dry-run, no writes)")
    ap.add_argument("--steps", default="identify,filename",
                    help="comma list of steps to run, from {identify,filename} "
                         "(default: both). Order is always identify-then-filename "
                         "regardless of how you list them — they write the same "
                         "fields and must not be reordered.")
    ap.add_argument("--stop-on-error", action="store_true",
                    help="abort the chain if a step exits non-zero (default: report "
                         "the failure and continue to the next step)")
    ap.add_argument("--no-snapshot", action="store_true",
                    help="skip the before/after GraphQL gap snapshot")
    ap.add_argument("--dry-run", action="store_true",
                    help="orchestrator-level dry-run: print the planned commands and "
                         "snapshots but do not launch the step subprocesses. (Not the "
                         "same as the steps' own dry-run, which is the default unless "
                         "you pass --apply.)")
    args = ap.parse_args(argv)

    # Filenames flow up from the child steps (emoji, en-dashes, …); keep our own
    # stdout UTF-8 so re-printed/captured text never aborts the run on a cp1252 console.
    for stream in (sys.stdout, sys.stderr):
        try:
            stream.reconfigure(encoding="utf-8", errors="replace")
        except (AttributeError, ValueError):
            pass

    # Resolve + validate the requested steps, normalised to the fixed safe order.
    requested = [s.strip().lower() for s in args.steps.split(",") if s.strip()]
    unknown = [s for s in requested if s not in STEPS]
    if unknown:
        ap.error(f"unknown step(s): {', '.join(unknown)} (valid: {', '.join(STEP_ORDER)})")
    if not requested:
        ap.error("no steps selected")
    ordered = [s for s in STEP_ORDER if s in requested]

    mode = "APPLY (writing changes)" if args.apply else "dry-run (no writes)"
    print("ingest_pipeline — off-queue metadata chain")
    print(f"  stash-url : {args.stash_url}")
    print(f"  steps     : {' -> '.join(ordered)}")
    print(f"  mode      : {mode}"
          + ("  [orchestrator --dry-run: steps not launched]" if args.dry_run else ""))

    # 1. starting gap snapshot (one light query; non-fatal on failure).
    start_snap = None
    if not args.no_snapshot:
        print("\ngap snapshot:")
        try:
            start_snap = gap_snapshot(args.stash_url, args.stash_api_key)
            print_snapshot("before", start_snap)
        except RuntimeError as e:
            print(f"  ! snapshot failed (continuing without it): {e}")

    # 2/3. run the steps in fixed order.
    results = []  # (name, rc, summary)
    aborted = False
    for name in ordered:
        rc, summary = run_step(
            name, args.stash_url, args.stash_api_key,
            apply_writes=args.apply, dry_run_subprocess=args.dry_run,
        )
        results.append((name, rc, summary))
        if rc != 0 and args.stop_on_error:
            print(f"\n! --stop-on-error: aborting after step '{name}' (rc={rc})")
            aborted = True
            break

    # 4. final gap snapshot + combined summary.
    end_snap = None
    if not args.no_snapshot and not args.dry_run:
        print("\ngap snapshot:")
        try:
            end_snap = gap_snapshot(args.stash_url, args.stash_api_key)
            print_snapshot("after", end_snap)
            if start_snap:
                d_sid = start_snap["missing_stash_id"] - end_snap["missing_stash_id"]
                d_perf = start_snap["missing_performers"] - end_snap["missing_performers"]
                print(f"  [delta]  stash_id gap closed: {d_sid:+d}   "
                      f"performer gap closed: {d_perf:+d}")
        except RuntimeError as e:
            print(f"  ! snapshot failed: {e}")

    print("\n=== combined summary ===")
    print(f"  mode: {mode}")
    for name, rc, summary in results:
        status = "ok" if rc == 0 else f"FAILED(rc={rc})"
        print(f"  - {name:<8} [{status}] {summary}")
    if aborted:
        print("  (chain aborted early via --stop-on-error)")
    if not args.apply and not args.dry_run:
        print("  (dry-run — re-run with --apply to write changes)")

    # exit non-zero if any step failed, so cron / a wrapper can detect trouble.
    any_failed = any(rc != 0 for _, rc, _ in results)
    return 1 if any_failed else 0


if __name__ == "__main__":
    sys.exit(main())
