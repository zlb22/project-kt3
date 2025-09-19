#!/usr/bin/env bash
set -euo pipefail

# One-click deployment script for project-kt3
# This script handles: build, package, nginx config, and deployment
# 
# Usage:
#   ./deploy/scripts/deploy.sh [--local|--remote]
#   
# Environment variables:
#   REMOTE_USER=ubuntu
#   REMOTE_HOST=172.24.125.63
#   SSH_PORT=22 (optional)
#   BACKEND_HOST=172.24.130.213 (optional, default from current setup)

usage() {
  cat <<EOF
Usage: $(basename "$0") [OPTIONS]

OPTIONS:
  --local     Deploy locally (for development/testing)
  --remote    Deploy to remote server (production)
  --help      Show this help

Environment variables for remote deployment:
  REMOTE_USER           SSH user (required for --remote)
  REMOTE_HOST           SSH host (required for --remote) 
  SSH_PORT              SSH port (default: 22)
  BACKEND_HOST          Backend server IP (default: 172.24.130.213)

Examples:
  # Local deployment
  ./deploy/scripts/deploy.sh --local
  
  # Remote deployment
  REMOTE_USER=ubuntu REMOTE_HOST=172.24.125.63 ./deploy/scripts/deploy.sh --remote
EOF
}

# Default values
DEPLOY_MODE=""
BACKEND_HOST=${BACKEND_HOST:-172.24.130.213}
SSH_PORT=${SSH_PORT:-22}

# Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --local)
      DEPLOY_MODE="local"
      shift
      ;;
    --remote)
      DEPLOY_MODE="remote"
      shift
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      usage
      exit 1
      ;;
  esac
done

if [[ -z "$DEPLOY_MODE" ]]; then
  echo "Error: Must specify --local or --remote"
  usage
  exit 1
fi

# Validate remote deployment requirements
if [[ "$DEPLOY_MODE" == "remote" ]]; then
  : "${REMOTE_USER:?REMOTE_USER is required for remote deployment}"
  : "${REMOTE_HOST:?REMOTE_HOST is required for remote deployment}"
fi

ROOT_DIR=$(cd -- "$(dirname -- "$0")"/../.. && pwd)
ART_DIR="$ROOT_DIR/dist_artifacts"
NGINX_CONF_SRC="$ROOT_DIR/deploy/nginx/kt3.conf"

echo "=== Project-KT3 One-Click Deployment ==="
echo "Mode: $DEPLOY_MODE"
echo "Root: $ROOT_DIR"
echo "Backend Host: $BACKEND_HOST"
if [[ "$DEPLOY_MODE" == "remote" ]]; then
  echo "Remote: ${REMOTE_USER}@${REMOTE_HOST}:${SSH_PORT}"
fi
echo

# Step 1: Clean and prepare
echo "[1/6] Cleaning previous artifacts..."
rm -rf "$ART_DIR"
mkdir -p "$ART_DIR"

# Step 2: Build frontend applications
echo "[2/6] Building frontend applications..."

echo "  Building React frontend..."
pushd "$ROOT_DIR/frontend" >/dev/null
if [ -f package-lock.json ]; then
  npm ci --silent
else
  npm install --silent
fi
PUBLIC_URL=/topic-three npm run build
popd >/dev/null

echo "  Building Vue sub-frontend..."
pushd "$ROOT_DIR/sub_project/sub_frontend" >/dev/null
if [ -f package-lock.json ]; then
  npm ci --silent
else
  npm install --silent
fi
npm run build-only
popd >/dev/null

# Step 3: Package artifacts
echo "[3/6] Packaging build artifacts..."
tar czf "$ART_DIR/topic-three_build.tgz" -C "$ROOT_DIR/frontend/build" .
tar czf "$ART_DIR/topic-three-online-experiment.tgz" -C "$ROOT_DIR/sub_project/sub_frontend/dist" .

echo "  Artifacts created:"
echo "    - topic-three_build.tgz"
echo "    - topic-three-online-experiment.tgz"

# Step 4: Update nginx configuration with current backend host
echo "[4/6] Updating nginx configuration..."
sed "s/172\.24\.130\.213/$BACKEND_HOST/g" "$ROOT_DIR/deploy/nginx/kt3.conf.template" > "$NGINX_CONF_SRC" 2>/dev/null || {
  # If template doesn't exist, update the current config
  sed -i.bak "s/172\.24\.130\.213/$BACKEND_HOST/g" "$NGINX_CONF_SRC"
}

# Step 5: Deploy based on mode
if [[ "$DEPLOY_MODE" == "local" ]]; then
  echo "[5/6] Deploying locally..."
  
  # Local web directories
  LOCAL_WEB_DIR_TOPIC="/var/www/topic-three"
  LOCAL_WEB_DIR_SUB="/var/www/topic-three-online-experiment"
  LOCAL_NGINX_CONF="/etc/nginx/conf.d/kt3.conf"
  
  # Create directories
  sudo mkdir -p "$LOCAL_WEB_DIR_TOPIC" "$LOCAL_WEB_DIR_SUB"
  
  # Backup existing deployments
  TS=$(date +%Y%m%d_%H%M%S)
  if [ -d "$LOCAL_WEB_DIR_TOPIC" ] && [ "$(ls -A "$LOCAL_WEB_DIR_TOPIC" 2>/dev/null | wc -l)" -gt 0 ]; then
    echo "  Backing up existing topic-three deployment..."
    sudo tar czf "${LOCAL_WEB_DIR_TOPIC}.bak_${TS}.tgz" -C "$LOCAL_WEB_DIR_TOPIC" . || true
  fi
  if [ -d "$LOCAL_WEB_DIR_SUB" ] && [ "$(ls -A "$LOCAL_WEB_DIR_SUB" 2>/dev/null | wc -l)" -gt 0 ]; then
    echo "  Backing up existing sub-frontend deployment..."
    sudo tar czf "${LOCAL_WEB_DIR_SUB}.bak_${TS}.tgz" -C "$LOCAL_WEB_DIR_SUB" . || true
  fi
  
  # Clean and extract
  sudo rm -rf "${LOCAL_WEB_DIR_TOPIC}"/* "${LOCAL_WEB_DIR_SUB}"/*
  sudo tar xzf "$ART_DIR/topic-three_build.tgz" -C "$LOCAL_WEB_DIR_TOPIC"
  sudo tar xzf "$ART_DIR/topic-three-online-experiment.tgz" -C "$LOCAL_WEB_DIR_SUB"
  
  # Set permissions
  sudo find "$LOCAL_WEB_DIR_TOPIC" -type d -exec chmod 755 {} +
  sudo find "$LOCAL_WEB_DIR_TOPIC" -type f -exec chmod 644 {} +
  sudo find "$LOCAL_WEB_DIR_SUB" -type d -exec chmod 755 {} +
  sudo find "$LOCAL_WEB_DIR_SUB" -type f -exec chmod 644 {} +
  
  # Update nginx config
  echo "  Updating nginx configuration..."
  sudo cp "$NGINX_CONF_SRC" "$LOCAL_NGINX_CONF"
  
else
  echo "[5/6] Deploying to remote server..."
  
  REMOTE_UPLOAD_DIR="~/.kt3_uploads"
  REMOTE_WEB_DIR_TOPIC="/var/www/topic-three"
  REMOTE_WEB_DIR_SUB="/var/www/topic-three-online-experiment"
  REMOTE_NGINX_CONF="/etc/nginx/conf.d/kt3.conf"
  
  # Create remote upload directory and upload artifacts
  echo "  Uploading artifacts..."
  ssh -p "$SSH_PORT" "${REMOTE_USER}@${REMOTE_HOST}" "mkdir -p $REMOTE_UPLOAD_DIR"
  scp -P "$SSH_PORT" "$ART_DIR/topic-three_build.tgz" "$ART_DIR/topic-three-online-experiment.tgz" \
      "${REMOTE_USER}@${REMOTE_HOST}:$REMOTE_UPLOAD_DIR/"
  
  # Upload nginx config
  echo "  Uploading nginx configuration..."
  scp -P "$SSH_PORT" "$NGINX_CONF_SRC" "${REMOTE_USER}@${REMOTE_HOST}:/tmp/kt3.conf"
  
  # Execute remote deployment
  echo "  Executing remote deployment..."
  cat <<EOSSH | ssh -p "$SSH_PORT" "${REMOTE_USER}@${REMOTE_HOST}" bash -s -- "$REMOTE_UPLOAD_DIR" "$REMOTE_WEB_DIR_TOPIC" "$REMOTE_WEB_DIR_SUB" "$REMOTE_NGINX_CONF"
set -euo pipefail

UPLOAD_DIR="\$1"
WEB_DIR_TOPIC="\$2"
WEB_DIR_SUB="\$3"
NGINX_CONF="\$4"

# Expand tilde
UPLOAD_DIR_EXPANDED=\$(eval echo "\$UPLOAD_DIR")

echo "  Creating web directories..."
sudo mkdir -p "\$WEB_DIR_TOPIC" "\$WEB_DIR_SUB"

# Backup existing deployments
TS=\$(date +%Y%m%d_%H%M%S)
if [ -d "\$WEB_DIR_TOPIC" ] && [ "\$(ls -A "\$WEB_DIR_TOPIC" 2>/dev/null | wc -l)" -gt 0 ]; then
  echo "  Backing up existing topic-three deployment..."
  sudo tar czf "\${WEB_DIR_TOPIC}.bak_\${TS}.tgz" -C "\$WEB_DIR_TOPIC" . || true
fi
if [ -d "\$WEB_DIR_SUB" ] && [ "\$(ls -A "\$WEB_DIR_SUB" 2>/dev/null | wc -l)" -gt 0 ]; then
  echo "  Backing up existing sub-frontend deployment..."
  sudo tar czf "\${WEB_DIR_SUB}.bak_\${TS}.tgz" -C "\$WEB_DIR_SUB" . || true
fi

# Clean and extract
echo "  Extracting new deployments..."
sudo rm -rf "\${WEB_DIR_TOPIC}"/* "\${WEB_DIR_SUB}"/*
sudo tar xzf "\$UPLOAD_DIR_EXPANDED/topic-three_build.tgz" -C "\$WEB_DIR_TOPIC"
sudo tar xzf "\$UPLOAD_DIR_EXPANDED/topic-three-online-experiment.tgz" -C "\$WEB_DIR_SUB"

# Set permissions
sudo find "\$WEB_DIR_TOPIC" -type d -exec chmod 755 {} +
sudo find "\$WEB_DIR_TOPIC" -type f -exec chmod 644 {} +
sudo find "\$WEB_DIR_SUB" -type d -exec chmod 755 {} +
sudo find "\$WEB_DIR_SUB" -type f -exec chmod 644 {} +

# Update nginx config
echo "  Installing nginx configuration..."
sudo mkdir -p "\$(dirname "\$NGINX_CONF")"
sudo mv /tmp/kt3.conf "\$NGINX_CONF"

echo "Remote deployment completed."
EOSSH

fi

# Step 6: Test and reload nginx
echo "[6/6] Testing and reloading nginx..."
if [[ "$DEPLOY_MODE" == "local" ]]; then
  sudo nginx -t
  sudo systemctl reload nginx
  echo "  Nginx reloaded successfully."
else
  ssh -p "$SSH_PORT" "${REMOTE_USER}@${REMOTE_HOST}" "sudo nginx -t && sudo systemctl reload nginx"
  echo "  Remote nginx reloaded successfully."
fi

echo
echo "=== Deployment Complete ==="
if [[ "$DEPLOY_MODE" == "local" ]]; then
  echo "Local URLs:"
  echo "  Main App: https://localhost/topic-three/"
  echo "  Sub App:  https://localhost/topic-three/online-experiment/"
  echo "  API:      https://localhost/api/"
else
  echo "Remote URLs:"
  echo "  Main App: https://${REMOTE_HOST}/topic-three/"
  echo "  Sub App:  https://${REMOTE_HOST}/topic-three/online-experiment/"
  echo "  API:      https://${REMOTE_HOST}/api/"
fi
echo
echo "Artifacts saved in: $ART_DIR"
echo "Nginx config: $NGINX_CONF_SRC"
echo
