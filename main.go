package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"math/big"

	unipool "github.com/ahmedtouahria/go-unipool/unipool"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const chain int32 = 1

// Add a global HTTP client with connection pooling
var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	},
}

type Token struct {
	Address    string `json:"address"`
	BuyFeeBps  string `json:"buyFeeBps"`
	ChainId    int    `json:"chainId"`
	Decimals   string `json:"decimals"`
	SellFeeBps string `json:"sellFeeBps"`
	Symbol     string `json:"symbol"`
}

type Reserve struct {
	Quotient string `json:"quotient"`
	Token    Token  `json:"token"`
}

type RouteEntry struct {
	Address   string  `json:"address"`
	AmountIn  string  `json:"amountIn"`
	AmountOut string  `json:"amountOut"`
	Reserve0  Reserve `json:"reserve0"`
	Reserve1  Reserve `json:"reserve1"`
	TokenIn   Token   `json:"tokenIn"`
	TokenOut  Token   `json:"tokenOut"`
	Type      string  `json:"type"`
}

type Quote struct {
	AllQuotes []struct {
		Quote struct {
			Amount                      string `json:"amount"`
			AmountDecimals              string `json:"amountDecimals"`
			BlockNumber                 string `json:"blockNumber"`
			GasPriceWei                 string `json:"gasPriceWei"`
			GasUseEstimate              string `json:"gasUseEstimate"`
			GasUseEstimateQuote         string `json:"gasUseEstimateQuote"`
			GasUseEstimateQuoteDecimals string `json:"gasUseEstimateQuoteDecimals"`
			GasUseEstimateUSD           string `json:"gasUseEstimateUSD"`
			HitsCachedRoutes            bool   `json:"hitsCachedRoutes"`
			MethodParameters            struct {
				Calldata string `json:"calldata"`
				To       string `json:"to"`
				Value    string `json:"value"`
			} `json:"methodParameters"`
			PortionAmount                      string `json:"portionAmount"`
			PortionAmountDecimals              string `json:"portionAmountDecimals"`
			PortionBips                        int    `json:"portionBips"`
			Quote                              string `json:"quote"`
			QuoteDecimals                      string `json:"quoteDecimals"`
			QuoteGasAdjusted                   string `json:"quoteGasAdjusted"`
			QuoteGasAdjustedDecimals           string `json:"quoteGasAdjustedDecimals"`
			QuoteGasAndPortionAdjusted         string `json:"quoteGasAndPortionAdjusted"`
			QuoteGasAndPortionAdjustedDecimals string `json:"quoteGasAndPortionAdjustedDecimals"`
			QuoteId                            string `json:"quoteId"`
			RequestId                          string `json:"requestId"`
			Route                              [][]struct {
				Address   string `json:"address"`
				AmountIn  string `json:"amountIn"`
				AmountOut string `json:"amountOut"`
				Reserve0  struct {
					Quotient string `json:"quotient"`
					Token    struct {
						Address    string `json:"address"`
						BuyFeeBps  string `json:"buyFeeBps"`
						ChainId    int    `json:"chainId"`
						Decimals   string `json:"decimals"`
						SellFeeBps string `json:"sellFeeBps"`
						Symbol     string `json:"symbol"`
					} `json:"token"`
				} `json:"reserve0"`
				Reserve1 struct {
					Quotient string `json:"quotient"`
					Token    struct {
						Address    string `json:"address"`
						BuyFeeBps  string `json:"buyFeeBps"`
						ChainId    int    `json:"chainId"`
						Decimals   string `json:"decimals"`
						SellFeeBps string `json:"sellFeeBps"`
						Symbol     string `json:"symbol"`
					} `json:"token"`
				} `json:"reserve1"`
				TokenIn  Token  `json:"tokenIn"`
				TokenOut Token  `json:"tokenOut"`
				Type     string `json:"type"`
			} `json:"route"`
			RouteString      string  `json:"routeString"`
			SimulationError  bool    `json:"simulationError"`
			SimulationStatus string  `json:"simulationStatus"`
			Slippage         float64 `json:"slippage"`
			TradeType        string  `json:"tradeType"`
		} `json:"quote"`
		Routing string `json:"routing"`
	} `json:"allQuotes"`
	Quote struct {
		Amount                      string `json:"amount"`
		AmountDecimals              string `json:"amountDecimals"`
		BlockNumber                 string `json:"blockNumber"`
		GasPriceWei                 string `json:"gasPriceWei"`
		GasUseEstimate              string `json:"gasUseEstimate"`
		GasUseEstimateQuote         string `json:"gasUseEstimateQuote"`
		GasUseEstimateQuoteDecimals string `json:"gasUseEstimateQuoteDecimals"`
		GasUseEstimateUSD           string `json:"gasUseEstimateUSD"`
		HitsCachedRoutes            bool   `json:"hitsCachedRoutes"`
		MethodParameters            struct {
			Calldata string `json:"calldata"`
			To       string `json:"to"`
			Value    string `json:"value"`
		} `json:"methodParameters"`
		PortionAmount                      string         `json:"portionAmount"`
		PortionAmountDecimals              string         `json:"portionAmountDecimals"`
		PortionBips                        int            `json:"portionBips"`
		Quote                              string         `json:"quote"`
		QuoteDecimals                      string         `json:"quoteDecimals"`
		QuoteGasAdjusted                   string         `json:"quoteGasAdjusted"`
		QuoteGasAdjustedDecimals           string         `json:"quoteGasAdjustedDecimals"`
		QuoteGasAndPortionAdjusted         string         `json:"quoteGasAndPortionAdjusted"`
		QuoteGasAndPortionAdjustedDecimals string         `json:"quoteGasAndPortionAdjustedDecimals"`
		QuoteId                            string         `json:"quoteId"`
		RequestId                          string         `json:"requestId"`
		Route                              [][]RouteEntry `json:"route"`
		RouteString                        string         `json:"routeString"`
		SimulationError                    bool           `json:"simulationError"`
		SimulationStatus                   string         `json:"simulationStatus"`
		Slippage                           float64        `json:"slippage"`
		TradeType                          string         `json:"tradeType"`
	} `json:"quote"`
	RequestId string `json:"requestId"`
	Routing   string `json:"routing"`
}

func getQuote(tokenIn string, tokenOut string, amount string) *Quote {
	url := "https://interface.gateway.uniswap.org/v2/quote"
	payload := map[string]interface{}{
		"tokenInChainId":     chain,
		"tokenIn":            tokenIn,
		"tokenOutChainId":    chain,
		"tokenOut":           tokenOut,
		"amount":             amount,
		"sendPortionEnabled": true,
		"type":               "EXACT_INPUT",
		"configs": []map[string]interface{}{
			{
				"protocols":                      []string{"V2"},
				"enableUniversalRouter":          true,
				"routingType":                    "CLASSIC",
				"enableFeeOnTransferFeeFetching": true,
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding payload:", err)
		return nil
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	req.Header.Set("Origin", "https://app.uniswap.org")
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req) // Use the global HTTP client
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	// fmt.Println("Response Body:", string(body)) // Print response body for debugging

	var quote Quote
	err = json.Unmarshal(body, &quote)
	if err != nil {
		fmt.Println("Error decoding response body:", err)
		return nil
	}

	return &quote
}

func main() {
	var tokenEthReserve *big.Int

	nodeURL := "https://eth.public-rpc.com"
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0xd3d2e2692501a5c9ca623199d38826e513033a17")
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

	// Divide tokenEthReserve by 10^18
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	tokenEthReserveHumanReadable := new(big.Int).Div(tokenEthReserve, divisor)

	minimumLiquidity, err := unipoolContract.MINIMUMLIQUIDITY(nil)
	if err != nil {
		log.Fatal(err)
	}
	name, err := unipoolContract.Name(nil)
	if err != nil {
		log.Fatal(err)
	}
	actualGasPrice, _ := client.SuggestGasPrice(context.Background())
	quote := getQuote(tokenIn, tokenOut, "100000000000000000")
	fmt.Println("Minimum Minimum Liquidity:", minimumLiquidity)
	fmt.Println("Token name:", name)
	fmt.Println("Token Liquidity:", tokenEthReserveHumanReadable)
	fmt.Println("Actual gas price:", actualGasPrice)
	fmt.Println("Avg gas price for a small tx:", quote.Quote.GasUseEstimate)
	fmt.Println("Avg gas price for a small tx usd:", quote.Quote.GasUseEstimateUSD)
}
