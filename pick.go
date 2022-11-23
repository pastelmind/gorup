package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func pick(entries []entry) (entry, float64, error) {
	if len(entries) == 0 {
		return entry{}, 0, errors.New("No entries found")
	}

	// Compute the sum of all q-values
	var total float64
	for _, e := range entries {
		if !(e.q > 0) {
			return entry{}, 0, fmt.Errorf("Q-value must be a positive number; got '%f' for '%s'", e.q, e.name)
		}
		total += e.q
	}

	// Pick a random number between [0, total)
	rand.Seed(time.Now().UnixNano())
	var roll = total * rand.Float64()

	// Pick the entry that matches the roll
	var max float64
	for _, e := range entries {
		max += e.q
		if roll < max {
			return e, roll, nil
		}
	}

	panic(fmt.Errorf("Roll (%f) exceeded max (%f)", roll, max))
}
