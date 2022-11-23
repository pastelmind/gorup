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

// Constant error type.
// See https://dave.cheney.net/2016/04/07/constant-errors for an explanation.
type constError string

func (e constError) Error() string { return string(e) }

// Constant error that indicates the user requested a help message
const eHelp = constError("Help")

func parseEntries() ([]entry, error) {
	args, err := parseArgs()
	if err != nil {
		return nil, err
	}

	var entries []entry

	const INITIAL_QUALITY float64 = 1.0
	quality := INITIAL_QUALITY
	for _, arg := range args {
		switch a := arg.(type) {
		case aName:
			entries = append(entries, entry{name: string(a), q: quality})
			quality = INITIAL_QUALITY
		case aQuality:
			quality = float64(a)
		default:
			panic(fmt.Errorf("Unhandled cliArg: %T", a))
		}
	}

	return entries, nil
}

func parseArgs() ([]cliArg, error) {
	var args []cliArg

	expectQuality := false
	for _, arg := range os.Args[1:] {
		if expectQuality {
			q, err := strconv.ParseFloat(arg, 64)
			if err != nil {
				return nil, fmt.Errorf("Expected a number after '-q', got %s", arg)
			}
			args = append(args, aQuality(q))
			expectQuality = false
			continue
		}

		if strings.HasPrefix(arg, "-") {
			flagBody := arg[1:]
			if strings.HasPrefix(arg, "--") {
				flagBody = arg[2:]
			}

			if flagBody == "" {
				return nil, errors.New("Missing flag after '-'")
			}

			if flagBody == "h" || flagBody == "help" {
				return nil, eHelp
			}

			if strings.HasPrefix(flagBody, "q") {
				suffix := flagBody[1:]
				if len(suffix) == 0 {
					expectQuality = true
					continue
				}

				q, err := strconv.ParseFloat(suffix, 64)
				if err != nil {
					return nil, fmt.Errorf("Expected a number after '-q', got %s", suffix)
				}
				args = append(args, aQuality(q))
				continue
			}

			return nil, fmt.Errorf("Unknown flag: %s", arg)
		} else {
			// Regular argument (i.e. name)
			args = append(args, aName(arg))
		}
	}

	if expectQuality {
		return nil, errors.New("No number after last '-q'")
	}

	return args, nil
}

// Dummy interface for sum type of CLI arguments
type cliArg interface {
	isCliArg()
}

type aName string
type aQuality float64

func (aName) isCliArg()    {}
func (aQuality) isCliArg() {}
