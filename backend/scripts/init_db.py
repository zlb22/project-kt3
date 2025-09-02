#!/usr/bin/env python3
"""
MongoDB数据库初始化脚本
创建必要的索引和集合
"""

import asyncio
import os
import sys
from motor.motor_asyncio import AsyncIOMotorClient
from dotenv import load_dotenv

# 添加父目录到路径以便导入配置
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

# 加载环境变量
load_dotenv()

# MongoDB配置
MONGODB_URL = os.getenv("MONGODB_URL", "mongodb://localhost:27017")
DATABASE_NAME = os.getenv("DATABASE_NAME", "project-kt3")

async def init_database():
    """初始化数据库和索引"""
    print(f"连接到MongoDB: {MONGODB_URL}")
    print(f"数据库名称: {DATABASE_NAME}")
    
    # 连接MongoDB
    client = AsyncIOMotorClient(MONGODB_URL)
    db = client[DATABASE_NAME]
    
    try:
        # 测试连接
        await client.admin.command('ping')
        print("✓ MongoDB连接成功")
        
        # 创建用户集合索引
        users_collection = db.users
        
        # 为用户名创建唯一索引
        await users_collection.create_index("username", unique=True)
        print("✓ 创建用户名唯一索引")
        
        # 为学号创建索引（可能需要查询）
        await users_collection.create_index("student_id")
        print("✓ 创建学号索引")
        
        # 为学校创建索引（可能需要按学校查询）
        await users_collection.create_index("school")
        print("✓ 创建学校索引")
        
        # 为创建时间创建索引
        await users_collection.create_index("created_at")
        print("✓ 创建创建时间索引")
        
        # 创建在线课程集合索引
        onlineclass_collection = db.onlineclass
        
        # 为用户创建索引
        await onlineclass_collection.create_index("user")
        print("✓ 创建在线课程用户索引")
        
        # 为测试ID创建索引
        await onlineclass_collection.create_index("test_id")
        print("✓ 创建测试ID索引")
        
        # 为提交时间创建索引
        await onlineclass_collection.create_index("submit_time")
        print("✓ 创建提交时间索引")
        
        # 复合索引：用户+测试ID
        await onlineclass_collection.create_index([("user", 1), ("test_id", 1)])
        print("✓ 创建用户+测试ID复合索引")
        
        # 显示所有索引
        print("\n=== 用户集合索引 ===")
        async for index in users_collection.list_indexes():
            print(f"- {index}")
            
        print("\n=== 在线课程集合索引 ===")
        async for index in onlineclass_collection.list_indexes():
            print(f"- {index}")
            
        print("\n✓ 数据库初始化完成！")
        
    except Exception as e:
        print(f"❌ 数据库初始化失败: {e}")
        return False
    finally:
        client.close()
    
    return True

if __name__ == "__main__":
    success = asyncio.run(init_database())
    sys.exit(0 if success else 1)
