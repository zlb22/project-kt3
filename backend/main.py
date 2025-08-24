from fastapi import FastAPI, HTTPException, Depends, status, UploadFile, File, Form
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import FileResponse
from pydantic import BaseModel
import jwt
import hashlib
import json
import os
from datetime import datetime, timedelta
from typing import Optional, List
import uvicorn
from motor.motor_asyncio import AsyncIOMotorClient
from dotenv import load_dotenv
from urllib.parse import urlparse
from minio import Minio
from minio.error import S3Error
import io

# Load environment variables
load_dotenv()

# JWT Configuration
SECRET_KEY = os.getenv("SECRET_KEY", "your-secret-key-change-in-production")
ALGORITHM = os.getenv("ALGORITHM", "HS256")
ACCESS_TOKEN_EXPIRE_MINUTES = int(os.getenv("ACCESS_TOKEN_EXPIRE_MINUTES", "30"))

# MongoDB Configuration
MONGODB_URL = os.getenv("MONGODB_URL", "mongodb://localhost:27017")
DATABASE_NAME = os.getenv("DATABASE_NAME", "project-kt3")

# MinIO Configuration
MINIO_URL = os.getenv("MINIO_URL", "http://localhost:9000")
MINIO_ACCESS_KEY = os.getenv("MINIO_ACCESS_KEY")
MINIO_SECRET_KEY = os.getenv("MINIO_SECRET_KEY")
MINIO_BUCKET_NAME = os.getenv("MINIO_BUCKET_NAME", "onlineclass")
# 可选：限制允许的 bucket，逗号分隔；为空则不限制
MINIO_ALLOWED_BUCKETS = [b.strip() for b in os.getenv("MINIO_ALLOWED_BUCKETS", "").split(",") if b.strip()]

app = FastAPI(title="Educational Assessment API", version="2.0")
security = HTTPBearer()

# MongoDB client
client = AsyncIOMotorClient(MONGODB_URL)
db = client[DATABASE_NAME]
users_collection = db.users

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=os.getenv("CORS_ORIGINS", "http://localhost:3000").split(","),
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Pydantic models
class UserLogin(BaseModel):
    username: str
    password: str

class UserCreate(BaseModel):
    username: str
    password: str

class PasswordChange(BaseModel):
    old_password: str
    new_password: str

class Token(BaseModel):
    access_token: str
    token_type: str

class User(BaseModel):
    username: str
    is_active: bool = True

# Utility functions
def make_hashed_password(password: str) -> str:
    return hashlib.sha256(str.encode(password)).hexdigest()

# MongoDB user operations
async def get_user_from_db(username: str):
    """从MongoDB获取用户信息"""
    user = await users_collection.find_one({"username": username})
    return user

async def update_user_password(username: str, new_password_hash: str):
    """更新用户密码"""
    await users_collection.update_one(
        {"username": username},
        {"$set": {"password": new_password_hash}}
    )

def check_password(password: str, hashed_password: str) -> bool:
    return make_hashed_password(password) == hashed_password

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

# Authentication endpoints
@app.post("/api/auth/login", response_model=Token)
async def login(user_data: UserLogin):
    # 从MongoDB获取用户信息
    user = await get_user_from_db(user_data.username)
    if not user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect username or password"
        )
    
    if not check_password(user_data.password, user["password"]):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect username or password"
        )
    
    access_token_expires = timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
    access_token = create_access_token(
        data={"sub": user_data.username}, expires_delta=access_token_expires
    )
    return {"access_token": access_token, "token_type": "bearer"}

@app.get("/api/auth/me", response_model=User)
async def read_users_me(current_user: str = Depends(verify_token)):
    return User(username=current_user)

@app.post("/api/auth/change-password")
async def change_password(
    password_data: PasswordChange,
    current_user: str = Depends(verify_token)
):
    # 从MongoDB获取当前用户信息
    user = await get_user_from_db(current_user)
    if not user:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="User not found"
        )
    
    if not check_password(password_data.old_password, user["password"]):
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Incorrect old password"
        )
    
    # 更新密码到MongoDB
    new_password_hash = make_hashed_password(password_data.new_password)
    await update_user_password(current_user, new_password_hash)
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
    current_user: str = Depends(verify_token)
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
        object_name = f"{folder}/{filename}"

        # 读取内容
        content = await file.read()

        # 解析目标 bucket（表单优先，其次环境变量）
        target_bucket = (bucket or MINIO_BUCKET_NAME).strip()
        if not target_bucket:
            raise HTTPException(status_code=400, detail="bucket 不能为空")
        if "/" in target_bucket or "\\" in target_bucket:
            raise HTTPException(status_code=400, detail="bucket 名称不合法")
        # 可选白名单校验
        if MINIO_ALLOWED_BUCKETS and target_bucket not in MINIO_ALLOWED_BUCKETS:
            raise HTTPException(status_code=403, detail=f"不允许的 bucket: {target_bucket}")

        # 初始化 MinIO 客户端
        parsed = urlparse(MINIO_URL)
        if not parsed.scheme or not parsed.netloc:
            raise HTTPException(status_code=500, detail="MINIO_URL 配置无效")
        secure = parsed.scheme == "https"
        endpoint = parsed.netloc
        client = Minio(endpoint, access_key=MINIO_ACCESS_KEY, secret_key=MINIO_SECRET_KEY, secure=secure)

        # 确保 bucket 存在
        try:
            if not client.bucket_exists(target_bucket):
                client.make_bucket(target_bucket)
        except S3Error as e:
            raise HTTPException(status_code=500, detail=f"检查/创建 bucket 失败: {e}")

        # 上传到 MinIO
        try:
            client.put_object(
                target_bucket,
                object_name,
                data=io.BytesIO(content),
                length=len(content),
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
            "file_size": len(content),
            "upload_time": datetime.now().isoformat(),
            "bucket": target_bucket,
            "object_name": object_name,
            "content_type": file.content_type,
            "file_path": f"minio://{target_bucket}/{object_name}",
        }

        # 保存上传记录到JSON文件
        uploads_log_file = os.path.join(UPLOAD_DIR, "uploads_log.json")
        uploads_log = []
        if os.path.exists(uploads_log_file):
            with open(uploads_log_file, 'r', encoding='utf-8') as f:
                uploads_log = json.load(f)

        uploads_log.append(upload_info)

        with open(uploads_log_file, 'w', encoding='utf-8') as f:
            json.dump(uploads_log, f, ensure_ascii=False, indent=2)

        return {
            "message": "视频上传成功",
            "video_type": upload_info["video_type"],
            "filename": filename,
            "file_size": len(content),
            "content_type": file.content_type,
            "bucket": target_bucket,
            "object_name": object_name,
            "file_path": upload_info["file_path"],
            "test_session_id": test_session_id,
            "upload_time": upload_info["upload_time"],
        }

    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"文件上传失败: {str(e)}")

@app.post("/api/tests/24point/submit")
async def submit_24point_test(
    test_data: dict,
    current_user: str = Depends(verify_token)
):
    """Submit 24-point test results and persist to MongoDB with video meta references."""
    try:
        # 构造文档
        doc = dict(test_data or {})
        doc["user"] = current_user
        doc.setdefault("test_id", "24point")
        doc["submit_time"] = datetime.now().isoformat()

        # videos 字段如果存在，确保是列表，元素推荐包含：
        # {type, bucket, object_name, file_path, file_size, content_type}
        vids = doc.get("videos")
        if vids is not None and not isinstance(vids, list):
            doc["videos"] = [vids]

        # 入库到 test_results 集合
        result = await db.onlineclass.insert_one(doc)

        return {
            "message": "测试结果提交成功",
            "id": str(result.inserted_id),
            "submit_time": doc["submit_time"],
            "user": current_user,
            "test_id": doc.get("test_id", "24point"),
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"测试结果提交失败: {str(e)}")

if __name__ == "__main__":
    host = os.getenv("HOST", "0.0.0.0")
    port = int(os.getenv("PORT", "8000"))
    uvicorn.run(app, host=host, port=port)
