from fastapi import FastAPI, HTTPException, Depends, Request, UploadFile, File, Form
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import FileResponse, RedirectResponse
from pydantic import BaseModel
import jwt
import hashlib
import json
import os
from datetime import datetime, timedelta
import io
from typing import Optional, List
import uvicorn
from dotenv import load_dotenv
from urllib.parse import urlparse
from minio import Minio
from minio.error import S3Error
import io
from sqlalchemy.orm import Session
from db_mysql import get_db, Student, OperationLog, TwentyFourRecord, create_tables
from keti3_middleware import verify_signature, keti3_response, error, ERRCODE_COMMON_ERROR, ERRCODE_INVALID_PARAMS
from keti3_storage import get_minio_client, ensure_bucket, presign_put_object

# Load environment variables
load_dotenv()

# JWT Configuration
SECRET_KEY = os.getenv("SECRET_KEY", "your-secret-key-change-in-production")
ALGORITHM = os.getenv("ALGORITHM", "HS256")
ACCESS_TOKEN_EXPIRE_MINUTES = int(os.getenv("ACCESS_TOKEN_EXPIRE_MINUTES", "30"))

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

app = FastAPI(title="Educational Assessment API", version="2.0")
security = HTTPBearer()

# Initialize MySQL tables
create_tables()

# Legacy MongoDB client (kept for backward compatibility in other modules if needed)
try:
    from motor.motor_asyncio import AsyncIOMotorClient
    client = AsyncIOMotorClient(MONGODB_URL)
    db = client[DATABASE_NAME]
    users_collection = db.users
except Exception:
    client = None
    db = None
    users_collection = None

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=os.getenv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173").split(","),
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
    screenshot_url: Optional[str] = None
    op_time: Optional[str] = None
    data_after: Optional[dict | list | str] = None

class Save24ptRecordReq(BaseModel):
    uid: int
    payload: Optional[dict | list | str] = None
    payload_str: Optional[str] = None

# Utility functions
def make_hashed_password(password: str) -> str:
    return hashlib.sha256(str.encode(password)).hexdigest()

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

async def update_user_password(username: str, new_password_hash: str):
    """更新用户密码"""
    await users_collection.update_one(
        {"username": username},
        {"$set": {"password": new_password_hash, "updated_at": datetime.utcnow()}}
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

# Adapter entry: upsert student and redirect to sub_project with JWT token
KETI3_FRONTEND_URL = os.getenv("KETI3_FRONTEND_URL", "https://localhost:5173/topic-three")

@app.get("/web/keti3/entry")
async def keti3_entry(username: str, school: str = "", grade: str = "", db: Session = Depends(get_db)):
    """Upsert student into MySQL and redirect to sub_project page with JWT token in query."""
    # Upsert student by username
    student = db.query(Student).filter(Student.username == username).first()
    if not student:
        student = Student(username=username, school=school or "", grade=grade or "")
        db.add(student)
        db.commit()
        db.refresh(student)
    else:
        # Optionally update school/grade if provided
        updated = False
        if school and student.school != school:
            student.school = school
            updated = True
        if grade and student.grade != grade:
            student.grade = grade
            updated = True
        if updated:
            db.commit()

    # Issue JWT using existing helper
    access_token_expires = timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
    token = create_access_token(data={"sub": username}, expires_delta=access_token_expires)

    # Redirect to sub_project (token in query)
    redirect_url = f"{KETI3_FRONTEND_URL}?token={token}"
    return RedirectResponse(url=redirect_url, status_code=302)

# Authentication endpoints
@app.post("/api/auth/login", response_model=Token)
async def login(user_data: UserLogin, db: Session = Depends(get_db)):
    """Login via MySQL students table"""
    stu = get_student_by_username(db, user_data.username)
    if not stu or not stu.password or not check_password(user_data.password, stu.password):
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Incorrect username or password")
    access_token_expires = timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
    access_token = create_access_token(data={"sub": user_data.username}, expires_delta=access_token_expires)
    return {"access_token": access_token, "token_type": "bearer"}

@app.post("/api/auth/register")
async def register_user(user_data: UserRegister, db: Session = Depends(get_db)):
    """用户注册"""
    # 验证输入数据
    if not user_data.username.strip():
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="用户名不能为空"
        )
    
    if not user_data.password.strip():
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="密码不能为空"
        )
    
    if len(user_data.password) < 6:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="密码长度至少6位"
        )
    
    # 哈希密码
    hashed_password = make_hashed_password(user_data.password)
    # 创建用户（MySQL）
    user_id = create_student_mysql(
        db,
        username=user_data.username,
        school=user_data.school,
        student_id=user_data.student_id,
        grade=user_data.grade,
        password_hash=hashed_password,
    )
    if user_id is None:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="用户名已存在"
        )
    
    return {
        "message": "用户注册成功",
        "user_id": str(user_id),
        "username": user_data.username
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
async def keti3_oss_upload(request: Request, img: UploadFile = File(None), audio: UploadFile = File(None)):
    """Proxy upload endpoint to avoid browser-to-MinIO CORS. Accepts multipart form-data with fields 'img' and 'audio'."""
    verify_signature(request)

    if img is None and audio is None:
        return error(ERRCODE_PARAM_ERROR, "No files provided. Expect 'img' and/or 'audio' fields.")

    try:
        client = get_minio_client()
        bucket = os.getenv("MINIO_BUCKET", "onlineclass")
        ensure_bucket(client, bucket)
        base_path = os.getenv("MINIO_BASE_PATH", "24game")

        now = datetime.utcnow()
        ts = int(now.timestamp())

        uid = 0
        try:
            uid = int(request.headers.get("X-User-Id", "0"))
        except Exception:
            uid = 0

        public_base = f"{MINIO_URL}/{bucket}"
        result: dict[str, str] = {}

        if img is not None:
            img_key = f"{base_path}/{uid}/{ts}/img{uid}{ts}"
            data = await img.read()
            client.put_object(bucket, img_key, data=io.BytesIO(data), length=len(data), content_type=img.content_type or "image/*")
            result["imgUrl"] = f"{public_base}/{img_key}"

        if audio is not None:
            audio_key = f"{base_path}/{uid}/{ts}/audio{uid}{ts}"
            data = await audio.read()
            client.put_object(bucket, audio_key, data=io.BytesIO(data), length=len(data), content_type=audio.content_type or "audio/*")
            result["audioUrl"] = f"{public_base}/{audio_key}"

        return keti3_response(data=result)
    except Exception as e:
        import traceback
        print(f"OSS Proxy Upload Error: {e}")
        print(f"Traceback: {traceback.format_exc()}")
        return error(ERRCODE_COMMON_ERROR, f"Upload failed: {str(e)}")

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
