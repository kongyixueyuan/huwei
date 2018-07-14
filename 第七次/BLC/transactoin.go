package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
	"strings"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/elliptic"
	"math/big"
)

// 创世区块，Token数量
const subsidy  = 10

type Hw_Transaction struct {
	Hw_ID   []byte
	Hw_Vin  []Hw_TXInput
	Hw_Vout []Hw_TXOutput
}

// 是否是创世区块交易
func (tx Hw_Transaction) Hw_IsCoinbase() bool {
	// Vin 只有一条
	// Vin 第一条数据的Txid 为 0
	// Vin 第一条数据的Vout 为 -1
	return len(tx.Hw_Vin) == 1 && len(tx.Hw_Vin[0].Hw_Txid) == 0 && tx.Hw_Vin[0].Hw_Vout == -1
}


// 将交易进行Hash
func (tx *Hw_Transaction) Hw_Hash() []byte  {
	var hash [32]byte

	txCopy := *tx
	txCopy.Hw_ID = []byte{}

	hash = sha256.Sum256(txCopy.Hw_Serialize())
	return hash[:]
}
// 新建创世区块的交易
func Hw_NewCoinbaseTX(to ,data string) *Hw_Transaction  {
	if data == ""{
		//如果数据为空，可以随机给默认数据,用于挖矿奖励
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}

		data = fmt.Sprintf("%x", randData)
	}
	txin := Hw_TXInput{[]byte{},-1,nil,[]byte(data)}
	txout := Hw_NewTXOutput(subsidy,to)

	tx := Hw_Transaction{nil,[]Hw_TXInput{txin},[]Hw_TXOutput{*txout}}
	tx.Hw_ID = tx.Hw_Hash()
	return &tx
}

// 转帐时生成交易
func Hw_NewUTXOTransaction(wallet *Hw_Wallet,to string,amount int,UTXOSet *Hw_UTXOSet,txs []*Hw_Transaction) *Hw_Transaction   {

	// 如果本区块中，多笔转账
	/**
	第一种情况：
	  A:10
	  A->B 2
	  A->C 4

	  tx1:
	      Vin:
	           ATxID  out ...
	      Vout:
	           A : 8
	           B : 2
	  tx1:
	      Vin:
	           ATxID  out ...
	      Vout:
	           A : 4
	           C : 4
	第二种情况：
	  A:10+10
	  A->B 4
	  A->C 8
	**/

	pubKeyHash := Hw_HashPubKey(wallet.Hw_PublicKey)
	if len(txs) > 0 {
		// 查的txs中的UTXO
		utxo := Hw_FindUTXOFromTransactions(txs)

		// 找出当前钱包已经花费的
		unspentOutputs := make(map[string][]int)
		acc := 0
		for txID,outs := range utxo {
			for outIdx, out := range outs.Hw_Outputs {
				if out.Hw_IsLockedWithKey(pubKeyHash) && acc < amount {
					acc += out.Hw_Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}

		if acc >= amount { // 当前交易中的剩余余额可以支付
			fmt.Println("txs>0 && acc >= amount")
			return Hw_NewUTXOTransactionEnd(wallet,to,amount,UTXOSet,acc,unspentOutputs,txs)
		}else{
			fmt.Println("txs>0 && acc < amount")
			accLeft, validOutputs := UTXOSet.Hw_FindSpendableOutputs(pubKeyHash,  amount - acc)
			for k,v := range unspentOutputs{
				validOutputs[k] = v
			}
			return Hw_NewUTXOTransactionEnd(wallet,to,amount,UTXOSet,acc + accLeft,validOutputs,txs)
		}
	} else { //只是当前一笔交易
		fmt.Println("txs==0")
		acc, validOutputs := UTXOSet.Hw_FindSpendableOutputs(pubKeyHash, amount)

		return Hw_NewUTXOTransactionEnd(wallet,to,amount,UTXOSet,acc,validOutputs,txs)
	}
}

func Hw_NewUTXOTransactionEnd(wallet *Hw_Wallet,to string,amount int,UTXOSet *Hw_UTXOSet,acc int,UTXO map[string][]int,txs []*Hw_Transaction) *Hw_Transaction {

	if acc < amount {
		log.Panic("账户余额不足")
	}

	var inputs []Hw_TXInput
	var outputs []Hw_TXOutput
	// 构造input
	for txid, outs := range UTXO {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := Hw_TXInput{txID, out, nil, wallet.Hw_PublicKey}
			inputs = append(inputs, input)
		}
	}
	// 生成交易输出
	outputs = append(outputs, *Hw_NewTXOutput(amount, to))
	// 生成余额
	if acc > amount {
		outputs = append(outputs, *Hw_NewTXOutput(acc-amount, string(wallet.Hw_GetAddress())))
	}

	tx := Hw_Transaction{nil, inputs, outputs}
	tx.Hw_ID = tx.Hw_Hash()
	// 签名

	//tx.String()
	UTXOSet.Hw_Blockchain.Hw_SignTransaction(&tx, wallet.Hw_PrivateKey,txs)

	return &tx
}


// 找出交易中的utxo
func Hw_FindUTXOFromTransactions(txs []*Hw_Transaction) map[string]Hw_TXOutputs {
	UTXO := make(map[string]Hw_TXOutputs)
	// 已经花费的交易txID : TXOutputs.index
	spentTXOs := make(map[string][]int)
	// 循环区块中的交易
	for _, tx := range txs {
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
	return UTXO

}

// 签名
func (tx *Hw_Transaction) Hw_Sign(privateKey ecdsa.PrivateKey,prevTXs map[string]Hw_Transaction)  {
	if tx.Hw_IsCoinbase() { // 创世区块不需要签名
		return
	}

	// 检查交易的输入是否正确
	for _,vin := range tx.Hw_Vin{
		if prevTXs[hex.EncodeToString(vin.Hw_Txid)].Hw_ID == nil{
			log.Panic("错误：之前的交易不正确")
		}
	}

	txCopy := tx.Hw_TrimmedCopy()

	for inID, vin := range txCopy.Hw_Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Hw_Txid)]
		txCopy.Hw_Vin[inID].Hw_Signature = nil
		txCopy.Hw_Vin[inID].Hw_PubKey = prevTx.Hw_Vout[vin.Hw_Vout].Hw_PubKeyHash

		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, []byte(dataToSign))
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Hw_Vin[inID].Hw_Signature = signature
		txCopy.Hw_Vin[inID].Hw_PubKey = nil
	}
}
// 验证签名
func (tx *Hw_Transaction) Hw_Verify(prevTXs map[string]Hw_Transaction) bool {
	if tx.Hw_IsCoinbase() {
		return true
	}

	for _, vin := range tx.Hw_Vin {
		if prevTXs[hex.EncodeToString(vin.Hw_Txid)].Hw_ID == nil {
			log.Panic("错误：之前的交易不正确")
		}
	}

	txCopy := tx.Hw_TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Hw_Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Hw_Txid)]
		txCopy.Hw_Vin[inID].Hw_Signature = nil
		txCopy.Hw_Vin[inID].Hw_PubKey = prevTx.Hw_Vout[vin.Hw_Vout].Hw_PubKeyHash

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Hw_Signature)
		r.SetBytes(vin.Hw_Signature[:(sigLen / 2)])
		s.SetBytes(vin.Hw_Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.Hw_PubKey)
		x.SetBytes(vin.Hw_PubKey[:(keyLen / 2)])
		y.SetBytes(vin.Hw_PubKey[(keyLen / 2):])

		dataToVerify := fmt.Sprintf("%x\n", txCopy)

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
			return false
		}
		txCopy.Hw_Vin[inID].Hw_PubKey = nil
	}

	return true
}

// 复制交易（输入的签名和公钥置为空）
func (tx *Hw_Transaction) Hw_TrimmedCopy() Hw_Transaction {
	var inputs []Hw_TXInput
	var outputs []Hw_TXOutput

	for _, vin := range tx.Hw_Vin {
		inputs = append(inputs, Hw_TXInput{vin.Hw_Txid, vin.Hw_Vout, nil, nil})
	}

	for _, vout := range tx.Hw_Vout {
		outputs = append(outputs, Hw_TXOutput{vout.Hw_Value, vout.Hw_PubKeyHash})
	}

	txCopy := Hw_Transaction{tx.Hw_ID, inputs, outputs}

	return txCopy
}
// 打印交易内容
func (tx Hw_Transaction) String()  {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction ID: %x", tx.Hw_ID))

	for i, input := range tx.Hw_Vin {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Hw_Txid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Hw_Vout))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Hw_Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.Hw_PubKey))
	}

	for i, output := range tx.Hw_Vout {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Hw_Value))
		lines = append(lines, fmt.Sprintf("       PubKeyHash: %x", output.Hw_PubKeyHash))
	}
	fmt.Println(strings.Join(lines, "\n"))
}


// 将交易序列化
func (tx Hw_Transaction) Hw_Serialize() []byte  {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)

	if err != nil{
		log.Panic(err)
	}
	return encoded.Bytes()
}
// 反序列化交易
func Hw_DeserializeTransaction(data []byte) Hw_Transaction {
	var transaction Hw_Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}

	return transaction
}

// 将交易数组序列化
func Hw_SerializeTransactions(txs []*Hw_Transaction) [][]byte  {

	var txsHash [][]byte
	for _,tx := range txs{
		txsHash = append(txsHash, tx.Hw_Serialize())
	}
	return txsHash
}

// 反序列化交易数组
func Hw_DeserializeTransactions(data [][]byte) []Hw_Transaction {
	var txs []Hw_Transaction
	for _,tx := range data {
		var transaction Hw_Transaction
		decoder := gob.NewDecoder(bytes.NewReader(tx))
		err := decoder.Decode(&transaction)
		if err != nil {
			log.Panic(err)
		}

		txs = append(txs, transaction)
	}
	return txs
}
