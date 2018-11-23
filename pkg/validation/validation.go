package validation

import "github.com/diamonddelt/dd-blockchain/pkg/block"

// ValidateBlockIntegrity ensures Blockchain integrity
// Validates that a new Block has incremented its Index, that the PreviousHash
// matches the CurrentHash, and that SHA256 hash generated matches the CurrentHash
func ValidateBlockIntegrity(new, old block.Block) bool {
	if old.Index+1 != new.Index {
		return false
	}

	if old.CurrentHash != new.PreviousHash {
		return false
	}

	if block.CalculateBlockHash(new) != new.CurrentHash {
		return false
	}

	return true
}
