package main

import (
	"fmt"
	"log"
	"os"
	unipool "github.com/ahmedtouahria/go-unipool/unipool"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)



func main() {
	// Replace with your Ethereum node URL
	nodeURL := "https://eth.public-rpc.com"
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		log.Fatal(err)
	}

	// Replace with the actual Unipool contract address
	contractAddress := common.HexToAddress("0xC75650fe4D14017b1e12341A97721D5ec51D5340")
	unipoolContract, err := unipool.NewUnipool(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// Call GetReserves function
	reserves, err := unipoolContract.GetReserves(nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Printf("Reserve0: %v\n", reserves.Reserve0)
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintf(os.Stdout, "Reserve1: %v\n", []any{reserves.Reserve1}...)
	if err != nil {
		panic(err)
	}
}
