#!/usr/bin/env python3
"""
external_identify.py — identify stash scenes against stash-box endpoints (StashDB,
ThePornDB, …) WITHOUT using stash's internal job queue.

Why this exists: stash runs Tasks (Scan/Generate/Identify) sequentially, so a native
Identify queues behind a long phash Generate. This tool talks directly to the stash
GraphQL API and the stash-box GraphQL API, so it runs in parallel with whatever stash
is doing — and can run on a different machine entirely.

It matches by FINGERPRINT (oshash + phash), exactly like stash's native Identify, via
the stash-box `findScenesBySceneFingerprints` query. oshash is computed at scan time
(no phash needed), so this works immediately.

SAFETY: defaults to --dry-run (no writes). Field strategy is MERGE-like: it fills
title/date/details/studio/performers when the scene doesn't already have them, always
stamps the stash-box stash_id, and (optionally) marks the scene organized. TAGS are added
by default using a MERGE strategy — the match's tags are added on top of the scene's
existing tags (never replacing them), creating any tag that doesn't exist yet (stamped with
its stash-box id). Pass --no-tags to disable. Tag/performer/studio linking is name-based
(no alias resolution), so it's slightly lower fidelity than stash's native Identify — use
native Identify when you can; use this when you need it to run in parallel/off-box.

Do NOT run this at the same time as a native Identify job over the same library.

Usage:
  python3 external_identify.py --stash-url http://localhost:9999 --dry-run
  python3 external_identify.py --stash-url http://localhost:9999 --apply --set-organized
  # over Tailscale / from another machine, with a stash API key if auth is enabled:
  python3 external_identify.py --stash-url http://overwatch-stash:9999 --stash-api-key XXX --apply
"""

import argparse
import json
import sys
import time
import urllib.request
import urllib.error

# ── GraphQL plumbing ──────────────────────────────────────────────────────────


def gql(url, query, variables=None, headers=None, timeout=60):
    body = json.dumps({"query": query, "variables": variables or {}}).encode()
    req = urllib.request.Request(url, data=body, method="POST")
    req.add_header("Content-Type", "application/json")
    # ThePornDB sits behind Cloudflare, which 403s (error 1010) the default
    # "Python-urllib/x.y" User-Agent. Present a browser UA so the stdlib client
    # isn't bot-blocked. StashDB doesn't care.
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
        raise RuntimeError(f"GraphQL errors from {url}: {json.dumps(payload['errors'])[:400]}")
    return payload.get("data") or {}


class Stash:
    def __init__(self, url, api_key=None):
        self.url = url.rstrip("/") + "/graphql"
        self.headers = {"ApiKey": api_key} if api_key else {}

    def q(self, query, variables=None):
        return gql(self.url, query, variables, self.headers)


class StashBox:
    def __init__(self, endpoint, api_key, name, max_rpm=0):
        self.endpoint = endpoint
        self.name = name
        self.headers = {"ApiKey": api_key} if api_key else {}
        # min seconds between requests (0 max_rpm => gentle default of ~4/s)
        self.min_interval = (60.0 / max_rpm) if max_rpm and max_rpm > 0 else 0.25
        self._last = 0.0

    def q(self, query, variables=None):
        wait = self.min_interval - (time.monotonic() - self._last)
        if wait > 0:
            time.sleep(wait)
        self._last = time.monotonic()
        return gql(self.endpoint, query, variables, self.headers)


# ── stash-box query (mirrors stash's graphql/stash-box/query.graphql) ──────────

SB_FP_QUERY = """
query($fingerprints: [[FingerprintQueryInput!]!]!) {
  findScenesBySceneFingerprints(fingerprints: $fingerprints) {
    id
    title
    details
    date
    urls { url }
    studio { id name }
    performers { as performer { id name } }
    tags { id name }
  }
}
"""


# ── stash queries ──────────────────────────────────────────────────────────────

STASH_BOXES_Q = "{ configuration { general { stashBoxes { endpoint api_key name max_requests_per_minute } } } }"

SCENES_Q = """
query($filter: FindFilterType, $scene_filter: SceneFilterType) {
  findScenes(filter: $filter, scene_filter: $scene_filter) {
    count
    scenes {
      id title date details organized
      studio { id }
      performers { id }
      tags { id }
      stash_ids { endpoint stash_id }
      files { fingerprints { type value } }
    }
  }
}
"""


def fetch_stash_boxes(stash):
    data = stash.q(STASH_BOXES_Q)
    return data["configuration"]["general"]["stashBoxes"] or []


def fetch_scenes(stash, only_unorganized, phashed_only=False, per_page=100):
    """Fetch candidate scenes. Filters by organized flag and (optionally) by phash
    presence so callers can target the high-yield (phashed) scenes first while stash
    is still generating phashes for the rest."""
    sf = {}
    if only_unorganized:
        sf["organized"] = False
    if phashed_only:
        sf["phash"] = {"modifier": "NOT_NULL", "value": ""}
    scene_filter = sf or None
    page, out = 1, []
    while True:
        data = stash.q(SCENES_Q, {
            "filter": {"per_page": per_page, "page": page, "sort": "id", "direction": "ASC"},
            "scene_filter": scene_filter,
        })
        scenes = data["findScenes"]["scenes"]
        out.extend(scenes)
        if len(scenes) < per_page:
            break
        page += 1
    return out


def scene_fingerprints(scene):
    """Return [{hash, algorithm}] for a stash scene (oshash + phash)."""
    algo = {"oshash": "OSHASH", "phash": "PHASH", "md5": "MD5"}
    fps, seen = [], set()
    for f in scene.get("files") or []:
        for fp in f.get("fingerprints") or []:
            a = algo.get((fp.get("type") or "").lower())
            v = fp.get("value")
            if a and v and (a, v) not in seen:
                seen.add((a, v))
                fps.append({"hash": v, "algorithm": a})
    return fps


# ── apply helpers (find-or-create studio/performer, then sceneUpdate) ──────────

def ensure_studio(stash, sb_endpoint, sb_studio, cache, dry):
    name = sb_studio["name"]
    if name in cache:
        return cache[name]
    found = stash.q(
        "query($f:StudioFilterType){findStudios(studio_filter:$f,filter:{per_page:1}){studios{id}}}",
        {"f": {"name": {"value": name, "modifier": "EQUALS"}}},
    )["findStudios"]["studios"]
    if found:
        cache[name] = found[0]["id"]
        return cache[name]
    if dry:
        cache[name] = f"(would-create:{name})"
        return cache[name]
    sid = stash.q(
        "mutation($i:StudioCreateInput!){studioCreate(input:$i){id}}",
        {"i": {"name": name, "stash_ids": [{"endpoint": sb_endpoint, "stash_id": sb_studio["id"]}]}},
    )["studioCreate"]["id"]
    cache[name] = sid
    return sid


def ensure_performer(stash, sb_endpoint, sb_perf, cache, dry):
    name = sb_perf["name"]
    if name in cache:
        return cache[name]
    found = stash.q(
        "query($f:PerformerFilterType){findPerformers(performer_filter:$f,filter:{per_page:1}){performers{id}}}",
        {"f": {"name": {"value": name, "modifier": "EQUALS"}}},
    )["findPerformers"]["performers"]
    if found:
        cache[name] = found[0]["id"]
        return cache[name]
    if dry:
        cache[name] = f"(would-create:{name})"
        return cache[name]
    pid = stash.q(
        "mutation($i:PerformerCreateInput!){performerCreate(input:$i){id}}",
        {"i": {"name": name, "stash_ids": [{"endpoint": sb_endpoint, "stash_id": sb_perf["id"]}]}},
    )["performerCreate"]["id"]
    cache[name] = pid
    return pid


def ensure_tag(stash, sb_endpoint, sb_tag, cache, dry):
    """Find a tag by name, creating it (stamped with its stash-box id) if missing."""
    name = sb_tag["name"]
    if name in cache:
        return cache[name]
    found = stash.q(
        "query($f:TagFilterType){findTags(tag_filter:$f,filter:{per_page:1}){tags{id}}}",
        {"f": {"name": {"value": name, "modifier": "EQUALS"}}},
    )["findTags"]["tags"]
    if found:
        cache[name] = found[0]["id"]
        return cache[name]
    if dry:
        cache[name] = f"(would-create:{name})"
        return cache[name]
    stash_ids = [{"endpoint": sb_endpoint, "stash_id": sb_tag["id"]}] if sb_tag.get("id") else []
    tid = stash.q(
        "mutation($i:TagCreateInput!){tagCreate(input:$i){id}}",
        {"i": {"name": name, "stash_ids": stash_ids}},
    )["tagCreate"]["id"]
    cache[name] = tid
    return tid


def build_update(stash, scene, match, sb_endpoint, set_organized,
                 studio_cache, perf_cache, tag_cache, add_tags, dry):
    """MERGE-like SceneUpdateInput: fill empty single-value fields, merge tags, stamp stash_id."""
    upd = {"id": scene["id"]}
    if not scene.get("title") and match.get("title"):
        upd["title"] = match["title"]
    if not scene.get("date") and match.get("date"):
        upd["date"] = match["date"]
    if not scene.get("details") and match.get("details"):
        upd["details"] = match["details"]
    if not scene.get("studio") and match.get("studio"):
        upd["studio_id"] = ensure_studio(stash, sb_endpoint, match["studio"], studio_cache, dry)
    if not scene.get("performers") and match.get("performers"):
        ids = []
        for ap in match["performers"]:
            p = ap.get("performer")
            if p:
                ids.append(ensure_performer(stash, sb_endpoint, p, perf_cache, dry))
        if ids:
            upd["performer_ids"] = ids
    # tags — MERGE strategy: add the match's tags on top of the scene's existing tags
    # (never replacing them), creating any tag that doesn't exist yet. Default on.
    if add_tags and match.get("tags"):
        existing_tag_ids = [t["id"] for t in scene.get("tags") or []]
        tag_ids = list(existing_tag_ids)
        for tg in match["tags"]:
            if not tg.get("name"):
                continue
            tid = ensure_tag(stash, sb_endpoint, tg, tag_cache, dry)
            if tid not in tag_ids:
                tag_ids.append(tid)
        if tag_ids != existing_tag_ids:
            upd["tag_ids"] = tag_ids
    # stamp stash-box id if not already linked
    existing = {(s["endpoint"], s["stash_id"]) for s in scene.get("stash_ids") or []}
    new_sid = (sb_endpoint, match["id"])
    if new_sid not in existing:
        upd["stash_ids"] = [{"endpoint": e, "stash_id": i} for (e, i) in existing] + \
                           [{"endpoint": sb_endpoint, "stash_id": match["id"]}]
    if set_organized and not scene.get("organized"):
        upd["organized"] = True
    return upd


# ── main ────────────────────────────────────────────────────────────────────

def main():
    ap = argparse.ArgumentParser(description="External fingerprint identify for stash (parallel/off-box).")
    ap.add_argument("--stash-url", default="http://localhost:9999")
    ap.add_argument("--stash-api-key", default=None)
    ap.add_argument("--apply", action="store_true", help="actually write changes (default is dry-run)")
    ap.add_argument("--all-scenes", action="store_true", help="consider all scenes (default: only unorganized)")
    ap.add_argument("--phashed-only", action="store_true", help="only consider scenes that have a phash (best match yield while phash gen is in progress)")
    ap.add_argument("--set-organized", action="store_true", help="mark matched scenes organized")
    ap.add_argument("--no-tags", action="store_true", help="do not add tags from the match (tags are added by default, merged with existing)")
    ap.add_argument("--allow-multiple", action="store_true", help="apply first match when a scene has several")
    ap.add_argument("--batch", type=int, default=40, help="scenes per stash-box fingerprint query")
    ap.add_argument("--limit", type=int, default=0, help="cap scenes processed (0 = no cap)")
    args = ap.parse_args()
    dry = not args.apply
    add_tags = not args.no_tags

    stash = Stash(args.stash_url, args.stash_api_key)
    boxes = [StashBox(b["endpoint"], b["api_key"], b["name"], b.get("max_requests_per_minute") or 0)
             for b in fetch_stash_boxes(stash)]
    if not boxes:
        sys.exit("No stash-box endpoints configured in stash.")
    print(f"stash-box sources (priority order): {', '.join(b.name for b in boxes)}")

    scenes = fetch_scenes(
        stash,
        only_unorganized=not args.all_scenes,
        phashed_only=args.phashed_only,
    )
    if args.limit:
        scenes = scenes[: args.limit]

    # only scenes that actually have a fingerprint
    pending = [s for s in scenes if scene_fingerprints(s)]

    # diagnostic: how many of these have phash vs only oshash?
    def has_algo(s, name):
        return any(fp.get("algorithm") == name for fp in scene_fingerprints(s))

    with_phash = sum(1 for s in pending if has_algo(s, "PHASH"))
    oshash_only = len(pending) - with_phash
    flags = []
    if args.phashed_only:
        flags.append("phashed-only")
    if args.all_scenes:
        flags.append("all-scenes")
    if add_tags:
        flags.append("tags:merge")
    if dry:
        flags.append("dry-run")
    print(f"{len(pending)} candidate scene(s) "
          f"[{with_phash} with phash+oshash, {oshash_only} with oshash only]"
          f"{' (' + ', '.join(flags) + ')' if flags else ''}")

    studio_cache, perf_cache, tag_cache = {}, {}, {}
    matched = skipped_multi = no_match = applied = 0

    for box in boxes:
        if not pending:
            break
        print(f"\n== querying {box.name} ({len(pending)} scenes) ==")
        still = []
        for i in range(0, len(pending), args.batch):
            chunk = pending[i : i + args.batch]
            fps = [scene_fingerprints(s) for s in chunk]
            try:
                results = box.q(SB_FP_QUERY, {"fingerprints": fps})["findScenesBySceneFingerprints"]
            except RuntimeError as e:
                print(f"  ! {box.name} query failed: {e}")
                still.extend(chunk)
                continue
            for scene, res in zip(chunk, results):
                res = res or []
                if not res:
                    still.append(scene)
                    continue
                if len(res) > 1 and not args.allow_multiple:
                    skipped_multi += 1
                    continue
                match = res[0]
                matched += 1
                upd = build_update(stash, scene, match, box.endpoint, args.set_organized,
                                   studio_cache, perf_cache, tag_cache, add_tags, dry)
                title = match.get("title") or "(untitled)"
                if dry:
                    print(f"  [match] scene {scene['id']} -> \"{title}\" via {box.name}; would set {sorted(k for k in upd if k!='id')}")
                else:
                    stash.q("mutation($i:SceneUpdateInput!){sceneUpdate(input:$i){id}}", {"i": upd})
                    applied += 1
                    print(f"  [applied] scene {scene['id']} -> \"{title}\" via {box.name}")
        pending = still
        no_match = len(pending)

    print(f"\nsummary: matched={matched} applied={applied} skipped_multiple={skipped_multi} "
          f"no_match={no_match} {'(dry-run — re-run with --apply)' if dry else ''}")


if __name__ == "__main__":
    main()
