#!/usr/bin/env python3
"""
scene_dedup.py — find near-duplicate scenes in stash using perceptual hashes (phash).

stash computes a perceptual hash (phash) per video file during phash Generate. Two
re-encodes / re-rips / resolution variants of the same scene end up with phashes that
are very close in Hamming distance even though their oshash/md5 differ completely. This
tool fetches every scene's primary-file phash, clusters scenes whose phashes are within
`--max-distance` bits of each other (default 8), picks a "keeper" per cluster (highest
resolution, else largest file), and reports the duplicate candidates + reclaimable space.

It does NOT delete anything. With --apply it merely TAGS the non-keeper scenes with a
`_dupe-candidate` tag (find-or-create) so you can review/cull them in the stash UI. The
tag is MERGED onto each scene's existing tags — never replacing them.

SAFETY: defaults to a dry-run (no writes). Pass --apply to tag non-keepers.
Clustering is single-linkage (union-find) over phash Hamming distance, bucketed by a
short hex prefix so we only compare scenes that share leading bits — a few thousand
scenes runs in well under a second on stdlib Python.

Usage:
  python3 scene_dedup.py --stash-url http://overwatch-stash:9999                  # dry-run report
  python3 scene_dedup.py --stash-url http://overwatch-stash:9999 --max-distance 6 # tighter
  python3 scene_dedup.py --stash-url http://overwatch-stash:9999 --apply          # tag non-keepers
"""

import argparse
import json
import sys
import urllib.error
import urllib.request

# ── GraphQL plumbing (mirrors tools/identify/external_identify.py) ──────────────


def gql(url, query, variables=None, headers=None, timeout=60):
    body = json.dumps({"query": query, "variables": variables or {}}).encode()
    req = urllib.request.Request(url, data=body, method="POST")
    req.add_header("Content-Type", "application/json")
    # stash returns an EMPTY body to the default "Python-urllib/x.y" User-Agent, so
    # present a browser UA. (Same gotcha as the other tools in this repo.)
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


# ── stash queries ──────────────────────────────────────────────────────────────

# Validated against stash (per_page 2): id/title resolve, and each file exposes
# id/path/size (bytes)/width/height/fingerprints{type,value}. The phash fingerprint
# value is a hex string of a uint64, e.g. "d662f6db041bce24".
SCENES_Q = """
query($filter: FindFilterType) {
  findScenes(filter: $filter) {
    count
    scenes {
      id
      title
      files {
        id
        path
        size
        width
        height
        fingerprints { type value }
      }
    }
  }
}
"""

TAG_FIND_Q = "query($f:TagFilterType){findTags(tag_filter:$f,filter:{per_page:1}){tags{id}}}"
TAG_CREATE_M = "mutation($i:TagCreateInput!){tagCreate(input:$i){id name}}"
SCENE_UPDATE_M = "mutation($i:SceneUpdateInput!){sceneUpdate(input:$i){id}}"
# Lightweight per-scene read so --apply merges (never replaces) existing tags.
SCENE_TAGS_Q = "query($id:ID!){findScene(id:$id){id tags{id}}}"

DUPE_TAG_NAME = "_dupe-candidate"


def fetch_scenes(stash, per_page=1000, limit=0):
    """Paginate findScenes. Stops early once `limit` scenes are collected (0 = all)."""
    page, out = 1, []
    while True:
        data = stash.q(SCENES_Q, {
            "filter": {"per_page": per_page, "page": page, "sort": "id", "direction": "ASC"},
        })
        block = data["findScenes"]["scenes"]
        out.extend(block)
        if limit and len(out) >= limit:
            return out[:limit]
        if len(block) < per_page:
            break
        page += 1
    return out


# ── phash helpers ──────────────────────────────────────────────────────────────


def primary_phash(scene):
    """Return the hex phash string of the scene's primary file (files[0]), or None."""
    files = scene.get("files") or []
    if not files:
        return None
    for fp in files[0].get("fingerprints") or []:
        if (fp.get("type") or "").lower() == "phash":
            v = fp.get("value")
            if v:
                return v
    return None


def hamming(a_int, b_int):
    """Hamming distance between two uint64s already parsed from hex."""
    return bin(a_int ^ b_int).count("1")


def resolution(scene):
    """width*height of the primary file, or 0 if unknown."""
    files = scene.get("files") or []
    if not files:
        return 0
    f = files[0]
    w, h = f.get("width") or 0, f.get("height") or 0
    return (w or 0) * (h or 0)


def file_size(scene):
    """Primary-file size in bytes (size can arrive as str from GraphQL Int64)."""
    files = scene.get("files") or []
    if not files:
        return 0
    try:
        return int(files[0].get("size") or 0)
    except (TypeError, ValueError):
        return 0


def primary_path(scene):
    files = scene.get("files") or []
    return files[0].get("path") if files else None


def human_size(n):
    n = float(n)
    for unit in ("B", "KiB", "MiB", "GiB", "TiB"):
        if n < 1024 or unit == "TiB":
            return f"{n:.2f} {unit}"
        n /= 1024


# ── clustering (single-linkage union-find over phash Hamming distance) ──────────


class UnionFind:
    def __init__(self, n):
        self.parent = list(range(n))

    def find(self, x):
        root = x
        while self.parent[root] != root:
            root = self.parent[root]
        while self.parent[x] != root:  # path compression
            self.parent[x], x = root, self.parent[x]
        return root

    def union(self, a, b):
        ra, rb = self.find(a), self.find(b)
        if ra != rb:
            self.parent[rb] = ra


def cluster_by_phash(items, max_distance, prefix_bits=12):
    """items: list of (scene, phash_int). Returns list of index-lists (clusters of >1).

    To avoid a full O(n^2) compare across thousands of scenes, bucket by the top
    `prefix_bits` of the phash. Two phashes within `max_distance` bits can differ in at
    most `max_distance` positions, so they will share a top-bits bucket UNLESS one of
    the differing bits falls in the prefix. To stay correct we therefore compare each
    item against every bucket reachable by flipping up to `max_distance` of its prefix
    bits — but that explodes combinatorially, so instead we use a cheaper, exact scheme:
    bucket by prefix, compare fully WITHIN each bucket, then also do one global pass that
    only compares items whose prefixes themselves are within `max_distance` Hamming. For
    a few thousand scenes the global pass is a small n^2 over distinct prefixes, which is
    trivial. This keeps results exact while skipping the bulk of pairwise work."""
    n = len(items)
    uf = UnionFind(n)
    shift = 64 - prefix_bits

    # bucket index lists by prefix
    buckets = {}
    for idx, (_, ph) in enumerate(items):
        buckets.setdefault(ph >> shift, []).append(idx)

    # 1) exact compare within each prefix bucket
    for members in buckets.values():
        for i in range(len(members)):
            ai = members[i]
            pa = items[ai][1]
            for j in range(i + 1, len(members)):
                bi = members[j]
                if hamming(pa, items[bi][1]) <= max_distance:
                    uf.union(ai, bi)

    # 2) cross-bucket: only compare pairs of buckets whose prefixes are themselves
    #    within max_distance bits (a near-prefix can still hide a near-full-hash).
    prefixes = list(buckets.keys())
    for i in range(len(prefixes)):
        pi = prefixes[i]
        for j in range(i + 1, len(prefixes)):
            pj = prefixes[j]
            if bin(pi ^ pj).count("1") > max_distance:
                continue
            for ai in buckets[pi]:
                pa = items[ai][1]
                for bi in buckets[pj]:
                    if hamming(pa, items[bi][1]) <= max_distance:
                        uf.union(ai, bi)

    groups = {}
    for idx in range(n):
        groups.setdefault(uf.find(idx), []).append(idx)
    return [g for g in groups.values() if len(g) > 1]


def pick_keeper(cluster_scenes):
    """Keeper = highest resolution (w*h); tie-break on largest file size, then lowest id."""
    return max(
        cluster_scenes,
        key=lambda s: (resolution(s), file_size(s), -int(s["id"])),
    )


# ── apply helpers ──────────────────────────────────────────────────────────────


def ensure_dupe_tag(stash, dry):
    """Find the _dupe-candidate tag by name, creating it if missing. Returns its id
    (or a placeholder string in dry-run when it doesn't exist yet)."""
    found = stash.q(TAG_FIND_Q, {"f": {"name": {"value": DUPE_TAG_NAME, "modifier": "EQUALS"}}})
    tags = found["findTags"]["tags"]
    if tags:
        return tags[0]["id"]
    if dry:
        return f"(would-create:{DUPE_TAG_NAME})"
    return stash.q(TAG_CREATE_M, {"i": {"name": DUPE_TAG_NAME}})["tagCreate"]["id"]


def tag_scene(stash, scene_id, dupe_tag_id):
    """MERGE the dupe tag onto a scene's existing tags (re-read tags first so we never
    clobber). No-op if the scene already carries the tag."""
    cur = stash.q(SCENE_TAGS_Q, {"id": scene_id})["findScene"]
    existing = [t["id"] for t in (cur.get("tags") or [])]
    if dupe_tag_id in existing:
        return False
    stash.q(SCENE_UPDATE_M, {"i": {"id": scene_id, "tag_ids": existing + [dupe_tag_id]}})
    return True


# ── main ────────────────────────────────────────────────────────────────────


def main():
    ap = argparse.ArgumentParser(description="Find near-duplicate scenes by phash; optionally tag non-keepers.")
    ap.add_argument("--stash-url", default="http://localhost:9999")
    ap.add_argument("--stash-api-key", default=None)
    ap.add_argument("--max-distance", type=int, default=8,
                    help="max phash Hamming distance to treat scenes as near-duplicates (default 8)")
    ap.add_argument("--apply", action="store_true",
                    help=f"tag non-keepers with '{DUPE_TAG_NAME}' (default: dry-run report only)")
    ap.add_argument("--limit", type=int, default=0,
                    help="cap scenes considered (0 = all)")
    args = ap.parse_args()
    dry = not args.apply

    stash = Stash(args.stash_url, args.stash_api_key)

    scenes = fetch_scenes(stash, per_page=1000, limit=args.limit)
    print(f"{len(scenes)} scene(s) fetched.")

    # keep only scenes whose PRIMARY file has a phash
    items = []  # (scene, phash_int)
    no_phash = 0
    for s in scenes:
        ph = primary_phash(s)
        if not ph:
            no_phash += 1
            continue
        try:
            items.append((s, int(ph, 16)))
        except ValueError:
            no_phash += 1
    print(f"{len(items)} with a primary-file phash, {no_phash} skipped (no phash / unparseable).")

    if not items:
        sys.exit("No phashed scenes to compare. (Has phash Generate run?)")

    clusters = cluster_by_phash(items, args.max_distance)
    clusters.sort(key=lambda g: -len(g))

    reclaimable = 0
    dup_count = 0
    to_tag = []  # (scene_id, title)

    print(f"\n=== near-duplicate clusters (max Hamming distance {args.max_distance}): {len(clusters)} ===")
    for gi, idxs in enumerate(clusters, 1):
        members = [items[i][0] for i in idxs]
        keeper = pick_keeper(members)
        print(f"\ncluster {gi} ({len(members)} scenes):")
        for s in sorted(members, key=lambda s: (-resolution(s), -file_size(s))):
            res = resolution(s)
            res_str = f"{(s.get('files') or [{}])[0].get('width')}x{(s.get('files') or [{}])[0].get('height')}" if res else "?x?"
            mark = "KEEP" if s["id"] == keeper["id"] else "dupe"
            print(f"  [{mark}] scene {s['id']}  {res_str}  {human_size(file_size(s))}  {primary_path(s)}")
            if s["id"] != keeper["id"]:
                reclaimable += file_size(s)
                dup_count += 1
                to_tag.append((s["id"], s.get("title") or "(untitled)"))

    print(f"\nfound {dup_count} duplicate candidate(s) across {len(clusters)} cluster(s); "
          f"potential reclaimable space: {human_size(reclaimable)}")

    if not to_tag:
        print("\nNothing to tag.")
        return

    if dry:
        print(f"\n(dry-run) would tag {len(to_tag)} non-keeper scene(s) with '{DUPE_TAG_NAME}'. "
              f"Re-run with --apply to write.")
        return

    dupe_tag_id = ensure_dupe_tag(stash, dry=False)
    print(f"\ntagging {len(to_tag)} non-keeper scene(s) with '{DUPE_TAG_NAME}' (id {dupe_tag_id}) ...")
    tagged = already = 0
    for sid, title in to_tag:
        try:
            if tag_scene(stash, sid, dupe_tag_id):
                tagged += 1
            else:
                already += 1
        except RuntimeError as e:
            print(f"  ! scene {sid} ({title}) failed: {e}")
    print(f"\nsummary: tagged {tagged}, already-tagged {already}, "
          f"reclaimable {human_size(reclaimable)} across {len(clusters)} clusters.")


if __name__ == "__main__":
    main()
