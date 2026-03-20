#!/bin/bash

# SearchCli 설치 스크립트 (macOS/Linux)

set -e

REPO="scian0204/SearchCli"

echo "🔍 SearchCli 설치 중..."

# 최신 버전 가져오기
VERSION=$(curl -s https://api.github.com/repos/${REPO}/releases/latest | grep '"tag_name"' | cut -d'"' -f4)

if [ -z "$VERSION" ]; then
    echo "❌ 버전 정보를 가져올 수 없습니다."
    exit 1
fi

echo "📦 버전: $VERSION"

# OS 와 아키텍처 감지
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
    darwin) OS="darwin" ;;
    linux) OS="linux" ;;
    *) echo "❌ Unsupported OS: $OS"; exit 1 ;;
esac

echo "🖥️ 플랫폼: $OS/$ARCH"

# 다운로드 URL (플랫폼별)
URL="https://github.com/${REPO}/releases/download/${VERSION}/searchcli_${OS}_${ARCH}"
TMP_FILE=$(mktemp)

echo "⬇️ 다운로드 중..."
curl -sL -o "$TMP_FILE" "$URL"

# 설치 경로
INSTALL_PATH="/usr/local/bin/searchcli"

if [ ! -d "/usr/local/bin" ]; then
    sudo mkdir -p /usr/local/bin
fi

echo "📁 설치 중..."
sudo mv "$TMP_FILE" "$INSTALL_PATH"
sudo chmod +x "$INSTALL_PATH"

echo "✅ 설치 완료!"
echo ""
echo "사용법:"
echo "  searchcli -q \"검색어\""
