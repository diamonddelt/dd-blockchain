package main

import (
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/diamonddelt/dd-blockchain/pkg/block"
	"github.com/diamonddelt/dd-blockchain/pkg/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load() // load all environment variables
	if err != nil {
		log.Fatal(err)
	}

	// spawn a goroutine to generate the initial Blockchain data
	go func() {
		t := time.Now()
		genesisBlock := block.Block{0, t.String(), 0, "", ""}
		spew.Dump(genesisBlock)
		block.Blockchain = append(block.Blockchain, genesisBlock)
	}()

	// initialize the Blockchain server
	log.Fatal(server.Run())
}
