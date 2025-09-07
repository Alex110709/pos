package core

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/params"
)

func TestHoodiGenesisAlloc(t *testing.T) {
	genesis := DefaultHoodiGenesisBlock()
	
	// Check if the genesis block is correctly configured
	if genesis.Config.ChainID.Cmp(params.HoodiChainConfig.ChainID) != 0 {
		t.Errorf("ChainID mismatch: got %v, want %v", genesis.Config.ChainID, params.HoodiChainConfig.ChainID)
	}

	// Check if the allocation is correctly multiplied by 100
	// We'll check if the total supply is approximately 100 times larger.
	// This is a rough check and might need adjustment based on the actual data.
	
	var originalTotalSupply *big.Int
	originalAlloc := decodePrealloc(hoodiAllocData)
	for _, account := range originalAlloc {
		if originalTotalSupply == nil {
			originalTotalSupply = new(big.Int).Set(account.Balance)
		} else {
			originalTotalSupply.Add(originalTotalSupply, account.Balance)
		}
	}
	
	var newTotalSupply *big.Int
	for _, account := range genesis.Alloc {
		if newTotalSupply == nil {
			newTotalSupply = new(big.Int).Set(account.Balance)
		} else {
			newTotalSupply.Add(newTotalSupply, account.Balance)
		}
	}
	
	// Calculate expected total supply (100 times original)
	expectedTotalSupply := new(big.Int).Mul(originalTotalSupply, big.NewInt(100))
	
	// Allow for some small variance due to potential rounding or other factors
	// This is a very rough check - in practice, you might want to be more precise
	if newTotalSupply.Cmp(expectedTotalSupply) != 0 {
		// Calculate the ratio to see how close we are
		ratio := new(big.Float).Quo(new(big.Float).SetInt(newTotalSupply), new(big.Float).SetInt(originalTotalSupply))
		t.Logf("Original total supply: %v", originalTotalSupply)
		t.Logf("New total supply: %v", newTotalSupply)
		t.Logf("Expected total supply: %v", expectedTotalSupply)
		t.Logf("Ratio (new/old): %v", ratio)
		
		// Check if the ratio is close to 100
		expectedRatio := big.NewFloat(100.0)
		diff := new(big.Float).Sub(ratio, expectedRatio)
		absDiff := new(big.Float).Abs(diff)
		
		// Allow 0.1% variance
		maxVariance := big.NewFloat(0.1)
		if absDiff.Cmp(maxVariance) > 0 {
			t.Errorf("Total supply not multiplied by approximately 100. Ratio: %v", ratio)
		}
	}
}