package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/ethereum/go-ethereum/cmd/pixelzx/commands"
)

func main() {
	app := &cli.App{
		Name:  "pixelzx",
		Usage: "PIXELZX POS EVM Chain CLI",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Initialize the blockchain network",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "network",
						Usage: "Network name",
					},
				},
				Action: func(c *cli.Context) error {
					network := c.String("network")
					if network == "" {
						network = "default"
					}
					fmt.Printf("Initializing PIXELZX network: %s\n", network)
					// Implementation would go here
					return nil
				},
			},
			{
				Name:  "start",
				Usage: "Start the PIXELZX node",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "validator",
						Usage: "Run as validator node",
					},
				},
				Action: func(c *cli.Context) error {
					validator := c.Bool("validator")
					if validator {
						fmt.Println("Starting PIXELZX node in validator mode")
					} else {
						fmt.Println("Starting PIXELZX node")
					}
					// Implementation would go here
					return nil
				},
			},
			commands.AdminCommand,
			commands.AccountCommand,
			commands.ValidatorCommand,
			commands.StakingCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}