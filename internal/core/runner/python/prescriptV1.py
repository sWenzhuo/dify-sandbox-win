import os
import sys
import traceback


# setup sys.excepthook 全局异常处理
def excepthook(type, value, tb):
    sys.stderr.write("".join(traceback.format_exception(type, value, tb)))
    sys.stderr.flush()
    sys.exit(-1)

sys.excepthook = excepthook


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

#代码运行目录
os.chdir(running_path)

{{preload}}


# base64解密
code = b64decode("{{code}}")

def decrypt(code, key):
    key_len = len(key)
    code_len = len(code)
    code = bytearray(code)
    for i in range(code_len):
        code[i] = code[i] ^ key[i % key_len]
    return bytes(code)

# 异或解密
code = decrypt(code, key)
exec(code)