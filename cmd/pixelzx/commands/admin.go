package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// AdminCommand defines the admin command structure
var AdminCommand = &cli.Command{
	Name:  "admin",
	Usage: "Admin commands for PIXELZX node",
	Subcommands: []*cli.Command{
		{
			Name:  "status",
			Usage: "Show node status",
			Action: func(c *cli.Context) error {
				fmt.Println("Node Status: Running")
				fmt.Println("Network: pixelzx-mainnet")
				fmt.Println("Peers: 12")
				return nil
			},
		},
		{
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
						return nil
					},
				},
			},
		},
	},
}