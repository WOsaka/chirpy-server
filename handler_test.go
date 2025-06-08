package main

import "testing"

func TestReplaceProfane(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"Hello, world!", "Hello, world!"},
		{"This is a test with a Kerfuffle.", "This is a test with a ****."},
		{"No bad words here.", "No bad words here."},
		{"Another sharbert in the text.", "Another **** in the text."},
	}

	for _, test := range tests {
		result := replaceProfane(test.input)
		if result != test.expected {
			t.Errorf("replaceProfane(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}