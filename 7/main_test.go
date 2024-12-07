package main

import "testing"

func Test_SumValid(t *testing.T) {
	testInput := `190: 10 19
3267: 81 40 27
83: 17 5
156: 15 6
7290: 6 8 6 15
161011: 16 10 13
192: 17 8 14
21037: 9 7 18 13
292: 11 6 16 20`
	exampleResult := 3749
	result := int(SumValid(testInput, Valid))
	if result != exampleResult {
		t.Errorf("got %d expected %d", result, exampleResult)
	}
}

func Test_SumValid2(t *testing.T) {
	testInput := `190: 10 19
3267: 81 40 27
83: 17 5
156: 15 6
7290: 6 8 6 15
161011: 16 10 13
192: 17 8 14
21037: 9 7 18 13
292: 11 6 16 20`
	exampleResult := 11387
	result := int(SumValid(testInput, Valid2))
	if result != exampleResult {
		t.Errorf("got %d expected %d", result, exampleResult)
	}
}
