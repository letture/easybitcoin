package ethrepc

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"golang.org/x/crypto/sha3"
	"log"
	"math/big"
)

type TransferData struct {
	fromAddress string
	privKey     string
	toAddress   string
	tokenAddress string
	amount  	float64
}

// 查询区块高度
func GetBlockNumber() uint64 {

	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return block.Number().Uint64()
}

// 获取nonce
func getNonce(address common.Address) uint64 {
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		log.Fatal(err)
	}
	return  nonce
}

// 获取燃气费
func getGasPrice() *big.Int{
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return gasPrice
}

// 转账
func Transfer(transferDate TransferData) {
	fromAddress := common.HexToAddress(transferDate.fromAddress)
	privateKey, err := crypto.HexToECDSA(transferDate.privKey)
	if err != nil {
		log.Fatal(err)
	}

	nonce := getNonce(fromAddress)

	toAddress := common.HexToAddress(transferDate.toAddress)

	amount := int64(transferDate.amount * params.Ether)
	value := big.NewInt(amount)

	gasLimit := uint64(21000) // in units
	gasPrice := getGasPrice()
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(signedTx)
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	hash := signedTx.Hash().Hex()
	fmt.Printf("tx sent: %s", hash)
	// 变更状态
}

func TransferToken(transferDate TransferData) {

	fromAddress := common.HexToAddress(transferDate.fromAddress)
	privateKey, err := crypto.HexToECDSA(transferDate.privKey)
	if err != nil {
		log.Fatal(err)
	}
	nonce := getNonce(fromAddress)

	value := big.NewInt(0) // in wei (0 eth)
	gasPrice := getGasPrice()

	toAddress := common.HexToAddress(transferDate.toAddress)
	tokenAddress := common.HexToAddress(transferDate.tokenAddress)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

	amount := new(big.Int)
	amount.SetString("1000000000000000", 18) // sets the value to 1000 tokens, in the token denomination

	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gasLimit) // 23256

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}
