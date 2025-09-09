import os
import hashlib
from typing import Callable
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.requests import Request
from starlette.responses import JSONResponse
from ..utils.response import error, ERRCODE_COMMON_ERROR
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

APP_ID = os.getenv("APP_ID", "keti-3")
SKEY = os.getenv("SKEY", "dev-skey")
# Optional: allow timestamp drift in seconds
SIGN_TS_DRIFT = int(os.getenv("SIGN_TS_DRIFT", "0"))


class SignMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next: Callable):
        # Health or static can bypass
        if request.url.path in ("/health",):
            return await call_next(request)

        app_id = request.headers.get("X-Auth-AppId", "")
        ts = request.headers.get("X-Auth-TimeStamp", "")
        sign = request.headers.get("X-Sign", "")

        # Basic checks
        if not app_id or not ts or not sign:
            return JSONResponse(error(ERRCODE_COMMON_ERROR, "missing signature headers"))
        if app_id != APP_ID:
            return JSONResponse(error(ERRCODE_COMMON_ERROR, "app id not match"))

        # Optional drift check
        if SIGN_TS_DRIFT:
            try:
                import time
                ts_i = int(ts)
                if abs(int(time.time()) - ts_i) > SIGN_TS_DRIFT:
                    return JSONResponse(error(ERRCODE_COMMON_ERROR, "timestamp drift too large"))
            except Exception:
                return JSONResponse(error(ERRCODE_COMMON_ERROR, "invalid timestamp"))

        expected = hashlib.md5(f"{SKEY}{ts}".encode("utf-8")).hexdigest()
        if expected != sign:
            return JSONResponse(error(ERRCODE_COMMON_ERROR, "signature not match"))

        return await call_next(request)
