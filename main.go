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
	MinTotalSelles := flag.Int("min_selles", 2, "min_selles parameter")
	// Parse command-line flags
	flag.Parse()
	// Check if the contract flag is provided
	if *smartContract == "" {
		fmt.Println("Error: smart contract address is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	err := trade.GetDataSelles(*smartContract, *MinTotalSelles)
	if err != nil {
		panic(err)
	}

}
