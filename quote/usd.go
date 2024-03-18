package quote

import (
	"encoding/json"
	"io"
	"net/http"
)

type CoinGeckoResponse struct {
	Ethereum struct {
		Usd float64 `json:"usd"`
	} `json:"ethereum"`
}

func GetETHPriceInUSD() (float64, error) {
	// Send GET request to CoinGecko API
	response, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd")
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	// Unmarshal JSON response
	var cgResponse CoinGeckoResponse
	err = json.Unmarshal(body, &cgResponse)
	if err != nil {
		return 0, err
	}

	// Extract USD price of ETH
	ethUSDPrice := cgResponse.Ethereum.Usd
	return ethUSDPrice, nil
}
