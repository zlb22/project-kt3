def success(data=None, errmsg="ok", trace=""):
    return {
        "errcode": 0,
        "errmsg": errmsg,
        "data": data or {},
        "trace": trace,
    }


def error(errcode: int, errmsg: str, trace: str = ""):
    return {
        "errcode": errcode,
        "errmsg": errmsg,
        "data": {},
        "trace": trace,
    }


# Common error codes to align with frontend handling
ERRCODE_USER_NOT_LOGIN = 20004
ERRCODE_NO_PERMISSION = 20005
ERRCODE_COMMON_ERROR = 60000
