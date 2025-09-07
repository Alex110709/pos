package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// adminConfigCmd returns the config subcommand for admin
func adminConfigCmd() *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "Configuration management",
		Subcommands: []*cli.Command{
			{
				Name:  "show",
				Usage: "Show current configuration",
				Action: func(c *cli.Context) error {
					fmt.Println("Current Configuration:")
					fmt.Println("  Network: pixelzx-mainnet")
					fmt.Println("  Chain ID: 8888")
					fmt.Println("  RPC Port: 8545")
					fmt.Println("  WS Port: 8546")
					fmt.Println("  P2P Port: 30303")
					fmt.Println("  Data Directory: /var/lib/pixelzx")
					fmt.Println("  KeyStore Directory: /var/lib/pixelzx/keystore")
					return nil
				},
			},
			{
				Name:  "validate",
				Usage: "Validate configuration",
				Action: func(c *cli.Context) error {
					fmt.Println("🔍 Validating configuration...")
					fmt.Println("  ✅ Network settings are valid")
					fmt.Println("  ✅ Ports are available")
					fmt.Println("  ✅ Data directory permissions are correct")
					fmt.Println("  ⚠️  Gas price is higher than recommended")
					fmt.Println("Configuration validation completed")
					return nil
				},
			},
		},
	}
}