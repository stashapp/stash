#!/usr/bin/env python3
"""
tag_consolidate.py — find and merge duplicate / near-duplicate tags in stash.

stash libraries accrete tag variants — case differences, punctuation, singular vs
plural, "Blond Hair" vs "Blonde Hair" — plus subjective families (POV variants,
hair-colour schemes, DP/creampie families). This tool clusters them and:

  - AUTO-MERGES the SAFE clusters (trivial formatting / plural / blond-blonde) via
    the `tagsMerge` mutation, preserving each merged name as an ALIAS on the
    surviving tag so future scraper matches still hit;
  - REPORTS the SUBJECTIVE families (POV / parenthetical / hair / age-bracket /
    1-edit near-dupes) for you to merge by hand (or via the assistant), since those
    are judgement calls.

SAFETY: defaults to a dry-run (no writes). Pass --apply to execute the SAFE merges.
Subjective clusters are NEVER auto-merged. Destination of a cluster = the member
with the most scenes (ties: longer name, then alphabetical; "blonde" preferred over
"blond").

Usage:
  python3 tag_consolidate.py --stash-url http://overwatch-stash:9999          # dry-run report
  python3 tag_consolidate.py --stash-url http://overwatch-stash:9999 --apply   # merge SAFE clusters
"""

import argparse
import json
import re
import sys
import urllib.error
import urllib.request

# ── GraphQL plumbing (mirrors tools/identify/external_identify.py) ──────────────


def gql(url, query, variables=None, headers=None, timeout=60):
    body = json.dumps({"query": query, "variables": variables or {}}).encode()
    req = urllib.request.Request(url, data=body, method="POST")
    req.add_header("Content-Type", "application/json")
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


class Stash:
    def __init__(self, url, api_key=None):
        self.url = url.rstrip("/") + "/graphql"
        self.headers = {"ApiKey": api_key} if api_key else {}

    def q(self, query, variables=None):
        return gql(self.url, query, variables, self.headers)


# ── clustering ──────────────────────────────────────────────────────────────

_BLOND = re.compile(r"\bblond\b")
_NONALNUM = re.compile(r"[^a-z0-9]+")
_WS = re.compile(r"\s+")


def safe_key(name):
    """Normalized key that collapses ONLY trivial formatting differences: case,
    punctuation/whitespace, en/em-dashes, blond->blonde, and a single trailing
    plural 's'. Crucially it KEEPS meaningful words, so "Blowjob" and
    "Blowjob - POV" map to different keys (POV is not a trivial variant)."""
    s = name.strip().lower().replace("–", "-").replace("—", "-").replace("’", "'")
    s = _NONALNUM.sub(" ", s).strip()
    s = _BLOND.sub("blonde", s)
    s = _WS.sub(" ", s)
    if len(s) > 3 and s.endswith("s") and not s.endswith("ss"):
        s = s[:-1]
    return s


def base_family(name):
    """Strip a trailing POV marker or parenthetical so e.g. 'Cowgirl', 'Cowgirl (DP)'
    and 'Cowgirl - POV' share a base — used to surface SUBJECTIVE families."""
    s = name.strip()
    s = re.sub(r"\s*[-(]\s*pov\s*\)?\s*$", "", s, flags=re.I)
    s = re.sub(r"\s*\([^)]*\)\s*$", "", s)
    s = re.sub(r"\s*-\s*pov\s*$", "", s, flags=re.I)
    return s.strip().lower()


def levenshtein1(a, b):
    """True if edit distance between a and b is exactly 1 (cheap early-out)."""
    if a == b:
        return False
    la, lb = len(a), len(b)
    if abs(la - lb) > 1:
        return False
    if la == lb:  # one substitution
        return sum(1 for x, y in zip(a, b) if x != y) == 1
    # one insertion/deletion
    if la > lb:
        a, b, la, lb = b, a, lb, la
    i = j = 0
    diff = 0
    while i < la and j < lb:
        if a[i] != b[j]:
            diff += 1
            if diff > 1:
                return False
            j += 1
        else:
            i += 1
            j += 1
    return True


def pick_destination(members):
    """Choose the surviving tag in a cluster: most scenes, then 'blonde' over
    'blond', then longest name, then alphabetical."""
    def rank(t):
        return (t["scene_count"], "blonde" in t["name"].lower(), len(t["name"]), )
    return max(members, key=lambda t: (rank(t), t["name"]))


def union_aliases(dest, sources):
    """Destination's existing aliases + every source name and alias, deduped
    case-insensitively, excluding the destination's own name."""
    seen = {}
    out = []
    for v in (dest.get("aliases") or []) + \
             [s["name"] for s in sources] + \
             [a for s in sources for a in (s.get("aliases") or [])]:
        v = v.strip()
        k = v.lower()
        if not v or k == dest["name"].strip().lower() or k in seen:
            continue
        seen[k] = True
        out.append(v)
    return out


# ── main ────────────────────────────────────────────────────────────────────

ALL_TAGS_Q = "{ findTags(filter:{per_page:-1}){ count tags { id name aliases scene_count } } }"
MERGE_M = ("mutation($i:TagsMergeInput!){ tagsMerge(input:$i){ id name } }")


def main():
    ap = argparse.ArgumentParser(description="Cluster + merge duplicate stash tags (SAFE auto, subjective report-only).")
    ap.add_argument("--stash-url", default="http://localhost:9999")
    ap.add_argument("--stash-api-key", default=None)
    ap.add_argument("--apply", action="store_true", help="execute the SAFE merges (default: dry-run report)")
    ap.add_argument("--show-subjective", type=int, default=40, help="max subjective groups to print")
    args = ap.parse_args()
    dry = not args.apply

    stash = Stash(args.stash_url, args.stash_api_key)
    tags = stash.q(ALL_TAGS_Q)["findTags"]["tags"]
    print(f"{len(tags)} tags loaded.")

    by_id = {t["id"]: t for t in tags}

    # ── SAFE clusters: group by safe_key ──
    safe_groups = {}
    for t in tags:
        safe_groups.setdefault(safe_key(t["name"]), []).append(t)
    safe_clusters = [g for g in safe_groups.values() if len(g) > 1]
    in_safe = {t["id"] for g in safe_clusters for t in g}

    print(f"\n=== SAFE clusters (auto-merge{'' if not dry else ' — DRY RUN'}): {len(safe_clusters)} ===")
    merged = failed = 0
    for g in sorted(safe_clusters, key=lambda g: -sum(t["scene_count"] for t in g)):
        dest = pick_destination(g)
        sources = [t for t in g if t["id"] != dest["id"]]
        names = ", ".join(f"{t['name']}({t['scene_count']})" for t in sources)
        print(f"  {names}  ->  {dest['name']}({dest['scene_count']}) [id {dest['id']}]")
        if dry:
            continue
        try:
            stash.q(MERGE_M, {"i": {
                "source": [s["id"] for s in sources],
                "destination": dest["id"],
                "values": {"id": dest["id"], "aliases": union_aliases(dest, sources)},
            }})
            merged += 1
        except RuntimeError as e:
            failed += 1
            print(f"    ! merge failed: {e}")

    # ── SUBJECTIVE families (report only) ──
    def report_group(title, groups):
        groups = [g for g in groups if len(g) > 1]
        if not groups:
            return
        print(f"\n=== {title} (review — NOT auto-merged): {len(groups)} groups ===")
        for g in sorted(groups, key=lambda g: -sum(t['scene_count'] for t in g))[: args.show_subjective]:
            print("  - " + " | ".join(f"{t['name']}({t['scene_count']})" for t in sorted(g, key=lambda t: -t['scene_count'])))

    fam = {}
    for t in tags:
        fam.setdefault(base_family(t["name"]), []).append(t)
    # only families that aren't already a single SAFE cluster
    fam_groups = [g for g in fam.values() if len(g) > 1 and not {t["id"] for t in g} <= in_safe]
    report_group("POV / parenthetical families", fam_groups)

    hair = [t for t in tags if t["name"].strip().lower().endswith("hair")]
    report_group("Hair-colour tags", [hair])

    age = [t for t in tags if re.search(r"\(\s*\d+\s*[-–+]", t["name"])]
    report_group("Age-bracket tags", [age])

    # 1-edit near-dupes not already grouped (bounded, informative)
    near = []
    keys = [(t, safe_key(t["name"])) for t in tags if t["id"] not in in_safe]
    seen_pairs = set()
    for i, (ta, ka) in enumerate(keys):
        if len(ka) < 6:
            continue
        for tb, kb in keys[i + 1:]:
            if len(kb) < 6 or ka[:3] != kb[:3]:
                continue
            if levenshtein1(ka, kb):
                pair = tuple(sorted((ta["id"], tb["id"])))
                if pair not in seen_pairs:
                    seen_pairs.add(pair)
                    near.append([ta, tb])
    report_group("1-edit near-duplicates", near)

    print(f"\nsummary: {len(safe_clusters)} safe clusters "
          f"({'would merge' if dry else f'merged {merged}, failed {failed}'}), "
          f"{len(fam_groups)} subjective families flagged. "
          f"{'(dry-run — re-run with --apply)' if dry else ''}")


if __name__ == "__main__":
    main()
