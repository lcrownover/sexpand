package main

import (
	"reflect"
	"testing"
)

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

func TestExpandRange(t *testing.T) {
	inputs := []string{
		"n[01-02]",
		"n[0-2]",
		"n[05-07,09]",
		"[05-07,09]",
	}
	expected := [][]string{
		{"n01", "n02"},
		{"n0", "n1", "n2"},
		{"05", "06", "07", "09"},
	}

	for i := range inputs {
		testInput := inputs[i]
		want := expected[i]
		got := expandRange(testInput)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %s, want %s", got, want)
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
