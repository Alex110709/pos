package consensus

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pixelzx/pos/core/types"
)

// Engine implements the PoS consensus engine
type Engine struct {
	chainID      *big.Int
	validators   *ValidatorSet
	proposer     common.Address
	privateKey   *ecdsa.PrivateKey
	blockTime    time.Duration
	epochLength  uint64
	currentEpoch uint64
	
	mu            sync.RWMutex
	isValidating  bool
	stopCh        chan struct{}
	events        chan Event
}

// Event represents consensus events
type Event struct {
	Type string
	Data interface{}
}

// EngineConfig represents consensus engine configuration
type EngineConfig struct {
	ChainID     *big.Int
	PrivateKey  *ecdsa.PrivateKey
	BlockTime   time.Duration
	EpochLength uint64
}

// NewEngine creates a new PoS consensus engine
func NewEngine(config *EngineConfig) *Engine {
	var proposer common.Address
	if config.PrivateKey != nil {
		proposer = crypto.PubkeyToAddress(config.PrivateKey.PublicKey)
	}

	return &Engine{
		chainID:     config.ChainID,
		proposer:    proposer,
		privateKey:  config.PrivateKey,
		blockTime:   config.BlockTime,
		epochLength: config.EpochLength,
		validators:  NewValidatorSet(),
		stopCh:      make(chan struct{}),
		events:      make(chan Event, 100),
	}
}

// Start starts the consensus engine
func (e *Engine) Start(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.isValidating {
		return fmt.Errorf("consensus engine already started")
	}

	e.isValidating = true
	
	// Start block production if we're a validator
	if e.privateKey != nil && e.validators.IsValidator(e.proposer) {
		go e.blockProducer(ctx)
	}

	// Start epoch manager
	go e.epochManager(ctx)

	// Start event handler
	go e.eventHandler(ctx)

	return nil
}

// Stop stops the consensus engine
func (e *Engine) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.isValidating {
		return fmt.Errorf("consensus engine not started")
	}

	close(e.stopCh)
	e.isValidating = false

	return nil
}

// ProposeBlock proposes a new block
func (e *Engine) ProposeBlock(parentHash common.Hash, number uint64, txs []*types.Transaction) (*types.Block, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Check if we're the current proposer
	currentProposer := e.validators.GetProposer(number)
	if currentProposer.Address != e.proposer {
		return nil, fmt.Errorf("not the current proposer")
	}

	// Create block header
	header := &types.Header{
		ParentHash:   parentHash,
		Number:       number,
		Timestamp:    uint64(time.Now().Unix()),
		GasLimit:     types.DefaultGasLimit,
		GasUsed:      0, // Will be calculated after transaction execution
		Proposer:     e.proposer,
		ValidatorSet: e.validators.Hash(),
		Extra:        []byte(fmt.Sprintf("PIXELZX PoS Block #%d", number)),
	}

	// Create block
	block := &types.Block{
		Header:       header,
		Transactions: txs,
		Validators:   e.validators.GetAddresses(),
	}

	// Sign the block
	if err := e.signBlock(block); err != nil {
		return nil, fmt.Errorf("failed to sign block: %w", err)
	}

	return block, nil
}

// ValidateBlock validates a proposed block
func (e *Engine) ValidateBlock(block *types.Block) error {
	// Basic validation
	if block.Header == nil {
		return fmt.Errorf("block header is nil")
	}

	// Validate proposer
	validator := e.validators.GetByAddress(block.Header.Proposer)
	if validator == nil {
		return fmt.Errorf("unknown proposer: %s", block.Header.Proposer.Hex())
	}

	// Validate validator set
	expectedValidatorSet := e.validators.Hash()
	if block.Header.ValidatorSet != expectedValidatorSet {
		return fmt.Errorf("invalid validator set hash")
	}

	// Validate timestamp
	now := time.Now().Unix()
	if int64(block.Header.Timestamp) > now+10 { // Allow 10 seconds future time
		return fmt.Errorf("block timestamp too far in future")
	}

	// Validate signature
	if err := e.validateBlockSignature(block); err != nil {
		return fmt.Errorf("invalid block signature: %w", err)
	}

	return nil
}

// UpdateValidatorSet updates the validator set
func (e *Engine) UpdateValidatorSet(validators []*types.Validator) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.validators.Update(validators)

	// Emit event
	e.events <- Event{
		Type: "ValidatorSetUpdated",
		Data: validators,
	}

	return nil
}

// GetValidators returns current validators
func (e *Engine) GetValidators() []*types.Validator {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.validators.GetValidators()
}

// IsValidator checks if an address is a validator
func (e *Engine) IsValidator(address common.Address) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.validators.IsValidator(address)
}

// GetCurrentProposer returns the current block proposer
func (e *Engine) GetCurrentProposer(blockNumber uint64) common.Address {
	e.mu.RLock()
	defer e.mu.RUnlock()

	proposer := e.validators.GetProposer(blockNumber)
	if proposer != nil {
		return proposer.Address
	}
	return common.Address{}
}

// blockProducer produces blocks when we're the proposer
func (e *Engine) blockProducer(ctx context.Context) {
	ticker := time.NewTicker(e.blockTime)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-e.stopCh:
			return
		case <-ticker.C:
			e.tryProduceBlock()
		}
	}
}

// epochManager manages epoch transitions
func (e *Engine) epochManager(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-e.stopCh:
			return
		case <-time.After(time.Duration(e.epochLength) * e.blockTime):
			e.nextEpoch()
		}
	}
}

// eventHandler handles consensus events
func (e *Engine) eventHandler(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-e.stopCh:
			return
		case event := <-e.events:
			e.handleEvent(event)
		}
	}
}

// tryProduceBlock attempts to produce a block if we're the proposer
func (e *Engine) tryProduceBlock() {
	// This would be connected to the actual blockchain state
	// For now, it's a placeholder
	fmt.Printf("Attempting to produce block at %s\n", time.Now().Format(time.RFC3339))
}

// nextEpoch transitions to the next epoch
func (e *Engine) nextEpoch() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.currentEpoch++
	
	// Update validator set based on stake
	e.validators.UpdateForNewEpoch()

	// Emit event
	e.events <- Event{
		Type: "EpochTransition",
		Data: e.currentEpoch,
	}

	fmt.Printf("Transitioned to epoch %d\n", e.currentEpoch)
}

// handleEvent handles consensus events
func (e *Engine) handleEvent(event Event) {
	switch event.Type {
	case "ValidatorSetUpdated":
		fmt.Printf("Validator set updated: %v\n", event.Data)
	case "EpochTransition":
		fmt.Printf("Epoch transitioned to: %v\n", event.Data)
	default:
		fmt.Printf("Unknown event: %s\n", event.Type)
	}
}

// signBlock signs a block with the engine's private key
func (e *Engine) signBlock(block *types.Block) error {
	if e.privateKey == nil {
		return fmt.Errorf("no private key available for signing")
	}

	// Create hash of block header for signing
	hash := e.blockHash(block.Header)
	
	// Sign the hash
	signature, err := crypto.Sign(hash.Bytes(), e.privateKey)
	if err != nil {
		return err
	}

	block.Header.Signature = signature
	return nil
}

// validateBlockSignature validates block signature
func (e *Engine) validateBlockSignature(block *types.Block) error {
	if len(block.Header.Signature) == 0 {
		return fmt.Errorf("block signature is empty")
	}

	// Get validator's public key
	validator := e.validators.GetByAddress(block.Header.Proposer)
	if validator == nil {
		return fmt.Errorf("unknown validator")
	}

	// Verify signature
	hash := e.blockHash(block.Header)
	pubKey, err := crypto.SigToPub(hash.Bytes(), block.Header.Signature)
	if err != nil {
		return err
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	if recoveredAddr != block.Header.Proposer {
		return fmt.Errorf("signature verification failed")
	}

	return nil
}

// blockHash calculates the hash of a block header
func (e *Engine) blockHash(header *types.Header) common.Hash {
	// Create a copy without signature for hashing
	headerCopy := *header
	headerCopy.Signature = nil
	
	// This is a simplified hash calculation
	// In a real implementation, this would use proper RLP encoding
	data := fmt.Sprintf("%s%d%d%s%s",
		header.ParentHash.Hex(),
		header.Number,
		header.Timestamp,
		header.Proposer.Hex(),
		header.ValidatorSet.Hex(),
	)
	
	return crypto.Keccak256Hash([]byte(data))
}

// GetCurrentEpoch returns the current epoch
func (e *Engine) GetCurrentEpoch() uint64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.currentEpoch
}

// IsValidating returns whether the engine is currently validating
func (e *Engine) IsValidating() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.isValidating
}