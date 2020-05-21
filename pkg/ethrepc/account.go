package ethrepc

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
	"log"
	"math"
	"math/big"
	"regexp"
)

func Generate() {
	// 检查现有地址

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	//address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	//fmt.Println(address) // 0x96216849c49358B10257cb55b28eA603c874b05E

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	address := hexutil.Encode(hash.Sum(nil)[12:])
	fmt.Println(address)

	// 入库

}

// 获取余额
func GetBalance(address string) *big.Float {
	client := Client()
	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	return ethValue
}

// 校验地址
func CheckAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}
