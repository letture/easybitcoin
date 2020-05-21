package main

import (
	"github.com/LETTRUE/easybitcoin/pkg/ethrepc"
	"github.com/LETTRUE/easybitcoin/pkg/setting"
)

func init() {
	setting.Setup()
}

func main() {
	ethrepc.Transfer()
}
