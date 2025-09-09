import hashlib
import time
import os
from fastapi import HTTPException, Request
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Signature middleware configuration
APP_ID = os.getenv("APP_ID", "keti3_app")
SKEY = os.getenv("SKEY", "your_secret_key")

def verify_signature(request: Request):
    """Verify request signature for keti3 APIs"""
    try:
        # Get signature from headers
        signature = request.headers.get("signature")
        timestamp = request.headers.get("timestamp")
        app_id = request.headers.get("appId")
        
        if not all([signature, timestamp, app_id]):
            raise HTTPException(status_code=401, detail="Missing signature headers")
        
        # Verify app_id
        if app_id != APP_ID:
            raise HTTPException(status_code=401, detail="Invalid app_id")
        
        # Check timestamp (within 5 minutes)
        current_time = int(time.time())
        request_time = int(timestamp)
        if abs(current_time - request_time) > 300:  # 5 minutes
            raise HTTPException(status_code=401, detail="Request expired")
        
        # Generate expected signature
        expected_signature = hashlib.md5(f"{SKEY}{timestamp}".encode()).hexdigest()
        
        if signature != expected_signature:
            raise HTTPException(status_code=401, detail="Invalid signature")
        
        return True
    except ValueError:
        raise HTTPException(status_code=401, detail="Invalid timestamp format")
    except Exception as e:
        raise HTTPException(status_code=401, detail=f"Signature verification failed: {str(e)}")

def keti3_response(data=None, message="success", code=0):
    """Standard keti3 API response format"""
    return {
        "code": code,
        "message": message,
        "data": data
    }

def error(code, message):
    """Error response helper"""
    return keti3_response(data=None, message=message, code=code)

# Error codes
ERRCODE_COMMON_ERROR = 1000
ERRCODE_INVALID_PARAMS = 1001
ERRCODE_AUTH_FAILED = 1002
