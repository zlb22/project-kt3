from fastapi import FastAPI, HTTPException, Depends, status
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

# JWT Configuration
SECRET_KEY = "your-secret-key-change-in-production"
ALGORITHM = "HS256"
ACCESS_TOKEN_EXPIRE_MINUTES = 30

app = FastAPI(title="Educational Assessment API", version="2.0")
security = HTTPBearer()

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],  # React dev server
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

def load_users():
    try:
        with open('users.json', 'r') as f:
            return json.load(f)
    except FileNotFoundError:
        default_users = {
            "admin": make_hashed_password("1")
        }
        save_users(default_users)
        return default_users

def save_users(users_data):
    with open('users.json', 'w') as f:
        json.dump(users_data, f)

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
    users = load_users()
    if user_data.username not in users:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect username or password"
        )
    
    if not check_password(user_data.password, users[user_data.username]):
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
    users = load_users()
    if not check_password(password_data.old_password, users[current_user]):
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Incorrect old password"
        )
    
    users[current_user] = make_hashed_password(password_data.new_password)
    save_users(users)
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

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
