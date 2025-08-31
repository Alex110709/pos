# PIXELZX POS EVM 체인 - 컨테이너 Exec 오류 문제 해결 설계

## 1. 문제 개요

PIXELZX POS EVM 체인 Docker 컨테이너 실행 중 다음과 같은 오류가 발생하는 문제를 해결합니다:

```
CI runtime exec failed: exec failed: unable to start container process: exec: "bash": executable file not found in $PATH: unknown
```

이 오류는 컨테이너 내부에 bash 실행 파일이 없거나 PATH에 포함되어 있지 않을 때 발생합니다. 주로 Alpine Linux 기반 이미지에서 발생하는 문제입니다.

## 2. 원인 분석

### 2.1 현재 시스템 아키텍처 확인

PIXELZX POS EVM 체인은 멀티 아키텍처를 지원하지만, 호스트 시스템과 컨테이너의 아키텍처가 일치하지 않으면 실행 오류가 발생할 수 있습니다.

| 빌드 환경 | 실행 환경 | 결과 | 원인 |
|----------|----------|------|------|
| AMD64 | AMD64 | ✅ 정상 | 아키텍처 일치 |
| AMD64 | ARM64 | ❌ exec format error | 아키텍처 불일치 |
| AMD64 | ARM/v7 | ❌ exec format error | 아키텍처 불일치 |

### 2.2 Dockerfile 분석

PIXELZX POS EVM 체인의 Dockerfile은 다음과 같은 구조로 되어 있습니다:

1. **빌드 스테이지**: golang:1.21-alpine 이미지 사용
2. **실행 스테이지**: alpine:latest 이미지 사용
3. **실행 환경**: Alpine Linux 기반으로 bash가 기본 설치되지 않음

### 2.3 오류 발생 지점

오류는 일반적으로 다음 상황에서 발생합니다:

1. 컨테이너 내부에서 bash를 필요로 하는 명령어 실행 시도
2. Docker Compose 파일에서 bash를 사용하는 명령어 정의
3. CI/CD 파이프라인에서 bash 기반 스크립트 실행

## 3. 해결 방안

### 3.1 즉시 해결 방안

#### 3.1.1 bash 대신 sh 사용
Alpine Linux는 기본적으로 bash 대신 sh를 사용합니다. 모든 스크립트를 sh 호환 구문으로 변경합니다.

#### 3.1.2 bash 설치
Dockerfile에 bash 설치 명령어 추가:

```dockerfile
# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata bash
```

### 3.2 플랫폼 명시적 지정

#### 3.2.1 Docker 실행 시 플랫폼 지정

AMD64 환경 (Intel/AMD PC):
```bash
docker run --platform linux/amd64 --rm yuchanshin/pixelzx-evm:latest
```

ARM64 환경 (Apple Silicon Mac, Raspberry Pi 4 이상):
```bash
docker run --platform linux/arm64 --rm yuchanshin/pixelzx-evm:latest
```

ARM/v7 환경 (Raspberry Pi 3 이하):
```bash
docker run --platform linux/arm/v7 --rm yuchanshin/pixelzx-evm:latest
```

#### 3.2.2 Docker Compose 사용 시 플랫폼 지정

```yaml
version: '3.8'
services:
  pixelzx-node:
    image: yuchanshin/pixelzx-evm:latest
    platform: ${DOCKER_DEFAULT_PLATFORM:-linux/amd64}
    # ... 나머지 설정
```

### 3.3 플랫폼 자동 감지 스크립트 사용

프로젝트에 포함된 플랫폼 자동 감지 스크립트를 사용:

```bash
# 기본 실행 (도움말 표시)
./scripts/detect-platform.sh

# 특정 명령어 실행
./scripts/detect-platform.sh admin status
./scripts/detect-platform.sh init
./scripts/detect-platform.sh start
```

## 4. 구현 계획

### 4.1 Dockerfile 수정

1. **실행 스테이지에 bash 설치 추가**:
   ```dockerfile
   # Install runtime dependencies
   RUN apk --no-cache add ca-certificates tzdata bash
   ```

2. **쉘 스크립트 호환성 검토 및 수정**:
   - 모든 sh 스크립트가 bash 의존성 없이 동작하도록 수정

### 4.2 Docker Compose 파일 수정

1. **docker-compose.yml에 플랫폼 정보 추가**:
   ```yaml
   services:
     pixelzx-node:
       # ... 기존 설정
       platform: ${DOCKER_DEFAULT_PLATFORM:-linux/amd64}
   ```

### 4.3 문서 업데이트

1. **DOCKER_HUB_GUIDE.md 업데이트**:
   - 플랫폼 관련 실행 방법 추가
   - bash 관련 문제 해결 가이드 추가

2. **README.md 업데이트**:
   - 플랫폼 호환성 관련 정보 추가
   - exec 오류 해결 방법 추가

## 5. 테스트 계획

### 5.1 단위 테스트

1. **bash 설치 여부 확인**:
   ```bash
   docker run --rm yuchanshin/pixelzx-evm:latest which bash
   ```

2. **sh 스크립트 호환성 테스트**:
   ```bash
   docker run --rm yuchanshin/pixelzx-evm:latest sh -c "echo 'Hello World'"
   ```

### 5.2 통합 테스트

1. **Docker Compose 환경 테스트**:
   ```bash
   docker-compose up -d
   docker-compose exec pixelzx-node pixelzx version
   ```

2. **플랫폼 자동 감지 스크립트 테스트**:
   ```bash
   ./scripts/detect-platform.sh version
   ```

### 5.3 CI/CD 파이프라인 테스트

1. **다양한 아키텍처에서 빌드 및 실행 테스트**
2. **exec 명령어 테스트**
3. **admin 명령어 테스트**

## 6. 보안 고려사항

1. **최소 권한 원칙**:
   - Alpine Linux 기반으로 보안 취약점 최소화
   - 불필요한 패키지 설치 금지

2. **사용자 권한 분리**:
   - root 사용자 대신 pixelzx 사용자로 실행
   - 필요한 디렉토리에만 쓰기 권한 부여

## 7. 성능 고려사항

1. **이미지 크기 최적화**:
   - 멀티 스테이지 빌드 유지
   - 불필요한 파일 제거

2. **실행 속도 최적화**:
   - Alpine Linux 기반으로 이미지 크기 최소화
   - 불필요한 서비스 비활성화

## 8. 모니터링 및 로깅

1. **컨테이너 상태 모니터링**:
   - HEALTHCHECK 추가
   - 로그 포맷 표준화

2. **오류 추적**:
   - exec 오류 발생 시 상세 로그 출력
   - 플랫폼 정보 로깅

## 9. 롤백 계획

1. **이전 버전으로의 롤백**:
   - Docker 이미지 태그를 통한 버전 관리
   - 실행 환경별 테스트 완료된 이미지 보관

2. **문제 발생 시 대응 절차**:
   - 오류 로그 분석
   - 플랫폼 호환성 재검증
   - 필요한 경우 bash 설치 제거 및 sh로 전환