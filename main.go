package main

import "github.com/diamonddelt/dd-blockchain/pkg/block"

// Blockchain is a slice of Blocks
var Blockchain []Block

func main() {
	block.GenerateBlock()
}
