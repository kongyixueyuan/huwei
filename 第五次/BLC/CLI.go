package BLC

import (
	"os"
	"flag"
	"log"
	"fmt"
)

type CLI struct {
	blockchain *Blockchain
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\taddresslists  -- 输出所有钱包地址.")
	fmt.Println("\tcreatewallet  -- 创建钱包.")
	fmt.Println("\tcreateblockchain -address -- 创建创世区块---交易数据.")
	fmt.Println("\tsend -from FROM -to TO -amount value --发起交易")
	fmt.Println("\tprintchain -- 输出区块信息.")
	fmt.Println("\tgetBalance -address -- 查询余额")
}
func IsvalidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

}
func (cli *CLI) createGenesisBlockChain(data string) {
	blockchain := CreateBlockchainWithGenenisBlock(data,1000)
	defer blockchain.DB.Close()
}
func (cli *CLI) Run() {
	IsvalidArgs()
	addresslistsCmd := flag.NewFlagSet("addresslists", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendBlockcmd := flag.NewFlagSet("send", flag.ExitOnError)
	getBalancecmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	flagFrom := sendBlockcmd.String("from", "", "转账源地址")
	flagTo := sendBlockcmd.String("to", "", "转账目的地址")
	flagAmount := sendBlockcmd.String("amount", "", "转账金额")
	//创建初始区块并生成初始地址
	flagCreateBlockchainWithAddress := createBlockchainCmd.String("address", "", "genesis data")
	getBalanceWithAddress := getBalancecmd.String("address", "", "查询某个地址对应的余额")

	switch os.Args[1] {
	case "send":
		err := sendBlockcmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "addresslists":
		err := addresslistsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getBalance":
		err := getBalancecmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}
	if sendBlockcmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			printUsage()
			os.Exit(1)
		}
		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		/*for _, fromAddress := range from {
			if Wallet.IsValidForAddress(fromAddress)
		}*/
		amount := JSONToArray(*flagAmount)
		cli.Send(from, to, amount)
	}
	if printChainCmd.Parsed() {
		cli.Printchain()
	}
	if addresslistsCmd.Parsed() {
		cli.addressLists()
	}
	if createWalletCmd.Parsed() {
		//创建钱包
		cli.createWallet()
	}
	if createBlockchainCmd.Parsed() {
		if *flagCreateBlockchainWithAddress == "" {
			fmt.Println("交易数据不能为空---地址不能为空。。。。。。")
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockChain(*flagCreateBlockchainWithAddress)
	}
	if getBalancecmd.Parsed() {
		if *getBalanceWithAddress == "" {
			fmt.Println("地址不能为空")
			printUsage()
			os.Exit(1)
		}
		cli.GetBalance(*getBalanceWithAddress)
	}
}

func (cli *CLI) Printchain() {

	if DBExists() == false {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}

	blockchain := BlockchainObject()

	defer blockchain.DB.Close()

	blockchain.Printchain()
}

// 转账
func (cli *CLI) Send(from []string, to []string, amount []string) {

	if DBExists() == false {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	defer blockchain.DB.Close()
	blockchain.MineNewBlock(from, to, amount)

}

// 查询余额
func (cli *CLI) GetBalance(address string) {

	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	amount := blockchain.GetBalance(address)

	fmt.Printf("%s一共有%d个Token\n", address, amount)

}
