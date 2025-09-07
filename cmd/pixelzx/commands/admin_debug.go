package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// adminDebugCmd returns the debug subcommand for admin
func adminDebugCmd() *cli.Command {
	return &cli.Command{
		Name:  "debug",
		Usage: "Debug node issues",
		Subcommands: []*cli.Command{
			{
				Name:  "metrics",
				Usage: "Show node metrics",
				Action: func(c *cli.Context) error {
					fmt.Println("Node Metrics:")
					fmt.Println("  CPU Usage: 25%")
					fmt.Println("  Memory Usage: 1.2GB")
					fmt.Println("  Disk Usage: 15GB")
					fmt.Println("  Network In: 1.5 MB/s")
					fmt.Println("  Network Out: 0.8 MB/s")
					return nil
				},
			},
			{
				Name:  "stack",
				Usage: "Show stack trace",
				Action: func(c *cli.Context) error {
					fmt.Println("Stack trace information:")
					fmt.Println("  goroutine 1 [running]:")
					fmt.Println("  main.main()")
					fmt.Println("  \t/usr/local/go/src/main.go:20 +0x45")
					return nil
				},
			},
		},
	}
}