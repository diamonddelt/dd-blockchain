package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/diamonddelt/dd-blockchain/pkg/block"
	"github.com/diamonddelt/dd-blockchain/pkg/message"
	"github.com/diamonddelt/dd-blockchain/pkg/validation"
	"github.com/gorilla/mux"
)

// Run create a locahost webserver on a port specified within the .env file
func Run() error {
	mux := createMuxRouter()
	addr := os.Getenv("ADDR")
	log.Println("Blockchain server listening on port", os.Getenv("ADDR"))

	s := &http.Server{
		Addr:           ":" + addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// createMuxRouter registers server routes and HTTP verbs for the Mux router
func createMuxRouter() http.Handler {
	mux := mux.NewRouter()

	// routes
	mux.HandleFunc("/", handleGetBlockchain).Methods("GET") // get current blockchain
	mux.HandleFunc("/", handleWriteBlock).Methods("POST")   // add to blockchain

	return mux
}

// handleGetBlockchain prints the current blockchain in JSON format to localhost:{ENV.ADDR}
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(block.Blockchain, "", "	")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(bytes))
}

// handleWriteBlock adds a new block to the blockchain
func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m message.Message

	// parse the HTTP response as JSON
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	// generate a new Block on the blockchain
	newBlock, err := block.GenerateBlock(block.Blockchain[len(block.Blockchain)-1], m.BPM)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}

	// append the new Block to the blockchain if it passes integrity checks
	if validation.ValidateBlockIntegrity(newBlock, block.Blockchain[len(block.Blockchain)-1]) {
		newBlockchain := append(block.Blockchain, newBlock)
		block.UpdateBlockchain(newBlockchain)
		spew.Dump(block.Blockchain) // spew is a useful library for debugging structs
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

// respondWithJSON is a wrapper function for handling HTTP responses to and from the Blockchain server
func respondWithJSON(w http.ResponseWriter, r *http.Request, resCode int, payload interface{}) {
	res, err := json.MarshalIndent(payload, "", "	")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(resCode)
	w.Write(res)
}
