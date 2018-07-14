package BLC

import (
	"os"
	"io/ioutil"
	"log"
	"encoding/gob"
	"crypto/elliptic"
	"bytes"
	"fmt"
)

const walletFile  = "wallet_%s.dat"

type Hw_Wallets struct {
	Hw_Wallets map[string]*Hw_Wallet
}

// 生成新的钱包
// 从数据库中读取，如果不存在
func Hw_NewWallets(nodeID string)(*Hw_Wallets,error)  {
	wallets := Hw_Wallets{}
	wallets.Hw_Wallets = make(map[string]*Hw_Wallet)

	err := wallets.Hw_LoadFromFile(nodeID)

	return &wallets,err
}
// 生成新的钱包地址列表
func (ws *Hw_Wallets) Hw_NewWallet() *Hw_Wallet {
	wallet := Hw_NewWallet()
	address := wallet.Hw_GetAddress()
	ws.Hw_Wallets[string(address)] = wallet
	return wallet
}
// 获取钱包地址
func (ws *Hw_Wallets) Hw_GetAddresses()[]string  {
	var addresses []string
	for address := range ws.Hw_Wallets{
		addresses = append(addresses,address)
	}
	return addresses
}

// 根据地址获取钱包的详细信息
func (ws Hw_Wallets) Hw_GetWallet(address string) Hw_Wallet {
	return *ws.Hw_Wallets[address]
}

// 从数据库中读取钱包列表
func (ws *Hw_Wallets)Hw_LoadFromFile(nodeID string) error  {
	 walletFile := fmt.Sprintf(walletFile, nodeID)
	 if _,err := os.Stat(walletFile) ; os.IsNotExist(err){
	 	return err
	 }

	 fileContent ,err := ioutil.ReadFile(walletFile)
	 if err !=nil{
	 	log.Panic(err)
	 }

	 var wallets Hw_Wallets
	 gob.Register(elliptic.P256())
	 decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	 err = decoder.Decode(&wallets)
	 if err !=nil{
	 	log.Panic(err)
	 }

	 ws.Hw_Wallets = wallets.Hw_Wallets

	 return nil
}

// 将钱包存到数据库中
func (ws *Hw_Wallets)Hw_SaveToFile(nodeID string)  {
	walletFile := fmt.Sprintf(walletFile, nodeID)
	var content bytes.Buffer

	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err !=nil{
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile,content.Bytes(),0644)
	if err !=nil{
		log.Panic(err)
	}
}
// 打印所有钱包的余额
func (ws *Hw_Wallets) Hw_GetBalanceAll(nodeID string) map[string]int {
	addresses := ws.Hw_GetAddresses()
	bc := Hw_NewBlockchain(nodeID)
	defer bc.Hw_db.Close()
	UTXOSet := Hw_UTXOSet{bc}

	result := make(map[string]int)
	for _,address := range addresses{
		if !Hw_ValidateAddress(address) {
			result[address] = -1
		}
		balance := UTXOSet.Hw_GetBalance(address)
		result[address] = balance
	}
	return result
}