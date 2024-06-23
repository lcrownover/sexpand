package main

import (
	"fmt"
	"strconv"
	"strings"
)

// splitOutsideRange splits a string on all commas that are not located
// inside bracket notations.
func splitOutsideRange(s string) []string {
	fmt.Printf("call splitOutsideRange\n")
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
	fmt.Printf("done calling splitOutsideRange\n")
	return outs
}

// splitPrefix takes a SLURM group range notation and returns the prefix and the notation separately
// this function expects that there is a range notation in the group and will panic if not
func splitPrefix(group string) (string, string) {
	fmt.Printf("call to split prefix:\n")
	fmt.Printf("	before group: 	%s\n", group)
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
	fmt.Printf("	after prefix: 	%s\n", prefix)
	fmt.Printf("	after group: 	%s\n", groupRange)
	fmt.Printf("done calling split prefix\n")
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
	fmt.Printf("calling expandGroup\n")
	fmt.Printf("	prefix: %s\n", prefix)
	fmt.Printf("	group: %s\n", group)
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
	fmt.Printf("done calling expandGroup\n")
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

// groupsToString simply joins the groups into a comma-separated string
func groupsToString(group []string) string {
	return strings.Join(group, ",")
}

// reverse returns its argument string reversed rune-wise left to right.
func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func unwrapRange(r string) string {
	fmt.Printf("	before unwrap: 	%s\n", r)
	r = strings.Replace(r, "[", "", 1)
	r = reverse(r)
	r = strings.Replace(r, "]", "", 1)
	r = reverse(r)
	fmt.Printf("	after unwrap: 	%s\n", r)
	return r
}

// func readyToExpand(r string) bool {
// 	nr := splitOutsideRange(r)
// 	return nr[0] == r && !needsToBeSplit(r)
// }

func needsToBeSplit(r string) bool {
	return strings.Contains(r, ",")
}

func readyToExpand(g string) bool {
	return !strings.Contains(g, "[")
}

// recurse is the main recursive runner for expanding the range notation.
// it should receive the outer prefix and the current group notation.
// it returns:
// prefix string
func recurse(nodes []string, prefix string, group string) ([]string, error) {
	var err error
	fmt.Printf("call to recurse\n")
	fmt.Printf("	nodes:		%v\n", nodes)
	fmt.Printf("	prefix:		%s\n", prefix)
	fmt.Printf("	groupString:	%s\n", group)

	//
	// base case: 	group provided has no brackets, we're read to expand
	//
	if readyToExpand(group) {
		fmt.Printf("ready to expand:\n")
		fmt.Printf("	prefix: 	%s\n", prefix)
		fmt.Printf("	group: 		%s\n", group)
		expandedNodes, err := expandGroup(prefix, group)
		fmt.Printf("	expanded: 	%v\n", expandedNodes)
		if err != nil {
			return []string{}, err
		}
		nodes = append(nodes, expandedNodes...)
		fmt.Printf("base case done\n")
		return nodes, nil
	}

	newGroups := splitOutsideRange(group)
	fmt.Printf("	newGroups:	%v\n", newGroups)

	for _, g := range newGroups {
		fmt.Printf("	old prefix: 	%s\n", prefix)
		np, ng := splitPrefix(g)
		newPrefix := fmt.Sprintf("%s%s", prefix, np)
		newGroup := unwrapRange(ng)
		fmt.Printf("	new prefix: 	%s\n", newPrefix)
		fmt.Printf("	new group: 	%s\n", newGroup)
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
