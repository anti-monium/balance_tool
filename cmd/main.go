package main

import (
	"balance_tool/balance"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	Address string `long:"address" required:"true" description:"Wallet address to receive balance"`
}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		flags.NewParser(&opts, flags.Default).WriteHelp(os.Stdout)
		return
	}
	log.SetFlags(log.Flags() &^ (log.Ldate))

	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(opts.Address) {
		log.Fatal("wallet address is not valid", opts.Address)
	}
	address := common.HexToAddress(opts.Address)
	log.Println("we get the address", address.Hex())

	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("we have a connection")

	bytecode, err := client.CodeAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal(err)
	}
	if len(bytecode) > 0 {
		log.Fatal("the received address is the address of the smart contract")
	}

	_, ethBalance, _, wethBalance, exchangeRate, usdBalance := balance.EthWethBalance(address, client)
	fmt.Println("ETH balance:", ethBalance)
	fmt.Println("WETH balance:", wethBalance)
	fmt.Println("ETH / USD:", exchangeRate)
	fmt.Printf("total USD balance: %.2f\n", usdBalance)
}
