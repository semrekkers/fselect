package qbuilder

import "testing"

func TestSliceContains(t *testing.T) {
	slice := []string{"aaa", "bbb", "ccc", "abc", "cba"}

	if sliceContains("abcd", slice) {
		t.Error(`"abcd" doesn't exist in slice`)
	}
	if !sliceContains("cba", slice) {
		t.Error(`"cba" does exist in slice`)
	}
}

func TestRepeatString(t *testing.T) {
	const expect = "test, test, test, test, test, test"

	if repeatString("test", ", ", 6) != expect {
		t.Error("failed")
	}
}

func TestJoinStringsWithSuffix(t *testing.T) {
	const expect = "test and succeed, play and succeed, program and succeed"

	set := []string{"test", "play", "program"}
	if joinStringsWithSuffix(set, " and succeed", ", ") != expect {
		t.Error("failed")
	}
}

func TestJoinTwoSlices(t *testing.T) {
	const expect = "a: 1, b: 2, c: 3, d: 4"
	a := []string{"a", "b", "c", "d"}
	b := []string{"1", "2", "3", "4"}

	if joinTwoSlices(a, b, ": ", ", ") != expect {
		t.Error("failed")
	}
}
