#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SERVICES_DIR="$ROOT_DIR/services"

"$ROOT_DIR/scripts/build-linux.sh"

docker build -t greenops-gateway:latest -f "$SERVICES_DIR/gateway/Dockerfile" "$SERVICES_DIR"
docker build -t greenops-auth-service:latest -f "$SERVICES_DIR/auth/Dockerfile" "$SERVICES_DIR"
docker build -t greenops-energy-service:latest -f "$SERVICES_DIR/energy/Dockerfile" "$SERVICES_DIR"
docker build -t greenops-alerts-service:latest -f "$SERVICES_DIR/alerts/Dockerfile" "$SERVICES_DIR"

echo "Images Kubernetes pretes :"
docker image ls "greenops-*"
