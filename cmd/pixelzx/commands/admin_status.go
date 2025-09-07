package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// adminStatusCmd returns the status subcommand for admin
func adminStatusCmd() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Show node status",
		Action: func(c *cli.Context) error {
			fmt.Println("Node Status: Running")
			fmt.Println("Network: pixelzx-mainnet")
			fmt.Println("Peers: 12")
			fmt.Println("Block Height: 12345")
			fmt.Println("Sync Status: Synced")
			return nil
		},
	}
}