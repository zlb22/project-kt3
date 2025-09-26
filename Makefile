# Makefile for project-kt3
# Commands:
#   Development: start, stop, status (HTTPS dev servers)
#   Production: production-build, production-deploy, production-start, production-stop
#   Legacy Dev: dev, dev-stop, dev-status (HTTP dev servers)
#   Deployment: deploy-local, deploy-remote

# 服务器配置
DB_SERVER       ?= 172.24.130.213
WEB_SERVER      ?= 172.24.125.63
DOMAIN          ?= kt3.bnu.edu.cn
REMOTE_USER     ?= ubuntu
SSH_PORT        ?= 22

# 生产环境路径
PROD_BACKEND_PATH   ?= /opt/project-kt3
PROD_FRONTEND_PATH  ?= /var/www/project-kt3

.PHONY: start stop status dev dev-stop dev-status build deploy-local deploy-remote
.PHONY: production-build production-deploy production-start production-stop production-status

# ---------------- Production HTTPS Services ----------------

start:
	@echo "Starting production HTTPS services..."
	@echo "Generating self-signed TLS certificate if missing..."
	@mkdir -p certs
	@[ -f certs/server.crt ] && [ -f certs/server.key ] || \
	  openssl req -x509 -nodes -newkey rsa:2048 -days 365 \
	    -keyout certs/server.key -out certs/server.crt \
	    -subj "/CN=localhost" \
	    -addext "subjectAltName=DNS:localhost,IP:127.0.0.1" >/dev/null 2>&1 || true
	@echo "[backend] Starting FastAPI HTTPS (https://0.0.0.0:8443)..."
	@cd backend && uv run uvicorn main:app --host 0.0.0.0 --port 8443 --reload \
	    --ssl-keyfile ../certs/server.key --ssl-certfile ../certs/server.crt \
	    > server-https.log 2>&1 & echo $$! > uvicorn-https.pid
	@echo "[frontend] Starting React HTTPS dev server (https://0.0.0.0:3000)..."
	@HOST=0.0.0.0 npm --prefix frontend run start-https > frontend/dev.log 2>&1 & echo $$! > frontend/dev.pid
	@echo "[sub-frontend] Starting Vue dev server with production URLs (https://0.0.0.0:5174)..."
	@cd sub_project/sub_frontend && npm run dev -- --host 0.0.0.0 --port 5174 --mode prod \
	    > dev-https.log 2>&1 & echo $$! > dev-https.pid
	@echo "HTTPS services started:"
	@echo "  Backend:      https://localhost:8443"
	@echo "  Frontend:     https://localhost:3000"
	@echo "  Sub-Frontend: https://localhost:5174"

# ---------------- Start backend on HTTP 8000 only ----------------
start-http:
	@echo "[backend] Starting FastAPI HTTP (http://0.0.0.0:8000)..."
	@cd backend && uv run uvicorn main:app --host 0.0.0.0 --port 8000 --reload \
	    > server-http.log 2>&1 & echo $$! > uvicorn-http.pid

stop:
	@echo "Stopping all services..."
	@# Frontend 3000
	@[ ! -f frontend/dev.pid ] || (kill $$(cat frontend/dev.pid) 2>/dev/null || true)
	@rm -f frontend/dev.pid
	@# Sub-frontend HTTPS 5174
	@[ ! -f sub_project/sub_frontend/dev-https.pid ] || (kill $$(cat sub_project/sub_frontend/dev-https.pid) 2>/dev/null || true)
	@rm -f sub_project/sub_frontend/dev-https.pid
	@# Sub-frontend HTTP 5173
	@[ ! -f sub_project/sub_frontend/dev.pid ] || (kill $$(cat sub_project/sub_frontend/dev.pid) 2>/dev/null || true)
	@rm -f sub_project/sub_frontend/dev.pid
	@# Backend HTTPS 8443
	@[ ! -f backend/uvicorn-https.pid ] || (kill $$(cat backend/uvicorn-https.pid) 2>/dev/null || true)
	@rm -f backend/uvicorn-https.pid
	@# Backend HTTP 8000
	@[ ! -f backend/uvicorn-http.pid ] || (kill $$(cat backend/uvicorn-http.pid) 2>/dev/null || true)
	@rm -f backend/uvicorn-http.pid
	@# Kill ports forcefully
	@{ command -v fuser >/dev/null && fuser -k 3000/tcp 5173/tcp 5174/tcp 8443/tcp 8000/tcp 2>/dev/null || true; } || { \
	    command -v lsof >/dev/null && for port in 3000 5173 5174 8443 8000; do \
	        pids=$$(lsof -ti:$$port 2>/dev/null || true); \
	        [ -z "$$pids" ] || kill $$pids 2>/dev/null || true; \
	    done; \
	}
	@echo "All services stopped."

status:
	@echo "Service status (3000,8443,8000):"
	@# Frontend 3000
	@pid_f1=""; [ -f frontend/dev.pid ] && pid_f1=$$(cat frontend/dev.pid) || true; \
	if [ -n "$$pid_f1" ] && ps -p $$pid_f1 >/dev/null 2>&1; then \
	  echo "[frontend 3000] RUNNING PID=$$pid_f1 URL=https://localhost:3000"; \
	else \
	  echo "[frontend 3000] STOPPED URL=https://localhost:3000"; \
	fi; \
	{ command -v lsof >/dev/null && ports_f1=$$(lsof -ti:3000 2>/dev/null || true) && [ -n "$$ports_f1" ] && echo "[frontend 3000] Port 3000 held by: $$ports_f1" || true; } || \
	{ command -v ss >/dev/null && ss -tulpn 2>/dev/null | grep ':3000' || true; } || true
	@# Backend HTTPS 8443
	@pid_b2=""; [ -f backend/uvicorn-https.pid ] && pid_b2=$$(cat backend/uvicorn-https.pid) || true; \
	if [ -n "$$pid_b2" ] && ps -p $$pid_b2 >/dev/null 2>&1; then \
	  echo "[backend https] RUNNING PID=$$pid_b2 URL=https://localhost:8443"; \
	else \
	  echo "[backend https] STOPPED URL=https://localhost:8443"; \
	fi; \
	{ command -v lsof >/dev/null && ports_b2=$$(lsof -ti:8443 2>/dev/null || true) && [ -n "$$ports_b2" ] && echo "[backend https] Port 8443 held by: $$ports_b2" || true; } || \
	{ command -v ss >/dev/null && ss -tulpn 2>/dev/null | grep ':8443' || true; } || true
	@# Backend HTTP 8000
	@pid_b1=""; [ -f backend/uvicorn-http.pid ] && pid_b1=$$(cat backend/uvicorn-http.pid) || true; \
	if [ -n "$$pid_b1" ] && ps -p $$pid_b1 >/dev/null 2>&1; then \
	  echo "[backend http] RUNNING PID=$$pid_b1 URL=http://localhost:8000"; \
	else \
	  echo "[backend http] STOPPED URL=http://localhost:8000"; \
	fi; \
	{ command -v lsof >/dev/null && ports_b1=$$(lsof -ti:8000 2>/dev/null || true) && [ -n "$$ports_b1" ] && echo "[backend http] Port 8000 held by: $$ports_b1" || true; } || \
	{ command -v ss >/dev/null && ss -tulpn 2>/dev/null | grep ':8000' || true; } || true
	@echo "Logs: backend/server-https.log | frontend/dev.log"


# ---------------- New: Local DEV (HTTP on 8000) ----------------
.PHONY: dev dev-stop dev-status

dev:
	@echo "[backend] Starting FastAPI HTTP (http://0.0.0.0:8000)..."
	@cd backend && uvicorn main:app --host 0.0.0.0 --port 8000 --reload \
	    > server-http.log 2>&1 & echo $$! > uvicorn-http.pid || true
	@echo "[frontend] Starting React dev server (http://0.0.0.0:3000)..."
	@HOST=0.0.0.0 npm --prefix frontend start > frontend/dev.log 2>&1 & echo $$! > frontend/dev.pid
	@echo "[sub-frontend] Starting Vite dev server (http://0.0.0.0:5173)..."
	@npm --prefix sub_project/sub_frontend run dev > sub_project/sub_frontend/dev.log 2>&1 & echo $$! > sub_project/sub_frontend/dev.pid || true
	@echo "DEV services started: API http://localhost:8000, CRA http://localhost:3000, Vite http://localhost:5173"

# ---------------- Frontend only commands ----------------
start-frontend:
	@echo "[frontend] Starting React HTTPS dev server (https://0.0.0.0:3000)..."
	@HOST=0.0.0.0 npm --prefix frontend run start-https > frontend/dev.log 2>&1 & echo $$! > frontend/dev.pid
	@echo "Frontend started: https://localhost:3000"

stop-frontend:
	@echo "Stopping frontend (3000)..."
	@[ ! -f frontend/dev.pid ] || (kill $$(cat frontend/dev.pid) 2>/dev/null || true)
	@rm -f frontend/dev.pid
	@{ command -v fuser >/dev/null && fuser -k 3000/tcp 2>/dev/null || true; } || { \
	    command -v lsof >/dev/null && pids=$$(lsof -ti:3000 2>/dev/null || true); \
	    [ -z "$$pids" ] || kill $$pids 2>/dev/null || true; \
	}
	@echo "Frontend stopped."

dev-stop:
	@echo "Stopping DEV services (3000,5173,8000)..."
	@[ ! -f frontend/dev.pid ] || (kill $$(cat frontend/dev.pid) 2>/dev/null || true)
	@rm -f frontend/dev.pid
	@[ ! -f sub_project/sub_frontend/dev.pid ] || (kill $$(cat sub_project/sub_frontend/dev.pid) 2>/dev/null || true)
	@rm -f sub_project/sub_frontend/dev.pid
	@[ ! -f backend/uvicorn-http.pid ] || (kill $$(cat backend/uvicorn-http.pid) 2>/dev/null || true)
	@rm -f backend/uvicorn-http.pid
	@{ command -v fuser >/dev/null && fuser -k 3000/tcp 2>/dev/null || true; } || true
	@{ command -v fuser >/dev/null && fuser -k 5173/tcp 2>/dev/null || true; } || true
	@{ command -v fuser >/dev/null && fuser -k 8000/tcp 2>/dev/null || true; } || true
	@echo "DEV services stopped."

dev-status:
	@echo "DEV Service status (3000,5173,8000):"
	@pid_f1=""; [ -f frontend/dev.pid ] && pid_f1=$$(cat frontend/dev.pid) || true; \
	if [ -n "$$pid_f1" ] && ps -p $$pid_f1 >/dev/null 2>&1; then \
	  echo "[frontend 3000] RUNNING PID=$$pid_f1 URL=http://localhost:3000"; \
	else \
	  echo "[frontend 3000] STOPPED URL=http://localhost:3000"; \
	fi; \
	{ command -v lsof >/dev/null && ports_f1=$$(lsof -ti:3000 2>/dev/null || true) && [ -n "$$ports_f1" ] && echo "[frontend 3000] Port 3000 held by: $$ports_f1" || true; } || true
	@pid_sf=""; [ -f sub_project/sub_frontend/dev.pid ] && pid_sf=$$(cat sub_project/sub_frontend/dev.pid) || true; \
	if [ -n "$$pid_sf" ] && ps -p $$pid_sf >/dev/null 2>&1; then \
	  echo "[sub-frontend 5173] RUNNING PID=$$pid_sf URL=http://localhost:5173"; \
	else \
	  echo "[sub-frontend 5173] STOPPED URL=http://localhost:5173"; \
	fi; \
	{ command -v lsof >/dev/null && ports_sf=$$(lsof -ti:5173 2>/dev/null || true) && [ -n "$$ports_sf" ] && echo "[sub-frontend 5173] Port 5173 held by: $$ports_sf" || true; } || true
	@pid_b1=""; [ -f backend/uvicorn-http.pid ] && pid_b1=$$(cat backend/uvicorn-http.pid) || true; \
	if [ -n "$$pid_b1" ] && ps -p $$pid_b1 >/dev/null 2>&1; then \
	  echo "[backend http] RUNNING PID=$$pid_b1 URL=http://localhost:8000"; \
	else \
	  echo "[backend http] STOPPED URL=http://localhost:8000"; \
	fi; \
	{ command -v lsof >/dev/null && ports_b1=$$(lsof -ti:8000 2>/dev/null || true) && [ -n "$$ports_b1" ] && echo "[backend http] Port 8000 held by: $$ports_b1" || true; } || true


# ---------------- New: Profile switching ----------------
.PHONY: switch-dev switch-prod

switch-dev:
	@echo "Switching profiles to DEV using A_LAN_IP=$(A_LAN_IP) ..."
	@mkdir -p backend frontend sub_project/sub_frontend
	@sed -e 's#__A_LAN_IP__#$(A_LAN_IP)#g' \
	    deploy/profiles/backend.env.dev > backend/.env
	@sed -e 's#__A_LAN_IP__#$(A_LAN_IP)#g' \
	    deploy/profiles/frontend.env.dev > frontend/.env.development
	@sed -e 's#__A_LAN_IP__#$(A_LAN_IP)#g' \
	    deploy/profiles/sub_frontend.env.dev > sub_project/sub_frontend/.env.development
	@echo "DEV profiles applied."

switch-prod:
	@echo "Switching profiles to PROD with domains: API=$(API_DOMAIN) APP=$(APP_DOMAIN) KETI3=$(KETI3_DOMAIN) MINIO=$(MINIO_DOMAIN) ..."
	@mkdir -p backend sub_project/sub_frontend
	@sed -e 's#__API_DOMAIN__#$(API_DOMAIN)#g' \
	     -e 's#__APP_DOMAIN__#$(APP_DOMAIN)#g' \
	     -e 's#__KETI3_DOMAIN__#$(KETI3_DOMAIN)#g' \
	     -e 's#__MINIO_DOMAIN__#$(MINIO_DOMAIN)#g' \
	    deploy/profiles/backend.env.prod > backend/.env
	@sed -e 's#__API_DOMAIN__#$(API_DOMAIN)#g' \
	     -e 's#__APP_DOMAIN__#$(APP_DOMAIN)#g' \
	     -e 's#__KETI3_DOMAIN__#$(KETI3_DOMAIN)#g' \
	     -e 's#__MINIO_DOMAIN__#$(MINIO_DOMAIN)#g' \
	    deploy/profiles/sub_frontend.env.prod > sub_project/sub_frontend/.env.production
	@echo "PROD profiles applied. Now build frontends with: make prod-build"


# ---------------- New: Production build helpers ----------------
.PHONY: prod-build

prod-build:
	@echo "Building CRA frontend (production)..."
	@npm --prefix frontend ci >/dev/null 2>&1 || npm --prefix frontend install
	@npm --prefix frontend run build
	@echo "Building Vite sub_frontend (production)..."
	@npm --prefix sub_project/sub_frontend ci >/dev/null 2>&1 || npm --prefix sub_project/sub_frontend install
	@npm --prefix sub_project/sub_frontend run build
	@echo "Builds completed. CRA=frontend/build, Vite=sub_project/sub_frontend/dist"

# ---------------- Deployment Commands ----------------

build:
	@echo "Building frontend applications..."
	@echo "  Building React frontend..."
	@cd frontend && npm ci --silent && PUBLIC_URL=/topic-three npm run build
	@echo "  Building Vue sub-frontend..."
	@cd sub_project/sub_frontend && npm ci --silent && npm run build-only
	@echo "Build completed."

deploy-local:
	@echo "Deploying locally..."
	@chmod +x deploy/scripts/deploy.sh
	@BACKEND_HOST=$(BACKEND_HOST) ./deploy/scripts/deploy.sh --local

deploy-remote:
	@echo "Deploying to remote server..."
	@chmod +x deploy/scripts/deploy.sh
	@REMOTE_USER=$(REMOTE_USER) REMOTE_HOST=$(REMOTE_HOST) SSH_PORT=$(SSH_PORT) BACKEND_HOST=$(BACKEND_HOST) ./deploy/scripts/deploy.sh --remote

# ---------------- Production Environment Commands ----------------

production-build:
	@echo "=== Building for Production Environment ==="
	@echo "Domain: $(DOMAIN)"
	@echo "DB Server: $(DB_SERVER)"
	@echo "Web Server: $(WEB_SERVER)"
	@chmod +x deploy/build-and-deploy.sh
	@./deploy/build-and-deploy.sh

production-deploy:
	@echo "=== Deploying to Production Servers ==="
	@if [ ! -f deploy/project-kt3-production-*.tar.gz ]; then \
		echo "No production build found. Running production-build first..."; \
		$(MAKE) production-build; \
	fi
	@latest_build=$$(ls -t deploy/project-kt3-production-*.tar.gz | head -1); \
	build_name=$$(basename "$$latest_build"); \
	echo "Using build: $$build_name"; \
	@if [ "$$(hostname -I | grep -o '172\.24\.130\.213')" ]; then \
		echo "=== Local deployment on DB Server (172.24.130.213) ==="; \
		echo "Deploying backend locally..."; \
		cp "$$latest_build" /tmp/; \
		cd /tmp && tar -xzf "$$build_name"; \
		echo "Running backend deployment script..."; \
		cd /tmp && sudo -S ./build/scripts/deploy-backend.sh < /dev/null; \
		echo "Uploading to Web Server ($(WEB_SERVER))..."; \
		scp -P $(SSH_PORT) "$$(pwd)/$$latest_build" $(REMOTE_USER)@$(WEB_SERVER):/tmp/; \
		echo "Executing deployment on Web Server..."; \
		ssh -t -p $(SSH_PORT) $(REMOTE_USER)@$(WEB_SERVER) "\
			cd /tmp && \
			tar -xzf $$build_name && \
			sudo ./build/scripts/deploy-web.sh"; \
	else \
		echo "=== Remote deployment mode ==="; \
		echo "Uploading to Web Server ($(WEB_SERVER))..."; \
		scp -P $(SSH_PORT) "$$latest_build" $(REMOTE_USER)@$(WEB_SERVER):/tmp/; \
		echo "Uploading to DB Server ($(DB_SERVER))..."; \
		scp -P $(SSH_PORT) "$$latest_build" $(REMOTE_USER)@$(DB_SERVER):/tmp/; \
		echo "Executing deployment on Web Server..."; \
		ssh -p $(SSH_PORT) $(REMOTE_USER)@$(WEB_SERVER) "\
			cd /tmp && \
			tar -xzf $$build_name && \
			sudo ./build/scripts/deploy-web.sh"; \
		echo "Executing deployment on DB Server..."; \
		ssh -p $(SSH_PORT) $(REMOTE_USER)@$(DB_SERVER) "\
			cd /tmp && \
			tar -xzf $$build_name && \
			sudo ./build/scripts/deploy-backend.sh"; \
	fi
	@echo "=== Production Deployment Complete ==="

production-start:
	@echo "=== Starting Production Services ==="
	@if [ "$$(hostname -I | grep -o '172\.24\.130\.213')" ]; then \
		echo "=== Local mode on DB Server (172.24.130.213) ==="; \
		echo "Starting backend service locally..."; \
		sudo systemctl start project-kt3-backend && \
		sudo systemctl enable project-kt3-backend; \
		echo "Reloading Nginx on Web Server ($(WEB_SERVER))..."; \
		ssh -p $(SSH_PORT) $(REMOTE_USER)@$(WEB_SERVER) "\
			sudo nginx -t && \
			sudo systemctl reload nginx"; \
	else \
		echo "=== Remote mode ==="; \
		echo "Starting backend service on DB Server ($(DB_SERVER))..."; \
		ssh -p $(SSH_PORT) $(REMOTE_USER)@$(DB_SERVER) "\
			sudo systemctl start project-kt3-backend && \
			sudo systemctl enable project-kt3-backend"; \
		echo "Reloading Nginx on Web Server ($(WEB_SERVER))..."; \
		ssh -p $(SSH_PORT) $(REMOTE_USER)@$(WEB_SERVER) "\
			sudo nginx -t && \
			sudo systemctl reload nginx"; \
	fi
	@echo "=== Production Services Started ==="
	@echo "Access URL: https://$(DOMAIN)"

production-stop:
	@echo "=== Stopping Production Services ==="
	@echo "Stopping backend service on DB Server ($(DB_SERVER))..."
	@ssh -p $(SSH_PORT) $(REMOTE_USER)@$(DB_SERVER) "\
		sudo systemctl stop project-kt3-backend" || true
	@echo "=== Production Services Stopped ==="

production-status:
	@echo "=== Production Services Status ==="
	@echo "Backend Service Status ($(DB_SERVER)):"
	@ssh -p $(SSH_PORT) $(REMOTE_USER)@$(DB_SERVER) "\
		sudo systemctl status project-kt3-backend --no-pager || true"
	@echo ""
	@echo "Nginx Status ($(WEB_SERVER)):"
	@ssh -p $(SSH_PORT) $(REMOTE_USER)@$(WEB_SERVER) "\
		sudo systemctl status nginx --no-pager || true"
	@echo ""
	@echo "Health Check:"
	@curl -k -s https://$(DOMAIN)/health || echo "Health check failed"

production-logs:
	@echo "=== Production Logs ==="
	@echo "Backend Logs ($(DB_SERVER)):"
	@ssh -p $(SSH_PORT) $(REMOTE_USER)@$(DB_SERVER) "\
		sudo journalctl -u project-kt3-backend -n 50 --no-pager"
	@echo ""
	@echo "Nginx Access Logs ($(WEB_SERVER)):"
	@ssh -p $(SSH_PORT) $(REMOTE_USER)@$(WEB_SERVER) "\
		sudo tail -n 20 /var/log/nginx/project-kt3-access.log || echo 'No access logs found'"
	@echo ""
	@echo "Nginx Error Logs ($(WEB_SERVER)):"
	@ssh -p $(SSH_PORT) $(REMOTE_USER)@$(WEB_SERVER) "\
		sudo tail -n 20 /var/log/nginx/project-kt3-error.log || echo 'No error logs found'"

production-clean:
	@echo "=== Cleaning Production Build Files ==="
	@rm -f deploy/project-kt3-production-*.tar.gz
	@rm -rf deploy/build
	@echo "Build files cleaned."

# 生产环境证书管理
production-ssl-setup:
	@echo "=== Setting up SSL Certificates ==="
	@echo "Creating self-signed certificates for $(DOMAIN)..."
	@ssh -p $(SSH_PORT) $(REMOTE_USER)@$(WEB_SERVER) "\
		sudo mkdir -p /etc/nginx/ssl && \
		sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
			-keyout /etc/nginx/ssl/project-kt3.key \
			-out /etc/nginx/ssl/project-kt3.crt \
			-subj '/CN=$(DOMAIN)' \
			-addext 'subjectAltName=DNS:$(DOMAIN)'"
	@echo "SSL certificates created."

# 生产环境帮助信息
production-help:
	@echo "=== Production Environment Commands ==="
	@echo "make production-build     - Build production packages"
	@echo "make production-deploy    - Deploy to production servers"
	@echo "make production-start     - Start production services"
	@echo "make production-stop      - Stop production services"
	@echo "make production-status    - Check production status"
	@echo "make production-logs      - View production logs"
	@echo "make production-clean     - Clean build files"
	@echo "make production-ssl-setup - Setup SSL certificates"
	@echo ""
	@echo "Environment Variables:"
	@echo "  DOMAIN=$(DOMAIN)"
	@echo "  DB_SERVER=$(DB_SERVER)"
	@echo "  WEB_SERVER=$(WEB_SERVER)"
	@echo "  REMOTE_USER=$(REMOTE_USER)"
