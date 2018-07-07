package BLC

import (
	"bytes"
)

type TXInput struct {
	// 交易的hash
	Txhash []byte
	// 存储TXOutput在Vout里面的索引
	Vout      int
	Signature []byte // 数字签名
	PubKey    []byte //公钥，钱包里面
}

// 判断当前的消费是否是该地址
func (txInput *TXInput) UnLockWithAddress(ripemd160Hash []byte) bool {
	publicKey := Ripemd160Hash(txInput.PubKey)
	return bytes.Compare(publicKey, ripemd160Hash) == 0

}
