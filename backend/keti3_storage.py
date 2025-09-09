import os
from datetime import datetime, timedelta
from minio import Minio
from minio.error import S3Error
from urllib.parse import urlparse
from dotenv import load_dotenv
import ssl
import urllib3

# Load environment variables
load_dotenv()

# MinIO Configuration for keti3
MINIO_URL = os.getenv("MINIO_URL", "http://localhost:9000")
MINIO_ACCESS_KEY = os.getenv("MINIO_ACCESS_KEY")
MINIO_SECRET_KEY = os.getenv("MINIO_SECRET_KEY")
MINIO_BUCKET_NAME = os.getenv("MINIO_BUCKET", "onlineclass")
MINIO_CA_CERT = os.getenv("MINIO_CA_CERT", "")
MINIO_INSECURE_TLS = os.getenv("MINIO_INSECURE_TLS", "false").lower() in ("1", "true", "yes")

def get_minio_client():
    """Get MinIO client instance"""
    parsed = urlparse(MINIO_URL)
    secure = parsed.scheme == "https"
    endpoint = parsed.netloc
    # Default http_client
    http_client = None
    if secure:
        # Prefer custom CA if provided
        if MINIO_CA_CERT:
            ctx = ssl.create_default_context(cafile=MINIO_CA_CERT)
            http_client = urllib3.PoolManager(ssl_context=ctx)
        elif MINIO_INSECURE_TLS:
            # Disable verification for self-signed certs in dev
            ctx = ssl._create_unverified_context()
            http_client = urllib3.PoolManager(ssl_context=ctx, assert_hostname=False)

    return Minio(
        endpoint,
        access_key=MINIO_ACCESS_KEY,
        secret_key=MINIO_SECRET_KEY,
        secure=secure,
        http_client=http_client
    )

def ensure_bucket(client, bucket_name):
    """Ensure bucket exists, create if not"""
    try:
        if not client.bucket_exists(bucket_name):
            client.make_bucket(bucket_name)
    except S3Error as e:
        raise Exception(f"Failed to create/check bucket: {e}")

def presign_put_object(client, bucket, object_name, content_type=None, expires=None):
    """Generate presigned PUT URL for object upload"""
    if expires is None:
        expires = timedelta(hours=1)
    
    try:
        # For MinIO SDK 7.x, don't use request_headers parameter
        url = client.presigned_put_object(
            bucket,
            object_name,
            expires=expires
        )
        return url
    except Exception as e:
        raise Exception(f"Failed to generate presigned URL: {e}")
