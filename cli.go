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
		if strings.HasPrefix(arg, "-") {
			flag, err := parseFlag(arg)
			if err != nil {
				return nil, false, err
			}

			switch v := flag.(type) {
			case fHelp:
				return nil, true, nil
			case fQualityBegin:
				expectQuality = true
			case fQuality:
				quality = float64(v)
			default:
				return nil, false, fmt.Errorf("Unknown flag: %s", arg)
			}
		} else {
			if expectQuality {
				q, err := strconv.ParseFloat(arg, 64)
				if err != nil {
					return nil, false, fmt.Errorf("Expected a number after '-q', got %s", arg)
				}
				quality = q
				expectQuality = false
			} else {
				// Regular argument
				entries = append(entries, entry{name: arg, q: quality})
				quality = INITIAL_QUALITY
			}
		}
	}

	if expectQuality {
		return nil, false, errors.New("No number after last '-q'")
	}

	return entries, false, nil
}

func parseFlag(flag string) (flagType, error) {
	var flagBody string
	if strings.HasPrefix(flag, "--") {
		flagBody = flag[2:]
	} else {
		flagBody = flag[1:]
	}

	if flagBody == "" {
		return nil, errors.New("Missing flag after '-'")
	}

	if flagBody == "h" || flagBody == "help" {
		return fHelp(struct{}{}), nil
	}

	if flagBody == "q" {
		return fQualityBegin(struct{}{}), nil
	}

	if strings.HasPrefix(flagBody, "q") {
		quality, err := strconv.ParseFloat(flagBody[1:], 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid suffix for -q flag: %w", err)
		}

		return fQuality(quality), nil
	}

	return fUnknown(struct{}{}), nil
}

// Interface for CLI flags
type flagType interface {
	// Dummy method for identifying flags
	isFlag()
}

// Flag type for help (-h, --help).
type fHelp struct{}

// Flag type for -q without a number suffix.
type fQualityBegin struct{}

// Flag type for -q with a number suffix (ex: -q1, -q2.5).
// The contained value is the number
type fQuality float64

// Flag type for unknown flags.
type fUnknown struct{}

func (fHelp) isFlag()         {}
func (fQualityBegin) isFlag() {}
func (fQuality) isFlag()      {}
func (fUnknown) isFlag()      {}
