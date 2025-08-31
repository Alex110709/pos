# PIXELZX POS EVM Chain - Docker Hub 배포 가이드

## 개요

PIXELZX POS EVM 체인이 Docker Hub에 성공적으로 배포되었습니다. 이제 전 세계 어디서나 간편하게 PIXELZX 노드를 실행할 수 있습니다.

**Docker Hub 저장소**: `yuchanshin/pixelzx-evm`

## 빠른 시작

### 1. 이미지 다운로드 및 실행

```bash
# 최신 이미지 다운로드
docker pull yuchanshin/pixelzx-evm:latest

# 간단한 노드 실행
docker run -d \
  --name pixelzx-node \
  -p 8545:8545 \
  -p 8546:8546 \
  -p 30303:30303 \
  -v pixelzx-data:/app/data \
  -v pixelzx-keystore:/app/keystore \
  -e PIXELZX_CHAIN_ID=8888 \
  yuchanshin/pixelzx-evm:latest

# 로그 확인
docker logs -f pixelzx-node
```

### 2. Docker Compose 사용 (권장)

프로덕션 환경에서는 Docker Compose를 사용하는 것을 권장합니다:

```bash
# docker-compose.production.yml 다운로드
wget https://raw.githubusercontent.com/pixelzx/pos/main/docker-compose.production.yml

# 서비스 시작
docker-compose -f docker-compose.production.yml up -d

# 상태 확인
docker-compose -f docker-compose.production.yml ps

# 로그 확인
docker-compose -f docker-compose.production.yml logs -f pixelzx-node
```

## 이미지 정보

### 태그 설명

- `latest`: 최신 릴리즈 버전 (프로덕션 환경 권장)
- `v{major}.{minor}.{patch}`: 특정 릴리즈 버전 (예: v1.0.0)
- `{commit-hash}`: 특정 커밋 버전 (개발/테스트용)

### 이미지 크기

- **최종 이미지 크기**: 약 15MB
- **아키텍처**: linux/arm64, linux/amd64 (멀티 아키텍처 지원)

## 네트워크 포트

| 포트 | 프로토콜 | 설명 |
|------|----------|------|
| 8545 | HTTP | JSON-RPC API |
| 8546 | WebSocket | WebSocket API |
| 30303 | TCP/UDP | P2P 네트워크 |

## 환경 변수

| 변수명 | 기본값 | 설명 |
|--------|--------|------|
| PIXELZX_CHAIN_ID | 8888 | 체인 ID (프로덕션: 8888, 개발: 7777) |
| PIXELZX_NETWORK | mainnet | 네트워크 타입 |
| PIXELZX_HOME | /app | 홈 디렉토리 |
| PIXELZX_CONFIG | /app/config.yaml | 설정 파일 경로 |
| PIXELZX_DATA_DIR | /app/data | 데이터 디렉토리 |

## 볼륨 매핑

| 컨테이너 경로 | 설명 | 권장 호스트 경로 |
|---------------|------|------------------|
| /app/data | 블록체인 데이터 | pixelzx-data (도커 볼륨) |
| /app/keystore | 키스토어 파일 | pixelzx-keystore (도커 볼륨) |
| /app/logs | 로그 파일 | pixelzx-logs (도커 볼륨) |

## 고급 사용법

### 1. 커스텀 설정으로 실행

```bash
# 커스텀 설정 파일 마운트
docker run -d \
  --name pixelzx-custom \
  -p 8545:8545 \
  -v $(pwd)/my-config.yaml:/app/config.yaml \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest
```

### 2. 개발 모드 실행

```bash
# 개발 체인 ID로 실행
docker run -d \
  --name pixelzx-dev \
  -p 8545:8545 \
  -e PIXELZX_CHAIN_ID=7777 \
  -e PIXELZX_NETWORK=devnet \
  yuchanshin/pixelzx-evm:latest
```

### 3. 노드 초기화

```bash
# 제네시스 초기화
docker run --rm \
  -v pixelzx-data:/app/data \
  yuchanshin/pixelzx-evm:latest \
  pixelzx init pixelzx-mainnet --chain-id 8888 --datadir /app/data
```

## 모니터링

### 1. 헬스체크

```bash
# 노드 상태 확인
docker exec pixelzx-node pixelzx version

# 헬스체크 상태 확인
docker inspect pixelzx-node | grep -A 5 Health
```

### 2. 메트릭스 수집

Docker Compose에 포함된 모니터링 스택 사용:

```bash
# 모니터링 포함하여 실행
docker-compose -f docker-compose.production.yml --profile monitoring up -d

# Grafana 대시보드 접속
open http://localhost:3000
# 기본 계정: admin / pixelzx-admin
```

## 업데이트 및 관리

### 1. 이미지 업데이트

```bash
# 최신 이미지 다운로드
docker pull yuchanshin/pixelzx-evm:latest

# 기존 컨테이너 중지 및 제거
docker stop pixelzx-node
docker rm pixelzx-node

# 새 이미지로 재시작
docker run -d \
  --name pixelzx-node \
  -p 8545:8545 \
  -p 8546:8546 \
  -p 30303:30303 \
  -v pixelzx-data:/app/data \
  -v pixelzx-keystore:/app/keystore \
  yuchanshin/pixelzx-evm:latest
```

### 2. 데이터 백업

```bash
# 데이터 볼륨 백업
docker run --rm \
  -v pixelzx-data:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/pixelzx-data-backup.tar.gz -C /data .

# 키스토어 백업
docker run --rm \
  -v pixelzx-keystore:/keystore \
  -v $(pwd):/backup \
  alpine tar czf /backup/pixelzx-keystore-backup.tar.gz -C /keystore .
```

## 트러블슈팅

### 1. 일반적인 문제

**포트 충돌**
```bash
# 사용 중인 포트 확인
netstat -tulpn | grep :8545

# 다른 포트로 바인딩
docker run -d -p 8555:8545 ... yuchanshin/pixelzx-evm:latest
```

**권한 문제**
```bash
# 볼륨 권한 확인
docker exec pixelzx-node ls -la /app/

# 권한 수정 (필요시)
docker exec -u root pixelzx-node chown -R pixelzx:pixelzx /app/data
```

### 2. 로그 확인

```bash
# 실시간 로그 모니터링
docker logs -f pixelzx-node

# 최근 100줄 로그
docker logs --tail 100 pixelzx-node

# 특정 시간 범위 로그
docker logs --since "2024-01-01T00:00:00" --until "2024-01-01T12:00:00" pixelzx-node
```

## 성능 최적화

### 1. 시스템 리소스

```bash
# CPU 및 메모리 제한 설정
docker run -d \
  --name pixelzx-node \
  --cpus="2.0" \
  --memory="4g" \
  -p 8545:8545 \
  yuchanshin/pixelzx-evm:latest
```

### 2. 스토리지 최적화

```bash
# SSD 볼륨 사용 (권장)
docker volume create --driver local \
  --opt type=none \
  --opt o=bind \
  --opt device=/path/to/ssd/storage \
  pixelzx-ssd-data
```

## 보안 권장사항

1. **방화벽 설정**: P2P 포트(30303)만 외부에 개방
2. **API 보안**: RPC 포트(8545, 8546)는 필요시에만 개방
3. **정기 업데이트**: 최신 보안 패치가 포함된 이미지로 업데이트
4. **키스토어 보안**: 키스토어 볼륨을 안전한 위치에 백업

## 지원 및 커뮤니티

- **GitHub**: https://github.com/pixelzx/pos
- **Discord**: https://discord.gg/pixelzx
- **문서**: https://docs.pixelzx.io
- **이슈 리포트**: https://github.com/pixelzx/pos/issues

---

**참고**: 이 문서는 PIXELZX POS EVM Chain v1.0.0을 기준으로 작성되었습니다. 최신 정보는 공식 문서를 확인하세요.