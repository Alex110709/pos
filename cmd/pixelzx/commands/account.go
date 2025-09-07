package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// AccountCommand defines the account command structure
var AccountCommand = &cli.Command{
	Name:  "account",
	Usage: "Account management commands",
	Subcommands: []*cli.Command{
		{
			Name:  "new",
			Usage: "Create a new account",
			Action: func(c *cli.Context) error {
				fmt.Println("Creating new account...")
				fmt.Println("Account created: 0x742d35Cc6634C0532925a3b844Bc454e4438f44e")
				return nil
			},
		},
		{
			Name:  "list",
			Usage: "List all accounts",
			Action: func(c *cli.Context) error {
				fmt.Println("Accounts:")
				fmt.Println("  0x742d35Cc6634C0532925a3b844Bc454e4438f44e")
				fmt.Println("  0x8ba1f109551bD432803012645Hac136c22C43137")
				return nil
			},
		},
	},
}