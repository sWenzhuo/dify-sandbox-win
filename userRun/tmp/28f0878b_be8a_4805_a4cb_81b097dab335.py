import ctypes
import os
import sys
import traceback
# setup sys.excepthook
def excepthook(type, value, tb):
    sys.stderr.write("".join(traceback.format_exception(type, value, tb)))
    sys.stderr.flush()
    sys.exit(-1)

sys.excepthook = excepthook

lib = ctypes.CDLL("./python.so")
lib.DifySeccomp.argtypes = [ctypes.c_uint32, ctypes.c_uint32, ctypes.c_bool]
lib.DifySeccomp.restype = None

# get running path
running_path = sys.argv[1]
if not running_path:
    exit(-1)

# get decrypt key
key = sys.argv[2]
if not key:
    exit(-1)

from base64 import b64decode
key = b64decode(key)

os.chdir(running_path)




lib.DifySeccomp(65537, 0, 0)

code = b64decode("5PU/IRbw1NJnT3pzOE1AgIi6N3tI7MTARQspKXGU45yj/eaN7XkZLh3Cps/gwem9AXVLHECUaGicG6W68Q/72pa/b2RY7dX6")

def decrypt(code, key):
    key_len = len(key)
    code_len = len(code)
    code = bytearray(code)
    for i in range(code_len):
        code[i] = code[i] ^ key[i % key_len]
    return bytes(code)

code = decrypt(code, key)
exec(code)