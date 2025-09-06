# PIXELZX POS EVM 체인 프로덕션 노드 설정

## 1. 개요

PIXELZX POS EVM 체인 프로덕션 노드는 이더리움 Geth 클라이언트와 유사한 CLI 인터페이스를 제공하면서도 PIXELZX 고유의 기능을 지원하는 블록체인 노드입니다. 이 문서는 프로덕션 환경에서 PIXELZX 노드를 설정하고 운영하기 위한 가이드를 제공합니다.

## 2. 아키텍처

### 2.1 시스템 구성도

```
graph TD
    A[PIXELZX CLI] --> B[Application Layer]
    B --> C[EVM Layer]
    B --> D[Consensus Layer]
    B --> E[Network Layer]
    B --> F[Storage Layer]
    
    C --> G[Smart Contract Execution]
    D --> H[PoS Consensus Engine]
    E --> I[P2P Networking]
    F --> J[Database Management]
    
    K[External Clients] --> L[JSON-RPC API]
    K --> M[WebSocket API]
    L --> B
    M --> B
    
    N[Admin Tools] --> O[Admin Commands]
    O --> B
```

### 2.2 주요 컴포넌트

| 컴포넌트 | 설명 | 기술 스택 |
|---------|------|----------|
| CLI 인터페이스 | 이더리움 Geth와 유사한 명령어 구조 | Cobra Framework |
| 합의 엔진 | Proof of Stake 알고리즘 | Go |
| 네트워크 레이어 | P2P 통신 및 블록 전파 | go-ethereum p2p |
| 저장소 레이어 | 블록체인 데이터 및 상태 저장 | LevelDB |
| API 레이어 | JSON-RPC 및 WebSocket 인터페이스 | Go HTTP Server |

## 3. CLI 명령어 구조

### 3.1 기본 명령어 구조

PIXELZX CLI는 이더리움 Geth와 유사한 명령어 구조를 따릅니다:

```
pixelzx [global options] command [command options] [arguments...]
```

### 3.2 주요 명령어

#### 3.2.1 노드 관리 명령어

| 명령어 | 설명 | 사용 예시 |
|--------|------|----------|
| `init` | 제네시스 블록 초기화 | `pixelzx init --datadir /data` |
| `start` | 노드 시작 | `pixelzx start --config /config.yaml` |
| `stop` | 노드 정지 | `pixelzx stop` |
| `reset` | 블록체인 데이터 리셋 | `pixelzx reset --datadir /data` |
| `version` | 버전 정보 표시 | `pixelzx version` |

#### 3.2.2 계정 관리 명령어

| 명령어 | 설명 | 사용 예시 |
|--------|------|----------|
| `account new` | 새 계정 생성 | `pixelzx account new --keystore /keystore` |
| `account list` | 계정 목록 표시 | `pixelzx account list --keystore /keystore` |
| `account import` | 개인키에서 계정 가져오기 | `pixelzx account import privatekey.txt` |
| `account update` | 계정 정보 업데이트 | `pixelzx account update <address>` |
| `account delete` | 계정 삭제 | `pixelzx account delete <address>` |

#### 3.2.3 블록체인 데이터 명령어

| 명령어 | 설명 | 사용 예시 |
|--------|------|----------|
| `export` | 블록체인 데이터 내보내기 | `pixelzx export chain.dat` |
| `import` | 블록체인 데이터 가져오기 | `pixelzx import chain.dat` |
| `dump` | 특정 블록 덤프 | `pixelzx dump 12345` |
| `prune` | 블록체인 데이터 정리 | `pixelzx prune --datadir /data` |

#### 3.2.4 네트워크 명령어

| 명령어 | 설명 | 사용 예시 |
|--------|------|----------|
| `attach` | 실행 중인 노드에 연결 | `pixelzx attach http://localhost:8545` |
| `console` | JavaScript 콘솔 시작 | `pixelzx console` |
| `peers` | 연결된 피어 목록 표시 | `pixelzx peers` |
| `peers add` | 새 피어 추가 | `pixelzx peers add enode://...` |
| `peers remove` | 피어 제거 | `pixelzx peers remove <nodeID>` |

#### 3.2.5 관리자 명령어

| 명령어 | 설명 | 사용 예시 |
|--------|------|----------|
| `admin status` | 노드 상태 확인 | `pixelzx admin status` |
| `admin backup` | 데이터 백업 | `pixelzx admin backup --target /backup` |
| `admin restore` | 데이터 복원 | `pixelzx admin restore --source /backup` |
| `admin config` | 설정 관리 | `pixelzx admin config show` |
| `admin debug` | 디버깅 도구 | `pixelzx admin debug trace` |
| `admin peer self` | 로컬 노드 enode 정보 조회 | `pixelzx admin peer self` |

##### 3.2.5.1 로컬 노드 정보 조회 (`admin peer self`)

로컬 PIXELZX 노드의 고유 식별자인 enode URL을 조회할 수 있습니다. 이 정보는 다른 노드와 P2P 연결을 설정할 때 사용됩니다.

```
# 기본 텍스트 형식으로 enode 정보 조회
pixelzx admin peer self

# JSON 형식으로 enode 정보 조회
pixelzx admin peer self --format json
```

출력 예시:
```
로컬 노드 enode 정보:
========================
enode URL: enode://a1b2c3d4e5f67890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef12345@192.168.1.100:30303
Node ID: a1b2c3d4e5f67890abcdef1234567890abcdef1234567890abcdef1234567890abcdef
IP 주소: 192.168.1.100
TCP 포트: 30303
UDP 포트: 30303
```

### 3.3 글로벌 옵션

| 옵션 | 설명 | 기본값 |
|------|------|--------|
| `--datadir` | 데이터 디렉토리 경로 | `$HOME/.pixelzx` |
| `--config` | 설정 파일 경로 | `config.yaml` |
| `--chain-id` | 체인 ID | `8888` |
| `--log-level` | 로그 레벨 | `info` |
| `--log-file` | 로그 파일 경로 | `/var/log/pixelzx.log` |

## 4. 설정 파일

### 4.1 프로덕션 설정 파일 (production.yaml)

```
# PIXELZX POS EVM Chain Production Configuration

# Network Configuration
network:
  name: "pixelzx-mainnet"
  chain_id: 8888
  block_time: 3s
  epoch_length: 200

# Node Configuration
node:
  name: "pixelzx-production-node"
  datadir: "/app/data"
  log_level: "info"

# P2P Network
p2p:
  port: 30303
  max_peers: 100
  bootnode_addrs: []
  enable_upnp: false
  nat: "none"

# JSON-RPC API
rpc:
  enabled: true
  host: "0.0.0.0"
  port: 8545
  cors_origins: ["*"]
  apis: ["eth", "net", "web3", "pxz"]
  gas_cap: 50000000
  tx_fee_cap: 1000000000000000000

# WebSocket API
websocket:
  enabled: true
  host: "0.0.0.0"
  port: 8546
  origins: ["*"]
  apis: ["eth", "net", "web3", "pxz"]

# Validator Configuration
validator:
  enabled: false
  address: ""
  password_file: ""

# Staking Configuration
staking:
  min_validator_stake: "1000000000000000000000000000"  # 1,000,000,000 PXZ (10억 PXZ)
  min_delegator_stake: "10000000000000000000000"       # 10,000 PXZ
  unbonding_period: "504h"                             # 21 days
  max_validators: 125

# Gas Configuration
gas:
  limit: 30000000
  price: 20000000000
  min_price: 1000000000
  price_bump: 10

# Security
security:
  keystore_dir: "/app/keystore"
  auto_unlock: false

# Metrics
metrics:
  enabled: true
  host: "0.0.0.0"
  port: 6060
  expensive: false

# Database
database:
  type: "leveldb"
  cache: 1024
  handles: 512
  ancient: "/app/data/ancient"
  freeze_threshold: 30000000

# Logging
logging:
  level: "info"
  file: "/app/logs/pixelzx.log"
  max_size: "100M"
  max_backups: 10
  max_age: 30
  compress: true

# Performance
performance:
  cache:
    database: 512
    trie: 256
    gc: 25
  workers: 4
  
# Sync Configuration
sync:
  mode: "full"
  gc_mode: "full"
  snapshot: true
  preimages: false

# Transaction Pool
txpool:
  locals: []
  no_locals: false
  journal: "/app/data/transactions.rlp"
  rejournal: "1h"
  price_limit: 1000000000
  price_bump: 10
  account_slots: 16
  global_slots: 4096
  account_queue: 64
  global_queue: 1024
  lifetime: "3h"

# Network Security
network_security:
  trusted_nodes: []
  banned_ips: []
  rate_limit: 1000
```

## 5. Docker 배포

### 5.1 Docker 이미지 빌드

```
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 의존성 설치
RUN apk add --no-cache git build-base

# 소스 코드 복사
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 빌드
RUN go build -o pixelzx ./cmd/pixelzx

# 실행 이미지
FROM alpine:latest

# 필요한 패키지 설치
RUN apk --no-cache add ca-certificates tzdata bash

WORKDIR /app

# 바이너리 복사
COPY --from=builder /app/pixelzx .

# 포트 노출
EXPOSE 8545 8546 30303

# 엔트리포인트 설정
ENTRYPOINT ["./pixelzx"]
```

### 5.2 Docker Compose 설정

```
# docker-compose.production.yml
version: '3.8'

services:
  pixelzx-node:
    image: yuchanshin/pixelzx-evm:latest
    container_name: pixelzx-production-node
    ports:
      - "8545:8545"   # JSON-RPC
      - "8546:8546"   # WebSocket
      - "30303:30303" # P2P
      - "6060:6060"   # Metrics
    volumes:
      - pixelzx-data:/app/data
      - pixelzx-keystore:/app/keystore
      - pixelzx-logs:/app/logs
      - ./configs/production.yaml:/app/config.yaml:ro
    environment:
      - PIXELZX_CHAIN_ID=8888
      - PIXELZX_LOG_LEVEL=info
    command: ["start", "--config", "/app/config.yaml"]
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8545"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 60s

volumes:
  pixelzx-data:
  pixelzx-keystore:
  pixelzx-logs:
```

### 5.3 멀티 아키텍처 지원

Docker Buildx를 사용하여 멀티 아키텍처 이미지를 빌드합니다:

```
# Buildx 설정
docker buildx create --name pixelzx-builder --use

# 멀티 플랫폼 빌드 및 푸시
docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 \
  -t yuchanshin/pixelzx-evm:latest \
  -t yuchanshin/pixelzx-evm:v1.0.0 \
  --push .
```

## 6. 보안 고려사항

### 6.1 키 관리

1. **키스토어 보안**:
   - 키스토어 디렉토리에 대한 적절한 파일 시스템 권한 설정
   - 민감한 키 정보는 환경 변수 또는 외부 시크릿 관리 시스템 사용

2. **계정 잠금 해제**:
   - 프로덕션 환경에서는 자동 잠금 해제 비활성화
   - 필요 시 수동 잠금 해제 사용

### 6.2 네트워크 보안

1. **API 엔드포인트 보안**:
   - 외부 접근이 필요한 경우 CORS 및 가상 호스트 설정
   - 민감한 API 엔드포인트는 방화벽 또는 reverse proxy로 보호

2. **피어 관리**:
   - 신뢰할 수 있는 피어만 연결 허용
   - IP 기반 접근 제어 설정

## 7. 모니터링 및 로깅

### 7.1 로깅 설정

프로덕션 환경에서는 다음과 같은 로깅 구성을 권장합니다:

```
logging:
  level: "info"
  file: "/app/logs/pixelzx.log"
  max_size: "100M"
  max_backups: 10
  max_age: 30
  compress: true
```

### 7.2 메트릭 수집

메트릭 엔드포인트를 활성화하여 Prometheus와 같은 모니터링 시스템과 연동:

```
metrics:
  enabled: true
  host: "0.0.0.0"
  port: 6060
  expensive: false
```

## 8. 백업 및 복구

### 8.1 데이터 백업

관리자 명령어를 사용하여 정기적인 데이터 백업 수행:

```
# 데이터 백업
pixelzx admin backup --target /backup/location

# 특정 블록체인 데이터만 백업
pixelzx export --datadir /app/data chain_backup.dat
```

### 8.2 데이터 복구

백업된 데이터를 사용하여 노드 복구:

```
# 전체 데이터 복구
pixelzx admin restore --source /backup/location

# 블록체인 데이터 복구
pixelzx import --datadir /app/data chain_backup.dat
```

## 9. 성능 튜닝

### 9.1 캐시 설정

데이터베이스 및 Trie 캐시를 적절히 조정하여 성능 최적화:

```
performance:
  cache:
    database: 512  # MB
    trie: 256      # MB
    gc: 25         # GC percentage
  workers: 4
```

### 9.2 트랜잭션 풀 설정

트랜잭션 풀 파라미터를 조정하여 네트워크 트래픽 및 처리량 최적화:

```
txpool:
  account_slots: 16
  global_slots: 4096
  account_queue: 64
  global_queue: 1024
```