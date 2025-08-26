package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	cmd.Flags().Uint64Var(&chainID, "chain-id", 1337, "체인 ID")
	cmd.Flags().StringVar(&networkName, "network", "pixelzx-pos", "네트워크 이름")

	return cmd
}

func initChain(dataDir, genesisPath string, chainID uint64, networkName string) error {
	fmt.Printf("PIXELZX POS EVM 체인 초기화 중...\n")
	fmt.Printf("데이터 디렉토리: %s\n", dataDir)
	fmt.Printf("체인 ID: %d\n", chainID)
	fmt.Printf("네트워크 이름: %s\n", networkName)

	// 데이터 디렉토리 생성
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("데이터 디렉토리 생성 실패: %w", err)
	}

	// 키스토어 디렉토리 생성
	keystoreDir := filepath.Join(dataDir, "keystore")
	if err := os.MkdirAll(keystoreDir, 0755); err != nil {
		return fmt.Errorf("키스토어 디렉토리 생성 실패: %w", err)
	}

	// 제네시스 파일 생성 또는 복사
	var genesis *Genesis
	var err error

	if genesisPath != "" {
		genesis, err = loadGenesis(genesisPath)
		if err != nil {
			return fmt.Errorf("제네시스 파일 로드 실패: %w", err)
		}
	} else {
		genesis = createDefaultGenesis(chainID, networkName)
	}

	// 제네시스 파일 저장
	genesisFile := filepath.Join(dataDir, "genesis.json")
	if err := saveGenesis(genesis, genesisFile); err != nil {
		return fmt.Errorf("제네시스 파일 저장 실패: %w", err)
	}

	// 설정 파일 생성
	configFile := filepath.Join(dataDir, "config.yaml")
	if err := createConfig(configFile, chainID, networkName); err != nil {
		return fmt.Errorf("설정 파일 생성 실패: %w", err)
	}

	fmt.Printf("✅ PIXELZX 체인 초기화 완료!\n")
	fmt.Printf("📁 제네시스 파일: %s\n", genesisFile)
	fmt.Printf("⚙️  설정 파일: %s\n", configFile)
	fmt.Printf("🔐 키스토어 디렉토리: %s\n", keystoreDir)
	fmt.Printf("\n다음 명령어로 노드를 시작할 수 있습니다:\n")
	fmt.Printf("  pixelzx start --datadir %s\n", dataDir)

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
	// PIXELZX 토큰 총 공급량: 1,000,000,000 PXZ (1e27 wei)
	totalSupply := "1000000000000000000000000000"
	
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
				MinValidatorStake:     "100000000000000000000000",           // 100,000 PXZ
				MinDelegatorStake:     "1000000000000000000",                // 1 PXZ
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
  min_validator_stake: "100000000000000000000000"  # 100,000 PXZ
  min_delegator_stake: "1000000000000000000"       # 1 PXZ
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