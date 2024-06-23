package main

import (
	"reflect"
	"testing"
)

func TestGroupsToString(t *testing.T) {
	inputs := [][]string{
		{"n01", "n02"},
		{"n[01-02]"},
		{"n[0-2]"},
		{"n[01,02]", "n03", "n[05-07,09]"},
	}
	expected := []string{
		"n01,n02",
		"n[01-02]",
		"n[0-2]",
		"n[01,02],n03,n[05-07,09]",
	}

	for i := range inputs {
		testInput := inputs[i]
		want := expected[i]
		got := groupsToString(testInput)
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	}
}

func TestSplitOutsideRange(t *testing.T) {
	inputs := []string{
		"n01,n02",
		"n[01-02]",
		"n[0-2]",
		"n[01,02],n03,n[05-07,09]",
	}
	expected := [][]string{
		{"n01", "n02"},
		{"n[01-02]"},
		{"n[0-2]"},
		{"n[01,02]", "n03", "n[05-07,09]"},
	}

	for i := range inputs {
		testInput := inputs[i]
		want := expected[i]
		got := splitOutsideRange(testInput)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %s, want %s", got, want)
		}
	}
}

func TestSplitPrefix(t *testing.T) {
	inputs := []string{
		"n[01-02]",
		"[0-2]",
		"np[05-07,09]",
	}
	expected := [][]string{
		{"n", "[01-02]"},
		{"", "[0-2]"},
		{"np", "[05-07,09]"},
	}

	for i := range inputs {
		testInput := inputs[i]
		want := expected[i]
		got0, got1 := splitPrefix(testInput)
		if got0 != want[0] || got1 != want[1] {
			t.Errorf("got (%s, %s), want (%s, %s)", got0, got1, want[0], want[1])
		}
	}
}

func TestExpandRange(t *testing.T) {
	inputs := []string{
		"01-02",
		"0-2",
		"05-07",
		"05-07",
	}
	expected := [][]string{
		{"01", "02"},
		{"0", "1", "2"},
		{"05", "06", "07"},
		{"05", "06", "07"},
	}

	for i := range inputs {
		testInput := inputs[i]
		want := expected[i]
		got, err := expandRange(testInput)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %s, want %s", got, want)
		}
	}
}

func TestExpandGroup(t *testing.T) {
	inputs := [][]string{
		{"n", "[01-02]"},
		{"n", "[0-2]"},
		{"np", "[05-07,09,11,23-25]"},
		{"", "[05-07,09]"},
	}
	expected := [][]string{
		{"n01", "n02"},
		{"n0", "n1", "n2"},
		{"np05", "np06", "np07", "np09", "np11", "np23", "np24", "np25"},
		{"05", "06", "07", "09"},
	}

	for i := range inputs {
		testInput := inputs[i]
		want := expected[i]
		got, err := expandGroup(testInput[0], testInput[1])
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %s, want %s", got, want)
		}
	}
}

func TestCheckFullyExpanded(t *testing.T) {
	inputs := [][]string{
		{"01-02"},
		{"0-2", "3"},
		{"05", "07"},
		{"05"},
		{"05-07,09"},
	}
	expected := []bool{
		false,
		false,
		true,
		true,
		false,
	}
	for i := range inputs {
		testInput := inputs[i]
		want := expected[i]
		got := checkFullyExpanded(testInput)
		if got != want {
			t.Errorf("input %s, got %v, want %v", testInput, got, want)
		}
	}
}

type recurseTestCase struct {
	input0 []string
	input1 string
	input2 string
}

func TestRecurse(t *testing.T) {
	inputs := []recurseTestCase{
		// {[]string{}, "n", "[01-02]"},
		// {[]string{}, "n", "t[05-07]"},
		{[]string{}, "n", "t[05-07,x[10-11]]"},
	}
	expected := [][]string{
		// {"n01", "n02"},
		// {"nt05", "nt06", "nt07"},
		{"nt05", "nt06", "nt07", "ntx10", "ntx11"},
	}

	for i := range inputs {
		testInput := inputs[i]
		want := expected[i]
		got, err := recurse(testInput.input0, testInput.input1, testInput.input2)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	}
}

/*
// expand_hostnames("n01,n02") -> ["n01", "n02"]
assert_eq!(expand_hostnames("n01,n02").unwrap(), ["n01", "n02"]);

// expand_hostnames("n[01-02]") -> ["n01", "n02"]
assert_eq!(expand_hostnames("n[01-02]").unwrap(), ["n01", "n02"]);

// expand_hostnames("n[0-2]") -> ["n0", "n1", "n2"]
assert_eq!(expand_hostnames("n[0-2]").unwrap(), ["n0", "n1", "n2"]);

// expand_hostnames("n[01-05]") -> ["n01", "n02", "n03", "n04", "n05"]
assert_eq!(expand_hostnames("n[01-05]").unwrap(), ["n01", "n02", "n03", "n04", "n05"]);

// expand_hostnames("n[01,02],n03,n[05-07,09]") -> ["n01", "n02", "n03", "n05", "n06", "n07"]
assert_eq!( expand_hostnames("n[01,02],n03,n[05-07,09]").unwrap(), ["n01", "n02", "n03", "n05", "n06", "n07", "n09"]);

// expand_hostnames("n[01,02],n03,n[05-07,09]") -> ["n01", "n02", "n03", "n05", "n06", "n07"]
assert_eq!( expand_hostnames("n[01,02],n03,n[05-07,09]").unwrap(), ["n01", "n02", "n03", "n05", "n06", "n07", "n09"]);

// expand_hostnames("n[[01,02]-03],n[05-07,09]") -> Err
let res = expand_hostnames("n[[01,02]-03],n[05-07,09]");
assert!(res.is_err())
*/

// func TestSExpand(t *testing.T) {
// 	inputs := []string{
// 		"n01,n02",
// 		"n[01-02]",
// 		"n[0-2]",
// 		"n[01-05]",
// 		"n[01,02],n03,n[05-07,09]",
// 		"n[01,02],n03,n[05-07,09]",
// 		"n[[01,02]-03],n[05-07,09]",
// 	}
// 	expected := [][]string{
// 		[]string{"n01", "n02"},
// 		[]string{"n01", "n02"},
// 		[]string{"n0", "n1", "n2"},
// 		[]string{"n01", "n02", "n03", "n04", "n05"},
// 		[]string{"n01", "n02", "n03", "n05", "n06", "n07", "n09"},
// 		[]string{"n01", "n02", "n03", "n05", "n06", "n07", "n09"},
// 		[]string{},
// 	}

// 	for i := range inputs {
// 		testInput := inputs[i]
// 		want := expected[i]
// 		got, err := SExpand(testInput)
// 		if err != nil {
// 			t.Errorf("unexpected error %v. got %s, want %s", got, want)
// 		}
// 		if !reflect.DeepEqual(got, want) {
// 			t.Errorf("got %s, want %s", got, want)
// 		}
// 	}
// }
