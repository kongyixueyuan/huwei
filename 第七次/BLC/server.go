package BLC

import (
	"fmt"
	"net"
	"log"
	"bytes"
	"encoding/gob"
	"io"
	"io/ioutil"
	"encoding/hex"
)

const protocol = "tcp"   // 节点协议
const nodeVersion = 1    // 节点版本
const commandLength = 12 // 命令长度

var nodeAddress string                         // 当前节点地址
var miningAddress string                       // 挖矿地址
var knownNodes = []string{"localhost:3000"}    // 存储所有已知节点
var blocksInTransit = [][]byte{}               // 存储接收到的区块hash
var mempool = make(map[string]Hw_Transaction) // 存储接收到的交易

type Hw_addr struct {
	Hw_AddrList []string
}

type Hw_block struct {
	Hw_AddrFrom string
	Hw_Block    []byte
}

type Hw_getblocks struct {
	Hw_AddrFrom string
}

type Hw_getdata struct {
	Hw_AddrFrom string
	Hw_Type     string
	Hw_ID       []byte
}

type Hw_inv struct {
	Hw_AddrFrom string
	Hw_Type     string
	Hw_Items    [][]byte
}

type Hw_txs struct {
	Hw_AddFrom     string
	Hw_Transactions [][]byte
}

type Hw_version struct {
	Hw_Version    int
	Hw_BestHeight int
	Hw_AddrFrom   string
}

//启动Server
func Hw_StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress
	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := Hw_NewBlockchain(nodeID)

	// 3000端口为：主节点
	// 3001端口为：钱包节点
	// 3002端口为：挖矿节点
	if nodeAddress != knownNodes[0] {
		// 此节点是钱包节点或者矿工节点，需要向主节点发送请求同步数据
		Hw_sendVersion(knownNodes[0], bc)
	}

	for { // 等待接收命令
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go Hw_handleConnecton(conn, bc)
	}
}

// 向中心节点发送 version 消息来查询是否自己的区块链已过时
func Hw_sendVersion(addr string, bc *Hw_Blockchain) {
	bestHeight := bc.Hw_GetBestHeight()
	payload := Hw_gobEncode(Hw_version{nodeVersion, bestHeight, nodeAddress})

	request := append(Hw_commandToBytes("version"), payload...)

	Hw_sendData(addr, request)
}

// 发送数据
func Hw_sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	// 如果连接失败，更新节点数据
	if err != nil {
		fmt.Sprintf("%s地址不可用\n", addr)
		var updatedNodes []string

		for _, node := range knownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		knownNodes = updatedNodes
		return
	}
	defer conn.Close()
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}

}

// 发送获取区块的的命令
func Hw_sendGetBlocks(address string) {
	payload := Hw_gobEncode(Hw_getblocks{nodeAddress})
	request := append(Hw_commandToBytes("getblocks"), payload...)

	Hw_sendData(address, request)
}

// 发送处理区块及交易hash列表请求
func Hw_sendInv(address, kind string, items [][]byte) {
	inventory := Hw_inv{nodeAddress, kind, items}
	payload := Hw_gobEncode(inventory)
	request := append(Hw_commandToBytes("inv"), payload...)

	Hw_sendData(address, request)
}

// 发送获取区块详细数据的命令
func Hw_sendGetData(address, kind string, id []byte) {
	payload := Hw_gobEncode(Hw_getdata{nodeAddress, kind, id})
	request := append(Hw_commandToBytes("getdata"), payload...)

	Hw_sendData(address, request)
}

// 发送区块内容命令
func Hw_sendBlock(addr string, b *Hw_Block) {
	data := Hw_block{nodeAddress, b.Hw_Serialize()}
	payload := Hw_gobEncode(data)
	request := append(Hw_commandToBytes("block"), payload...)

	Hw_sendData(addr, request)
}

// 发送交易内容命令
func Hw_sendTx(addr string, tx *Hw_Transaction) {
	txs := []*Hw_Transaction{tx}
	Hw_sendTxs(addr,txs)
}
// 发送交易内容命令
func Hw_sendTxs(addr string, txs []*Hw_Transaction) {

	data := Hw_txs{nodeAddress, Hw_SerializeTransactions(txs)}
	payload := Hw_gobEncode(data)
	request := append(Hw_commandToBytes("tx"), payload...)

	Hw_sendData(addr, request)
}

//================================================================
// 命令收集并分发
func Hw_handleConnecton(conn net.Conn, bc *Hw_Blockchain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := Hw_bytesToCommand(request[:commandLength])
	fmt.Printf("接收到命令：%s\n", command)

	switch command {
	case "addr": // 添加新节点
		Hw_handleAddr(request)
	case "block": // 添加新区块
		Hw_handleBlock(request, bc)
	case "inv": // 接收区块及交易hash列表 ，区块获取到内容到存储到 blocksInTransit ， 交易存储到 mempool
		Hw_handleInv(request, bc)
	case "getblocks": // 将当前节点区块链中的所有区块hash，返回给请求节点
		Hw_handleGetBlocks(request, bc)
	case "getdata": // 将单个交易或区块的内容 返回给请求节点
		Hw_handleGetData(request, bc)
	case "tx": // 添加新的交易,交易数量大于2，矿工节点挖矿,如果是主节点，进行分发交易
		Hw_handleTx(request, bc)
	case "version": // 检查是否需要同步数据，根据区块的height
		Hw_handleVersion(request, bc)
	default:
		fmt.Println("未知命令!")
	}

	conn.Close()

}

// 添加新节点
func Hw_handleAddr(request []byte) {
	var buff bytes.Buffer
	var payload Hw_addr

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	knownNodes = append(knownNodes, payload.Hw_AddrList...)
	fmt.Printf("有%d个节点加入\n", len(knownNodes))
	// 把新节点推送给其他所有节点
	for _, node := range knownNodes {
		Hw_sendGetBlocks(node)
	}
}

/*
当接收到一个新块时，我们把它放到区块链里面。
如果还有更多的区块需要下载，我们继续从上一个下载的块的那个节点继续请求。
当最后把所有块都下载完后，对 UTXO 集进行重新索引

TODO: 并非无条件信任，我们应该在将每个块加入到区块链之前对它们进行验证。
TODO: 并非运行 UTXOSet.Reindex()， 而是应该使用 UTXOSet.Update(block)，
TODO: 因为如果区块链很大，它将需要很多时间来对整个 UTXO 集重新索引。
 */
func Hw_handleBlock(request []byte, bc *Hw_Blockchain) {
	var buff bytes.Buffer
	var payload Hw_block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Hw_Block
	block := Hw_DeserializeBlock(blockData)

	fmt.Println("收到一个新的区块!")
	bc.Hw_AddBlock(block)

	fmt.Printf("Added block %x\n", block.Hw_Hash)

	// 如果还有更多的区块需要下载，继续从上一个下载的块的那个节点继续请求
	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		Hw_sendGetData(payload.Hw_AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := Hw_UTXOSet{bc}
		UTXOSet.Hw_Reindex()
	}
}

// 向其他节点展示当前节点有什么块和交易,接收区块及交易列表
func Hw_handleInv(request []byte, bc *Hw_Blockchain) {
	var buff bytes.Buffer
	var payload Hw_inv

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("接收到列表,长度为：%d，类型： %s\n", len(payload.Hw_Items), payload.Hw_Type)

	// 如果数据是 区块
	if payload.Hw_Type == "block" {
		blocksInTransit = payload.Hw_Items

		blockHash := payload.Hw_Items[0]
		// 发送获取区块详细数据的命令
		Hw_sendGetData(payload.Hw_AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}
	// 如果数据是 交易
	// 转账时，未立即挖矿，将交易保存到内存池中
	if payload.Hw_Type == "tx" {
		txID := payload.Hw_Items[0]
		// 如果内存池中，交易内容为空
		if mempool[hex.EncodeToString(txID)].Hw_ID == nil {
			// 发送获取交易详细内容请求
			Hw_sendGetData(payload.Hw_AddrFrom, "tx", txID)
		}
	}
}

// 处理获取区块命令，返回区块链中的所有区块hash
func Hw_handleGetBlocks(request []byte, bc *Hw_Blockchain) {
	var buff bytes.Buffer
	var payload Hw_getblocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.Hw_GetBlockHashes()
	Hw_sendInv(payload.Hw_AddrFrom, "block", blocks)
}

//  将单个交易或区块的内容 返回给请求节点
func Hw_handleGetData(request []byte, bc *Hw_Blockchain) {
	var buff bytes.Buffer
	var payload Hw_getdata

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Hw_Type == "block" {
		block, err := bc.Hw_GetBlock([]byte(payload.Hw_ID))
		if err != nil {
			return
		}

		Hw_sendBlock(payload.Hw_AddrFrom, &block)
	}

	if payload.Hw_Type == "tx" {
		txID := hex.EncodeToString(payload.Hw_ID)
		tx := mempool[txID]

		Hw_sendTx(payload.Hw_AddrFrom, &tx)
		// delete(mempool, txID)
	}
}

// 处理交易
func Hw_handleTx(request []byte, bc *Hw_Blockchain) {
	var buff bytes.Buffer
	var payload Hw_txs

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Hw_Transactions
	txsDes := Hw_DeserializeTransactions(txData)

	for _,tx := range txsDes {
		mempool[hex.EncodeToString(tx.Hw_ID)] = tx

		// 如果是主节点
		if nodeAddress == knownNodes[0] {
			for _, node := range knownNodes {
				// 给其他节点分发，添加交易
				if node != nodeAddress && node != payload.Hw_AddFrom {
					Hw_sendInv(node, "tx", [][]byte{tx.Hw_ID})
				}
			}
		} else {
			// 如果交易池中有两条交易 并且当前是挖矿节点
			if len(mempool) >= 2 && len(miningAddress) > 0 {
			MineTransactions:
				var txs []*Hw_Transaction

				for id := range mempool {
					tx := mempool[id]
					if bc.Hw_VerifyTransaction(&tx, txs) {
						txs = append(txs, &tx)
					}
				}

				if len(txs) == 0 {
					fmt.Println("交易不可用...")
					break
				}

				cbTx := Hw_NewCoinbaseTX(miningAddress, "")
				txs = append(txs, cbTx)

				newBlock := bc.Hw_MineBlock(txs)
				UTXOSet := Hw_UTXOSet{bc}
				UTXOSet.Update(newBlock)

				fmt.Println("挖到新的区块!")

				for _, tx := range txs {
					txID := hex.EncodeToString(tx.Hw_ID)
					delete(mempool, txID)
				}

				for _, node := range knownNodes {
					if node != nodeAddress {
						Hw_sendInv(node, "block", [][]byte{newBlock.Hw_Hash})
					}
				}

				if len(mempool) > 0 {
					goto MineTransactions
				}
			}
		}
	}
}

// 检查是否需要同步数据
func Hw_handleVersion(request []byte, bc *Hw_Blockchain) {
	var buff bytes.Buffer
	var payload Hw_version
	// 获取数据
	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)

	if err != nil {
		log.Panic(err)
	}

	// 获取当前节点的最后height
	myBestHeight := bc.Hw_GetBestHeight()
	// 数据中的最后height
	foreignerBestHeight := payload.Hw_BestHeight

	// 节点将从消息中提取的 BestHeight 与自身进行比较
	// 当前的height比对方的小
	// 发送获取区块的的命令
	if myBestHeight < foreignerBestHeight {
		Hw_sendGetBlocks(payload.Hw_AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		// 当前的height比对方的大
		// 通知对方节点，同步数据
		Hw_sendVersion(payload.Hw_AddrFrom, bc)
	}

	// 如果不是已知节点，将节点添加到已知节点中
	if !Hw_nodeIsKnown(payload.Hw_AddrFrom) {
		knownNodes = append(knownNodes, payload.Hw_AddrFrom)
	}
}

// 是否是已知节点
func Hw_nodeIsKnown(addr string) bool {
	for _, node := range knownNodes {
		if node == addr {
			return true
		}
	}

	return false
}

//================================================================

// 命令转字节数组
func Hw_commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

// 将字节数组转字符串命令
func Hw_bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

// 加密
func Hw_gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
