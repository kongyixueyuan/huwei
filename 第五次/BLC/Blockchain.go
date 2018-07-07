package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"os"
	"encoding/hex"
	"math/big"
	"strconv"
	"time"
	"crypto/ecdsa"
	"bytes"
)

// 数据库名字
const dbName = "blockchain.db"
const blockTableName = "blocks"

type Blockchain struct {
	Tip []byte //最新的区块hash
	DB  *bolt.DB
}


// 遍历输出所有区块的信息
func (blc *Blockchain) Printchain() {

	blockchainIterator := blc.Iterator()

	for {
		block := blockchainIterator.Next()

		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("PreBlockHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Timestamp: %s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Println("Transaction:")
		for _, tx := range block.Txs {
			fmt.Printf("TxHash:%x\n", tx.TxHash)
			fmt.Println("Vins:")
			for _, in := range tx.Vins {
				fmt.Printf("%x\n", in.Txhash)
				fmt.Printf("%d\n", in.Vout)
				fmt.Printf("%s\n", in.PubKey)
			}
			fmt.Println("Vouts:")
			for _, out := range tx.Vouts {
				fmt.Println(out.Value)
				fmt.Println(out.Ripemd160Hash)
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}

}

// 判断数据库是否存在
func DBExists() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}

	return true
}
// 增加区块到区块链里面
func (blc *Blockchain) AddBlockToBlockchain(txs []*Transaction) {

	err := blc.DB.Update(func(tx *bolt.Tx) error {
		// 1. 获取表
		b := tx.Bucket([]byte(blockTableName))

		// 2. 创建新区块
		if b != nil {
			// 从数据库中取到上一个区块的信息（获取最新区块）
			blockBytes := b.Get(blc.Tip)
			// 反序列化
			block := DeserializeBlock(blockBytes)
			// 3. 将区块序列化，存储到数据库中
			newBlock := NewBlock(txs, block.Height+1, block.Hash)
			// 保存生成新的区块到数据库中
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			// 4. 更新数据库里面 "H" 对应的hash
			err = b.Put([]byte("H"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
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
// 创建带有创世区块的区块链
func CreateBlockchainWithGenenisBlock(address string,amount int64) *Blockchain {

	if DBExists() {
		fmt.Println("创世区块已经存在")
		os.Exit(1)
	}
	fmt.Println("正在创建创世区块")
	// 创建或者打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var genesisHash []byte

	// 更新数据库
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			// 创建创世区块
			// 创建了一个coinbase Transaction
			txCoinbase := NewCoinbaseTransaction(address,amount)
			//创建创世区块
			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
			// 将创世区块存储到表中
			err := b.Put(genesisBlock.Hash, genesisBlock.Serialize())

			if err != nil {
				log.Panic(err)
			}

			// 存储最新的区块的Hash
			err = b.Put([]byte("H"), genesisBlock.Hash)

			if err != nil {
				log.Panic(err)
			}

			genesisHash = genesisBlock.Hash

		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{genesisHash, db}

}

// 返回Blockchain对象
func BlockchainObject() *Blockchain {
	//如果数据库不存在 ，无法获取区块链
	if ! DBExists() {
		fmt.Println("数据库不存在,创世区块还未生成")
		os.Exit(1)
	}
	// 创建或者打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var tip []byte
    //查看表
	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			// 读取最新区块的hash
			tip = b.Get([]byte("H"))
		}

		return nil
	})
	return &Blockchain{tip, db}

}

// 如果一个地址对应的TXOutput未花费，那么这个Transaction就应该添加到数组中返回
func (blockchain *Blockchain) UnUTXOs(address string, txs []*Transaction) []*UTXO {
	// 用于存储未花费的Transaction
	var unUTXOs []*UTXO

	spentTXOutputs := make(map[string][]int)

	for _, tx := range txs {

		// Vins
		if tx.IsCoinbaseTransaction() == false {
			for _, in := range tx.Vins {
				publicKeyHash := Base58Decode([]byte(address))
				ripemd160Hash := publicKeyHash[1 : len(publicKeyHash)-addressChecksumLen]
				// 是否能够解锁
				if in.UnLockWithAddress(ripemd160Hash) {

					key := hex.EncodeToString(in.Txhash)

					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)

				}

			}
		}

	}

	for _, tx := range txs {
	label:
		for index, out := range tx.Vouts {

			if out.UnLockScriptPubKeyWithAddress(address) {
				fmt.Println(address)
				fmt.Println(spentTXOutputs)
				if len(spentTXOutputs) == 0 {
					utxo := &UTXO{tx.TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash, indexArray := range spentTXOutputs {
						txHashStr := hex.EncodeToString(tx.TxHash)
						if hash == txHashStr {

							var isSpentUTXO bool

							for _, outIndex := range indexArray {
								if index == outIndex {
									isSpentUTXO = true
									continue label
								}

								if isSpentUTXO == false {
									utxo := &UTXO{tx.TxHash, index, out}
									unUTXOs = append(unUTXOs, utxo)
								}
							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}

					}

				}

			}

		}

	}

	blockIterator := blockchain.Iterator()

	for {

		block := blockIterator.Next()
		fmt.Println(block)
		fmt.Println()
		for i := len(block.Txs) - 1; i >= 0; i-- {

			tx := block.Txs[i]
			// txHash

			// Vins
			if tx.IsCoinbaseTransaction() == false {
				for _, in := range tx.Vins {
					// 是否能够解锁
					publicKeyHash := Base58Decode([]byte(address))
					ripemd160Hash := publicKeyHash[1 : len(publicKeyHash)-addressChecksumLen]
					if in.UnLockWithAddress(ripemd160Hash) {

						key := hex.EncodeToString(in.Txhash)
						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)

					}

				}
			}

			// Vouts
		label1:
			for index, out := range tx.Vouts {

				if out.UnLockScriptPubKeyWithAddress(address) {

					if spentTXOutputs != nil {

						if len(spentTXOutputs) != 0 {

							var isSpentUTXO bool

							for txHash, indexArray := range spentTXOutputs {

								for _, i := range indexArray {
									if index == i && txHash == hex.EncodeToString(tx.TxHash) {
										isSpentUTXO = true
										continue label1
									}
								}

							}

							if isSpentUTXO == false {
								utxo := &UTXO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)
							}

						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}

					}

				}

			}

		}
		fmt.Println(spentTXOutputs)
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}

	}

	return unUTXOs

}

// 转账时查找可用的UTXO
func (blockchain *Blockchain) FindSpendableUTXOS(from string, amount int, txs []*Transaction) (int64, map[string][]int) {

	// 获取所有的UTXO
	utxos := blockchain.UnUTXOs(from, txs)

	spendableUTXO := make(map[string][]int)

	var value int64

	// 遍历utxos
	for _, utxo := range utxos {

		value = value + utxo.Output.Value

		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)

		if value >= int64(amount) {
			break
		}

	}

	if value < int64(amount) {
		fmt.Printf("%s's fund is not enough\n", from)
		os.Exit(1)
	}

	return value, spendableUTXO
}

// 挖掘新的区块
func (blockchain *Blockchain) MineNewBlock(from []string, to []string, amount []string) {
	// 通过相关算法建立Transantion数组
	var txs []*Transaction
	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := NewSimpleTransaction(address, to[index], value, blockchain, txs)
		txs = append(txs, tx)
	}

	var block *Block
	blockchain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			hash := b.Get([]byte("H"))
			blockBytes := b.Get(hash)
			block = DeserializeBlock(blockBytes)
		}
		return nil

	})

	// 建立新的区块
	block = NewBlock(txs, block.Height+1, block.Hash)

	// 将新区块存储到数据库
	blockchain.DB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {

			b.Put(block.Hash, block.Serialize())
			b.Put([]byte("H"), block.Hash)
			blockchain.Tip = block.Hash
		}
		return nil

	})

}

// 查询余额
func (blockchain *Blockchain) GetBalance(address string) int64 {
	utxos := blockchain.UnUTXOs(address, []*Transaction{})
	var amount int64
	for _, out := range utxos {
		amount = amount + out.Output.Value
	}
	return amount
}
//通过 ID 找到一笔交易（这需要在区块链上迭代所有区块）
func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {

	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Txs {
			if bytes.Compare(tx.TxHash, ID) == 0 {
				return *tx, nil
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}
	}

	return Transaction{},nil
}
//传入一笔交易，找到它引用的交易，然后对它进行签名
func (blockchain *Blockchain) SignTransaction(tx *Transaction, privateKey ecdsa.PrivateKey)  {
	if tx.IsCoinbaseTransaction(){
		return
	}

	prevTXs := make(map[string]Transaction)

	for _,vin := range tx.Vins{
		prevTX ,err := blockchain.FindTransaction(vin.Txhash)
		if err !=nil{
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}
	tx.Sign(privateKey, prevTXs)
}


//传入一笔交易，找到它引用的交易，然后对它进行验证
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vins {
		prevTX, _ := bc.FindTransaction(vin.Txhash)
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}

	return tx.Verify(prevTXs)
}
