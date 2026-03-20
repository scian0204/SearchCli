# SearchCli 설치 스크립트 (Windows)

Write-Host "🔍 SearchCli 설치 중..." -ForegroundColor Cyan

$REPO = "scian0204/SearchCli"

# 최신 버전 가져오기
try {
    $latestRelease = Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest"
    $VERSION = $latestRelease.tag_name
} catch {
    Write-Host "❌ 버전 정보를 가져올 수 없습니다." -ForegroundColor Red
    exit 1
}

Write-Host "📦 버전: $VERSION" -ForegroundColor Green

# 아키텍처 감지
$ARCH = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

Write-Host "🖥️ 플랫폼: windows/$ARCH" -ForegroundColor Yellow

# 다운로드 URL
$URL = "https://github.com/$REPO/releases/download/$VERSION/searchcli_windows_amd64.exe"
$TMP_FILE = Join-Path $env:TEMP "searchcli.exe"

# 설치 경로 (사용자별 또는 전역)
$INSTALL_PATH = "$env:LOCALAPPDATA\searchcli\searchcli.exe"

Write-Host "⬇️ 다운로드 중..." -ForegroundColor Cyan
try {
    Invoke-WebRequest -Uri $URL -OutFile $TMP_FILE -UseBasicParsing
} catch {
    # amd64 실패 시 386 시도
    $URL = "https://github.com/$REPO/releases/download/$VERSION/searchcli_windows_386.exe"
    Write-Host "⬇️ 32bit 버전 다운로드 중..." -ForegroundColor Cyan
    try {
        Invoke-WebRequest -Uri $URL -OutFile $TMP_FILE -UseBasicParsing
    } catch {
        Write-Host "❌ 다운로드 실패" -ForegroundColor Red
        exit 1
    }
}

# 설치 디렉토리 생성
if (-not (Test-Path $env:LOCALAPPDATA\searchcli)) {
    New-Item -ItemType Directory -Path $env:LOCALAPPDATA\searchcli | Out-Null
}

# 설치
Write-Host "📁 설치 중..." -ForegroundColor Cyan
Copy-Item $TMP_FILE -Destination $INSTALL_PATH -Force

# PATH 에 추가 (현재 사용자)
$currentUserPath = [Environment]::GetEnvironmentVariable("Path", "User")
$searchcliPath = "$env:LOCALAPPDATA\searchcli"

if ($currentUserPath -notlike "*$searchcliPath*") {
    $newPath = "$searchcliPath;$currentUserPath"
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    # 현재 셸에도 적용
    $env:Path = "$searchcliPath;" + $env:Path
    Write-Host "🔄 PATH 에 추가됨" -ForegroundColor Yellow
}

# 임시 파일 삭제
Remove-Item $TMP_FILE -Force

Write-Host ""
Write-Host "✅ 설치 완료!" -ForegroundColor Green
Write-Host ""
Write-Host "사용법:" -ForegroundColor Cyan
Write-Host "  searchcli -q `"검색어`""
Write-Host ""
Write-Host "참고: 설치된 경로 - $INSTALL_PATH"
