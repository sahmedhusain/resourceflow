package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"stock-exchange-sim/pkg/models"
	"stock-exchange-sim/pkg/parser"
)

type LogEvent struct {
	Cycle       int
	ProcessName string
}

type FinishEvent struct {
	Cycle   int
	Process *models.Process
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ./checker <file> <log_file>")
		os.Exit(1)
	}

	configPath := os.Args[1]
	logPath := os.Args[2]

	config, err := parser.ParseFile(configPath)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Exiting...")
		os.Exit(1)
	}

	f, err := os.Open(logPath)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		os.Exit(1)
	}
	defer f.Close()

	logs := []LogEvent{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			cycle, err := strconv.Atoi(parts[0])
			if err == nil {
				logs = append(logs, LogEvent{Cycle: cycle, ProcessName: parts[1]})
			}
		}
	}

	stocks := make(map[string]int)
	for k, v := range config.InitialStocks {
		stocks[k] = v
	}

	processesMap := make(map[string]*models.Process)
	for i := range config.Processes {
		p := &config.Processes[i]
		processesMap[p.Name] = p
	}

	finishEvents := []FinishEvent{}

	for _, log := range logs {
		fmt.Printf("Evaluating: %d:%s\n", log.Cycle, log.ProcessName)

		// 1. Apply finish events that complete at or before this log's cycle
		newFinishes := []FinishEvent{}
		for _, fe := range finishEvents {
			if fe.Cycle <= log.Cycle {
				for k, v := range fe.Process.Results {
					stocks[k] += v
				}
			} else {
				newFinishes = append(newFinishes, fe)
			}
		}
		finishEvents = newFinishes

		// 2. Validate process existence
		p, exists := processesMap[log.ProcessName]
		if !exists {
			fmt.Println("Error detected")
			fmt.Printf("at %d:%s process unknown\n", log.Cycle, log.ProcessName)
			fmt.Println("Exiting...")
			os.Exit(1)
		}

		// 3. Check stock sufficiency
		for k, v := range p.Needs {
			if stocks[k] < v {
				fmt.Println("Error detected")
				fmt.Printf("at %d:%s stock insufficient\n", log.Cycle, log.ProcessName)
				fmt.Println("Exiting...")
				os.Exit(1)
			}
		}

		// 4. Deduct stocks
		for k, v := range p.Needs {
			stocks[k] -= v
		}

		// 5. Add finish event
		finishEvents = append(finishEvents, FinishEvent{
			Cycle:   log.Cycle + p.Duration,
			Process: p,
		})
	}

	fmt.Println("Trace completed, no error detected.")
}
