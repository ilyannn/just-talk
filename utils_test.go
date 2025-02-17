package main

import (
	"testing"
)

func TestToValidName(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"simpleName", "simpleName"},
		{"invalid chars?!", "invalid_chars__"},
		{"123-abc_DEF", "123-abc_DEF"},
		{
			"thisIsWayTooLongNameThatExceeds64CharactersByAddingALotOfUnnecessaryText12345",
			"thisIsWayTooLongNameThatExceeds64CharactersByAddingALotOfUnneces",
		},
	}

	for _, c := range cases {
		output := toValidName(c.input)
		if output != c.expected {
			t.Errorf("toValidName(%q) = %q; want %q", c.input, output, c.expected)
		}
		if len(output) > 64 {
			t.Errorf("Output length exceeded 64")
		}
	}
}
