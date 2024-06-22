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
