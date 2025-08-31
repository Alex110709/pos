package main

import (
	"fmt"
	"os"

	"github.com/pixelzx/pos/cmd/pixelzx/commands"
	"github.com/spf13/cobra"
)

const (
	AppName    = "pixelzx"
	AppVersion = "1.0.0"
	AppDesc    = "PIXELZX POS EVM Chain Node"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     AppName,
		Short:   AppDesc,
		Version: AppVersion,
		Long: fmt.Sprintf(`%s v%s

PIXELZX는 Proof of Stake (PoS) 기반의 EVM 호환 블록체인입니다.
PIXELZX 토큰을 네이티브 토큰으로 사용하며, 높은 성능과 낮은 수수료를 제공합니다.

주요 특징:
- PoS 합의 메커니즘
- EVM 완전 호환
- 3초 블록 타임
- 1000+ TPS 처리량
- 낮은 가스 수수료`, AppDesc, AppVersion),
	}

	// 글로벌 플래그 추가
	rootCmd.PersistentFlags().String("config", "", "설정 파일 경로")
	rootCmd.PersistentFlags().String("datadir", "./data", "데이터 디렉토리 경로")
	rootCmd.PersistentFlags().String("log-level", "info", "로그 레벨 (debug, info, warn, error)")
	rootCmd.PersistentFlags().Bool("testnet", false, "테스트넷 모드로 실행")

	// 하위 명령어 추가
	rootCmd.AddCommand(
		commands.InitCmd(),
		commands.StartCmd(),
		commands.ValidatorCmd(),
		commands.StakingCmd(),
		commands.GovernanceCmd(),
		commands.AccountCmd(),
		commands.ConfigCmd(),
		commands.AdminCmd(),
		commands.VersionCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}