package BLC

import (
	"fmt"
	"os"
	"flag"
	"log"
)

type HW_CLI struct{}

// 打印使用说明
func (cli *HW_CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createwallet - 创建钱包")
	fmt.Println("  listaddresses - 打印钱包地址")
	fmt.Println("  createblockchain -address ADDRESS - 创建区块链")
	fmt.Println("  getbalance -address ADDRESS - 获取地址的余额")
	fmt.Println("  getbalanceall - 打印钱包地址的余额")
	fmt.Println("  printchain - 打印区块链中的所有区块数据")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT 转账")
	fmt.Println("  reindexutxo - 重建UTXO set")
	fmt.Println("  printutxo - 打印所有的UTXO set")

}

// 验证参数
func (cli *HW_CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli HW_CLI) HW_Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printUTXOCmd := flag.NewFlagSet("printutxo", flag.ExitOnError)
	getBalanceAllCmd := flag.NewFlagSet("getbalanceall", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "查询余额地址")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "创建创世区块地址")
	sendFrom := sendCmd.String("from", "", "转出账地址")
	sendTo := sendCmd.String("to", "", "转到账地址")
	sendAmount := sendCmd.String("amount", "", "转账金额")

	var err error
	switch os.Args[1] {
	case "getbalance":
		err = getBalanceCmd.Parse(os.Args[2:])
	case "createblockchain":
		err = createBlockchainCmd.Parse(os.Args[2:])
	case "createwallet":
		err = createWalletCmd.Parse(os.Args[2:])
	case "listaddresses":
		err = listAddressesCmd.Parse(os.Args[2:])
	case "printchain":
		err = printChainCmd.Parse(os.Args[2:])
	case "reindexutxo":
		err = reindexUTXOCmd.Parse(os.Args[2:])
	case "send":
		err = sendCmd.Parse(os.Args[2:])
	case "printutxo":
		err = printUTXOCmd.Parse(os.Args[2:])
	case "getbalanceall":
		err = getBalanceAllCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if err != nil {
		log.Panic(err)
	}

	if createWalletCmd.Parsed() {
		cli.HW_createWallet()
	}

	if listAddressesCmd.Parsed() {
		cli.HW_listAddrsss()
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.HW_createblockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.HW_printchain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount == "" {
			sendCmd.Usage()
			os.Exit(1)
		}

		// 检查参数，有效性
		from := JSONToArray(*sendFrom)
		to := JSONToArray(*sendTo)

		for index, fromAdress := range from {
			if !HW_ValidateAddress(fromAdress) || !HW_ValidateAddress(to[index]) {
				fmt.Println("地址无效。。")
				os.Exit(1)
			}
		}
		amount := JSONToArray(*sendAmount)
		cli.HW_send(from, to, amount)
	}
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.HW_getBalance(*getBalanceAddress)
	}
	if reindexUTXOCmd.Parsed() {
		cli.HW_reindexUTXO()
	}

	if printUTXOCmd.Parsed() {
		cli.HW_printutxo()
	}
	if getBalanceAllCmd.Parsed() {
		cli.HW_getBalanceAll()
	}
}
