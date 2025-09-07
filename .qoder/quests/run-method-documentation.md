# PIXELZX POS EVM 체인

PIXELZX는 Proof of Stake(PoS) 합의 알고리즘을 기반으로 하는 EVM(Ethereum Virtual Machine) 호환 블록체인입니다. 이 프로젝트는 고성능(3초 블록 타임, 1000+ TPS)과 낮은 트랜잭션 수수료를 제공하며, PIXELZX(PXZ)를 네이티브 토큰으로 사용합니다.

## 목차

1. [실행 방법](#실행-방법)
2. [Peer 연결 방법](#peer-연결-방법)
3. [Docker에서 실행](#docker에서-실행)
4. [Docker Compose 사용](#docker-compose-사용)

## 실행 방법

### 요구 사항

- Go 1.21 이상
- Linux, macOS 또는 Windows (WSL 권장)

### 소스 코드 빌드

1. 저장소를 클론합니다:
   ```bash
   git clone <repository-url>
   cd pos
   ```

2. 의존성 설치:
   ```bash
   make deps
   # 또는
   go mod download
   go mod tidy
   ```

3. 바이너리 빌드:
   ```bash
   make build
   # 또는
   go build -o bin/pixelzx ./cmd/pixelzx
   ```

### 노드 초기화

제네시스 블록으로 네트워크를 초기화합니다:

```bash
# 기본 초기화
./bin/pixelzx init

# 특정 네트워크 이름으로 초기화
./bin/pixelzx init --network pixelzx-mainnet
```

### 노드 실행

초기화된 노드를 시작합니다:

```bash
# 기본 실행
./bin/pixelzx start

# 검증자 모드로 실행
./bin/pixelzx start --validator
```

## Peer 연결 방법

PIXELZX 노드는 P2P 네트워크를 통해 다른 노드와 연결됩니다. 노드는 시작 시 자동으로 부트스트랩 피어에 연결하고, 네트워크 내의 다른 피어들과 통신합니다.

### 네트워크 포트

PIXELZX는 다음 포트를 사용합니다:

- **JSON-RPC**: 8545
- **WebSocket**: 8546
- **P2P**: 30303

### 피어 연결 설정

피어 연결은 설정 파일을 통해 구성할 수 있습니다. `configs/development.yaml` 또는 `configs/production.yaml` 파일에서 네트워크 설정을 조정할 수 있습니다:

```yaml
network:
  name: "pixelzx-mainnet"
  chain_id: 8888
  bootnodes:
    - "enode://bootnode1@192.168.1.100:30303"
    - "enode://bootnode2@192.168.1.101:30303"
```

## Docker에서 실행

PIXELZX는 Docker를 통해 쉽게 배포할 수 있습니다.

### Docker 이미지 빌드

로컬에서 Docker 이미지를 빌드합니다:

```bash
# 기본 이미지 빌드
make docker-build

# Docker Hub용 이미지 빌드
make docker-build-hub
```

### Docker 컨테이너 실행

빌드된 이미지를 사용하여 컨테이너를 실행합니다:

```bash
docker run -d \
  --name pixelzx-node \
  -p 8545:8545 \
  -p 8546:8546 \
  -p 30303:30303 \
  -v $(pwd)/data:/app/data \
  pixelzx-node:latest
```

### Docker를 통한 노드 초기화 및 실행

```bash
# 초기화
docker run --rm \
  -v $(pwd)/data:/app/data \
  pixelzx-node:latest \
  pixelzx init

# 실행
docker run -d \
  --name pixelzx-node \
  -p 8545:8545 \
  -p 8546:8546 \
  -p 30303:30303 \
  -v $(pwd)/data:/app/data \
  pixelzx-node:latest \
  pixelzx start
```

## Docker Compose 사용

PIXELZX는 Docker Compose를 통해 멀티 노드 환경을 쉽게 구성할 수 있습니다.

### 개발 환경 실행

```bash
# 개발 환경 시작
make compose-dev-up

# 또는 직접 실행
docker-compose -f docker-compose.dev.yml up -d
```

### 프로덕션 환경 실행

```bash
# 프로덕션 환경 시작
make compose-up

# 또는 직접 실행
docker-compose -f docker-compose.yml up -d
```

### 로그 확인

```bash
# 개발 환경 로그
make compose-dev-logs

# 프로덕션 환경 로그
make compose-logs
```

### Docker Compose 설정 예시

`docker-compose.yml` 파일 예시:

```yaml
version: '3.8'

services:
  pixelzx-node-1:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: pixelzx-node-1
    ports:
      - "8545:8545"
      - "8546:8546"
      - "30303:30303"
    volumes:
      - ./data/node1:/app/data
    command: ["pixelzx", "start", "--validator"]

  pixelzx-node-2:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: pixelzx-node-2
    ports:
      - "8547:8545"
      - "8548:8546"
      - "30304:30303"
    volumes:
      - ./data/node2:/app/data
    command: ["pixelzx", "start"]

volumes:
  node1-data:
  node2-data:
```

### 설정 파일

Docker Compose 환경에서는 각 노드의 설정을 `configs/` 디렉토리에 있는 파일을 통해 구성할 수 있습니다:

- `configs/development.yaml`: 개발 환경 설정
- `configs/production.yaml`: 프로덕션 환경 설정

## 추가 명령어

PIXELZX CLI는 다양한 관리 기능을 제공합니다:

```bash
# 계정 관리
./bin/pixelzx account new
./bin/pixelzx account list

# 검증자 관리
./bin/pixelzx validator list
./bin/pixelzx validator register

# 스테이킹
./bin/pixelzx staking delegate
./bin/pixelzx staking undelegate

# 설정 관리
./bin/pixelzx config show
./bin/pixelzx config validate
```

## 문제 해결

### 일반적인 문제

1. **포트 충돌**: 다른 프로세스가 포트를 사용 중일 수 있습니다. 포트를 변경하거나 다른 프로세스를 종료하세요.

2. **권한 문제**: 데이터 디렉토리에 대한 읽기/쓰기 권한을 확인하세요.

3. **네트워크 연결 문제**: 방화벽 설정과 부트스트랩 노드 설정을 확인하세요.

### 로그 확인

실행 중인 노드의 로그를 확인하여 문제를 진단할 수 있습니다:

```bash
# 로컬 실행 시
tail -f data/logs/pixelzx.log

# Docker 실행 시
docker logs -f pixelzx-node

# Docker Compose 실행 시
docker-compose logs -f pixelzx-node-1
```

## 라이선스

이 프로젝트는 MIT 라이선스 하에 있습니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.
