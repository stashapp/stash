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
    if isinstance(levelChar, bytes):
        try:
            levelChar = levelChar.decode('ascii', errors='ignore')
        except Exception:
            levelChar = ''
    startLevelChar = '\x01'
    endLevelChar = '\x02'

    ret = startLevelChar + levelChar + endLevelChar
    return ret

def __log(levelChar, s):
    if not levelChar:
        return

    print(__prefix(levelChar) + str(s), file=sys.stderr, flush=True)

def LogTrace(s):
    __log('t', s)

def LogDebug(s):
    __log('d', s)

def LogInfo(s):
    __log('i', s)

def LogWarning(s):
    __log('w', s)

def LogError(s):
    __log('e', s)

def LogProgress(p):
    progress = min(max(0, p), 1)
    __log('p', str(progress))