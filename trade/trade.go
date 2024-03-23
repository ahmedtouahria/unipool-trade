package quote

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Define structs to represent the response body

type Response struct {
	Data Data `json:"data"`
}

type Data struct {
	EVM EVM `json:"EVM"`
}

type EVM struct {
	BuySide []BuySide `json:"buyside"`
}

type BuySide struct {
	Block                Block  `json:"Block"`
	Trade                Trade  `json:"Trade"`
	DistinctBuyer        string `json:"distinctBuyer"`
	DistinctSeller       string `json:"distinctSeller"`
	DistinctSender       string `json:"distinctSender"`
	DistinctTransactions string `json:"distinctTransactions"`
	TotalBuys            string `json:"total_buys"`
	TotalCount           string `json:"total_count"`
	TotalSales           string `json:"total_sales"`
	Volume               string `json:"volume"`
}

type Block struct {
	Time string `json:"Time"`
}

type Trade struct {
	Currency Currency `json:"Currency"`
	Side     Side     `json:"Side"`
}

type Currency struct {
	Name string `json:"Name"`
}

type Side struct {
	Currency Currency `json:"Currency"`
	High     float64  `json:"high"`
	Low      float64  `json:"low"`
	Open     float64  `json:"open"`
	Close    float64  `json:"close"`
}

// ExecuteGraphQLQuery takes a query string and executes it
func ExecuteGraphQLQuery(query *strings.Reader, apiKey string, token string) ([]byte, error) {
	url := "https://streaming.bitquery.io/graphql"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, query)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-KEY", apiKey)
	req.Header.Add("Authorization", token)

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

func GetDataSelles(smartContract string, MinTotalSelles int) error {
	var trade Response
	isSelled := false
	query := strings.NewReader("{\"query\":\"query MyQuery {\\n        EVM(dataset: archive, network: eth) {\\n            buyside: DEXTradeByTokens(\\n                limit: {count: 30}\\n                orderBy: {descending: Block_Time}\\n                where: {Trade: {Currency: {SmartContract: {is: \\\"" + smartContract + "\\\"}}}, Block: {Date: {since: \\\"2023-07-01\\\", till: \\\"2023-08-01\\\"}}}\\n            ) {\\n                Block {\\n                    Time(interval: {in: days, count: 1})\\n                }\\n                volume: sum(of: Trade_Amount)\\n                distinctBuyer: count(distinct: Trade_Buyer)\\n                distinctSeller: count(distinct: Trade_Seller)\\n                distinctSender: count(distinct: Trade_Sender)\\n                distinctTransactions: count(distinct: Transaction_Hash)\\n                total_sales: count(\\n                    if: {Trade: {Side: {Currency: {SmartContract: {is: \\\"" + smartContract + "\\\"}}}}}\\n                )\\n                total_buys: count(\\n                    if: {Trade: {Currency: {SmartContract: {is: \\\"" + smartContract + "\\\"}}}}\\n                )\\n                total_count: count\\n                Trade {\\n                    Currency {\\n                        Name\\n                    }\\n                    Side {\\n                        Currency {\\n                            Name\\n                        }\\n                    }\\n                    high: Price(maximum: Trade_Price)\\n                    low: Price(minimum: Trade_Price)\\n                    open: Price(minimum: Block_Number)\\n                    close: Price(maximum: Block_Number)\\n                }\\n            }\\n        }\\n    }\",\"variables\":{}}")
	// Execute the GraphQL query
	body, err := ExecuteGraphQLQuery(query, "BQY8Asjo0kL3NFMwshSAIS0iF1l3Yg2S", "Bearer ory_at_m4EWphV-sh_l7Q3JJc4dnETUoa5Hu8kimgAPnXL0GH0.fFP-i9GyIYexQiPSYK4Um1gp6j2MeOMTJ2ohLpJkVYQ")
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &trade)
	for _, buySide := range trade.Data.EVM.BuySide {
		TotalSales, err := strconv.Atoi(buySide.TotalSales)
		if err != nil {
			fmt.Println("Error decoding response body:", err)
			return err
		}
		if TotalSales > MinTotalSelles {
			isSelled = true
			break
		}
	}
	if err != nil {
		fmt.Println("Error decoding response body:", err)

		return err
	}

	if err != nil {
		fmt.Println("Error decoding response body:", err)

		return err
	}
	if isSelled {
		fmt.Printf("\\This token has selles")

	} else {
		fmt.Printf("Not yet available selles\n")

	}

	return nil
}
