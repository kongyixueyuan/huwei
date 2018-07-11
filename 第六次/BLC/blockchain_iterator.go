package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type HW_BlockchainIterator struct {
	hw_currentHash []byte
	hw_db          *bolt.DB
}

func (i *HW_BlockchainIterator) HW_Next() *HW_Block {
	var block *HW_Block

	err := i.hw_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.hw_currentHash)
		block = HW_DeserializeBlock(encodedBlock)
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.hw_currentHash = block.Hw_PrevBlockHash

	return block
}
