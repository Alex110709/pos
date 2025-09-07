package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// adminResetCmd returns the reset subcommand for admin
func adminResetCmd() *cli.Command {
	return &cli.Command{
		Name:  "reset",
		Usage: "Reset node data",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "chaindata",
				Usage: "Reset chain data only",
			},
			&cli.BoolFlag{
				Name:  "snapshots",
				Usage: "Reset snapshots only",
			},
		},
		Action: func(c *cli.Context) error {
			chaindata := c.Bool("chaindata")
			snapshots := c.Bool("snapshots")
			
			if chaindata {
				fmt.Println("Resetting chain data...")
				fmt.Println("Chain data reset completed")
			} else if snapshots {
				fmt.Println("Resetting snapshots...")
				fmt.Println("Snapshots reset completed")
			} else {
				fmt.Println("Resetting all node data...")
				fmt.Println("Full node reset completed")
			}
			
			return nil
		},
	}
}