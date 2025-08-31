#!/bin/sh
# detect-platform.sh - PIXELZX 플랫폼 자동 감지 및 실행 스크립트

set -e

echo "🚀 PIXELZX POS EVM 체인 플랫폼 감지 및 실행 스크립트"
echo "==============================================="

# 현재 시스템 아키텍처 감지
ARCH=$(uname -m)
OS=$(uname -s)

echo "🖥️  시스템 정보:"
echo "  - 운영체제: $OS"
echo "  - 아키텍처: $ARCH"

# 플랫폼 매핑
case $ARCH in
    x86_64|amd64)
        PLATFORM="linux/amd64"
        echo "  - 매핑 플랫폼: $PLATFORM"
        ;;
    aarch64|arm64)
        PLATFORM="linux/arm64"
        echo "  - 매핑 플랫폼: $PLATFORM"
        ;;
    armv7l|armv7)
        PLATFORM="linux/arm/v7"
        echo "  - 매핑 플랫폼: $PLATFORM"
        ;;
    *)
        echo "❌ 지원하지 않는 아키텍처: $ARCH"
        echo "💡 지원 플랫폼: x86_64, aarch64, armv7l"
        exit 1
        ;;
esac

echo ""
echo "🐳 Docker 컨테이너 실행 중..."
echo "   플랫폼: $PLATFORM"
echo "   이미지: yuchanshin/pixelzx-evm:latest"
echo ""

# Docker 실행
if [ "$#" -eq 0 ]; then
    # 인자가 없으면 기본 실행 (도움말 표시)
    echo "📋 기본 명령어 실행 (도움말 표시)"
    docker run --platform $PLATFORM --rm yuchanshin/pixelzx-evm:latest
else
    # 인자가 있으면 해당 명령어로 pixelzx 실행
    echo "📋 PIXELZX 명령어 실행: $@"
    docker run --platform $PLATFORM --rm yuchanshin/pixelzx-evm:latest pixelzx "$@"
fi

echo ""
echo "✅ PIXELZX 노드 실행 완료"