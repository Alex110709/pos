#!/bin/bash
# PIXELZX Docker 권한 문제 해결 스크립트

set -e

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 로고 출력
echo -e "${CYAN}"
echo "╔═══════════════════════════════════════════════╗"
echo "║              PIXELZX POS EVM                  ║"
echo "║         Docker 권한 문제 해결 도구            ║"
echo "╚═══════════════════════════════════════════════╝"
echo -e "${NC}"

# 함수 정의
show_help() {
    echo -e "${BLUE}사용법:${NC}"
    echo "  $0 [명령어]"
    echo ""
    echo -e "${BLUE}명령어:${NC}"
    echo -e "  ${GREEN}check${NC}          - 현재 권한 상태 확인"
    echo -e "  ${GREEN}fix${NC}            - 권한 문제 자동 수정"
    echo -e "  ${GREEN}init${NC}           - Docker로 체인 초기화"
    echo -e "  ${GREEN}start${NC}          - Docker 노드 시작"
    echo -e "  ${GREEN}stop${NC}           - Docker 노드 중지"
    echo -e "  ${GREEN}logs${NC}           - 로그 확인"
    echo -e "  ${GREEN}shell${NC}          - 컨테이너 쉘 접속"
    echo -e "  ${GREEN}clean${NC}          - 모든 데이터 정리"
    echo -e "  ${GREEN}status${NC}         - 서비스 상태 확인"
    echo -e "  ${GREEN}help${NC}           - 이 도움말 표시"
    echo ""
    echo -e "${YELLOW}예시:${NC}"
    echo "  $0 check    # 권한 상태 확인"
    echo "  $0 fix      # 권한 문제 수정"
    echo "  $0 init     # 체인 초기화"
    echo "  $0 start    # 노드 시작"
}

check_permissions() {
    echo -e "${BLUE}🔍 권한 상태 확인 중...${NC}"
    
    # 현재 사용자 정보
    echo -e "\n${CYAN}현재 사용자 정보:${NC}"
    echo "  사용자: $(whoami) (UID: $(id -u), GID: $(id -g))"
    
    # Docker 상태 확인
    echo -e "\n${CYAN}Docker 상태:${NC}"
    if ! command -v docker &> /dev/null; then
        echo -e "  ${RED}❌ Docker가 설치되지 않았습니다${NC}"
        return 1
    else
        echo -e "  ${GREEN}✅ Docker 설치됨${NC} (버전: $(docker --version | cut -d' ' -f3))"
    fi
    
    if ! docker info &> /dev/null; then
        echo -e "  ${RED}❌ Docker 데몬이 실행되지 않았습니다${NC}"
        return 1
    else
        echo -e "  ${GREEN}✅ Docker 데몬 실행 중${NC}"
    fi
    
    # 디렉터리 상태 확인
    echo -e "\n${CYAN}디렉터리 상태:${NC}"
    for dir in "data" "keystore" "logs"; do
        if [ -d "./$dir" ]; then
            owner=$(ls -ld "./$dir" | awk '{print $3":"$4}')
            perms=$(ls -ld "./$dir" | awk '{print $1}')
            echo -e "  ${GREEN}✅ ./$dir${NC} (소유자: $owner, 권한: $perms)"
        else
            echo -e "  ${YELLOW}⚠️  ./$dir 디렉터리가 존재하지 않음${NC}"
        fi
    done
    
    # Docker 볼륨 상태 확인
    echo -e "\n${CYAN}Docker 볼륨 상태:${NC}"
    for volume in "pixelzx-data" "pixelzx-keystore" "pixelzx-logs"; do
        if docker volume inspect $volume &> /dev/null; then
            echo -e "  ${GREEN}✅ $volume${NC}"
        else
            echo -e "  ${YELLOW}⚠️  $volume 볼륨이 존재하지 않음${NC}"
        fi
    done
}

fix_permissions() {
    echo -e "${BLUE}🔧 권한 문제 수정 중...${NC}"
    
    # 디렉터리 생성
    echo -e "\n${CYAN}📁 디렉터리 생성...${NC}"
    mkdir -p ./data ./keystore ./logs
    
    # 권한 설정
    echo -e "${CYAN}🔒 권한 설정...${NC}"
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        sudo chown -R $(id -u):$(id -g) ./data ./keystore ./logs
        chmod -R 755 ./data ./keystore ./logs
    else
        # Linux
        sudo chown -R $(id -u):$(id -g) ./data ./keystore ./logs
        chmod -R 755 ./data ./keystore ./logs
    fi
    
    echo -e "${GREEN}✅ 권한 수정 완료${NC}"
}

init_chain() {
    echo -e "${BLUE}🚀 PIXELZX 체인 초기화...${NC}"
    
    # 권한 먼저 수정
    fix_permissions
    
    # Docker로 초기화 실행
    echo -e "\n${CYAN}🐳 Docker 초기화 실행...${NC}"
    docker-compose --profile init up pixelzx-init
    
    echo -e "${GREEN}✅ 체인 초기화 완료${NC}"
}

start_node() {
    echo -e "${BLUE}🚀 PIXELZX 노드 시작...${NC}"
    
    # 초기화가 되어있는지 확인
    if [ ! -f "./data/genesis.json" ]; then
        echo -e "${YELLOW}⚠️  제네시스 파일이 없습니다. 먼저 초기화를 실행합니다.${NC}"
        init_chain
    fi
    
    docker-compose up -d pixelzx-node
    echo -e "${GREEN}✅ 노드 시작 완료${NC}"
    
    # 상태 확인
    echo -e "\n${CYAN}📊 서비스 상태:${NC}"
    docker-compose ps
}

stop_node() {
    echo -e "${BLUE}🛑 PIXELZX 노드 중지...${NC}"
    docker-compose down
    echo -e "${GREEN}✅ 노드 중지 완료${NC}"
}

show_logs() {
    echo -e "${BLUE}📝 로그 확인...${NC}"
    docker-compose logs -f pixelzx-node
}

enter_shell() {
    echo -e "${BLUE}🐚 컨테이너 쉘 접속...${NC}"
    docker-compose exec pixelzx-node sh
}

clean_all() {
    echo -e "${RED}⚠️  모든 데이터를 삭제합니다. 계속하시겠습니까? (y/N)${NC}"
    read -r response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        echo -e "${BLUE}🧹 데이터 정리 중...${NC}"
        
        # 컨테이너 중지 및 삭제
        docker-compose down -v
        
        # 로컬 디렉터리 삭제
        sudo rm -rf ./data ./keystore ./logs
        
        # Docker 볼륨 삭제
        docker volume rm pixelzx-data pixelzx-keystore pixelzx-logs 2>/dev/null || true
        
        echo -e "${GREEN}✅ 데이터 정리 완료${NC}"
    else
        echo -e "${YELLOW}작업 취소됨${NC}"
    fi
}

show_status() {
    echo -e "${BLUE}📊 서비스 상태 확인...${NC}"
    
    echo -e "\n${CYAN}Docker Compose 서비스:${NC}"
    docker-compose ps
    
    echo -e "\n${CYAN}Docker 볼륨:${NC}"
    docker volume ls | grep pixelzx || echo "  볼륨 없음"
    
    echo -e "\n${CYAN}포트 사용 상태:${NC}"
    for port in 8545 8546 30303 6060; do
        if lsof -i :$port &> /dev/null; then
            echo -e "  ${GREEN}✅ 포트 $port 사용 중${NC}"
        else
            echo -e "  ${YELLOW}⚠️  포트 $port 사용 안됨${NC}"
        fi
    done
}

# 메인 로직
case "${1:-help}" in
    "check")
        check_permissions
        ;;
    "fix")
        fix_permissions
        ;;
    "init")
        init_chain
        ;;
    "start")
        start_node
        ;;
    "stop")
        stop_node
        ;;
    "logs")
        show_logs
        ;;
    "shell")
        enter_shell
        ;;
    "clean")
        clean_all
        ;;
    "status")
        show_status
        ;;
    "help"|*)
        show_help
        ;;
esac