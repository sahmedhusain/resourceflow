package simulator

import (
	"fmt"
	"os"
	"sort"
	"time"

	"resourceflow/pkg/models"
)

type Event struct {
	Cycle   int
	Process *models.Process
}

func Run(config *models.Config, waitTimeSeconds float64, configPath string) {
	logFile := configPath + ".log"
	f, err := os.Create(logFile)
	if err != nil {
		fmt.Printf("Error creating log file: %v\n", err)
		return
	}
	defer f.Close()

	cycle := 0
	startTime := time.Now()
	timeout := time.Duration(waitTimeSeconds * float64(time.Second))

	events := []Event{}
	stocks := make(map[string]int)
	inProgress := make(map[string]int)
	
	for k, v := range config.InitialStocks {
		stocks[k] = v
	}

	fmt.Println("Main Processes:")

	virtualDemand := make(map[string]int)

	for {
		if time.Since(startTime) > timeout {
			break
		}

		// Process finished events for the current cycle
		newEvents := []Event{}
		for _, e := range events {
			if e.Cycle == cycle {
				for k, v := range e.Process.Results {
					stocks[k] += v
					inProgress[k] -= v
				}
			} else {
				newEvents = append(newEvents, e)
			}
		}
		events = newEvents

		// Setup base demand for this cycle
		for _, t := range config.Optimize.Stocks {
			if virtualDemand[t] <= 0 {
				virtualDemand[t] = 1
			}
		}
		
		if config.Optimize.Time && len(config.Optimize.Stocks) == 0 {
			for _, p := range config.Processes {
				for k := range p.Results {
					if virtualDemand[k] <= 0 {
						virtualDemand[k] = 1
					}
				}
			}
		}

		anyStarted := false
		for {
			if time.Since(startTime) > timeout {
				break
			}
			startedInThisPass := false
			demandChanged := false

			for i := range config.Processes {
				p := &config.Processes[i]

				producesDemand := false
				for k := range p.Results {
					if virtualDemand[k] > 0 {
						producesDemand = true
						break
					}
				}
				
				if !producesDemand {
					continue
				}

				canRun := true
				for k, v := range p.Needs {
					if stocks[k] < v {
						canRun = false
						break
					}
				}

				if canRun {
					for k, v := range p.Needs {
						stocks[k] -= v
					}
					for k, v := range p.Results {
						virtualDemand[k] -= v
						inProgress[k] += v
					}
					
					events = append(events, Event{
						Cycle:   cycle + p.Duration,
						Process: p,
					})
					
					fmt.Printf(" %d:%s\n", cycle, p.Name)
					f.WriteString(fmt.Sprintf("%d:%s\n", cycle, p.Name))
					startedInThisPass = true
					anyStarted = true
					break 
				} else {
					for k, v := range p.Needs {
						totalAvailable := stocks[k] + inProgress[k]
						if totalAvailable < v {
							deficit := v - totalAvailable
							if virtualDemand[k] < deficit {
								virtualDemand[k] = deficit
								demandChanged = true
							}
						}
					}
				}
			}

			if !startedInThisPass && !demandChanged {
				break
			}
		}

		if len(events) == 0 && !anyStarted {
			// To match the original outputs, it prints the cycle after the last event
			fmt.Printf("No more process doable at cycle %d\n", cycle+1)
			f.WriteString(fmt.Sprintf("# No more process doable at cycle %d\n", cycle+1))
			break
		}

		if len(events) > 0 {
			nextCycle := events[0].Cycle
			for _, e := range events {
				if e.Cycle < nextCycle {
					nextCycle = e.Cycle
				}
			}
			cycle = nextCycle
		} else {
			break
		}
	}

	fmt.Println("Stock:")
	
	// Add pending stocks to the final stocks list as per typical scheduler output 
	// Wait, original output doesn't add inProgress to stock. Let's just print stocks[k].
	
	keys := []string{}
	for k := range stocks {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	
	for _, k := range keys {
		fmt.Printf(" %s => %d\n", k, stocks[k])
	}
}
