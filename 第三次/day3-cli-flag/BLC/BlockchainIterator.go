package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"math/big"
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

}
func (blc *Blockchain) PrintChain() {
	var block *Block
	//var currentHash []byte = blc.Tip
	blcIterator := blc.Iterator()
	for {
		block := blcIterator.next()
		var hashInt big.Int
		hashInt.SetBytes(block.prevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}
