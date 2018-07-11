package BLC

import (
	"os"
	"io/ioutil"
	"log"
	"encoding/gob"
	"crypto/elliptic"
	"bytes"
)

const walletFile = "wallet.dat"

type HW_Wallets struct {
	HW_Wallets map[string]*HW_Wallet
}

// 生成新的钱包
// 从数据库中读取，如果不存在
func HW_NewWallets() (*HW_Wallets, error) {
	wallets := HW_Wallets{}
	wallets.HW_Wallets = make(map[string]*HW_Wallet)

	err := wallets.HW_LoadFromFile()

	return &wallets, err
}

// 生成新的钱包地址列表
func (ws *HW_Wallets) HW_NewWallet() *HW_Wallet {
	wallet := HW_NewWallet()
	address := wallet.HW_GetAddress()
	ws.HW_Wallets[string(address)] = wallet
	return wallet
}

// 获取钱包地址
func (ws *HW_Wallets) HW_GetAddresses() []string {
	var addresses []string
	for address := range ws.HW_Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

// 根据地址获取钱包的详细信息
func (ws HW_Wallets) HW_GetWallet(address string) HW_Wallet {
	return *ws.HW_Wallets[address]
}

// 从数据库中读取钱包列表
func (ws *HW_Wallets) HW_LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets HW_Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	ws.HW_Wallets = wallets.HW_Wallets

	return nil
}

// 将钱包存到数据库中
func (ws *HW_Wallets) HW_SaveToFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

// 打印所有钱包的余额
func (ws *HW_Wallets) HW_GetBalanceAll() map[string]int {
	addresses := ws.HW_GetAddresses()
	bc := HW_NewBlockchain()
	defer bc.hw_db.Close()
	UTXOSet := HW_UTXOSet{bc}

	result := make(map[string]int)
	for _, address := range addresses {
		if !HW_ValidateAddress(address) {
			result[address] = -1
		}
		balance := UTXOSet.HW_GetBalance(address)
		result[address] = balance
	}
	return result
}
