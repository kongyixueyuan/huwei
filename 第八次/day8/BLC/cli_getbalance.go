package BLC

import (
	"log"
	"fmt"
)

func (cli *Hw_CLI) Hw_getBalance(address string,nodeID string) {
	if !Hw_ValidateAddress(address) {
		log.Panic("错误：地址无效")
	}

	bc := Hw_NewBlockchain(nodeID)
	defer bc.Hw_db.Close()
	UTXOSet := Hw_UTXOSet{bc}

	balance := UTXOSet.Hw_GetBalance(address)
	fmt.Printf("地址:%s的余额为：%d\n", address, balance)
}

func (cli *Hw_CLI) Hw_getBalanceAll(nodeID string) {
	wallets,err := Hw_NewWallets(nodeID)
	if err!=nil{
		log.Panic(err)
	}
	balances := wallets.Hw_GetBalanceAll(nodeID)
	for address,balance := range balances{
		fmt.Printf("地址:%s的余额为：%d\n", address, balance)
	}
}