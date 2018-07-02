package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
)

//int64 转成[]byte
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if (err != nil) {
		log.Panic(err)
	}
	return buff.Bytes()

}

// 标准的jsonString转成数组
func JSONToArray(jsonString string) []string {

	//json 到 []string
	var sArr []string
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil {
		log.Panic(err)
	}

	return sArr
}
