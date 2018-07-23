> 注： windows 下设置临时环境变量命令为 set

----
#### 主节点：
```
# 重新编译
go build main.go

# 设置临时环境变量，用于网络端口
export NODE_ID=3000

# 创建钱包 wallet_3000.dat
./main createwallet
1FUdTVBF6e5u9bdwF5fBHdFyVD4aneygeR

# 使用上面生成的钱包地址 创建创世区块 blockchain_3000.db
./main createblockchain -address 1FUdTVBF6e5u9bdwF5fBHdFyVD4aneygeR

# 查看余额
./main getbalance -address 1FUdTVBF6e5u9bdwF5fBHdFyVD4aneygeR
10

# 转到钱包节点

# 从钱包节点回来
# 查看所有钱包地址
./main listaddresses
ADDRESS

# 转账
./main send -from 1FUdTVBF6e5u9bdwF5fBHdFyVD4aneygeR -to 1HJS1wZYqxM78g8kRFsyXmfHT8SwecvEeB -amount 8 -mine
./main send -from 1FUdTVBF6e5u9bdwF5fBHdFyVD4aneygeR -to 1BDcTrNULhy1yarFvKTeYYivCjwR5y39Fr -amount 6 -mine

# 查询余额
# ADDRESS = 10 - 8 + 10 - 6 + 10
./main getbalance -address 1FUdTVBF6e5u9bdwF5fBHdFyVD4aneygeR
16
./main getbalance -address 1HJS1wZYqxM78g8kRFsyXmfHT8SwecvEeB
8
./main getbalance -address 1BDcTrNULhy1yarFvKTeYYivCjwR5y39Fr
6

# 转到钱包节点

# 从钱包节点回来
# 启动节点服务器
./main startnode
```


#### 钱包节点
```
# 备份创建区块内容备用
cp blockchain_3000.db blockchain_genesis.db

# 复制只带创世区块的数据库
cp blockchain_genesis.db blockchain_3001.db

# 设置临时环境变量，用于网络端口
export NODE_ID=3001

# 生成四个钱包地址
./main createwallet
1HJS1wZYqxM78g8kRFsyXmfHT8SwecvEeB
./main createwallet
1BDcTrNULhy1yarFvKTeYYivCjwR5y39Fr
./main createwallet
146igzxoYs1VcTfR1H2wzxd1cFA35K5Usb
./main createwallet
12nJP2HkvT7MsknXEaU8twJpmeuS6CAgY2

# 查看所有钱包地址
./main listaddresses
1HJS1wZYqxM78g8kRFsyXmfHT8SwecvEeB
1BDcTrNULhy1yarFvKTeYYivCjwR5y39Fr
146igzxoYs1VcTfR1H2wzxd1cFA35K5Usb
12nJP2HkvT7MsknXEaU8twJpmeuS6CAgY2

# 切换到主节点

# 从主节点回来
./main getbalance -address 1FUdTVBF6e5u9bdwF5fBHdFyVD4aneygeR
10
./main getbalance -address 1HJS1wZYqxM78g8kRFsyXmfHT8SwecvEeB
0
./main getbalance -address 1BDcTrNULhy1yarFvKTeYYivCjwR5y39Fr
0

# 切换到主节点

# 从主节点回来
# 启动节点服务器
./main startnode

# 停止节点，查询余额，正常
./main getbalance -address 1FUdTVBF6e5u9bdwF5fBHdFyVD4aneygeR
16
./main getbalance -address 1HJS1wZYqxM78g8kRFsyXmfHT8SwecvEeB
8
./main getbalance -address 1BDcTrNULhy1yarFvKTeYYivCjwR5y39Fr
6
```

#### 矿工节点
```
# 设置临时环境变量，用于网络端口
export NODE_ID=3002

# 复制只带创世区块的数据库
cp blockchain_genesis.db blockchain_3002.db

# 生成钱包地址(挖矿奖励)
./main createwallet
16K5hepB4nWTJhy1jqfAk2pR5jVSxyyQEi

# 启动矿工节点服务器
./main startnode -miner 16K5hepB4nWTJhy1jqfAk2pR5jVSxyyQEi
```

#### 测试
```

#运行主节点(3000)和挖矿节点(3002),在钱包节点(3001)进行转账测试

# ADDRESS1 转 ADDRESS3  5个币
./main send -from 1HJS1wZYqxM78g8kRFsyXmfHT8SwecvEeB -to 146igzxoYs1VcTfR1H2wzxd1cFA35K5Usb -amount 5

# ADDRESS2 转 ADDRESS4  5个币
./main send -from 1BDcTrNULhy1yarFvKTeYYivCjwR5y39Fr -to 12nJP2HkvT7MsknXEaU8twJpmeuS6CAgY2 -amount 5

# 查询余额
./main getbalanceall
地址:1BDcTrNULhy1yarFvKTeYYivCjwR5y39Fr的余额为：1
地址:146igzxoYs1VcTfR1H2wzxd1cFA35K5Usb的余额为：5
地址:12nJP2HkvT7MsknXEaU8twJpmeuS6CAgY2的余额为：5
地址:1HJS1wZYqxM78g8kRFsyXmfHT8SwecvEeB的余额为：3

# 现在查看钱包节点钱包余额，没有任何变化
# 因为没有启动服务器
./main startnode

# 钱包节点停止服务，再查询余额，发现正常

# 继续测试同时多笔转账
# ADDRESS3 转 ADDRESS1 2个币，ADDRESS4 转 ADDRESS2 3个币
./main send -from '["146igzxoYs1VcTfR1H2wzxd1cFA35K5Usb","12nJP2HkvT7MsknXEaU8twJpmeuS6CAgY2"]' -to '["1HJS1wZYqxM78g8kRFsyXmfHT8SwecvEeB","1BDcTrNULhy1yarFvKTeYYivCjwR5y39Fr"]' -amount '["2","3"]'

# 启动服务器
./main startnode

# 结束服务器后
# 查询余额
./main getbalanceall
地址:146igzxoYs1VcTfR1H2wzxd1cFA35K5Usb的余额为：3
地址:12nJP2HkvT7MsknXEaU8twJpmeuS6CAgY2的余额为：2
地址:1HJS1wZYqxM78g8kRFsyXmfHT8SwecvEeB的余额为：5
地址:1BDcTrNULhy1yarFvKTeYYivCjwR5y39Fr的余额为：4

# 至此测试内容结束
```