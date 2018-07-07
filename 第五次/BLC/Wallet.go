package BLC

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
	"publicchain/day5-base58/BLC"
	"bytes"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
)

const version = byte(0x00)
const addressChecksumLen = 4

type Wallet struct {
	Privatekey ecdsa.PrivateKey
	PublicKey  []byte
}

/*创建一个钱包地址
1.生成一对公钥和私钥
2.想获取地址，可以通过公钥进行Base58编码
3.想要别人给我转账，把地址给别人，别人将地址进行反编码变成公钥，将公钥和数据进行签名
4.通过私钥进行解密，解密是单方向的，只有用私钥的人才能进行解密
*/

func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

//创建私钥、公钥
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubkey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubkey
}
func (w *Wallet) IsValidForAddress(address []byte) bool {
	//解码 25字节
	version_public_checksumBytes := Base58Decode(address)
	fmt.Println(version_public_checksumBytes)
	//截取后面的4个字节
	checkSumBytes := version_public_checksumBytes[len(version_public_checksumBytes)-addressChecksumLen:]
	//截取前面的21个字节s
	version_ripemd160 := version_public_checksumBytes[:len(version_public_checksumBytes)-addressChecksumLen]
	fmt.Printf("%s\n", len(checkSumBytes))
	fmt.Printf("%s\n", len(version_ripemd160))

	checkBytes := checksum(version_ripemd160)
	if bytes.Compare(checkSumBytes, checkBytes) == 0 {
		return true
	}
	return false
}

//返回钱包地址
func (w Wallet) GetAddress() []byte {
	// 1. hash160
	// 20字节
	pubkeyHash := Ripemd160Hash(w.PublicKey)
	// 21 字节
	versionedPayload := append([]byte{version}, pubkeyHash...)
	// 两次的256 hash
	checksum := checksum(versionedPayload)
	//25 字节
	fullPayload := append(versionedPayload, checksum...)
	//base58 编码加密
	address := BLC.Base58Encoder(fullPayload)
	return address
}

// ripemd160(sha256(PubKey))
func Ripemd160Hash(pubkey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubkey)
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160

}

//验证地址的有效性
func ValidateAddress(address string) bool {
	pubkeyHash := BLC.Base58Decode([]byte(address))
	actualChecksum := pubkeyHash[len(pubkeyHash)-addressChecksumLen:]
	version := pubkeyHash[0]
	pubkeyHash = pubkeyHash[1: len(pubkeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubkeyHash...))
	return bytes.Compare(actualChecksum, targetChecksum) == 0

}

//Checksum 为一个公钥生成 checksum
func checksum(payload []byte) []byte {
	//两次进行Sum256 Hash
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:addressChecksumLen]
}
