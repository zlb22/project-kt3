import os
from dotenv import load_dotenv
from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from .middleware.sign import SignMiddleware
from .middleware.auth import AuthMiddleware
from .api.keti3 import router as keti3_router

load_dotenv()

APP_TITLE = os.getenv("APP_TITLE", "keti3 Python API")
APP_VERSION = os.getenv("APP_VERSION", "0.1.0")

app = FastAPI(title=APP_TITLE, version=APP_VERSION)

# CORS (adjust as needed)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"]
)

# Middlewares (order matters): signature first, then auth
app.add_middleware(SignMiddleware)
app.add_middleware(AuthMiddleware)

# Routers
app.include_router(keti3_router, prefix="/web/keti3")


@app.get("/health")
async def health():
    return {"status": "ok"}
