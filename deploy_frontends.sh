#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# --- Configuration ---
SSH_USER="ubuntu"
SSH_HOST="172.24.125.63"
SSH_PASS="Bnu@fe2024"
MAIN_FRONTEND_PATH="~/www/kt3-main-frontend"
SUB_FRONTEND_PATH="~/www/kt3-sub-frontend"

# --- Build Main Frontend (React) ---
echo "Building main frontend (React)..."
cd /home/ubuntu/zlb/project-kt3/frontend
npm install
npm run build
echo "Main frontend built successfully."

# --- Build Sub Frontend (Vue) ---
echo "Building sub frontend (Vue)..."
cd /home/ubuntu/zlb/project-kt3/sub_project/sub_frontend
npm install
npm run build-only
echo "Sub frontend built successfully."

# --- Create remote directories and set permissions ---
echo "Creating remote directories and setting permissions on ${SSH_HOST}..."
sshpass -p "${SSH_PASS}" ssh -o StrictHostKeyChecking=no "${SSH_USER}@${SSH_HOST}" "mkdir -p ${MAIN_FRONTEND_PATH} && mkdir -p ${SUB_FRONTEND_PATH} && echo '${SSH_PASS}' | sudo -S chown -R ${SSH_USER}:${SSH_USER} ~/www"
echo "Remote directories ready."

# --- Deploy Main Frontend ---
echo "Deploying main frontend to ${SSH_HOST}..."
sshpass -p "${SSH_PASS}" rsync -rtlvz --delete /home/ubuntu/zlb/project-kt3/frontend/build/ ${SSH_USER}@${SSH_HOST}:${MAIN_FRONTEND_PATH}/
echo "Main frontend deployed successfully."

# --- Deploy Sub Frontend ---
echo "Deploying sub frontend to ${SSH_HOST}..."
sshpass -p "${SSH_PASS}" rsync -rtlvz --delete /home/ubuntu/zlb/project-kt3/sub_project/sub_frontend/dist/ ${SSH_USER}@${SSH_HOST}:${SUB_FRONTEND_PATH}/
echo "Sub frontend deployed successfully."

echo "Deployment complete!"