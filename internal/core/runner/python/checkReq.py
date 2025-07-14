# checkReq.py
import sys
import os
import importlib

if len(sys.argv) != 2:
    print("Usage: python checkReq.py <requirements.txt path>")
    sys.exit(1)

requirements_path = sys.argv[1]

if not os.path.exists(requirements_path):
    print(f"File not found: {requirements_path}")
    sys.exit(1)

with open(requirements_path, "r", encoding="utf-8") as f:
    lines = f.readlines()

missing = []

for line in lines:
    line = line.strip()
    if not line or line.startswith("#"):
        continue
    # 去除版本号等，如 "pandas==2.1.0" -> "pandas"
    pkg = line.split("==")[0].split(">=")[0].split("<=")[0].strip()

    try:
        importlib.import_module(pkg)
    except ImportError:
        missing.append(pkg)

if missing:
    print("Missing or incompatible packages from requirements:", missing)
    sys.exit(1)

print("All required packages are installed.")
sys.exit(0)
