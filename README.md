# Search CLI

Go 로 작성된 명령어 라인 기반 웹 검색 도구입니다. Bing 과 DuckDuckGo 검색 엔진을 지원하며, 검색 결과 페이지를 크롤링하여 상세 내용을 추출할 수 있습니다.

## 🎯 빠른 시작

```bash
# 1 줄로 설치 (macOS/Linux)
curl -sSf https://raw.githubusercontent.com/scian0204/SearchCli/main/install.sh | sh

# 바로 사용 (Bing 검색)
searchcli -q "go programming language"
```

## 기능

- **다중 검색 엔진 지원**: Bing, DuckDuckGo 검색 가능
- **HTML 파싱**: 검색 결과 페이지를 자동으로 파싱하여 구조화된 데이터 추출
- **웹 크롤링**: 검색 결과 링크를 방문하여 페이지 내용 추출
- **JSON 출력**: 결과를 JSON 형식으로 출력 또는 파일 저장
- **gzip 압축 지원**: 압축된 HTTP 응답 자동 처리

## 설치

### Homebrew (macOS/Linux)

```bash
brew tap scian0204/SearchCli https://github.com/scian0204/SearchCli.git
brew install searchcli
```

### 자동 설치 스크립트

#### macOS/Linux

```bash
curl -sSf https://raw.githubusercontent.com/scian0204/SearchCli/main/install.sh | sh
```

#### Windows (PowerShell)

```powershell
iwr https://raw.githubusercontent.com/scian0204/SearchCli/main/install.ps1 -UseBasicParsing | iex
```

### GitHub Releases

[릴리스 페이지](https://github.com/scian0204/SearchCli/releases) 에서 플랫폼에 맞는 이진 파일을 다운로드하세요.

### 소스에서 빌드

#### 사전 요구사항

- Go 1.16 이상
- 인터넷 연결

```bash
# 의존성 설치
go mod download

# 이진 파일 빌드
go build -o searchcli .

# 또는 빌드 스크립트 사용
./build.sh

# 또는 Makefile 사용
make build
```

### Makefile 명령

```bash
make build        # 현재 플랫폼용 빌드
make build-all    # 다중 플랫폼 빌드
make install      # 시스템에 설치 (/usr/local/bin)
make uninstall    # 시스템에서 제거
make clean        # 빌드 파일 정리
make test         # 테스트 실행
make help         # 도움말 표시
```

## 사용법

### 기본 검색

```bash
# Bing 으로 검색 (기본)
searchcli -q "go programming language"

# DuckDuckGo 로 검색
searchcli -q "go programming language" -engine ddg
```

### 링크 크롤링

검색 결과 페이지를 방문하여 상세 내용 추출:

```bash
# 최대 5 개의 링크 크롤링
searchcli -q "go programming language" -crawl -max-links 5

# 10 개의 링크 크롤링
searchcli -q "go programming language" -crawl -max-links 10
```

### 결과 저장

```bash
# JSON 파일로 저장
searchcli -q "go programming language" -output results.json

# 크롤링 후 파일 저장
searchcli -q "go programming language" -crawl -output results.json
```

### 도움말

```bash
searchcli -help
```

## 명령어 옵션

| 옵션 | 설명 | 기본값 |
|------|------|--------|
| `-q string` | 검색어 (필수) | - |
| `-engine string` | 검색 엔진: `bing`, `ddg` | `bing` |
| `-crawl` | 링크 크롤링 활성화 | `false` |
| `-max-links int` | 크롤링할 최대 링크 수 | `5` |
| `-output string` | 출력 파일 경로 (비우면 stdout) | - |
| `-help` | 도움말 표시 | `false` |

## 출력 형식

### JSON 구조

```json
{
  "search_info": {
    "query": "go programming language"
  },
  "results": [
    {
      "title": "The Go Programming Language",
      "link": "https://go.dev/",
      "snippet": "Go is an open source programming language...",
      "display_link": "go.dev",
      "crawled_content": {
        "title": "Go Programming Language",
        "description": "The Go programming language",
        "keywords": "go, programming, language",
        "headings": ["Introduction", "Installation"],
        "paragraphs": ["Go is a compiled language..."],
        "links": ["https://go.dev/doc", "https://go.dev/tutorial"],
        "source_url": "https://go.dev/"
      }
    }
  ]
}
```

### 필드 설명

| 필드 | 설명 |
|------|------|
| `search_info.query` | 검색어 |
| `results.title` | 결과 제목 |
| `results.link` | 결과 URL |
| `results.snippet` | 결과 미리보기 텍스트 |
| `results.display_link` | 표시용 도메인 이름 |
| `results.crawled_content` | 크롤링된 페이지 내용 (`-crawl` 옵션 사용 시) |

## 프로젝트 구조

```
SearchCli/
├── main.go           # CLI 엔트리 포인트, 검색 엔진 통합
├── models.go         # 데이터 구조 정의
├── search_parser.go  # 검색 결과 HTML/XML 파서
├── content_crawler.go # 웹 페이지 크롤러
├── build.sh          # 빌드 스크립트
├── install.sh        # macOS/Linux 설치 스크립트
├── install.ps1       # Windows 설치 스크립트
├── Makefile          # 빌드/설치 자동화
├── go.mod            # Go 모듈 정의
├── go.sum            # 의존성 체크섬
└── searchcli         # 컴파일된 이진 파일
```

### 파일별 역할

| 파일 | 설명 |
|------|------|
| `main.go` | CLI 파싱, 검색 엔진 선택, HTTP 요청 |
| `models.go` | SearchResult, Result, CrawledContent 구조체 |
| `search_parser.go` | DuckDuckGo/Bing HTML 파싱 로직 |
| `content_crawler.go` | 페이지 내용 추출 (제목, 본문, 링크 등) |
| `build.sh` | 크로스 플랫폼 빌드 스크립트 |
| `install.sh` | macOS/Linux 자동 설치 |
| `install.ps1` | Windows 자동 설치 |
| `Makefile` | 빌드/설치/테스트 자동화 |

## 기술 스택

- **언어**: Go (Golang)
- **HTML 파싱**: `golang.org/x/net/html`
- **HTTP 클라이언트**: 표준 라이브러리 `net/http`

## 예시

### Python 튜토리얼 검색

```bash
searchcli -q "python tutorial for beginners"
```

### 크롤링하여 블로그 내용 추출

```bash
searchcli -q "golang best practices" -crawl -max-links 3 -output golang-practices.json
```

### Bing 으로 특정 주제 검색

```bash
searchcli -q "machine learning basics" -engine bing -crawl
```

## 주의사항

- 일부 웹사이트는 크롤링을 제한할 수 있습니다
- 대량 크롤링 시 검색 엔진의 사용 정책을 준수하세요
- `-crawl` 옵션은 네트워크 속도에 따라 시간이 걸릴 수 있습니다

## 라이선스

MIT License
