-- 创建数据库
CREATE DATABASE IF NOT EXISTS kt3 CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE kt3;

-- 创建学生表
CREATE TABLE IF NOT EXISTS students (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    school VARCHAR(100) NOT NULL,
    grade VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_username (username)
);

-- 创建操作日志表
CREATE TABLE IF NOT EXISTS operation_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid INT NOT NULL,
    action VARCHAR(100) NOT NULL,
    details TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_uid (uid),
    INDEX idx_timestamp (timestamp)
);
