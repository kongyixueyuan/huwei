package BLC

import "bytes"

type HW_TXInput struct {
	Hw_Txid      []byte
	Hw_Vout      int    // Vout的index
	Hw_Signature []byte // 签名
	Hw_PubKey    []byte // 公钥
}

func (in HW_TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HW_HashPubKey(in.Hw_PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
