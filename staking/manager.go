package staking

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pixelzx/pos/core/types"
)

// Manager manages staking operations
type Manager struct {
	validators     map[common.Address]*types.Validator
	delegations    map[common.Address]map[common.Address]*types.Delegation // delegator -> validator -> delegation
	unbonding      map[common.Address][]*types.UnbondingDelegation          // delegator -> unbonding delegations
	rewards        map[common.Address]*big.Int                              // address -> reward amount
	totalStaked    *big.Int
	
	mu             sync.RWMutex
	blockReward    *big.Int
	inflationRate  *big.Int  // Annual inflation rate in basis points
	commissionRate *big.Int  // Default commission rate in basis points
	
	// Configuration
	minValidatorStake *big.Int
	minDelegatorStake *big.Int
	maxValidators     uint64
	unbondingPeriod   time.Duration
	
	// Events
	events chan Event
}

// Event represents staking events
type Event struct {
	Type string
	Data interface{}
}

// Config represents staking manager configuration
type Config struct {
	MinValidatorStake *big.Int
	MinDelegatorStake *big.Int
	MaxValidators     uint64
	UnbondingPeriod   time.Duration
	BlockReward       *big.Int
	InflationRate     *big.Int      // Annual inflation rate (e.g., 800 = 8%)
	CommissionRate    *big.Int      // Default commission rate (e.g., 1000 = 10%)
}

// NewManager creates a new staking manager
func NewManager(config *Config) *Manager {
	return &Manager{
		validators:        make(map[common.Address]*types.Validator),
		delegations:       make(map[common.Address]map[common.Address]*types.Delegation),
		unbonding:         make(map[common.Address][]*types.UnbondingDelegation),
		rewards:           make(map[common.Address]*big.Int),
		totalStaked:       big.NewInt(0),
		blockReward:       config.BlockReward,
		inflationRate:     config.InflationRate,
		commissionRate:    config.CommissionRate,
		minValidatorStake: config.MinValidatorStake,
		minDelegatorStake: config.MinDelegatorStake,
		maxValidators:     config.MaxValidators,
		unbondingPeriod:   config.UnbondingPeriod,
		events:           make(chan Event, 100),
	}
}

// CreateValidator creates a new validator
func (m *Manager) CreateValidator(
	address common.Address,
	pubKey []byte,
	selfStake *big.Int,
	commission *big.Int,
	details string,
	website string,
) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if validator already exists
	if _, exists := m.validators[address]; exists {
		return ErrValidatorExists
	}

	// Check minimum stake requirement
	if selfStake.Cmp(m.minValidatorStake) < 0 {
		return ErrInsufficientStake
	}

	// Check maximum validators limit
	if uint64(len(m.validators)) >= m.maxValidators {
		return ErrMaxValidatorsReached
	}

	// Create validator
	validator := &types.Validator{
		Address:     address,
		PubKey:      pubKey,
		VotingPower: new(big.Int).Set(selfStake),
		Commission:  commission,
		Jailed:      false,
		Details:     details,
		Website:     website,
	}

	m.validators[address] = validator

	// Create self-delegation
	if m.delegations[address] == nil {
		m.delegations[address] = make(map[common.Address]*types.Delegation)
	}

	delegation := &types.Delegation{
		Delegator: address,
		Validator: address,
		Amount:    new(big.Int).Set(selfStake),
		Shares:    new(big.Int).Set(selfStake), // Initially 1:1 ratio
	}

	m.delegations[address][address] = delegation
	m.totalStaked.Add(m.totalStaked, selfStake)

	// Initialize rewards
	m.rewards[address] = big.NewInt(0)

	// Emit event
	m.events <- Event{
		Type: "ValidatorCreated",
		Data: map[string]interface{}{
			"validator": validator,
			"selfStake": selfStake,
		},
	}

	return nil
}

// Delegate delegates tokens to a validator
func (m *Manager) Delegate(delegator, validator common.Address, amount *big.Int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if validator exists
	val, exists := m.validators[validator]
	if !exists {
		return ErrValidatorNotFound
	}

	// Check if validator is jailed
	if val.Jailed {
		return ErrValidatorJailed
	}

	// Check minimum delegation amount
	if amount.Cmp(m.minDelegatorStake) < 0 {
		return ErrInsufficientStake
	}

	// Initialize delegations map for delegator if needed
	if m.delegations[delegator] == nil {
		m.delegations[delegator] = make(map[common.Address]*types.Delegation)
	}

	// Calculate shares to issue
	shares := m.calculateSharesToIssue(validator, amount)

	// Update or create delegation
	if existingDelegation, exists := m.delegations[delegator][validator]; exists {
		existingDelegation.Amount.Add(existingDelegation.Amount, amount)
		existingDelegation.Shares.Add(existingDelegation.Shares, shares)
	} else {
		delegation := &types.Delegation{
			Delegator: delegator,
			Validator: validator,
			Amount:    new(big.Int).Set(amount),
			Shares:    shares,
		}
		m.delegations[delegator][validator] = delegation
	}

	// Update validator voting power
	val.VotingPower.Add(val.VotingPower, amount)
	m.totalStaked.Add(m.totalStaked, amount)

	// Initialize rewards for delegator if needed
	if m.rewards[delegator] == nil {
		m.rewards[delegator] = big.NewInt(0)
	}

	// Emit event
	m.events <- Event{
		Type: "Delegated",
		Data: map[string]interface{}{
			"delegator": delegator,
			"validator": validator,
			"amount":    amount,
			"shares":    shares,
		},
	}

	return nil
}

// Undelegate starts the unbonding process for delegated tokens
func (m *Manager) Undelegate(delegator, validator common.Address, amount *big.Int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if delegation exists
	delegatorDelegations, exists := m.delegations[delegator]
	if !exists {
		return ErrDelegationNotFound
	}

	delegation, exists := delegatorDelegations[validator]
	if !exists {
		return ErrDelegationNotFound
	}

	// Check if amount is valid
	if amount.Cmp(delegation.Amount) > 0 {
		return ErrInsufficientDelegation
	}

	// Calculate shares to remove
	sharesPercent := new(big.Int).Mul(amount, big.NewInt(10000))
	sharesPercent.Div(sharesPercent, delegation.Amount)
	sharesToRemove := new(big.Int).Mul(delegation.Shares, sharesPercent)
	sharesToRemove.Div(sharesToRemove, big.NewInt(10000))

	// Update delegation
	delegation.Amount.Sub(delegation.Amount, amount)
	delegation.Shares.Sub(delegation.Shares, sharesToRemove)

	// Remove delegation if amount becomes zero
	if delegation.Amount.Cmp(big.NewInt(0)) == 0 {
		delete(delegatorDelegations, validator)
		if len(delegatorDelegations) == 0 {
			delete(m.delegations, delegator)
		}
	}

	// Update validator voting power
	if val, exists := m.validators[validator]; exists {
		val.VotingPower.Sub(val.VotingPower, amount)
	}

	m.totalStaked.Sub(m.totalStaked, amount)

	// Create unbonding delegation
	unbondingDelegation := &types.UnbondingDelegation{
		Delegator:      delegator,
		Validator:      validator,
		Amount:         new(big.Int).Set(amount),
		CompletionTime: time.Now().Add(m.unbondingPeriod),
	}

	m.unbonding[delegator] = append(m.unbonding[delegator], unbondingDelegation)

	// Emit event
	m.events <- Event{
		Type: "Undelegated",
		Data: map[string]interface{}{
			"delegator":      delegator,
			"validator":      validator,
			"amount":         amount,
			"completionTime": unbondingDelegation.CompletionTime,
		},
	}

	return nil
}

// CompleteUnbonding completes unbonding delegations that have passed the unbonding period
func (m *Manager) CompleteUnbonding(delegator common.Address) ([]*types.UnbondingDelegation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	unbondingDelegations, exists := m.unbonding[delegator]
	if !exists {
		return nil, nil
	}

	now := time.Now()
	var completed []*types.UnbondingDelegation
	var remaining []*types.UnbondingDelegation

	for _, unbonding := range unbondingDelegations {
		if now.After(unbonding.CompletionTime) {
			completed = append(completed, unbonding)
			
			// Emit event for completed unbonding
			m.events <- Event{
				Type: "UnbondingCompleted",
				Data: map[string]interface{}{
					"delegator": unbonding.Delegator,
					"validator": unbonding.Validator,
					"amount":    unbonding.Amount,
				},
			}
		} else {
			remaining = append(remaining, unbonding)
		}
	}

	// Update unbonding list
	if len(remaining) == 0 {
		delete(m.unbonding, delegator)
	} else {
		m.unbonding[delegator] = remaining
	}

	return completed, nil
}

// DistributeRewards distributes block rewards to validators and delegators
func (m *Manager) DistributeRewards(blockHeight uint64, blockReward *big.Int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Get active validators
	activeValidators := make([]*types.Validator, 0, len(m.validators))
	for _, validator := range m.validators {
		if !validator.Jailed && validator.VotingPower.Cmp(big.NewInt(0)) > 0 {
			activeValidators = append(activeValidators, validator)
		}
	}

	if len(activeValidators) == 0 {
		return nil // No active validators
	}

	// Calculate total voting power of active validators
	totalVotingPower := big.NewInt(0)
	for _, validator := range activeValidators {
		totalVotingPower.Add(totalVotingPower, validator.VotingPower)
	}

	// Distribute rewards proportionally
	for _, validator := range activeValidators {
		// Calculate validator's share of rewards
		validatorReward := new(big.Int).Mul(blockReward, validator.VotingPower)
		validatorReward.Div(validatorReward, totalVotingPower)

		// Calculate commission
		commission := new(big.Int).Mul(validatorReward, validator.Commission)
		commission.Div(commission, big.NewInt(10000)) // Commission is in basis points

		// Validator gets commission
		if m.rewards[validator.Address] == nil {
			m.rewards[validator.Address] = big.NewInt(0)
		}
		m.rewards[validator.Address].Add(m.rewards[validator.Address], commission)

		// Remaining reward goes to delegators (including validator's self-delegation)
		delegatorReward := new(big.Int).Sub(validatorReward, commission)
		
		// Distribute to delegators proportionally
		if err := m.distributeToDelegators(validator.Address, delegatorReward); err != nil {
			return fmt.Errorf("failed to distribute rewards to delegators: %w", err)
		}
	}

	// Emit event
	m.events <- Event{
		Type: "RewardsDistributed",
		Data: map[string]interface{}{
			"blockHeight": blockHeight,
			"blockReward": blockReward,
			"validators":  len(activeValidators),
		},
	}

	return nil
}

// distributeToDelegators distributes rewards to delegators of a validator
func (m *Manager) distributeToDelegators(validator common.Address, reward *big.Int) error {
	// Get total shares for this validator
	totalShares := big.NewInt(0)
	delegatorShares := make(map[common.Address]*big.Int)

	for delegator, delegations := range m.delegations {
		if delegation, exists := delegations[validator]; exists {
			totalShares.Add(totalShares, delegation.Shares)
			delegatorShares[delegator] = delegation.Shares
		}
	}

	if totalShares.Cmp(big.NewInt(0)) == 0 {
		return nil // No delegators
	}

	// Distribute rewards proportionally to shares
	for delegator, shares := range delegatorShares {
		delegatorReward := new(big.Int).Mul(reward, shares)
		delegatorReward.Div(delegatorReward, totalShares)

		if m.rewards[delegator] == nil {
			m.rewards[delegator] = big.NewInt(0)
		}
		m.rewards[delegator].Add(m.rewards[delegator], delegatorReward)
	}

	return nil
}

// calculateSharesToIssue calculates shares to issue for a delegation amount
func (m *Manager) calculateSharesToIssue(validator common.Address, amount *big.Int) *big.Int {
	// For simplicity, using 1:1 ratio initially
	// In a more sophisticated system, this would account for accumulated rewards
	return new(big.Int).Set(amount)
}

// GetValidator returns validator information
func (m *Manager) GetValidator(address common.Address) (*types.Validator, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	validator, exists := m.validators[address]
	if !exists {
		return nil, ErrValidatorNotFound
	}

	// Return a copy to prevent external modification
	validatorCopy := *validator
	return &validatorCopy, nil
}

// GetValidators returns all validators
func (m *Manager) GetValidators() []*types.Validator {
	m.mu.RLock()
	defer m.mu.RUnlock()

	validators := make([]*types.Validator, 0, len(m.validators))
	for _, validator := range m.validators {
		validatorCopy := *validator
		validators = append(validators, &validatorCopy)
	}

	return validators
}

// GetDelegation returns delegation information
func (m *Manager) GetDelegation(delegator, validator common.Address) (*types.Delegation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	delegations, exists := m.delegations[delegator]
	if !exists {
		return nil, ErrDelegationNotFound
	}

	delegation, exists := delegations[validator]
	if !exists {
		return nil, ErrDelegationNotFound
	}

	// Return a copy
	delegationCopy := *delegation
	return &delegationCopy, nil
}

// GetDelegations returns all delegations for a delegator
func (m *Manager) GetDelegations(delegator common.Address) ([]*types.Delegation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	delegations, exists := m.delegations[delegator]
	if !exists {
		return []*types.Delegation{}, nil
	}

	result := make([]*types.Delegation, 0, len(delegations))
	for _, delegation := range delegations {
		delegationCopy := *delegation
		result = append(result, &delegationCopy)
	}

	return result, nil
}

// GetUnbondingDelegations returns unbonding delegations for a delegator
func (m *Manager) GetUnbondingDelegations(delegator common.Address) ([]*types.UnbondingDelegation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	unbondings, exists := m.unbonding[delegator]
	if !exists {
		return []*types.UnbondingDelegation{}, nil
	}

	result := make([]*types.UnbondingDelegation, len(unbondings))
	for i, unbonding := range unbondings {
		unbondingCopy := *unbonding
		result[i] = &unbondingCopy
	}

	return result, nil
}

// GetRewards returns accumulated rewards for an address
func (m *Manager) GetRewards(address common.Address) *big.Int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	rewards, exists := m.rewards[address]
	if !exists {
		return big.NewInt(0)
	}

	return new(big.Int).Set(rewards)
}

// ClaimRewards claims accumulated rewards
func (m *Manager) ClaimRewards(address common.Address) (*big.Int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	rewards, exists := m.rewards[address]
	if !exists || rewards.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0), nil
	}

	claimedAmount := new(big.Int).Set(rewards)
	m.rewards[address] = big.NewInt(0)

	// Emit event
	m.events <- Event{
		Type: "RewardsClaimed",
		Data: map[string]interface{}{
			"address": address,
			"amount":  claimedAmount,
		},
	}

	return claimedAmount, nil
}

// GetTotalStaked returns total staked amount
func (m *Manager) GetTotalStaked() *big.Int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return new(big.Int).Set(m.totalStaked)
}

// Errors
var (
	ErrValidatorExists        = fmt.Errorf("validator already exists")
	ErrValidatorNotFound      = fmt.Errorf("validator not found")
	ErrValidatorJailed        = fmt.Errorf("validator is jailed")
	ErrInsufficientStake      = fmt.Errorf("insufficient stake amount")
	ErrMaxValidatorsReached   = fmt.Errorf("maximum validators limit reached")
	ErrDelegationNotFound     = fmt.Errorf("delegation not found")
	ErrInsufficientDelegation = fmt.Errorf("insufficient delegation amount")
)