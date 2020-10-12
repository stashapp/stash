import sys

# Log messages sent from a plugin instance are transmitted via stderr and are
# encoded with a prefix consisting of special character SOH, then the log
# level (one of t, d, i, w, e, or p - corresponding to trace, debug, info,
# warning, error and progress levels respectively), then special character
# STX.
#
# The LogTrace, LogDebug, LogInfo, LogWarning, and LogError methods, and their equivalent
# formatted methods are intended for use by plugin instances to transmit log
# messages. The LogProgress method is also intended for sending progress data.
#

def __prefix(levelChar):
    startLevelChar = b'\x01'
    endLevelChar = b'\x02'

    ret = startLevelChar + levelChar + endLevelChar
    return ret.decode()

def __log(levelChar, s):
    if levelChar == "":
        return

    print(__prefix(levelChar) + s + "\n", file=sys.stderr, flush=True)

def LogTrace(s):
    __log(b't', s)

def LogDebug(s):
    __log(b'd', s)

def LogInfo(s):
    __log(b'i', s)

def LogWarning(s):
    __log(b'w', s)

def LogError(s):
    __log(b'e', s)

def LogProgress(p):
    progress = min(max(0, p), 1)
    __log(b'p', str(progress))
