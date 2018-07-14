package BLC

import "fmt"

func (cli *Hw_CLI) Hw_createWallet(nodeID string) {
	//wallet := Hw_NewWallet()
	//address := wallet.Hw_GetAddress()
	//fmt.Printf("钱包地址：%s\n",address)

	wallets, _ := Hw_NewWallets(nodeID)
	address := wallets.Hw_NewWallet().Hw_GetAddress()
	wallets.Hw_SaveToFile(nodeID)
	fmt.Printf("钱包地址：%s\n", address)

}
