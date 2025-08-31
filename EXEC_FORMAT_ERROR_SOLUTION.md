# PIXELZX POS EVM 체인 - Exec Format Error 문제 해결 가이드

## 문제 설명

PIXELZX POS EVM 체인을 Docker로 실행할 때 다음과 같은 오류가 발생하는 경우가 있습니다:

```
exec /usr/local/bin/pixelzx: exec format error
```

이 오류는 Docker 컨테이너의 아키텍처와 호스트 시스템의 아키텍처가 일치하지 않을 때 발생합니다.

## 원인 분석

### 일반적인 발생 시나리오

| 빌드 환경 | 실행 환경 | 결과 | 원인 |
|----------|----------|------|------|
| AMD64 | AMD64 | ✅ 정상 | 아키텍처 일치 |
| AMD64 | ARM64 | ❌ exec format error | 아키텍처 불일치 |
| AMD64 | ARM/v7 | ❌ exec format error | 아키텍처 불일치 |

## 해결 방법

### 1. 즉시 해결 방안 - 플랫폼 명시적 지정

가장 빠른 해결 방법은 Docker 실행 시 플랫폼을 명시적으로 지정하는 것입니다:

#### AMD64 환경 (Intel/AMD PC)
```bash
docker run --platform linux/amd64 --rm yuchanshin/pixelzx-evm:latest
```

#### ARM64 환경 (Apple Silicon Mac, Raspberry Pi 4 이상)
```bash
docker run --platform linux/arm64 --rm yuchanshin/pixelzx-evm:latest
```

#### ARM/v7 환경 (Raspberry Pi 3 이하)
```bash
docker run --platform linux/arm/v7 --rm yuchanshin/pixelzx-evm:latest
```

### 2. 플랫폼 자동 감지 스크립트 사용

프로젝트에 포함된 플랫폼 자동 감지 스크립트를 사용하면 시스템에 맞는 플랫폼을 자동으로 감지하여 실행할 수 있습니다:

```bash
# 기본 실행 (도움말 표시)
./scripts/detect-platform.sh

# 특정 명령어 실행
./scripts/detect-platform.sh admin status
./scripts/detect-platform.sh init
./scripts/detect-platform.sh start
```

### 3. Docker Compose 사용 시

docker-compose.yml 파일에 플랫폼 정보를 명시적으로 추가할 수 있습니다:

```yaml
version: '3.8'
services:
  pixelzx-node:
    image: yuchanshin/pixelzx-evm:latest
    platform: linux/arm64  # 또는 linux/amd64, linux/arm/v7
    ports:
      - "8545:8545"
      - "8546:8546"
      - "30303:30303"
    volumes:
      - ./data:/app/data
```

## 근본적 해결 - 멀티 아키텍처 이미지 사용

PIXELZX POS EVM 체인 이미지는 이제 멀티 아키텍처를 지원합니다. Docker Hub의 최신 이미지(yuchanshin/pixelzx-evm:latest)는 다음 플랫폼을 모두 지원합니다:

- linux/amd64 (Intel/AMD 64비트)
- linux/arm64 (Apple Silicon, Raspberry Pi 4 이상)
- linux/arm/v7 (Raspberry Pi 3 이하)

Docker는 자동으로 호스트 시스템에 맞는 이미지를 선택하여 실행합니다.

## 진단 및 검증

### 시스템 아키텍처 확인

현재 시스템의 아키텍처를 확인하려면 다음 명령어를 사용합니다:

```bash
uname -m
```

결과:
- `x86_64` 또는 `amd64`: AMD64 아키텍처
- `aarch64` 또는 `arm64`: ARM64 아키텍처
- `armv7l`: ARM/v7 아키텍처

### Docker 이미지 매니페스트 확인

Docker 이미지가 지원하는 플랫폼을 확인하려면 다음 명령어를 사용합니다:

```bash
docker buildx imagetools inspect yuchanshin/pixelzx-evm:latest
```

## 문제 해결 체크리스트

1. [ ] 시스템 아키텍처 확인 (`uname -m`)
2. [ ] Docker 버전 확인 (20.10 이상 필요)
3. [ ] Docker Buildx 설치 여부 확인
4. [ ] 플랫폼 명시적 지정으로 실행 시도
5. [ ] 플랫폼 자동 감지 스크립트 사용
6. [ ] Docker 이미지 매니페스트 확인

## 추가 정보

### Docker 버전 업데이트

Docker 버전이 낮은 경우 exec format error가 발생할 수 있습니다. Docker를 최신 버전으로 업데이트하세요:

```bash
# Docker Desktop (Mac/Windows)
# https://www.docker.com/products/docker-desktop 에서 다운로드

# Ubuntu/Debian
sudo apt update
sudo apt install docker-ce docker-ce-cli containerd.io

# CentOS/RHEL
sudo yum update
sudo yum install docker-ce docker-ce-cli containerd.io
```

### Buildx 활성화

Docker Buildx가 설치되어 있지 않은 경우 다음 명령어로 설치합니다:

```bash
# Docker Buildx 설치
mkdir -p ~/.docker/cli-plugins
wget -qO ~/.docker/cli-plugins/docker-buildx https://github.com/docker/buildx/releases/download/v0.11.2/buildx-v0.11.2.linux-amd64
chmod +x ~/.docker/cli-plugins/docker-buildx

# 설치 확인
docker buildx version
```

## 문의 및 지원

문제가 지속되는 경우 다음 정보를 포함하여 문의해 주세요:

1. 운영체제 및 버전
2. 시스템 아키텍처 (`uname -m` 결과)
3. Docker 버전 (`docker version` 결과)
4. 발생하는 정확한 오류 메시지
5. 실행한 명령어

GitHub Issues: [링크]
이메일: [이메일 주소]