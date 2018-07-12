package BLC

import "fmt"

func (cli *HW_CLI) HW_createWallet() {

	wallets, _ := HW_NewWallets()
	address := wallets.HW_NewWallet().HW_GetAddress()
	wallets.HW_SaveToFile()
	fmt.Printf("钱包地址：%s\n", address)

}
