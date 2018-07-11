package BLC

import (
	"log"
	"fmt"
)

func (cli *Rwq_CLI) HW_getBalance(address string) {
	if !HW_ValidateAddress(address) {
		log.Panic("错误：地址无效")
	}

	bc := HW_NewBlockchain()
	defer bc.hw_db.Close()
	UTXOSet := HW_UTXOSet{bc}

	balance := UTXOSet.HW_GetBalance(address)
	fmt.Printf("地址:%s的余额为：%d\n", address, balance)
}

func (cli *Rwq_CLI) rwq_getBalanceAll() {
	wallets, err := HW_NewWallets()
	if err != nil {
		log.Panic(err)
	}
	balances := wallets.HW_GetBalanceAll()
	for address, balance := range balances {
		fmt.Printf("地址:%s的余额为：%d\n", address, balance)
	}
}
