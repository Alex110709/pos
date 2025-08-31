#!/bin/sh
# validate-multiarch.sh - PIXELZX 멀티 아키텍처 이미지 검증 스크립트

set -e

echo "=== PIXELZX 멀티 아키텍처 이미지 검증 스크립트 ==="
echo ""

# 1. 시스템 정보
echo "1. 시스템 정보 수집"
echo "   - 시스템 아키텍처: $(uname -m)"
echo "   - 운영체제: $(uname -s)"
echo ""

# 2. Docker 환경 확인
echo "2. Docker 환경 확인"
echo "   - Docker 버전: $(docker version --format '{{.Client.Version}}' 2>/dev/null || echo '미설치')"
echo "   - Buildx 지원: $(docker buildx version 2>/dev/null || echo '미지원')"
echo ""

# 3. 이미지 매니페스트 확인
echo "3. PIXELZX 이미지 매니페스트 검사"
if docker buildx imagetools inspect yuchanshin/pixelzx-evm:latest >/dev/null 2>&1; then
    echo "   ✅ 이미지 매니페스트 확인 성공"
    echo ""
    echo "   지원하는 플랫폼:"
    docker buildx imagetools inspect yuchanshin/pixelzx-evm:latest | grep "Platform:" | sed 's/^/     /' | head -3
else
    echo "   ❌ 이미지 매니페스트 확인 실패"
    echo "   💡 Docker Hub에 이미지가 존재하는지 확인하세요"
    exit 1
fi
echo ""

# 4. 플랫폼별 실행 테스트
echo "4. 플랫폼별 실행 테스트"
for platform in "linux/amd64" "linux/arm64" "linux/arm/v7"; do
    echo "   테스트 중: $platform"
    # version 명령어 대신 간단한 테스트
    if docker run --platform $platform --rm yuchanshin/pixelzx-evm:latest >/dev/null 2>&1; then
        echo "     ✅ $platform: 성공"
    else
        echo "     ❌ $platform: 실패"
        echo "     💡 해당 플랫폼의 이미지가 누락되었을 수 있습니다"
    fi
done
echo ""

# 5. 플랫폼 자동 감지 테스트
echo "5. 플랫폼 자동 감지 스크립트 테스트"
if [ -f "./scripts/detect-platform.sh" ]; then
    echo "   플랫폼 감지 스크립트 존재 확인: ✅"
    if ./scripts/detect-platform.sh version >/dev/null 2>&1; then
        echo "   자동 감지 스크립트 실행 테스트: ✅"
    else
        echo "   자동 감지 스크립트 실행 테스트: ❌"
        echo "   💡 스크립트에 실행 권한이 있는지 확인하세요 (chmod +x)"
    fi
else
    echo "   플랫폼 감지 스크립트 존재 확인: ❌"
    echo "   💡 scripts/detect-platform.sh 파일이 존재하지 않습니다"
fi
echo ""

echo "=== 검증 완료 ==="
echo ""
echo "📋 요약:"
echo "  - 시스템 아키텍처: $(uname -m)"
echo "  - Docker 버전: $(docker version --format '{{.Client.Version}}' 2>/dev/null || echo 'N/A')"
echo ""
echo "💡 다음 단계:"
echo "  1. 문제가 발생한 경우 문제 해결 가이드를 참고하세요 (EXEC_FORMAT_ERROR_SOLUTION.md)"
echo "  2. 추가 지원이 필요하면 시스템 정보와 함께 문의해 주세요"