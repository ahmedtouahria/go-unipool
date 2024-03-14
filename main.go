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
    contractAddress := common.HexToAddress("0x2Bc3368e7876BE81C359dfAd7453ddDE1D13a4EB")
    unipoolContract, err := unipool.NewUnipool(contractAddress, client)
    if err != nil {
        log.Fatal(err)
    }

    // Call Token0 and Token1 functions to get token addresses
    token0, err := unipoolContract.Token0(nil)
    if err != nil {
        log.Fatal(err)
    }

    token1, err := unipoolContract.Token1(nil)
    if err != nil {
        log.Fatal(err)
    }

    var ethTokenAddress = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")

    if token0 == ethTokenAddress {
        fmt.Println("Token 0 represents ETH")
    } else if token1 == ethTokenAddress {
        fmt.Println("Token 1 represents ETH")
    } else {
        fmt.Println("Neither token represents ETH")
    }
}