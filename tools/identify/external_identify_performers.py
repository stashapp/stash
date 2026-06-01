#!/usr/bin/env python3
"""
external_identify_performers.py — enrich stash performers with metadata from
stash-box endpoints (StashDB, ThePornDB, ...) WITHOUT using stash's internal
job queue.

Companion to external_identify.py (which identifies scenes). Where the scene
tool fingerprint-matches scenes, this tool either:

  * refresh mode — for each performer that already has a stash_id but is
    missing fields, fetches the canonical record via stash-box findPerformer
    and MERGE-fills the empty fields.
  * search  mode — for each unlinked performer, runs stash-box
    searchPerformer(term=name), post-filters by case-insensitive exact name
    (mirroring strings.EqualFold in pkg/stashbox/performer.go:311-325), and
    merges the single survivor (or first if --allow-multiple).
  * both    — refresh-then-search.

Translation logic (measurements format, body-mod join, breast_type ->
fake_tits, ethnicity/eye_color/hair_color enum -> string, height>0 skip,
alias case-equal-name filter, merged_into_id one-hop redirect, images[0]
convention) is a direct port of pkg/stashbox/performer.go. Do not re-derive
formats — the line refs are cited inline.

Multi-value fields (alias_list, urls, stash_ids) are passed in Set mode, which
REPLACES wholesale (see internal/api/resolver_mutation_performer.go:401), so
this tool always unions the new values with the performer's existing values
before sending performerUpdate.

SAFETY: defaults to dry-run. Use --apply to write.

Do NOT run this concurrently with native stash performer scraping
(last write wins under performerUpdate's single txn).

Usage:
  python3 external_identify_performers.py --stash-url http://localhost:9999
  python3 external_identify_performers.py --stash-url http://localhost:9999 --apply --mode both
  python3 external_identify_performers.py --stash-url http://overwatch-stash:9999 \
      --stash-api-key XXX --apply --mode refresh
"""

import argparse
import json
import sys
import time
import urllib.request
import urllib.error

# ── GraphQL plumbing (copied verbatim from external_identify.py) ──────────────


def gql(url, query, variables=None, headers=None, timeout=60):
    body = json.dumps({"query": query, "variables": variables or {}}).encode()
    req = urllib.request.Request(url, data=body, method="POST")
    req.add_header("Content-Type", "application/json")
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


# ── stash-box fragment + queries (mirrors graphql/stash-box/query.graphql) ────
#
# PerformerFragment matches the bundled fragment in stash exactly so the result
# can be fed straight into the port of pkg/stashbox/performer.go.

SB_PERFORMER_FRAGMENT = """
fragment PerformerFragment on Performer {
  id
  name
  disambiguation
  aliases
  gender
  deleted
  merged_into_id
  urls { url type }
  images { id url width height }
  birth_date
  death_date
  ethnicity
  country
  eye_color
  hair_color
  height
  measurements { band_size cup_size waist hip }
  breast_type
  career_start_year
  career_end_year
  tattoos { location description }
  piercings { location description }
}
"""

SB_FIND_PERFORMER_Q = SB_PERFORMER_FRAGMENT + """
query($id: ID!) {
  findPerformer(id: $id) { ...PerformerFragment }
}
"""

SB_SEARCH_PERFORMER_Q = SB_PERFORMER_FRAGMENT + """
query($term: String!) {
  searchPerformer(term: $term) { ...PerformerFragment }
}
"""


# ── stash queries ─────────────────────────────────────────────────────────────

STASH_BOXES_Q = "{ configuration { general { stashBoxes { endpoint api_key name max_requests_per_minute } } } }"

# Pull every field this tool might write so we can decide what to send.
# alias_list/urls/stash_ids are pulled so we can union (Set replaces wholesale —
# see internal/api/resolver_mutation_performer.go:401).
PERFORMERS_Q = """
query($filter: FindFilterType) {
  findPerformers(filter: $filter) {
    count
    performers {
      id
      name
      disambiguation
      gender
      birthdate
      death_date
      ethnicity
      country
      eye_color
      hair_color
      height_cm
      measurements
      fake_tits
      career_start
      career_end
      tattoos
      piercings
      details
      image_path
      alias_list
      urls
      stash_ids { endpoint stash_id }
    }
  }
}
"""


def fetch_stash_boxes(stash):
    data = stash.q(STASH_BOXES_Q)
    return data["configuration"]["general"]["stashBoxes"] or []


def fetch_all_performers(stash, per_page=100):
    """Fetch every performer in stash via findPerformers (paginated).

    PerformerFilterType.is_missing accepts only one value per call
    (pkg/sqlite/performer_filter.go:331-365), with no compound "all of these
    are empty" filter, so we filter client-side after pulling everything.
    """
    page, out = 1, []
    while True:
        data = stash.q(PERFORMERS_Q, {
            "filter": {"per_page": per_page, "page": page, "sort": "id", "direction": "ASC"},
        })
        performers = data["findPerformers"]["performers"]
        out.extend(performers)
        if len(performers) < per_page:
            break
        page += 1
    return out


# ── port of pkg/stashbox/performer.go ─────────────────────────────────────────


def _is_present(v):
    """Treat None and (string) empty/whitespace as absent.

    Mirrors stash's "IS NULL OR TRIM(field) = ''" emptiness check used in
    pkg/sqlite/performer_filter.go:362 for the is_missing filter.
    """
    if v is None:
        return False
    if isinstance(v, str) and v.strip() == "":
        return False
    if isinstance(v, (list, tuple)) and len(v) == 0:
        return False
    return True


def enum_to_string(e):
    """Port of enumToStringPtr(e, titleCase=true) — performer.go:91-102.

    stash-box enum names are SHOUTING_SNAKE. The Go path replaces underscores
    with spaces and then Title-cases the result (e.g. EYE_COLOR_BLUE -> "Blue",
    NATURAL -> "Natural", LIGHT_BROWN -> "Light Brown"). We do the same:
    lowercase first, then title-case each space-separated word.

    NOTE: this is called for ethnicity, eye_color, hair_color, breast_type
    in performer.go:257,261,265,269 — all with titleCase=true. The design doc
    says "lowercase string"; the Go file says title-case. We follow the Go
    file (CLAUDE.md says port byte-for-byte).
    """
    if e is None or e == "":
        return None
    s = str(e).replace("_", " ").lower()
    # Title-case each whitespace-separated word (mirrors golang.org/x/text/cases.Title)
    return " ".join(w[:1].upper() + w[1:] for w in s.split(" ") if w != "") or None


def format_measurements(m):
    """Port of formatMeasurements — performer.go:128-135.

    Only returns a value when ALL FOUR sub-fields are non-null. Format is
    "<band><cup>-<waist>-<hip>" with band/waist/hip as ints and cup as a
    string (e.g. "32D-24-34").
    """
    if not m:
        return None
    band = m.get("band_size")
    cup = m.get("cup_size")
    waist = m.get("waist")
    hip = m.get("hip")
    if band is None or cup is None or waist is None or hip is None:
        return None
    return f"{band}{cup}-{waist}-{hip}"


def format_body_modifications(mods):
    """Port of formatBodyModifications — performer.go:155-171.

    Each entry is either "<location>" (when description is nil) or
    "<location>, <description>". Entries joined by "; ".
    """
    if not mods:
        return None
    parts = []
    for f in mods:
        loc = f.get("location") or ""
        desc = f.get("description")
        if desc is None or desc == "":
            parts.append(loc)
        else:
            parts.append(f"{loc}, {desc}")
    return "; ".join(parts) if parts else None


def translate_gender(g):
    """Port of translateGender — performer.go:104-126.

    stash-box GenderEnum names map 1:1 to stash GenderEnum names. The Go file
    explicitly enumerates them; we whitelist the same set and pass through.
    """
    if g is None or g == "":
        return None
    allowed = {
        "MALE", "FEMALE", "INTERSEX",
        "TRANSGENDER_MALE", "TRANSGENDER_FEMALE",
        "NON_BINARY",
    }
    return g if g in allowed else None


def filter_aliases(aliases, name):
    """Port of the aliases filter — performer.go:272-284.

    #4437: drop entries case-equal to the performer's own name.
    #4596: case-fold dedupe (UniqueFold).
    """
    if not aliases:
        return []
    out = []
    seen = set()
    name_lc = (name or "").lower()
    for a in aliases:
        if not a:
            continue
        a_lc = a.lower()
        if a_lc == name_lc:
            continue
        if a_lc in seen:
            continue
        seen.add(a_lc)
        out.append(a)
    return out


# ── image pre-HEAD ────────────────────────────────────────────────────────────


def head_ok(url, timeout=5):
    """HEAD the image URL; return True only on 2xx.

    Pre-flighting per design doc §4: avoid full performerUpdate rollback when
    the CDN URL is dead. Some CDNs reject HEAD with 405; treat that as failure
    too — we'd rather skip the image than blow up the whole mutation.
    """
    try:
        req = urllib.request.Request(url, method="HEAD")
        with urllib.request.urlopen(req, timeout=timeout) as resp:
            return 200 <= resp.status < 300
    except Exception:
        return False


# ── candidate filtering ───────────────────────────────────────────────────────

# Fields whose emptiness on stash's side makes a performer a refresh-mode
# candidate. Mirrors the stash-side enrichable set (the design doc §4 table).
ENRICHABLE_FIELDS = (
    "disambiguation",
    "gender",
    "birthdate",
    "death_date",
    "ethnicity",
    "country",
    "eye_color",
    "hair_color",
    "height_cm",
    "measurements",
    "fake_tits",
    "career_start",
    "career_end",
    "tattoos",
    "piercings",
    "image_path",
    "alias_list",
    "urls",
)


def some_field_empty(performer):
    return any(not _is_present(performer.get(f)) for f in ENRICHABLE_FIELDS)


# ── build update from PerformerFragment (the MERGE + union pass) ──────────────


def build_performer_update(stash_perf, sb_perf, sb_endpoint, dry):
    """Construct a minimal PerformerUpdateInput.

    Only includes fields where:
      - stash's current value is empty (MERGE rule, design doc §4), AND
      - stash-box has a non-empty value after the performer.go port.

    For the multi-value fields (alias_list, urls, stash_ids) we always send
    the union with existing values because Set mode replaces wholesale
    (internal/api/resolver_mutation_performer.go:401). They appear in the
    update whenever the union differs from what's currently stored.

    Returns (upd, changes) where upd is the PerformerUpdateInput dict and
    changes is the sorted list of fields touched (for diagnostic output).
    """
    upd = {"id": stash_perf["id"]}

    def fill_if_empty(stash_key, sb_value):
        if _is_present(sb_value) and not _is_present(stash_perf.get(stash_key)):
            upd[stash_key] = sb_value

    # Never overwrite the existing name in refresh mode (design doc §4 table).
    # We don't include `name` at all — stash already has it.

    fill_if_empty("disambiguation", sb_perf.get("disambiguation"))
    fill_if_empty("gender", translate_gender(sb_perf.get("gender")))
    fill_if_empty("birthdate", sb_perf.get("birth_date"))
    fill_if_empty("death_date", sb_perf.get("death_date"))
    fill_if_empty("ethnicity", enum_to_string(sb_perf.get("ethnicity")))
    fill_if_empty("country", sb_perf.get("country"))
    fill_if_empty("eye_color", enum_to_string(sb_perf.get("eye_color")))
    fill_if_empty("hair_color", enum_to_string(sb_perf.get("hair_color")))

    # height > 0 skip — performer.go:229
    h = sb_perf.get("height")
    if isinstance(h, int) and h > 0:
        fill_if_empty("height_cm", h)

    fill_if_empty("measurements", format_measurements(sb_perf.get("measurements")))

    # breast_type -> fake_tits — performer.go:268-270
    fill_if_empty("fake_tits", enum_to_string(sb_perf.get("breast_type")))

    # career_start/end are strings on stash's side — performer.go:234-242
    cs = sb_perf.get("career_start_year")
    if isinstance(cs, int):
        fill_if_empty("career_start", str(cs))
    ce = sb_perf.get("career_end_year")
    if isinstance(ce, int):
        fill_if_empty("career_end", str(ce))

    fill_if_empty("tattoos", format_body_modifications(sb_perf.get("tattoos")))
    fill_if_empty("piercings", format_body_modifications(sb_perf.get("piercings")))

    # ── multi-value union fields (Set mode replaces — must union) ─────────
    # alias_list: filter the sb aliases (drop name-equal, fold-dedupe), then
    # union with existing stash aliases under case-fold dedupe.
    sb_aliases = filter_aliases(sb_perf.get("aliases") or [], sb_perf.get("name") or stash_perf.get("name"))
    existing_aliases = stash_perf.get("alias_list") or []
    union_aliases = _casefold_union(existing_aliases, sb_aliases, drop_equal=stash_perf.get("name"))
    if _set_differs(existing_aliases, union_aliases):
        upd["alias_list"] = union_aliases

    # urls: flatten sb's [{url,type}] to [url], union with existing.
    sb_urls = [u.get("url") for u in (sb_perf.get("urls") or []) if u and u.get("url")]
    existing_urls = stash_perf.get("urls") or []
    union_urls = _casefold_union(existing_urls, sb_urls)
    if _set_differs(existing_urls, union_urls):
        upd["urls"] = union_urls

    # stash_ids: keep every existing {endpoint, stash_id}; add the new one if
    # the (endpoint, id) pair isn't already there. Design doc §5: never drop
    # entries from other endpoints.
    existing_sids = stash_perf.get("stash_ids") or []
    existing_pairs = {(s["endpoint"], s["stash_id"]) for s in existing_sids}
    new_pair = (sb_endpoint, sb_perf["id"])
    if new_pair not in existing_pairs:
        upd["stash_ids"] = (
            [{"endpoint": s["endpoint"], "stash_id": s["stash_id"]} for s in existing_sids]
            + [{"endpoint": sb_endpoint, "stash_id": sb_perf["id"]}]
        )

    # ── image: pre-HEAD only if we'd actually set it ──────────────────────
    # images[0] convention — performer.go:225-227. We only consider sending
    # an image when stash has none (image_path empty).
    if not _is_present(stash_perf.get("image_path")):
        imgs = sb_perf.get("images") or []
        if imgs and imgs[0].get("url"):
            url = imgs[0]["url"]
            if head_ok(url):
                upd["image"] = url
            # else: silently skip just the image, keep the rest of the update.

    changes = sorted(k for k in upd if k != "id")
    return upd, changes


def _casefold_union(existing, new_values, drop_equal=None):
    """Union two string lists, case-fold dedupe, preserve order (existing first).

    If drop_equal is set, drop entries case-equal to that string (used for
    alias_list to honor the performer.go #4437 rule against name-equal
    aliases even when they come from the existing stash data).
    """
    out = []
    seen = set()
    drop_lc = drop_equal.lower() if drop_equal else None
    for v in list(existing) + list(new_values):
        if v is None:
            continue
        lc = v.lower()
        if drop_lc is not None and lc == drop_lc:
            continue
        if lc in seen:
            continue
        seen.add(lc)
        out.append(v)
    return out


def _set_differs(a, b):
    """True iff the case-fold sets differ (used to skip no-op multi-value writes)."""
    return {x.lower() for x in (a or [])} != {x.lower() for x in (b or [])}


# ── stash-box lookups (with merged_into_id 1-hop redirect) ────────────────────


def sb_find_performer(box, sb_id, _hop=0):
    """findPerformer(id) with one-hop merged_into_id redirect.

    If the returned record has deleted=true and merged_into_id is set, refetch
    once at the redirect target. If that one is also deleted, return None.
    Design doc §5: "1-hop redirect only, no recursion".
    """
    data = box.q(SB_FIND_PERFORMER_Q, {"id": sb_id})
    p = data.get("findPerformer")
    if not p:
        return None
    if p.get("deleted"):
        merged = p.get("merged_into_id")
        if merged and _hop == 0:
            data2 = box.q(SB_FIND_PERFORMER_Q, {"id": merged})
            p2 = data2.get("findPerformer")
            if not p2 or p2.get("deleted"):
                return None
            return p2
        return None
    return p


def sb_search_performer(box, name):
    """searchPerformer(term) with EqualFold(name) post-filter.

    Mirrors performer.go:311-325 (FindPerformerByName). Returns the list of
    surviving fragments — caller decides what to do with len 0/1/many.
    """
    data = box.q(SB_SEARCH_PERFORMER_Q, {"term": name})
    results = data.get("searchPerformer") or []
    name_lc = (name or "").lower()
    return [r for r in results if r and (r.get("name") or "").lower() == name_lc]


# ── main ──────────────────────────────────────────────────────────────────────


def main():
    ap = argparse.ArgumentParser(
        description="External performer identifier for stash (parallel/off-box). "
                    "Do NOT run concurrently with native stash performer scraping.",
    )
    ap.add_argument("--stash-url", default="http://localhost:9999")
    ap.add_argument("--stash-api-key", default=None)
    ap.add_argument("--apply", action="store_true", help="actually write changes (default is dry-run)")
    ap.add_argument("--mode", choices=("refresh", "search", "both"), default="refresh",
                    help="refresh: enrich linked performers with empty fields; "
                         "search: find linkage for unlinked performers; "
                         "both: refresh then search")
    ap.add_argument("--limit", type=int, default=0, help="cap performers processed (0 = no cap)")
    ap.add_argument("--allow-multiple", action="store_true",
                    help="search mode: take first remaining result after EqualFold(name) filter "
                         "instead of skipping ambiguous")
    ap.add_argument("--per-page", type=int, default=100, help="stash pagination page size")
    args = ap.parse_args()
    dry = not args.apply

    stash = Stash(args.stash_url, args.stash_api_key)
    boxes = [StashBox(b["endpoint"], b["api_key"], b["name"], b.get("max_requests_per_minute") or 0)
             for b in fetch_stash_boxes(stash)]
    if not boxes:
        sys.exit("No stash-box endpoints configured in stash.")
    print(f"stash-box sources (priority order): {', '.join(b.name for b in boxes)}")
    boxes_by_endpoint = {b.endpoint: b for b in boxes}

    # Pull all performers; filter client-side per mode.
    all_performers = fetch_all_performers(stash, per_page=args.per_page)

    linked = [p for p in all_performers if (p.get("stash_ids") or [])]
    unlinked = [p for p in all_performers if not (p.get("stash_ids") or [])]

    if args.mode == "refresh":
        candidates = [p for p in linked if some_field_empty(p)]
    elif args.mode == "search":
        candidates = unlinked
    else:  # both
        candidates = [p for p in linked if some_field_empty(p)] + unlinked

    if args.limit:
        candidates = candidates[: args.limit]

    refresh_count = sum(1 for p in candidates if p.get("stash_ids"))
    search_count = len(candidates) - refresh_count
    mode_tag = f"{args.mode}-mode"
    print(f"{len(candidates)} candidate performer(s) "
          f"[linked: {refresh_count}, unlinked: {search_count}] ({mode_tag})")

    matched = applied = skipped_ambiguous = no_match = 0

    for perf in candidates:
        is_linked = bool(perf.get("stash_ids"))

        sb_perf = None
        chosen_endpoint = None
        ambiguous = False

        # ── refresh path: direct findPerformer per linked stash_id ────────
        if is_linked and args.mode in ("refresh", "both"):
            for sid in perf["stash_ids"]:
                box = boxes_by_endpoint.get(sid["endpoint"])
                if not box:
                    # stash-id from an endpoint not currently configured — skip it.
                    continue
                try:
                    found = sb_find_performer(box, sid["stash_id"])
                except RuntimeError as e:
                    print(f"  ! {box.name} findPerformer({sid['stash_id']}) failed: {e}")
                    continue
                if found:
                    sb_perf = found
                    chosen_endpoint = box.endpoint
                    break

        # ── search path: name-based lookup against each box in order ──────
        if sb_perf is None and (not is_linked) and args.mode in ("search", "both"):
            name = perf.get("name") or ""
            if not name:
                no_match += 1
                continue
            for box in boxes:
                try:
                    results = sb_search_performer(box, name)
                except RuntimeError as e:
                    print(f"  ! {box.name} searchPerformer('{name}') failed: {e}")
                    continue
                if not results:
                    continue
                if len(results) > 1 and not args.allow_multiple:
                    ambiguous = True
                    # Stop scanning further boxes — keep the box where the
                    # collision happened as the diagnostic source.
                    chosen_endpoint = box.endpoint
                    break
                sb_perf = results[0]
                chosen_endpoint = box.endpoint
                break

        if ambiguous:
            skipped_ambiguous += 1
            print(f"  [skip] performer {perf['id']} \"{perf.get('name')}\": "
                  f"multiple stash-box matches after EqualFold filter — "
                  f"re-run with --allow-multiple to take first")
            continue

        if sb_perf is None:
            no_match += 1
            continue

        matched += 1
        upd, changes = build_performer_update(perf, sb_perf, chosen_endpoint, dry)

        if not changes:
            # Nothing to do — stash-box record had no new info to merge.
            print(f"  [noop] performer {perf['id']} \"{perf.get('name')}\": "
                  f"no empty fields to fill from {chosen_endpoint}")
            continue

        if dry:
            print(f"  [match] performer {perf['id']} \"{perf.get('name')}\" "
                  f"<- {sb_perf.get('name')} via {chosen_endpoint}; would set {changes}")
        else:
            try:
                stash.q(
                    "mutation($i:PerformerUpdateInput!){performerUpdate(input:$i){id}}",
                    {"i": upd},
                )
                applied += 1
                print(f"  [applied] performer {perf['id']} \"{perf.get('name')}\" "
                      f"<- {sb_perf.get('name')} via {chosen_endpoint}; set {changes}")
            except RuntimeError as e:
                print(f"  ! performerUpdate({perf['id']}) failed: {e}")

    suffix = " (dry-run — re-run with --apply)" if dry else ""
    print(f"\nsummary: matched={matched} applied={applied} "
          f"skipped_ambiguous={skipped_ambiguous} no_match={no_match}{suffix}")


if __name__ == "__main__":
    main()
