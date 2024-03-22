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
	timesmap := flag.Int("timesmap", 10, "timesmap parameter")
	MinTotalSelles := flag.Int("min_selles", 2, "MinTotalSelles parameter")

	quote_p.GetPoolInfo("0xd3d2e2692501a5c9ca623199d38826e513033a17")
	// Parse command-line flags
	flag.Parse()
	// Check if the contract flag is provided

	if *smartContract == "" {
		fmt.Println("Error: smart contract address is required")
		flag.PrintDefaults()
		os.Exit(1)
	}
	err := quote_p.GetDataSelles(*smartContract, *timesmap, *MinTotalSelles)

	if err != nil {
		panic(err)
	}
}
