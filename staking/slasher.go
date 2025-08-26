package staking

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pixelzx/pos/core/types"
)

// Slasher handles slashing operations
type Slasher struct {
	manager *Manager
	
	// Slashing parameters
	doubleSignSlashRate   *big.Int // Slash rate for double signing (in basis points)
	downtimeSlashRate     *big.Int // Slash rate for downtime (in basis points)
	maxDowntimeBlocks     uint64   // Maximum consecutive blocks a validator can miss
	slashingWindow        uint64   // Window for tracking missed blocks
	minSlashAmount        *big.Int // Minimum amount to slash
	
	// Tracking
	missedBlocks          map[common.Address]uint64    // validator -> consecutive missed blocks
	lastSeenBlock         map[common.Address]uint64    // validator -> last block they were seen
	slashingHistory       map[common.Address][]SlashRecord // validator -> slash history
	
	mu sync.RWMutex
}

// SlashRecord represents a slashing record
type SlashRecord struct {
	Height      uint64    `json:"height"`
	Type        SlashType `json:"type"`
	Amount      *big.Int  `json:"amount"`
	Reason      string    `json:"reason"`
	Timestamp   time.Time `json:"timestamp"`
}

// SlashType represents types of slashing
type SlashType uint8

const (
	SlashTypeDoubleSign SlashType = iota
	SlashTypeDowntime
	SlashTypeMaliciousBehavior
)

// SlasherConfig represents slasher configuration
type SlasherConfig struct {
	DoubleSignSlashRate *big.Int // e.g., 500 = 5%
	DowntimeSlashRate   *big.Int // e.g., 100 = 1%
	MaxDowntimeBlocks   uint64   // e.g., 100 consecutive blocks
	SlashingWindow      uint64   // e.g., 10000 blocks window
	MinSlashAmount      *big.Int // Minimum amount to slash
}

// NewSlasher creates a new slasher
func NewSlasher(manager *Manager, config *SlasherConfig) *Slasher {
	return &Slasher{
		manager:               manager,
		doubleSignSlashRate:   config.DoubleSignSlashRate,
		downtimeSlashRate:     config.DowntimeSlashRate,
		maxDowntimeBlocks:     config.MaxDowntimeBlocks,
		slashingWindow:        config.SlashingWindow,
		minSlashAmount:        config.MinSlashAmount,
		missedBlocks:          make(map[common.Address]uint64),
		lastSeenBlock:         make(map[common.Address]uint64),
		slashingHistory:       make(map[common.Address][]SlashRecord),
	}
}

// SlashDoubleSign slashes a validator for double signing
func (s *Slasher) SlashDoubleSign(validator common.Address, height uint64, evidence DoubleSignEvidence) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get validator
	val, err := s.manager.GetValidator(validator)
	if err != nil {
		return fmt.Errorf("validator not found: %w", err)
	}

	// Check if validator is already jailed
	if val.Jailed {
		return fmt.Errorf("validator already jailed")
	}

	// Calculate slash amount
	slashAmount := s.calculateSlashAmount(val.VotingPower, s.doubleSignSlashRate)
	
	// Apply slashing
	if err := s.applySlash(validator, slashAmount, SlashTypeDoubleSign, 
		fmt.Sprintf("Double signing at height %d", height)); err != nil {
		return fmt.Errorf("failed to apply slash: %w", err)
	}

	// Jail the validator permanently for double signing
	if err := s.jailValidator(validator, time.Time{}); err != nil {
		return fmt.Errorf("failed to jail validator: %w", err)
	}

	return nil
}

// SlashDowntime slashes a validator for excessive downtime
func (s *Slasher) SlashDowntime(validator common.Address, height uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get validator
	val, err := s.manager.GetValidator(validator)
	if err != nil {
		return fmt.Errorf("validator not found: %w", err)
	}

	// Check if validator is already jailed
	if val.Jailed {
		return fmt.Errorf("validator already jailed")
	}

	// Check if downtime threshold is met
	missedBlocks := s.missedBlocks[validator]
	if missedBlocks < s.maxDowntimeBlocks {
		return fmt.Errorf("downtime threshold not met: %d < %d", missedBlocks, s.maxDowntimeBlocks)
	}

	// Calculate slash amount
	slashAmount := s.calculateSlashAmount(val.VotingPower, s.downtimeSlashRate)
	
	// Apply slashing
	if err := s.applySlash(validator, slashAmount, SlashTypeDowntime, 
		fmt.Sprintf("Excessive downtime: %d consecutive missed blocks", missedBlocks)); err != nil {
		return fmt.Errorf("failed to apply slash: %w", err)
	}

	// Jail the validator temporarily for downtime
	jailUntil := time.Now().Add(24 * time.Hour) // Jail for 24 hours
	if err := s.jailValidator(validator, jailUntil); err != nil {
		return fmt.Errorf("failed to jail validator: %w", err)
	}

	// Reset missed blocks counter
	s.missedBlocks[validator] = 0

	return nil
}

// TrackValidatorActivity tracks validator activity for downtime detection
func (s *Slasher) TrackValidatorActivity(validator common.Address, height uint64, signed bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if signed {
		// Validator signed, reset missed blocks
		s.missedBlocks[validator] = 0
		s.lastSeenBlock[validator] = height
	} else {
		// Validator missed block, increment counter
		s.missedBlocks[validator]++
	}
}

// CheckDowntime checks if any validators should be slashed for downtime
func (s *Slasher) CheckDowntime(currentHeight uint64) []common.Address {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var toSlash []common.Address
	
	for validator, missedBlocks := range s.missedBlocks {
		if missedBlocks >= s.maxDowntimeBlocks {
			// Check if validator is not already jailed
			val, err := s.manager.GetValidator(validator)
			if err != nil || val.Jailed {
				continue
			}
			
			toSlash = append(toSlash, validator)
		}
	}

	return toSlash
}

// applySlash applies slashing to a validator and their delegators
func (s *Slasher) applySlash(validator common.Address, slashAmount *big.Int, slashType SlashType, reason string) error {
	// Get validator
	val, err := s.manager.GetValidator(validator)
	if err != nil {
		return err
	}

	// Ensure slash amount doesn't exceed voting power
	if slashAmount.Cmp(val.VotingPower) > 0 {
		slashAmount = new(big.Int).Set(val.VotingPower)
	}

	// Ensure minimum slash amount
	if slashAmount.Cmp(s.minSlashAmount) < 0 {
		slashAmount = new(big.Int).Set(s.minSlashAmount)
	}

	// Reduce validator's voting power
	val.VotingPower.Sub(val.VotingPower, slashAmount)

	// Slash delegations proportionally
	if err := s.slashDelegations(validator, slashAmount); err != nil {
		return fmt.Errorf("failed to slash delegations: %w", err)
	}

	// Record slashing
	slashRecord := SlashRecord{
		Height:    0, // Would be set to current block height
		Type:      slashType,
		Amount:    new(big.Int).Set(slashAmount),
		Reason:    reason,
		Timestamp: time.Now(),
	}

	s.slashingHistory[validator] = append(s.slashingHistory[validator], slashRecord)

	// Emit slashing event
	s.manager.events <- Event{
		Type: "ValidatorSlashed",
		Data: map[string]interface{}{
			"validator":   validator,
			"amount":      slashAmount,
			"type":        slashType,
			"reason":      reason,
			"timestamp":   slashRecord.Timestamp,
		},
	}

	return nil
}

// slashDelegations slashes delegations proportionally
func (s *Slasher) slashDelegations(validator common.Address, totalSlashAmount *big.Int) error {
	// Get total voting power before slashing
	val, err := s.manager.GetValidator(validator)
	if err != nil {
		return err
	}

	originalVotingPower := new(big.Int).Add(val.VotingPower, totalSlashAmount)

	// Slash each delegation proportionally
	for delegator, delegations := range s.manager.delegations {
		if delegation, exists := delegations[validator]; exists {
			// Calculate proportional slash amount
			delegationSlash := new(big.Int).Mul(totalSlashAmount, delegation.Amount)
			delegationSlash.Div(delegationSlash, originalVotingPower)

			// Apply slash to delegation
			delegation.Amount.Sub(delegation.Amount, delegationSlash)
			
			// Calculate shares to remove proportionally
			shareSlash := new(big.Int).Mul(delegation.Shares, delegationSlash)
			shareSlash.Div(shareSlash, delegation.Amount.Add(delegation.Amount, delegationSlash))
			delegation.Shares.Sub(delegation.Shares, shareSlash)

			// Remove delegation if amount becomes zero or negative
			if delegation.Amount.Cmp(big.NewInt(0)) <= 0 {
				delete(delegations, validator)
				if len(delegations) == 0 {
					delete(s.manager.delegations, delegator)
				}
			}

			// Emit event
			s.manager.events <- Event{
				Type: "DelegationSlashed",
				Data: map[string]interface{}{
					"delegator": delegator,
					"validator": validator,
					"amount":    delegationSlash,
				},
			}
		}
	}

	return nil
}

// jailValidator jails a validator
func (s *Slasher) jailValidator(validator common.Address, jailUntil time.Time) error {
	val, err := s.manager.GetValidator(validator)
	if err != nil {
		return err
	}

	val.Jailed = true
	val.JailedUntil = jailUntil

	// Emit event
	s.manager.events <- Event{
		Type: "ValidatorJailed",
		Data: map[string]interface{}{
			"validator": validator,
			"jailUntil": jailUntil,
		},
	}

	return nil
}

// UnjailValidator unjails a validator if the jail period has expired
func (s *Slasher) UnjailValidator(validator common.Address) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, err := s.manager.GetValidator(validator)
	if err != nil {
		return err
	}

	if !val.Jailed {
		return fmt.Errorf("validator is not jailed")
	}

	// Check if jail period has expired (for non-permanent jails)
	if !val.JailedUntil.IsZero() && time.Now().Before(val.JailedUntil) {
		return fmt.Errorf("jail period has not expired")
	}

	// For permanent jails (double signing), require governance approval
	if val.JailedUntil.IsZero() {
		return fmt.Errorf("permanent jail requires governance approval")
	}

	val.Jailed = false
	val.JailedUntil = time.Time{}

	// Reset missed blocks counter
	s.missedBlocks[validator] = 0

	// Emit event
	s.manager.events <- Event{
		Type: "ValidatorUnjailed",
		Data: map[string]interface{}{
			"validator": validator,
		},
	}

	return nil
}

// calculateSlashAmount calculates the amount to slash based on voting power and rate
func (s *Slasher) calculateSlashAmount(votingPower, slashRate *big.Int) *big.Int {
	slashAmount := new(big.Int).Mul(votingPower, slashRate)
	slashAmount.Div(slashAmount, big.NewInt(10000)) // Slash rate is in basis points
	return slashAmount
}

// GetSlashingHistory returns slashing history for a validator
func (s *Slasher) GetSlashingHistory(validator common.Address) []SlashRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	history, exists := s.slashingHistory[validator]
	if !exists {
		return []SlashRecord{}
	}

	// Return a copy
	result := make([]SlashRecord, len(history))
	copy(result, history)
	return result
}

// GetMissedBlocks returns the number of consecutive missed blocks for a validator
func (s *Slasher) GetMissedBlocks(validator common.Address) uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.missedBlocks[validator]
}

// DoubleSignEvidence represents evidence of double signing
type DoubleSignEvidence struct {
	ValidatorAddress common.Address `json:"validatorAddress"`
	Height          uint64         `json:"height"`
	Round           int64          `json:"round"`
	Signature1      []byte         `json:"signature1"`
	Signature2      []byte         `json:"signature2"`
	BlockHash1      common.Hash    `json:"blockHash1"`
	BlockHash2      common.Hash    `json:"blockHash2"`
}

// ValidateDoubleSignEvidence validates double sign evidence
func (s *Slasher) ValidateDoubleSignEvidence(evidence DoubleSignEvidence) error {
	// Validate that both signatures are from the same validator
	// Validate that they are for the same height but different blocks
	// Validate signatures
	
	if evidence.Height == 0 {
		return fmt.Errorf("invalid height")
	}

	if evidence.BlockHash1 == evidence.BlockHash2 {
		return fmt.Errorf("block hashes are the same")
	}

	if len(evidence.Signature1) == 0 || len(evidence.Signature2) == 0 {
		return fmt.Errorf("signatures cannot be empty")
	}

	// Additional validation would be implemented here
	return nil
}