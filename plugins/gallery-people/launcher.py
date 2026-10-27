"""Start the plugin with its private Python environment; never log plugin input."""
import json
import os
from pathlib import Path
import site
import sys


def main():
    raw = sys.stdin.buffer.read()
    try:
        payload = json.loads(raw)
        data = Path(os.environ.get("GALLERY_PEOPLE_DATA") or
                    str(Path(payload["server_connection"]["Dir"]) / "gallery-people"))
        python = data / "venv" / ("Scripts/python.exe" if os.name == "nt" else "bin/python")
        if not python.is_file():
            raise RuntimeError("Gallery People needs its local models. Run the included setup command (see START-HERE.md), then try again.")
        version = f"python{sys.version_info.major}.{sys.version_info.minor}"
        packages = data / "venv" / ("Lib/site-packages" if os.name == "nt" else f"lib/{version}/site-packages")
        if not packages.is_dir():
            raise RuntimeError("Python changed since setup. Re-run Gallery People setup with the Python used by Stash.")
        site.addsitedir(str(packages))
        from plugin import run
        print(json.dumps({"output": run(payload)}))
        return 0
    except Exception as exc:
        print(json.dumps({"error": str(exc)}))
        return 1


if __name__ == "__main__":
    sys.exit(main())
