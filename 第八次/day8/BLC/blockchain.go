package BLC

import (
	"github.com/boltdb/bolt"
	"os"
	"fmt"
	"log"
	"encoding/hex"
	"strconv"
	"crypto/ecdsa"
	"bytes"
	"github.com/pkg/errors"
)

const dbFile = "blockchain_%s.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "genesis data 08/07/2018 by viky"

type Hw_Blockchain struct {
	Hw_tip []byte
	Hw_db  *bolt.DB
}

// 打印区块链内容
func (bc *Hw_Blockchain) Hw_Printchain() {
	bci := bc.Hw_Iterator()

	for {
		block := bci.Hw_Next()
		block.String()
		if len(block.Hw_PrevBlockHash) == 0 {
			break
		}
	}

}

// 通过交易hash,查找交易
func (bc *Hw_Blockchain) Hw_FindTransaction(ID []byte) (Hw_Transaction, error) {
	bci := bc.Hw_Iterator()
	for {
		block := bci.Hw_Next()
		for _, tx := range block.Hw_Transactions {
			if bytes.Compare(tx.Hw_ID, ID) == 0 {
				return *tx, nil
			}
		}
		if len(block.Hw_PrevBlockHash) == 0 {
			break
		}
	}
	fmt.Printf("查找%x的交易失败",ID)
	return Hw_Transaction{}, errors.New("未找到交易")
}

// FindUTXO finds all unspent transaction outputs and returns transactions with spent outputs removed
func (bc *Hw_Blockchain) FindUTXO() map[string]Hw_TXOutputs {
	// 未花费的交易输出
	// key:交易hash   txID
	UTXO := make(map[string]Hw_TXOutputs)
	// 已经花费的交易txID : TXOutputs.index
	spentTXOs := make(map[string][]int)
	bci := bc.Hw_Iterator()

	for {
		block := bci.Hw_Next()

		// 循环区块中的交易
		for _, tx := range block.Hw_Transactions {
			// 将区块中的交易hash，转为字符串
			txID := hex.EncodeToString(tx.Hw_ID)

		Outputs:
			for outIdx, out := range tx.Hw_Vout { // 循环交易中的 TXOutputs
				// Was the output spent?
				// 如果已经花费的交易输出中，有此输出，证明已经花费
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx { // 如果花费的正好是此笔输出
							continue Outputs // 继续下一次循环
						}
					}
				}

				outs := UTXO[txID] // 获取UTXO指定txID对应的TXOutputs
				outs.Hw_Outputs = append(outs.Hw_Outputs, out)
				UTXO[txID] = outs
			}

			if tx.Hw_IsCoinbase() == false { // 非创世区块
				for _, in := range tx.Hw_Vin {
					inTxID := hex.EncodeToString(in.Hw_Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Hw_Vout)
				}
			}
		}
		// 如果上一区块的hash为0，代表已经到创世区块，循环结束
		if len(block.Hw_PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

// 获取迭代器
func (bc *Hw_Blockchain) Hw_Iterator() *Hw_BlockchainIterator {
	return &Hw_BlockchainIterator{bc.Hw_tip, bc.Hw_db}
}

// 新建区块链(包含创世区块)
func Hw_CreateBlockchain(address string,nodeID string) *Hw_Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if Hw_dbExists(dbFile) {
		fmt.Println("blockchain数据库已经存在.")
		os.Exit(1)
	}

	var tip []byte
	cbtx := Hw_NewCoinbaseTX(address, genesisCoinbaseData)
	genesis := Hw_NewGenesisBlock(cbtx)

	//genesis.String()

	// 打开数据库，如果不存在自动创建
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		// 新区块存入数据库
		err = b.Put(genesis.Hw_Hash, genesis.Hw_Serialize())
		if err != nil {
			log.Panic(err)
		}
		// 将创世区块的hash存入数据库
		err = b.Put([]byte("l"), genesis.Hw_Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hw_Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Hw_Blockchain{tip, db}
}

// 获取blockchain对象
func Hw_NewBlockchain(nodeID string) *Hw_Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if !Hw_dbExists(dbFile) {
		log.Panic("区块链还未创建")
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Hw_Blockchain{tip, db}
}

// 生成新的区块(挖矿)
func (bc *Hw_Blockchain) MineNewBlock(from []string, to []string, amount []string,nodeID string , mineNow bool) *Hw_Block {
	UTXOSet := Hw_UTXOSet{bc}

	wallets, err := Hw_NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}

	var txs []*Hw_Transaction

	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])
		if value<=0 {
			log.Panic("错误：转账金额需要大于0")
		}
		wallet := wallets.Hw_GetWallet(address)
		tx := Hw_NewUTXOTransaction(&wallet, to[index], value, &UTXOSet, txs)
		txs = append(txs, tx)
	}

	if mineNow {
		// 挖矿奖励
		tx := Hw_NewCoinbaseTX(from[0], "")
		txs = append(txs, tx)

		//=====================================
		newBlock := bc.Hw_MineBlock(txs)
		UTXOSet.Update(newBlock)
		return newBlock
	}else{
		// 如果不立即挖矿，将交易写到内存中
		//var txs_all []Hw_Transaction
		//for _,value := range txs{
		//	txs_all= append(txs_all, *value)
		//}
		Hw_sendTxs(knownNodes[0],txs)
		return nil
	}


}

// 挖矿
func (bc *Hw_Blockchain) Hw_MineBlock(txs []*Hw_Transaction) *Hw_Block  {
	var lashHash []byte
	var lastHeight int

	// 检查交易是否有效，验证签名
	for _, tx := range txs {
		if !bc.Hw_VerifyTransaction(tx, txs) {
			log.Panic("错误：无效的交易")
		}
	}
	// 获取最后一个区块的hash,然后获取最后一个区块的信息，进而获得height
	err := bc.Hw_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lashHash = b.Get([]byte("l"))
		blockData := b.Get(lashHash)
		block := Hw_DeserializeBlock(blockData)
		lastHeight = block.Hw_Height
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	// 生成新的区块
	newBlock := Hw_NewBlock(txs, lashHash, lastHeight+1)

	// 将新区块的内容更新到数据库中
	err = bc.Hw_db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hw_Hash, newBlock.Hw_Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"), newBlock.Hw_Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.Hw_tip = newBlock.Hw_Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	return newBlock
}

// 签名
func (bc *Hw_Blockchain) Hw_SignTransaction(tx *Hw_Transaction, privKey ecdsa.PrivateKey,txs []*Hw_Transaction) {
	prevTXs := make(map[string]Hw_Transaction)

	// 找到交易输入中，之前的交易
	Vin:
	for _, vin := range tx.Hw_Vin {
		for _, tx := range txs {
			if bytes.Compare(tx.Hw_ID, vin.Hw_Txid) == 0 {
				prevTX := *tx
				prevTXs[hex.EncodeToString(prevTX.Hw_ID)] = prevTX
				continue Vin
			}
		}

		prevTX, err := bc.Hw_FindTransaction(vin.Hw_Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Hw_ID)] = prevTX

	}

	tx.Hw_Sign(privKey, prevTXs)
}

// 验证签名
func (bc *Hw_Blockchain) Hw_VerifyTransaction(tx *Hw_Transaction,txs []*Hw_Transaction) bool {
	if tx.Hw_IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]Hw_Transaction)
	Vin:
	for _, vin := range tx.Hw_Vin {
		for _, tx := range txs {
			if bytes.Compare(tx.Hw_ID, vin.Hw_Txid) == 0 {
				prevTX := *tx
				prevTXs[hex.EncodeToString(prevTX.Hw_ID)] = prevTX
				continue Vin
			}
		}
		prevTX, err := bc.Hw_FindTransaction(vin.Hw_Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Hw_ID)] = prevTX
	}

	return tx.Hw_Verify(prevTXs)
}

// 判断数据库是否存在
func Hw_dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// 获取BestHeight
func (bc *Hw_Blockchain) Hw_GetBestHeight() int {
	var lastBlock Hw_Block

	err := bc.Hw_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		lastBlock = *Hw_DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return lastBlock.Hw_Height
}

// 获取所有区块的hash
func (bc *Hw_Blockchain) Hw_GetBlockHashes() [][]byte {
	var blocks [][]byte
	bci := bc.Hw_Iterator()

	for {
		block := bci.Hw_Next()

		blocks = append(blocks, block.Hw_Hash)

		if len(block.Hw_PrevBlockHash) == 0 {
			break
		}
	}

	return blocks
}

// 根据hash获取某个区块的内容
func (bc *Hw_Blockchain) Hw_GetBlock(blockHash []byte) (Hw_Block, error) {
	var block Hw_Block

	err := bc.Hw_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		blockData := b.Get(blockHash)

		if blockData == nil {
			return errors.New("未找到区块")
		}

		block = *Hw_DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		return block, err
	}

	return block, nil
}

// 将区块添加到链中
func (bc *Hw_Blockchain) Hw_AddBlock(block *Hw_Block) {
	err := bc.Hw_db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		blockInDb := b.Get(block.Hw_Hash)

		if blockInDb != nil {
			return nil
		}

		blockData := block.Hw_Serialize()
		err := b.Put(block.Hw_Hash, blockData)
		if err != nil {
			log.Panic(err)
		}

		lastHash := b.Get([]byte("l"))
		lastBlockData := b.Get(lastHash)
		lastBlock := Hw_DeserializeBlock(lastBlockData)

		if block.Hw_Height > lastBlock.Hw_Height {
			err = b.Put([]byte("l"), block.Hw_Hash)
			if err != nil {
				log.Panic(err)
			}
			bc.Hw_tip = block.Hw_Hash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}