#!/usr/bin/env python3
"""
db_backup.py — consistent, scheduled backup of stash's SQLite metadata DB + config.

The stash metadata DB (`stash-go.sqlite`) is the single irreplaceable asset: it holds
every scene/image/tag/performer relationship, scrape result, marker, and view-count.
The media files on disk can be re-scanned; this DB cannot be rebuilt. This tool makes a
CONSISTENT point-in-time copy of it — even while stash is running and writing — by using
SQLite's ONLINE BACKUP API (`sqlite3.Connection.backup`), NOT a plain file copy. A plain
`cp` of a live SQLite file can capture a torn write (a DB checkpoint mid-flush, or a
half-applied transaction with a separate `-wal`/`-shm`), yielding a backup that opens but
is silently corrupt. The online backup API copies page-by-page under SQLite's own locking,
so the result is always a valid, self-consistent database.

What it does:
  1. Snapshot the DB via the online backup API into a temp file.
  2. Copy `config.yml` alongside it.
  3. Bundle both into a timestamped gzip tar:  stash-backup-YYYYMMDD-HHMMSS.tar.gz
  4. Prune: keep the most recent --keep backups in --out, delete older ones.
  5. Optionally rsync the new tar to an off-NAS host (--remote).

Safety: DRY-RUN BY DEFAULT. It reports what it WOULD do (DB size, config presence, output
path, prune candidates) and writes NOTHING unless you pass --apply.

  ──────────────────────────────────────────────────────────────────────────────
  WHERE TO RUN: on the NAS host (overwatch-stash), where the DB file actually lives
  at /volume1/docker/stash/config/stash-go.sqlite. The stash-llm Docker container
  mounts that same config dir, so the host path is the source of truth. Reach the
  NAS with:  ssh -p 3239 shadowshark@overwatch-stash
  (shadowshark is not a sudoer and has no scp/SFTP — see --remote notes below.)
  ──────────────────────────────────────────────────────────────────────────────

  RECOMMENDED NAS CRON (nightly at 03:30, keep 14 days, write for real):
    30 3 * * * /usr/bin/python3 /volume1/docker/stash/tools/maintenance/db_backup.py --apply >> /volume1/docker/stash/backups/db_backup.log 2>&1

  On Synology DSM, prefer adding this as a "Scheduled Task" (user-defined script) in
  Control Panel → Task Scheduler rather than a raw crontab line, since DSM manages the
  user crontab and may rewrite hand-edited entries. The command body is identical.

Usage:
  # report only (default) — safe to run anywhere, writes nothing:
  python3 db_backup.py
  # actually take + bundle + prune a backup:
  python3 db_backup.py --apply
  # also push the new tar off-box:
  python3 db_backup.py --apply --remote charlesd@charles-lemurpro:/home/charlesd/stash-backups/
"""

import argparse
import os
import shutil
import sqlite3
import subprocess
import sys
import tarfile
import tempfile
import time

DEFAULT_DB = "/volume1/docker/stash/config/stash-go.sqlite"
DEFAULT_CONFIG = "/volume1/docker/stash/config/config.yml"
DEFAULT_OUT = "/volume1/docker/stash/backups"
DEFAULT_KEEP = 14
BACKUP_PREFIX = "stash-backup-"
BACKUP_SUFFIX = ".tar.gz"


def human(nbytes):
    """Render a byte count compactly (1023 -> '1023 B', 1.4M -> '1.4 MiB')."""
    if nbytes is None:
        return "?"
    n = float(nbytes)
    for unit in ("B", "KiB", "MiB", "GiB", "TiB"):
        if n < 1024.0 or unit == "TiB":
            return ("%d %s" % (int(n), unit)) if unit == "B" else ("%.1f %s" % (n, unit))
        n /= 1024.0
    return "%d B" % int(nbytes)


def snapshot_sqlite(src_path, dst_path):
    """Consistent online-backup copy of a (possibly live) SQLite DB.

    Opens the source read-only via URI so we never take a write lock or touch the WAL,
    then streams every page into a fresh DB at dst_path under SQLite's own locking.
    """
    src_uri = "file:%s?mode=ro" % src_path
    src = sqlite3.connect(src_uri, uri=True, timeout=30)
    try:
        dst = sqlite3.connect(dst_path)
        try:
            # pages=-1 copies the whole DB in one step; backup() handles locking/retries.
            src.backup(dst, pages=-1)
        finally:
            dst.close()
    finally:
        src.close()


def list_backups(out_dir):
    """Existing backup tars in out_dir, newest first, as (path, mtime, size)."""
    items = []
    if not os.path.isdir(out_dir):
        return items
    for name in os.listdir(out_dir):
        if name.startswith(BACKUP_PREFIX) and name.endswith(BACKUP_SUFFIX):
            p = os.path.join(out_dir, name)
            try:
                st = os.stat(p)
            except OSError:
                continue
            items.append((p, st.st_mtime, st.st_size))
    items.sort(key=lambda t: t[1], reverse=True)
    return items


def main():
    ap = argparse.ArgumentParser(
        description="Consistent backup of stash's SQLite metadata DB + config (dry-run by default)."
    )
    ap.add_argument("--db", default=DEFAULT_DB,
                    help="path to stash-go.sqlite (default: %(default)s)")
    ap.add_argument("--config", default=DEFAULT_CONFIG,
                    help="path to config.yml (default: %(default)s)")
    ap.add_argument("--out", default=DEFAULT_OUT,
                    help="output directory for backup tarballs (default: %(default)s)")
    ap.add_argument("--keep", type=int, default=DEFAULT_KEEP,
                    help="number of most-recent backups to retain (default: %(default)s)")
    ap.add_argument("--remote", default=None,
                    help="optional rsync destination for the new tar, e.g. "
                         "user@host:/path/ (scp/SFTP are unavailable on the NAS — "
                         "this shells out to rsync over ssh)")
    mode = ap.add_mutually_exclusive_group()
    mode.add_argument("--dry-run", dest="apply", action="store_false",
                      help="report only; write nothing (DEFAULT)")
    mode.add_argument("--apply", dest="apply", action="store_true",
                      help="actually take the backup, prune, and (if set) push --remote")
    ap.set_defaults(apply=False)
    args = ap.parse_args()

    mode_label = "APPLY" if args.apply else "DRY-RUN"
    print("=== stash db_backup (%s) ===" % mode_label)

    # ── source checks ─────────────────────────────────────────────────────────
    if not os.path.isfile(args.db):
        sys.exit("DB not found: %s  (run this on the NAS host where the DB lives)" % args.db)
    db_size = os.path.getsize(args.db)
    print("  DB:      %s  (%s)" % (args.db, human(db_size)))

    # WAL/SHM sidecars are informational — the online backup API folds them in for us.
    for side in ("-wal", "-shm"):
        sp = args.db + side
        if os.path.exists(sp):
            try:
                print("  sidecar: %s  (%s)" % (sp, human(os.path.getsize(sp))))
            except OSError:
                pass

    have_config = os.path.isfile(args.config)
    if have_config:
        print("  config:  %s  (%s)" % (args.config, human(os.path.getsize(args.config))))
    else:
        print("  config:  %s  (MISSING — will back up DB only)" % args.config)

    ts = time.strftime("%Y%m%d-%H%M%S")
    tar_name = "%s%s%s" % (BACKUP_PREFIX, ts, BACKUP_SUFFIX)
    tar_path = os.path.join(args.out, tar_name)
    print("  out:     %s" % tar_path)

    # ── retention preview ───────────────────────────────────────────────────────
    existing = list_backups(args.out)
    # In --apply we add one new tar, so the prune cutoff is computed against keep-1
    # of the existing set (the new one always survives). In dry-run we just show the
    # current over-limit tail.
    if args.apply:
        survivors = existing[: max(args.keep - 1, 0)]
        to_prune = existing[max(args.keep - 1, 0):]
    else:
        survivors = existing[: args.keep]
        to_prune = existing[args.keep:]
    print("  retain:  keep=%d  (currently %d backup(s) in out dir)"
          % (args.keep, len(existing)))
    if to_prune:
        print("  prune:   %d old backup(s) would be deleted:" % len(to_prune))
        for p, mt, sz in to_prune:
            print("             - %s  (%s, %s)"
                  % (os.path.basename(p), human(sz),
                     time.strftime("%Y-%m-%d %H:%M", time.localtime(mt))))
    else:
        print("  prune:   nothing to prune")

    if args.remote:
        print("  remote:  would rsync new tar -> %s" % args.remote)

    # ── dry-run stops here ──────────────────────────────────────────────────────
    if not args.apply:
        print("\n(dry-run — nothing written. Re-run with --apply to take the backup.)")
        return

    # ── apply: snapshot -> bundle -> prune -> remote ────────────────────────────
    os.makedirs(args.out, exist_ok=True)
    tmpdir = tempfile.mkdtemp(prefix="stash-backup-")
    try:
        snap_path = os.path.join(tmpdir, "stash-go.sqlite")
        print("\nsnapshotting DB via SQLite online-backup API …")
        snapshot_sqlite(args.db, snap_path)
        snap_size = os.path.getsize(snap_path)
        print("  snapshot OK (%s)" % human(snap_size))

        print("bundling -> %s" % tar_path)
        # Write to a .partial then rename, so a crash mid-write never leaves a
        # truncated tar that retention/restore might mistake for a good backup.
        partial = tar_path + ".partial"
        with tarfile.open(partial, "w:gz") as tar:
            tar.add(snap_path, arcname="stash-go.sqlite")
            if have_config:
                tar.add(args.config, arcname="config.yml")
        os.replace(partial, tar_path)
        print("  wrote %s (%s)" % (tar_name, human(os.path.getsize(tar_path))))
    finally:
        shutil.rmtree(tmpdir, ignore_errors=True)

    # prune AFTER the new tar exists and survives by definition
    if to_prune:
        print("pruning %d old backup(s) …" % len(to_prune))
        for p, _mt, _sz in to_prune:
            try:
                os.remove(p)
                print("  removed %s" % os.path.basename(p))
            except OSError as e:
                print("  ! could not remove %s: %s" % (os.path.basename(p), e))

    # optional off-box copy via rsync (scp/SFTP are not available on the NAS)
    if args.remote:
        print("pushing off-box via rsync -> %s" % args.remote)
        cmd = ["rsync", "-av", "--partial", tar_path, args.remote]
        try:
            rc = subprocess.call(cmd)
        except OSError as e:
            print("  ! rsync could not be launched: %s" % e)
            rc = 1
        if rc == 0:
            print("  rsync OK")
        else:
            print("  ! rsync exited %d — local backup is intact; off-box copy FAILED" % rc)

    print("\nsummary: backup %s (%s)%s%s"
          % (tar_name, human(os.path.getsize(tar_path)),
             ", pruned %d" % len(to_prune) if to_prune else "",
             ", pushed remote" if args.remote else ""))


if __name__ == "__main__":
    main()
