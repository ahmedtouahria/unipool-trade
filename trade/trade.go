package quote

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	//"strconv"
	"strings"
)

// Define structs to represent the response body

// Response represents the structure of the JSON response
type Response struct {
	Data struct {
		EVM struct {
			Buyside []TradeData `json:"buyside"`
		} `json:"EVM"`
	} `json:"data"`
}

// TradeData represents the structure of the trade data
type TradeData struct {
	Trade              Trade   `json:"Trade"`
	DistinctBuyer      string  `json:"distinctBuyer"`
	DistinctSeller     string  `json:"distinctSeller"`
	DistinctSender     string  `json:"distinctSender"`
	DistinctTransactions string `json:"distinctTransactions"`
	TotalBuys          string  `json:"total_buys"`
	TotalCount         string  `json:"total_count"`
	TotalSales         string  `json:"total_sales"`
	Volume             string  `json:"volume"`
}

// Trade represents the structure of the trade
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
	query:="{\"query\":\"query MyQuery {\\n        EVM(dataset: archive, network: eth) {\\n            buyside: DEXTradeByTokens(\\n                limit: {count: 30}\\n                orderBy: {descending: Block_Time}\\n                where: {Trade: {Currency: {SmartContract: {is: \\\""+smartContract+"\\\"}}}}\\n            ) {\\n\\n                volume: sum(of: Trade_Amount)\\n                distinctBuyer: count(distinct: Trade_Buyer)\\n                distinctSeller: count(distinct: Trade_Seller)\\n                distinctSender: count(distinct: Trade_Sender)\\n                distinctTransactions: count(distinct: Transaction_Hash)\\n                total_sales: count(\\n                    if: {Trade: {Side: {Currency: {SmartContract: {is: \\\""+smartContract+"\\\"}}}}}\\n                )\\n                total_buys: count(\\n                    if: {Trade: {Currency: {SmartContract: {is: \\\""+smartContract+"\\\"}}}}\\n                )\\n                total_count: count\\n                Trade {\\n                    Currency {\\n                        Name\\n                    }\\n                    Side {\\n                        Currency {\\n                            Name\\n                        }\\n                    }\\n                   high: Price(maximum: Trade_Price)\\n                    low: Price(minimum: Trade_Price)\\n                    open: Price(minimum: Block_Number)\\n                    close: Price(maximum: Block_Number)\\n                }\\n            }\\n        }\\n    }\",\"variables\":{}}"
	reader := strings.NewReader(query)
	// Execute the GraphQL query
	//print(query)
	body, err := ExecuteGraphQLQuery(reader, "BQY8Asjo0kL3NFMwshSAIS0iF1l3Yg2S", "Bearer ory_at_m4EWphV-sh_l7Q3JJc4dnETUoa5Hu8kimgAPnXL0GH0.fFP-i9GyIYexQiPSYK4Um1gp6j2MeOMTJ2ohLpJkVYQ")
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &trade)
	totalSelles:=0
	totalBuys:=0
	for _, buySide := range trade.Data.EVM.Buyside {
		sales, _ := strconv.Atoi(buySide.TotalSales)
		buys, _ := strconv.Atoi(buySide.TotalBuys)
		totalSelles += sales
		totalBuys += buys
		if err != nil {
			fmt.Println("Error decoding response body:", err)
			return err
		}

	}
	// Output total sales and total buys in JSON format
	totalData := map[string]any{
		"contract":smartContract,
		"TotalSales": totalSelles,
		"TotalBuys":  totalBuys,
	}
	totalJSON, err := json.Marshal(totalData)
	if err != nil {
		return err
	}
	fmt.Println(string(totalJSON))
	if err != nil {
		fmt.Println("Error decoding response body:", err)

		return err
	}

	if err != nil {
		fmt.Println("Error decoding response body:", err)

		return err
	}

	return nil
}
