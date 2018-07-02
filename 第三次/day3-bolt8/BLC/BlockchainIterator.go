package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"fmt"
)

type BlockchainIterator struct {
	CurrentHash []byte
	DB          *bolt.DB
}

func (blcIterator *BlockchainIterator) next() *Block {
	var block *Block
	err := blcIterator.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			blockBytes := b.Get(blcIterator.CurrentHash)
			block = DeserializeBlock(blockBytes)
			blcIterator.CurrentHash = block.prevBlockHash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return block
}
func (blc *Blockchain) PrintChain() {
	//var block *Block
	blcIterator := blc.Iterator()
	for {
		block := blcIterator.next()
		fmt.Printf("Height:%d\n", block.Height)
		fmt.Printf("PrevBlockHash:%x\n", block.prevBlockHash)
		fmt.Printf("Data:%s\n", block.Data)
		fmt.Printf("Timestamp:%d\n", block.Timestamp)
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Nonce:%d\n", block.Nonce)
		var hashInt big.Int
		hashInt.SetBytes(block.prevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}

	}

}
