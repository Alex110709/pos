package network

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/pixelzx/pos/core/types"
)

// Manager handles P2P networking for PIXELZX
type Manager struct {
	server     *p2p.Server
	nodeDB     *enode.DB
	config     *Config
	peers      map[enode.ID]*Peer
	protocols  []p2p.Protocol
	
	// Channels for communication
	blockCh    chan *types.Block
	txCh       chan *types.Transaction
	
	mu         sync.RWMutex
	running    bool
	stopCh     chan struct{}
}

// Config represents network configuration
type Config struct {
	Name         string
	Port         int
	MaxPeers     int
	BootNodes    []string
	PrivateKey   []byte
	NoDiscovery  bool
	ListenAddr   string
}

// Peer represents a connected peer
type Peer struct {
	ID       enode.ID
	Node     *enode.Node
	Conn     *p2p.Peer
	Version  uint
	Network  uint64
	Head     common.Hash
	Height   uint64
	
	// Communication channels
	blockCh  chan *types.Block
	txCh     chan *types.Transaction
	
	mu       sync.RWMutex
	lastSeen time.Time
}

// NewManager creates a new network manager
func NewManager(config *Config) (*Manager, error) {
	// Create node database for peer discovery
	nodeDB, err := enode.OpenDB("")
	if err != nil {
		return nil, fmt.Errorf("failed to open node database: %w", err)
	}

	nm := &Manager{
		config:    config,
		nodeDB:    nodeDB,
		peers:     make(map[enode.ID]*Peer),
		blockCh:   make(chan *types.Block, 100),
		txCh:      make(chan *types.Transaction, 1000),
		stopCh:    make(chan struct{}),
	}

	// Setup protocols
	nm.setupProtocols()

	// Create P2P server
	if err := nm.createServer(); err != nil {
		return nil, fmt.Errorf("failed to create P2P server: %w", err)
	}

	return nm, nil
}

// Start starts the network manager
func (nm *Manager) Start(ctx context.Context) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	if nm.running {
		return fmt.Errorf("network manager already running")
	}

	// Start P2P server
	if err := nm.server.Start(); err != nil {
		return fmt.Errorf("failed to start P2P server: %w", err)
	}

	nm.running = true

	// Start goroutines
	go nm.messageHandler(ctx)
	go nm.peerManager(ctx)
	go nm.blockBroadcaster(ctx)
	go nm.txBroadcaster(ctx)

	return nil
}

// Stop stops the network manager
func (nm *Manager) Stop() error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	if !nm.running {
		return nil
	}

	// Stop P2P server
	nm.server.Stop()

	// Stop goroutines
	close(nm.stopCh)
	
	nm.running = false
	return nil
}

// BroadcastBlock broadcasts a block to all peers
func (nm *Manager) BroadcastBlock(block *types.Block) error {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	if !nm.running {
		return fmt.Errorf("network manager not running")
	}

	// Send to block channel for broadcasting
	select {
	case nm.blockCh <- block:
		return nil
	default:
		return fmt.Errorf("block channel full")
	}
}

// BroadcastTransaction broadcasts a transaction to all peers
func (nm *Manager) BroadcastTransaction(tx *types.Transaction) error {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	if !nm.running {
		return fmt.Errorf("network manager not running")
	}

	// Send to transaction channel for broadcasting
	select {
	case nm.txCh <- tx:
		return nil
	default:
		return fmt.Errorf("transaction channel full")
	}
}

// GetPeers returns connected peers
func (nm *Manager) GetPeers() []*Peer {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	peers := make([]*Peer, 0, len(nm.peers))
	for _, peer := range nm.peers {
		peers = append(peers, peer)
	}

	return peers
}

// GetPeerCount returns the number of connected peers
func (nm *Manager) GetPeerCount() int {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	return len(nm.peers)
}

// setupProtocols sets up the network protocols
func (nm *Manager) setupProtocols() {
	protocol := p2p.Protocol{
		Name:    "pixelzx",
		Version: 1,
		Length:  10, // Number of message types
		Run:     nm.handlePeer,
	}

	nm.protocols = []p2p.Protocol{protocol}
}

// createServer creates the P2P server
func (nm *Manager) createServer() error {
	// Parse listen address
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", nm.config.ListenAddr, nm.config.Port))
	if err != nil {
		return fmt.Errorf("invalid listen address: %w", err)
	}

	// Create server configuration
	config := &p2p.Config{
		Name:       nm.config.Name,
		MaxPeers:   nm.config.MaxPeers,
		Protocols:  nm.protocols,
		ListenAddr: addr.String(),
		NodeDatabase: nm.nodeDB,
	}

	// Add boot nodes
	if len(nm.config.BootNodes) > 0 {
		bootNodes := make([]*enode.Node, 0, len(nm.config.BootNodes))
		for _, bootNode := range nm.config.BootNodes {
			if node, err := enode.Parse(enode.ValidSchemes, bootNode); err == nil {
				bootNodes = append(bootNodes, node)
			}
		}
		config.BootstrapNodes = bootNodes
	}

	// Disable discovery if requested
	if nm.config.NoDiscovery {
		config.NoDiscovery = true
	}

	// Create server
	nm.server = &p2p.Server{Config: *config}

	return nil
}

// handlePeer handles a new peer connection
func (nm *Manager) handlePeer(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
	// Create peer instance
	p := &Peer{
		ID:       peer.ID(),
		Node:     peer.Node(),
		Conn:     peer,
		blockCh:  make(chan *types.Block, 10),
		txCh:     make(chan *types.Transaction, 100),
		lastSeen: time.Now(),
	}

	// Add peer to manager
	nm.addPeer(p)
	defer nm.removePeer(p.ID)

	// Handshake
	if err := nm.handshake(p, rw); err != nil {
		return fmt.Errorf("handshake failed: %w", err)
	}

	// Handle messages
	return nm.handlePeerMessages(p, rw)
}

// addPeer adds a peer to the manager
func (nm *Manager) addPeer(peer *Peer) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	nm.peers[peer.ID] = peer
}

// removePeer removes a peer from the manager
func (nm *Manager) removePeer(id enode.ID) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	delete(nm.peers, id)
}

// handshake performs handshake with a peer
func (nm *Manager) handshake(peer *Peer, rw p2p.MsgReadWriter) error {
	// Send handshake message
	handshake := &HandshakeMsg{
		Version: 1,
		Network: 1337, // PIXELZX chain ID
		Genesis: common.Hash{}, // Genesis block hash
		Head:    common.Hash{}, // Current head block hash
		Height:  0,             // Current block height
	}

	if err := p2p.Send(rw, HandshakeMsgCode, handshake); err != nil {
		return fmt.Errorf("failed to send handshake: %w", err)
	}

	// Receive handshake response
	msg, err := rw.ReadMsg()
	if err != nil {
		return fmt.Errorf("failed to read handshake response: %w", err)
	}
	defer msg.Discard()

	if msg.Code != HandshakeMsgCode {
		return fmt.Errorf("expected handshake message, got %d", msg.Code)
	}

	var response HandshakeMsg
	if err := msg.Decode(&response); err != nil {
		return fmt.Errorf("failed to decode handshake response: %w", err)
	}

	// Validate handshake
	if response.Network != 1337 {
		return fmt.Errorf("different network: %d", response.Network)
	}

	// Update peer info
	peer.Version = response.Version
	peer.Network = response.Network
	peer.Head = response.Head
	peer.Height = response.Height

	return nil
}

// handlePeerMessages handles messages from a peer
func (nm *Manager) handlePeerMessages(peer *Peer, rw p2p.MsgReadWriter) error {
	for {
		// Read message
		msg, err := rw.ReadMsg()
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		// Handle message based on type
		switch msg.Code {
		case BlockMsgCode:
			if err := nm.handleBlockMessage(peer, msg); err != nil {
				msg.Discard()
				continue
			}

		case TransactionMsgCode:
			if err := nm.handleTransactionMessage(peer, msg); err != nil {
				msg.Discard()
				continue
			}

		case RequestBlocksMsgCode:
			if err := nm.handleRequestBlocksMessage(peer, rw, msg); err != nil {
				msg.Discard()
				continue
			}

		default:
			msg.Discard()
		}
	}
}

// Message handling functions
func (nm *Manager) handleBlockMessage(peer *Peer, msg p2p.Msg) error {
	defer msg.Discard()

	var block types.Block
	if err := msg.Decode(&block); err != nil {
		return fmt.Errorf("failed to decode block: %w", err)
	}

	// Update peer's head
	peer.mu.Lock()
	peer.Head = block.Header.Hash()
	peer.Height = block.Header.Number
	peer.lastSeen = time.Now()
	peer.mu.Unlock()

	// TODO: Validate and process block
	// For now, just log
	fmt.Printf("Received block #%d from peer %s\n", block.Header.Number, peer.ID.TerminalString())

	return nil
}

func (nm *Manager) handleTransactionMessage(peer *Peer, msg p2p.Msg) error {
	defer msg.Discard()

	var tx types.Transaction
	if err := msg.Decode(&tx); err != nil {
		return fmt.Errorf("failed to decode transaction: %w", err)
	}

	// Update peer's last seen
	peer.mu.Lock()
	peer.lastSeen = time.Now()
	peer.mu.Unlock()

	// TODO: Validate and process transaction
	// For now, just log
	fmt.Printf("Received transaction %s from peer %s\n", tx.Hash().Hex(), peer.ID.TerminalString())

	return nil
}

func (nm *Manager) handleRequestBlocksMessage(peer *Peer, rw p2p.MsgReadWriter, msg p2p.Msg) error {
	defer msg.Discard()

	var req RequestBlocksMsg
	if err := msg.Decode(&req); err != nil {
		return fmt.Errorf("failed to decode request: %w", err)
	}

	// TODO: Fetch and send requested blocks
	// For now, just acknowledge
	fmt.Printf("Received block request from %d to %d from peer %s\n", 
		req.From, req.To, peer.ID.TerminalString())

	return nil
}

// Background goroutines
func (nm *Manager) messageHandler(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-nm.stopCh:
			return
		default:
			// Handle incoming messages
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (nm *Manager) peerManager(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-nm.stopCh:
			return
		case <-ticker.C:
			nm.cleanupPeers()
		}
	}
}

func (nm *Manager) blockBroadcaster(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-nm.stopCh:
			return
		case block := <-nm.blockCh:
			nm.broadcastBlock(block)
		}
	}
}

func (nm *Manager) txBroadcaster(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-nm.stopCh:
			return
		case tx := <-nm.txCh:
			nm.broadcastTransaction(tx)
		}
	}
}

// Helper functions
func (nm *Manager) cleanupPeers() {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	now := time.Now()
	for id, peer := range nm.peers {
		peer.mu.RLock()
		if now.Sub(peer.lastSeen) > 5*time.Minute {
			delete(nm.peers, id)
		}
		peer.mu.RUnlock()
	}
}

func (nm *Manager) broadcastBlock(block *types.Block) {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	for _, peer := range nm.peers {
		select {
		case peer.blockCh <- block:
		default:
			// Channel full, skip this peer
		}
	}
}

func (nm *Manager) broadcastTransaction(tx *types.Transaction) {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	for _, peer := range nm.peers {
		select {
		case peer.txCh <- tx:
		default:
			// Channel full, skip this peer
		}
	}
}

// Message types and codes
const (
	HandshakeMsgCode      = 0x00
	BlockMsgCode          = 0x01
	TransactionMsgCode    = 0x02
	RequestBlocksMsgCode  = 0x03
	ResponseBlocksMsgCode = 0x04
)

// Message structures
type HandshakeMsg struct {
	Version uint64      `json:"version"`
	Network uint64      `json:"network"`
	Genesis common.Hash `json:"genesis"`
	Head    common.Hash `json:"head"`
	Height  uint64      `json:"height"`
}

type RequestBlocksMsg struct {
	From uint64 `json:"from"`
	To   uint64 `json:"to"`
}

type ResponseBlocksMsg struct {
	Blocks []*types.Block `json:"blocks"`
}