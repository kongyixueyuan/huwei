package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockchainIterator struct {
	CurrentHash []byte
	DB *bolt.DB
}

// 迭代器
func (blockchain *Blockchain) Iterator() *BlockchainIterator {

	return &BlockchainIterator{blockchain.Tip, blockchain.DB}

}

func (blockchainIterator *BlockchainIterator) Next() *Block {

	var block *Block

	err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			currentBlockBytes := b.Get(blockchainIterator.CurrentHash)

			// 获取到当前迭代器里面的currentHash所对于的区块
			block = DeserializeBlock(currentBlockBytes)

			// 更新迭代里面的currentHash
			blockchainIterator.CurrentHash = block.prevBlockHash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return block

}