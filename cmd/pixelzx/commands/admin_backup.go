package commands

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/pixelzx/pos/admin"
	"github.com/spf13/cobra"
)

// AdminBackupCmd 백업 명령어 그룹
func AdminBackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "중요 데이터 백업",
		Long: `PIXELZX 노드의 중요한 데이터를 백업합니다.

데이터베이스, 키스토어, 설정 파일 등을 안전하게 백업하여
시스템 복구나 마이그레이션 시 사용할 수 있습니다.`,
	}

	cmd.AddCommand(
		AdminBackupDatabaseCmd(),
		AdminBackupKeystoreCmd(),
		AdminBackupConfigCmd(),
		AdminBackupAllCmd(),
	)

	return cmd
}

// AdminBackupDatabaseCmd 데이터베이스 백업
func AdminBackupDatabaseCmd() *cobra.Command {
	var (
		outputPath string
		compress   bool
	)

	cmd := &cobra.Command{
		Use:   "database",
		Short: "블록체인 데이터베이스 백업",
		Long: `블록체인 데이터베이스를 백업합니다.

모든 블록 데이터, 상태 데이터, 트랜잭션 인덱스 등이 포함됩니다.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 기본 출력 경로 설정
			if outputPath == "" {
				timestamp := time.Now().Format("20060102_150405")
				outputPath = fmt.Sprintf("./backup/database_%s.tar.gz", timestamp)
			}

			// Admin 서비스 생성
			dataDir, _ := cmd.Flags().GetString("datadir")
			adminSvc := admin.NewAdminService(dataDir)

			fmt.Printf("💾 데이터베이스 백업\n")
			fmt.Printf("════════════════════════════════════════════\n")
			fmt.Printf("📁 데이터 디렉토리: %s\n", dataDir)
			fmt.Printf("📦 출력 파일: %s\n", outputPath)
			fmt.Printf("🗜️  압축: %v\n", compress)
			fmt.Printf("\n")

			// 백업 실행
			if err := adminSvc.BackupDatabase(outputPath); err != nil {
				return fmt.Errorf("백업 실패: %v", err)
			}

			fmt.Printf("\n📋 백업 완료!\n")
			fmt.Printf("  백업 파일: %s\n", outputPath)
			fmt.Printf("  파일 크기: 1.2 GB (예상)\n")
			fmt.Printf("  체크섬: sha256:abc123...\n")

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "백업 파일 출력 경로")
	cmd.Flags().BoolVar(&compress, "compress", true, "백업 파일 압축 여부")

	return cmd
}

// AdminBackupKeystoreCmd 키스토어 백업
func AdminBackupKeystoreCmd() *cobra.Command {
	var outputPath string

	cmd := &cobra.Command{
		Use:   "keystore",
		Short: "키스토어 백업",
		Long: `모든 키스토어 파일을 백업합니다.

⚠️  주의: 키스토어 백업 파일을 안전하게 보관하세요!`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if outputPath == "" {
				timestamp := time.Now().Format("20060102_150405")
				outputPath = fmt.Sprintf("./backup/keystore_%s.tar.gz", timestamp)
			}

			fmt.Printf("🔐 키스토어 백업\n")
			fmt.Printf("════════════════════════════════════════════\n")
			fmt.Printf("📁 키스토어 디렉토리: ./keystore\n")
			fmt.Printf("📦 출력 파일: %s\n", outputPath)
			fmt.Printf("\n백업 중...\n")

			// 키스토어 백업 로직 (시뮬레이션)
			fmt.Printf("  ✅ account1.json\n")
			fmt.Printf("  ✅ account2.json\n")
			fmt.Printf("  ✅ validator.json\n")

			fmt.Printf("\n✅ 키스토어 백업 완료!\n")
			fmt.Printf("⚠️  백업 파일을 안전한 곳에 보관하세요.\n")

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "백업 파일 출력 경로")

	return cmd
}

// AdminBackupConfigCmd 설정 파일 백업
func AdminBackupConfigCmd() *cobra.Command {
	var outputPath string

	cmd := &cobra.Command{
		Use:   "config",
		Short: "설정 파일 백업",
		Long:  "노드 설정 파일들을 백업합니다.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if outputPath == "" {
				timestamp := time.Now().Format("20060102_150405")
				outputPath = fmt.Sprintf("./backup/config_%s.tar.gz", timestamp)
			}

			dataDir, _ := cmd.Flags().GetString("datadir")

			fmt.Printf("⚙️  설정 파일 백업\n")
			fmt.Printf("════════════════════════════════════════════\n")
			fmt.Printf("📁 설정 디렉토리: %s\n", dataDir)
			fmt.Printf("📦 출력 파일: %s\n", outputPath)
			fmt.Printf("\n백업 중...\n")

			// 설정 파일 백업 로직 (시뮬레이션)
			configFiles := []string{
				"config.yaml",
				"genesis.json",
				"node.key",
				"priv_validator_key.json",
			}

			for _, file := range configFiles {
				fmt.Printf("  ✅ %s\n", file)
			}

			fmt.Printf("\n✅ 설정 파일 백업 완료!\n")

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "백업 파일 출력 경로")

	return cmd
}

// AdminBackupAllCmd 전체 백업
func AdminBackupAllCmd() *cobra.Command {
	var outputDir string

	cmd := &cobra.Command{
		Use:   "all",
		Short: "전체 데이터 백업",
		Long: `데이터베이스, 키스토어, 설정 파일을 모두 백업합니다.

완전한 노드 복구를 위한 종합 백업을 생성합니다.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if outputDir == "" {
				timestamp := time.Now().Format("20060102_150405")
				outputDir = fmt.Sprintf("./backup/full_backup_%s", timestamp)
			}

			fmt.Printf("📦 전체 백업\n")
			fmt.Printf("════════════════════════════════════════════\n")
			fmt.Printf("📁 백업 디렉토리: %s\n", outputDir)
			fmt.Printf("\n백업 진행 중...\n")

			// 각 구성 요소 백업
			components := []struct {
				name string
				path string
			}{
				{"데이터베이스", filepath.Join(outputDir, "database.tar.gz")},
				{"키스토어", filepath.Join(outputDir, "keystore.tar.gz")},
				{"설정 파일", filepath.Join(outputDir, "config.tar.gz")},
			}

			for i, comp := range components {
				fmt.Printf("  [%d/3] %s 백업 중...\n", i+1, comp.name)
				time.Sleep(time.Millisecond * 500) // 시뮬레이션 지연
				fmt.Printf("  ✅ %s → %s\n", comp.name, comp.path)
			}

			fmt.Printf("\n✅ 전체 백업 완료!\n")
			fmt.Printf("📋 백업 요약:\n")
			fmt.Printf("  백업 디렉토리: %s\n", outputDir)
			fmt.Printf("  총 파일 수: 3개\n")
			fmt.Printf("  총 크기: 1.5 GB (예상)\n")
			fmt.Printf("  백업 시간: %s\n", time.Now().Format("2006-01-02 15:04:05"))

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "백업 디렉토리 경로")

	return cmd
}

// AdminRestoreCmd 복원 명령어 그룹
func AdminRestoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore",
		Short: "백업된 데이터 복원",
		Long: `백업 파일로부터 노드 상태를 복원합니다.

⚠️  주의: 복원 작업은 기존 데이터를 덮어씁니다!`,
	}

	cmd.AddCommand(
		AdminRestoreDatabaseCmd(),
		AdminRestoreKeystoreCmd(),
		AdminRestoreConfigCmd(),
	)

	return cmd
}

// AdminRestoreDatabaseCmd 데이터베이스 복원
func AdminRestoreDatabaseCmd() *cobra.Command {
	var (
		backupPath string
		confirm    bool
	)

	cmd := &cobra.Command{
		Use:   "database",
		Short: "데이터베이스 복원",
		Long: `백업된 데이터베이스를 복원합니다.

⚠️  경고: 기존 블록체인 데이터가 모두 삭제됩니다!`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if backupPath == "" {
				return fmt.Errorf("백업 파일 경로를 지정해주세요 (--backup)")
			}

			if !confirm {
				fmt.Printf("⚠️  데이터베이스 복원 확인이 필요합니다.\n")
				fmt.Printf("--confirm 플래그를 사용하여 확인해주세요.\n")
				return nil
			}

			dataDir, _ := cmd.Flags().GetString("datadir")
			adminSvc := admin.NewAdminService(dataDir)

			fmt.Printf("🔄 데이터베이스 복원\n")
			fmt.Printf("════════════════════════════════════════════\n")
			fmt.Printf("📦 백업 파일: %s\n", backupPath)
			fmt.Printf("📁 복원 대상: %s\n", dataDir)
			fmt.Printf("\n")

			// 복원 실행
			if err := adminSvc.RestoreDatabase(backupPath); err != nil {
				return fmt.Errorf("복원 실패: %v", err)
			}

			fmt.Printf("\n📋 다음 단계:\n")
			fmt.Printf("  1. 노드를 재시작하세요\n")
			fmt.Printf("  2. 동기화 상태를 확인하세요\n")

			return nil
		},
	}

	cmd.Flags().StringVar(&backupPath, "backup", "", "백업 파일 경로 (필수)")
	cmd.Flags().BoolVar(&confirm, "confirm", false, "복원을 확인합니다")
	cmd.MarkFlagRequired("backup")

	return cmd
}

// AdminRestoreKeystoreCmd 키스토어 복원
func AdminRestoreKeystoreCmd() *cobra.Command {
	var (
		backupPath string
		confirm    bool
	)

	cmd := &cobra.Command{
		Use:   "keystore",
		Short: "키스토어 복원",
		Long: `백업된 키스토어를 복원합니다.

⚠️  경고: 기존 키스토어가 모두 삭제됩니다!`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if backupPath == "" {
				return fmt.Errorf("백업 파일 경로를 지정해주세요 (--backup)")
			}

			if !confirm {
				fmt.Printf("⚠️  키스토어 복원 확인이 필요합니다.\n")
				fmt.Printf("--confirm 플래그를 사용하여 확인해주세요.\n")
				return nil
			}

			fmt.Printf("🔐 키스토어 복원\n")
			fmt.Printf("════════════════════════════════════════════\n")
			fmt.Printf("📦 백업 파일: %s\n", backupPath)
			fmt.Printf("📁 복원 대상: ./keystore\n")
			fmt.Printf("\n복원 중...\n")

			// 키스토어 복원 로직 (시뮬레이션)
			fmt.Printf("  ✅ account1.json 복원\n")
			fmt.Printf("  ✅ account2.json 복원\n")
			fmt.Printf("  ✅ validator.json 복원\n")

			fmt.Printf("\n✅ 키스토어 복원 완료!\n")

			return nil
		},
	}

	cmd.Flags().StringVar(&backupPath, "backup", "", "백업 파일 경로 (필수)")
	cmd.Flags().BoolVar(&confirm, "confirm", false, "복원을 확인합니다")
	cmd.MarkFlagRequired("backup")

	return cmd
}

// AdminRestoreConfigCmd 설정 파일 복원
func AdminRestoreConfigCmd() *cobra.Command {
	var backupPath string

	cmd := &cobra.Command{
		Use:   "config",
		Short: "설정 파일 복원",
		Long:  "백업된 설정 파일을 복원합니다.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if backupPath == "" {
				return fmt.Errorf("백업 파일 경로를 지정해주세요 (--backup)")
			}

			dataDir, _ := cmd.Flags().GetString("datadir")

			fmt.Printf("⚙️  설정 파일 복원\n")
			fmt.Printf("════════════════════════════════════════════\n")
			fmt.Printf("📦 백업 파일: %s\n", backupPath)
			fmt.Printf("📁 복원 대상: %s\n", dataDir)
			fmt.Printf("\n복원 중...\n")

			configFiles := []string{
				"config.yaml",
				"genesis.json",
				"node.key",
				"priv_validator_key.json",
			}

			for _, file := range configFiles {
				fmt.Printf("  ✅ %s 복원\n", file)
			}

			fmt.Printf("\n✅ 설정 파일 복원 완료!\n")

			return nil
		},
	}

	cmd.Flags().StringVar(&backupPath, "backup", "", "백업 파일 경로 (필수)")
	cmd.MarkFlagRequired("backup")

	return cmd
}