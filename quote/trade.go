package quote

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	EVM EVM `json:"EVM"`
}

type EVM struct {
	Buyside []TradeData `json:"buyside"`
}

type TradeData struct {
	Block              Block     `json:"Block"`
	Trade              Trade     `json:"Trade"`
	DistinctBuyer      string    `json:"distinctBuyer"`
	DistinctSeller     string    `json:"distinctSeller"`
	DistinctSender     string    `json:"distinctSender"`
	DistinctTransactions string  `json:"distinctTransactions"`
	TotalBuys          string    `json:"total_buys"`
	TotalCount         string    `json:"total_count"`
	TotalSales         string    `json:"total_sales"`
	Volume             string    `json:"volume"`
}

type Block struct {
	Time string `json:"Time"`
}

type Trade struct {
	Currency struct {
		Name string `json:"Name"`
	} `json:"Currency"`
	Side struct {
		Currency struct {
			Name string `json:"Name"`
		} `json:"Currency"`
		Close float64 `json:"close"`
		High  float64 `json:"high"`
		Low   float64 `json:"low"`
		Open  float64 `json:"open"`
	} `json:"Side"`
}

// ExecuteGraphQLQuery takes a query string and executes it
func ExecuteGraphQLQuery(query string, apiKey string, token string) ([]byte, error) {
	url := "https://streaming.bitquery.io/graphql"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{"query":"%s"}`, query))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-KEY", apiKey)
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetDataSelles(smartContract string, since string, till string) error {
	query := `query MyQuery {
		EVM(dataset: archive, network: bsc) {
			buyside: DEXTradeByTokens(
				limit: {count: 30}
				orderBy: {descending: Block_Time}
				where: {Trade: {Currency: {SmartContract: {is: "` + smartContract + `"}}}, Block: {Time: {since: "` + since + `", till: "` + till + `"}}}
			) {
				Block {
					Time(interval: {in: days, count: 1})
				}
				volume: sum(of: Trade_Amount)
				distinctBuyer: count(distinct: Trade_Buyer)
				distinctSeller: count(distinct: Trade_Seller)
				distinctSender: count(distinct: Trade_Sender)
				distinctTransactions: count(distinct: Transaction_Hash)
				total_sales: count(
					if: {Trade: {Side: {Currency: {SmartContract: {is: "` + smartContract + `"}}}}}
				)
				total_buys: count(
					if: {Trade: {Currency: {SmartContract: {is: "` + smartContract + `"}}}}
				)
				total_count: count
				Trade {
					Currency {
						Name
					}
					Side {
						Currency {
							Name
						}
						high: Price(maximum: Trade_Price)
						low: Price(minimum: Trade_Price)
						open: Price(minimum: Block_Number)
						close: Price(maximum: Block_Number)
					}
				}
			}
		}
	}`

	// Execute the GraphQL query
	body, err := ExecuteGraphQLQuery(query, "BQY8Asjo0kL3NFMwshSAIS0iF1l3Yg2S", "ory_at_SQq4FXsEqBT5jsNfnljCL3wtx4X80yaSU8JgJYpSahU.VrLQ68qfrmc-rEkCSyram3A0mX94BqrGMRjSD1eEnuo")
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	var trade Response
	err = json.Unmarshal(body, &trade)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", trade)

	return nil
}
