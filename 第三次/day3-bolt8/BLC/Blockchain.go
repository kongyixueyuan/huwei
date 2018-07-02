package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

// 数据库名字
const dbName = "blockchain.db"
const blockTableName = "blocks"

type Blockchain struct {
	Tip []byte //最新区块的hash
	DB  *bolt.DB
}

//编写一个Iterator
func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockchain.Tip, blockchain.DB}
}

//增加区块到区块链里面
func (blc *Blockchain) AddBlockToBlockChain(data string) {
	err := blc.DB.Update(func(tx *bolt.Tx) error {
		//1.获取表
		b := tx.Bucket([]byte(blockTableName))
		//2 .创建新区块
		if b != nil {
			//先取最新区块
			blockBytes := b.Get(blc.Tip)
			//反序列化
			block := DeserializeBlock(blockBytes)
			// 将新区块序列化保存到数据库
			newBlock := NewBlock(data, block.Height+1, block.Hash)
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Fatal(err)
			}
			// 4. 更新数据库里面"H"对应的Hash
			err = b.Put([]byte("H"), newBlock.Hash)
			if err != nil {
				log.Fatal(err)
			}
			// 5. 更新blockchain的Tip
			blc.Tip = newBlock.Hash
		}
		return nil

	})
	if err != nil {
		log.Panic(err)
	}
}

//1. 创建带有创世区块的区块链
func CreateBlockChainWithGeneisBlock() *Blockchain {
	//创建创世区块
	/*
	genesisBlock := CreateGenesisBlock("Genenis block......")
	//返回区块链对象
	return &Blockchain{[]*Block{genesisBlock}}
	*/
	//2018-06-30 创建或打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var blockHash []byte
	err = db.Update(func(tx *bolt.Tx) error {
		//创建数据库表
		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil {
			log.Fatal(err)
		}
		if b != nil {
			//创建创世区块
			genesisBlock := CreateGenesisBlock("Genesis Data ......")
			//将创世区块保存在区块链上
			err := b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Fatal(err)
			}
			//保存最新区块的hash
			err = b.Put([]byte("H"), genesisBlock.Hash)
			if err != nil {
				log.Fatal(err)
			}
			blockHash = genesisBlock.Hash
		}
		return nil
	})
	return &Blockchain{blockHash, db}
}
