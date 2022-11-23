package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	os.Exit(run())
}

func run() int {
	entries, showHelp, err := parseArgs()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Pass -h or --help to see usage")
		return 1
	}

	if showHelp || len(entries) == 0 {
		printHelp()
		return 0
	}

	e, roll, err := pick(entries)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	// Print chance to pick each item
	var total float64
	for _, e := range entries {
		total += e.q
	}
	for _, e := range entries {
		fmt.Printf("%-20s: %7.3f%%\n", e.name, e.q/total*100)
	}

	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("%s (rolled %f)\n", e.name, roll)
	return 0
}

type entry struct {
	name string
	q    float64
}
