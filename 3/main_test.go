package main

import "testing"

func Test_Mul(t *testing.T) {

	tests := []struct {
		name   string
		input  string
		result int
	}{
		{
			name:   "from instruction",
			input:  "xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))",
			result: 161,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := []byte(test.input)
			result := Mul(input)
			if result != test.result {
				t.Errorf("%v expected %v", result, test.result)
			}
		})
	}
}

func Test_MulConditional(t *testing.T) {

	tests := []struct {
		name   string
		input  string
		result int
	}{
		{
			name:   "from instruction",
			input:  "xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))",
			result: 161,
		},

		{
			name:   "from 2 instruction",
			input:  "xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))",
			result: 48,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := []byte(test.input)
			result := MulConditional(input)
			if result != test.result {
				t.Errorf("%v expected %v", result, test.result)
			}
		})
	}
}
