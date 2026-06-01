#!/usr/bin/env python3
"""
facerec.py — last-resort performer tagging by FACE RECOGNITION.

The fingerprint identifier (external_identify.py) and the filename parser
(external_filename_parse.py) between them clear most of the library. What's left
is the dark tail: ~446 scenes that no stash-box fingerprint and no filename token
could identify. For those, the only remaining signal is *who is on screen*. This
tool recognizes performers' faces.

It does this in two steps (two subcommands):

  build-embeddings
      Fetch every performer from stash (id + name + image), download each
      performer image, run InsightFace (buffalo_l) face detection + 512-d ArcFace
      embedding, and write a per-performer reference embedding library to disk
      (an .npz). Performers with no/blank image are skipped. If a performer image
      yields several faces, the largest face is used (a portrait's subject); the
      stored vector is L2-normalized so cosine similarity is a dot product.

  tag-scenes
      For target scenes (default `is_missing:"performers"`), extract a few frames
      via ffmpeg at fixed fractions of the duration (default 10/30/50/70/90%),
      detect + embed every face in each frame, match each face against the
      embedding library by cosine similarity (threshold --min-similarity, default
      0.5), and propose performer links for confident matches. It LINKS EXISTING
      PERFORMERS ONLY — never creates — exactly like external_filename_parse.py.

SAFETY: tag-scenes defaults to --dry-run (prints proposed scene -> performers,
writes nothing). --apply writes via sceneUpdate, MERGING the proposed performer
ids on top of the scene's existing performers (never replacing). Face recognition
is the lowest-fidelity identifier here — always review a dry-run before --apply,
and tune --min-similarity on a known-good sample first (see the README).

This runs on the Windows + RTX box (GPU), not on the stash host: stash media paths
are translated to the worker-side SMB/local path with --media-prefix
"/data=\\\\overwatch-stash\\share" (mirrors the worker's PrefixRewriter).

The heavy deps (insightface, onnxruntime-gpu, cv2, numpy) are imported lazily
inside the functions that need them, so `python -m py_compile facerec.py` and
`--help` work on a box where they aren't installed yet.

Usage:
  # Step 1 (once, and after adding performers/images): build the reference library
  python facerec.py build-embeddings \
      --stash-url http://overwatch-stash:9999 \
      --embeddings-file performers.npz

  # Step 2: dry-run face tagging over the unidentified tail
  python facerec.py tag-scenes \
      --stash-url http://overwatch-stash:9999 \
      --embeddings-file performers.npz \
      --media-prefix "/data=\\\\overwatch-stash\\share" \
      --ffmpeg ffmpeg

  # apply once the dry-run looks right
  python facerec.py tag-scenes ... --apply
"""

import argparse
import io
import json
import os
import re
import subprocess
import sys
import urllib.error
import urllib.request


# ── GraphQL plumbing (mirrors tools/identify/external_identify.py gql()/Stash) ──

def gql(url, query, variables=None, headers=None, timeout=60):
    body = json.dumps({"query": query, "variables": variables or {}}).encode()
    req = urllib.request.Request(url, data=body, method="POST")
    req.add_header("Content-Type", "application/json")
    # stash itself doesn't require it, but a browser User-Agent keeps us off any
    # bot-blocking middleware (same reasoning as external_identify.py).
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
        self.base = url.rstrip("/")
        self.url = self.base + "/graphql"
        self.api_key = api_key
        self.headers = {"ApiKey": api_key} if api_key else {}

    def q(self, query, variables=None):
        return gql(self.url, query, variables, self.headers)

    def fetch_bytes(self, url, timeout=60):
        """GET raw bytes from a stash URL (e.g. a performer image). Sends the API
        key both as the ApiKey header and as the ?apikey= query param, since the
        image endpoint reads the query param when there's no header session."""
        if self.api_key and "apikey=" not in url:
            sep = "&" if "?" in url else "?"
            url = f"{url}{sep}apikey={self.api_key}"
        req = urllib.request.Request(url, method="GET")
        req.add_header("User-Agent",
                       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "
                       "(KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
        for k, v in self.headers.items():
            req.add_header(k, v)
        with urllib.request.urlopen(req, timeout=timeout) as resp:
            return resp.read()


# ── stash queries ───────────────────────────────────────────────────────────────

PERFORMERS_Q = """
query($filter: FindFilterType) {
  findPerformers(filter: $filter) {
    count
    performers { id name image_path }
  }
}
"""

# Scenes missing performers, with the file path + duration we need for frame
# extraction. duration lives on the scene's video file.
SCENES_Q = """
query($filter: FindFilterType, $scene_filter: SceneFilterType) {
  findScenes(filter: $filter, scene_filter: $scene_filter) {
    count
    scenes {
      id title
      performers { id }
      files { path duration }
    }
  }
}
"""


def fetch_performers(stash, per_page=200):
    out, page = [], 1
    while True:
        data = stash.q(PERFORMERS_Q, {
            "filter": {"per_page": per_page, "page": page, "sort": "id", "direction": "ASC"},
        })
        items = data["findPerformers"]["performers"]
        out.extend(items)
        if len(items) < per_page:
            break
        page += 1
    return out


def fetch_scenes(stash, scene_filter, limit=0, per_page=100):
    out, page = [], 1
    while True:
        data = stash.q(SCENES_Q, {
            "filter": {"per_page": per_page, "page": page, "sort": "id", "direction": "ASC"},
            "scene_filter": scene_filter,
        })
        items = data["findScenes"]["scenes"]
        out.extend(items)
        if len(items) < per_page or (limit and len(out) >= limit):
            break
        page += 1
    if limit:
        out = out[:limit]
    return out


def performer_image_url(stash, performer):
    """Prefer the image_path stash returns (absolute, already includes any cache
    token); fall back to the documented /performer/<id>/image endpoint."""
    ip = (performer.get("image_path") or "").strip()
    if ip:
        if ip.startswith("http://") or ip.startswith("https://"):
            return ip
        return stash.base + (ip if ip.startswith("/") else "/" + ip)
    return f"{stash.base}/performer/{performer['id']}/image"


# ── InsightFace wrapper (heavy deps imported lazily) ────────────────────────────

def build_face_app(det_size=640, ctx_id=0):
    """Construct an InsightFace FaceAnalysis app on the GPU (buffalo_l).

    onnxruntime-gpu must be installed for CUDAExecutionProvider; if CUDA isn't
    available it transparently falls back to CPU (slow). ctx_id>=0 selects the
    GPU device; ctx_id=-1 forces CPU. Imported here so py_compile/--help don't
    need the package present."""
    from insightface.app import FaceAnalysis  # noqa: WPS433 (lazy import on purpose)
    app = FaceAnalysis(
        name="buffalo_l",
        providers=["CUDAExecutionProvider", "CPUExecutionProvider"],
    )
    app.prepare(ctx_id=ctx_id, det_size=(det_size, det_size))
    return app


def decode_image(raw_bytes):
    """Decode image bytes (jpg/png/webp) to a BGR numpy array for InsightFace.
    Returns None if the bytes aren't a decodable image."""
    import numpy as np
    import cv2
    arr = np.frombuffer(raw_bytes, dtype=np.uint8)
    img = cv2.imdecode(arr, cv2.IMREAD_COLOR)  # BGR, drops alpha
    return img


def normalized_embedding(face):
    """Return the L2-normalized 512-d embedding for an InsightFace face.
    Prefer the precomputed normed_embedding; fall back to embedding/||embedding||."""
    import numpy as np
    emb = getattr(face, "normed_embedding", None)
    if emb is None:
        emb = face.embedding
        n = np.linalg.norm(emb)
        emb = emb / n if n else emb
    return np.asarray(emb, dtype=np.float32)


def largest_face(faces):
    """Pick the face with the largest bounding-box area (the portrait subject)."""
    def area(f):
        x1, y1, x2, y2 = f.bbox
        return max(0.0, (x2 - x1)) * max(0.0, (y2 - y1))
    return max(faces, key=area)


# ── build-embeddings ─────────────────────────────────────────────────────────────

def cmd_build_embeddings(args):
    import numpy as np

    stash = Stash(args.stash_url, args.stash_api_key)
    print("loading performers from stash …")
    performers = fetch_performers(stash)
    print(f"  {len(performers)} performers returned.")

    print(f"initializing InsightFace (buffalo_l) on ctx_id={args.ctx_id} …")
    app = build_face_app(det_size=args.det_size, ctx_id=args.ctx_id)

    ids, names, vectors = [], [], []
    no_image = no_face = decode_fail = fetch_fail = 0

    for i, p in enumerate(performers, 1):
        pid, name = p["id"], p.get("name") or f"performer {p['id']}"
        url = performer_image_url(stash, p)
        if not url:
            no_image += 1
            continue
        try:
            raw = stash.fetch_bytes(url)
        except Exception as e:  # noqa: BLE001 — network is best-effort per performer
            fetch_fail += 1
            print(f"  ! {name} (id {pid}): image fetch failed: {e}")
            continue
        # stash returns a tiny default avatar for performers with no real image;
        # treat suspiciously small payloads as "no image".
        if not raw or len(raw) < args.min_image_bytes:
            no_image += 1
            continue
        img = decode_image(raw)
        if img is None:
            decode_fail += 1
            print(f"  ! {name} (id {pid}): could not decode image")
            continue
        faces = app.get(img)
        if not faces:
            no_face += 1
            continue
        emb = normalized_embedding(largest_face(faces))
        ids.append(str(pid))
        names.append(name)
        vectors.append(emb)
        if i % 50 == 0:
            print(f"  … processed {i}/{len(performers)} (embedded {len(ids)})")

    if not vectors:
        sys.exit("No performer embeddings were produced — nothing to write. "
                 "Check that performers have images and that InsightFace loaded.")

    mat = np.stack(vectors).astype(np.float32)  # (N, 512), each row L2-normalized
    np.savez_compressed(
        args.embeddings_file,
        ids=np.asarray(ids, dtype=object),
        names=np.asarray(names, dtype=object),
        embeddings=mat,
    )
    print(f"\nwrote {len(ids)} performer embeddings -> {args.embeddings_file}")
    print(f"summary: embedded={len(ids)} no_image={no_image} no_face_detected={no_face} "
          f"decode_fail={decode_fail} fetch_fail={fetch_fail}")


def load_embeddings(path):
    """Load the .npz into (ids[list], names[list], matrix[(N,512) float32])."""
    import numpy as np
    if not os.path.exists(path):
        sys.exit(f"embeddings file not found: {path} — run `build-embeddings` first.")
    data = np.load(path, allow_pickle=True)
    ids = [str(x) for x in data["ids"]]
    names = [str(x) for x in data["names"]]
    mat = np.asarray(data["embeddings"], dtype=np.float32)
    return ids, names, mat


# ── frame extraction (mirrors worker/internal/cover.go ExtractCover) ────────────

class PrefixRewriter:
    """Translate a stash-side media path to the worker-side (SMB/local) path.

    Mirrors worker/internal/paths.go: parse "STASH_PREFIX=WORKER_PREFIX", strip the
    stash prefix, prepend the worker prefix, and (when the worker side is a Windows
    UNC/backslash path) flip separators. Both shells collapse "\\\\" to "\\" before
    argv, so a malformed leading single-backslash UNC is normalized back."""

    def __init__(self, from_prefix, to_prefix):
        self.from_prefix = from_prefix.rstrip("/\\")
        self.to_prefix = self._normalize_unc(to_prefix).rstrip("/\\")

    @classmethod
    def parse(cls, spec):
        if not spec:
            return None
        idx = spec.find("=")
        if idx < 1 or idx == len(spec) - 1:
            raise ValueError(f"expected STASH_PREFIX=WORKER_PREFIX, got {spec!r}")
        return cls(spec[:idx].strip(), spec[idx + 1:].strip())

    @staticmethod
    def _normalize_unc(p):
        if len(p) < 2:
            return p
        if p[0].isalpha() and p[1] == ":":  # drive letter C:\ — not UNC
            return p
        if (p[0] == "\\" and p[1] == "\\") or (p[0] == "/" and p[1] == "/"):
            return p  # already UNC
        if p[0] == "\\" and (p[1].isalnum()):  # half-escaped \server -> \\server
            return "\\" + p
        return p

    def rewrite(self, stash_path):
        if not self.from_prefix:
            return stash_path
        if stash_path == self.from_prefix:
            return self.to_prefix
        if stash_path.startswith(self.from_prefix + "/") or stash_path.startswith(self.from_prefix + "\\"):
            rest = stash_path[len(self.from_prefix):]
            if "\\" in self.to_prefix:
                rest = rest.replace("/", "\\")
            return self.to_prefix + rest
        return stash_path


def parse_frame_fractions(spec):
    """'10,30,50,70,90' -> [0.10, 0.30, 0.50, 0.70, 0.90] (percent of duration)."""
    out = []
    for tok in spec.split(","):
        tok = tok.strip()
        if not tok:
            continue
        v = float(tok)
        out.append(v / 100.0 if v > 1.0 else v)
    return out or [0.10, 0.30, 0.50, 0.70, 0.90]


def extract_frame(ffmpeg_path, source, at_seconds):
    """Seek to at_seconds and return one JPEG frame as raw bytes (or None).

    Mirrors worker/internal/cover.go: `-ss <t> -i <src> -frames:v 1 -q:v 2
    -f mjpeg pipe:1`, with -ss before -i for a fast input seek."""
    if at_seconds < 0:
        at_seconds = 0
    args = [
        ffmpeg_path,
        "-loglevel", "error",
        "-ss", f"{at_seconds:.3f}",
        "-i", source,
        "-frames:v", "1",
        "-q:v", "2",
        "-f", "mjpeg",
        "pipe:1",
    ]
    try:
        proc = subprocess.run(args, capture_output=True, timeout=120)
    except (subprocess.TimeoutExpired, FileNotFoundError) as e:
        return None, str(e)
    if proc.returncode != 0 or not proc.stdout:
        return None, (proc.stderr.decode(errors="replace")[:200] if proc.stderr else "no output")
    return proc.stdout, None


# ── tag-scenes ───────────────────────────────────────────────────────────────────

def scene_path_and_duration(scene):
    files = scene.get("files") or []
    if not files:
        return None, 0.0
    f = files[0]
    return f.get("path"), float(f.get("duration") or 0.0)


def cmd_tag_scenes(args):
    import numpy as np

    # Filenames/titles carry unicode; force lossy UTF-8 console output so a print
    # never aborts the run (same guard as external_filename_parse.py).
    for stream in (sys.stdout, sys.stderr):
        try:
            stream.reconfigure(encoding="utf-8", errors="replace")
        except (AttributeError, ValueError):
            pass

    dry = not args.apply
    rewriter = PrefixRewriter.parse(args.media_prefix)
    fractions = parse_frame_fractions(args.frames)

    ids, names, mat = load_embeddings(args.embeddings_file)
    print(f"loaded {len(ids)} performer embeddings from {args.embeddings_file}")

    stash = Stash(args.stash_url, args.stash_api_key)

    # default target: scenes missing performers. --scene-filter overrides with raw
    # JSON for the "dark tail" (e.g. '{"is_missing":"performers","organized":false}').
    if args.scene_filter:
        scene_filter = json.loads(args.scene_filter)
    else:
        scene_filter = {"is_missing": "performers"}
    scenes = fetch_scenes(stash, scene_filter, limit=args.limit)
    print(f"{len(scenes)} candidate scene(s){' (dry-run)' if dry else ''}\n")

    print(f"initializing InsightFace (buffalo_l) on ctx_id={args.ctx_id} …")
    app = build_face_app(det_size=args.det_size, ctx_id=args.ctx_id)

    proposed = applied = no_path = no_duration = no_face = no_match = 0

    for s in scenes:
        sid = s["id"]
        stash_path, duration = scene_path_and_duration(s)
        if not stash_path:
            no_path += 1
            continue
        if duration <= 0:
            no_duration += 1
            print(f"  ! scene {sid}: no duration on file, skipping ({stash_path})")
            continue
        source = rewriter.rewrite(stash_path) if rewriter else stash_path

        # best cosine similarity per performer across all sampled frames/faces
        best = {}  # performer_index -> best cosine similarity seen
        faces_seen = 0
        for frac in fractions:
            t = duration * frac
            raw, err = extract_frame(args.ffmpeg, source, t)
            if raw is None:
                if err:
                    print(f"  ! scene {sid}: frame at {t:.1f}s failed: {err}")
                continue
            img = decode_image(raw)
            if img is None:
                continue
            faces = app.get(img)
            faces_seen += len(faces)
            for face in faces:
                emb = normalized_embedding(face)
                # both sides L2-normalized => cosine similarity is a dot product;
                # vectorized over the whole library at once.
                sims = mat @ emb  # (N,)
                j = int(np.argmax(sims))
                score = float(sims[j])
                if score >= args.min_similarity and score > best.get(j, -1.0):
                    best[j] = score

        if not faces_seen:
            no_face += 1
            continue
        if not best:
            no_match += 1
            continue

        # rank matched performers by confidence; cap with --max-performers
        ranked = sorted(best.items(), key=lambda kv: -kv[1])[: args.max_performers]
        match_ids = [ids[j] for j, _ in ranked]
        desc = ", ".join(f"{names[j]}({score:.2f})" for j, score in ranked)

        existing = [p["id"] for p in (s.get("performers") or [])]
        merged = existing + [pid for pid in match_ids if pid not in existing]
        if merged == existing:
            no_match += 1
            continue

        proposed += 1
        title = s.get("title") or "(untitled)"
        if dry:
            print(f"  scene {sid} \"{title[:50]}\": performers <- [{desc}]")
        else:
            try:
                stash.q("mutation($i:SceneUpdateInput!){sceneUpdate(input:$i){id}}",
                        {"i": {"id": sid, "performer_ids": merged}})
                applied += 1
                print(f"  [applied] scene {sid} \"{title[:50]}\": performers <- [{desc}]")
            except RuntimeError as e:
                print(f"  ! scene {sid} update failed: {e}")

    print(f"\nsummary: proposed={proposed} "
          f"{'applied=' + str(applied) if not dry else '(dry-run — re-run with --apply)'} "
          f"no_face={no_face} no_match={no_match} no_path={no_path} no_duration={no_duration}")


# ── CLI ──────────────────────────────────────────────────────────────────────────

def build_parser():
    ap = argparse.ArgumentParser(
        description="Face-recognition performer tagger for stash (InsightFace buffalo_l, GPU).")
    sub = ap.add_subparsers(dest="command", required=True)

    # shared flags helper
    def add_common(p):
        p.add_argument("--stash-url", default="http://overwatch-stash:9999")
        p.add_argument("--stash-api-key", default=None)
        p.add_argument("--embeddings-file", default="performers.npz",
                       help="path to the performer embedding library (.npz)")
        p.add_argument("--ctx-id", type=int, default=0,
                       help="InsightFace device id: 0 = first GPU, -1 = CPU")
        p.add_argument("--det-size", type=int, default=640,
                       help="InsightFace detector input size (square)")

    be = sub.add_parser("build-embeddings",
                        help="fetch performers, embed their images, write the .npz library")
    add_common(be)
    be.add_argument("--min-image-bytes", type=int, default=2048,
                    help="treat image payloads smaller than this as 'no image' (default avatars)")
    be.set_defaults(func=cmd_build_embeddings)

    ts = sub.add_parser("tag-scenes",
                        help="match scene frames against the library and propose/apply performer links")
    add_common(ts)
    ts.add_argument("--media-prefix", default=None,
                    help='stash->worker path rewrite, e.g. "/data=\\\\overwatch-stash\\share"')
    ts.add_argument("--ffmpeg", default="ffmpeg", help="path to the ffmpeg binary")
    ts.add_argument("--frames", default="10,30,50,70,90",
                    help="comma-separated percentages of duration to sample")
    ts.add_argument("--min-similarity", type=float, default=0.5,
                    help="cosine-similarity threshold for a confident match (0..1)")
    ts.add_argument("--max-performers", type=int, default=4,
                    help="cap performers proposed per scene (highest-confidence first)")
    ts.add_argument("--scene-filter", default=None,
                    help='raw JSON SceneFilterType overriding the default is_missing:"performers"')
    ts.add_argument("--limit", type=int, default=0, help="cap scenes processed (0 = no cap)")
    ts.add_argument("--apply", action="store_true", help="write changes (default: dry-run)")
    ts.set_defaults(func=cmd_tag_scenes)

    return ap


def main():
    args = build_parser().parse_args()
    args.func(args)


if __name__ == "__main__":
    main()
