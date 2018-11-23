package block

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Block represents a blockchain unit of work
type Block struct {
	Index        int    // position of the data record on the blockchain
	Timestamp    string // when the data is written
	BPM          int    // the pulse rate
	CurrentHash  string // SHA256 identifer of this record
	PreviousHash string // SHA256 identifier of previous record
}

// Blockchain is an ever-growing, source of truth array of Blocks
var Blockchain []Block

// GenerateBlock creates a new Block based on a previous Block in the Blockchain
func GenerateBlock(old Block, BPM int) (Block, error) {
	var new Block

	t := time.Now()

	new.Index = old.Index + 1
	new.Timestamp = t.String()
	new.BPM = BPM
	new.PreviousHash = old.CurrentHash
	new.CurrentHash = CalculateBlockHash(new)

	return new, nil
}

// CalculateBlockHash creates a SHA256 hash of a given Block
func CalculateBlockHash(b Block) string {
	record := string(b.Index) + b.Timestamp + string(b.BPM) + b.PreviousHash
	hash := sha256.New()
	hash.Write([]byte(record))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

// UpdateBlockchain updates the current Blockchain with a potentially newer Blockchain
// if the candidate Blockchain has a longer chain length
func UpdateBlockchain(potential []Block) {
	if len(potential) > len(Blockchain) {
		Blockchain = potential
	}
}
