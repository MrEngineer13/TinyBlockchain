package models

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
)

//Block is what gets chained together to form the blockchain
type Block struct {
	Index        int       `json:"index"`
	Timestamp    string    `json:"timestamp"`
	Data         BlockData `json:"data"`
	PreviousHash []byte    `json:"previous_hash"`
	Hash         []byte    `json:"hash"`
}

func (b Block) ToString() string {
	return fmt.Sprintf("Index: %d Timestamp: %s Data: %s PrevHash: %x Hash: %x", b.Index, b.Timestamp, b.Data.ToString(), b.PreviousHash, b.Hash)
}

func (block Block) GetHash() []byte {
	s := []string{strconv.Itoa(block.Index), block.Timestamp, block.Data.ToString(), string(block.PreviousHash[:]), "\n"}
	fmt.Printf(strings.Join(s, ""))
	h := sha256.New()
	h.Write([]byte(strings.Join(s, "")))
	fmt.Printf("Hash : %x", h.Sum(nil))
	return h.Sum(nil)
}

type Blockchain []Block

func (chain Blockchain) ToString() string {
	var buffer bytes.Buffer
	for _, block := range chain {
		buffer.WriteString(block.ToString())
	}
	return buffer.String()
}

//Transaction is to record the interactions to fill the data in a block
type Transaction struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

func (t Transaction) ToString() string {
	return fmt.Sprintf("From: %s To: %s Amount: %d", t.From, t.To, t.Amount)
}

type BlockData struct {
	Proof        int           `json:"proof-of-work"`
	Transactions []Transaction `json:"transactions"`
}

func (d BlockData) ToString() string {
	var buffer bytes.Buffer
	for i := 0; i < len(d.Transactions); i++ {
		buffer.WriteString(d.Transactions[i].ToString())
	}
	return fmt.Sprintf(" Proof: %d Transactions %s ", d.Proof, buffer.String())
}
