package utils

import (
	"time"

	"../models"
)

func GetTime() string {
	current_time := time.Now().Local()
	return current_time.Format("Jan 2 15:04:05 -0700 MST 2006")
}

func CreateGenesisBlock() models.Block {
	blockData := models.BlockData{Proof: 9}
	block := models.Block{Index: 0, Timestamp: GetTime(), Data: blockData, PreviousHash: []byte{0}}
	block.Hash = block.GetHash()
	return block
}

func GetNextBlock(lastBlock models.Block) models.Block {
	blockData := models.BlockData{Proof: lastBlock.Data.Proof + 1}
	newBlock := models.Block{Index: lastBlock.Index + 1, Timestamp: GetTime(), Data: blockData, PreviousHash: lastBlock.Hash}
	newBlock.Hash = newBlock.GetHash()
	return newBlock
}
