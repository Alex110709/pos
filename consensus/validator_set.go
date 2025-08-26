package consensus

import (
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pixelzx/pos/core/types"
)

// ValidatorSet manages the set of validators
type ValidatorSet struct {
	validators map[common.Address]*types.Validator
	sorted     []*types.Validator
	totalPower *big.Int
	proposer   *types.Validator
	mu         sync.RWMutex
}

// NewValidatorSet creates a new validator set
func NewValidatorSet() *ValidatorSet {
	return &ValidatorSet{
		validators: make(map[common.Address]*types.Validator),
		sorted:     make([]*types.Validator, 0),
		totalPower: big.NewInt(0),
	}
}

// Add adds a validator to the set
func (vs *ValidatorSet) Add(validator *types.Validator) {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	// Remove old validator if exists
	if old, exists := vs.validators[validator.Address]; exists {
		vs.totalPower.Sub(vs.totalPower, old.VotingPower)
	}

	// Add new validator
	vs.validators[validator.Address] = validator
	vs.totalPower.Add(vs.totalPower, validator.VotingPower)

	// Update sorted list
	vs.updateSorted()
}

// Remove removes a validator from the set
func (vs *ValidatorSet) Remove(address common.Address) {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	if validator, exists := vs.validators[address]; exists {
		delete(vs.validators, address)
		vs.totalPower.Sub(vs.totalPower, validator.VotingPower)
		vs.updateSorted()
	}
}

// Update updates the entire validator set
func (vs *ValidatorSet) Update(validators []*types.Validator) {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	// Clear existing validators
	vs.validators = make(map[common.Address]*types.Validator)
	vs.totalPower = big.NewInt(0)

	// Add new validators
	for _, validator := range validators {
		vs.validators[validator.Address] = validator
		vs.totalPower.Add(vs.totalPower, validator.VotingPower)
	}

	// Update sorted list
	vs.updateSorted()
}

// GetValidators returns all validators
func (vs *ValidatorSet) GetValidators() []*types.Validator {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	validators := make([]*types.Validator, len(vs.sorted))
	copy(validators, vs.sorted)
	return validators
}

// GetByAddress returns validator by address
func (vs *ValidatorSet) GetByAddress(address common.Address) *types.Validator {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	return vs.validators[address]
}

// IsValidator checks if address is a validator
func (vs *ValidatorSet) IsValidator(address common.Address) bool {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	_, exists := vs.validators[address]
	return exists
}

// GetProposer returns the proposer for a given block number
func (vs *ValidatorSet) GetProposer(blockNumber uint64) *types.Validator {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	if len(vs.sorted) == 0 {
		return nil
	}

	// Round-robin selection based on block number
	index := blockNumber % uint64(len(vs.sorted))
	return vs.sorted[index]
}

// GetTotalVotingPower returns total voting power
func (vs *ValidatorSet) GetTotalVotingPower() *big.Int {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	return new(big.Int).Set(vs.totalPower)
}

// Size returns the number of validators
func (vs *ValidatorSet) Size() int {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	return len(vs.validators)
}

// Hash returns hash of the validator set
func (vs *ValidatorSet) Hash() common.Hash {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	if len(vs.sorted) == 0 {
		return common.Hash{}
	}

	// Create deterministic hash based on sorted validators
	var data []byte
	for _, validator := range vs.sorted {
		data = append(data, validator.Address.Bytes()...)
		data = append(data, validator.PubKey...)
		data = append(data, validator.VotingPower.Bytes()...)
	}

	return crypto.Keccak256Hash(data)
}

// GetAddresses returns addresses of all validators
func (vs *ValidatorSet) GetAddresses() []common.Address {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	addresses := make([]common.Address, len(vs.sorted))
	for i, validator := range vs.sorted {
		addresses[i] = validator.Address
	}
	return addresses
}

// UpdateForNewEpoch updates validator set for new epoch
func (vs *ValidatorSet) UpdateForNewEpoch() {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	// Remove jailed validators
	for addr, validator := range vs.validators {
		if validator.Jailed {
			delete(vs.validators, addr)
			vs.totalPower.Sub(vs.totalPower, validator.VotingPower)
		}
	}

	// Update sorted list
	vs.updateSorted()

	// Limit to maximum validators
	if len(vs.sorted) > int(types.MaxValidators) {
		// Keep top validators by voting power
		vs.sorted = vs.sorted[:types.MaxValidators]
		
		// Update map to match sorted list
		newValidators := make(map[common.Address]*types.Validator)
		newTotalPower := big.NewInt(0)
		
		for _, validator := range vs.sorted {
			newValidators[validator.Address] = validator
			newTotalPower.Add(newTotalPower, validator.VotingPower)
		}
		
		vs.validators = newValidators
		vs.totalPower = newTotalPower
	}
}

// GetTopValidators returns top N validators by voting power
func (vs *ValidatorSet) GetTopValidators(n int) []*types.Validator {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	if n > len(vs.sorted) {
		n = len(vs.sorted)
	}

	top := make([]*types.Validator, n)
	copy(top, vs.sorted[:n])
	return top
}

// ValidateSignatures validates signatures from 2/3+ validators
func (vs *ValidatorSet) ValidateSignatures(signatures map[common.Address][]byte, hash common.Hash) error {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	// Calculate required voting power (2/3 + 1)
	required := new(big.Int).Mul(vs.totalPower, big.NewInt(2))
	required.Div(required, big.NewInt(3))
	required.Add(required, big.NewInt(1))

	// Validate signatures and accumulate voting power
	validPower := big.NewInt(0)
	for addr, signature := range signatures {
		validator := vs.validators[addr]
		if validator == nil {
			continue
		}

		// Verify signature
		pubKey, err := crypto.SigToPub(hash.Bytes(), signature)
		if err != nil {
			continue
		}

		recoveredAddr := crypto.PubkeyToAddress(*pubKey)
		if recoveredAddr != addr {
			continue
		}

		validPower.Add(validPower, validator.VotingPower)
	}

	// Check if we have enough voting power
	if validPower.Cmp(required) < 0 {
		return ErrInsufficientVotingPower
	}

	return nil
}

// updateSorted updates the sorted validator list (must be called with lock held)
func (vs *ValidatorSet) updateSorted() {
	vs.sorted = make([]*types.Validator, 0, len(vs.validators))
	for _, validator := range vs.validators {
		if !validator.Jailed {
			vs.sorted = append(vs.sorted, validator)
		}
	}

	// Sort by voting power (descending), then by address (ascending) for determinism
	sort.Slice(vs.sorted, func(i, j int) bool {
		if vs.sorted[i].VotingPower.Cmp(vs.sorted[j].VotingPower) == 0 {
			return vs.sorted[i].Address.Hex() < vs.sorted[j].Address.Hex()
		}
		return vs.sorted[i].VotingPower.Cmp(vs.sorted[j].VotingPower) > 0
	})
}

// Errors
var (
	ErrInsufficientVotingPower = fmt.Errorf("insufficient voting power")
	ErrValidatorNotFound       = fmt.Errorf("validator not found")
	ErrValidatorJailed         = fmt.Errorf("validator is jailed")
)

// ValidatorSetSnapshot represents a snapshot of validator set at specific height
type ValidatorSetSnapshot struct {
	Height     uint64                           `json:"height"`
	Validators map[common.Address]*types.Validator `json:"validators"`
	TotalPower *big.Int                         `json:"totalPower"`
	Hash       common.Hash                      `json:"hash"`
}

// CreateSnapshot creates a snapshot of current validator set
func (vs *ValidatorSet) CreateSnapshot(height uint64) *ValidatorSetSnapshot {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	// Deep copy validators
	validators := make(map[common.Address]*types.Validator)
	for addr, validator := range vs.validators {
		validatorCopy := *validator
		validators[addr] = &validatorCopy
	}

	return &ValidatorSetSnapshot{
		Height:     height,
		Validators: validators,
		TotalPower: new(big.Int).Set(vs.totalPower),
		Hash:       vs.Hash(),
	}
}

// RestoreSnapshot restores validator set from snapshot
func (vs *ValidatorSet) RestoreSnapshot(snapshot *ValidatorSetSnapshot) {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	vs.validators = snapshot.Validators
	vs.totalPower = new(big.Int).Set(snapshot.TotalPower)
	vs.updateSorted()
}

// GetValidatorByRank returns validator by rank (0-based)
func (vs *ValidatorSet) GetValidatorByRank(rank int) *types.Validator {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	if rank < 0 || rank >= len(vs.sorted) {
		return nil
	}

	return vs.sorted[rank]
}

// GetValidatorRank returns the rank of a validator (0-based)
func (vs *ValidatorSet) GetValidatorRank(address common.Address) int {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	for i, validator := range vs.sorted {
		if validator.Address == address {
			return i
		}
	}

	return -1
}