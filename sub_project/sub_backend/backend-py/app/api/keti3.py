from fastapi import APIRouter, Request, Depends
from pydantic import BaseModel
from typing import Optional, Dict
from ..utils.response import success, error, ERRCODE_COMMON_ERROR
from ..services.storage_minio import get_minio_client, presign_put_object, ensure_bucket
from datetime import datetime, timedelta
import os
from sqlalchemy.orm import Session
from sqlalchemy import select, text
from sqlalchemy.exc import IntegrityError
from ..db.session import get_db
from ..orm.student_info import StudentInfo
from ..orm.op_log import OpLog

router = APIRouter()


class StudentLoginReq(BaseModel):
    school_id: int
    student_name: str
    student_num: str
    grade_id: int
    class_name: str


@router.post("/student/login")
async def student_login(req: StudentLoginReq, request: Request, db: Session = Depends(get_db)):
    # 查找是否已有记录（唯一键：school_id, student_num, grade_id, class_name）
    stmt = (
        select(StudentInfo)
        .where(StudentInfo.school_id == req.school_id)
        .where(StudentInfo.student_num == req.student_num)
        .where(StudentInfo.grade_id == req.grade_id)
        .where(StudentInfo.class_name == req.class_name)
        .limit(1)
    )
    row = db.execute(stmt).scalars().first()

    if row is None:
        # 插入并用自增 id 赋值为 uid
        row = StudentInfo(
            uid=0,
            school_id=req.school_id,
            student_name=req.student_name,
            student_num=req.student_num,
            grade_id=req.grade_id,
            class_name=req.class_name,
        )
        db.add(row)
        try:
            db.flush()  # 获取自增 id
        except IntegrityError:
            db.rollback()
            # 并发下唯一键冲突，再查一次
            row = db.execute(stmt).scalars().first()
            if row is None:
                # 仍异常则抛出
                raise
        if row.uid in (None, 0):
            row.uid = row.id
            try:
                db.commit()
            except IntegrityError:
                db.rollback()
        else:
            db.commit()
    else:
        # 已存在，必要时补写 uid
        if not row.uid:
            row.uid = row.id
            try:
                db.commit()
            except IntegrityError:
                db.rollback()

    return success({"uid": int(row.uid)})


@router.post("/config/list")
async def config_list():
    # Mock data - replace with actual database queries
    schools = [
        {"id": 1, "name": "School A"},
        {"id": 2, "name": "School B"},
        {"id": 3, "name": "School C"},
    ]
    grades = [
        {"id": 1, "name": "一年级"},
        {"id": 2, "name": "二年级"},
        {"id": 3, "name": "三年级"},
        {"id": 4, "name": "四年级"},
        {"id": 5, "name": "五年级"},
        {"id": 6, "name": "六年级"},
    ]
    # 前端期望 school 和 grade 字段名（单数）
    return success({"school": schools, "grade": grades})


class OssAuthReq(BaseModel):
    uid: Optional[int]
    content_types: Optional[Dict[str, str]] = None  # {"img": "image/png", "audio": "audio/mpeg"}


@router.post("/oss/auth")
async def oss_auth(req: OssAuthReq):
    try:
        # Generate two presigned PUT URLs for image and audio
        client = get_minio_client()
        bucket = os.getenv("MINIO_BUCKET", os.getenv("MINIO_BUCKET_NAME", "onlineclass"))
        ensure_bucket(client, bucket)
        base_path = os.getenv("MINIO_BASE_PATH", "24game")
    except Exception as e:
        import traceback
        print(f"OSS Auth Error: {e}")
        print(f"Traceback: {traceback.format_exc()}")
        return error(ERRCODE_COMMON_ERROR, f"MinIO error: {str(e)}")

    now = datetime.utcnow()
    ts = int(now.timestamp())
    uid = req.uid or 0

    # object keys
    img_key = f"{base_path}/{uid}/{ts}/img{uid}{ts}"
    audio_key = f"{base_path}/{uid}/{ts}/audio{uid}{ts}"

    img_ct = (req.content_types or {}).get("img", "image/*")
    audio_ct = (req.content_types or {}).get("audio", "audio/*")

    expires = timedelta(hours=1)
    img_put = presign_put_object(client, bucket, img_key, content_type=img_ct, expires=expires)
    audio_put = presign_put_object(client, bucket, audio_key, content_type=audio_ct, expires=expires)

    # Public URL base (for constructing access URL); if behind proxy, set MINIO_PUBLIC_BASE
    public_base = os.getenv("MINIO_PUBLIC_BASE")
    if not public_base:
        endpoint = os.getenv("MINIO_ENDPOINT", "localhost:9000")
        scheme = "https" if os.getenv("MINIO_USE_SSL", "false").lower() == "true" else "http"
        public_base = f"{scheme}://{endpoint}/{bucket}"

    resp = {
        "storage": "minio",
        "Bucket": bucket,
        "HttpsDomain": public_base,
        "Presigned": [
            {"key": img_key, "url": img_put, "method": "PUT", "content_type": img_ct, "public_url": f"{public_base}/{img_key}"},
            {"key": audio_key, "url": audio_put, "method": "PUT", "content_type": audio_ct, "public_url": f"{public_base}/{audio_key}"},
        ],
        # compatibility fields for old OSS STS (kept empty)
        "AccessKeyId": "",
        "AccessKeySecret": "",
        "SecurityToken": "",
        "Expiration": (now + expires).isoformat() + "Z",
    }
    return success(resp)


# ----- UID INCREASE -----
class UidReq(BaseModel):
    count: Optional[int] = 1


@router.post("/uid/get-by-increase")
async def uid_get_by_increase(req: UidReq, db: Session = Depends(get_db)):
    # 简易自增序列表：keti3_uid_seq(id bigint auto_increment primary key)
    # 如不存在则创建
    db.execute(text(
        """
        CREATE TABLE IF NOT EXISTS keti3_uid_seq (
            id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
            PRIMARY KEY (id)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
        """
    ))
    db.commit()

    need = max(1, int(req.count or 1))
    uids = []
    for _ in range(need):
        db.execute(text("INSERT INTO keti3_uid_seq VALUES ()"))
        r = db.execute(text("SELECT LAST_INSERT_ID()"))
        new_id = r.scalar_one()
        uids.append(int(new_id))
    db.commit()
    return success({"uids": uids})


# ----- LOG SAVE -----
class LogSaveReq(BaseModel):
    uid: int
    op_time: Optional[str] = None  # ISO string
    op_type: str
    op_object: Optional[str] = None
    object_no: Optional[str] = None
    object_name: Optional[str] = None
    data_before: Optional[dict] = None
    data_after: Optional[dict] = None
    voice_url: Optional[str] = None
    screenshot_url: Optional[str] = None
    submit_id: Optional[int] = None


@router.post("/log/save")
async def log_save(req: LogSaveReq, db: Session = Depends(get_db)):
    try:
        op_time = None
        if req.op_time:
            try:
                op_time = datetime.fromisoformat(req.op_time.replace("Z", "+00:00"))
            except Exception:
                pass
        row = OpLog(
            submit_id=req.submit_id,
            uid=req.uid,
            op_time=op_time,
            op_type=req.op_type,
            op_object=req.op_object,
            object_no=req.object_no,
            object_name=req.object_name,
            data_before=req.data_before,
            data_after=req.data_after,
            voice_url=req.voice_url,
            screenshot_url=req.screenshot_url,
        )
        db.add(row)
        db.commit()
        return success({"id": int(row.id)})
    except Exception as e:
        db.rollback()
        return error(ERRCODE_COMMON_ERROR, f"log save failed: {e}")


# ----- LOG LIST -----
class LogListReq(BaseModel):
    uid: int
    page: Optional[int] = 1
    size: Optional[int] = 20


@router.post("/log/list")
async def log_list(req: LogListReq, db: Session = Depends(get_db)):
    page = max(1, int(req.page or 1))
    size = max(1, min(100, int(req.size or 20)))
    offset = (page - 1) * size

    q = select(OpLog).where(OpLog.uid == req.uid).order_by(OpLog.op_time.desc()).offset(offset).limit(size)
    rows = db.execute(q).scalars().all()
    items = []
    for r in rows:
        items.append({
            "id": int(r.id),
            "uid": int(r.uid),
            "op_time": (r.op_time.isoformat() if r.op_time else None),
            "op_type": r.op_type,
            "op_object": r.op_object,
            "object_no": r.object_no,
            "object_name": r.object_name,
            "data_before": r.data_before,
            "data_after": r.data_after,
            "voice_url": r.voice_url,
            "screenshot_url": r.screenshot_url,
        })
    return success({"list": items, "page": page, "size": size})


# ----- LOG EXPORT -----
class LogExportReq(BaseModel):
    uid: int


@router.get("/log/export")
async def log_export(uid: int, db: Session = Depends(get_db)):
    # 查询指定 uid 的日志
    rows = db.execute(select(OpLog).where(OpLog.uid == uid).order_by(OpLog.op_time.asc())).scalars().all()
    # 生成 XLSX 并打包 ZIP，上传 MinIO，返回下载链接
    from openpyxl import Workbook
    import tempfile, os as _os, zipfile
    from ..services.storage_minio import ensure_bucket, put_object_from_file, presign_get_object

    wb = Workbook()
    ws = wb.active
    ws.title = "op_log"
    ws.append(["id", "uid", "op_time", "op_type", "op_object", "object_no", "object_name", "voice_url", "screenshot_url"])  # 简化列
    for r in rows:
        ws.append([
            int(r.id), int(r.uid),
            (r.op_time.isoformat() if r.op_time else ""), r.op_type or "",
            r.op_object or "", r.object_no or "", r.object_name or "",
            r.voice_url or "", r.screenshot_url or "",
        ])

    tmpdir = tempfile.mkdtemp(prefix="export_")
    xlsx_path = _os.path.join(tmpdir, f"op_log_{uid}.xlsx")
    wb.save(xlsx_path)

    zip_path = _os.path.join(tmpdir, f"op_log_{uid}.zip")
    with zipfile.ZipFile(zip_path, 'w', zipfile.ZIP_DEFLATED) as zf:
        zf.write(xlsx_path, _os.path.basename(xlsx_path))

    client = get_minio_client()
    bucket = os.getenv("MINIO_BUCKET", os.getenv("MINIO_BUCKET_NAME", "onlineclass"))
    ensure_bucket(client, bucket)
    base_path = os.getenv("MINIO_BASE_PATH", "24game")
    object_name = f"{base_path}/export/{uid}/op_log_{uid}_{int(datetime.utcnow().timestamp())}.zip"
    put_object_from_file(client, bucket, object_name, zip_path, content_type="application/zip")

    # 生成 GET 预签名
    url = presign_get_object(client, bucket, object_name)
    return success({"download_url": url, "bucket": bucket, "object": object_name})
