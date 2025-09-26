# Import statements first
from fastapi import FastAPI, HTTPException, Depends, Request, UploadFile, File, Form, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import FileResponse, RedirectResponse
from pydantic import BaseModel
import jwt
from jwt import PyJWTError
import hashlib
import bcrypt
from datetime import datetime, timedelta
from typing import Optional
from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.primitives import serialization, hashes
import base64
import io
import random
import secrets
import traceback
from typing import Optional, List
import uvicorn
import os
from dotenv import load_dotenv
import logging
from urllib.parse import urlparse
from PIL import Image, ImageDraw, ImageFont, ImageFilter
import re
import json
from minio import Minio
from minio.error import S3Error
from sqlalchemy.orm import Session
from db_mysql import get_db, Student, OperationLog, TwentyFourRecord, VideoAsset, AuthLock, create_tables
from keti3_middleware import verify_signature, keti3_response, error, ERRCODE_COMMON_ERROR, ERRCODE_INVALID_PARAMS
from keti3_storage import get_minio_client, ensure_bucket, presign_put_object

# Load environment variables
load_dotenv()

# Configuration constants
CAPTCHA_EXPIRE_SECONDS = int(os.getenv("CAPTCHA_EXPIRE_SECONDS", "120"))
LOCK_THRESHOLD = int(os.getenv("AUTH_LOCK_THRESHOLD", "5"))
LOCK_DURATION_MINUTES = int(os.getenv("AUTH_LOCK_DURATION_MINUTES", "15"))

# In-memory captcha store: id -> { code, expires }
_captcha_store: dict[str, dict] = {}

def _captcha_prune_now():
    now = datetime.utcnow().timestamp()
    expired = [k for k, v in _captcha_store.items() if v.get("expires_ts", 0) <= now]
    for k in expired:
        try:
            del _captcha_store[k]
        except Exception:
            pass

def _gen_captcha_text(length: int = 5) -> str:
    alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
    return ''.join(secrets.choice(alphabet) for _ in range(length))

def _render_captcha_image(text: str, width: int = 140, height: int = 50) -> bytes:
    img = Image.new('RGB', (width, height), (255, 255, 255))
    draw = ImageDraw.Draw(img)
    # Try to load a common font; fallback to default
    try:
        font = ImageFont.truetype("DejaVuSans-Bold.ttf", 28)
    except Exception:
        font = ImageFont.load_default()
    # Add noise lines
    for _ in range(8):
        x1, y1 = random.randint(0, width), random.randint(0, height)
        x2, y2 = random.randint(0, width), random.randint(0, height)
        draw.line(((x1, y1), (x2, y2)), fill=(random.randint(150,200), random.randint(150,200), random.randint(150,200)), width=1)
    # Draw text with slight jitter per char
    x = 10
    for ch in text:
        y = random.randint(5, 15)
        draw.text((x, y), ch, fill=(random.randint(0,80), random.randint(0,80), random.randint(0,80)), font=font)
        x += 24
    # Light blur
    try:
        img = img.filter(ImageFilter.SMOOTH)
    except Exception:
        pass
    buf = io.BytesIO()
    img.save(buf, format='PNG')
    return buf.getvalue()

def _validate_password_strength(password: str) -> tuple[bool, str]:
    """Validate password meets security requirements: 8+ chars, uppercase, lowercase, numbers, special chars"""
    if len(password) < 8:
        return False, "密码长度至少8位"
    
    if not re.search(r'[A-Z]', password):
        return False, "密码必须包含大写字母"
    
    if not re.search(r'[a-z]', password):
        return False, "密码必须包含小写字母"
    
    if not re.search(r'[0-9]', password):
        return False, "密码必须包含数字"
    
    if not re.search(r'[!@#$%^&*()_+\-=\[\]{};:\'"\\|,.<>\/?~`]', password):
        return False, "密码必须包含特殊字符（如!@#$%^&*等）"
    
    return True, ""

def _validate_and_consume_captcha(captcha_id: str, captcha_code: str) -> bool:
    _captcha_prune_now()
    if not captcha_id or not captcha_code:
        return False
    rec = _captcha_store.get(captcha_id)
    if not rec:
        return False
    # One-time use
    try:
        del _captcha_store[captcha_id]
    except Exception:
        pass
    if rec.get("expires_ts", 0) < datetime.utcnow().timestamp():
        return False
    return str(captcha_code).strip().lower() == rec.get("code", "")

# JWT Configuration
SECRET_KEY = "your-secret-key-here"  # Change this in production
ALGORITHM = "HS256"
ACCESS_TOKEN_EXPIRE_MINUTES = 30

# Generate RSA key pair for password encryption
private_key = rsa.generate_private_key(
    public_exponent=65537,
    key_size=2048,
)
public_key = private_key.public_key()

# Serialize public key for transmission
public_key_pem = public_key.public_bytes(
    encoding=serialization.Encoding.PEM,
    format=serialization.PublicFormat.SubjectPublicKeyInfo
).decode('utf-8')

# MongoDB Configuration (legacy; not used for /api/auth/* anymore)
MONGODB_URL = os.getenv("MONGODB_URL", "mongodb://localhost:27017")
DATABASE_NAME = os.getenv("DATABASE_NAME", "project-kt3")

# MinIO Configuration
MINIO_URL = os.getenv("MINIO_URL", "http://localhost:9000")
MINIO_ACCESS_KEY = os.getenv("MINIO_ACCESS_KEY", "minioadmin")
MINIO_SECRET_KEY = os.getenv("MINIO_SECRET_KEY", "minioadmin")
MINIO_BUCKET_NAME = os.getenv("MINIO_BUCKET_NAME", "onlineclass")
# Dedicated bucket for Keti3 "在线实验" sub-project (uploads, assets)
MINIO_BUCKET_KETI3 = os.getenv("MINIO_BUCKET_KETI3", os.getenv("MINIO_BUCKET", "online-experiment"))
# 可选：限制允许的 bucket，逗号分隔；为空则不限制
MINIO_ALLOWED_BUCKETS = [b.strip() for b in os.getenv("MINIO_ALLOWED_BUCKETS", "").split(",") if b.strip()]

# Optional: enable append-only JSON logging for uploads (for debugging/audit). Default: disabled
UPLOADS_JSON_LOG_ENABLED = os.getenv("UPLOADS_JSON_LOG_ENABLED", "false").lower() in ("1", "true", "yes")

app = FastAPI(title="Educational Assessment API", version="2.0")
security = HTTPBearer()

# Logger for auth debugging
logger = logging.getLogger("auth")
logger.setLevel(logging.INFO)

# Initialize MySQL tables
create_tables()

# Removed legacy MongoDB client and collections; switched to MySQL-only metadata storage

cors_origins = os.getenv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173").split(",")

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=cors_origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Add security headers middleware
@app.middleware("http")
async def add_security_headers(request: Request, call_next):
    response = await call_next(request)
    # Force HTTPS and prevent downgrade attacks
    response.headers["Strict-Transport-Security"] = "max-age=31536000; includeSubDomains"
    # Prevent clickjacking
    response.headers["X-Frame-Options"] = "DENY"
    # Prevent MIME type sniffing
    response.headers["X-Content-Type-Options"] = "nosniff"
    # XSS protection
    response.headers["X-XSS-Protection"] = "1; mode=block"
    return response

# HTTP to HTTPS redirect middleware
@app.middleware("http")
async def force_https(request: Request, call_next):
    if request.headers.get("x-forwarded-proto") == "http":
        url = str(request.url).replace("http://", "https://", 1)
        return RedirectResponse(url=url, status_code=301)
    return await call_next(request)

@app.get("/api/auth/captcha")
async def get_captcha():
    """Generate a new CAPTCHA and return { captcha_id, image_base64 }"""
    _captcha_prune_now()
    text = _gen_captcha_text()
    cid = secrets.token_urlsafe(16)
    img_bytes = _render_captcha_image(text)
    b64 = base64.b64encode(img_bytes).decode('ascii')
    _captcha_store[cid] = {
        "code": text.lower(),
        "expires_ts": (datetime.utcnow().timestamp() + CAPTCHA_EXPIRE_SECONDS),
    }
    return {"captcha_id": cid, "image_base64": f"data:image/png;base64,{b64}", "expires_in": CAPTCHA_EXPIRE_SECONDS}

@app.get("/api/auth/public-key")
async def get_public_key():
    """Get RSA public key for password encryption"""
    return {"public_key": public_key_pem}

# Pydantic models

# Pydantic models
class UserLogin(BaseModel):
    username: str
    password: str
    captcha_id: str
    captcha_code: str

class UserCreate(BaseModel):
    username: str
    password: str

class UserRegister(BaseModel):
    username: str
    school: str
    student_id: str
    grade: str
    password: str

class PasswordChange(BaseModel):
    old_password: str
    new_password: str

class Token(BaseModel):
    access_token: str
    token_type: str

class User(BaseModel):
    username: str
    school: str
    student_id: str
    grade: str
    is_active: bool = True
    created_at: Optional[datetime] = None

# Keti3 API Models
class StudentLoginReq(BaseModel):
    username: str
    school: str
    grade: str

class OssAuthReq(BaseModel):
    uid: int | None = None
    content_types: dict[str, str] | None = None

class LogSaveReq(BaseModel):
    uid: int
    action: str
    details: Optional[str] = None

# Keti3 specific log save payload compatible with frontend
class Keti3LogSaveReq(BaseModel):
    uid: int
    op_type: Optional[str] = None
    voice_url: Optional[str] = None
    voice_text: Optional[str] = None
    screenshot_url: Optional[str] = None
    op_time: Optional[str] = None
    data_after: Optional[dict | list | str] = None

class Save24ptRecordReq(BaseModel):
    uid: int
    payload: Optional[dict | list | str] = None
    payload_str: Optional[str] = None

# Video upload (presign/commit) models
class VideoPresignReq(BaseModel):
    video_type: str
    test_session_id: str
    bucket: Optional[str] = None
    content_type: Optional[str] = None

class VideoCommitReq(BaseModel):
    bucket: str
    object_name: str
    test_session_id: str
    content_type: Optional[str] = None
    file_size: Optional[int] = None
    video_type: Optional[str] = None

# Utility functions
def make_hashed_password(password: str) -> str:
    """Hash password using bcrypt with salt"""
    salt = bcrypt.gensalt()
    hashed = bcrypt.hashpw(password.encode('utf-8'), salt)
    return hashed.decode('utf-8')

# MySQL user operations (replace legacy Mongo for /api/auth/*)
def get_student_by_username(db: Session, username: str) -> Student | None:
    return db.query(Student).filter(Student.username == username).first()

def create_student_mysql(db: Session, username: str, school: str, student_id: str, grade: str, password_hash: str) -> int | None:
    if get_student_by_username(db, username):
        return None
    stu = Student(
        username=username.strip(),
        school=school.strip(),
        grade=grade.strip(),
        password=password_hash,
        is_active=1,
        created_at=datetime.utcnow(),
        updated_at=datetime.utcnow(),
    )
    db.add(stu)
    db.commit()
    db.refresh(stu)
    return stu.id

def update_student_password_mysql(db: Session, username: str, new_password_hash: str):
    stu = get_student_by_username(db, username)
    if not stu:
        return False
    stu.password = new_password_hash
    stu.updated_at = datetime.utcnow()
    db.commit()
    return True

# Removed legacy MongoDB password update helper

def decrypt_password(encrypted_password: str) -> str:
    """
    Decrypt RSA encrypted password from frontend
    """
    try:
        # Decode base64 encrypted password
        encrypted_bytes = base64.b64decode(encrypted_password)
        
        # Decrypt using private key
        decrypted_bytes = private_key.decrypt(
            encrypted_bytes,
            padding.OAEP(
                mgf=padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )
        
        return decrypted_bytes.decode('utf-8')
    except Exception as e:
        raise HTTPException(status_code=400, detail="Password decryption failed")

def check_password(encrypted_password: str, stored_password_hash: str) -> bool:
    """
    Verify encrypted password against stored hash
    Frontend sends RSA encrypted password, backend decrypts and compares
    """
    try:
        # Decrypt password first
        password_from_frontend = decrypt_password(encrypted_password)
        try:
            logger.info("auth.check_password: decrypted password length=%d, stored_prefix=%s", len(password_from_frontend), stored_password_hash[:4])
        except Exception:
            pass
        
        # If stored hash is bcrypt format ($2a$, $2b$, $2y$)
        if stored_password_hash.startswith(('$2a$', '$2b$', '$2y$')):
            # 1) Try bcrypt against plaintext password (current correct behavior)
            if bcrypt.checkpw(password_from_frontend.encode('utf-8'), stored_password_hash.encode('utf-8')):
                logger.info("auth.check_password: bcrypt verification success (plaintext)")
                return True
            # 2) Backward compatibility: some historical accounts may have stored bcrypt(hash_sha256(password))
            legacy_sha = hashlib.sha256(password_from_frontend.encode('utf-8')).hexdigest()
            ok = bcrypt.checkpw(legacy_sha.encode('utf-8'), stored_password_hash.encode('utf-8'))
            logger.info("auth.check_password: bcrypt(sha256) fallback result=%s", ok)
            return ok
        else:
            # Legacy SHA256 hash
            ok = hashlib.sha256(str.encode(password_from_frontend)).hexdigest() == stored_password_hash
            logger.info("auth.check_password: legacy sha256 compare result=%s", ok)
            return ok
    except Exception:
        return False

def create_access_token(data: dict, expires_delta: Optional[timedelta] = None):
    to_encode = data.copy()
    if expires_delta:
        expire = datetime.utcnow() + expires_delta
    else:
        expire = datetime.utcnow() + timedelta(minutes=15)
    to_encode.update({"exp": expire})
    encoded_jwt = jwt.encode(to_encode, SECRET_KEY, algorithm=ALGORITHM)
    return encoded_jwt

def verify_token(credentials: HTTPAuthorizationCredentials = Depends(security)):
    try:
        payload = jwt.decode(credentials.credentials, SECRET_KEY, algorithms=[ALGORITHM])
        username: str = payload.get("sub")
        if username is None:
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Could not validate credentials",
                headers={"WWW-Authenticate": "Bearer"},
            )
        return username
    except jwt.PyJWTError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Could not validate credentials",
            headers={"WWW-Authenticate": "Bearer"},
        )

# Adapter entry: upsert student and redirect to sub_project with JWT token
KETI3_FRONTEND_URL = os.getenv("KETI3_FRONTEND_URL", "https://localhost:5173/topic-three")

def get_keti3_frontend_url(request: Request) -> str:
    """
    Auto-detect frontend URL based on environment and request
    """
    # Check if we have explicit config
    if KETI3_FRONTEND_URL and not KETI3_FRONTEND_URL.startswith("auto://"):
        return KETI3_FRONTEND_URL
    
    # Auto-detect based on request origin
    origin = request.headers.get("origin") or f"{request.url.scheme}://{request.url.hostname}"
    if request.url.port and request.url.port not in [80, 443]:
        origin = f"{origin}:{request.url.port}"
    
    # Development: use sub-frontend dev server
    if "localhost" in origin or "127.0.0.1" in origin or ":3000" in origin:
        return f"{request.url.scheme}://{request.url.hostname}:5174/topic-three/online-experiment"
    
    # Production: assume Nginx proxy at /topic-three/online-experiment
    return f"{origin}/topic-three/online-experiment"

@app.get("/web/keti3/entry")
async def keti3_entry(request: Request, current_user: str = Depends(verify_token), db: Session = Depends(get_db)):
    """
    Secure entry point - requires authentication, no URL parameters
    Redirects to sub_project with JWT token
    """
    # Get user info from token (already authenticated)
    student = db.query(Student).filter(Student.username == current_user).first()
    if not student:
        raise HTTPException(status_code=404, detail="User not found")

    # Issue JWT for sub-frontend
    access_token_expires = timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
    token = create_access_token(data={"sub": current_user}, expires_delta=access_token_expires)

    # Auto-detect frontend URL
    frontend_url = get_keti3_frontend_url(request)
    redirect_url = f"{frontend_url}?token={token}"
    
    return RedirectResponse(url=redirect_url, status_code=302)

@app.post("/api/auth/create-sub-token")
async def create_sub_token(current_user: str = Depends(verify_token), db: Session = Depends(get_db)):
    """
    Create a JWT token for sub-frontend access
    """
    # Verify user exists
    student = db.query(Student).filter(Student.username == current_user).first()
    if not student:
        raise HTTPException(status_code=404, detail="User not found")

    # Issue JWT for sub-frontend
    access_token_expires = timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
    token = create_access_token(data={"sub": current_user}, expires_delta=access_token_expires)
    
    return {"token": token}

# Authentication endpoints
@app.post("/api/auth/login", response_model=Token)
async def login(user_data: UserLogin, db: Session = Depends(get_db)):
    """Login via MySQL students table with CAPTCHA and lockout"""
    # 1) Lockout check
    lock = db.query(AuthLock).filter(AuthLock.username == user_data.username).first()
    now_dt = datetime.utcnow()
    if lock and lock.locked_until and lock.locked_until > now_dt:
        remaining = int((lock.locked_until - now_dt).total_seconds() // 60) + 1
        raise HTTPException(status_code=status.HTTP_429_TOO_MANY_REQUESTS, detail=f"账号已锁定，请稍后再试（约{remaining}分钟）")

    # 2) CAPTCHA validation
    if not _validate_and_consume_captcha(user_data.captcha_id, user_data.captcha_code):
        # Count as a failure
        if not lock:
            lock = AuthLock(username=user_data.username, failed_count=1, locked_until=None)
            db.add(lock)
        else:
            lock.failed_count = (lock.failed_count or 0) + 1
        if lock.failed_count >= LOCK_THRESHOLD:
            lock.locked_until = now_dt + timedelta(minutes=LOCK_DURATION_MINUTES)
        db.commit()
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="验证码错误或过期")

    # 3) Password validation
    stu = get_student_by_username(db, user_data.username)
    if not stu or not stu.password or not check_password(user_data.password, stu.password):
        if not lock:
            lock = AuthLock(username=user_data.username, failed_count=1, locked_until=None)
            db.add(lock)
        else:
            lock.failed_count = (lock.failed_count or 0) + 1
        if lock.failed_count >= LOCK_THRESHOLD:
            lock.locked_until = now_dt + timedelta(minutes=LOCK_DURATION_MINUTES)
        db.commit()
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="用户名或密码错误")

    # 4) Success -> reset lock
    if lock:
        lock.failed_count = 0
        lock.locked_until = None
        db.commit()

    access_token_expires = timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
    access_token = create_access_token(data={"sub": user_data.username}, expires_delta=access_token_expires)
    return {"access_token": access_token, "token_type": "bearer"}

@app.post("/api/auth/register")
async def register_user(user_data: UserRegister, request: Request, db: Session = Depends(get_db)):
    """用户注册"""
    
    # 直接返回注册失败，不进行任何实际操作
    return {
        "success": False,
        "message": "注册失败"
    }

@app.get("/api/auth/me", response_model=User)
async def read_users_me(current_user: str = Depends(verify_token), db: Session = Depends(get_db)):
    stu = get_student_by_username(db, current_user)
    if not stu:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="用户不存在")
    return User(
        username=stu.username,
        school=stu.school,
        student_id="",
        grade=stu.grade,
        is_active=bool(getattr(stu, "is_active", 1)),
        created_at=getattr(stu, "created_at", None),
    )

@app.post("/api/auth/change-password")
async def change_password(
    password_data: PasswordChange,
    current_user: str = Depends(verify_token),
    db: Session = Depends(get_db),
):
    stu = get_student_by_username(db, current_user)
    if not stu or not stu.password or not check_password(password_data.old_password, stu.password):
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="Incorrect old password")
    new_password_hash = make_hashed_password(password_data.new_password)
    ok = update_student_password_mysql(db, current_user, new_password_hash)
    if not ok:
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail="Failed to update password")
    return {"message": "Password changed successfully"}

# Health check
@app.get("/api/health")
async def health_check():
    return {"status": "healthy", "timestamp": datetime.utcnow()}

# Test endpoints (placeholders for your existing functionality)
@app.get("/api/tests/available")
async def get_available_tests(current_user: str = Depends(verify_token)):
    """Get list of available tests"""
    return {
        "tests": [
            {"id": "problem_solving", "name": "问题解决测评", "description": "在线实验场景下问题解决能力测评"},
            {"id": "interactive_discussion", "name": "互动讨论", "description": "互动讨论场景下智能化测评"},
            {"id": "aut_test", "name": "AUT测试", "description": "问题解决、团队协作能力智能化测评"},
            {"id": "emotion_test", "name": "情绪测评", "description": "学习与抗挫情绪智能化测评"},
            {"id": "online_class", "name": "在线课程", "description": "在线微课程场景下智能化测评"}
        ]
    }

@app.get("/api/tests/{test_id}")
async def get_test_details(test_id: str, current_user: str = Depends(verify_token)):
    """Get details for a specific test"""
    # This would integrate with your existing test logic
    return {"test_id": test_id, "status": "available", "description": f"Details for {test_id}"}

# 创建上传目录
UPLOAD_DIR = os.getenv("UPLOAD_DIR", "uploads")
os.makedirs(UPLOAD_DIR, exist_ok=True)
os.makedirs(f"{UPLOAD_DIR}/videos", exist_ok=True)

@app.post("/api/upload/video")
async def upload_video(
    file: UploadFile = File(...),
    video_type: str = Form(...),  # "camera" or "screen"（自动纠正拼写）
    test_session_id: str = Form(...),
    bucket: Optional[str] = Form(None),
    current_user: str = Depends(verify_token),
    db: Session = Depends(get_db),
):
    """Upload recorded video files to MinIO (bucket/<camera|screen>/filename)."""
    try:
        # 验证文件类型
        if not file.content_type.startswith('video/'):
            raise HTTPException(status_code=400, detail="只允许上传视频文件")

        # 规范化 video_type（纠正常见拼写）
        vt = (video_type or "").strip().lower()
        folder_map = {
            "camera": "camera",
            "cam": "camera",
            "camaer": "camera",
            "screen": "screen",
            "scrren": "screen",
        }
        folder = folder_map.get(vt)
        if folder is None:
            raise HTTPException(status_code=400, detail="video_type 只能是 camera 或 screen")

        # 生成文件名（使用纠正后的 folder 名称）
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        filename = f"{current_user}_{test_session_id}_{folder}_{timestamp}.webm"
        # 先解析并写入 MySQL 元数据需要的 uid，因此在命名前先定位学生
        student = get_student_by_username(db, current_user)
        if not student:
            raise HTTPException(status_code=404, detail="当前用户未找到")
        # 每个 uid 单独文件夹，结构：<camera|screen>/<uid>/<test_session_id>/<filename>
        object_name = f"{folder}/{student.id}/{test_session_id}/{filename}"

        # 流式方式：直接使用底层临时文件对象，避免整文件读入内存
        try:
            raw = file.file  # SpooledTemporaryFile or file-like
            raw.seek(0, os.SEEK_END)
            file_size = raw.tell()
            raw.seek(0)
            if file_size <= 0:
                raise HTTPException(status_code=400, detail="上传文件为空")
        except Exception as e:
            raise HTTPException(status_code=400, detail=f"无法读取上传文件: {e}")

        # 解析目标 bucket（表单优先，其次环境变量）
        target_bucket = (bucket or MINIO_BUCKET_NAME).strip()
        if not target_bucket:
            raise HTTPException(status_code=400, detail="bucket 不能为空")
        if "/" in target_bucket or "\\" in target_bucket:
            raise HTTPException(status_code=400, detail="bucket 名称不合法")
        # 可选白名单校验
        if MINIO_ALLOWED_BUCKETS and target_bucket not in MINIO_ALLOWED_BUCKETS:
            raise HTTPException(status_code=403, detail=f"不允许的 bucket: {target_bucket}")

        # 初始化 MinIO 客户端（支持自签名/配置化 SSL）
        client = get_minio_client()

        # 确保 bucket 存在
        try:
            ensure_bucket(client, target_bucket)
        except Exception as e:
            raise HTTPException(status_code=500, detail=f"检查/创建 bucket 失败: {e}")

        # 上传到 MinIO（使用文件对象与已知长度，避免高内存占用）
        try:
            client.put_object(
                target_bucket,
                object_name,
                data=raw,
                length=file_size,
                content_type=file.content_type,
            )
        except S3Error as e:
            raise HTTPException(status_code=500, detail=f"上传到 MinIO 失败: {e}")

        # 记录上传信息（沿用原结构，file_path 写成 MinIO 路径）
        upload_info = {
            "user": current_user,
            "test_session_id": test_session_id,
            "video_type": folder,
            "filename": filename,
            "file_size": file_size,
            "upload_time": datetime.now().isoformat(),
            "bucket": target_bucket,
            "object_name": object_name,
            "content_type": file.content_type,
            "file_path": f"minio://{target_bucket}/{object_name}",
        }

        # 解析并写入 MySQL 元数据（按 uid 索引）
        public_base = f"{MINIO_URL}/{target_bucket}"
        public_url = f"{public_base}/{object_name}"
        asset = VideoAsset(
            uid=student.id,
            username=current_user,
            test_session_id=test_session_id,
            video_type=folder,
            bucket=target_bucket,
            object_name=object_name,
            content_type=file.content_type,
            file_size=file_size,
        )
        db.add(asset)
        db.commit()
        db.refresh(asset)

        # 可选：保存上传记录到 JSON 文件（仅用于调试/审计）
        if UPLOADS_JSON_LOG_ENABLED:
            uploads_log_file = os.path.join(UPLOAD_DIR, "uploads_log.json")
            uploads_log = []
            if os.path.exists(uploads_log_file):
                try:
                    with open(uploads_log_file, 'r', encoding='utf-8') as f:
                        uploads_log = json.load(f)
                except Exception:
                    # 如果旧日志损坏，则从空开始
                    uploads_log = []

            uploads_log.append(upload_info)

            with open(uploads_log_file, 'w', encoding='utf-8') as f:
                json.dump(uploads_log, f, ensure_ascii=False, indent=2)

        return {
            "message": "视频上传成功",
            "video_type": upload_info["video_type"],
            "filename": filename,
            "file_size": file_size,
            "content_type": file.content_type,
            "bucket": target_bucket,
            "object_name": object_name,
            "file_path": upload_info["file_path"],
            "test_session_id": test_session_id,
            "upload_time": upload_info["upload_time"],
            "uid": student.id,
            "id": asset.id,
            "public_url": public_url,
        }

    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"文件上传失败: {str(e)}")

# Removed legacy MongoDB-based 24pt submit endpoint; use MySQL-based endpoints below instead

@app.post("/api/tests/24point/submit")
async def submit_24point_test(
    test_data: dict,
    current_user: str = Depends(verify_token),
    db: Session = Depends(get_db),
):
    """Submit 24-point test results and persist to MySQL (twentyfour_records) with video meta references.
    Compatible with old frontend that posts to /api/tests/24point/submit.
    """
    try:
        # Assemble document
        doc = dict(test_data or {})
        doc["user"] = current_user
        doc.setdefault("test_id", "24point")
        doc["submit_time"] = datetime.now().isoformat()

        # 绑定 uid
        stu = get_student_by_username(db, current_user)
        if not stu:
            raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="用户不存在")

        # 合并视频元数据：按 uid + test_session_id 查询 VideoAsset，写入 doc['videos']
        tsid = doc.get("test_session_id") or doc.get("session_id")
        videos_list = []
        if tsid:
            q = (
                db.query(VideoAsset)
                .filter(VideoAsset.uid == stu.id, VideoAsset.test_session_id == tsid)
                .order_by(VideoAsset.upload_time.asc())
            )
            for a in q.all():
                public_url = f"{MINIO_URL}/{a.bucket}/{a.object_name}"
                videos_list.append({
                    "type": a.video_type,
                    "bucket": a.bucket,
                    "object_name": a.object_name,
                    "public_url": public_url,
                    "file_size": a.file_size,
                    "content_type": a.content_type,
                    "upload_time": a.upload_time.isoformat() if a.upload_time else None,
                })
        doc["videos"] = videos_list

        payload_text = json.dumps(doc, ensure_ascii=False)
        rec = TwentyFourRecord(uid=stu.id, payload=payload_text)
        db.add(rec)
        db.commit()
        db.refresh(rec)

        return {
            "message": "测试结果提交成功",
            "id": rec.id,
            "submit_time": doc["submit_time"],
            "user": current_user,
            "test_id": doc.get("test_id", "24point"),
        }
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"测试结果提交失败: {str(e)}")
# Keti3 API Endpoints
@app.options("/web/keti3/student/login")
async def keti3_student_login_options():
    """Handle OPTIONS for student login"""
    return {}

@app.post("/web/keti3/student/login")
async def keti3_student_login(request: Request, req: StudentLoginReq, db: Session = Depends(get_db)):
    """Keti3 student login endpoint"""
    try:
        verify_signature(request)
    except HTTPException as e:
        # For debugging, log the signature error but continue
        print(f"Signature verification failed: {e.detail}")
        # Comment out the next line to bypass signature verification during debugging
        # raise e
    
    try:
        # Check if student exists
        student = db.query(Student).filter(Student.username == req.username).first()
        
        if not student:
            # Create new student
            student = Student(
                username=req.username,
                school=req.school,
                grade=req.grade
            )
            db.add(student)
            db.commit()
            db.refresh(student)
        
        return keti3_response(data={
            "uid": student.id,
            "username": student.username,
            "school": student.school,
            "grade": student.grade
        })
    except Exception as e:
        return error(ERRCODE_COMMON_ERROR, f"Login failed: {str(e)}")

@app.options("/web/keti3/config/list")
async def keti3_config_list_options():
    """Handle OPTIONS for config list"""
    return {}

@app.get("/web/keti3/config/list")
@app.post("/web/keti3/config/list")
async def keti3_config_list(request: Request):
    """Get configuration list for keti3"""
    verify_signature(request)
    
    try:
        # Return static configuration data
        config_data = {
            "school": ["北京小学", "上海小学", "广州小学", "深圳小学"],
            "grade": ["一年级", "二年级", "三年级", "四年级", "五年级", "六年级"]
        }
        return keti3_response(data=config_data)
    except Exception as e:
        return error(ERRCODE_COMMON_ERROR, f"Config list failed: {str(e)}")

@app.options("/web/keti3/oss/auth")
async def keti3_oss_auth_options():
    """Handle OPTIONS for oss auth"""
    return {}

@app.post("/web/keti3/oss/auth")
async def keti3_oss_auth(request: Request, req: OssAuthReq | None = None):
    """Get OSS authentication for file uploads. Body is optional."""
    verify_signature(request)
    
    try:
        # Tolerate missing body or different shapes
        uid = 0
        img_ct = "image/*"
        audio_ct = "audio/*"
        if req is not None:
            uid = req.uid or 0
            if req.content_types:
                img_ct = req.content_types.get("img", img_ct)
                audio_ct = req.content_types.get("audio", audio_ct)

        client = get_minio_client()
        bucket = os.getenv("MINIO_BUCKET_KETI3", MINIO_BUCKET_KETI3)
        ensure_bucket(client, bucket)
        base_path = os.getenv("MINIO_BASE_PATH", "24game")
        
        now = datetime.utcnow()
        ts = int(now.timestamp())
        
        img_key = f"{base_path}/{uid}/{ts}/img{uid}{ts}"
        audio_key = f"{base_path}/{uid}/{ts}/audio{uid}{ts}"
        
        expires = timedelta(hours=1)
        img_put = presign_put_object(client, bucket, img_key, content_type=img_ct, expires=expires)
        audio_put = presign_put_object(client, bucket, audio_key, content_type=audio_ct, expires=expires)
        
        public_base = f"{MINIO_URL}/{bucket}"
        
        resp = {
            "storage": "minio",
            "Bucket": bucket,
            "HttpsDomain": public_base,
            "Presigned": [
                {"key": img_key, "url": img_put, "method": "PUT", "content_type": img_ct, "public_url": f"{public_base}/{img_key}"},
                {"key": audio_key, "url": audio_put, "method": "PUT", "content_type": audio_ct, "public_url": f"{public_base}/{audio_key}"},
            ],
            "AccessKeyId": "",
            "AccessKeySecret": "",
            "SecurityToken": "",
            "Expiration": (now + expires).isoformat() + "Z",
        }
        return keti3_response(data=resp)
    except Exception as e:
        import traceback
        print(f"OSS Auth Error: {e}")
        print(f"Traceback: {traceback.format_exc()}")
        return error(ERRCODE_COMMON_ERROR, f"MinIO error: {str(e)}")

@app.post("/web/keti3/oss/upload")
async def keti3_oss_upload(
    request: Request,
    img: UploadFile = File(None),
    audio: UploadFile = File(None),
    uid: Optional[int] = Form(None),
    db: Session = Depends(get_db),
):
    """Proxy upload endpoint to avoid browser-to-MinIO CORS. Accepts multipart form-data with fields 'img' and 'audio'."""
    verify_signature(request)

    if img is None and audio is None:
        return error(ERRCODE_INVALID_PARAMS, "No files provided. Expect 'img' and/or 'audio' fields.")

    try:
        client = get_minio_client()
        bucket = os.getenv("MINIO_BUCKET_KETI3", "online-experiment")
        ensure_bucket(client, bucket)
        base_path = os.getenv("MINIO_BASE_PATH_KETI3", "creative_work")

        now = datetime.utcnow()
        ts = int(now.timestamp())

        # Resolve uid in a robust way: form field -> header -> token username -> 0
        resolved_uid = 0
        # 1) form field uid
        if uid is not None:
            try:
                resolved_uid = int(uid)
            except Exception:
                resolved_uid = 0
        # 2) header X-User-Id
        if resolved_uid <= 0:
            try:
                resolved_uid = int(request.headers.get("X-User-Id", "0"))
            except Exception:
                resolved_uid = 0
        # 3) Authorization token -> username -> students.id
        if resolved_uid <= 0:
            try:
                auth = request.headers.get("Authorization", "")
                if auth.startswith("Bearer "):
                    token = auth.split(" ", 1)[1]
                    payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
                    username = payload.get("sub")
                    if username:
                        stu = get_student_by_username(db, username)
                        if stu:
                            resolved_uid = int(stu.id)
            except Exception:
                resolved_uid = 0
        # Debug print minimal info
        try:
            print(f"[oss/upload] resolved_uid={resolved_uid}")
        except Exception:
            pass

        public_base = f"{MINIO_URL}/{bucket}"
        result: dict[str, str] = {}

        if img is not None:
            img_key = f"{base_path}/{resolved_uid}/{ts}/img{resolved_uid}{ts}"
            try:
                raw = img.file
                raw.seek(0, os.SEEK_END)
                img_size = raw.tell()
                raw.seek(0)
                if img_size <= 0:
                    return error(ERRCODE_INVALID_PARAMS, "img is empty")
                client.put_object(bucket, img_key, data=raw, length=img_size, content_type=img.content_type or "image/*")
            finally:
                try:
                    await img.close()
                except Exception:
                    pass
            result["imgUrl"] = f"{public_base}/{img_key}"

        if audio is not None:
            audio_key = f"{base_path}/{resolved_uid}/{ts}/audio{resolved_uid}{ts}"
            try:
                raw = audio.file
                raw.seek(0, os.SEEK_END)
                audio_size = raw.tell()
                raw.seek(0)
                if audio_size <= 0:
                    return error(ERRCODE_INVALID_PARAMS, "audio is empty")
                client.put_object(bucket, audio_key, data=raw, length=audio_size, content_type=audio.content_type or "audio/*")
            finally:
                try:
                    await audio.close()
                except Exception:
                    pass
            result["audioUrl"] = f"{public_base}/{audio_key}"

        return keti3_response(data=result)
    except Exception as e:
        import traceback
        print(f"OSS Proxy Upload Error: {e}")
        print(f"Traceback: {traceback.format_exc()}")
        return error(ERRCODE_COMMON_ERROR, f"Upload failed: {str(e)}")

@app.options("/web/keti3/oss/upload")
async def keti3_oss_upload_options():
    """Handle OPTIONS for oss upload"""
    return {}

# Presigned video upload (recommended path)
@app.post("/api/videos/presign")
async def presign_video_upload(req: VideoPresignReq, current_user: str = Depends(verify_token), db: Session = Depends(get_db)):
    """Generate a presigned PUT URL for direct video upload to MinIO, and return planned object_name and public_url.
    Frontend should PUT the file to url, then call /api/videos/commit to register metadata.
    """
    try:
        # Normalize type
        vt = (req.video_type or "").strip().lower()
        folder_map = {
            "camera": "camera",
            "cam": "camera",
            "camaer": "camera",
            "screen": "screen",
            "scrren": "screen",
        }
        folder = folder_map.get(vt)
        if folder is None:
            raise HTTPException(status_code=400, detail="video_type 只能是 camera 或 screen")

        # Resolve student uid
        stu = get_student_by_username(db, current_user)
        if not stu:
            raise HTTPException(status_code=404, detail="当前用户未找到")

        # Target bucket
        target_bucket = (req.bucket or MINIO_BUCKET_NAME).strip()
        if not target_bucket:
            raise HTTPException(status_code=400, detail="bucket 不能为空")
        if "/" in target_bucket or "\\" in target_bucket:
            raise HTTPException(status_code=400, detail="bucket 名称不合法")
        if MINIO_ALLOWED_BUCKETS and target_bucket not in MINIO_ALLOWED_BUCKETS:
            raise HTTPException(status_code=403, detail=f"不允许的 bucket: {target_bucket}")

        # Build object name
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        # default content_type
        ct = (req.content_type or "video/webm").strip()
        ext = ".webm" if "webm" in ct else ".mp4" if "mp4" in ct else ".webm"
        filename = f"{current_user}_{req.test_session_id}_{folder}_{timestamp}{ext}"
        object_name = f"{folder}/{stu.id}/{req.test_session_id}/{filename}"

        client = get_minio_client()
        ensure_bucket(client, target_bucket)
        url = presign_put_object(client, target_bucket, object_name, content_type=ct, expires=timedelta(hours=1))
        public_base = f"{MINIO_URL}/{target_bucket}"
        return {
            "bucket": target_bucket,
            "object_name": object_name,
            "put_url": url,
            "public_url": f"{public_base}/{object_name}",
            "content_type": ct,
            "uid": stu.id,
            "username": current_user,
        }
    except HTTPException:
        raise
    except Exception as e:
        # Print detailed traceback to server log for diagnostics
        print(f"[presign] Error: {e}")
        print(traceback.format_exc())
        raise HTTPException(status_code=500, detail=f"presign failed: {e}")


@app.post("/api/videos/commit")
async def commit_video_upload(req: VideoCommitReq, current_user: str = Depends(verify_token), db: Session = Depends(get_db)):
    """Register an uploaded video object into MySQL video_assets. Frontend should call this after a successful direct upload."""
    stu = get_student_by_username(db, current_user)
    if not stu:
        raise HTTPException(status_code=404, detail="当前用户未找到")

    # derive video_type if not provided
    vtype = (req.video_type or "").strip().lower()
    if not vtype:
        try:
            vtype = req.object_name.split("/", 1)[0]
        except Exception:
            vtype = ""
    if vtype not in ("camera", "screen"):
        vtype = "camera"

    content_type = (req.content_type or "video/webm").strip()
    file_size = int(req.file_size or 0)

    asset = VideoAsset(
        uid=stu.id,
        username=current_user,
        test_session_id=req.test_session_id,
        video_type=vtype,
        bucket=req.bucket,
        object_name=req.object_name,
        content_type=content_type,
        file_size=file_size,
    )
    db.add(asset)
    db.commit()
    db.refresh(asset)

    public_base = f"{MINIO_URL}/{req.bucket}"
    return {
        "message": "登记成功",
        "id": asset.id,
        "uid": stu.id,
        "bucket": req.bucket,
        "object_name": req.object_name,
        "public_url": f"{public_base}/{req.object_name}",
        "content_type": content_type,
        "file_size": file_size,
        "test_session_id": req.test_session_id,
        "video_type": vtype,
    }

@app.options("/web/keti3/log/save")
async def keti3_log_save_options():
    return {}

@app.post("/web/keti3/log/save")
async def keti3_log_save(request: Request, req: Keti3LogSaveReq, db: Session = Depends(get_db)):
    """Save operation log compatible with frontend payload"""
    verify_signature(request)
    try:
        details_obj = {
            "voice_url": req.voice_url,
            "voice_text": req.voice_text,
            "screenshot_url": req.screenshot_url,
            "op_time": req.op_time,
            "data_after": req.data_after,
        }
        # 仅保留有值的键
        details_obj = {k: v for k, v in details_obj.items() if v is not None}
        details_str = json.dumps(details_obj, ensure_ascii=False) if details_obj else None

        log_entry = OperationLog(
            uid=req.uid,
            action=(req.op_type or "submit"),
            details=details_str
        )
        db.add(log_entry)
        db.commit()
        return keti3_response(data={"message": "Log saved successfully"})
    except Exception as e:
        return error(ERRCODE_COMMON_ERROR, f"Log save failed: {str(e)}")

# 24-point big JSON storage endpoints (MySQL)
@app.post("/api/24pt/record/save")
async def save_24pt_record(request: Request, req: Save24ptRecordReq, db: Session = Depends(get_db)):
    """Persist a 24-point record as a big JSON/text payload bound to uid."""
    verify_signature(request)
    try:
        if req.payload is None and (req.payload_str is None or req.payload_str.strip() == ""):
            return error(ERRCODE_INVALID_PARAMS, "payload is required")
        payload_text = req.payload_str if req.payload_str not in (None, "") else json.dumps(req.payload, ensure_ascii=False)
        rec = TwentyFourRecord(uid=req.uid, payload=payload_text)
        db.add(rec)
        db.commit()
        db.refresh(rec)
        return keti3_response(data={
            "id": rec.id,
            "uid": rec.uid,
            "created_at": rec.created_at.isoformat() if rec.created_at else None
        })
    except Exception as e:
        return error(ERRCODE_COMMON_ERROR, f"Save 24pt record failed: {str(e)}")

@app.get("/api/24pt/record/list")
async def list_24pt_records(request: Request, uid: int, limit: int = 20, offset: int = 0, db: Session = Depends(get_db)):
    """List 24-point records for a user, newest first. Payload returned as parsed JSON if possible, otherwise string."""
    verify_signature(request)
    try:
        limit = max(1, min(100, int(limit)))
        offset = max(0, int(offset))
        q = db.query(TwentyFourRecord).filter(TwentyFourRecord.uid == uid).order_by(TwentyFourRecord.created_at.desc()).offset(offset).limit(limit)
        items = []
        for rec in q.all():
            payload_value = rec.payload
            parsed = None
            try:
                parsed = json.loads(payload_value)
            except Exception:
                parsed = payload_value
            items.append({
                "id": rec.id,
                "uid": rec.uid,
                "created_at": rec.created_at.isoformat() if rec.created_at else None,
                "payload": parsed,
            })
        return keti3_response(data={"items": items, "limit": limit, "offset": offset})
    except Exception as e:
        return error(ERRCODE_COMMON_ERROR, f"List 24pt records failed: {str(e)}")

if __name__ == "__main__":
    host = os.getenv("HOST", "0.0.0.0")
    port = int(os.getenv("PORT", "8000"))
    uvicorn.run(app, host=host, port=port)
