package main

import (
	"fmt"
	"log"

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
	contractAddress := common.HexToAddress("0xd3d2e2692501a5c9ca623199d38826e513033a17")
	unipoolContract, err := unipool.NewUnipool(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// Call GetReserves function
	reserves, err := unipoolContract.GetReserves(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("reserve 0 : %v, reserver 1 : %v \n", reserves.Reserve0, reserves.Reserve1)
}
