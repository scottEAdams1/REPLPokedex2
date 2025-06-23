package main

import (
	"fmt"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander    Bulbasaur PIKACHU      ",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "       ",
			expected: []string{},
		},
		{
			input:    "hello",
			expected: []string{"hello"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			fmt.Printf("Expected = %s, Actual = %s\n", expectedWord, word)
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("%s does not match %s\n", word, expectedWord)
			}
		}
	}
}
