#!/bin/bash

# Project-kt3 部署验证脚本
# 验证生产环境部署是否正确，特别是跳转逻辑

set -e

# 配置
DOMAIN="${DOMAIN:-172-24-125-63.sslip.io}"
DB_SERVER="${DB_SERVER:-172.24.130.213}"
WEB_SERVER="${WEB_SERVER:-172.24.125.63}"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Project-kt3 部署验证 ===${NC}"
echo -e "${YELLOW}域名: $DOMAIN${NC}"
echo -e "${YELLOW}数据库服务器: $DB_SERVER${NC}"
echo -e "${YELLOW}Web服务器: $WEB_SERVER${NC}"

# 验证函数
check_url() {
    local url=$1
    local expected_status=$2
    local description=$3
    
    echo -n "检查 $description ($url)... "
    
    if response=$(curl -k -s -o /dev/null -w "%{http_code}" "$url" 2>/dev/null); then
        if [ "$response" = "$expected_status" ]; then
            echo -e "${GREEN}✓ $response${NC}"
            return 0
        else
            echo -e "${RED}✗ $response (期望: $expected_status)${NC}"
            return 1
        fi
    else
        echo -e "${RED}✗ 连接失败${NC}"
        return 1
    fi
}

check_redirect() {
    local url=$1
    local description=$2
    
    echo -n "检查重定向 $description ($url)... "
    
    if response=$(curl -k -s -I "$url" 2>/dev/null); then
        status=$(echo "$response" | grep "HTTP" | awk '{print $2}')
        location=$(echo "$response" | grep -i "location:" | cut -d' ' -f2- | tr -d '\r')
        
        if [ "$status" = "302" ] && [ -n "$location" ]; then
            echo -e "${GREEN}✓ 302 → $location${NC}"
            return 0
        else
            echo -e "${RED}✗ 状态: $status, Location: $location${NC}"
            return 1
        fi
    else
        echo -e "${RED}✗ 连接失败${NC}"
        return 1
    fi
}

# 开始验证
echo -e "\n${BLUE}1. 基础连通性测试${NC}"

# HTTP重定向到HTTPS
check_url "http://$DOMAIN" "301" "HTTP重定向"

# 主前端
check_url "https://$DOMAIN" "200" "主前端首页"
check_url "https://$DOMAIN/health" "200" "健康检查"

# 静态资源
check_url "https://$DOMAIN/static/css/main.css" "200" "CSS资源" || echo "  (CSS文件可能不存在，这是正常的)"
check_url "https://$DOMAIN/favicon.ico" "200" "Favicon" || echo "  (Favicon可能不存在，这是正常的)"

echo -e "\n${BLUE}2. API接口测试${NC}"

# 后端API (通过Nginx代理)
check_url "https://$DOMAIN/api/auth/captcha" "200" "验证码API"
check_url "https://$DOMAIN/api/auth/public-key" "200" "公钥API"

echo -e "\n${BLUE}3. 子前端测试${NC}"

# 子前端直接访问
check_url "https://$DOMAIN/topic-three/online-experiment/" "200" "子前端首页"

echo -e "\n${BLUE}4. 跳转逻辑测试${NC}"

# 注意：这些测试需要有效的JWT token，所以可能会失败
echo "注意：以下跳转测试需要认证，可能显示401错误"

# 后端入口 (需要认证)
check_url "https://$DOMAIN/web/keti3/entry" "401" "后端入口 (未认证)"

echo -e "\n${BLUE}5. 服务器服务状态${NC}"

echo "检查数据库服务器后端服务..."
if ssh -o ConnectTimeout=5 ubuntu@$DB_SERVER "systemctl is-active project-kt3-backend" >/dev/null 2>&1; then
    echo -e "${GREEN}✓ 后端服务运行中${NC}"
else
    echo -e "${RED}✗ 后端服务未运行${NC}"
fi

echo "检查Web服务器Nginx状态..."
if ssh -o ConnectTimeout=5 ubuntu@$WEB_SERVER "systemctl is-active nginx" >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Nginx服务运行中${NC}"
else
    echo -e "${RED}✗ Nginx服务未运行${NC}"
fi

echo -e "\n${BLUE}6. 跳转配置验证${NC}"

echo "验证环境变量配置..."

# 检查后端配置
echo "后端KETI3_FRONTEND_URL配置:"
ssh ubuntu@$DB_SERVER "grep KETI3_FRONTEND_URL /opt/project-kt3/.env" 2>/dev/null || echo "  配置文件不存在或无法访问"

# 检查子前端配置
echo "子前端VITE_LOGIN_URL配置:"
echo "  应该指向: https://$DOMAIN/login"

echo -e "\n${BLUE}7. 完整跳转流程测试${NC}"

echo "模拟完整跳转流程："
echo "1. 用户访问: https://$DOMAIN"
echo "2. 登录后点击'在线实验'"
echo "3. 前端调用: /api/auth/create-sub-token"
echo "4. 跳转到: https://$DOMAIN/topic-three/online-experiment/?token=..."
echo "5. 子前端'返回主页'跳转到: https://$DOMAIN/dashboard"

echo -e "\n${GREEN}=== 验证完成 ===${NC}"
echo -e "${YELLOW}请手动测试完整的用户流程以确保跳转正确${NC}"
echo -e "${YELLOW}访问 https://$DOMAIN 进行完整测试${NC}"
