package quote

import(
	"context"
	"fmt"

	"log"
	"math/big"
	unipool "github.com/ahmedtouahria/go-unipool/unipool"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

)
// GetPoolInfo its a function to get the token informations
func GetPoolInfo(address string){
	var tokenEthReserve *big.Int

	nodeURL := "https://eth.public-rpc.com"
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(address)
	unipoolContract, err := unipool.NewUnipool(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	token0, err := unipoolContract.Token0(nil)
	if err != nil {
		log.Fatal(err)
	}

	token1, err := unipoolContract.Token1(nil)

	if err != nil {
		log.Fatal(err)
	}
	tokenIn := "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	var ethTokenAddress = common.HexToAddress(tokenIn)
	//println(ethTokenAddress)
	// Call GetReserves function
	reserves, err := unipoolContract.GetReserves(nil)
	if err != nil {
		log.Fatal(err)
	}

	var tokenOut string

	if token0 == ethTokenAddress {
		fmt.Println("Token 0 represents ETH == ", reserves.Reserve0)
		tokenEthReserve = reserves.Reserve0
		tokenOutRaw, _ := unipoolContract.Token1(nil)
		tokenOut = fmt.Sprint(tokenOutRaw)
	} else if token1 == ethTokenAddress {
		fmt.Println("Token 1 represents ETH")
		tokenEthReserve = reserves.Reserve1
		tokenOutRaw, _ := unipoolContract.Token0(nil)
		tokenOut = fmt.Sprint(tokenOutRaw)
	} else {
		fmt.Println("Neither token represents ETH")
	}
	totalSuply,err := unipoolContract.TotalSupply(nil)
	if err != nil {
		log.Fatal(err)
	}
	UsdUnitPrice,err:= GetETHPriceInUSD()
	if err != nil {
		log.Fatal(err)
	}
	// Divide tokenEthReserve by 10^18
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	tokenEthReserveHumanReadable := new(big.Int).Div(tokenEthReserve, divisor)
	//in usd
	tokenEthReserveHumanReadableUsd := new(big.Int).Mul(tokenEthReserveHumanReadable,big.NewInt(int64(UsdUnitPrice)))
	minimumLiquidity, err := unipoolContract.MINIMUMLIQUIDITY(nil)
	if err != nil {
		log.Fatal(err)
	}
	name, err := unipoolContract.Name(nil)
	if err != nil {
		log.Fatal(err)
	}
	actualGasPrice, _ := client.SuggestGasPrice(context.Background())
	quote := GetQuote(tokenIn, tokenOut, "100000000000000000")

	fmt.Println("Minimum Minimum Liquidity:", minimumLiquidity)
	fmt.Println("TotalSupply:",totalSuply)

	fmt.Println("Token name:", name)
	fmt.Println("Token Liquidity ETH:", tokenEthReserveHumanReadable,"eth")
	fmt.Println("Token Liquidity USD:", tokenEthReserveHumanReadableUsd,"$")

	fmt.Println("Actual gas price:", actualGasPrice)
	fmt.Println("Avg gas price for a small tx:", quote.Quote.GasUseEstimate)
	fmt.Println("Avg gas price for a small tx usd:", quote.Quote.GasUseEstimateUSD)
	fmt.Println("1ETH =", UsdUnitPrice, "usd")
}
