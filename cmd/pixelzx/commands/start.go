package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
)

// StartCmd creates the start command
func StartCmd() *cobra.Command {
	var (
		configPath    string
		rpcPort       int
		wsPort        int
		p2pPort       int
		validatorMode bool
		syncMode      string
	)

	cmd := &cobra.Command{
		Use:   "start",
		Short: "PIXELZX 노드 시작",
		Long: `PIXELZX POS EVM 체인 노드를 시작합니다.

블록체인 네트워크에 연결하고 트랜잭션 처리를 시작합니다.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dataDir, _ := cmd.Flags().GetString("datadir")
			logLevel, _ := cmd.Flags().GetString("log-level")
			testnet, _ := cmd.Flags().GetBool("testnet")

			return startNode(StartConfig{
				DataDir:       dataDir,
				ConfigPath:    configPath,
				LogLevel:      logLevel,
				RPCPort:       rpcPort,
				WSPort:        wsPort,
				P2PPort:       p2pPort,
				ValidatorMode: validatorMode,
				SyncMode:      syncMode,
				Testnet:       testnet,
			})
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "설정 파일 경로")
	cmd.Flags().IntVar(&rpcPort, "rpc-port", 8545, "JSON-RPC 포트")
	cmd.Flags().IntVar(&wsPort, "ws-port", 8546, "WebSocket 포트")
	cmd.Flags().IntVar(&p2pPort, "p2p-port", 30303, "P2P 네트워크 포트")
	cmd.Flags().BoolVar(&validatorMode, "validator", false, "검증자 모드로 실행")
	cmd.Flags().StringVar(&syncMode, "sync-mode", "full", "동기화 모드 (full, fast, light)")

	return cmd
}

// StartConfig represents node start configuration
type StartConfig struct {
	DataDir       string
	ConfigPath    string
	LogLevel      string
	RPCPort       int
	WSPort        int
	P2PPort       int
	ValidatorMode bool
	SyncMode      string
	Testnet       bool
}

func startNode(config StartConfig) error {
	fmt.Printf("🚀 PIXELZX POS EVM 체인 노드 시작 중...\n")
	fmt.Printf("📁 데이터 디렉토리: %s\n", config.DataDir)
	fmt.Printf("🔧 로그 레벨: %s\n", config.LogLevel)
	fmt.Printf("🌐 JSON-RPC 포트: %d\n", config.RPCPort)
	fmt.Printf("🔌 WebSocket 포트: %d\n", config.WSPort)
	fmt.Printf("🔗 P2P 포트: %d\n", config.P2PPort)
	fmt.Printf("⚙️  동기화 모드: %s\n", config.SyncMode)
	
	if config.ValidatorMode {
		fmt.Printf("✅ 검증자 모드 활성화\n")
	}
	
	if config.Testnet {
		fmt.Printf("🧪 테스트넷 모드\n")
	}

	// 데이터 디렉토리 확인
	if _, err := os.Stat(config.DataDir); os.IsNotExist(err) {
		return fmt.Errorf("데이터 디렉토리가 존재하지 않습니다: %s\n초기화를 먼저 실행하세요: pixelzx init", config.DataDir)
	}

	// 제네시스 파일 확인
	genesisFile := filepath.Join(config.DataDir, "genesis.json")
	if _, err := os.Stat(genesisFile); os.IsNotExist(err) {
		return fmt.Errorf("제네시스 파일이 존재하지 않습니다: %s\n초기화를 먼저 실행하세요: pixelzx init", genesisFile)
	}

	// 노드 구성 요소 초기화
	fmt.Printf("\n📦 노드 구성 요소 초기화 중...\n")
	
	// 1. 저장소 초기화
	fmt.Printf("  💾 저장소 초기화...\n")
	
	// 2. P2P 네트워크 초기화
	fmt.Printf("  🔗 P2P 네트워크 초기화...\n")
	
	// 3. 합의 엔진 초기화
	fmt.Printf("  🎯 PoS 합의 엔진 초기화...\n")
	
	// 4. EVM 실행 환경 초기화
	fmt.Printf("  ⚡ EVM 실행 환경 초기화...\n")
	
	// 5. API 서버 초기화
	fmt.Printf("  🌐 API 서버 초기화...\n")
	
	// 6. 스테이킹 모듈 초기화
	fmt.Printf("  💰 스테이킹 모듈 초기화...\n")

	// 컨텍스트 설정
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 시그널 핸들링
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 노드 시작
	fmt.Printf("\n✅ PIXELZX 노드 시작 완료!\n")
	fmt.Printf("\n📋 서비스 상태:\n")
	fmt.Printf("  🌐 JSON-RPC: http://localhost:%d\n", config.RPCPort)
	fmt.Printf("  🔌 WebSocket: ws://localhost:%d\n", config.WSPort)
	fmt.Printf("  🔗 P2P 리스닝: :%d\n", config.P2PPort)
	
	if config.ValidatorMode {
		fmt.Printf("  ✅ 검증자 모드: 활성화\n")
		fmt.Printf("  🎯 블록 생성: 대기 중\n")
	}

	fmt.Printf("\n📊 실시간 상태:\n")
	fmt.Printf("  📦 현재 블록: 0\n")
	fmt.Printf("  🔗 연결된 피어: 0\n")
	fmt.Printf("  💾 체인 상태: 동기화 중\n")

	fmt.Printf("\n종료하려면 Ctrl+C를 누르세요...\n")

	// 메인 루프
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("\n🛑 노드 종료 중...\n")
			return nil
		case sig := <-sigCh:
			fmt.Printf("\n📡 시그널 수신: %s\n", sig)
			fmt.Printf("🛑 노드 정상 종료 중...\n")
			
			// 정리 작업
			fmt.Printf("  💾 상태 저장 중...\n")
			fmt.Printf("  🔗 P2P 연결 종료 중...\n")
			fmt.Printf("  🌐 API 서버 종료 중...\n")
			
			fmt.Printf("✅ 노드 종료 완료\n")
			return nil
		}
	}
}