# Keti3 Project Makefile
# 一键启动和管理前后端服务

.PHONY: help start stop restart backend frontend logs clean status

# 默认目标
help:
	@echo "Keti3 Project Management"
	@echo "========================"
	@echo "make start     - 启动后端和前端服务"
	@echo "make stop      - 停止所有服务"
	@echo "make restart   - 重启所有服务"
	@echo "make backend   - 仅启动后端服务"
	@echo "make frontend  - 仅启动前端服务"
	@echo "make logs      - 查看服务日志"
	@echo "make status    - 查看服务状态"
	@echo "make clean     - 清理日志文件"

# 启动所有服务
start:
	@echo "🚀 启动 Keti3 项目服务..."
	@make backend &
	@sleep 3
	@make frontend &
	@echo "✅ 服务启动完成!"
	@echo "📱 前端: http://localhost:5173/topic-three"
	@echo "🔧 后端: http://localhost:8096"

# 停止所有服务
stop:
	@echo "🛑 停止所有服务..."
	@pkill -f "uvicorn app.main:app" || true
	@pkill -f "vite.*sub_frontend" || true
	@echo "✅ 所有服务已停止"

# 重启服务
restart: stop
	@sleep 2
	@make start

# 启动后端服务
backend:
	@echo "🔧 启动后端服务 (端口 8096)..."
	@cd sub_backend/backend-py && \
	nohup uvicorn app.main:app --host 0.0.0.0 --port 8096 --reload > ../../logs/backend.log 2>&1 &
	@echo "✅ 后端服务已启动"

# 启动前端服务
frontend:
	@echo "📱 启动前端服务 (端口 5173)..."
	@cd sub_frontend && \
	nohup npm run dev > ../logs/frontend.log 2>&1 &
	@echo "✅ 前端服务已启动"

# 查看日志
logs:
	@echo "📋 查看服务日志..."
	@echo "=== 后端日志 ==="
	@tail -n 20 logs/backend.log 2>/dev/null || echo "后端日志文件不存在"
	@echo ""
	@echo "=== 前端日志 ==="
	@tail -n 20 logs/frontend.log 2>/dev/null || echo "前端日志文件不存在"

# 查看服务状态
status:
	@echo "📊 服务状态检查..."
	@echo "=== 后端服务 (8096端口) ==="
	@if pgrep -f "uvicorn app.main:app" > /dev/null; then \
		echo "✅ 后端服务运行中"; \
		curl -s http://localhost:8096/health > /dev/null && echo "✅ 后端健康检查通过" || echo "❌ 后端健康检查失败"; \
	else \
		echo "❌ 后端服务未运行"; \
	fi
	@echo ""
	@echo "=== 前端服务 (5173端口) ==="
	@if pgrep -f "vite.*sub_frontend" > /dev/null; then \
		echo "✅ 前端服务运行中"; \
		curl -s http://localhost:5173 > /dev/null && echo "✅ 前端服务可访问" || echo "❌ 前端服务不可访问"; \
	else \
		echo "❌ 前端服务未运行"; \
	fi
	@echo ""
	@echo "=== Docker 依赖服务 ==="
	@docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep -E "(mysql|minio|mongodb)" || echo "❌ 未发现 Docker 依赖服务"

# 清理日志文件
clean:
	@echo "🧹 清理日志文件..."
	@rm -f logs/*.log
	@echo "✅ 日志文件已清理"

# 初始化项目（首次使用）
init:
	@echo "🔧 初始化项目环境..."
	@mkdir -p logs
	@echo "📦 安装前端依赖..."
	@cd sub_frontend && npm install
	@echo "📦 安装后端依赖..."
	@cd sub_backend/backend-py && pip install -r requirements.txt
	@echo "✅ 项目初始化完成"

# 开发模式（实时查看日志）
dev:
	@echo "🔧 开发模式启动..."
	@mkdir -p logs
	@make backend
	@sleep 3
	@make frontend
	@echo "📋 实时查看日志 (Ctrl+C 退出)..."
	@tail -f logs/backend.log logs/frontend.log

# 测试接口
test:
	@echo "🧪 测试后端接口..."
	@echo "测试健康检查..."
	@curl -s http://localhost:8096/health || echo "❌ 健康检查失败"
	@echo ""
	@echo "测试配置接口..."
	@TS=$$(date +%s); \
	SKEY='YJNGFHGASBDG#Y8d'; \
	SIGN=$$(printf "%s" "$${SKEY}$${TS}" | md5sum | awk '{print $$1}'); \
	curl -s -X POST 'http://localhost:8096/web/keti3/config/list' \
		-H "X-Auth-AppId: keti-3" \
		-H "X-Auth-TimeStamp: $$TS" \
		-H "X-Sign: $$SIGN" \
		-H 'Content-Type: application/json' \
		-d '{}' | jq . || echo "❌ 配置接口测试失败"
