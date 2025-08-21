# Makefile for project-kt3
# One-click install, start, stop, and clean for backend (FastAPI) and frontend (React)

# Variables
PYTHON := python3
BACKEND_VENV := backend/.venv
PIP := $(BACKEND_VENV)/bin/pip
UVICORN := $(BACKEND_VENV)/bin/uvicorn
NODE := node
NPM := npm

# Phony targets
.PHONY: help install install-backend install-frontend start start-backend start-frontend stop stop-backend stop-frontend restart status logs clean

help:
	@echo "Available targets:"
	@echo "  make install          - Install backend and frontend dependencies"
	@echo "  make start            - Start backend (8000) and frontend (3000) in background"
	@echo "  make stop             - Stop backend and frontend dev servers"
	@echo "  make restart          - Restart both servers"
	@echo "  make status           - Show running PIDs and URLs"
	@echo "  make logs             - Tail both logs"
	@echo "  make clean            - Stop servers and remove caches/build artifacts"

install: install-backend install-frontend

install-backend:
	@echo "[backend] Creating venv and installing dependencies..."
	@test -d $(BACKEND_VENV) || $(PYTHON) -m venv $(BACKEND_VENV)
	@$(PIP) install --upgrade pip
	@$(PIP) install -r backend/requirements.txt
	@echo "[backend] Done."

install-frontend:
	@echo "[frontend] Installing npm dependencies..."
	@$(NPM) --prefix frontend install
	@echo "[frontend] Done."

start: start-backend start-frontend status

start-backend:
	@echo "[backend] Starting FastAPI (http://localhost:8000)..."
	@mkdir -p backend
	@rm -f backend/uvicorn.pid
	@$(UVICORN) backend.main:app --host 0.0.0.0 --port 8000 --reload > backend/server.log 2>&1 & echo $$! > backend/uvicorn.pid
	@echo "[backend] PID: $$(cat backend/uvicorn.pid) | Logs: backend/server.log"

start-frontend:
	@echo "[frontend] Starting React dev server (http://localhost:3000)..."
	@mkdir -p frontend
	@rm -f frontend/dev.pid
	@$(NPM) --prefix frontend start > frontend/dev.log 2>&1 & echo $$! > frontend/dev.pid
	@echo "[frontend] PID: $$(cat frontend/dev.pid) | Logs: frontend/dev.log"

stop: stop-backend stop-frontend

stop-backend:
	@echo "[backend] Stopping..."
	@# Try kill by PID file first (graceful)
	@[ ! -f backend/uvicorn.pid ] || (kill $$(cat backend/uvicorn.pid) 2>/dev/null || true)
	@sleep 0.5
	@# Kill anything listening on 8000 (graceful then force), using lsof or fuser as fallback
	@{ command -v lsof >/dev/null && ( \
		pids=$$(lsof -ti:8000 2>/dev/null || true); \
		if [ -n "$$pids" ]; then kill $$pids 2>/dev/null || true; fi; \
		sleep 0.5; \
		pids2=$$(lsof -ti:8000 2>/dev/null || true); \
		[ -z "$$pids2" ] || kill -9 $$pids2 2>/dev/null || true \
	) || { command -v fuser >/dev/null && fuser -k 8000/tcp 2>/dev/null || true; }; }
	@rm -f backend/uvicorn.pid
	@echo "[backend] Stopped."

stop-frontend:
	@echo "[frontend] Stopping..."
	@# Try kill by PID file first (graceful)
	@[ ! -f frontend/dev.pid ] || (kill $$(cat frontend/dev.pid) 2>/dev/null || true)
	@sleep 0.5
	@# Kill anything listening on 3000 (graceful then force), using lsof or fuser as fallback
	@{ command -v lsof >/dev/null && ( \
		pids=$$(lsof -ti:3000 2>/dev/null || true); \
		if [ -n "$$pids" ]; then kill $$pids 2>/dev/null || true; fi; \
		sleep 0.5; \
		pids2=$$(lsof -ti:3000 2>/dev/null || true); \
		[ -z "$$pids2" ] || kill -9 $$pids2 2>/dev/null || true \
	) || { command -v fuser >/dev/null && fuser -k 3000/tcp 2>/dev/null || true; }; }
	@rm -f frontend/dev.pid
	@echo "[frontend] Stopped."

restart: stop start

status:
	@echo "=== Status ==="
	@echo "Backend:  PID=$$(cat backend/uvicorn.pid 2>/dev/null || echo '-')  URL=http://localhost:8000"
	@echo "Frontend: PID=$$(cat frontend/dev.pid 2>/dev/null || echo '-')  URL=http://localhost:3000"

logs:
	@echo "Tailing logs (Ctrl-C to stop)..."
	@touch backend/server.log frontend/dev.log
	@tail -n 50 -f backend/server.log frontend/dev.log

clean: stop
	@echo "[clean] Removing caches and build artifacts..."
	@find backend -type d -name __pycache__ -exec rm -rf {} + 2>/dev/null || true
	@find . -type d -name .pytest_cache -exec rm -rf {} + 2>/dev/null || true
	@rm -rf frontend/build 2>/dev/null || true
	@rm -f backend/server.log frontend/dev.log backend/uvicorn.pid frontend/dev.pid
	@echo "[clean] Done."
