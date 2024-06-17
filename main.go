package main

import (
	"fmt"
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

func expandRange(s string) []string {

}

func SExpand(s string) ([]string, error) {
	return []string{}, fmt.Errorf("not implemented")
}

func main() {

}
