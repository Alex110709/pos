package commands

import (
	"github.com/urfave/cli/v2"
)

// AdminCommand defines the admin command structure
var AdminCommand = &cli.Command{
	Name:  "admin",
	Usage: "Admin commands for PIXELZX node",
	Subcommands: []*cli.Command{
		adminStatusCmd(),
		adminBackupCmd(),
		adminRestoreCmd(),
		adminConfigCmd(),
		adminDebugCmd(),
		adminResetCmd(),
		adminPeerCmd(),
	},
}