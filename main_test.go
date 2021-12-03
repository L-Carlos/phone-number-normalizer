package main

import "testing"

var normalizeTestCases = []struct {
	input string
	want  string
}{
	{"1234567890", "1234567890"},
	{"123 456 7891", "1234567891"},
	{"(123) 456 7892", "1234567892"},
	{"(123) 456-7893", "1234567893"},
	{"123-456-7894", "1234567894"},
	{"123-456-7890", "1234567890"},
	{"1234567892", "1234567892"},
	{"(123)456-7892", "1234567892"},
}

func TestNormalize(t *testing.T) {
	for _, tc := range normalizeTestCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := normalize(tc.input)
			if actual != tc.want {
				t.Errorf("got %s; want %s", actual, tc.want)
			}
		})
	}
}

func TestNormalizeRegex(t *testing.T) {
	for _, tc := range normalizeTestCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := normalizeRegex(tc.input)
			if actual != tc.want {
				t.Errorf("got %s; want %s", actual, tc.want)
			}
		})
	}
}
