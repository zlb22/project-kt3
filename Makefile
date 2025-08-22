# Makefile for project-kt3
# Simple commands: install (update environment), start, stop

.PHONY: install start stop status gen-cert start-https stop-https start-all-https stop-all-https

install:
	@echo "Updating environment..."
	@echo "[backend] Installing dependencies with uv..."
	@cd backend && uv sync
	@echo "[frontend] Installing npm dependencies..."
	@npm --prefix frontend install
	@echo "Environment updated."

start:
	@echo "Starting services..."
	@echo "[backend] Starting FastAPI (http://localhost:8000)..."
	@cd backend && uv run uvicorn main:app --host 0.0.0.0 --port 8000 --reload > server.log 2>&1 & echo $$! > uvicorn.pid
	@echo "[frontend] Starting React dev server (http://localhost:3000)..."
	@npm --prefix frontend start > frontend/dev.log 2>&1 & echo $$! > frontend/dev.pid
	@echo "Services started. Backend: http://localhost:8000 | Frontend: http://localhost:3000"

gen-cert:
	@echo "Generating self-signed TLS certificate (valid for localhost)..."
	@mkdir -p certs
	@openssl req -x509 -nodes -newkey rsa:2048 -days 365 \
		-keyout certs/server.key -out certs/server.crt \
		-subj "/CN=localhost" \
		-addext "subjectAltName=DNS:localhost,IP:127.0.0.1" >/dev/null 2>&1 || true
	@echo "Certificate created: certs/server.crt, key: certs/server.key"

start-https:
	@echo "[backend] Starting FastAPI HTTPS (https://localhost:8443)..."
	@cd backend && uv run uvicorn main:app --host 0.0.0.0 --port 8443 --reload \
		--ssl-keyfile ../certs/server.key --ssl-certfile ../certs/server.crt \
		> server-https.log 2>&1 & echo $$! > uvicorn-https.pid
	@echo "HTTPS backend started at https://localhost:8443"

start-all-https:
	@echo "Generating TLS certificate if missing..."
	@[ -f certs/server.crt ] && [ -f certs/server.key ] || $(MAKE) gen-cert
	@echo "Starting backend (HTTPS :8443) and frontend (HTTPS :3000)..."
	@cd backend && uv run uvicorn main:app --host 0.0.0.0 --port 8443 --reload \
		--ssl-keyfile ../certs/server.key --ssl-certfile ../certs/server.crt \
		> server-https.log 2>&1 & echo $$! > uvicorn-https.pid
	@npm --prefix frontend run start-https > frontend/dev.log 2>&1 & echo $$! > frontend/dev.pid
	@echo "Services started. Backend: https://localhost:8443 | Frontend: https://localhost:3000"

stop:
	@echo "Stopping services..."
	@# Backend: try PID file first
	@[ ! -f backend/uvicorn.pid ] || (kill $$(cat backend/uvicorn.pid) 2>/dev/null || true)
	@rm -f backend/uvicorn.pid
	@# Backend: kill anything on port 8000 using fuser or lsof
	@{ command -v fuser >/dev/null && fuser -k 8000/tcp 2>/dev/null || true; } || { \
		command -v lsof >/dev/null && pids=$$(lsof -ti:8000 2>/dev/null || true); \
		[ -z "$$pids" ] || kill $$pids 2>/dev/null || true; \
	}
	@# Frontend: try PID file first
	@[ ! -f frontend/dev.pid ] || (kill $$(cat frontend/dev.pid) 2>/dev/null || true)
	@rm -f frontend/dev.pid
	@# Frontend: kill anything on port 3000 using fuser or lsof
	@{ command -v fuser >/dev/null && fuser -k 3000/tcp 2>/dev/null || true; } || { \
		command -v lsof >/dev/null && pids=$$(lsof -ti:3000 2>/dev/null || true); \
		[ -z "$$pids" ] || kill $$pids 2>/dev/null || true; \
	}
	@echo "Services stopped."

stop-https:
	@echo "Stopping HTTPS backend..."
	@[ ! -f backend/uvicorn-https.pid ] || (kill $$(cat backend/uvicorn-https.pid) 2>/dev/null || true)
	@rm -f backend/uvicorn-https.pid
	@{ command -v fuser >/dev/null && fuser -k 8443/tcp 2>/dev/null || true; } || { \
		command -v lsof >/dev/null && pids=$$(lsof -ti:8443 2>/dev/null || true); \
		[ -z "$$pids" ] || kill $$pids 2>/dev/null || true; \
	}
	@echo "HTTPS backend stopped."

stop-all-https:
	@echo "Stopping HTTPS backend and frontend..."
	@[ ! -f backend/uvicorn-https.pid ] || (kill $$(cat backend/uvicorn-https.pid) 2>/dev/null || true)
	@rm -f backend/uvicorn-https.pid
	@{ command -v fuser >/dev/null && fuser -k 8443/tcp 2>/dev/null || true; } || { \
		command -v lsof >/dev/null && pids=$$(lsof -ti:8443 2>/dev/null || true); \
		[ -z "$$pids" ] || kill $$pids 2>/dev/null || true; \
	}
	@[ ! -f frontend/dev.pid ] || (kill $$(cat frontend/dev.pid) 2>/dev/null || true)
	@rm -f frontend/dev.pid
	@{ command -v fuser >/dev/null && fuser -k 3000/tcp 2>/dev/null || true; } || { \
		command -v lsof >/dev/null && pids=$$(lsof -ti:3000 2>/dev/null || true); \
		[ -z "$$pids" ] || kill $$pids 2>/dev/null || true; \
	}
	@echo "HTTPS services stopped."

status:
	@echo "Service status:"
	@# Backend status
	@pid_b=""; [ -f backend/uvicorn.pid ] && pid_b=$$(cat backend/uvicorn.pid) || true; \
	if [ -n "$$pid_b" ] && ps -p $$pid_b >/dev/null 2>&1; then \
	  echo "[backend] RUNNING PID=$$pid_b URL=http://localhost:8000"; \
	else \
	  echo "[backend] STOPPED URL=http://localhost:8000"; \
	fi; \
	{ command -v lsof >/dev/null && ports_b=$$(lsof -ti:8000 2>/dev/null || true) && [ -n "$$ports_b" ] && echo "[backend] Port 8000 held by: $$ports_b" || true; } || \
	{ command -v ss >/dev/null && ss -tulpn 2>/dev/null | grep ':8000' || true; } || true
	@# Frontend status
	@pid_f=""; [ -f frontend/dev.pid ] && pid_f=$$(cat frontend/dev.pid) || true; \
	if [ -n "$$pid_f" ] && ps -p $$pid_f >/dev/null 2>&1; then \
	  echo "[frontend] RUNNING PID=$$pid_f URL=http://localhost:3000"; \
	else \
	  echo "[frontend] STOPPED URL=http://localhost:3000"; \
	fi; \
	{ command -v lsof >/dev/null && ports_f=$$(lsof -ti:3000 2>/dev/null || true) && [ -n "$$ports_f" ] && echo "[frontend] Port 3000 held by: $$ports_f" || true; } || \
	{ command -v ss >/dev/null && ss -tulpn 2>/dev/null | grep ':3000' || true; } || true
	@echo "Hint: tail -n 20 backend/server.log | tail -n 20; tail -n 20 frontend/dev.log | tail -n 20"
