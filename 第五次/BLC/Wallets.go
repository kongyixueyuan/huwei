package BLC

import (
	"fmt"
	"os"
	"io/ioutil"
	"log"
	"encoding/gob"
	"crypto/elliptic"
	"bytes"
)

const walletFile = "Wallets.dat"

type Wallets struct {
	WalletsMap map[string]*Wallet
}

// 创建钱包集合
func NewWallets() (*Wallets, error) {
	//如果钱包文件不存在
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &Wallets{}
		wallets.WalletsMap = make(map[string]*Wallet)
		return wallets, err
	}
	//读取文件
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}
	var wallets Wallets
	//反序列化
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	return &wallets, err
}

// 创建新钱包
func (w *Wallets) CreateNewWallet() {
	wallet := NewWallet()
	fmt.Printf("钱包Address为: %s\n", wallet.GetAddress())
	w.WalletsMap[string(wallet.GetAddress())] = wallet
}

// 将钱包信息写入到文件
func (w *Wallets) SaveWallets() {
	var content bytes.Buffer

	// 注册的目的是为了，可以序列化所有的类型
	gob.Register(elliptic.P256())

	// 系列化钱包
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&w)
	if err != nil {
		log.Panic(err)
	}
	// 将序列化的钱包，写入到文件
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
