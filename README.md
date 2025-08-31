# PIXELZX POS EVM Chain

PIXELZX를 네이티브 토큰으로 하는 Proof of Stake (POS) 기반 Ethereum Virtual Machine (EVM) 호환 블록체인 네트워크입니다.

## 목차

- [주요 특징](#주요-특징)
- [토큰 사양](#토큰-사양)
- [네트워크 파라미터](#네트워크-파라미터)
- [아키텍처](#아키텍처)
- [디렉토리 구조](#디렉토리-구조)
- [빌드 및 실행](#빌드-및-실행)
  - [의존성](#의존성)
  - [로컬 빌드](#로컬-빌드)
  - [Docker 빌드](#docker-빌드)
  - [실행](#실행)
- [**Docker 빠른 시작**](#docker-빠른-시작) 🚀
  - [기본 Docker 명령어](#기본-docker-명령어)
  - [환경 변수 설정](#환경-변수-설정)
  - [볼륨 마운트 가이드](#볼륨-마운트-가이드)
  - [헬스체크 및 상태 확인](#헬스체크-및-상태-확인)
- [**P2P 네트워크 연결**](#p2p-네트워크-연결) 🌐
  - [P2P 포트 설정](#p2p-포트-설정)
  - [부트노드 연결](#부트노드-연결)
  - [네트워크 상태 모니터링](#네트워크-상태-모니터링)
  - [P2P 연결 트러블슈팅](#p2p-연결-트러블슈팅)
- [**노드 초기화 및 설정**](#노드-초기화-및-설정) ⚙️
  - [제네시스 파일 초기화](#제네시스-파일-초기화)
  - [데이터 디렉토리 설정](#데이터-디렉토리-설정)
  - [설정 파일 커스터마이징](#설정-파일-커스터마이징)
  - [키스토어 관리](#키스토어-관리)
  - [초기화 검증](#초기화-검증)
- [API 엔드포인트](#api-엔드포인트)
- [**문제 해결**](#문제-해결) 🚑
  - [Docker 관련 문제](#docker-관련-문제)
  - [P2P 네트워크 문제](#p2p-네트워크-문제)
  - [API 연결 문제](#api-연결-문제)
  - [성능 문제](#성능-문제)
  - [로그 분석](#로그-분석)
- [Docker Hub](#docker-hub)
- [라이센스](#라이센스)

## 주요 특징

- **네이티브 토큰**: PIXELZX (PXZ)
- **합의 메커니즘**: Proof of Stake (PoS)
- **EVM 호환성**: 완전한 Ethereum 스마트 컨트랙트 지원
- **높은 성능**: 3초 블록 타임, 1000+ TPS
- **낮은 수수료**: 가스비 최적화
- **멀티 아키텍처**: linux/amd64, linux/arm64, linux/arm/v7 지원

## 토큰 사양

| 속성 | 값 |
|------|-----|
| 토큰명 | PIXELZX |
| 심볼 | PXZ |
| 총 공급량 | 10,000,000,000,000,000 PXZ |
| 소수점 자리수 | 18 |
| 토큰 타입 | 네이티브 토큰 |

## 네트워크 파라미터

| 파라미터 | 값 |
|----------|-----|
| 블록 타임 | 3초 |
| 블록 크기 제한 | 30MB |
| 가스 제한 | 30,000,000 |
| 최대 검증자 수 | 125 |
| 언본딩 기간 | 21일 |

## 아키텍처

### 계층 구조

1. **Application Layer**: DApp 인터페이스, API 엔드포인트
2. **EVM Layer**: Ethereum 가상 머신, 스마트 컨트랙트 실행
3. **Consensus Layer**: PoS 합의 알고리즘, 블록 생성/검증
4. **Network Layer**: P2P 통신, 블록 전파
5. **Storage Layer**: 상태 저장소, 블록체인 데이터베이스

## 디렉토리 구조

```
pos/
├── cmd/                    # 실행 가능한 바이너리
├── consensus/              # PoS 합의 메커니즘
├── core/                   # 코어 블록체인 로직
├── evm/                    # EVM 통합 및 실행 환경
├── network/                # P2P 네트워킹
├── api/                    # JSON-RPC, WebSocket API
├── staking/                # 스테이킹 및 검증자 관리
├── governance/             # 거버넌스 시스템
├── storage/                # 데이터 저장 및 상태 관리
├── crypto/                 # 암호화 및 보안 기능
├── tests/                  # 테스트 코드
├── docs/                   # 문서
└── scripts/                # 유틸리티 스크립트
```

## 빌드 및 실행

### 의존성

- Go 1.21+
- Git
- Docker (선택사항)
- Docker Buildx (멀티 플랫폼 빌드용)

### 로컬 빌드

```bash
go mod tidy
go build -o bin/pixelzx ./cmd/pixelzx
```

### Docker 빌드

#### 단일 플랫폼 빌드
```bash
# 현재 플랫폼용 이미지 빌드
make docker-build-local

# 또는 직접 빌드
docker build -t pixelzx-pos:latest .
```

#### 멀티 플랫폼 빌드
```bash
# Docker Buildx 설정
make buildx-setup

# 모든 플랫폼용 빌드 및 배포
make docker-push-multi

# 플랫폼별 테스트
make docker-test-multi
```

#### 지원 플랫폼
- **linux/amd64**: Intel/AMD 64비트 프로세서
- **linux/arm64**: ARM 64비트 프로세서 (Apple Silicon, ARM 서버)
- **linux/arm/v7**: ARM 32비트 프로세서 (라즈베리파이 등)

### 실행

#### 로컬 실행
```bash
# 제네시스 파일 초기화
./bin/pixelzx init

# 노드 시작
./bin/pixelzx start
```

#### Docker 실행
```bash
# 프로덕션 환경
docker-compose -f docker-compose.production.yml up -d

# 개발 환경
docker-compose -f docker-compose.dev.yml up -d

# 직접 실행
docker run -d --name pixelzx-node \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  yuchanshin/pixelzx-evm:latest
```

## Docker 빠른 시작

PixelZX 노드를 Docker로 쉽고 빠르게 시작할 수 있는 방법을 안내합니다.

### 기본 Docker 명령어

#### 1. 이미지 다운로드
```bash
# 최신 이미지 다운로드
docker pull yuchanshin/pixelzx-evm:latest

# 특정 버전 다운로드
docker pull yuchanshin/pixelzx-evm:v1.0.0
```

#### 2. 노드 초기화 (선택사항)
```bash
# 메인넷 초기화
docker run --rm \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx init --datadir /app/data

# 테스트넷 초기화
docker run --rm \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx init --datadir /app/data --chain-id 8889
```

#### 3. 노드 실행
```bash
# 기본 실행
docker run -d \
  --name pixelzx-node \
  -p 8545:8545 \
  -p 8546:8546 \
  -p 30303:30303 \
  -v pixelzx-data:/app/data \
  -v pixelzx-keystore:/app/keystore \
  yuchanshin/pixelzx-evm:latest

# 환경 변수와 함께 실행
docker run -d \
  --name pixelzx-node \
  -p 8545:8545 \
  -p 8546:8546 \
  -p 30303:30303 \
  -e PIXELZX_CHAIN_ID=8888 \
  -e PIXELZX_NETWORK=mainnet \
  -v pixelzx-data:/app/data \
  -v pixelzx-keystore:/app/keystore \
  yuchanshin/pixelzx-evm:latest
```

### 환경 변수 설정

| 변수명 | 기본값 | 설명 | 예시 |
|--------|--------|------|------|
| PIXELZX_CHAIN_ID | 8888 | 체인 ID | 8888 (메인넷), 8889 (테스트넷) |
| PIXELZX_NETWORK | mainnet | 네트워크 타입 | mainnet, testnet, devnet |
| PIXELZX_P2P_PORT | 30303 | P2P 통신 포트 | 30303 |
| PIXELZX_RPC_PORT | 8545 | JSON-RPC API 포트 | 8545 |
| PIXELZX_WS_PORT | 8546 | WebSocket API 포트 | 8546 |
| PIXELZX_DATA_DIR | /app/data | 데이터 디렉토리 | /app/data |
| PIXELZX_KEYSTORE_DIR | /app/keystore | 키스토어 디렉토리 | /app/keystore |

### 볼륨 마운트 가이드

#### Docker 볼륨 생성
```bash
# 데이터 및 키스토어 볼륨 생성
docker volume create pixelzx-data
docker volume create pixelzx-keystore

# 볼륨 위치 확인
docker volume inspect pixelzx-data
docker volume inspect pixelzx-keystore
```

#### 호스트 디렉토리 마운트
```bash
# 호스트 디렉토리 생성
mkdir -p $HOME/pixelzx/{data,keystore}

# 호스트 디렉토리로 마운트
docker run -d \
  --name pixelzx-node \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v $HOME/pixelzx/data:/app/data \
  -v $HOME/pixelzx/keystore:/app/keystore \
  yuchanshin/pixelzx-evm:latest
```

### 헬스체크 및 상태 확인

```bash
# 컨테이너 상태 확인
docker ps | grep pixelzx-node

# 로그 확인
docker logs pixelzx-node

# 실시간 로그 확인
docker logs -f pixelzx-node

# 컨테이너 내부 접속
docker exec -it pixelzx-node /bin/sh

# 노드 버전 확인
docker exec pixelzx-node pixelzx version

# 블록 높이 확인
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://localhost:8545
```

## P2P 네트워크 연결

PixelZX 노드가 네트워크의 다른 노드들과 P2P 연결을 설정하는 방법을 안내합니다.

### P2P 포트 설정

#### 방화벽 설정
```bash
# Ubuntu/Debian 방화벽 설정
sudo ufw allow 30303/tcp
sudo ufw allow 30303/udp

# CentOS/RHEL 방화벽 설정
sudo firewall-cmd --permanent --add-port=30303/tcp
sudo firewall-cmd --permanent --add-port=30303/udp
sudo firewall-cmd --reload
```

#### Docker 포트 확인
```bash
# P2P 포트 확인
docker exec pixelzx-node netstat -tulpn | grep 30303

# 포트 바인딩 확인
docker port pixelzx-node
```

### 부트노드 연결

#### 네트워크 정보 확인
```bash
# 현재 노드 정보 확인
docker exec pixelzx-node pixelzx admin nodeInfo

# 연결된 피어 목록 확인
docker exec pixelzx-node pixelzx admin peers

# 피어 수 확인
docker exec pixelzx-node pixelzx admin peerCount
```

#### 수동 피어 추가
```bash
# 특정 피어에 연결
docker exec pixelzx-node pixelzx admin addPeer "enode://[PEER_ID]@[IP]:[PORT]"

# 예시: 부트노드 연결
docker exec pixelzx-node pixelzx admin addPeer "enode://abcd1234@52.123.45.67:30303"
```

### 네트워크 상태 모니터링

#### 동기화 상태 확인
```bash
# 블록 동기화 상태 확인
docker exec pixelzx-node pixelzx eth syncing

# 현재 블록 번호 확인
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://localhost:8545

# 네트워크 ID 확인
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"net_version","params":[],"id":1}' \
  http://localhost:8545
```

#### 피어 연결 상태 확인
```bash
# 연결된 피어 수 확인
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}' \
  http://localhost:8545

# 리스닝 상태 확인
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"net_listening","params":[],"id":1}' \
  http://localhost:8545
```

### P2P 연결 트러블슈팅

| 문제 | 증상 | 원인 | 해결방법 |
|------|------|------|----------|
| 피어 연결 실패 | peerCount가 0 | 방화벽 차단 | 포트 30303 개방 |
| 느린 동기화 | 블록 높이 증가 안함 | 부트노드 응답 없음 | 다른 부트노드 시도 |
| NAT 문제 | 인바운드 연결 불가 | 공인 IP 없음 | --nat 옵션 사용 |
| 포트 충돌 | 노드 시작 실패 | 포트 이미 사용 중 | 다른 포트 사용 |

#### 상세 디버깅
```bash
# 네트워크 연결 상태 확인
docker exec pixelzx-node ss -tulpn | grep 30303

# 외부에서 포트 접근 테스트
telnet [YOUR_PUBLIC_IP] 30303

# Docker 네트워크 설정 확인
docker inspect pixelzx-node | grep -A 10 "NetworkSettings"

# 방화벽 상태 확인 (Ubuntu)
sudo ufw status verbose

# NAT 설정으로 노드 재시작
docker run -d \
  --name pixelzx-node-nat \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --nat extip:[YOUR_PUBLIC_IP]
```

## 노드 초기화 및 설정

노드를 처음 시작할 때 필요한 초기화 과정과 설정 방법을 안내합니다.

### 제네시스 파일 초기화

#### 기본 초기화
```bash
# 메인넷 제네시스 초기화
docker run --rm \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx init --datadir /app/data

# 테스트넷 제네시스 초기화
docker run --rm \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx init --datadir /app/data --chain-id 8889 --network testnet
```

#### 커스텀 제네시스 파일 사용
```bash
# 커스텀 제네시스 파일 준비
cat > custom-genesis.json << EOF
{
  "chainId": 8888,
  "homesteadBlock": 0,
  "eip150Block": 0,
  "eip155Block": 0,
  "eip158Block": 0,
  "byzantiumBlock": 0,
  "constantinopleBlock": 0,
  "petersburgBlock": 0,
  "istanbulBlock": 0,
  "berlinBlock": 0,
  "londonBlock": 0,
  "alloc": {
    "0x742d35cc6672c0532925a3b8d6f7b71b47c0062f": {
      "balance": "1000000000000000000000000"
    }
  },
  "difficulty": "0x1",
  "gasLimit": "0x1c9c380"
}
EOF

# 커스텀 제네시스로 초기화
docker run --rm \
  -v $(pwd)/custom-genesis.json:/app/genesis.json \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx init /app/genesis.json --datadir /app/data
```

### 데이터 디렉토리 설정

#### 볼륨 관리
```bash
# 데이터 볼륨 생성
docker volume create pixelzx-data
docker volume create pixelzx-keystore

# 볼륨 백업
docker run --rm \
  -v pixelzx-data:/source \
  -v $(pwd):/backup \
  alpine tar czf /backup/pixelzx-data-backup.tar.gz -C /source .

# 볼륨 복원
docker run --rm \
  -v pixelzx-data:/target \
  -v $(pwd):/backup \
  alpine tar xzf /backup/pixelzx-data-backup.tar.gz -C /target

# 볼륨 내용 확인
docker run --rm \
  -v pixelzx-data:/data \
  alpine ls -la /data
```

#### 디렉토리 구조
```
pixelzx-data/
├── chaindata/          # 블록체인 데이터
├── nodes/              # 노드 정보
├── trie/               # 상태 트라이
└── ancient/            # 아카이브 데이터

pixelzx-keystore/
├── UTC--[timestamp]--[address]  # 키 파일들
└── ...
```

### 설정 파일 커스터마이징

#### 기본 설정 파일 추출
```bash
# 설정 파일 확인
docker run --rm yuchanshin/pixelzx-evm:latest ls -la /app/configs/

# 프로덕션 설정 파일 추출
docker run --rm \
  -v $(pwd):/backup \
  yuchanshin/pixelzx-evm:latest \
  cp /app/configs/production.yaml /backup/

# 개발 설정 파일 추출
docker run --rm \
  -v $(pwd):/backup \
  yuchanshin/pixelzx-evm:latest \
  cp /app/configs/development.yaml /backup/
```

#### 커스텀 설정으로 실행
```bash
# 설정 파일 수정 (예시)
cat > custom-config.yaml << EOF
chain_id: 8888
network_id: 8888
data_dir: "/app/data"
keystore_dir: "/app/keystore"

rpc:
  enabled: true
  host: "0.0.0.0"
  port: 8545
  cors: ["*"]
  api: ["eth", "net", "web3", "personal", "admin"]

ws:
  enabled: true
  host: "0.0.0.0"
  port: 8546
  origins: ["*"]
  api: ["eth", "net", "web3"]

p2p:
  enabled: true
  host: "0.0.0.0"
  port: 30303
  max_peers: 50
  
logging:
  level: "info"
  format: "json"
EOF

# 커스텀 설정으로 노드 실행
docker run -d \
  --name pixelzx-custom \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v $(pwd)/custom-config.yaml:/app/configs/production.yaml \
  -v pixelzx-data:/app/data \
  -v pixelzx-keystore:/app/keystore \
  yuchanshin/pixelzx-evm:latest
```

### 키스토어 관리

#### 계정 생성
```bash
# 새 계정 생성
docker exec -it pixelzx-node pixelzx account new

# 계정 목록 확인
docker exec pixelzx-node pixelzx account list

# 계정 정보 확인
docker exec pixelzx-node pixelzx account info [ADDRESS]
```

#### 키스토어 파일 관리
```bash
# 키스토어 파일 확인
docker exec pixelzx-node ls -la /app/keystore/

# 키스토어 파일 백업
docker cp pixelzx-node:/app/keystore/ ./keystore-backup/

# 키스토어 파일 복원
docker cp ./keystore-backup/ pixelzx-node:/app/keystore/
```

### 초기화 검증

#### 시스템 상태 확인
```bash
# 노드 버전 확인
docker exec pixelzx-node pixelzx version

# 체인 ID 확인
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
  http://localhost:8545

# 제네시스 블록 확인
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x0",true],"id":1}' \
  http://localhost:8545

# 계정 잔액 확인
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x742d35cc6672c0532925a3b8d6f7b71b47c0062f","latest"],"id":1}' \
  http://localhost:8545
```

#### 초기화 문제 해결
```bash
# 데이터 디렉토리 초기화 (주의: 모든 데이터 삭제)
docker volume rm pixelzx-data
docker volume create pixelzx-data

# 권한 문제 해결
docker exec pixelzx-node chown -R 1000:1000 /app/data
docker exec pixelzx-node chown -R 1000:1000 /app/keystore

# 로그에서 오류 확인
docker logs pixelzx-node | grep -i error
docker logs pixelzx-node | grep -i fatal
```

## API 엔드포인트

### JSON-RPC API

- **포트**: 8545
- **URL**: http://localhost:8545

### WebSocket API

- **포트**: 8546
- **URL**: ws://localhost:8546

### P2P 네트워크

- **포트**: 30303
- **프로토콜**: TCP/UDP

## 문제 해결

### Docker 관련 문제

#### Exec Format Error

Docker 컨테이너 실행 시 `exec format error`가 발생하는 경우:

1. **멀티 플랫폼 이미지 사용**: 
   ```bash
   docker run --rm yuchanshin/pixelzx-evm:latest /usr/local/bin/pixelzx version
   ```

2. **플랫폼 명시적 지정**:
   ```bash
   docker run --rm --platform linux/amd64 yuchanshin/pixelzx-evm:latest /usr/local/bin/pixelzx version
   ```

3. **로컬 빌드 사용**:
   ```bash
   make docker-build-local
   docker run --rm yuchanshin/pixelzx-evm:local /usr/local/bin/pixelzx version
   ```

자세한 내용은 [EXEC_FORMAT_ERROR_SOLUTION.md](./EXEC_FORMAT_ERROR_SOLUTION.md) 문서를 참조하세요.

#### 컨테이너 시작 실패
```bash
# 컨테이너 로그 확인
docker logs pixelzx-node

# 컨테이너 상태 확인
docker ps -a | grep pixelzx

# 포트 충돌 확인
sudo netstat -tulpn | grep -E '(8545|8546|30303)'

# 컨테이너 재시작
docker restart pixelzx-node

# 컨테이너 완전 재생성
docker stop pixelzx-node
docker rm pixelzx-node
docker run -d --name pixelzx-node \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest
```

#### 볼륨 권한 문제
```bash
# 볼륨 권한 확인
docker exec pixelzx-node ls -la /app/

# 권한 수정
docker exec pixelzx-node chown -R 1000:1000 /app/data
docker exec pixelzx-node chown -R 1000:1000 /app/keystore

# SELinux 환경에서 볼륨 마운트 문제
docker run -d --name pixelzx-node \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data:Z \
  -v pixelzx-keystore:/app/keystore:Z \
  yuchanshin/pixelzx-evm:latest
```

### P2P 네트워크 문제

#### 피어 연결 불가
```bash
# 방화벽 상태 확인
sudo ufw status
sudo firewall-cmd --list-ports

# NAT 환경에서 포트 포워딩 확인
# 라우터 설정에서 30303 포트를 노드 IP로 포워딩

# 네트워크 연결 테스트
telnet [REMOTE_NODE_IP] 30303

# P2P 디버깅 모드로 노드 시작
docker run -d --name pixelzx-debug \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --verbosity 5
```

#### 동기화 문제
```bash
# 블록 동기화 상태 상세 확인
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}' \
  http://localhost:8545

# 현재 블록과 네트워크 최신 블록 비교
# 1. 현재 노드 블록 높이
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://localhost:8545

# 2. 다른 노드에서 최신 블록 확인
# 공식 블록 익스플로러나 다른 노드 API 사용

# 동기화 재시작
docker restart pixelzx-node

# 빠른 동기화 모드 (스냅샷 사용)
docker run -d --name pixelzx-fast \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --syncmode fast
```

### API 연결 문제

#### JSON-RPC API 연결 실패
```bash
# API 서비스 상태 확인
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"web3_clientVersion","params":[],"id":1}' \
  http://localhost:8545

# 포트 리스닝 상태 확인
docker exec pixelzx-node netstat -tulpn | grep 8545

# 방화벽에서 API 포트 허용
sudo ufw allow 8545/tcp
sudo ufw allow 8546/tcp

# CORS 문제 해결
docker run -d --name pixelzx-cors \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --http.corsdomain "*" --ws.origins "*"
```

#### WebSocket 연결 문제
```bash
# WebSocket 연결 테스트
wscat ws://localhost:8546

# 또는 JavaScript로 테스트
node -e "
  const WebSocket = require('ws');
  const ws = new WebSocket('ws://localhost:8546');
  ws.on('open', () => {
    console.log('WebSocket 연결 성공');
    ws.close();
  });
  ws.on('error', (err) => {
    console.log('WebSocket 연결 실패:', err.message);
  });
"

# WebSocket 서비스 상태 확인
docker exec pixelzx-node netstat -tulpn | grep 8546
```

### 성능 문제

#### 메모리 부족
```bash
# 컨테이너 리소스 사용량 확인
docker stats pixelzx-node

# 메모리 제한 설정
docker run -d --name pixelzx-limited \
  --memory="2g" --memory-swap="4g" \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest

# 가비지 컨렉션 설정 조정
docker run -d --name pixelzx-gc \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --cache 1024 --gcmode archive
```

#### 느린 응답 시간
```bash
# 캐시 크기 증가
docker run -d --name pixelzx-cache \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --cache 2048

# SSD 사용 권장 (호스트 디렉토리 마운트 시)
mkdir -p /fast-ssd/pixelzx/{data,keystore}
docker run -d --name pixelzx-ssd \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v /fast-ssd/pixelzx/data:/app/data \
  -v /fast-ssd/pixelzx/keystore:/app/keystore \
  yuchanshin/pixelzx-evm:latest
```

### 로그 분석

#### 주요 로그 패턴
```bash
# 오류 로그 확인
docker logs pixelzx-node 2>&1 | grep -i error
docker logs pixelzx-node 2>&1 | grep -i fatal
docker logs pixelzx-node 2>&1 | grep -i panic

# P2P 연결 로그
docker logs pixelzx-node 2>&1 | grep -i peer
docker logs pixelzx-node 2>&1 | grep -i "connection"

# 동기화 로그
docker logs pixelzx-node 2>&1 | grep -i sync
docker logs pixelzx-node 2>&1 | grep -i "block"

# API 요청 로그
docker logs pixelzx-node 2>&1 | grep -i "rpc"
docker logs pixelzx-node 2>&1 | grep -i "http"
```

#### 로그 레벨 조정
```bash
# 디버그 로그 모드
docker run -d --name pixelzx-debug \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --verbosity 5

# 조용한 로그 모드
docker run -d --name pixelzx-quiet \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --verbosity 1

# JSON 형식 로그
docker run -d --name pixelzx-json \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx start --log.json
```

## Docker Hub

공식 이미지: [yuchanshin/pixelzx-evm](https://hub.docker.com/r/yuchanshin/pixelzx-evm)

```bash
# 최신 버전 다운로드
docker pull yuchanshin/pixelzx-evm:latest

# 특정 버전 다운로드
docker pull yuchanshin/pixelzx-evm:v1.0.0

# 매니페스트 확인 (지원 플랫폼 목록)
docker buildx imagetools inspect yuchanshin/pixelzx-evm:latest
```

## 라이센스

MIT License