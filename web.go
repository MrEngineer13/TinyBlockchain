package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"./models"
	"./utils"
)

//// A completely random address of the owner of this node
var minerAddress = "q3nf394hjg-random-miner-address-34nf3i4nflkn3oi"

// Store the transactions that
// this node has in a list
var thisNodesTransactions = []models.Transaction{}

// Store the url data of every
// other node in the network
// so that we can communicate
// with them, this list is prepopulated
// with the local host for testing
var peer_nodes = []string{"http://localhost:9090"}

// This node's blockchain copy
var blockchain = []models.Block{utils.CreateGenesisBlock()}

func main() {
	mux := &MyMux{}
	http.ListenAndServe(":9090", mux)
}

type MyMux struct {
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		nextBlock()
		return
	} else if r.URL.Path == "/search" {
		findNewChains()
		return
	} else if r.Method == "POST" && r.URL.Path == "/txion" {
		fmt.Fprintf(w, transaction(w, r))
		return
	} else if r.URL.Path == "/mine" {
		mine(w, r)
		return
	} else if r.URL.Path == "/blocks" {
		getBlocks(w, r)
		return
	}
	http.NotFound(w, r)
	return
}

func transaction(w http.ResponseWriter, r *http.Request) string {
	// On each new POST request,
	// we extract the transaction data
	decoder := json.NewDecoder(r.Body)
	var newTxion models.Transaction
	err := decoder.Decode(&newTxion)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	// Then we add the transaction to our list
	thisNodesTransactions = append(thisNodesTransactions[:], newTxion)
	// Because the transaction was successfully
	// submitted, we log it to our console

	fmt.Fprintf(w, "\nNew transaction %s \n", newTxion.ToString())
	// Then we let the client know it worked out
	return "Transaction submission successful\n"
}

func getBlocks(w http.ResponseWriter, r *http.Request) []models.Block {
	consensus()
	json.NewEncoder(w).Encode(blockchain)
	return blockchain
}

func findNewChains() []models.Blockchain {
	// Get the blockchains of every
	// other node
	otherChains := []models.Blockchain{}

	for _, nodeUrl := range peer_nodes {
		// Convert the JSON object to a block
		chain := models.Blockchain{}
		// Get their chains using a GET request
		url := fmt.Sprintf("%s/blocks", nodeUrl)
		getJson(url, &chain)

		// Add it to our list
		otherChains = append(otherChains[:], chain)
	}
	return otherChains
}

func getJson(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 60 * time.Second}

	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func consensus() {
	// Get the blocks from other nodes
	otherChains := []models.Blockchain{}
	otherChains = findNewChains()
	// If our chain isn't longest,
	// then we store the longest chain
	longestChain := blockchain
	for _, chain := range otherChains {
		if len(longestChain) < len(chain) {
			longestChain = chain
		}
	}
	// If the longest chain isn't ours,
	// then we stop mining and set
	// our chain to the longest one
	blockchain = longestChain
}

func proofOfWork(lastproof int) int {
	// Create a variable that we will use to find
	// our next proof of work
	incrementor := lastproof + 1
	// Keep incrementing the incrementor until
	// it's equal to a number divisible by 9
	// and the proof of work of the previous
	// block in the chain
	for incrementor%9 != 0 && incrementor%lastproof != 0 {
		incrementor += 1
	}

	// Once that number is found,
	// we can return it as a proof
	// of our work
	return incrementor
}

func mine(w http.ResponseWriter, r *http.Request) {
	// Get the last proof of work
	lastBlock := blockchain[len(blockchain)-1]
	lastProof := lastBlock.Data.Proof
	// Find the proof of work for
	// the current block being mined
	// Note: The program will hang here until a new
	//       proof of work is found
	proof := proofOfWork(lastProof)
	// Once we find a valid proof of work,
	// we know we can mine a block so
	// we reward the miner by adding a transaction
	newTxion := models.Transaction{From: "network", To: minerAddress, Amount: 1}
	thisNodesTransactions = append(thisNodesTransactions[:], newTxion)
	// Now we can gather the data needed
	// to create the new block
	newBlockData := models.BlockData{Proof: proof, Transactions: thisNodesTransactions}
	newBlockIndex := lastBlock.Index + 1
	newBlockTimestamp := utils.GetTime()
	lastBlockHash := lastBlock.Hash
	// Empty transaction list
	thisNodesTransactions = []models.Transaction{}
	// Now create the
	// new block!

	minedBlock := models.Block{Index: newBlockIndex, Timestamp: newBlockTimestamp, Data: newBlockData, PreviousHash: lastBlockHash}
	minedBlock.Hash = minedBlock.GetHash()
	blockchain = append(blockchain[:], minedBlock)
	// Let the client know we mined a block
	fmt.Printf("Mined block %s\n", minedBlock.ToString())
	json.NewEncoder(w).Encode(minedBlock)
}

func nextBlock() {
	const NUMBER_OF_BLOCKS = 20

	previousblock := blockchain[0]

	for i := 0; i < NUMBER_OF_BLOCKS; i++ {
		blockToAdd := utils.GetNextBlock(previousblock)
		blockchain = append(blockchain[:], blockToAdd)
		previousblock = blockToAdd
		fmt.Printf("\nBlock %d has been added to the blockchain!\n", blockToAdd.Index)
		fmt.Printf("Hash: %x\n\n", string(blockToAdd.Hash[:]))
	}
}
