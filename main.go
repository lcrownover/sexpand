package main

import (
	"fmt"
	"strconv"
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

// splitPrefix takes a SLURM group range notation and returns the prefix and the notation separately
// this function expects that there is a range notation in the group and will panic if not
func splitPrefix(group string) (string, string) {
	prefix := ""
	groupRange := ""
	beforePrefix := true

	if !strings.Contains(group, "[") {
		beforePrefix = false
	}
	for _, c := range group {
		if c == '[' {
			beforePrefix = false
		}
		if beforePrefix {
			prefix = prefix + string(c)
		} else {
			groupRange = groupRange + string(c)
		}
	}
	return prefix, groupRange
}

// expandRange expands a single bracketed range notation into a
// list of strings
func expandRange(r string) ([]string, error) {
	s := strings.Split(r, "-")
	groupStart := s[0]
	groupEnd := s[1]
	rangeLength := len(groupStart)
	// now convert start and end to integers
	startNum, err := strconv.Atoi(groupStart)
	if err != nil {
		return nil, fmt.Errorf("invalid number used in range: %d", startNum)
	}
	endNum, err := strconv.Atoi(groupEnd)
	if err != nil {
		return nil, fmt.Errorf("invalid number used in range: %d", endNum)
	}
	// validate that the start num comes before the end num for a valid range
	if startNum >= endNum {
		return nil, fmt.Errorf("range start must be smaller than range end: [%d-%d]", startNum, endNum)
	}
	// generate the list of integers
	nums := []int{}
	for i := startNum; i <= endNum; i++ {
		nums = append(nums, i)
	}
	// reformat that list to strings of the correct padded length
	// from the original length
	numStrs := []string{}
	for _, n := range nums {
		paddedNum := fmt.Sprintf("%0*d", rangeLength, n)
		numStrs = append(numStrs, paddedNum)
	}
	return numStrs, nil
}

// expandGroup expands a single prefix and SLURM range notation
// into a list of strings
func expandGroup(prefix string, group string) ([]string, error) {
	// first we find the length of runes for the start of the range
	group = strings.ReplaceAll(group, "[", "")
	group = strings.ReplaceAll(group, "]", "")
	expandedParts := []string{}
	csParts := strings.Split(group, ",")
	for _, p := range csParts {
		if strings.Contains(p, "-") {
			expandedRange, err := expandRange(p)
			if err != nil {
				return nil, fmt.Errorf("invalid range: %s", p)
			}
			expandedParts = append(expandedParts, expandedRange...)
		} else {
			// simple part, no range notation "-"
			expandedParts = append(expandedParts, p)
		}
	}
	// add the prefix to all of them
	outs := []string{}
	for _, n := range expandedParts {
		outs = append(outs, fmt.Sprintf("%s%s", prefix, n))
	}
	return outs, nil
}

// checkFullyExpanded returns true if no elements in the provided list of stirngs
// contain a comma or brackets
func checkFullyExpanded(l []string) bool {
	for _, g := range l {
		if strings.Contains(g, ",") {
			return false
		}
		if strings.Contains(g, "-") {
			return false
		}
		if strings.Contains(g, "[") {
			return false
		}
		if strings.Contains(g, "]") {
			return false
		}
	}
	return true
}

// reverse returns its argument string reversed rune-wise left to right.
func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// unwrapRange removes the leading and trailing brackets and
// returns the resulting string.
func unwrapRange(r string) string {
	r = strings.Replace(r, "[", "", 1)
	r = reverse(r)
	r = strings.Replace(r, "]", "", 1)
	r = reverse(r)
	return r
}

// readyToExpand returns true if the group string does not contain any brackets
// which signifies that the range is ready to be expanded.
func readyToExpand(g string) bool {
	return !strings.Contains(g, "[")
}

// recurse is the main recursive runner for expanding the range notation.
// it should receive the outer prefix and the current group notation.
// it returns:
// prefix string
func recurse(nodes []string, prefix string, group string) ([]string, error) {
	var err error
	//
	// base case: 	group provided has no brackets, we're read to expand
	//
	if readyToExpand(group) {
		expandedNodes, err := expandGroup(prefix, group)
		if err != nil {
			return []string{}, err
		}
		nodes = append(nodes, expandedNodes...)
		return nodes, nil
	}

	newGroups := splitOutsideRange(group)
	for _, g := range newGroups {
		np, ng := splitPrefix(g)
		newPrefix := fmt.Sprintf("%s%s", prefix, np)
		newGroup := unwrapRange(ng)
		nodes, err = recurse(nodes, newPrefix, newGroup)
		if err != nil {
			return []string{}, err
		}
	}

	return nodes, nil
}

func SExpand(s string) ([]string, error) {

	return []string{}, fmt.Errorf("not implemented")
}

func main() {

}
