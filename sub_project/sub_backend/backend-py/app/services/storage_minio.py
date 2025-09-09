import os
from datetime import timedelta
from typing import Optional
from minio import Minio
from minio.commonconfig import CopySource
from dotenv import load_dotenv

# Load environment variables
load_dotenv()


def get_minio_client() -> Minio:
    endpoint = os.getenv("MINIO_ENDPOINT", "localhost:9000")
    access_key = os.getenv("MINIO_ACCESS_KEY")
    secret_key = os.getenv("MINIO_SECRET_KEY")
    use_ssl = os.getenv("MINIO_USE_SSL", "false").lower() == "true"
    # Allow both http://host:9000 and host:9000
    if endpoint.startswith("http://"):
        endpoint = endpoint[len("http://") :]
        use_ssl = False
    elif endpoint.startswith("https://"):
        endpoint = endpoint[len("https://") :]
        use_ssl = True
    if not access_key or not secret_key:
        raise RuntimeError("MINIO_ACCESS_KEY / MINIO_SECRET_KEY not set")
    return Minio(endpoint, access_key=access_key, secret_key=secret_key, secure=use_ssl)


def presign_put_object(
    client: Minio,
    bucket: str,
    object_name: str,
    content_type: Optional[str] = None,
    expires: timedelta = timedelta(hours=1),
) -> str:
    # MinIO SDK 7.x+ doesn't support request_headers in presigned_put_object
    return client.presigned_put_object(bucket_name=bucket, object_name=object_name, expires=expires)


def presign_get_object(
    client: Minio,
    bucket: str,
    object_name: str,
    expires: timedelta = timedelta(hours=1),
) -> str:
    return client.presigned_get_object(bucket_name=bucket, object_name=object_name, expires=expires)


def ensure_bucket(client: Minio, bucket: str):
    found = client.bucket_exists(bucket)
    if not found:
        client.make_bucket(bucket)


def put_object_from_file(client: Minio, bucket: str, object_name: str, file_path: str, content_type: str = "application/octet-stream"):
    import os as _os
    size = _os.path.getsize(file_path)
    with open(file_path, "rb") as f:
        client.put_object(bucket, object_name, f, length=size, content_type=content_type)
