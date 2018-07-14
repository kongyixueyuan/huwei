package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"encoding/hex"
	"fmt"
	"strings"
)

const utxoBucket = "chainstate"

type Hw_UTXOSet struct {
	Hw_Blockchain *Hw_Blockchain
}

// 查询可花费的交易输出
func (u Hw_UTXOSet) Hw_FindSpendableOutputs(pubkeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Hw_Blockchain.Hw_db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := Hw_DeserializeOutputs(v)

			for outIdx, out := range outs.Hw_Outputs {
				if out.Hw_IsLockedWithKey(pubkeyHash) && accumulated < amount {
					accumulated += out.Hw_Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return accumulated, unspentOutputs
}

func (u Hw_UTXOSet) Hw_Reindex() {
	db := u.Hw_Blockchain.Hw_db
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		// 删除旧的bucket
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic()
		}
		_, err = tx.CreateBucket(bucketName)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	UTXO := u.Hw_Blockchain.FindUTXO()

	err = db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucketName)

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(key, outs.Hw_Serialize())
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
}

// 生成新区块的时候，更新UTXO数据库
func (u Hw_UTXOSet) Update(block *Hw_Block) {
	err := u.Hw_Blockchain.Hw_db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		for _, tx := range block.Hw_Transactions {
			if !tx.Hw_IsCoinbase() {
				for _, vin := range tx.Hw_Vin {
					updatedOuts := Hw_TXOutputs{}
					outsBytes := b.Get(vin.Hw_Txid)
					outs := Hw_DeserializeOutputs(outsBytes)

					// 找出Vin对应的outputs,过滤掉花费的
					for outIndex, out := range outs.Hw_Outputs {
						if outIndex != vin.Hw_Vout {
							updatedOuts.Hw_Outputs = append(updatedOuts.Hw_Outputs, out)
						}
					}
					// 未花费的交易输出TXOutput为0
					if len(updatedOuts.Hw_Outputs) == 0 {
						err := b.Delete(vin.Hw_Txid)
						if err != nil {
							log.Panic(err)
						}
					} else { // 未花费的交易输出TXOutput>0
						err := b.Put(vin.Hw_Txid, updatedOuts.Hw_Serialize())
						if err != nil {
							log.Panic(err)
						}
					}
				}
			}

			// 将所有的交易输出TXOutput存入数据库中
			newOutputs := Hw_TXOutputs{}
			for _, out := range tx.Hw_Vout {
				newOutputs.Hw_Outputs = append(newOutputs.Hw_Outputs, out)
			}
			err := b.Put(tx.Hw_ID, newOutputs.Hw_Serialize())
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// 打出某个公钥hash对应的所有未花费输出
func (u *Hw_UTXOSet) Hw_FindUTXO(pubKeyHash []byte) []Hw_TXOutput {
	var UTXOs []Hw_TXOutput

	err := u.Hw_Blockchain.Hw_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := Hw_DeserializeOutputs(v)

			for _, out := range outs.Hw_Outputs {
				if out.Hw_IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

// 查询某个地址的余额
func (u *Hw_UTXOSet) Hw_GetBalance(address string) int {
	balance := 0
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := u.Hw_FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Hw_Value
	}
	return balance
}

// 打印所有的UTXO
func (u *Hw_UTXOSet) String() {
	//outputs := make(map[string][]Hw_TXOutput)

	var lines []string
	lines = append(lines, "---ALL UTXO:")
	err := u.Hw_Blockchain.Hw_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := Hw_DeserializeOutputs(v)

			lines = append(lines, fmt.Sprintf("     Key: %s", txID))
			for i, out := range outs.Hw_Outputs {
				//outputs[txID] = append(outputs[txID], out)
				lines = append(lines, fmt.Sprintf("     Output: %d", i))
				lines = append(lines, fmt.Sprintf("         value:  %d", out.Hw_Value))
				lines = append(lines, fmt.Sprintf("         PubKeyHash:  %x", out.Hw_PubKeyHash))
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(strings.Join(lines, "\n"))
}
