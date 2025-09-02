#!/usr/bin/env python3
"""
测试用户注册和登录功能
"""

import asyncio
import aiohttp
import json
import sys

BASE_URL = "http://localhost:8000"

async def test_registration():
    """测试用户注册"""
    print("=== 测试用户注册 ===")
    
    # 测试数据
    test_users = [
        {
            "username": "testuser1",
            "school": "清华大学",
            "student_id": "2021001001",
            "grade": "大三",
            "password": "password123"
        },
        {
            "username": "testuser2", 
            "school": "北京大学",
            "student_id": "2021002001",
            "grade": "大二",
            "password": "password456"
        }
    ]
    
    async with aiohttp.ClientSession() as session:
        for user in test_users:
            try:
                async with session.post(
                    f"{BASE_URL}/api/auth/register",
                    json=user
                ) as response:
                    result = await response.json()
                    if response.status == 200:
                        print(f"✓ 用户 {user['username']} 注册成功")
                        print(f"  用户ID: {result.get('user_id')}")
                    else:
                        print(f"❌ 用户 {user['username']} 注册失败: {result}")
            except Exception as e:
                print(f"❌ 注册请求失败: {e}")

async def test_login():
    """测试用户登录"""
    print("\n=== 测试用户登录 ===")
    
    login_data = {
        "username": "testuser1",
        "password": "password123"
    }
    
    async with aiohttp.ClientSession() as session:
        try:
            async with session.post(
                f"{BASE_URL}/api/auth/login",
                json=login_data
            ) as response:
                result = await response.json()
                if response.status == 200:
                    print(f"✓ 用户 {login_data['username']} 登录成功")
                    token = result.get('access_token')
                    print(f"  Token: {token[:50]}...")
                    
                    # 测试获取用户信息
                    await test_user_info(session, token)
                else:
                    print(f"❌ 登录失败: {result}")
        except Exception as e:
            print(f"❌ 登录请求失败: {e}")

async def test_user_info(session, token):
    """测试获取用户信息"""
    print("\n=== 测试获取用户信息 ===")
    
    headers = {"Authorization": f"Bearer {token}"}
    
    try:
        async with session.get(
            f"{BASE_URL}/api/auth/me",
            headers=headers
        ) as response:
            result = await response.json()
            if response.status == 200:
                print("✓ 获取用户信息成功")
                print(f"  用户名: {result.get('username')}")
                print(f"  学校: {result.get('school')}")
                print(f"  学号: {result.get('student_id')}")
                print(f"  年级: {result.get('grade')}")
                print(f"  状态: {'活跃' if result.get('is_active') else '非活跃'}")
            else:
                print(f"❌ 获取用户信息失败: {result}")
    except Exception as e:
        print(f"❌ 获取用户信息请求失败: {e}")

async def test_duplicate_registration():
    """测试重复注册"""
    print("\n=== 测试重复注册 ===")
    
    duplicate_user = {
        "username": "testuser1",  # 已存在的用户名
        "school": "复旦大学",
        "student_id": "2021003001",
        "grade": "大一",
        "password": "newpassword"
    }
    
    async with aiohttp.ClientSession() as session:
        try:
            async with session.post(
                f"{BASE_URL}/api/auth/register",
                json=duplicate_user
            ) as response:
                result = await response.json()
                if response.status == 400:
                    print("✓ 重复注册正确被拒绝")
                    print(f"  错误信息: {result.get('detail')}")
                else:
                    print(f"❌ 重复注册应该被拒绝，但返回: {result}")
        except Exception as e:
            print(f"❌ 重复注册测试失败: {e}")

async def main():
    """主测试函数"""
    print("开始测试用户认证功能...")
    print("确保后端服务正在运行在 http://localhost:8000\n")
    
    await test_registration()
    await test_login()
    await test_duplicate_registration()
    
    print("\n测试完成！")

if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print("\n测试被用户中断")
    except Exception as e:
        print(f"测试过程中发生错误: {e}")
        sys.exit(1)
