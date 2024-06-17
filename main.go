package main

import (
	"fmt"
)

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
	fmt.Printf("input string: 		%s\n", s)
	fmt.Printf("comma positions: 	%v\n", outerCommaPositions)
	for i := range outerCommaPositions {
		didSomething := false
		// first element
		if i == 0 {
			outs = append(outs, s[:outerCommaPositions[i]])
			didSomething = true
		}
		// last element
		if i == len(outerCommaPositions)-1 {
			outs = append(outs, s[outerCommaPositions[i+1]:])
			didSomething = true
		}
		if didSomething {
			continue
		}
		outs = append(outs, s[outerCommaPositions[i+1]:outerCommaPositions[i]])
	}
	// if nothing to be done, just return the single string back inside a slice
	if len(outs) == 0 {
		outs = append(outs, s)
	}
	fmt.Printf("outputs:		%v\n\n", outs)
	return outs
}

func SExpand(s string) ([]string, error) {
	return []string{}, fmt.Errorf("not implemented")
}

func main() {

}
