package main

import (
	"flag"
	"fmt"
	"os"

	quote_p "github.com/ahmedtouahria/go-unipool/quote"
)

func main() {
	// Define command-line flags
	smartContract := flag.String("contract", "", "Smart contract address")
	MinTotalSelles := flag.Int("min_selles", 2, "min_selles parameter")
	DataType := flag.String("type", "info", "type parameter must be one of 'info' or 'selles'")

	// Parse command-line flags
	flag.Parse()
	// Check if the contract flag is provided

	if *smartContract == "" {
		fmt.Println("Error: smart contract address is required")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *DataType == "info" {
		quote_p.GetPoolInfo("0xd3d2e2692501a5c9ca623199d38826e513033a17")
	} else if *DataType == "selles" {
		err := quote_p.GetDataSelles(*smartContract, *MinTotalSelles)
		if err != nil {
			panic(err)
		}
	} else {
		flag.PrintDefaults()
		os.Exit(1)
	}

}
