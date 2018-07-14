> 注： windows 下设置临时环境变量命令为 set

----
#### 主节点：
```
# 设置临时环境变量，用于网络端口
export NODE_ID=3000

# 创建钱包 wallet_3000.dat
./main createwallet
ADDRESS

# 使用上面生成的钱包地址 创建创世区块 blockchain_3000.db
./main createblockchain -address ADDRESS

# 查看余额
./main getbalance -address ADDRESS
10

# 转到钱包节点

# 从钱包节点回来
# 查看所有钱包地址
./main listaddresses
ADDRESS

# 转账
./main send -from ADDRESS -to ADDRESS1 -amount 8 -mine
./main send -from ADDRESS -to ADDRESS2 -amount 6 -mine

# 查询余额
# ADDRESS = 10 - 8 + 10 - 6 + 10
./main getbalance -address ADDRESS
16
./main getbalance -address ADDRESS1
8
./main getbalance -address ADDRESS2
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
ADDRESS1
./main createwallet
ADDRESS2
./main createwallet
ADDRESS3
./main createwallet
ADDRESS3

# 查看所有钱包地址
./main listaddresses
ADDRESS1
ADDRESS2
ADDRESS3
ADDRESS4

# 切换到主节点

# 从主节点回来
./main getbalance -address ADDRESS
10
./main getbalance -address ADDRESS1
0
./main getbalance -address ADDRESS2
0

# 切换到主节点

# 从主节点回来
# 启动节点服务器
./main startnode

# 停止节点，查询余额，正常
./main getbalance -address ADDRESS
16
./main getbalance -address ADDRESS1
8
./main getbalance -address ADDRESS2
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
ADDRESS5

# 启动矿工节点服务器
./main startnode -miner ADDRESS5
```

#### 测试
```
#运行主节点(3000)和挖矿节点(3002),在钱包节点(3001)进行转账测试

# ADDRESS1 转 ADDRESS3  5个币
./main send -from ADDRESS1 -to ADDRESS3 -amount 5

# ADDRESS2 转 ADDRESS4  5个币
./main send -from ADDRESS2 -to ADDRESS4 -amount 5

# 查询余额
./main getbalance -address ADDRESS1
3
./main getbalance -address ADDRESS2
1
./main getbalance -address ADDRESS3
5
./main getbalance -address ADDRESS4
5

# 现在查看钱包节点钱包余额，没有任何变化
# 因为没有启动服务器
./main startnode

# 钱包节点停止服务，再查询余额，发现正常
# 至此测试内容结束
```