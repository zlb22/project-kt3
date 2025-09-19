#!/bin/bash

# Project-kt3 一键打包部署脚本
# 支持两台服务器环境：
# - 172.24.130.213: 数据库、MinIO、后端服务
# - 172.24.125.63: 前端服务、Nginx代理

set -e  # 遇到错误立即退出

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置变量
PROJECT_ROOT="/home/ubuntu/zlb/project-kt3"
DEPLOY_DIR="$PROJECT_ROOT/deploy"
BUILD_DIR="$DEPLOY_DIR/build"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# 服务器配置
DB_SERVER="172.24.130.213"
WEB_SERVER="172.24.125.63"
DOMAIN="172-24-125-63.sslip.io"  # 临时域名，后续替换为正式域名

echo -e "${BLUE}=== Project-kt3 生产环境部署脚本 ===${NC}"
echo -e "${YELLOW}时间戳: $TIMESTAMP${NC}"
echo -e "${YELLOW}数据库服务器: $DB_SERVER${NC}"
echo -e "${YELLOW}Web服务器: $WEB_SERVER${NC}"
echo -e "${YELLOW}域名: $DOMAIN${NC}"

# 函数：打印步骤
print_step() {
    echo -e "\n${GREEN}[步骤] $1${NC}"
}

# 函数：打印错误
print_error() {
    echo -e "\n${RED}[错误] $1${NC}"
}

# 函数：检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        print_error "$1 命令未找到，请先安装"
        exit 1
    fi
}

# 检查必要工具
print_step "检查必要工具"
check_command "node"
check_command "npm"
check_command "rsync"

# 创建构建目录
print_step "准备构建目录"
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"/{frontend,sub_frontend,backend,nginx,scripts}

# 1. 构建主前端
print_step "构建主前端 (React)"
cd "$PROJECT_ROOT/frontend"

# 创建生产环境配置
cat > .env.production << EOF
# 生产环境配置
REACT_APP_API_URL=https://$DOMAIN
GENERATE_SOURCEMAP=false
EOF

# 更新依赖并构建
echo "更新前端依赖..."
npm install
npm run build

# 复制构建结果
cp -r build/* "$BUILD_DIR/frontend/"
echo -e "${GREEN}✓ 主前端构建完成${NC}"

# 2. 构建子前端
print_step "构建子前端 (Vue)"
cd "$PROJECT_ROOT/sub_project/sub_frontend"

# 创建生产环境配置
cat > .env.production << EOF
VITE_SKEY='JYTFGFJUYT9#Q7d'
VITE_APP_ID=keti-3
VITE_LOGIN_URL=https://$DOMAIN/login
EOF

# 更新依赖并构建
echo "更新子前端依赖..."
npm install
# 跳过类型检查，直接构建 (生产环境优先)
npm run build-only

# 复制构建结果
cp -r dist/* "$BUILD_DIR/sub_frontend/"
echo -e "${GREEN}✓ 子前端构建完成${NC}"

# 3. 准备后端文件
print_step "准备后端文件"
cd "$PROJECT_ROOT/backend"

# 创建生产环境配置
cat > .env.production << EOF
# 生产环境配置
SECRET_KEY=kt3-production-secret-$(openssl rand -hex 32)
ALGORITHM=HS256
ACCESS_TOKEN_EXPIRE_MINUTES=30

# MongoDB Configuration
MONGODB_URL=mongodb://$DB_SERVER:27017
DATABASE_NAME=project-kt3

# 生产CORS配置
CORS_ORIGINS=https://$DOMAIN,http://$DOMAIN

# Server Configuration
HOST=0.0.0.0
PORT=8000

# Upload Configuration
UPLOAD_DIR=uploads

# MinIO Configuration
MINIO_URL=https://$DB_SERVER:9000
MINIO_ACCESS_KEY=admin
MINIO_SECRET_KEY=12345678
MINIO_BUCKET_KETI3=online-experiment

# 生产环境前端URL
KETI3_FRONTEND_URL=https://$DOMAIN/topic-three/online-experiment

MINIO_BASE_PATH_KETI3=creative_work
MINIO_ALLOWED_BUCKETS=onlineclass,another-bucket

# Keti3 API Configuration
APP_ID=keti-3
SKEY=YJNGFHGASBDG#Y8d
MYSQL_DSN=mysql+pymysql://root:123456@$DB_SERVER:3306/kt3
MINIO_BASE_PATH=24game

# 生产环境日志配置
UPLOADS_JSON_LOG_ENABLED=true
EOF

# 复制后端文件（排除开发文件）
rsync -av --exclude='__pycache__' --exclude='*.pyc' --exclude='.env' \
    --exclude='server-*.log' --exclude='*.pid' \
    ./ "$BUILD_DIR/backend/"

# 复制生产配置
cp .env.production "$BUILD_DIR/backend/.env"
echo -e "${GREEN}✓ 后端文件准备完成${NC}"

# 4. 生成Nginx配置
print_step "生成Nginx配置"
cat > "$BUILD_DIR/nginx/project-kt3.conf" << EOF
# Project-kt3 Nginx配置
# 适用于服务器: $WEB_SERVER
# 域名: $DOMAIN

# HTTP重定向到HTTPS
server {
    listen 80;
    server_name $DOMAIN;
    return 301 https://\$server_name\$request_uri;
}

# HTTPS主配置
server {
    listen 443 ssl http2;
    server_name $DOMAIN;
    
    # SSL配置 (使用Let's Encrypt或自签名证书)
    ssl_certificate /etc/nginx/ssl/project-kt3.crt;
    ssl_certificate_key /etc/nginx/ssl/project-kt3.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    
    # 安全头
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    
    # Gzip压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/javascript application/xml+rss application/json;
    
    # 主前端 (默认路由)
    location / {
        root /var/www/project-kt3/frontend;
        try_files \$uri \$uri/ /index.html;
        
        # 缓存静态资源
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)\$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }
    
    # 后端API代理到数据库服务器
    location /api/ {
        proxy_pass https://$DB_SERVER:8443;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_ssl_verify off;
        
        # 超时设置
        proxy_connect_timeout 30s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;
    }
    
    # 后端入口代理
    location /web/ {
        proxy_pass https://$DB_SERVER:8443;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_ssl_verify off;
    }
    
    # 子前端 (在线实验)
    location /topic-three/online-experiment/ {
        root /var/www/project-kt3/sub_frontend;
        try_files \$uri \$uri/ /index.html;
        
        # 移除路径前缀
        rewrite ^/topic-three/online-experiment/(.*)\$ /\$1 break;
    }
    
    # 健康检查
    location /health {
        access_log off;
        return 200 "healthy\n";
        add_header Content-Type text/plain;
    }
}
EOF

echo -e "${GREEN}✓ Nginx配置生成完成${NC}"

# 5. 生成部署脚本
print_step "生成部署脚本"

# Web服务器部署脚本
cat > "$BUILD_DIR/scripts/deploy-web.sh" << 'EOF'
#!/bin/bash
# Web服务器部署脚本

set -e

DEPLOY_ROOT="/var/www/project-kt3"
NGINX_CONF="/etc/nginx/sites-available/project-kt3"
BACKUP_DIR="/var/backups/project-kt3/$(date +%Y%m%d_%H%M%S)"

echo "=== Web服务器部署开始 ==="

# 创建备份
if [ -d "$DEPLOY_ROOT" ]; then
    echo "创建备份到 $BACKUP_DIR"
    sudo mkdir -p "$BACKUP_DIR"
    sudo cp -r "$DEPLOY_ROOT" "$BACKUP_DIR/"
fi

# 创建部署目录
sudo mkdir -p "$DEPLOY_ROOT"/{frontend,sub_frontend}

# 部署前端文件
echo "部署前端文件..."
sudo cp -r frontend/* "$DEPLOY_ROOT/frontend/"
sudo cp -r sub_frontend/* "$DEPLOY_ROOT/sub_frontend/"

# 设置权限
sudo chown -R www-data:www-data "$DEPLOY_ROOT"
sudo chmod -R 755 "$DEPLOY_ROOT"

# 部署Nginx配置
echo "部署Nginx配置..."
sudo cp nginx/project-kt3.conf "$NGINX_CONF"
sudo ln -sf "$NGINX_CONF" /etc/nginx/sites-enabled/project-kt3

# 测试Nginx配置
sudo nginx -t

# 重载Nginx
sudo systemctl reload nginx

echo "=== Web服务器部署完成 ==="
EOF

# 数据库服务器部署脚本
cat > "$BUILD_DIR/scripts/deploy-backend.sh" << 'EOF'
#!/bin/bash
# 数据库服务器部署脚本

set -e

DEPLOY_ROOT="/opt/project-kt3"
SERVICE_NAME="project-kt3-backend"
BACKUP_DIR="/var/backups/project-kt3/$(date +%Y%m%d_%H%M%S)"

echo "=== 数据库服务器部署开始 ==="

# 停止现有服务
if systemctl is-active --quiet "$SERVICE_NAME"; then
    echo "停止现有服务..."
    sudo systemctl stop "$SERVICE_NAME"
fi

# 创建备份
if [ -d "$DEPLOY_ROOT" ]; then
    echo "创建备份到 $BACKUP_DIR"
    sudo mkdir -p "$BACKUP_DIR"
    sudo cp -r "$DEPLOY_ROOT" "$BACKUP_DIR/"
fi

# 创建部署目录
sudo mkdir -p "$DEPLOY_ROOT"

# 部署后端文件
echo "部署后端文件..."
sudo cp -r backend/* "$DEPLOY_ROOT/"

# 安装Python依赖
cd "$DEPLOY_ROOT"
echo "安装Python依赖..."
# 检查是否安装了 uv
if command -v uv >/dev/null 2>&1; then
    echo "使用 uv 安装依赖..."
    sudo uv sync --frozen
else
    echo "使用 pip 安装依赖..."
    if [ -f "requirements.txt" ]; then
        sudo python3 -m pip install -r requirements.txt
    else
        echo "错误: 未找到 requirements.txt 文件"
        exit 1
    fi
fi

# 设置权限
sudo chown -R ubuntu:ubuntu "$DEPLOY_ROOT"
sudo chmod +x "$DEPLOY_ROOT"/*.py

# 创建systemd服务文件
sudo tee "/etc/systemd/system/$SERVICE_NAME.service" > /dev/null << SYSTEMD_EOF
[Unit]
Description=Project-kt3 Backend Service
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=$DEPLOY_ROOT
Environment=PATH=/usr/local/bin:/usr/bin:/bin
ExecStart=/usr/local/bin/uv run uvicorn main:app --host 0.0.0.0 --port 8443 --ssl-keyfile ../certs/server.key --ssl-certfile ../certs/server.crt
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
SYSTEMD_EOF

# 重载systemd并启动服务
sudo systemctl daemon-reload
sudo systemctl enable "$SERVICE_NAME"
sudo systemctl start "$SERVICE_NAME"

echo "=== 数据库服务器部署完成 ==="
echo "服务状态:"
sudo systemctl status "$SERVICE_NAME" --no-pager
EOF

chmod +x "$BUILD_DIR/scripts"/*.sh

echo -e "${GREEN}✓ 部署脚本生成完成${NC}"

# 6. 生成部署包
print_step "打包部署文件"
cd "$DEPLOY_DIR"
tar -czf "project-kt3-production-$TIMESTAMP.tar.gz" build/

echo -e "${GREEN}✓ 部署包生成完成: project-kt3-production-$TIMESTAMP.tar.gz${NC}"

# 7. 生成部署说明
cat > "$BUILD_DIR/README.md" << EOF
# Project-kt3 生产环境部署包

**构建时间**: $TIMESTAMP
**目标环境**: 生产环境
**数据库服务器**: $DB_SERVER
**Web服务器**: $WEB_SERVER
**域名**: $DOMAIN

## 部署步骤

### 1. Web服务器 ($WEB_SERVER) 部署
\`\`\`bash
# 解压部署包
tar -xzf project-kt3-production-$TIMESTAMP.tar.gz
cd build

# 执行部署脚本
sudo ./scripts/deploy-web.sh
\`\`\`

### 2. 数据库服务器 ($DB_SERVER) 部署
\`\`\`bash
# 复制后端文件到数据库服务器
scp -r build/backend ubuntu@$DB_SERVER:/tmp/
scp build/scripts/deploy-backend.sh ubuntu@$DB_SERVER:/tmp/

# 在数据库服务器上执行
ssh ubuntu@$DB_SERVER
cd /tmp
sudo ./deploy-backend.sh
\`\`\`

### 3. 验证部署
- 访问: https://$DOMAIN
- 检查服务状态: \`sudo systemctl status project-kt3-backend\`
- 查看日志: \`sudo journalctl -u project-kt3-backend -f\`

## 回滚步骤
如需回滚，可使用 /var/backups/project-kt3/ 下的备份文件。
EOF

print_step "部署完成"
echo -e "${GREEN}=== 构建和打包完成 ===${NC}"
echo -e "${YELLOW}部署包位置: $DEPLOY_DIR/project-kt3-production-$TIMESTAMP.tar.gz${NC}"
echo -e "${YELLOW}部署说明: $BUILD_DIR/README.md${NC}"
echo ""
echo -e "${BLUE}下一步操作:${NC}"
echo -e "1. 将部署包传输到目标服务器"
echo -e "2. 在Web服务器上执行 deploy-web.sh"
echo -e "3. 在数据库服务器上执行 deploy-backend.sh"
echo -e "4. 验证部署结果"
