import os
import jwt
import re
import asyncio
from typing import Callable
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.requests import Request
from starlette.responses import JSONResponse
from ..utils.response import error, ERRCODE_USER_NOT_LOGIN
try:
    import redis  # optional
except Exception:  # pragma: no cover
    redis = None

JWT_SECRET = os.getenv("SECRET_KEY", "your-secret-key-change-in-production")
JWT_ALG = os.getenv("ALGORITHM", "HS256")
REDIS_URL = os.getenv("REDIS_URL", "")  # empty means disabled

# Whitelist paths that can bypass auth (still require signature)
AUTH_WHITELIST = {
    "/health",
    "/web/keti3/student/login",
    "/web/keti3/config/list",
    "/web/keti3/oss/auth",
}

# Create a global redis client (blocking client is fine in middleware)
_redis_client = None

def get_redis():
    global _redis_client
    if not REDIS_URL or redis is None:
        return None
    if _redis_client is None:
        _redis_client = redis.from_url(REDIS_URL)
    return _redis_client


class AuthMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next: Callable):
        path = request.url.path
        if path in AUTH_WHITELIST:
            return await call_next(request)

        token = request.headers.get("Authorization", "").strip()
        if not token:
            return JSONResponse(error(ERRCODE_USER_NOT_LOGIN, "请登陆"))

        # Try JWT first (contains two dots or Bearer prefix)
        if token.startswith("Bearer "):
            token_candidate = token[len("Bearer "):]
        else:
            token_candidate = token

        if token_candidate.count(".") == 2:
            try:
                payload = jwt.decode(token_candidate, JWT_SECRET, algorithms=[JWT_ALG])
                # attach to request state
                request.state.user = payload
                return await call_next(request)
            except Exception:
                # fall back to redis token
                pass

        # Fallback: legacy teacher token in redis (optional)
        r = get_redis()
        if r is not None:
            try:
                key = f"base-keti:teacher-login:{token}"
                val = r.get(key)
                if val:
                    request.state.user = {"role": "teacher", "token": token}
                    return await call_next(request)
            except Exception:
                # treat as not logged in on redis errors
                pass

        return JSONResponse(error(ERRCODE_USER_NOT_LOGIN, "请登陆"))
