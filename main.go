package main

import (
	"fmt"
	"strings"
)

// splitOutsideRange splits a string on all commas that are not located
// inside bracket notations.
func splitOutsideRange(s string) []string {
	nestLevel := 0
	outerCommaPositions := []int{}
	outs := []string{}
	for i, c := range s {
		switch c {
		case '[':
			nestLevel += 1
		case ']':
			nestLevel -= 1
		case ',':
			if nestLevel == 0 {
				outerCommaPositions = append(outerCommaPositions, i)
			}
		}
	}
	for i := range outerCommaPositions {
		// first element
		if i == 0 {
			outs = append(outs, s[:outerCommaPositions[i]])
			continue
		}
		// middle elements
		outs = append(outs, s[outerCommaPositions[i-1]+1:outerCommaPositions[i]])
	}
	// last element
	if len(outerCommaPositions) > 0 {
		outs = append(outs, s[outerCommaPositions[len(outerCommaPositions)-1]+1:])
	}
	// if nothing to be done, just return the single string back inside a slice
	if len(outs) == 0 {
		outs = append(outs, s)
	}
	return outs
}

// expandRange expands a single SLURM range notation into a list of strings
func expandRange(s string) []string {
	prefix := strings.Split(s, "[")[0]
	fmt.Printf("prefix:		%s\n", prefix)
	return nil
}

func SExpand(s string) ([]string, error) {
	return []string{}, fmt.Errorf("not implemented")
}

func main() {

}
