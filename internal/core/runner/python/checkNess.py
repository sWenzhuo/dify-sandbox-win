# check_sys_packages.py
import sys
import importlib

if len(sys.argv) < 2:
    print("请传入至少一个要检测的模块名，例如：python check_sys_packages.py numpy pandas", file=sys.stderr)
    sys.exit(1)

required_packages = sys.argv[1:]  # 从命令行参数获取包名列表
missing = []

for pkg in required_packages:
    try:
        importlib.import_module(pkg)
    except ImportError:
        missing.append(pkg)

if missing:
    print(f"缺失以下系统包：{missing}", file=sys.stderr)
    sys.exit(1)
else:
    print("所有系统包都已安装。")
