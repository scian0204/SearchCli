#!/bin/bash

# SearchCli 빌드 스크립트

set -e

echo "🔨 SearchCli 빌드 중..."

# OS 와 아키텍처 감지
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
    darwin) GOOS="darwin" ;;
    linux) GOOS="linux" ;;
    msys_nt|mingw*) GOOS="windows" ;;
    *) echo "❌ Unsupported OS: $OS"; exit 1 ;;
esac

echo "🖥️ 타겟 플랫폼: $GOOS/$ARCH"

# Go 모듈 다운로드
echo "📦 의존성 다운로드 중..."
go mod download

# 이진 파일 빌드
echo "🏗️ 컴파일 중..."
go build -ldflags="-s -w" -o searchcli .

echo "✅ 빌드 완료!"
echo "📁 출력 파일: ./searchcli"
