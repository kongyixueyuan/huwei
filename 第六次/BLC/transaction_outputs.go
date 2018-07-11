package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type HW_TXOutputs struct {
	Hw_Outputs []HW_TXOutput
}

//  序列化 TXOutputs
func (outs HW_TXOutputs) Hw_Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// 反序列化 TXOutputs
func Rwq_DeserializeOutputs(data []byte) HW_TXOutputs {
	var outputs HW_TXOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}
