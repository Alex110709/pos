package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// StakingCommand defines the staking command structure
var StakingCommand = &cli.Command{
	Name:  "staking",
	Usage: "Staking management commands",
	Subcommands: []*cli.Command{
		{
			Name:  "delegate",
			Usage: "Delegate tokens to a validator",
			Action: func(c *cli.Context) error {
				fmt.Println("Delegating tokens to validator...")
				fmt.Println("Delegation successful")
				return nil
			},
		},
		{
			Name:  "undelegate",
			Usage: "Undelegate tokens from a validator",
			Action: func(c *cli.Context) error {
				fmt.Println("Undelegating tokens from validator...")
				fmt.Println("Undelegation successful")
				return nil
			},
		},
	},
}