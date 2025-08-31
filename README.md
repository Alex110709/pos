# PIXELZX POS EVM Chain

PIXELZX를 네이티브 토큰으로 하는 Proof of Stake (POS) 기반 Ethereum Virtual Machine (EVM) 호환 블록체인 네트워크입니다.

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

### Exec Format Error

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