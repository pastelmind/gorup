package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func printHelp() {
	fmt.Println(`Usage: gorup [...names]

Randomly picks a name from a list.

By default, all names have an equal chance (q=1.0) of being picked.
To override the chance for a name, place a '-q' and a number before it.

Examples:
    gorup apple banana pear
        apple : banana : pear = 1 : 1 : 1 chance of being picked

    gorup dog -q 1.2 cat -q0.1 mouse
        dog : cat : mouse = 1 : 1.2 : 0.1 chance of being picked

To see this message, use -h or --help.`)
}

func parseArgs() ([]entry, bool, error) {
	const INITIAL_QUALITY float64 = 1.0

	var entries []entry
	expectQuality := false
	quality := INITIAL_QUALITY

	for _, arg := range os.Args[1:] {
		if expectQuality {
			q, err := strconv.ParseFloat(arg, 64)
			if err != nil {
				return nil, false, fmt.Errorf("Expected a number after '-q', got %s", arg)
			}
			quality = q
			expectQuality = false
			continue
		}

		if strings.HasPrefix(arg, "-") {
			flagBody := arg[1:]
			if strings.HasPrefix(arg, "--") {
				flagBody = arg[2:]
			}

			if flagBody == "" {
				return nil, false, errors.New("Missing flag after '-'")
			}

			if flagBody == "h" || flagBody == "help" {
				return nil, true, nil
			}

			if strings.HasPrefix(flagBody, "q") {
				suffix := flagBody[1:]
				if len(suffix) == 0 {
					expectQuality = true
					continue
				}

				q, err := strconv.ParseFloat(suffix, 64)
				if err != nil {
					return nil, false, fmt.Errorf("Expected a number after '-q', got %s", suffix)
				}
				quality = q
				continue
			}

			return nil, false, fmt.Errorf("Unknown flag: %s", arg)
		} else {
			// Regular argument
			entries = append(entries, entry{name: arg, q: quality})
			quality = INITIAL_QUALITY
		}
	}

	if expectQuality {
		return nil, false, errors.New("No number after last '-q'")
	}

	return entries, false, nil
}
