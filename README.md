# unipool-trade

`unipool-trade` is a Go package for fetching and processing trade data from Uniswap pools.

## Installation

To install `unipool-trade`, use `go get`:

```bash
go get github.com/ahmedtouahria/unipool-trade/trade
```

## Usage

```go
package main

import (
	"flag"
	"fmt"
	"os"

	trade "github.com/ahmedtouahria/unipool-trade/trade"
)

func main() {
	// Define command-line flags
	smartContract := flag.String("contract", "", "Smart contract address")
	minTotalSelles := flag.Int("min_selles", 2, "Minimum number of selles parameter")
	// Parse command-line flags
	flag.Parse()
	// Check if the contract flag is provided
	if *smartContract == "" {
		fmt.Println("Error: smart contract address is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	err := trade.GetDataSelles(*smartContract, *minTotalSelles)
	if err != nil {
		panic(err)
	}
}
```

Run the binary and provide the smart contract address using the `-contract` flag.

## Parameters

- `contract`: Smart contract address of the Uniswap pool.
- `min_selles`: Minimum number of sellers required.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

Special thanks to [ahmedtouahria](https://github.com/ahmedtouahria) for developing this package.

--
