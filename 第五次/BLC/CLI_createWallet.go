package BLC

import "fmt"

func (cli *CLI) createWallet() {
	wallets,_ := NewWallets()
	wallets.CreateNewWallet()
	wallets.SaveWallets()
	fmt.Println(wallets.WalletsMap)
}
