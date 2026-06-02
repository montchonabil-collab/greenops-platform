#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/../services"
mkdir -p bin

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

go mod tidy
go mod download
go build -trimpath -ldflags="-s -w" -o bin/gateway ./gateway
go build -trimpath -ldflags="-s -w" -o bin/auth-service ./auth
go build -trimpath -ldflags="-s -w" -o bin/energy-service ./energy
go build -trimpath -ldflags="-s -w" -o bin/alerts-service ./alerts

ls -lh bin
