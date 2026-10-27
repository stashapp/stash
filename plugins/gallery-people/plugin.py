"""Local face grouping for Stash. The only runtime HTTP destination is Stash itself."""
from contextlib import contextmanager
import hashlib
import io
import json
import math
import os
from pathlib import Path
import sqlite3
import sys
import time
import urllib.parse
import urllib.request
import uuid

ROOT = Path(__file__).resolve().parent
MANIFEST = json.loads((ROOT / "models.json").read_text())
VERSION = MANIFEST["version"]
IMAGE_FIELDS = "id title paths { image thumbnail preview } galleries { id }"


class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):
        raise RuntimeError("Stash redirected a local request; check its local server address")


class Stash:
    def __init__(self, connection):
        scheme = connection.get("Scheme", "http")
        if scheme not in ("http", "https"):
            raise ValueError("Unsupported Stash connection scheme")
        port = int(connection["Port"])
        if not 1 <= port <= 65535:
            raise ValueError("Invalid Stash port")
        self.base = f"{scheme}://localhost:{port}"
        cookie = connection.get("SessionCookie") or {}
        self.headers = {"Content-Type": "application/json"}
        if cookie.get("Value"):
            self.headers["Cookie"] = f"{cookie['Name']}={cookie['Value']}"
        # Local image data and authentication never go through an HTTP proxy.
        self.opener = urllib.request.build_opener(urllib.request.ProxyHandler({}), NoRedirect())

    def graphql(self, query, variables=None):
        request = urllib.request.Request(self.base + "/graphql",
            json.dumps({"query": query, "variables": variables or {}}).encode(), self.headers)
        with self.opener.open(request, timeout=120) as response:
            result = json.load(response)
        if result.get("errors"):
            raise RuntimeError("; ".join(e["message"] for e in result["errors"]))
        return result["data"]

    def images(self, gallery_id=None, ids=None, page=1, per_page=100):
        image_filter = {"galleries": {"value": [str(gallery_id)], "modifier": "INCLUDES"}} if gallery_id else {}
        return self.graphql("""query($f: ImageFilterType, $ids: [ID!], $page: Int!, $limit: Int!) {
          findImages(image_filter: $f, ids: $ids, filter: {page: $page, per_page: $limit, sort: "id", direction: ASC}) {
            count images { """ + IMAGE_FIELDS + """ }
          }
        }""", {"f": image_filter, "ids": ids, "page": page, "limit": per_page})["findImages"]

    def image_bytes(self, image):
        path = urllib.parse.urlsplit(image["paths"]["image"] or "")
        if not path.path.startswith("/image/") or ".." in path.path:
            raise ValueError("Stash did not supply a local image route")
        # Discard external_host: always fetch from this Stash process.
        url = self.base + path.path + ("?" + path.query if path.query else "")
        with self.opener.open(urllib.request.Request(url, headers=self.headers), timeout=120) as response:
            data = response.read(64 * 1024 * 1024 + 1)
        if len(data) > 64 * 1024 * 1024:
            raise ValueError("Image exceeds the 64 MB processing limit")
        return data


@contextmanager
def exclusive_lock(path):
    """OS-owned lock releases even when Stash forcibly cancels a raw plugin task."""
    with path.open("a+b") as lock:
        if os.name == "nt":
            import msvcrt
            lock.seek(0)
            if not lock.read(1):
                lock.write(b"0")
                lock.flush()
            lock.seek(0)
            try:
                msvcrt.locking(lock.fileno(), msvcrt.LK_NBLCK, 1)
            except OSError as exc:
                raise RuntimeError("A face scan or edit is already running. Try again when it finishes.") from exc
        else:
            import fcntl
            try:
                fcntl.flock(lock, fcntl.LOCK_EX | fcntl.LOCK_NB)
            except BlockingIOError as exc:
                raise RuntimeError("A face scan or edit is already running. Try again when it finishes.") from exc
        try:
            yield
        finally:
            if os.name == "nt":
                lock.seek(0)
                msvcrt.locking(lock.fileno(), msvcrt.LK_UNLCK, 1)
            else:
                fcntl.flock(lock, fcntl.LOCK_UN)


def cosine(a, b):
    return sum(x * y for x, y in zip(a, b)) / max(1e-12, math.sqrt(sum(x*x for x in a) * sum(y*y for y in b)))


class Store:
    def __init__(self, data):
        data.mkdir(parents=True, exist_ok=True, mode=0o700)
        self.data = data
        self.db = sqlite3.connect(data / "faces.sqlite", timeout=30)
        self.db.row_factory = sqlite3.Row
        self.db.execute("PRAGMA journal_mode=WAL")
        self.db.execute("PRAGMA foreign_keys=ON")
        self.db.executescript("""
          CREATE TABLE IF NOT EXISTS people(id TEXT PRIMARY KEY, name TEXT NOT NULL, performer_id TEXT);
          CREATE TABLE IF NOT EXISTS images(id TEXT PRIMARY KEY, digest TEXT, model TEXT, error TEXT);
          CREATE TABLE IF NOT EXISTS galleries(image_id TEXT REFERENCES images(id) ON DELETE CASCADE, gallery_id TEXT,
            PRIMARY KEY(image_id, gallery_id));
          CREATE TABLE IF NOT EXISTS faces(id TEXT PRIMARY KEY, image_id TEXT REFERENCES images(id) ON DELETE CASCADE,
            person_id TEXT REFERENCES people(id), box TEXT NOT NULL, embedding TEXT NOT NULL);
          CREATE INDEX IF NOT EXISTS face_person ON faces(person_id);
          CREATE INDEX IF NOT EXISTS face_image ON faces(image_id);
          CREATE INDEX IF NOT EXISTS gallery_lookup ON galleries(gallery_id, image_id);
          CREATE TABLE IF NOT EXISTS state(key TEXT PRIMARY KEY, value TEXT NOT NULL);
        """)

    def state(self, value=None):
        if value is not None:
            with self.db:
                self.db.execute("INSERT OR REPLACE INTO state VALUES('scan', ?)", (json.dumps(value),))
        row = self.db.execute("SELECT value FROM state WHERE key='scan'").fetchone()
        return json.loads(row[0]) if row else {"phase": "idle"}

    def new_person(self):
        key = uuid.uuid4().hex
        number = self.db.execute("SELECT count(*) FROM people").fetchone()[0] + 1
        self.db.execute("INSERT INTO people VALUES(?,?,NULL)", (key, f"Person {number}"))
        return key

    def candidates(self):
        # Use a fixed exemplar per group, avoiding transitive clustering drift.
        return [(r["person_id"], json.loads(r["embedding"])) for r in self.db.execute(
            "SELECT person_id, embedding FROM faces WHERE rowid IN (SELECT min(rowid) FROM faces GROUP BY person_id)")]

    def save_image(self, image, digest, detections):
        old = self.db.execute("SELECT digest, model FROM images WHERE id=?", (image["id"],)).fetchone()
        with self.db:
            self.db.execute("INSERT INTO images VALUES(?,?,?,NULL) ON CONFLICT(id) DO UPDATE SET digest=excluded.digest, model=excluded.model, error=NULL",
                            (image["id"], digest, VERSION))
            self.db.execute("DELETE FROM galleries WHERE image_id=?", (image["id"],))
            self.db.executemany("INSERT INTO galleries VALUES(?,?)", [(image["id"], g["id"]) for g in image["galleries"]])
            if old and old["digest"] == digest and old["model"] == VERSION:
                return False
            self.db.execute("DELETE FROM faces WHERE image_id=?", (image["id"],))
            candidates = self.candidates()
            used = set()
            for box, embedding in detections:
                ranking = sorted(((cosine(embedding, ref), pid) for pid, ref in candidates if pid not in used), reverse=True)
                # Scores are similarity values, not probabilities or identity guarantees.
                if ranking and ranking[0][0] >= 0.55 and (len(ranking) < 2 or ranking[0][0] - ranking[1][0] >= 0.08):
                    person = ranking[0][1]
                else:
                    person = self.new_person()
                    candidates.append((person, embedding))
                used.add(person)
                self.db.execute("INSERT INTO faces VALUES(?,?,?,?,?)",
                    (uuid.uuid4().hex, image["id"], person, json.dumps(box), json.dumps(embedding)))
            self.db.execute("DELETE FROM people WHERE id NOT IN (SELECT person_id FROM faces)")
        return True

    def record_error(self, image, message):
        with self.db:
            self.db.execute("INSERT INTO images(id,error) VALUES(?,?) ON CONFLICT(id) DO UPDATE SET error=excluded.error", (image["id"], message))
            self.db.execute("DELETE FROM galleries WHERE image_id=?", (image["id"],))
            self.db.executemany("INSERT INTO galleries VALUES(?,?)", [(image["id"], g["id"]) for g in image["galleries"]])

    @staticmethod
    def scope(gallery):
        return (" AND f.image_id IN (SELECT image_id FROM galleries WHERE gallery_id=?)", [gallery]) if gallery else ("", [])

    def overview(self, gallery=None, page=1):
        scope, params = self.scope(gallery)
        base = " FROM people p JOIN faces f ON f.person_id=p.id WHERE 1=1" + scope
        count = self.db.execute("SELECT count(DISTINCT p.id)" + base, params).fetchone()[0]
        rows = self.db.execute("SELECT p.*,count(DISTINCT f.image_id) AS image_count,min(f.image_id) AS cover_id" + base +
            " GROUP BY p.id ORDER BY image_count DESC,p.id LIMIT 60 OFFSET ?", [*params, (page-1)*60]).fetchall()
        image_scope = " WHERE id IN (SELECT image_id FROM galleries WHERE gallery_id=?)" if gallery else ""
        stats = self.db.execute("SELECT count(*) AS scanned,sum(error IS NOT NULL) AS errors FROM images" + image_scope, params).fetchone()
        return {"people": [dict(r) for r in rows], "count": count, "page": page,
                "scanned": stats["scanned"], "errors": stats["errors"] or 0, "scan": self.state()}

    def faces(self, person, gallery=None, page=1):
        scope, params = self.scope(gallery)
        where = " FROM faces f WHERE f.person_id=?" + scope
        count = self.db.execute("SELECT count(*)" + where, [person, *params]).fetchone()[0]
        rows = self.db.execute("SELECT f.id,f.image_id,f.box" + where + " ORDER BY f.image_id,f.id LIMIT 60 OFFSET ?",
                               [person, *params, (page-1)*60]).fetchall()
        return {"count": count, "page": page, "faces": [{**dict(r), "box": json.loads(r["box"])} for r in rows]}

    def person(self, key):
        row = self.db.execute("SELECT * FROM people WHERE id=?", (key,)).fetchone()
        if not row:
            raise ValueError("This person group no longer exists. Refresh the page.")
        return dict(row)

    def rename(self, key, name):
        self.person(key)
        if not isinstance(name, str) or not name.strip() or len(name.strip()) > 160:
            raise ValueError("Enter a name between 1 and 160 characters")
        with self.db:
            self.db.execute("UPDATE people SET name=? WHERE id=?", (name.strip(), key))

    def split(self, face_ids):
        if not isinstance(face_ids, list) or not face_ids or len(face_ids) > 60:
            raise ValueError("Select between 1 and 60 faces")
        face_ids = list(dict.fromkeys(str(x) for x in face_ids))
        marks = ",".join("?" for _ in face_ids)
        with self.db:
            rows = self.db.execute(f"SELECT id FROM faces WHERE id IN ({marks})", face_ids).fetchall()
            if len(rows) != len(face_ids):
                raise ValueError("Some faces no longer exist. Refresh the page.")
            person = self.new_person()
            self.db.execute(f"UPDATE faces SET person_id=? WHERE id IN ({marks})", [person, *face_ids])
            self.db.execute("DELETE FROM people WHERE id NOT IN (SELECT person_id FROM faces)")
        return person

    def merge(self, source, target):
        if source == target:
            raise ValueError("Choose a different person group")
        a, b = self.person(source), self.person(target)
        if a["performer_id"] and a["performer_id"] != b["performer_id"]:
            raise ValueError("These groups have different performer assignments. Clear the source link before merging.")
        with self.db:
            self.db.execute("UPDATE faces SET person_id=? WHERE person_id=?", (target, source))
            self.db.execute("DELETE FROM people WHERE id=?", (source,))


class FaceEngine:
    def __init__(self, data):
        import cv2
        for model in MANIFEST["models"]:
            file = data / "models" / model["name"]
            if not file.is_file() or hashlib.sha256(file.read_bytes()).hexdigest() != model["sha256"]:
                raise RuntimeError("Local model missing or damaged. Run Gallery People setup again.")
        cv2.setNumThreads(2)
        self.cv = cv2
        self.detector = cv2.FaceDetectorYN.create(str(data / "models/detector.onnx"), "", (320,320), 0.9, 0.3, 5000)
        self.recognizer = cv2.FaceRecognizerSF.create(str(data / "models/recognizer.onnx"), "")

    def detect(self, content):
        import numpy as np
        from PIL import Image, ImageOps
        with Image.open(io.BytesIO(content)) as original:
            if original.width * original.height > 80_000_000:
                raise ValueError("Image exceeds the 80 megapixel processing limit")
            oriented = ImageOps.exif_transpose(original).convert("RGB")
            oriented.thumbnail((1600, 1600))
            image = self.cv.cvtColor(np.asarray(oriented), self.cv.COLOR_RGB2BGR)
        height, width = image.shape[:2]
        self.detector.setInputSize((width, height))
        _, faces = self.detector.detect(image)
        results = []
        for face in faces if faces is not None else []:
            if min(face[2], face[3]) < 24:
                continue
            aligned = self.recognizer.alignCrop(image, face)
            embedding = self.recognizer.feature(aligned).flatten()
            embedding /= max(float(np.linalg.norm(embedding)), 1e-12)
            results.append(([float(face[0]/width), float(face[1]/height), float(face[2]/width), float(face[3]/height)],
                            embedding.tolist()))
        return results


def scan(store, client, gallery=None):
    with exclusive_lock(store.data / "write.lock"):
        engine = FaceEngine(store.data)
        state = {"phase": "running", "done": 0, "total": 0, "failed": 0, "gallery_id": gallery}
        store.state(state)
        seen = set()
        try:
            page = 1
            while True:
                result = client.images(gallery, page=page)
                state["total"] = result["count"]
                if not result["images"]:
                    break
                for image in result["images"]:
                    seen.add(image["id"])
                    try:
                        content = client.image_bytes(image)
                        digest = hashlib.sha256(content).hexdigest()
                        old = store.db.execute("SELECT digest,model FROM images WHERE id=?", (image["id"],)).fetchone()
                        detections = [] if old and old["digest"] == digest and old["model"] == VERSION else engine.detect(content)
                        store.save_image(image, digest, detections)
                    except Exception as exc:
                        # Never log connection details, URLs containing signatures, or image bytes.
                        store.record_error(image, type(exc).__name__)
                        state["failed"] += 1
                    state["done"] += 1
                    store.state(state)
                    print(f"\x01p{state['done']/max(1,state['total']):.4f}", file=sys.stderr, flush=True)
                if len(result["images"]) < 100:
                    break
                page += 1
            # Reconcile deletions only after a complete, successful listing.
            with store.db:
                if gallery:
                    indexed = store.db.execute("SELECT image_id FROM galleries WHERE gallery_id=?", (gallery,)).fetchall()
                    store.db.executemany("DELETE FROM galleries WHERE gallery_id=? AND image_id=?",
                                         [(gallery, r[0]) for r in indexed if r[0] not in seen])
                else:
                    indexed = store.db.execute("SELECT id FROM images").fetchall()
                    store.db.executemany("DELETE FROM images WHERE id=?", [(r[0],) for r in indexed if r[0] not in seen])
                    store.db.execute("DELETE FROM people WHERE id NOT IN (SELECT person_id FROM faces)")
            state["phase"] = "finished"
        except Exception:
            state["phase"] = "failed"
            store.state(state)
            raise
        store.state(state)
        if state["failed"]:
            raise RuntimeError(f"Scan completed with {state['failed']} unreadable images. Successful images were saved; retry to process failures.")
        return state


def run(payload):
    args = payload.get("args") or {}
    data = Path(os.environ.get("GALLERY_PEOPLE_DATA") or str(Path(payload["server_connection"]["Dir"]) / "gallery-people"))
    store = Store(data)
    client = Stash(payload["server_connection"])
    op = args.get("operation", "overview")
    gallery = str(args["gallery_id"]) if args.get("gallery_id") else None
    page = max(1, min(int(args.get("page", 1)), 100000))
    try:
        if op == "scan":
            return scan(store, client, gallery)
        if op == "overview":
            result = store.overview(gallery, page)
            try:
                with exclusive_lock(data / "write.lock"):
                    if result["scan"]["phase"] == "running":
                        result["scan"]["phase"] = "interrupted"
                result["busy"] = False
            except RuntimeError:
                result["busy"] = True
            ids = list({p["cover_id"] for p in result["people"]})
            covers = client.images(ids=ids, per_page=60)["images"] if ids else []
            result["covers"] = {i["id"]: i for i in covers}
            return result
        if op == "faces":
            result = store.faces(str(args["person_id"]), gallery, page)
            result["person"] = store.person(str(args["person_id"]))
            ids = list({f["image_id"] for f in result["faces"]})
            result["images"] = client.images(ids=ids, per_page=60)["images"] if ids else []
            return result
        with exclusive_lock(data / "write.lock"):
            if op == "rename":
                store.rename(str(args["person_id"]), args.get("name"))
            elif op == "split":
                return {"person_id": store.split(args.get("face_ids"))}
            elif op == "merge":
                store.merge(str(args["person_id"]), str(args["target_id"]))
            elif op == "unlink":
                with store.db:
                    store.db.execute("UPDATE people SET performer_id=NULL WHERE id=?", (str(args["person_id"]),))
            elif op == "assign":
                if args.get("confirmed") is not True:
                    raise ValueError("Review the images and confirm before assigning a performer")
                person = store.person(str(args["person_id"]))
                performer = client.graphql("query($id: ID!) { findPerformer(id: $id) { id name } }", {"id": str(args["performer_id"])})["findPerformer"]
                if not performer:
                    raise ValueError("The selected performer no longer exists")
                if person["performer_id"] and person["performer_id"] != performer["id"]:
                    raise ValueError("Clear this group's existing performer link first")
                scope, params = store.scope(gallery)
                ids = [r[0] for r in store.db.execute("SELECT DISTINCT f.image_id FROM faces f WHERE f.person_id=?" + scope,
                                                     [person["id"], *params])]
                # ADD is idempotent and preserves all manually assigned performers.
                for start in range(0, len(ids), 100):
                    client.graphql("""mutation($input: BulkImageUpdateInput!) { bulkImageUpdate(input: $input) { id } }""",
                        {"input": {"ids": ids[start:start+100], "performer_ids": {"ids": [performer["id"]], "mode": "ADD"}}})
                with store.db:
                    store.db.execute("UPDATE people SET performer_id=?,name=? WHERE id=?", (performer["id"], performer["name"], person["id"]))
                return {"assigned": len(ids)}
            else:
                raise ValueError("Unknown Gallery People operation")
            return {"ok": True}
    finally:
        store.db.close()


if __name__ == "__main__":
    if len(sys.argv) == 3 and sys.argv[1] == "--check":
        FaceEngine(Path(sys.argv[2]))
        print("Face detector and recognizer loaded successfully (CPU).")
    else:
        try:
            print(json.dumps({"output": run(json.load(sys.stdin))}))
        except Exception as exc:
            print(json.dumps({"error": str(exc)}))
            sys.exit(1)
