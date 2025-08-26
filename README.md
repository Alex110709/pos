# PIXELZX POS EVM Chain

PIXELZX를 네이티브 토큰으로 하는 Proof of Stake (POS) 기반 Ethereum Virtual Machine (EVM) 호환 블록체인 네트워크입니다.

## 주요 특징

- **네이티브 토큰**: PIXELZX (PXZ)
- **합의 메커니즘**: Proof of Stake (PoS)
- **EVM 호환성**: 완전한 Ethereum 스마트 컨트랙트 지원
- **높은 성능**: 3초 블록 타임, 1000+ TPS
- **낮은 수수료**: 가스비 최적화

## 토큰 사양

| 속성 | 값 |
|------|-----|
| 토큰명 | PIXELZX |
| 심볼 | PXZ |
| 총 공급량 | 1,000,000,000 PXZ |
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

### 빌드

```bash
go mod tidy
go build -o bin/pixelzx ./cmd/pixelzx
```

### 실행

```bash
# 제네시스 파일 초기화
./bin/pixelzx init

# 노드 시작
./bin/pixelzx start
```

## API 엔드포인트

### JSON-RPC API

- **포트**: 8545
- **URL**: http://localhost:8545

### WebSocket API

- **포트**: 8546
- **URL**: ws://localhost:8546

## 라이센스

MIT License