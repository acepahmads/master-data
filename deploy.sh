#!/bin/bash

# ==============================================================================
# IoT Product R&D Control Center - Automated Deployment Script (Ubuntu Linux)
# ==============================================================================

set -e # Exit immediately if a command exits with a non-zero status

# Color output formatting
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}======================================================${NC}"
echo -e "${BLUE} Starting IoT Product R&D Control Center Deployment... ${NC}"
echo -e "${BLUE}======================================================${NC}"

# 1. Pull Latest Code from Git
echo -e "\n${YELLOW}[1/4] Pulling latest code from Git repository...${NC}"
git pull

# 2. Build Frontend
echo -e "\n${YELLOW}[2/4] Installing dependencies & Building Frontend (Vue + Vite)...${NC}"
npm install --silent
npm run build

# 3. Build Backend Binary (Go)
echo -e "\n${YELLOW}[3/4] Compiling Go Backend for Linux...${NC}"
cd backend
mkdir -p uploads
go mod tidy
go build -o server ./cmd/main.go
chmod +x server
cd ..

# Ensure root uploads folder or permissions if needed
mkdir -p uploads
mkdir -p dist

# 4. Restart Backend Service / Process
echo -e "\n${YELLOW}[4/4] Restarting Application Service...${NC}"
if systemctl is-active --quiet iot-backend; then
    sudo systemctl restart iot-backend
    echo -e "${GREEN}✓ Systemd service 'iot-backend' successfully restarted.${NC}"
else
    echo -e "${YELLOW}! Systemd service 'iot-backend' is not currently active.${NC}"
    echo -e "  If using systemd, run: sudo systemctl start iot-backend"
    echo -e "  Or run binary manually: cd backend && ./server"
fi

# Optional Nginx reload
if systemctl is-active --quiet nginx; then
    sudo systemctl reload nginx
    echo -e "${GREEN}✓ Nginx successfully reloaded.${NC}"
fi

echo -e "\n${GREEN}======================================================${NC}"
echo -e "${GREEN} Deployment Completed Successfully! ${NC}"
echo -e "${GREEN}======================================================${NC}\n"
