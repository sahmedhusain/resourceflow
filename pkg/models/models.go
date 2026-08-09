package models

type Process struct {
	Name     string
	Needs    map[string]int
	Results  map[string]int
	Duration int
}

type OptimizeGoal struct {
	Time   bool
	Stocks []string
}

type Config struct {
	InitialStocks map[string]int
	Processes     []Process
	Optimize      OptimizeGoal
}
