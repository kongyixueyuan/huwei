package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type Hw_BlockchainIterator struct {
	Hw_currentHash []byte
	Hw_db          *bolt.DB
}

func (i *Hw_BlockchainIterator) Hw_Next() *Hw_Block {
	var block *Hw_Block

	err := i.Hw_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.Hw_currentHash)
		block = Hw_DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.Hw_currentHash = block.Hw_PrevBlockHash

	return block
}