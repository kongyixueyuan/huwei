package BLC

import "bytes"

type Hw_TXInput struct {
	Hw_Txid      []byte
	Hw_Vout      int      // Vout的index
	Hw_Signature []byte   // 签名
	Hw_PubKey    []byte   // 公钥
}

func (in Hw_TXInput) UsesKey(pubKeyHash []byte) bool  {
	lockingHash := Hw_HashPubKey(in.Hw_PubKey)

	return bytes.Compare(lockingHash,pubKeyHash) == 0
}
