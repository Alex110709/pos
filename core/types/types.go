package types

import (
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// Hash represents a 32-byte hash
type Hash = common.Hash

// Address represents a 20-byte address
type Address = common.Address

// BlockNumber represents a block number
type BlockNumber = *big.Int

// ChainID represents chain identifier
type ChainID = *big.Int

// Transaction represents PIXELZX transaction
type Transaction = types.Transaction

// Block represents PIXELZX block
type Block struct {
	Header       *Header        `json:"header"`
	Transactions []*Transaction `json:"transactions"`
	Validators   []Address      `json:"validators"`
}

// Header represents block header
type Header struct {
	ParentHash   Hash      `json:"parentHash"`
	Number       uint64    `json:"number"`
	Timestamp    uint64    `json:"timestamp"`
	GasLimit     uint64    `json:"gasLimit"`
	GasUsed      uint64    `json:"gasUsed"`
	StateRoot    Hash      `json:"stateRoot"`
	TxRoot       Hash      `json:"transactionRoot"`
	ReceiptRoot  Hash      `json:"receiptRoot"`
	ValidatorSet Hash      `json:"validatorSet"`
	Proposer     Address   `json:"proposer"`
	Signature    []byte    `json:"signature"`
	Extra        []byte    `json:"extra"`
}

// Validator represents a validator in the network
type Validator struct {
	Address     Address    `json:"address"`
	PubKey      []byte     `json:"pubkey"`
	VotingPower *big.Int   `json:"votingPower"`
	Commission  *big.Int   `json:"commission"`  // Commission rate in basis points (10000 = 100%)
	Jailed      bool       `json:"jailed"`
	JailedUntil time.Time  `json:"jailedUntil"`
	Details     string     `json:"details"`
	Website     string     `json:"website"`
}

// Delegation represents stake delegation
type Delegation struct {
	Delegator Address  `json:"delegator"`
	Validator Address  `json:"validator"`
	Amount    *big.Int `json:"amount"`
	Shares    *big.Int `json:"shares"`
}

// UnbondingDelegation represents an unbonding delegation
type UnbondingDelegation struct {
	Delegator     Address   `json:"delegator"`
	Validator     Address   `json:"validator"`
	Amount        *big.Int  `json:"amount"`
	CompletionTime time.Time `json:"completionTime"`
}

// Proposal represents a governance proposal
type Proposal struct {
	ID               uint64            `json:"id"`
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	Proposer         Address           `json:"proposer"`
	Status           ProposalStatus    `json:"status"`
	SubmitTime       time.Time         `json:"submitTime"`
	VotingStartTime  time.Time         `json:"votingStartTime"`
	VotingEndTime    time.Time         `json:"votingEndTime"`
	TotalDeposit     *big.Int          `json:"totalDeposit"`
	FinalTallyResult TallyResult       `json:"finalTallyResult"`
}

// ProposalStatus represents proposal status
type ProposalStatus uint8

const (
	StatusNil ProposalStatus = iota
	StatusDepositPeriod
	StatusVotingPeriod
	StatusPassed
	StatusRejected
	StatusFailed
)

// Vote represents a governance vote
type Vote struct {
	ProposalID uint64      `json:"proposalId"`
	Voter      Address     `json:"voter"`
	Option     VoteOption  `json:"option"`
	Weight     *big.Int    `json:"weight"`
}

// VoteOption represents vote option
type VoteOption uint8

const (
	OptionEmpty VoteOption = iota
	OptionYes
	OptionAbstain
	OptionNo
	OptionNoWithVeto
)

// TallyResult represents vote tally result
type TallyResult struct {
	Yes        *big.Int `json:"yes"`
	Abstain    *big.Int `json:"abstain"`
	No         *big.Int `json:"no"`
	NoWithVeto *big.Int `json:"noWithVeto"`
}

// Account represents PIXELZX account
type Account struct {
	Address common.Address `json:"address"`
	Balance *big.Int       `json:"balance"`
	Nonce   uint64         `json:"nonce"`
}

// ValidatorSet represents a set of validators
type ValidatorSet struct {
	Validators []*Validator `json:"validators"`
	Proposer   *Validator   `json:"proposer"`
}

// Constants for PIXELZX network
const (
	// Token constants
	TokenName     = "PIXELZX"
	TokenSymbol   = "PXZ"
	TokenDecimals = 18
	TotalSupply   = "1000000000000000000000000000" // 1 billion PXZ in wei

	// Network constants
	DefaultChainID = 1337
	BlockTime      = 3 * time.Second
	EpochLength    = 200

	// Staking constants
	MinValidatorStake = "100000000000000000000000" // 100,000 PXZ in wei
	MinDelegatorStake = "1000000000000000000"      // 1 PXZ in wei
	MaxValidators     = 125
	UnbondingPeriod   = 21 * 24 * time.Hour        // 21 days

	// Gas constants
	DefaultGasLimit = 30000000
	DefaultGasPrice = 20000000000 // 20 Gwei

	// Governance constants
	MinDeposit         = "100000000000000000000000" // 100,000 PXZ in wei
	DepositPeriod      = 7 * 24 * time.Hour         // 7 days
	VotingPeriod       = 14 * 24 * time.Hour        // 14 days
	Quorum             = "200000000000000000"        // 20% (0.2)
	Threshold          = "500000000000000000"        // 50% (0.5)
	VetoThreshold      = "334000000000000000"        // 33.4% (0.334)
)

// KeyPair represents a cryptographic key pair
type KeyPair struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    common.Address
}

// NewKeyPair generates a new key pair
func NewKeyPair() (*KeyPair, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	publicKey := &privateKey.PublicKey
	address := crypto.PubkeyToAddress(*publicKey)

	return &KeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}, nil
}

// Genesis represents genesis configuration
type Genesis struct {
	ChainID    uint64                 `json:"chainId"`
	Timestamp  uint64                 `json:"timestamp"`
	Alloc      map[string]*big.Int    `json:"alloc"`
	Validators []GenesisValidator     `json:"validators"`
	Config     *ChainConfig           `json:"config"`
}

// GenesisValidator represents validator in genesis
type GenesisValidator struct {
	Address Address  `json:"address"`
	PubKey  []byte   `json:"pubkey"`
	Power   *big.Int `json:"power"`
}

// ChainConfig represents chain configuration
type ChainConfig struct {
	ChainID     uint64     `json:"chainId"`
	POS         *POSConfig `json:"pos"`
	EIP150Block uint64     `json:"eip150Block"`
	EIP155Block uint64     `json:"eip155Block"`
	EIP158Block uint64     `json:"eip158Block"`
}

// POSConfig represents Proof of Stake configuration
type POSConfig struct {
	Period          uint64   `json:"period"`
	Epoch           uint64   `json:"epoch"`
	MinStake        *big.Int `json:"minStake"`
	MaxValidators   uint64   `json:"maxValidators"`
	UnbondingPeriod uint64   `json:"unbondingPeriod"`
}

// Event represents blockchain event
type Event struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	BlockHash Hash        `json:"blockHash"`
	TxHash    Hash        `json:"txHash"`
	TxIndex   uint        `json:"txIndex"`
	LogIndex  uint        `json:"logIndex"`
}

// StakingEvent represents staking-related events
type StakingEvent struct {
	Type      StakingEventType `json:"type"`
	Validator Address          `json:"validator"`
	Delegator Address          `json:"delegator"`
	Amount    *big.Int         `json:"amount"`
}

// StakingEventType represents staking event types
type StakingEventType uint8

const (
	EventDelegate StakingEventType = iota
	EventUndelegate
	EventValidatorCreated
	EventValidatorSlashed
	EventRewardDistributed
)

// GovernanceEvent represents governance-related events
type GovernanceEvent struct {
	Type       GovernanceEventType `json:"type"`
	ProposalID uint64              `json:"proposalId"`
	Voter      Address             `json:"voter"`
	Option     VoteOption          `json:"option"`
}

// GovernanceEventType represents governance event types
type GovernanceEventType uint8

const (
	EventProposalSubmitted GovernanceEventType = iota
	EventVoteCast
	EventProposalPassed
	EventProposalRejected
)