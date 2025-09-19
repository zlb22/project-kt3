# Project-kt3 部署与跳转配置指南

## 概述

本文档详细说明了 project-kt3 在开发环境和生产环境下的跳转逻辑，以及如何配置 Nginx 实现统一的访问体验。

## 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Main Frontend │    │     Backend     │    │  Sub Frontend   │
│   (React CRA)   │    │   (FastAPI)     │    │   (Vue Vite)    │
│     Port 3000   │    │   Port 8443     │    │   Port 5174     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 开发环境配置

### 端口分配
- **主前端 (React)**: `https://172.24.130.213:3000`
- **后端 (FastAPI)**: `https://172.24.130.213:8443`
- **子前端 (Vue)**: `https://172.24.130.213:5174`

### 跳转流程

#### 1. 主前端 → 子前端
```
用户点击"在线实验" 
→ Dashboard.tsx 调用 /api/auth/create-sub-token
→ 后端返回 JWT token
→ 前端构造 URL: https://172.24.130.213:5174/topic-three/online-experiment/?token=xxx
→ 浏览器跳转到子前端
```

#### 2. 子前端 → 主前端
```
用户点击"返回主页"
→ NaviBar.vue 检测当前端口为 5174
→ 自动跳转到: https://172.24.130.213:3000/dashboard
```

### 环境配置文件

#### backend/.env (开发环境)
```env
# 自动检测模式
KETI3_FRONTEND_URL=https://172.24.130.213:5174/topic-three/online-experiment/

# CORS 配置
CORS_ORIGINS=https://172.24.130.213:3000,https://172.24.130.213:5174
```

#### sub_project/sub_frontend/.env.development
```env
VITE_LOGIN_URL=https://172.24.130.213:3000/login
```

**注意**: 子前端的 axios 配置会自动检测环境：
- 开发环境：5174 → 8443 (API 请求直接到后端)
- 生产环境：通过 Nginx 代理

## 生产环境配置

### Nginx 反向代理架构
```
Internet → Nginx (443/80) → Backend Services
                ├── / → Main Frontend (3000)
                ├── /api → Backend (8443)
                └── /topic-three/online-experiment → Sub Frontend (5174)
```

### 推荐的 Nginx 配置

```nginx
server {
    listen 443 ssl http2;
    server_name yourdomain.com;
    
    # SSL 配置
    ssl_certificate /path/to/your/cert.pem;
    ssl_certificate_key /path/to/your/key.pem;
    
    # 主前端 (默认路由)
    location / {
        proxy_pass https://127.0.0.1:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_ssl_verify off;
    }
    
    # 后端 API
    location /api/ {
        proxy_pass https://127.0.0.1:8443;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_ssl_verify off;
    }
    
    # 后端入口
    location /web/ {
        proxy_pass https://127.0.0.1:8443;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_ssl_verify off;
    }
    
    # 子前端
    location /topic-three/online-experiment/ {
        proxy_pass https://127.0.0.1:5174/topic-three/online-experiment/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_ssl_verify off;
    }
}

# HTTP 重定向到 HTTPS
server {
    listen 80;
    server_name yourdomain.com;
    return 301 https://$server_name$request_uri;
}
```

### 生产环境跳转流程

#### 1. 主前端 → 子前端
```
用户访问: https://yourdomain.com/dashboard
用户点击"在线实验"
→ 调用 /api/auth/create-sub-token (通过 Nginx 代理到后端)
→ 后端返回 JWT token
→ 前端构造 URL: https://yourdomain.com/topic-three/online-experiment/?token=xxx
→ Nginx 代理到子前端 5174
```

#### 2. 子前端 → 主前端
```
用户在子前端点击"返回主页"
→ NaviBar.vue 检测非开发环境
→ 跳转到: https://yourdomain.com/dashboard
→ Nginx 代理到主前端 3000
```

### 生产环境配置文件

#### backend/.env.production
```env
# 生产环境配置
KETI3_FRONTEND_URL=https://yourdomain.com/topic-three/online-experiment

# 生产 CORS 配置
CORS_ORIGINS=https://yourdomain.com

# 其他生产配置
SECRET_KEY=your-production-secret-key
MYSQL_DSN=mysql+pymysql://root:password@localhost:3306/kt3
```

#### sub_project/sub_frontend/.env.production
```env
VITE_LOGIN_URL=https://yourdomain.com/login
```

## 部署步骤

### 1. 开发环境启动
```bash
# 启动所有服务
make start

# 验证端口监听
ss -tulpn | grep -E "(3000|5174|8443)"
```

### 2. 生产环境部署

#### 步骤 1: 准备配置文件
```bash
# 复制生产配置
cp backend/.env.production backend/.env
cp sub_project/sub_frontend/.env.production sub_project/sub_frontend/.env

# 修改域名和密钥
vim backend/.env
vim sub_project/sub_frontend/.env
```

#### 步骤 2: 构建前端
```bash
# 构建主前端
cd frontend
npm run build

# 构建子前端
cd ../sub_project/sub_frontend
npm run build
```

#### 步骤 3: 配置 Nginx
```bash
# 复制配置文件
sudo cp /path/to/nginx.conf /etc/nginx/sites-available/project-kt3
sudo ln -s /etc/nginx/sites-available/project-kt3 /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重载配置
sudo nginx -s reload
```

#### 步骤 4: 启动后端服务
```bash
# 启动后端 (生产模式)
cd backend
uvicorn main:app --host 0.0.0.0 --port 8443 --ssl-keyfile ../certs/server.key --ssl-certfile ../certs/server.crt
```

## 故障排除

### 常见问题

#### 1. CORS 错误
**症状**: 控制台显示 "Access to XMLHttpRequest ... has been blocked by CORS policy"
**解决**: 检查 `backend/.env` 中的 `CORS_ORIGINS` 配置

#### 2. 跳转到错误地址
**症状**: 点击按钮后跳转到开发端口或错误域名
**解决**: 检查环境变量配置，确保 `VITE_LOGIN_URL` 正确

#### 3. 子前端无法访问
**症状**: 5174 端口无响应或证书错误
**解决**: 
```bash
# 检查端口监听
ss -tulpn | grep 5174

# 检查子前端日志
tail -f sub_project/sub_frontend/dev-https.log
```

#### 4. Token 验证失败
**症状**: 跳转后显示未登录或认证失败
**解决**: 检查 JWT 密钥配置，确保前后端使用相同的 `SECRET_KEY`

### 调试命令

```bash
# 检查所有服务状态
make status

# 查看后端日志
tail -f backend/server-https.log

# 测试 API 端点
curl -k -H "Authorization: Bearer $TOKEN" https://localhost:8443/api/auth/me

# 测试跳转端点
curl -k -I -H "Authorization: Bearer $TOKEN" https://localhost:8443/web/keti3/entry
```

## 安全考虑

1. **HTTPS 强制**: 生产环境必须使用 HTTPS
2. **JWT 安全**: 使用强密钥，设置合理的过期时间
3. **CORS 限制**: 生产环境严格限制 CORS 来源
4. **证书管理**: 使用有效的 SSL 证书，避免自签名证书

## 维护建议

1. **定期更新**: 保持依赖包和系统更新
2. **监控日志**: 监控错误日志和访问日志
3. **备份配置**: 定期备份配置文件和数据库
4. **性能优化**: 使用 CDN 和缓存优化静态资源加载

---

**最后更新**: 2025-09-19
**版本**: v1.0
**维护者**: 开发团队
