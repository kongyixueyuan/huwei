package BLC

import (
	"log"
	"fmt"
)

func (cli *HW_CLI) HW_listAddrsss() {
	wallets, err := HW_NewWallets()

	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.HW_GetAddresses()

	for _, address := range addresses {
		fmt.Printf("%s\n", address)
	}
}
