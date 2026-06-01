#!/usr/bin/env python3
"""
external_filename_parse.py — fill scene metadata from filenames/paths, LINK-ONLY.

The fingerprint identifier (external_identify.py) can't match the long tail of
scenes that aren't in any stash-box. Their studio/performers are, however, usually
right there in the path ("Studio - Performer A & Performer B - Title.mp4", or a
performer-named folder). This tool reads the path and links scenes to EXISTING
studios/performers — it NEVER creates new ones (filenames are noisy), so it can't
pollute your library with junk records.

Strategy (deliberately conservative): instead of guessing the filename's structure,
it matches the path against the set of names/aliases you ALREADY have. A known
performer or studio name appearing in the path (whole-word) is a link; nothing else
is touched. MERGE semantics — only fills a scene's EMPTY studio, adds performers on
top of existing, and (with --set-title) sets a cleaned title when missing. A clear
date in the path is filled when missing.

SAFETY: defaults to --dry-run. Lower fidelity than fingerprint identify — review the
dry-run before --apply. Do NOT run alongside a native Identify over the same scenes.

Usage:
  python3 external_filename_parse.py --stash-url http://overwatch-stash:9999             # dry-run
  python3 external_filename_parse.py --stash-url http://overwatch-stash:9999 --apply
"""

import argparse
import json
import re
import sys
import urllib.error
import urllib.request

# ── GraphQL plumbing ───────────────────────────────────────────────────────────


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


# ── name index (existing performers + studios) ─────────────────────────────────

# A name/alias is only "matchable" if it is MULTI-WORD (contains a space). Single
# tokens — even long ones like "Alexis" or "Nicole" — false-match constantly (a
# first-name alias hits any path mentioning that first name, linking the wrong
# performer). For a link-only tool a missed link is far cheaper than a wrong one,
# so we require a full "First Last"-style name to appear contiguously.
def matchable(name, min_len):
    n = name.strip()
    return " " in n and len(n) >= 5


def build_matchers(entries, min_len):
    """entries: [{id,name,aliases}]. Returns [(compiled_regex, id, display_name)]
    sorted longest-first so the most specific match wins."""
    out = []
    for e in entries:
        names = [e["name"]] + (e.get("aliases") or [])
        best = max((n for n in names), key=len, default="")
        for n in names:
            if not matchable(n, min_len):
                continue
            # whole-word, case-insensitive, tolerant of . _ - separators between words
            pat = r"(?<![a-z0-9])" + re.escape(n.strip()).replace(r"\ ", r"[\s._-]+") + r"(?![a-z0-9])"
            out.append((re.compile(pat, re.I), e["id"], best, len(n)))
    out.sort(key=lambda t: -t[3])
    return out


def normalize_path_text(path):
    """The filename + its immediate parent folder, as searchable text."""
    parts = re.split(r"[\\/]+", path)
    stem = re.sub(r"\.[a-z0-9]{2,4}$", "", parts[-1], flags=re.I) if parts else ""
    folder = parts[-2] if len(parts) >= 2 else ""
    return f"{folder} {stem}", stem, folder


DATE_RES = [
    (re.compile(r"(20\d{2})[.\-_](\d{2})[.\-_](\d{2})"), (1, 2, 3)),      # 2009-04-06
    (re.compile(r"(?<!\d)(\d{2})[.\-_](\d{2})[.\-_](\d{2})(?!\d)"), None),  # 23.04.06 (ambiguous → YY.MM.DD)
]


def parse_date(text):
    m = DATE_RES[0][0].search(text)
    if m:
        y, mo, d = m.group(1), m.group(2), m.group(3)
        if 1 <= int(mo) <= 12 and 1 <= int(d) <= 31:
            return f"{y}-{mo}-{d}"
    m = DATE_RES[1][0].search(text)
    if m:
        yy, mo, d = m.group(1), m.group(2), m.group(3)
        if 1 <= int(mo) <= 12 and 1 <= int(d) <= 31:
            return f"20{yy}-{mo}-{d}"
    return None


def clean_title(stem, studio_name, performer_names, date_str):
    t = stem
    for n in [studio_name] + list(performer_names):
        if n:
            t = re.sub(re.escape(n), " ", t, flags=re.I)
    t = re.sub(r"\b(19|20)\d{2}[.\-_]\d{2}[.\-_]\d{2}\b", " ", t)
    t = re.sub(r"\b\d{3,4}p\b|\b(4k|2160p|uhd|xxx|hd)\b", " ", t, flags=re.I)
    t = re.sub(r"\b[A-Z]{2,}\d{2,}\b", " ", t)            # scene codes GIO1240/AH297
    t = re.sub(r"[._\-&]+", " ", t)
    t = re.sub(r"\s+", " ", t).strip(" -_.")
    return t if len(t) >= 4 else ""


# ── apply ──────────────────────────────────────────────────────────────────────

SCENES_Q = """
query($filter: FindFilterType, $scene_filter: SceneFilterType) {
  findScenes(filter: $filter, scene_filter: $scene_filter) {
    count
    scenes {
      id title date
      studio { id }
      performers { id }
      files { path }
    }
  }
}
"""


def fetch_all(stash, query, key, per_page=500):
    out, page = [], 1
    while True:
        data = stash.q(query, {"filter": {"per_page": per_page, "page": page, "sort": "id", "direction": "ASC"}})
        items = data[key][list(data[key].keys())[-1]]
        out.extend(items)
        if len(items) < per_page:
            break
        page += 1
    return out


def main():
    ap = argparse.ArgumentParser(description="Link unidentified scenes to EXISTING studios/performers from their filenames.")
    ap.add_argument("--stash-url", default="http://localhost:9999")
    ap.add_argument("--stash-api-key", default=None)
    ap.add_argument("--apply", action="store_true", help="write changes (default: dry-run)")
    ap.add_argument("--all-scenes", action="store_true", help="consider all scenes (default: only scenes with no stash_id)")
    ap.add_argument("--set-title", action="store_true", help="also set a cleaned title when the scene has none (noisier)")
    ap.add_argument("--limit", type=int, default=0, help="cap scenes processed (0 = no cap)")
    ap.add_argument("--min-name-len", type=int, default=5, help="min length for a single-token name to be matchable")
    args = ap.parse_args()
    dry = not args.apply

    stash = Stash(args.stash_url, args.stash_api_key)

    print("loading existing performers + studios …")
    performers = fetch_all(stash, "query($filter:FindFilterType){findPerformers(filter:$filter){performers{id name aliases: alias_list}}}", "findPerformers")
    studios = fetch_all(stash, "query($filter:FindFilterType){findStudios(filter:$filter){studios{id name aliases}}}", "findStudios")
    perf_match = build_matchers(performers, args.min_name_len)
    studio_match = build_matchers(studios, args.min_name_len)
    print(f"  {len(performers)} performers, {len(studios)} studios indexed.")

    sf = None if args.all_scenes else {"is_missing": "stash_id"}
    page, scenes = 1, []
    while True:
        data = stash.q(SCENES_Q, {"filter": {"per_page": 200, "page": page, "sort": "path", "direction": "ASC"}, "scene_filter": sf})
        chunk = data["findScenes"]["scenes"]
        scenes.extend(chunk)
        if len(chunk) < 200 or (args.limit and len(scenes) >= args.limit):
            break
        page += 1
    if args.limit:
        scenes = scenes[: args.limit]
    print(f"{len(scenes)} candidate scene(s){' (dry-run)' if dry else ''}\n")

    linked_studio = linked_perf = set_title = set_date = applied = skipped = 0

    for s in scenes:
        path = (s.get("files") or [{}])[0].get("path", "")
        if not path:
            continue
        text, stem, folder = normalize_path_text(path)

        # studio: first (longest) known studio name appearing in the path
        studio_id = studio_name = None
        for rx, sid, disp, _ in studio_match:
            if rx.search(text):
                studio_id, studio_name = sid, disp
                break

        # performers: all known performer names appearing in the path
        perf_ids, perf_names = [], []
        for rx, pid, disp, _ in perf_match:
            if pid in perf_ids:
                continue
            if rx.search(text):
                perf_ids.append(pid)
                perf_names.append(disp)

        date_str = parse_date(text)
        title = clean_title(stem, studio_name, perf_names, date_str) if args.set_title else ""

        upd = {"id": s["id"]}
        if studio_id and not s.get("studio"):
            upd["studio_id"] = studio_id
        if perf_ids:
            existing = [p["id"] for p in (s.get("performers") or [])]
            merged = existing + [p for p in perf_ids if p not in existing]
            if merged != existing:
                upd["performer_ids"] = merged
        if date_str and not s.get("date"):
            upd["date"] = date_str
        if title and not s.get("title"):
            upd["title"] = title

        if len(upd) == 1:  # nothing to set
            skipped += 1
            continue

        desc = []
        if "studio_id" in upd:
            desc.append(f"studio={studio_name}")
            linked_studio += 1
        if "performer_ids" in upd:
            desc.append(f"performers=[{', '.join(perf_names)}]")
            linked_perf += 1
        if "date" in upd:
            desc.append(f"date={date_str}")
            set_date += 1
        if "title" in upd:
            desc.append(f"title=\"{title}\"")
            set_title += 1

        if dry:
            print(f"  scene {s['id']}: {'; '.join(desc)}   <- {stem[:70]}")
        else:
            try:
                stash.q("mutation($i:SceneUpdateInput!){sceneUpdate(input:$i){id}}", {"i": upd})
                applied += 1
                print(f"  [applied] scene {s['id']}: {'; '.join(desc)}")
            except RuntimeError as e:
                print(f"  ! scene {s['id']} update failed: {e}")

    print(f"\nsummary: studio-links={linked_studio} performer-links={linked_perf} "
          f"dates={set_date} titles={set_title} skipped(no-match)={skipped} "
          f"{'applied=' + str(applied) if not dry else '(dry-run — re-run with --apply)'}")


if __name__ == "__main__":
    main()
