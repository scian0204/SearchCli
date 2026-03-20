.PHONY: build clean install uninstall test help

# 변수
BINARY_NAME=searchcli
VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "v1.0.0")
LDFLAGS=-ldflags="-s -w -X main.Version=$(VERSION)"

# 기본 타겟
all: build

# 현재 플랫폼용 빌드
build:
	@echo "🔨 빌드 중..."
	go mod download
	go build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "✅ 빌드 완료: ./$(BINARY_NAME)"

# 다중 플랫폼 빌드
build-all:
	@echo "🔨 다중 플랫폼 빌드 중..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o searchcli_darwin_amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o searchcli_darwin_arm64 .
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o searchcli_linux_amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o searchcli_linux_arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o searchcli_windows_amd64.exe .
	@echo "✅ 다중 플랫폼 빌드 완료"

# macOS/Linux 설치
install: build
	@echo "📁 설치 중..."
	@if [ -d "/usr/local/bin" ]; then \
		sudo cp $(BINARY_NAME) /usr/local/bin/ && sudo chmod +x /usr/local/bin/$(BINARY_NAME); \
	else \
		sudo mkdir -p /usr/local/bin && \
		sudo cp $(BINARY_NAME) /usr/local/bin/ && sudo chmod +x /usr/local/bin/$(BINARY_NAME); \
	fi
	@echo "✅ 설치 완료"

# macOS/Linux 제거
uninstall:
	@echo "🗑️ 제거 중..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ 제거 완료"

# 빌드 파일 정리
clean:
	@echo "🧹 정리 중..."
	rm -f $(BINARY_NAME) searchcli_*
	@echo "✅ 정리 완료"

# 테스트 실행
test:
	@echo "🧪 테스트 중..."
	go test -v ./...

# 도움말
help:
	@echo "SearchCli Makefile"
	@echo ""
	@echo "사용 가능한 명령:"
	@echo "  make build       - 현재 플랫폼용 빌드"
	@echo "  make build-all   - 다중 플랫폼 빌드"
	@echo "  make install     - 시스템에 설치 (/usr/local/bin)"
	@echo "  make uninstall   - 시스템에서 제거"
	@echo "  make clean       - 빌드 파일 정리"
	@echo "  make test        - 테스트 실행"
	@echo "  make help        - 이 도움말 표시"
