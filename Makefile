# Makefile for project-kt3 (HTTPS simplified)
# Commands: start, stop, status

.PHONY: start stop status

start:
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
	@echo "Services started. Backend: https://localhost:8443 | Frontend: https://localhost:3000"


stop:
	@echo "Stopping services (3000,8443)..."
	@# Frontend 3000
	@[ ! -f frontend/dev.pid ] || (kill $$(cat frontend/dev.pid) 2>/dev/null || true)
	@rm -f frontend/dev.pid
	@{ command -v fuser >/dev/null && fuser -k 3000/tcp 2>/dev/null || true; } || { \
	    command -v lsof >/dev/null && pids=$$(lsof -ti:3000 2>/dev/null || true); \
	    [ -z "$$pids" ] || kill $$pids 2>/dev/null || true; \
	}
	@# Backend HTTPS 8443
	@[ ! -f backend/uvicorn-https.pid ] || (kill $$(cat backend/uvicorn-https.pid) 2>/dev/null || true)
	@rm -f backend/uvicorn-https.pid
	@{ command -v fuser >/dev/null && fuser -k 8443/tcp 2>/dev/null || true; } || { \
	    command -v lsof >/dev/null && pids=$$(lsof -ti:8443 2>/dev/null || true); \
	    [ -z "$$pids" ] || kill $$pids 2>/dev/null || true; \
	}
	@echo "All services stopped."

status:
	@echo "Service status (3000,8443):"
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
	@echo "Logs: backend/server-https.log | frontend/dev.log"
