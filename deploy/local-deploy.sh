#!/bin/bash

# Project-kt3 本地部署脚本
# 在172.24.130.213上执行，部署后端到本地，前端到172.24.125.63

set -e

# 配置
WEB_SERVER="172.24.125.63"
REMOTE_USER="ubuntu"
PROJECT_ROOT="/home/ubuntu/zlb/project-kt3"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Project-kt3 本地部署 ===${NC}"

# 检查是否在正确的服务器上
if ! hostname -I | grep -q "172.24.130.213"; then
    echo -e "${RED}错误: 此脚本只能在172.24.130.213上运行${NC}"
    exit 1
fi

# 查找最新的构建包
latest_build=$(ls -t "$PROJECT_ROOT/deploy/project-kt3-production-"*.tar.gz 2>/dev/null | head -1)
if [ -z "$latest_build" ]; then
    echo -e "${RED}错误: 找不到构建包，请先运行 make production-build${NC}"
    exit 1
fi

build_name=$(basename "$latest_build")
echo -e "${YELLOW}使用构建包: $build_name${NC}"

# 1. 本地部署后端
echo -e "\n${GREEN}[步骤1] 本地部署后端服务${NC}"
echo "解压构建包到 /tmp..."
cp "$latest_build" /tmp/
cd /tmp
tar -xzf "$build_name"

echo "执行后端部署..."
cd /tmp/build
sudo ./scripts/deploy-backend.sh

# 2. 上传并部署前端到Web服务器
echo -e "\n${GREEN}[步骤2] 部署前端到Web服务器${NC}"
echo "上传构建包到 $WEB_SERVER..."
scp "$latest_build" "$REMOTE_USER@$WEB_SERVER:/tmp/"

echo "在Web服务器上执行部署..."
ssh -t "$REMOTE_USER@$WEB_SERVER" "
    cd /tmp && 
    tar -xzf $build_name && 
    cd build &&
    sudo ./scripts/deploy-web.sh
"

echo -e "\n${GREEN}=== 部署完成 ===${NC}"
echo -e "${YELLOW}下一步:${NC}"
echo "1. 启动后端服务: sudo systemctl start project-kt3-backend"
echo "2. 检查服务状态: sudo systemctl status project-kt3-backend"
echo "3. 访问网站: https://172-24-125-63.sslip.io"
