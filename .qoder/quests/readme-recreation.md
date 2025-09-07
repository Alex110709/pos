# PIXELZX POS EVM Chain

PIXELZX는 이더리움 가상 머신(EVM)과 호환되는 증명금지(PoS) 기반 블록체인 네트워크입니다. PIXELZX(PXZ)는 이 네트워크의 네이티브 토큰으로, 높은 성능(3초 블록 시간, 1000+ TPS)과 낮은 수수료를 제공합니다.

## 목차

1. [주요 특징](#주요-특징)
2. [기술 스택](#기술-스택)
3. [빠른 시작](#빠른-시작)
   - [소스 코드로부터 빌드](#소스-코드로부터-빌드)
   - [Docker를 사용한 실행](#docker를-사용한-실행)
   - [Docker Compose를 사용한 멀티 노드 설정](#docker-compose를-사용한-멀티-노드-설정)
4. [CLI 명령어](#cli-명령어)
5. [Docker 이미지](#docker-이미지)
6. [구성](#구성)
7. [포트](#포트)
8. [볼륨](#볼륨)
9. [환경 변수](#환경-변수)
10. [문제 해결](#문제-해결)

## 주요 특징

- **EVM 호환성**: 이더리움 스마트 계약과 완전히 호환됩니다.
- **PoS 합의 메커니즘**: 에너지 효율적인 증명금지 합의 알고리즘을 사용합니다.
- **고성능**: 3초 블록 시간과 1000+ TPS 처리량을 제공합니다.
- **낮은 수수료**: 최적화된 가스 수수료 구조로 낮은 트랜잭션 비용을 제공합니다.
- **멀티 아키텍처 지원**: linux/amd64, linux/arm64, linux/arm/v7 아키텍처를 지원합니다.

## 기술 스택

- **언어**: Go 1.21+
- **프레임워크**: Ethereum/go-ethereum 라이브러리
- **CLI 프레임워크**: Cobra
- **컨테이너화**: Docker, Docker Buildx
- **의존성 관리**: Go Modules

## 빠른 시작

### 소스 코드로부터 빌드

PIXELZX CLI를 빌드하려면 다음 명령어를 실행하세요:

```shell
make pixelzx
```

빌드가 완료되면 실행 파일은 `build/bin/pixelzx`에 생성됩니다.

### 네트워크 초기화

```shell
./build/bin/pixelzx init
```

### 노드 시작

```shell
./build/bin/pixelzx start
```

### 검증자 노드 시작

```shell
./build/bin/pixelzx start --validator
```

### Docker를 사용한 실행

Docker를 사용하여 PIXELZX를 실행할 수도 있습니다:

```shell
# Docker 이미지 빌드
make docker-build

# PIXELZX 노드 실행
docker run -d \
  --name pixelzx-node \
  -p 8545:8545 \
  -p 8546:8546 \
  -p 30303:30303 \
  -v $(pwd)/data:/app/data \
  pixelzx-node:latest
```

### Docker Compose를 사용한 멀티 노드 설정

멀티 노드 환경을 구성하려면 Docker Compose를 사용할 수 있습니다:

```shell
# 프로덕션 환경 시작
make compose-up

# 개발 환경 시작
make compose-dev-up
```

## CLI 명령어

PIXELZX CLI는 다양한 명령어를 제공하여 블록체인 네트워크를 관리할 수 있습니다:

- `pixelzx init` - 블록체인 네트워크 초기화
- `pixelzx start` - PIXELZX 노드 시작
- `pixelzx start --validator` - 검증자 모드로 PIXELZX 노드 시작
- `pixelzx account` - 계정 관리
- `pixelzx admin` - 관리자 명령어
- `pixelzx validator` - 검증자 관리
- `pixelzx staking` - 스테이킹 관리

각 명령어는 추가적인 하위 명령어를 가질 수 있습니다. 예를 들어, `pixelzx account --help`를 실행하여 계정 관련 명령어를 확인할 수 있습니다.

## Docker 이미지

Docker Hub에서 제공하는 PIXELZX 이미지:

- `yuchanshin/pixelzx-evm:latest` - 최신 안정 릴리스
- `yuchanshin/pixelzx-evm:develop` - 최신 개발 빌드

### 이미지 가져오기

```bash
docker pull yuchanshin/pixelzx-evm:latest
```

### PIXELZX 노드 실행

```bash
docker run -d \
  --name pixelzx-node \
  -p 8545:8545 \
  -p 8546:8546 \
  -p 30303:30303 \
  -v $(pwd)/data:/app/data \
  yuchanshin/pixelzx-evm:latest
```

### 네트워크 초기화

```bash
docker run --rm \
  -v $(pwd)/data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx init
```

### 검증자 모드로 실행

```bash
docker run -d \
  --name pixelzx-validator \
  -p 8545:8545 \
  -p 8546:8546 \
  -p 30303:30303 \
  -v $(pwd)/data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --validator
```

## 구성

사용자 정의 구성 파일을 마운트하여 노드 구성을 변경할 수 있습니다:

```bash
docker run -d \
  --name pixelzx-node \
  -p 8545:8545 \
  -p 8546:8546 \
  -p 30303:30303 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/configs/production.yaml:/app/configs/production.yaml \
  yuchanshin/pixelzx-evm:latest
```

## 포트

PIXELZX 노드에서 사용하는 포트:

- `8545` - JSON-RPC
- `8546` - WebSocket
- `30303` - P2P (TCP 및 UDP)

## 볼륨

권장하는 볼륨:

- `/app/data` - 노드 데이터 디렉토리
- `/app/configs` - 구성 파일
- `/app/keystore` - 키스토어 파일

## 환경 변수

- `CONFIG_ENV` - "development" 또는 "production"으로 설정하여 다른 구성 파일을 사용

## 문제 해결

### "exec format error"

이 오류는 호환되지 않는 아키텍처에서 이미지를 실행할 때 발생합니다. 플랫폼에 맞는 올바른 이미지를 사용하고 있는지 확인하세요.

### 연결 문제

필요한 포트가 방화벽에 의해 차단되지 않았는지 확인하세요.

### 권한 문제

데이터 디렉토리가 컨테이너에서 읽고 쓸 수 있는 올바른 권한을 가지고 있는지 확인하세요.