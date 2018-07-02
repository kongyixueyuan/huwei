package main

import (
	"fmt"
	"publicchain/day3-bolt8/BLC"
	"encoding/hex"
)

func main() {
	/*db, err := bolt.Open("block.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//view 读取数据
	err = db.View(func(tx *bolt.Tx) error {
		table := tx.Bucket([]byte("block"))
		if table != nil {
			blockBytes := table.Get([]byte("huwei"))
			fmt.Println(blockBytes)
			block := BLC.DeserializeBlock(blockBytes)
			fmt.Println(block.Nonce)

		}
		return nil

	})
	if err != nil {
		log.Panic(err)
	}*/
	blockchain := BLC.CreateBlockChainWithGeneisBlock()
	// 新区块
	blockchain.AddBlockToBlockChain("send 1000RMB to hovi")
	blockchain.AddBlockToBlockChain("send 2000RMB to hovi")
	blockchain.AddBlockToBlockChain("send 3000RMB to hovi")
	blockchain.AddBlockToBlockChain("send 4000RMB to hovi")
	fmt.Println(hex.EncodeToString(blockchain.Tip))
	// 打印所有区块的信息
	blockchain.PrintChain()

}
