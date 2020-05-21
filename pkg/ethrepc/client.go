package ethrepc

import (
	"github.com/LETTRUE/easybitcoin/pkg/setting"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

var client *ethclient.Client

func Client() {
	var err error
	client, err = ethclient.Dial(setting.EthSetting.Url)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func Close(client *ethclient.Client) {
	defer client.Close()
}
