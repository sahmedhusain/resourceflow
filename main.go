package main

import (
	"fmt"
	"os"
	"strconv"

	"resourceflow/internal/simulator"
	"resourceflow/pkg/parser"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ./stock_exchange <file> <waiting_time>")
		os.Exit(1)
	}

	configPath := os.Args[1]
	waitTimeStr := os.Args[2]

	waitTime, err := strconv.ParseFloat(waitTimeStr, 64)
	if err != nil {
		fmt.Printf("Invalid waiting time: %s\n", waitTimeStr)
		os.Exit(1)
	}

	config, err := parser.ParseFile(configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0) // Note: error1 example expects "Exiting..." for some errors, or just prints and exits.
	}

	simulator.Run(config, waitTime, configPath)
}
