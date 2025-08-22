#!/usr/bin/env python3
import os
import sys
from urllib.parse import urlparse
from dotenv import load_dotenv
from minio import Minio
from minio.error import S3Error


def get_minio_client():
    load_dotenv()

    url = os.getenv("MINIO_URL", "http://localhost:9000")
    access_key = os.getenv("MINIO_ACCESS_KEY")
    secret_key = os.getenv("MINIO_SECRET_KEY")
    access_key = "admin"
    secret_key = "12345678"

    if not access_key or not secret_key:
        print("[ERROR] MINIO_ACCESS_KEY or MINIO_SECRET_KEY is not set in .env", file=sys.stderr)
        sys.exit(2)

    parsed = urlparse(url)
    if not parsed.scheme or not parsed.netloc:
        print(f"[ERROR] MINIO_URL is invalid: {url}", file=sys.stderr)
        sys.exit(2)

    secure = parsed.scheme == "https"
    endpoint = parsed.netloc  # host:port

    client = Minio(
        endpoint,
        access_key=access_key,
        secret_key=secret_key,
        secure=secure,
    )
    return client


def main():
    bucket = os.getenv("MINIO_BUCKET_NAME", "onlineclass")

    # Warn if using console port 9001
    url = os.getenv("MINIO_URL", "http://localhost:9001")
    if ":9001" in url:
        print("[WARN] MINIO_URL appears to use port 9001 (console). The S3 API usually runs on 9000.")

    try:
        client = get_minio_client()

        # Simple operation to verify connectivity: list buckets
        try:
            buckets = client.list_buckets()
            print(f"[OK] Connected to MinIO. Buckets: {[b.name for b in buckets]}")
        except S3Error as e:
            print(f"[ERROR] Failed to list buckets: {e}", file=sys.stderr)
            sys.exit(1)

        # Ensure target bucket exists
        try:
            exists = client.bucket_exists(bucket)
            if not exists:
                client.make_bucket(bucket)
                print(f"[OK] Bucket created: {bucket}")
            else:
                print(f"[OK] Bucket exists: {bucket}")
        except S3Error as e:
            print(f"[ERROR] Bucket check/create failed for '{bucket}': {e}", file=sys.stderr)
            sys.exit(1)

        print("[SUCCESS] MinIO connectivity verified.")
        sys.exit(0)

    except Exception as e:
        print(f"[FATAL] Unexpected error: {e}", file=sys.stderr)
        sys.exit(3)


if __name__ == "__main__":
    main()
