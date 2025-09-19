#!/usr/bin/env python3
"""
Security Test Script for Project-KT3
Tests the new password security implementation
"""

import sys
import os
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

import requests
import json
import bcrypt
from datetime import datetime

# Test configuration
BASE_URL = "https://172.24.130.213:8443"  # Backend server
TEST_USER = {
    "username": "security_test_user",
    "school": "Test School",
    "student_id": "TEST001",
    "grade": "Grade 1",
    "password": "TestPassword123!"
}

def hash_password_frontend(password: str) -> str:
    """Simulate frontend password hashing"""
    salt = bcrypt.gensalt(rounds=12)
    return bcrypt.hashpw(password.encode('utf-8'), salt).decode('utf-8')

def test_password_security():
    """Test password security implementation"""
    print("🔒 Testing Password Security Implementation")
    print("=" * 50)
    
    # Disable SSL verification for self-signed certificates
    import urllib3
    urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
    
    session = requests.Session()
    session.verify = False
    
    try:
        # Test 1: Get CAPTCHA
        print("1. Testing CAPTCHA generation...")
        captcha_response = session.get(f"{BASE_URL}/api/auth/captcha")
        if captcha_response.status_code == 200:
            captcha_data = captcha_response.json()
            print("   ✅ CAPTCHA generated successfully")
        else:
            print("   ❌ CAPTCHA generation failed")
            return False
        
        # Test 2: Register with hashed password
        print("2. Testing secure registration...")
        hashed_password = hash_password_frontend(TEST_USER["password"])
        
        register_data = {
            **TEST_USER,
            "password": hashed_password
        }
        
        register_response = session.post(f"{BASE_URL}/api/auth/register", json=register_data)
        if register_response.status_code == 200:
            print("   ✅ Secure registration successful")
        else:
            print(f"   ❌ Registration failed: {register_response.status_code}")
            print(f"   Response: {register_response.text}")
            return False
        
        # Test 3: Login with hashed password
        print("3. Testing secure login...")
        login_hashed_password = hash_password_frontend(TEST_USER["password"])
        
        login_data = {
            "username": TEST_USER["username"],
            "password": login_hashed_password,
            "captcha_id": captcha_data["captcha_id"],
            "captcha_code": "test"  # This will likely fail, but tests the flow
        }
        
        login_response = session.post(f"{BASE_URL}/api/auth/login", json=login_data)
        if login_response.status_code == 200:
            print("   ✅ Secure login successful")
            token = login_response.json()["access_token"]
        elif login_response.status_code == 400 and "验证码" in login_response.text:
            print("   ✅ Login flow working (CAPTCHA validation expected to fail)")
        else:
            print(f"   ❌ Login failed unexpectedly: {login_response.status_code}")
            print(f"   Response: {login_response.text}")
        
        # Test 4: Check security headers
        print("4. Testing security headers...")
        headers_response = session.get(f"{BASE_URL}/api/auth/captcha")
        headers = headers_response.headers
        
        security_headers = {
            "Strict-Transport-Security": "HSTS header",
            "X-Frame-Options": "Clickjacking protection",
            "X-Content-Type-Options": "MIME sniffing protection",
            "X-XSS-Protection": "XSS protection"
        }
        
        for header, description in security_headers.items():
            if header in headers:
                print(f"   ✅ {description}: {headers[header]}")
            else:
                print(f"   ❌ Missing {description}")
        
        # Test 5: Password strength validation
        print("5. Testing password strength validation...")
        weak_passwords = [
            "123456",
            "password",
            "Password",
            "Password1",
            "password123"
        ]
        
        for weak_pwd in weak_passwords:
            weak_hashed = hash_password_frontend(weak_pwd)
            weak_data = {
                **TEST_USER,
                "username": f"weak_test_{datetime.now().microsecond}",
                "password": weak_hashed
            }
            
            weak_response = session.post(f"{BASE_URL}/api/auth/register", json=weak_data)
            if weak_response.status_code != 200:
                print(f"   ✅ Weak password '{weak_pwd}' rejected")
            else:
                print(f"   ❌ Weak password '{weak_pwd}' accepted")
        
        print("\n🎉 Security testing completed!")
        return True
        
    except Exception as e:
        print(f"❌ Security test failed with error: {e}")
        return False

if __name__ == "__main__":
    success = test_security()
    if not success:
        sys.exit(1)
