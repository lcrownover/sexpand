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
	if !strings.Contains(group, "[") {
		panic("splitPrefix group must contain range notation. check before using.")
	}
	prefix := ""
	groupRange := ""
	beforePrefix := true
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

// groupsToString simply joins the groups into a comma-separated string
func groupsToString(group []string) string {
	return strings.Join(group, ",")
}

// recurse is the main recursive runner for expanding the range notation.
// it should receive the outer prefix and the current group notation.
// it returns:
// prefix string
func recurse(allGroups []string, outerPrefix string, groupString string) ([]string, error) {
	fmt.Printf("call to recurse\n")
	fmt.Printf("	allGroups:		%v\n", allGroups)
	fmt.Printf("	outerPrefix:		%s\n", outerPrefix)
	fmt.Printf("	groupString:		%s\n", groupString)
	newGroupStrings := splitOutsideRange(groupString)
	fmt.Printf("	newGroupString:		%v\n", newGroupStrings)
	// base case: 	group provided is fully expanded
	// 				so we know that allGroups should contain everuthing
	if checkFullyExpanded(newGroupStrings) {
		fmt.Printf("fully expanded:\n")
		fmt.Printf("	groupString: 		%s\n", groupString)
		nodes := strings.Split(groupString, ",")
		for _, n := range nodes {
			fmt.Printf("	node: %s\n", n)
			n = strings.TrimSpace(n)
			if n != "" {
				allGroups = append(allGroups, n)
			}
		}
		return allGroups, nil
	}
	for _, ngs := range newGroupStrings {
		fmt.Printf("	processing gs:		%s\n", ngs)
		// if it's not fully expanded
		// add the outer prefix to the inner prefix
		prefix, ngs := splitPrefix(ngs)
		prefix = fmt.Sprintf("%s%s", outerPrefix, prefix)
		// then expand the group using the new larger prefix
		newGroups, err := expandGroup(prefix, ngs)
		if err != nil {
			return []string{}, nil
		}
		newGroupsString := groupsToString(newGroups)
		newAllGroups, err := recurse(allGroups, prefix, newGroupsString)
		if err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		for _,g := range newAllGroups {
			fmt.Printf("	adding: %s\n", g)
			allGroups = append(allGroups, g)
		}
		// allGroups = append(allGroups, newAllGroups...)
	}
	return allGroups, nil
}

func SExpand(s string) ([]string, error) {

	return []string{}, fmt.Errorf("not implemented")
}

func main() {

}
