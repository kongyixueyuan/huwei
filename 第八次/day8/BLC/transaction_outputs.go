package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Hw_TXOutputs struct {
	Hw_Outputs []Hw_TXOutput
}

//  序列化 TXOutputs
func (outs Hw_TXOutputs) Hw_Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// 反序列化 TXOutputs
func Hw_DeserializeOutputs(data []byte) Hw_TXOutputs {
	var outputs Hw_TXOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}
