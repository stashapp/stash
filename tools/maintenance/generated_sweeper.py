#!/usr/bin/env python3
"""
generated_sweeper.py — audit (and optionally prune) stash's generated/ directory.

stash writes derived artifacts (previews, sprites, transcodes, thumbnails, marker
images, heatmaps) under a generated/ root, keyed by each file's hash (oshash or
md5). When scenes/images are deleted or re-hashed, their artifacts are left behind
as ORPHANS. This tool:

  - builds the set of LIVE file hashes from the stash API (scene + image
    fingerprints),
  - walks generated/ and flags files whose hash is no longer live (orphans),
    zero-byte/corrupt artifacts, and stale tmp/ leftovers,
  - prints a coverage summary (how many scenes have a preview / sprite / transcode),
  - and, with --prune, MOVES orphans into generated/_quarantine/ (reversible — never
    hard-deletes).

Run it ON the NAS for a fast local walk (`--generated-dir /volume1/docker/stash/generated`)
or remotely against an SMB path. Read-only by default.

Usage:
  python3 generated_sweeper.py --stash-url http://overwatch-stash:9999 \
      --generated-dir /volume1/docker/stash/generated
  # add --prune to quarantine the orphans it finds
"""

import argparse
import json
import os
import shutil
import sys
import time
import urllib.error
import urllib.request

# ── GraphQL ──────────────────────────────────────────────────────────────────


def gql(url, query, variables=None, headers=None, timeout=120, retries=3):
    body = json.dumps({"query": query, "variables": variables or {}}).encode()
    last = None
    for attempt in range(retries):
        req = urllib.request.Request(url, data=body, method="POST")
        req.add_header("Content-Type", "application/json")
        req.add_header("User-Agent",
                       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "
                       "(KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
        for k, v in (headers or {}).items():
            req.add_header(k, v)
        try:
            with urllib.request.urlopen(req, timeout=timeout) as resp:
                raw = resp.read().decode()
            if not raw.strip():
                raise ValueError("empty response body")
            payload = json.loads(raw)
        except (urllib.error.URLError, ValueError) as e:
            # transient under load (stash busy with other jobs) — back off + retry
            last = e
            time.sleep(1.5 * (attempt + 1))
            continue
        if payload.get("errors"):
            raise RuntimeError(f"GraphQL errors: {json.dumps(payload['errors'])[:400]}")
        return payload.get("data") or {}
    raise RuntimeError(f"request to {url} failed after {retries} tries: {last}")


def fetch_live_hashes(url, headers):
    """All scene + image file hashes (oshash + md5) — the set of valid artifact keys."""
    hashes = set()
    scene_hashes = set()
    for kind, q in (
        ("findScenes", "query($f:FindFilterType){findScenes(filter:$f){scenes{files{fingerprints{type value}}}}}"),
        ("findImages", "query($f:FindFilterType){findImages(filter:$f){images{visual_files{... on ImageFile{fingerprints{type value}} ... on VideoFile{fingerprints{type value}}}}}}"),
    ):
        page = 1
        while True:
            data = gql(url, q, {"f": {"per_page": 1000, "page": page, "sort": "id", "direction": "ASC"}}, headers)
            items = data[kind][list(data[kind].keys())[0]]
            if not items:
                break
            for it in items:
                files = it.get("files") or it.get("visual_files") or []
                for f in files:
                    for fp in (f.get("fingerprints") or []):
                        t, v = (fp.get("type") or "").lower(), fp.get("value")
                        if v and t in ("oshash", "md5"):
                            hashes.add(v)
                            if kind == "findScenes":
                                scene_hashes.add(v)
            if len(items) < 1000:
                break
            page += 1
    return hashes, scene_hashes


# ── generated/ layout → hash extraction (mirrors pkg/models/paths) ─────────────

VIDEO_EXTS = (".mp4", ".webp", ".webm", ".jpg", ".png", ".vtt")


def hash_for(subdir, fname):
    """Return the hash a generated file is keyed on, or None if unrecognized."""
    stem = fname
    if subdir == "screenshots":            # {hash}.mp4 | {hash}.webp
        for e in (".mp4", ".webp"):
            if fname.endswith(e):
                return fname[: -len(e)]
    elif subdir == "vtt":                  # {hash}_sprite.jpg | {hash}_thumbs.vtt
        for suf in ("_sprite.jpg", "_thumbs.vtt"):
            if fname.endswith(suf):
                return fname[: -len(suf)]
    elif subdir == "transcodes":           # {hash}.mp4
        if fname.endswith(".mp4"):
            return fname[:-4]
    elif subdir == "interactive_heatmaps":  # {hash}.png
        if fname.endswith(".png"):
            return fname[:-4]
    elif subdir == "thumbnails":           # {hash[:2]}/{hash}_{width}.jpg|.webm
        for e in (".jpg", ".webm"):
            if fname.endswith(e) and "_" in fname:
                return fname[: -len(e)].rsplit("_", 1)[0]
    return None


# subdirs whose IMMEDIATE child is a per-hash directory (markers/{hash}/{secs}.*)
DIR_KEYED = {"markers"}
# subdirs we walk for files keyed by filename
FILE_KEYED = {"screenshots", "vtt", "transcodes", "interactive_heatmaps", "thumbnails"}


def main():
    ap = argparse.ArgumentParser(description="Audit/prune stash's generated/ directory for orphans + integrity.")
    ap.add_argument("--stash-url", default="http://localhost:9999")
    ap.add_argument("--stash-api-key", default=None)
    ap.add_argument("--generated-dir", required=True, help="filesystem path to stash's generated/ root")
    ap.add_argument("--prune", action="store_true", help="move orphans into generated/_quarantine/ (default: report only)")
    args = ap.parse_args()
    headers = {"ApiKey": args.stash_api_key} if args.stash_api_key else {}
    gql_url = args.stash_url.rstrip("/") + "/graphql"
    root = args.generated_dir.rstrip("/\\")

    if not os.path.isdir(root):
        sys.exit(f"generated dir not found: {root}")

    print("fetching live file hashes from stash …")
    live, scene_hashes = fetch_live_hashes(gql_url, headers)
    print(f"  {len(live)} live hashes ({len(scene_hashes)} scene).")

    orphans, zero_byte, present = [], [], {}
    quarantine = os.path.join(root, "_quarantine")

    def check(subdir, path, h):
        present[subdir] = present.get(subdir, 0) + 1
        try:
            if os.path.getsize(path) == 0:
                zero_byte.append(path)
        except OSError:
            pass
        if h is not None and h not in live:
            orphans.append((subdir, path))

    for subdir in sorted(FILE_KEYED | DIR_KEYED):
        d = os.path.join(root, subdir)
        if not os.path.isdir(d):
            continue
        if subdir in DIR_KEYED:                 # markers/{hash}/...
            for entry in os.scandir(d):
                if entry.is_dir():
                    h = entry.name
                    for f in os.scandir(entry.path):
                        if f.is_file():
                            check(subdir, f.path, h)
        elif subdir == "thumbnails":            # nested one level
            for sub in os.scandir(d):
                if sub.is_dir():
                    for f in os.scandir(sub.path):
                        if f.is_file():
                            check(subdir, f.path, hash_for(subdir, f.name))
        else:
            for f in os.scandir(d):
                if f.is_file():
                    check(subdir, f.path, hash_for(subdir, f.name))

    # stale tmp / download_stage
    stale_tmp = []
    for t in ("tmp", "download_stage"):
        td = os.path.join(root, t)
        if os.path.isdir(td):
            for dirpath, _, files in os.walk(td):
                for f in files:
                    stale_tmp.append(os.path.join(dirpath, f))

    # ── report ──
    print("\n=== coverage (generated files present per type) ===")
    for k in sorted(present):
        print(f"  {k:22} {present[k]}")
    print(f"  (scenes in library: {len(scene_hashes)})")

    print(f"\n=== orphans (hash no longer in library): {len(orphans)} ===")
    for subdir, p in orphans[:50]:
        print(f"  [{subdir}] {os.path.relpath(p, root)}")
    if len(orphans) > 50:
        print(f"  … and {len(orphans) - 50} more")

    print(f"\n=== zero-byte / corrupt artifacts: {len(zero_byte)} ===")
    for p in zero_byte[:30]:
        print(f"  {os.path.relpath(p, root)}")

    print(f"\n=== stale tmp/download_stage files: {len(stale_tmp)} ===")
    for p in stale_tmp[:20]:
        print(f"  {os.path.relpath(p, root)}")

    if args.prune and orphans:
        print(f"\n=== pruning {len(orphans)} orphans -> {quarantine}/ ===")
        moved = 0
        for subdir, p in orphans:
            dest_dir = os.path.join(quarantine, subdir)
            os.makedirs(dest_dir, exist_ok=True)
            try:
                shutil.move(p, os.path.join(dest_dir, os.path.basename(p)))
                moved += 1
            except OSError as e:
                print(f"  ! could not move {p}: {e}")
        print(f"  moved {moved} files into {quarantine}/ (reversible).")
    elif orphans:
        print("\n(orphans listed only — re-run with --prune to quarantine them)")

    print(f"\nsummary: {sum(present.values())} artifacts, {len(orphans)} orphans, "
          f"{len(zero_byte)} zero-byte, {len(stale_tmp)} stale tmp"
          f"{' — pruned' if (args.prune and orphans) else ''}")


if __name__ == "__main__":
    main()
