package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"os"
	unipool "github.com/ahmedtouahria/go-unipool/unipool"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Replace this with the actual ABI of your Unipool contract
var unipoolABI = unipool.UnipoolABI

// UnipoolContract represents the contract interface.
type UnipoolContract struct {
	client          *ethclient.Client
	contract        *abi.ABI
	contractAddress common.Address
}

// NewUnipoolContract creates a new instance of the UnipoolContract.
func NewUnipoolContract(client *ethclient.Client, contractAddress string) (*UnipoolContract, error) {
	contractABI, err := abi.JSON(strings.NewReader(unipoolABI))
	if err != nil {
		return nil, err
	}

	address := common.HexToAddress(contractAddress)

	return &UnipoolContract{
		client:          client,
		contract:        &contractABI,
		contractAddress: address,
	}, nil
}

// GetReserves retrieves reserve values (reserve0 and reserve1) from the Unipool contract.
func (u *UnipoolContract) GetReserves() (*big.Int, *big.Int, error) {
	callData, err := u.contract.Pack("getReserves")
	if err != nil {
		return nil, nil, err
	}

	message := ethereum.CallMsg{
		To:   &u.contractAddress,
		Data: callData,
	}

	result, err := u.client.CallContract(context.Background(), message, nil)
	if err != nil {
		return nil, nil, err
	}

	reserves, err := u.contract.Unpack("getReserves", result)
	if err != nil {
		return nil, nil, err
	}

	// Assuming the result is a tuple [reserve0, reserve1]
	reserve0, ok := reserves[0].(*big.Int)
	if !ok {
		return nil, nil, fmt.Errorf("failed to convert reserve0 to *big.Int")
	}

	reserve1, ok := reserves[1].(*big.Int)
	if !ok {
		return nil, nil, fmt.Errorf("failed to convert reserve1 to *big.Int")
	}

	return reserve0, reserve1, nil
}

func main() {
	// Replace with your Ethereum node URL
	nodeURL := "http://localhost:8545"
	client, err := rpc.Dial(nodeURL)
	if err != nil {
		log.Fatal(err)
	}

	// Replace with the actual Unipool contract address
	contractAddress := "0x1234567890123456789012345678901234567890"
	unipoolContract, err := NewUnipoolContract(ethclient.NewClient(client), contractAddress)
	if err != nil {
		log.Fatal(err)
	}

	// Call GetReserves function
	reserve0, reserve1, err := unipoolContract.GetReserves()
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Printf("Reserve0: %v\n", reserve0)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stdout, "Reserve1: %v\n", []any{reserve1}...)
}
