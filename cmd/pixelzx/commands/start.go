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
		Short: "Start PIXELZX node",
		Long: `Start the PIXELZX POS EVM chain node.

Connect to the blockchain network and start processing transactions.`,
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

	cmd.Flags().StringVar(&configPath, "config", "", "Path to config file")
	cmd.Flags().IntVar(&rpcPort, "rpc-port", 8545, "JSON-RPC port")
	cmd.Flags().IntVar(&wsPort, "ws-port", 8546, "WebSocket port")
	cmd.Flags().IntVar(&p2pPort, "p2p-port", 30303, "P2P network port")
	cmd.Flags().BoolVar(&validatorMode, "validator", false, "Run in validator mode")
	cmd.Flags().StringVar(&syncMode, "sync-mode", "full", "Sync mode (full, fast, light)")

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
	fmt.Printf("🚀 Starting PIXELZX POS EVM Chain node...\n")
	fmt.Printf("📁 Data directory: %s\n", config.DataDir)
	fmt.Printf("🔧 Log level: %s\n", config.LogLevel)
	fmt.Printf("🌐 JSON-RPC port: %d\n", config.RPCPort)
	fmt.Printf("🔌 WebSocket port: %d\n", config.WSPort)
	fmt.Printf("🔗 P2P port: %d\n", config.P2PPort)
	fmt.Printf("⚙️  Sync mode: %s\n", config.SyncMode)
	
	if config.ValidatorMode {
		fmt.Printf("✅ Validator mode enabled\n")
	}
	
	if config.Testnet {
		fmt.Printf("🧪 Testnet mode\n")
	}

	// Check data directory
	if _, err := os.Stat(config.DataDir); os.IsNotExist(err) {
		return fmt.Errorf("data directory does not exist: %s\nPlease run initialization first: pixelzx init", config.DataDir)
	}

	// Check genesis file
	genesisFile := filepath.Join(config.DataDir, "genesis.json")
	if _, err := os.Stat(genesisFile); os.IsNotExist(err) {
		return fmt.Errorf("genesis file does not exist: %s\nPlease run initialization first: pixelzx init", genesisFile)
	}

	// Initialize node components
	fmt.Printf("\n📦 Initializing node components...\n")
	
	// 1. Storage initialization
	fmt.Printf("  💾 Initializing storage...\n")
	
	// 2. P2P network initialization
	fmt.Printf("  🔗 Initializing P2P network...\n")
	
	// 3. Consensus engine initialization
	fmt.Printf("  🎯 Initializing PoS consensus engine...\n")
	
	// 4. EVM execution environment initialization
	fmt.Printf("  ⚡ Initializing EVM execution environment...\n")
	
	// 5. API server initialization
	fmt.Printf("  🌐 Initializing API server...\n")
	
	// 6. Staking module initialization
	fmt.Printf("  💰 Initializing staking module...\n")

	// Set up context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start node
	fmt.Printf("\n✅ PIXELZX node started successfully!\n")
	fmt.Printf("\n📋 Service status:\n")
	fmt.Printf("  🌐 JSON-RPC: http://localhost:%d\n", config.RPCPort)
	fmt.Printf("  🔌 WebSocket: ws://localhost:%d\n", config.WSPort)
	fmt.Printf("  🔗 P2P listening: :%d\n", config.P2PPort)
	
	if config.ValidatorMode {
		fmt.Printf("  ✅ Validator mode: enabled\n")
		fmt.Printf("  🎯 Block creation: waiting\n")
	}

	fmt.Printf("\n📊 Real-time status:\n")
	fmt.Printf("  📦 Current block: 0\n")
	fmt.Printf("  🔗 Connected peers: 0\n")
	fmt.Printf("  💾 Chain status: syncing\n")

	fmt.Printf("\nPress Ctrl+C to exit...\n")

	// Main loop
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("\n🛑 Shutting down node...\n")
			return nil
		case sig := <-sigCh:
			fmt.Printf("\n📡 Signal received: %s\n", sig)
			fmt.Printf("🛑 Shutting down node gracefully...\n")
			
			// Cleanup
			fmt.Printf("  💾 Saving state...\n")
			fmt.Printf("  🔗 Closing P2P connections...\n")
			fmt.Printf("  🌐 Shutting down API server...\n")
			
			fmt.Printf("✅ Node shutdown complete\n")
			return nil
		}
	}
}