package evm

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	pixelzxtypes "github.com/pixelzx/pos/core/types"
)

// Executor handles EVM transaction execution
type Executor struct {
	chainConfig *params.ChainConfig
	vmConfig    vm.Config
	stateDB     *state.StateDB
	chainID     *big.Int
	
	mu sync.RWMutex
}

// ExecutorConfig represents EVM executor configuration
type ExecutorConfig struct {
	ChainID     *big.Int
	ChainConfig *params.ChainConfig
	VMConfig    vm.Config
	StateDB     *state.StateDB
}

// NewExecutor creates a new EVM executor
func NewExecutor(config *ExecutorConfig) *Executor {
	return &Executor{
		chainConfig: config.ChainConfig,
		vmConfig:    config.VMConfig,
		stateDB:     config.StateDB,
		chainID:     config.ChainID,
	}
}

// ExecuteTransactions executes a batch of transactions
func (e *Executor) ExecuteTransactions(
	ctx context.Context,
	txs []*types.Transaction,
	header *pixelzxtypes.Header,
	stateDB *state.StateDB,
) (*ExecutionResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Create block context
	blockContext := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     func(uint64) common.Hash { return common.Hash{} }, // Simplified
		Coinbase:    header.Proposer,
		BlockNumber: new(big.Int).SetUint64(header.Number),
		Time:        new(big.Int).SetUint64(header.Timestamp),
		Difficulty:  big.NewInt(1), // PoS has difficulty 1
		GasLimit:    header.GasLimit,
		BaseFee:     big.NewInt(pixelzxtypes.DefaultGasPrice),
	}

	// Execute transactions
	var (
		receipts     []*types.Receipt
		totalGasUsed uint64
		allLogs      []*types.Log
		failed       []int // Indices of failed transactions
	)

	for i, tx := range txs {
		receipt, gasUsed, err := e.executeTransaction(ctx, tx, i, stateDB, blockContext)
		if err != nil {
			// Log error but continue with other transactions
			failed = append(failed, i)
			continue
		}

		receipts = append(receipts, receipt)
		totalGasUsed += gasUsed
		allLogs = append(allLogs, receipt.Logs...)

		// Check gas limit
		if totalGasUsed > header.GasLimit {
			return nil, fmt.Errorf("block gas limit exceeded")
		}
	}

	// Calculate roots
	receiptRoot := types.DeriveSha(types.Receipts(receipts), types.NewTrie())
	stateRoot := stateDB.IntermediateRoot(true)

	return &ExecutionResult{
		StateRoot:    stateRoot,
		ReceiptRoot:  receiptRoot,
		Receipts:     receipts,
		GasUsed:      totalGasUsed,
		Logs:         allLogs,
		FailedTxs:    failed,
	}, nil
}

// executeTransaction executes a single transaction
func (e *Executor) executeTransaction(
	ctx context.Context,
	tx *types.Transaction,
	txIndex int,
	stateDB *state.StateDB,
	blockContext vm.BlockContext,
) (*types.Receipt, uint64, error) {
	// Create transaction context
	txContext := vm.TxContext{
		Origin:   common.Address{}, // Will be set from signature
		GasPrice: tx.GasPrice(),
	}

	// Recover sender from signature
	signer := types.LatestSignerForChainID(e.chainID)
	sender, err := types.Sender(signer, tx)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid transaction signature: %w", err)
	}
	txContext.Origin = sender

	// Create EVM instance
	evm := vm.NewEVM(blockContext, txContext, stateDB, e.chainConfig, e.vmConfig)

	// Apply transaction
	result, err := core.ApplyTransaction(
		e.chainConfig,
		nil, // chain context
		&blockContext.Coinbase,
		core.NewGasPool(blockContext.GasLimit),
		stateDB,
		nil, // header
		tx,
		&blockContext.GasLimit,
		evm,
	)

	if err != nil {
		return nil, 0, fmt.Errorf("transaction execution failed: %w", err)
	}

	return result, result.GasUsed, nil
}

// EstimateGas estimates gas required for a transaction
func (e *Executor) EstimateGas(
	ctx context.Context,
	from common.Address,
	to *common.Address,
	value *big.Int,
	data []byte,
	stateDB *state.StateDB,
) (uint64, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Create a dummy transaction for estimation
	tx := types.NewTransaction(
		stateDB.GetNonce(from),
		*to,
		value,
		21000, // Start with minimum gas
		big.NewInt(pixelzxtypes.DefaultGasPrice),
		data,
	)

	// Binary search for optimal gas limit
	lo, hi := uint64(21000), uint64(pixelzxtypes.DefaultGasLimit)
	
	for lo < hi {
		mid := (lo + hi) / 2
		
		// Create test transaction with current gas limit
		testTx := types.NewTransaction(
			tx.Nonce(),
			*tx.To(),
			tx.Value(),
			mid,
			tx.GasPrice(),
			tx.Data(),
		)

		// Test execution
		if e.canExecuteTransaction(testTx, from, stateDB) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}

	return lo, nil
}

// canExecuteTransaction checks if a transaction can be executed with given gas
func (e *Executor) canExecuteTransaction(tx *types.Transaction, sender common.Address, stateDB *state.StateDB) bool {
	// Create a snapshot for testing
	snapshot := stateDB.Snapshot()
	defer stateDB.RevertToSnapshot(snapshot)

	// Basic validation
	balance := stateDB.GetBalance(sender)
	cost := new(big.Int).Mul(new(big.Int).SetUint64(tx.Gas()), tx.GasPrice())
	cost.Add(cost, tx.Value())

	if balance.Cmp(cost) < 0 {
		return false
	}

	// Try to execute
	blockContext := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     func(uint64) common.Hash { return common.Hash{} },
		Coinbase:    common.Address{},
		BlockNumber: big.NewInt(1),
		Time:        big.NewInt(1),
		Difficulty:  big.NewInt(1),
		GasLimit:    pixelzxtypes.DefaultGasLimit,
		BaseFee:     big.NewInt(pixelzxtypes.DefaultGasPrice),
	}

	txContext := vm.TxContext{
		Origin:   sender,
		GasPrice: tx.GasPrice(),
	}

	evm := vm.NewEVM(blockContext, txContext, stateDB, e.chainConfig, e.vmConfig)

	// Apply transaction
	_, err := core.ApplyTransaction(
		e.chainConfig,
		nil,
		&blockContext.Coinbase,
		core.NewGasPool(blockContext.GasLimit),
		stateDB,
		nil,
		tx,
		&blockContext.GasLimit,
		evm,
	)

	return err == nil
}

// CallContract executes a read-only contract call
func (e *Executor) CallContract(
	ctx context.Context,
	from common.Address,
	to common.Address,
	data []byte,
	stateDB *state.StateDB,
) ([]byte, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Create block and transaction contexts
	blockContext := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     func(uint64) common.Hash { return common.Hash{} },
		Coinbase:    common.Address{},
		BlockNumber: big.NewInt(1),
		Time:        big.NewInt(1),
		Difficulty:  big.NewInt(1),
		GasLimit:    pixelzxtypes.DefaultGasLimit,
		BaseFee:     big.NewInt(pixelzxtypes.DefaultGasPrice),
	}

	txContext := vm.TxContext{
		Origin:   from,
		GasPrice: big.NewInt(pixelzxtypes.DefaultGasPrice),
	}

	// Create EVM instance
	evm := vm.NewEVM(blockContext, txContext, stateDB, e.chainConfig, e.vmConfig)

	// Execute call
	result, leftOverGas, err := evm.Call(
		vm.AccountRef(from),
		to,
		data,
		pixelzxtypes.DefaultGasLimit,
		big.NewInt(0),
	)

	if err != nil {
		return nil, fmt.Errorf("contract call failed: %w", err)
	}

	_ = leftOverGas // Can be used for gas calculation

	return result, nil
}

// CreateDefaultChainConfig creates a default chain configuration for PIXELZX
func CreateDefaultChainConfig(chainID *big.Int) *params.ChainConfig {
	return &params.ChainConfig{
		ChainID:             chainID,
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		BerlinBlock:         big.NewInt(0),
		LondonBlock:         big.NewInt(0),
		// Disable difficulty bomb for PoS
		TerminalTotalDifficulty: big.NewInt(0),
	}
}

// ExecutionResult represents the result of transaction execution
type ExecutionResult struct {
	StateRoot   common.Hash      `json:"stateRoot"`
	ReceiptRoot common.Hash      `json:"receiptRoot"`
	Receipts    []*types.Receipt `json:"receipts"`
	GasUsed     uint64           `json:"gasUsed"`
	Logs        []*types.Log     `json:"logs"`
	FailedTxs   []int            `json:"failedTxs"`
}

// TransactionPool manages pending transactions
type TransactionPool struct {
	pending map[common.Address]*types.Transaction
	queued  map[common.Address][]*types.Transaction
	
	mu sync.RWMutex
}

// NewTransactionPool creates a new transaction pool
func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		pending: make(map[common.Address]*types.Transaction),
		queued:  make(map[common.Address][]*types.Transaction),
	}
}

// AddTransaction adds a transaction to the pool
func (tp *TransactionPool) AddTransaction(tx *types.Transaction) error {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	// Recover sender
	signer := types.LatestSigner(CreateDefaultChainConfig(big.NewInt(pixelzxtypes.DefaultChainID)))
	sender, err := types.Sender(signer, tx)
	if err != nil {
		return fmt.Errorf("invalid transaction signature: %w", err)
	}

	// Add to pending
	tp.pending[sender] = tx

	return nil
}

// GetPendingTransactions returns pending transactions
func (tp *TransactionPool) GetPendingTransactions() []*types.Transaction {
	tp.mu.RLock()
	defer tp.mu.RUnlock()

	txs := make([]*types.Transaction, 0, len(tp.pending))
	for _, tx := range tp.pending {
		txs = append(txs, tx)
	}

	return txs
}

// RemoveTransactions removes transactions from the pool
func (tp *TransactionPool) RemoveTransactions(txs []*types.Transaction) {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	signer := types.LatestSigner(CreateDefaultChainConfig(big.NewInt(pixelzxtypes.DefaultChainID)))
	
	for _, tx := range txs {
		if sender, err := types.Sender(signer, tx); err == nil {
			delete(tp.pending, sender)
		}
	}
}

// GetTransactionCount returns the number of pending transactions
func (tp *TransactionPool) GetTransactionCount() int {
	tp.mu.RLock()
	defer tp.mu.RUnlock()

	return len(tp.pending)
}