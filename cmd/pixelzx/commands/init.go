package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

// Genesis represents the genesis block configuration
type Genesis struct {
	ChainID     uint64            `json:"chainId"`
	Timestamp   int64             `json:"timestamp"`
	Difficulty  string            `json:"difficulty"`
	GasLimit    string            `json:"gasLimit"`
	Alloc       map[string]Alloc  `json:"alloc"`
	Config      ChainConfig       `json:"config"`
	ExtraData   string            `json:"extraData"`
	Validators  []ValidatorInfo   `json:"validators"`
}

// Alloc represents initial account allocation
type Alloc struct {
	Balance string `json:"balance"`
}

// ChainConfig represents chain configuration
type ChainConfig struct {
	ChainID             uint64 `json:"chainId"`
	HomesteadBlock      uint64 `json:"homesteadBlock"`
	EIP150Block         uint64 `json:"eip150Block"`
	EIP155Block         uint64 `json:"eip155Block"`
	EIP158Block         uint64 `json:"eip158Block"`
	ByzantiumBlock      uint64 `json:"byzantiumBlock"`
	ConstantinopleBlock uint64 `json:"constantinopleBlock"`
	PetersburgBlock     uint64 `json:"petersburgBlock"`
	IstanbulBlock       uint64 `json:"istanbulBlock"`
	BerlinBlock         uint64 `json:"berlinBlock"`
	LondonBlock         uint64 `json:"londonBlock"`
	POS                 POSConfig `json:"pos"`
}

// POSConfig represents Proof of Stake configuration
type POSConfig struct {
	Period                uint64 `json:"period"`                // 블록 생성 간격 (초)
	Epoch                 uint64 `json:"epoch"`                 // 에포크 길이 (블록 수)
	MinValidatorStake     string `json:"minValidatorStake"`     // 최소 검증자 스테이킹
	MinDelegatorStake     string `json:"minDelegatorStake"`     // 최소 위임자 스테이킹
	MaxValidators         uint64 `json:"maxValidators"`         // 최대 검증자 수
	UnbondingPeriod       uint64 `json:"unbondingPeriod"`       // 언본딩 기간 (초)
	SlashingPenalty       string `json:"slashingPenalty"`       // 슬래싱 페널티 비율
}

// ValidatorInfo represents initial validator information
type ValidatorInfo struct {
	Address string `json:"address"`
	PubKey  string `json:"pubkey"`
	Power   string `json:"power"`
}

// PermissionValidator handles permission validation and directory creation
type PermissionValidator struct {
	dataDir    string
	currentUID int
	currentGID int
	userName   string
}

// UserInfo contains current user information
type UserInfo struct {
	UID      int
	GID      int
	Username string
	HomeDir  string
	IsRoot   bool
}

// DirectoryInfo contains directory access information
type DirectoryInfo struct {
	Exists      bool
	CanRead     bool
	CanWrite    bool
	CanExecute  bool
	Owner       string
	Group       string
	Permissions string
}

// InitCmd creates the init command
func InitCmd() *cobra.Command {
	var (
		genesisPath string
		chainID     uint64
		networkName string
	)

	cmd := &cobra.Command{
		Use:   "init [network-name]",
		Short: "PIXELZX 체인 초기화",
		Long: `PIXELZX POS EVM 체인을 초기화합니다.

제네시스 블록을 생성하고 초기 설정을 구성합니다.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				networkName = args[0]
			}

			dataDir, _ := cmd.Flags().GetString("datadir")
			
			return initChain(dataDir, genesisPath, chainID, networkName)
		},
	}

	cmd.Flags().StringVar(&genesisPath, "genesis", "", "제네시스 파일 경로")
	cmd.Flags().Uint64Var(&chainID, "chain-id", 8888, "체인 ID")
	cmd.Flags().StringVar(&networkName, "network", "pixelzx-pos", "네트워크 이름")

	return cmd
}

// NewPermissionValidator creates a new PermissionValidator instance
func NewPermissionValidator(dataDir string) (*PermissionValidator, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("현재 사용자 정보 조회 실패: %w", err)
	}

	uid, err := strconv.Atoi(currentUser.Uid)
	if err != nil {
		return nil, fmt.Errorf("UID 파싱 실패: %w", err)
	}

	gid, err := strconv.Atoi(currentUser.Gid)
	if err != nil {
		return nil, fmt.Errorf("GID 파싱 실패: %w", err)
	}

	return &PermissionValidator{
		dataDir:    dataDir,
		currentUID: uid,
		currentGID: gid,
		userName:   currentUser.Username,
	}, nil
}

// GetCurrentUserInfo returns information about the current user
func (pv *PermissionValidator) GetCurrentUserInfo() *UserInfo {
	currentUser, _ := user.Current()
	return &UserInfo{
		UID:      pv.currentUID,
		GID:      pv.currentGID,
		Username: pv.userName,
		HomeDir:  currentUser.HomeDir,
		IsRoot:   pv.currentUID == 0,
	}
}

// CheckDirectoryAccess checks if we can access the specified directory
func (pv *PermissionValidator) CheckDirectoryAccess(path string) (*DirectoryInfo, error) {
	info := &DirectoryInfo{}
	
	// Check if directory exists
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			info.Exists = false
			return info, nil
		}
		return nil, fmt.Errorf("디렉터리 정보 조회 실패: %w", err)
	}

	info.Exists = true
	info.Permissions = stat.Mode().String()

	// Check read access
	if _, err := os.Open(path); err == nil {
		info.CanRead = true
	}

	// Check write access by trying to create a temporary file
	tempFile := filepath.Join(path, ".pixelzx_permission_test")
	if file, err := os.Create(tempFile); err == nil {
		file.Close()
		os.Remove(tempFile)
		info.CanWrite = true
	}

	// Get file info for Unix systems
	if sysInfo, ok := stat.Sys().(*syscall.Stat_t); ok {
		info.Owner = strconv.Itoa(int(sysInfo.Uid))
		info.Group = strconv.Itoa(int(sysInfo.Gid))
	}

	return info, nil
}

// ValidatePermissions performs comprehensive permission validation
func (pv *PermissionValidator) ValidatePermissions() error {
	userInfo := pv.GetCurrentUserInfo()
	
	// Check if data directory exists
	dirInfo, err := pv.CheckDirectoryAccess(pv.dataDir)
	if err != nil {
		return fmt.Errorf("디렉터리 접근 검사 실패: %w", err)
	}

	if dirInfo.Exists {
		// Directory exists, check write permission
		if !dirInfo.CanWrite {
			return pv.createPermissionError("existing_dir_no_write", pv.dataDir, userInfo, dirInfo)
		}
	} else {
		// Directory doesn't exist, check parent directory
		parentDir := filepath.Dir(pv.dataDir)
		parentInfo, err := pv.CheckDirectoryAccess(parentDir)
		if err != nil {
			return fmt.Errorf("상위 디렉터리 접근 검사 실패: %w", err)
		}

		if !parentInfo.Exists {
			return pv.createPermissionError("parent_dir_not_exists", parentDir, userInfo, parentInfo)
		}

		if !parentInfo.CanWrite {
			return pv.createPermissionError("parent_dir_no_write", parentDir, userInfo, parentInfo)
		}
	}

	return nil
}

// createPermissionError creates detailed permission error with solutions
func (pv *PermissionValidator) createPermissionError(errorType, path string, userInfo *UserInfo, dirInfo *DirectoryInfo) error {
	errorMsg := "\n❌ 권한 오류가 발생했습니다!\n\n"
	
	// Error details
	errorMsg += fmt.Sprintf("🔍 문제 분석:\n")
	errorMsg += fmt.Sprintf("  - 현재 사용자: %s (UID: %d, GID: %d)\n", userInfo.Username, userInfo.UID, userInfo.GID)
	errorMsg += fmt.Sprintf("  - 대상 경로: %s\n", path)
	
	if dirInfo.Exists {
		errorMsg += fmt.Sprintf("  - 디렉터리 권한: %s\n", dirInfo.Permissions)
		if dirInfo.Owner != "" {
			errorMsg += fmt.Sprintf("  - 소유자: UID %s, GID %s\n", dirInfo.Owner, dirInfo.Group)
		}
	} else {
		errorMsg += fmt.Sprintf("  - 상태: 디렉터리가 존재하지 않음\n")
	}

	// Solutions based on error type
	errorMsg += "\n💡 해결 방법:\n"
	
	switch errorType {
	case "existing_dir_no_write":
		errorMsg += "  1. 관리자 권한으로 실행:\n"
		errorMsg += "     sudo ./pixelzx init\n\n"
		errorMsg += "  2. 디렉터리 소유권 변경:\n"
		errorMsg += fmt.Sprintf("     sudo chown -R %s:%s %s\n\n", userInfo.Username, userInfo.Username, path)
		errorMsg += "  3. 다른 디렉터리 사용:\n"
		errorMsg += fmt.Sprintf("     ./pixelzx init --datadir %s/pixelzx-data\n", userInfo.HomeDir)
		errorMsg += "     또는\n"
		errorMsg += "     ./pixelzx init --datadir /tmp/pixelzx-data\n"
	
	case "parent_dir_not_exists":
		errorMsg += "  1. 상위 디렉터리 생성:\n"
		errorMsg += fmt.Sprintf("     mkdir -p %s\n\n", filepath.Dir(path))
		errorMsg += "  2. 권한 설정 후 재시도:\n"
		errorMsg += fmt.Sprintf("     chmod 755 %s\n", filepath.Dir(path))
		errorMsg += "     ./pixelzx init\n"
	
	case "parent_dir_no_write":
		errorMsg += "  1. 상위 디렉터리 권한 변경:\n"
		errorMsg += fmt.Sprintf("     sudo chmod 755 %s\n\n", filepath.Dir(path))
		errorMsg += "  2. 홈 디렉터리 사용 (권장):\n"
		errorMsg += fmt.Sprintf("     ./pixelzx init --datadir %s/pixelzx-data\n\n", userInfo.HomeDir)
		errorMsg += "  3. 임시 디렉터리 사용 (테스트용):\n"
		errorMsg += "     ./pixelzx init --datadir /tmp/pixelzx-data\n"
	}

	// Docker specific guidance
	errorMsg += "\n🐳 Docker 환경에서 실행 중인 경우:\n"
	errorMsg += "  1. 호스트 볼륨 권한 설정:\n"
	errorMsg += "     sudo chown -R 1001:1001 ./data\n\n"
	errorMsg += "  2. Docker Compose 사용:\n"
	errorMsg += "     docker-compose up -d\n\n"
	errorMsg += "  3. 컨테이너 내부 경로 사용:\n"
	errorMsg += "     docker run -it yuchanshin/pixelzx-evm:latest init\n"

	return fmt.Errorf(errorMsg)
}

// CreateDirectorySafely creates directories with proper permissions
func (pv *PermissionValidator) CreateDirectorySafely(path string, perm os.FileMode) error {
	// Validate permissions first
	if err := pv.ValidatePermissions(); err != nil {
		return err
	}

	// Create directory with proper permissions
	if err := os.MkdirAll(path, perm); err != nil {
		// If creation fails, provide detailed error
		userInfo := pv.GetCurrentUserInfo()
		dirInfo, _ := pv.CheckDirectoryAccess(filepath.Dir(path))
		return pv.createPermissionError("creation_failed", path, userInfo, dirInfo)
	}

	// Verify creation was successful
	if info, err := pv.CheckDirectoryAccess(path); err != nil || !info.Exists {
		return fmt.Errorf("디렉터리 생성 검증 실패: %s", path)
	}

	return nil
}

func initChain(dataDir, genesisPath string, chainID uint64, networkName string) error {
	fmt.Printf("PIXELZX POS EVM 체인 초기화 중...\n")
	fmt.Printf("데이터 디렉토리: %s\n", dataDir)
	fmt.Printf("체인 ID: %d\n", chainID)
	fmt.Printf("네트워크 이름: %s\n", networkName)

	// 권한 검증자 생성
	validator, err := NewPermissionValidator(dataDir)
	if err != nil {
		return fmt.Errorf("권한 검증자 생성 실패: %w", err)
	}

	fmt.Printf("\n🔍 권한 검사 중...\n")
	
	// 사전 권한 검증
	if err := validator.ValidatePermissions(); err != nil {
		return err // 상세한 오류 메시지는 ValidatePermissions에서 제공
	}

	fmt.Printf("✅ 권한 검사 완료\n")

	// 안전한 데이터 디렉터리 생성
	fmt.Printf("\n📁 데이터 디렉터리 생성 중...\n")
	if err := validator.CreateDirectorySafely(dataDir, 0755); err != nil {
		return fmt.Errorf("데이터 디렉터리 생성 실패: %w", err)
	}

	// 키스토어 디렉터리 생성
	fmt.Printf("🔐 키스토어 디렉터리 생성 중...\n")
	keystoreDir := filepath.Join(dataDir, "keystore")
	if err := os.MkdirAll(keystoreDir, 0755); err != nil {
		return fmt.Errorf("키스토어 디렉터리 생성 실패: %w", err)
	}

	// 로그 디렉터리 생성
	fmt.Printf("📝 로그 디렉터리 생성 중...\n")
	logDir := filepath.Join(dataDir, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("로그 디렉터리 생성 실패: %w", err)
	}

	// 제네시스 파일 생성 또는 복사
	fmt.Printf("🌍 제네시스 파일 처리 중...\n")
	var genesis *Genesis

	if genesisPath != "" {
		genesis, err = loadGenesis(genesisPath)
		if err != nil {
			return fmt.Errorf("제네시스 파일 로드 실패: %w", err)
		}
		fmt.Printf("✅ 사용자 제네시스 파일 로드 완료: %s\n", genesisPath)
	} else {
		genesis = createDefaultGenesis(chainID, networkName)
		fmt.Printf("✅ 기본 제네시스 생성 완료\n")
	}

	// 제네시스 파일 저장
	genesisFile := filepath.Join(dataDir, "genesis.json")
	if err := saveGenesis(genesis, genesisFile); err != nil {
		return fmt.Errorf("제네시스 파일 저장 실패: %w", err)
	}

	// 설정 파일 생성
	fmt.Printf("⚙️ 설정 파일 생성 중...\n")
	configFile := filepath.Join(dataDir, "config.yaml")
	if err := createConfig(configFile, chainID, networkName); err != nil {
		return fmt.Errorf("설정 파일 생성 실패: %w", err)
	}

	// 초기화 완료 메시지
	fmt.Printf("\n✨ PIXELZX 체인 초기화 완료!\n")
	fmt.Printf("📁 제네시스 파일: %s\n", genesisFile)
	fmt.Printf("⚙️  설정 파일: %s\n", configFile)
	fmt.Printf("🔐 키스토어 디렉터리: %s\n", keystoreDir)
	fmt.Printf("📝 로그 디렉터리: %s\n", logDir)
	fmt.Printf("\n🚀 다음 명령어로 노드를 시작할 수 있습니다:\n")
	fmt.Printf("  pixelzx start --datadir %s\n", dataDir)
	fmt.Printf("\n📚 더 많은 정보는 README.md를 참고하세요.\n")

	return nil
}

func loadGenesis(path string) (*Genesis, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var genesis Genesis
	if err := json.Unmarshal(data, &genesis); err != nil {
		return nil, err
	}

	return &genesis, nil
}

func createDefaultGenesis(chainID uint64, networkName string) *Genesis {
	// PIXELZX 토큰 총 공급량: 10,000,000,000,000,000 PXZ (1경 PXZ, 1e34 wei)
	totalSupply := "10000000000000000000000000000000000"
	
	// 기본 계정에 전체 공급량 할당
	defaultAccount := "0x742d35Cc6635C0532925a3b8D5C0532925b8D5C05"

	return &Genesis{
		ChainID:    chainID,
		Timestamp:  time.Now().Unix(),
		Difficulty: "0x1",
		GasLimit:   "0x1c9c380", // 30,000,000
		Alloc: map[string]Alloc{
			defaultAccount: {
				Balance: totalSupply,
			},
		},
		Config: ChainConfig{
			ChainID:             chainID,
			HomesteadBlock:      0,
			EIP150Block:         0,
			EIP155Block:         0,
			EIP158Block:         0,
			ByzantiumBlock:      0,
			ConstantinopleBlock: 0,
			PetersburgBlock:     0,
			IstanbulBlock:       0,
			BerlinBlock:         0,
			LondonBlock:         0,
			POS: POSConfig{
				Period:                3,                                    // 3초 블록 타임
				Epoch:                 200,                                  // 200 블록 에포크
				MinValidatorStake:     "1000000000000000000000000000",           // 1,000,000,000 PXZ (10억 PXZ)
				MinDelegatorStake:     "10000000000000000000000",                // 10,000 PXZ
				MaxValidators:         125,                                  // 최대 125명 검증자
				UnbondingPeriod:       1814400,                              // 21일 (초)
				SlashingPenalty:       "50000000000000000",                  // 5% (0.05)
			},
		},
		ExtraData:  fmt.Sprintf("PIXELZX POS EVM Chain - %s", networkName),
		Validators: []ValidatorInfo{}, // 초기에는 검증자 없음
	}
}

func saveGenesis(genesis *Genesis, path string) error {
	data, err := json.MarshalIndent(genesis, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func createConfig(path string, chainID uint64, networkName string) error {
	config := fmt.Sprintf(`# PIXELZX POS EVM Chain Configuration

# Network Configuration
network:
  name: "%s"
  chain_id: %d
  block_time: 3s
  epoch_length: 200

# Node Configuration
node:
  name: "pixelzx-node"
  datadir: "./data"
  log_level: "info"

# P2P Network
p2p:
  port: 30303
  max_peers: 50
  bootnode_addrs: []

# JSON-RPC API
rpc:
  enabled: true
  host: "0.0.0.0"
  port: 8545
  cors_origins: ["*"]
  apis: ["eth", "net", "web3", "personal", "admin", "debug", "txpool", "pxz"]

# WebSocket API
websocket:
  enabled: true
  host: "0.0.0.0"
  port: 8546
  origins: ["*"]
  apis: ["eth", "net", "web3", "pxz"]

# Validator Configuration
validator:
  enabled: false
  address: ""
  password_file: ""

# Staking Configuration
staking:
  min_validator_stake: "1000000000000000000000000000"  # 1,000,000,000 PXZ (10억 PXZ)
  min_delegator_stake: "10000000000000000000000"       # 10,000 PXZ
  unbonding_period: "504h"                         # 21 days
  max_validators: 125

# Gas Configuration
gas:
  limit: 30000000
  price: 20000000000
  min_price: 1000000000

# Security
security:
  keystore_dir: "./keystore"

# Metrics
metrics:
  enabled: true
  host: "0.0.0.0"
  port: 6060

# Database
database:
  type: "leveldb"
  cache: 512
  handles: 256
`, networkName, chainID)

	return os.WriteFile(path, []byte(config), 0644)
}