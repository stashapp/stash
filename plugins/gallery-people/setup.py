"""Install a private runtime and checksum-verified public models (no library access)."""
import argparse
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys
import urllib.request
import venv

ROOT = Path(__file__).resolve().parent
COMMIT = "48b10fa505aeb51fbd127fc60a5c660040e823b7"


def download_models(data):
    model_dir = data / "models"
    model_dir.mkdir(parents=True, exist_ok=True)
    manifest = json.loads((ROOT / "models.json").read_text())
    for model in manifest["models"]:
        dest = model_dir / model["name"]
        if dest.is_file() and hashlib.sha256(dest.read_bytes()).hexdigest() == model["sha256"]:
            continue
        url = f"https://media.githubusercontent.com/media/opencv/opencv_zoo/{COMMIT}/{model['path']}"
        temp = dest.with_suffix(".download")
        try:
            print(f"Downloading {model['name']}…", flush=True)
            digest = hashlib.sha256()
            with urllib.request.urlopen(url, timeout=120) as response, temp.open("wb") as output:
                while chunk := response.read(1024 * 1024):
                    digest.update(chunk)
                    output.write(chunk)
            if digest.hexdigest() != model["sha256"]:
                raise RuntimeError(f"Checksum failed for {model['name']}; the model was not installed")
            os.replace(temp, dest)
        finally:
            temp.unlink(missing_ok=True)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--data-dir", required=True, type=Path)
    args = parser.parse_args()
    data = args.data_dir.expanduser().resolve()
    data.mkdir(parents=True, exist_ok=True, mode=0o700)
    python = data / "venv" / ("Scripts/python.exe" if os.name == "nt" else "bin/python")
    if not python.is_file():
        venv.EnvBuilder(with_pip=True).create(data / "venv")
    subprocess.run([str(python), "-m", "pip", "install", "--disable-pip-version-check",
                    "-r", str(ROOT / "requirements.txt")], check=True)
    download_models(data)
    subprocess.run([str(python), str(ROOT / "plugin.py"), "--check", str(data)], check=True)
    print("Gallery People is installed. Image processing will work offline.")


if __name__ == "__main__":
    main()
