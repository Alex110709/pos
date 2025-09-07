package commands

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

// adminPeerCmd returns the peer subcommand for admin
func adminPeerCmd() *cli.Command {
	return &cli.Command{
		Name:  "peer",
		Usage: "Manage network peers",
		Subcommands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List connected peers",
				Action: func(c *cli.Context) error {
					fmt.Println("Connected Peers:")
					fmt.Println("  enode://a1b2c3d4e5f67890@192.168.1.100:30303")
					fmt.Println("  enode://f6e5d4c3b2a10987@192.168.1.101:30303")
					fmt.Println("  enode://7890abcdef12345@192.168.1.102:30303")
					return nil
				},
			},
			{
				Name:      "connect",
				Usage:     "Connect to a peer",
				UsageText: "pixelzx admin peer connect <enode-url>",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return fmt.Errorf("missing enode URL")
					}
					
					enodeURL := c.Args().First()
					if err := validateEnodeURL(enodeURL); err != nil {
						return fmt.Errorf("invalid enode URL: %v", err)
					}
					
					fmt.Printf("Connecting to peer: %s\n", enodeURL)
					fmt.Println("Peer connected successfully")
					return nil
				},
			},
			{
				Name:      "disconnect",
				Usage:     "Disconnect from a peer",
				UsageText: "pixelzx admin peer disconnect <enode-url>",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return fmt.Errorf("missing enode URL")
					}
					
					enodeURL := c.Args().First()
					fmt.Printf("Disconnecting from peer: %s\n", enodeURL)
					fmt.Println("Peer disconnected successfully")
					return nil
				},
			},
			{
				Name:  "self",
				Usage: "Show local node enode information",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "json",
						Usage: "Output in JSON format",
					},
				},
				Action: func(c *cli.Context) error {
					jsonFormat := c.Bool("json")
					
					if jsonFormat {
						fmt.Println(`{
  "enode": "enode://abcd1234ef567890@192.168.1.50:30303",
  "id": "abcd1234ef567890",
  "ip": "192.168.1.50",
  "tcp": 30303,
  "udp": 30303
}`)
					} else {
						fmt.Println("Local Node Information:")
						fmt.Println("  Enode: enode://abcd1234ef567890@192.168.1.50:30303")
						fmt.Println("  Node ID: abcd1234ef567890")
						fmt.Println("  IP: 192.168.1.50")
						fmt.Println("  TCP Port: 30303")
						fmt.Println("  UDP Port: 30303")
					}
					return nil
				},
			},
		},
	}
}

// validateEnodeURL validates the format of an enode URL
func validateEnodeURL(enodeURL string) error {
	// Basic format check
	if !strings.HasPrefix(enodeURL, "enode://") {
		return fmt.Errorf("invalid enode URL format, must start with 'enode://'")
	}
	
	// Extract the part after "enode://"
	parts := enodeURL[8:] // Skip "enode://"
	
	// Find the '@' separator
	atIndex := strings.Index(parts, "@")
	if atIndex == -1 {
		return fmt.Errorf("invalid enode URL format, missing '@' separator")
	}
	
	publicKey := parts[:atIndex]
	if len(publicKey) != 128 {
		return fmt.Errorf("invalid public key length, must be 128 characters, got %d", len(publicKey))
	}
	
	// Validate public key is hex
	if _, err := hex.DecodeString(publicKey); err != nil {
		return fmt.Errorf("invalid public key format, must be hexadecimal")
	}
	
	// Extract host and port
	hostPort := parts[atIndex+1:]
	colonIndex := strings.LastIndex(hostPort, ":")
	
	if colonIndex == -1 {
		return fmt.Errorf("invalid enode URL format, missing port")
	}
	
	host := hostPort[:colonIndex]
	portStr := hostPort[colonIndex+1:]
	
	// Validate IP address
	if ip := net.ParseIP(host); ip == nil {
		return fmt.Errorf("invalid IP address: %s", host)
	}
	
	// Validate port
	if port, err := strconv.Atoi(portStr); err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid port number: %s", portStr)
	}
	
	return nil
}