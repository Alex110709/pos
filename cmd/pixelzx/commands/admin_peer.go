package commands

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/spf13/cobra"
)

// adminPeerCmd creates the peer command group
func adminPeerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "peer",
		Short: "P2P 네트워크 피어 관리",
		Long: `P2P 네트워크에 연결된 피어를 관리합니다.

연결된 피어 목록을 조회하고, 특정 피어의 상세 정보를 확인하며,
피어와의 연결을 관리할 수 있습니다.`,
	}

	// 하위 명령어 추가
	cmd.AddCommand(
		adminPeerListCmd(),
		adminPeerInfoCmd(),
		adminPeerStatsCmd(),
		adminPeerConnectCmd(),
		adminPeerDisconnectCmd(),
		adminPeerSelfCmd(), // Newly added command
	)

	return cmd
}

// adminPeerListCmd lists connected peers
func adminPeerListCmd() *cobra.Command {
	var (
		verbose bool
		format  string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "연결된 피어 목록 조회",
		Long:  "현재 P2P 네트워크에 연결된 모든 피어의 목록을 표시합니다.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if format == "json" {
				return listPeersJSON(verbose)
			}
			return listPeersTable(cmd, verbose)
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "상세 정보 표시")
	cmd.Flags().StringVar(&format, "format", "table", "출력 형식 (table, json)")

	return cmd
}

// adminPeerInfoCmd shows detailed information about a specific peer
func adminPeerInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info [peer-id]",
		Short: "특정 피어의 상세 정보 확인",
		Long:  "지정된 피어 ID에 대한 상세 정보를 표시합니다.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			peerID := args[0]
			return showPeerInfo(peerID)
		},
	}

	return cmd
}

// adminPeerStatsCmd shows peer connection statistics
func adminPeerStatsCmd() *cobra.Command {
	var duration time.Duration

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "피어 연결 통계 정보",
		Long:  "피어 연결 통계 정보를 표시합니다.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return showPeerStats(duration)
		},
	}

	cmd.Flags().DurationVar(&duration, "duration", time.Minute, "통계 수집 기간")

	return cmd
}

// adminPeerConnectCmd connects to a new peer
func adminPeerConnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect [enode-url]",
		Short: "새로운 피어에 연결",
		Long:  "지정된 ENODE URL을 사용하여 새로운 피어에 연결을 시도합니다.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			enodeURL := args[0]
			return connectToPeer(enodeURL)
		},
	}

	return cmd
}

// adminPeerDisconnectCmd disconnects from a peer
func adminPeerDisconnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disconnect [peer-id]",
		Short: "특정 피어와 연결 해제",
		Long:  "지정된 피어 ID와의 연결을 해제합니다.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			peerID := args[0]
			return disconnectFromPeer(peerID)
		},
	}

	return cmd
}

// adminPeerSelfCmd shows local node enode information
func adminPeerSelfCmd() *cobra.Command {
	var format string
	
	cmd := &cobra.Command{
		Use:   "self",
		Short: "로컬 노드의 enode 정보 조회",
		Long:  "현재 PIXELZX 노드의 enode URL과 관련 정보를 표시합니다.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return showLocalEnode(format)
		},
	}
	
	cmd.Flags().StringVar(&format, "format", "text", "출력 형식 (text, json)")
	
	return cmd
}

// showLocalEnode displays local node enode information
func showLocalEnode(format string) error {
	// TODO: In a real implementation, this would fetch actual enode data from the network manager
	// For now, we'll use mock data for demonstration purposes
	
	// Create a mock enode for demonstration purposes
	// In a real implementation, we would call:
	// node, err := getNetworkManager().GetLocalEnode()
	// if err != nil {
	//     return fmt.Errorf("enode 정보 조회 실패: %w", err)
	// }
	
	if format == "json" {
		return showLocalEnodeJSON()
	}
	
	return showLocalEnodeText()
}

// showLocalEnodeText displays enode information in text format
func showLocalEnodeText() error {
	fmt.Printf("로컬 노드 enode 정보:\n")
	fmt.Printf("========================\n")
	// For demonstration, we'll use a mock enode URL
	fmt.Printf("enode URL: enode://a1b2c3d4e5f67890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef12345@192.168.1.100:30303\n")
	fmt.Printf("Node ID: a1b2c3d4e5f67890abcdef1234567890abcdef1234567890abcdef1234567890abcdef\n")
	fmt.Printf("IP 주소: 192.168.1.100\n")
	fmt.Printf("TCP 포트: 30303\n")
	fmt.Printf("UDP 포트: 30303\n")
	
	return nil
}

// showLocalEnodeJSON displays enode information in JSON format
func showLocalEnodeJSON() error {
	info := struct {
		EnodeURL string `json:"enode_url"`
		NodeID   string `json:"node_id"`
		IP       string `json:"ip"`
		TCPPort  int    `json:"tcp_port"`
		UDPPort  int    `json:"udp_port"`
	}{
		EnodeURL: "enode://a1b2c3d4e5f67890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef12345@192.168.1.100:30303",
		NodeID:   "a1b2c3d4e5f67890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		IP:       "192.168.1.100",
		TCPPort:  30303,
		UDPPort:  30303,
	}
	
	output, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON 직렬화 실패: %w", err)
	}
	
	fmt.Println(string(output))
	return nil
}

// PeerInfo represents information about a peer
type PeerInfo struct {
	ID          string    `json:"id"`
	Enode       string    `json:"enode"`
	IPAddress   string    `json:"ipAddress"`
	Port        int       `json:"port"`
	Direction   string    `json:"direction"`
	Latency     string    `json:"latency"`
	LastSeen    time.Time `json:"lastSeen"`
	ConnectedAt time.Time `json:"connectedAt"`
	Capabilities []string `json:"capabilities"`
}

// PeerStats represents peer connection statistics
type PeerStats struct {
	TotalPeers     int       `json:"totalPeers"`
	InboundPeers   int       `json:"inboundPeers"`
	OutboundPeers  int       `json:"outboundPeers"`
	AvgLatency     string    `json:"avgLatency"`
	BytesReceived  int64     `json:"bytesReceived"`
	BytesSent      int64     `json:"bytesSent"`
	StartTime      time.Time `json:"startTime"`
}

// listPeersTable displays peers in a table format
func listPeersTable(cmd *cobra.Command, verbose bool) error {
	// In a real implementation, this would fetch actual peer data from the network manager
	// For now, we'll use mock data for demonstration
	
	peers := getMockPeers()
	
	if verbose {
		// Verbose table format
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tIP 주소\t포트\t방향\t지연시간\t마지막 통신\t연결 시간")
		fmt.Fprintln(w, "────────────────────────────────────────────────────────────────────────────────")
		
		for _, peer := range peers {
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\t%s\t%s\n",
				peer.ID[:16]+"...",
				peer.IPAddress,
				peer.Port,
				peer.Direction,
				peer.Latency,
				peer.LastSeen.Format("15:04:05"),
				peer.ConnectedAt.Format("15:04:05"))
		}
		
		return w.Flush()
	} else {
		// Simple table format
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tIP 주소\t포트\t방향\t지연시간")
		fmt.Fprintln(w, "────────────────────────────────────────────────────────")
		
		for _, peer := range peers {
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n",
				peer.ID[:16]+"...",
				peer.IPAddress,
				peer.Port,
				peer.Direction,
				peer.Latency)
		}
		
		return w.Flush()
	}
}

// listPeersJSON displays peers in JSON format
func listPeersJSON(verbose bool) error {
	// In a real implementation, this would fetch actual peer data from the network manager
	peers := getMockPeers()
	
	output, err := json.MarshalIndent(peers, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal peers to JSON: %w", err)
	}
	
	fmt.Println(string(output))
	return nil
}

// showPeerInfo displays detailed information about a specific peer
func showPeerInfo(peerID string) error {
	// In a real implementation, this would fetch actual peer data from the network manager
	// For now, we'll use mock data for demonstration
	
	peer := getMockPeerByID(peerID)
	if peer == nil {
		return fmt.Errorf("peer not found: %s", peerID)
	}
	
	fmt.Printf("피어 ID: %s\n", peer.ID)
	fmt.Printf("ENODE: %s\n", peer.Enode)
	fmt.Printf("IP 주소: %s\n", peer.IPAddress)
	fmt.Printf("포트: %d\n", peer.Port)
	fmt.Printf("연결 방향: %s\n", peer.Direction)
	fmt.Printf("지연 시간: %s\n", peer.Latency)
	fmt.Printf("마지막 통신: %s\n", peer.LastSeen.Format("2006-01-02 15:04:05 UTC"))
	fmt.Printf("연결 시간: %s\n", peer.ConnectedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Printf("지원 프로토콜: %s\n", strings.Join(peer.Capabilities, ", "))
	
	return nil
}

// showPeerStats displays peer connection statistics
func showPeerStats(duration time.Duration) error {
	// In a real implementation, this would fetch actual statistics from the network manager
	// For now, we'll use mock data for demonstration
	
	stats := getMockPeerStats()
	
	fmt.Printf("총 피어 수: %d개\n", stats.TotalPeers)
	fmt.Printf("수신 연결: %d개\n", stats.InboundPeers)
	fmt.Printf("송신 연결: %d개\n", stats.OutboundPeers)
	fmt.Printf("평균 지연 시간: %s\n", stats.AvgLatency)
	fmt.Printf("수신 데이터: %.1f GB\n", float64(stats.BytesReceived)/1000000000)
	fmt.Printf("송신 데이터: %.1f GB\n", float64(stats.BytesSent)/1000000000)
	
	return nil
}

// connectToPeer attempts to connect to a new peer
func connectToPeer(enodeURL string) error {
	// First, perform detailed validation of the enode URL
	if err := validateEnodeURL(enodeURL); err != nil {
		return err
	}

	// In a real implementation, this would use the network manager to connect to a peer
	// For now, we'll just validate the enode URL and show a message
	
	_, err := enode.Parse(enode.ValidSchemes, enodeURL)
	if err != nil {
		return fmt.Errorf("invalid enode URL: %w", err)
	}
	
	fmt.Printf("피어에 연결 시도 중: %s\n", enodeURL)
	fmt.Printf("연결이 성공적으로 시작되었습니다.\n")
	
	return nil
}

// validateEnodeURL performs detailed validation of an enode URL
func validateEnodeURL(enodeURL string) error {
	// Check if the URL starts with "enode://"
	if !strings.HasPrefix(enodeURL, "enode://") {
		return fmt.Errorf("enode URL must start with 'enode://'")
	}

	// Remove the "enode://" prefix
	urlWithoutPrefix := strings.TrimPrefix(enodeURL, "enode://")

	// Split the URL into public key and address parts
	parts := strings.Split(urlWithoutPrefix, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid enode URL format: expected 'enode://<public-key>@<ip>:<port>'")
	}

	publicKey := parts[0]
	address := parts[1]

	// Validate the public key
	if err := validatePublicKey(publicKey); err != nil {
		return err
	}

	// Handle IPv6 addresses which are enclosed in brackets
	var ip, port string
	if strings.HasPrefix(address, "[") {
		// IPv6 format: [ipv6]:port
		endBracket := strings.Index(address, "]")
		if endBracket == -1 {
			return fmt.Errorf("invalid IPv6 address format: missing closing bracket")
		}
		
		ip = address[1:endBracket]
		if len(address) <= endBracket+2 || address[endBracket+1] != ':' {
			return fmt.Errorf("invalid address format: expected '[ipv6]:port'")
		}
		
		port = address[endBracket+2:]
	} else {
		// IPv4 format: ip:port
		addrParts := strings.Split(address, ":")
		if len(addrParts) != 2 {
			return fmt.Errorf("invalid address format: expected '<ip>:<port>'")
		}
		
		ip = addrParts[0]
		port = addrParts[1]
	}

	// Validate IP address
	if err := validateIP(ip); err != nil {
		return err
	}

	// Validate port
	if err := validatePort(port); err != nil {
		return err
	}

	return nil
}

// validatePublicKey validates the public key part of an enode URL
func validatePublicKey(publicKey string) error {
	// Check if the public key is a valid hex string
	for _, r := range publicKey {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return fmt.Errorf("공개 키는 16진수 문자열이어야 합니다")
		}
	}

	// Check if the public key length is exactly 128 characters (64 bytes)
	if len(publicKey) != 128 {
		return fmt.Errorf("공개 키는 128자 길이의 16진수 문자열이어야 합니다 (현재 길이: %d)", len(publicKey))
	}

	return nil
}

// validateIP validates the IP address part of an enode URL
func validateIP(ip string) error {
	// Try to parse the IP address
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("IP 주소 형식이 올바르지 않습니다: %s", ip)
	}

	return nil
}

// validatePort validates the port part of an enode URL
func validatePort(portStr string) error {
	// Convert string to integer
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("포트 번호는 1-65535 범위의 숫자여야 합니다: %s", portStr)
	}

	// Check if port is in valid range
	if port < 1 || port > 65535 {
		return fmt.Errorf("포트 번호는 1-65535 범위의 숫자여야 합니다: %s", portStr)
	}

	return nil
}

// disconnectFromPeer disconnects from a peer
func disconnectFromPeer(peerID string) error {
	// In a real implementation, this would use the network manager to disconnect from a peer
	// For now, we'll just show a message
	
	fmt.Printf("피어와의 연결 해제 중: %s\n", peerID)
	fmt.Printf("연결이 성공적으로 해제되었습니다.\n")
	
	return nil
}

// Mock data functions for demonstration purposes
func getMockPeers() []PeerInfo {
	return []PeerInfo{
		{
			ID:          "16Uiu2HAm9Q8R7nxZ8GWqK9vK1b2X3h4mL5nF8dY3...",
			Enode:       "enode://9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7@192.168.1.100:30303",
			IPAddress:   "192.168.1.100",
			Port:        30303,
			Direction:   "수신",
			Latency:     "45ms",
			LastSeen:    time.Now().Add(-10 * time.Second),
			ConnectedAt: time.Now().Add(-2 * time.Hour),
			Capabilities: []string{"pixelzx/1.0"},
		},
		{
			ID:          "16Uiu2HAm5P7M6nxY9FVqH8vJ1a2W2g3hK4mE7cX2...",
			Enode:       "enode://1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c@203.123.45.67:30303",
			IPAddress:   "203.123.45.67",
			Port:        30303,
			Direction:   "송신",
			Latency:     "120ms",
			LastSeen:    time.Now().Add(-30 * time.Second),
			ConnectedAt: time.Now().Add(-1 * time.Hour),
			Capabilities: []string{"pixelzx/1.0", "eth/66"},
		},
		{
			ID:          "16Uiu2HAm8N6L5mxX8EVqG7vI0z1V1f2gJ3lD6bW1...",
			Enode:       "enode://2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d@151.101.1.140:30303",
			IPAddress:   "151.101.1.140",
			Port:        30303,
			Direction:   "수신",
			Latency:     "89ms",
			LastSeen:    time.Now().Add(-5 * time.Second),
			ConnectedAt: time.Now().Add(-3 * time.Hour),
			Capabilities: []string{"pixelzx/1.0"},
		},
	}
}

func getMockPeerByID(peerID string) *PeerInfo {
	peers := getMockPeers()
	for _, peer := range peers {
		if strings.HasPrefix(peer.ID, peerID) || peer.ID == peerID {
			return &peer
		}
	}
	return nil
}

func getMockPeerStats() PeerStats {
	return PeerStats{
		TotalPeers:     24,
		InboundPeers:   12,
		OutboundPeers:  12,
		AvgLatency:     "122ms",
		BytesReceived:  2300000000,
		BytesSent:      1800000000,
		StartTime:      time.Now().Add(-24 * time.Hour),
	}
}