package main

import (
	"testing"
)

func TestParseJustOutput(t *testing.T) {
	raw := `Available recipes:
    bootstrap           # Bootstrap Kibana.
    ensure-node-version # Update node version with asdf if necessary.
    command arg *varg   # Command description.
    with-default arg=foo # Command with default value.
`
	recipes := parseJustOutput(raw)

	if len(recipes) != 4 {
		t.Fatalf("Expected 4 recipes, got %d", len(recipes))
	}

	if recipes[0].Name != "bootstrap" {
		t.Errorf("Expected bootstrap, got %s", recipes[0].Name)
	}

	if recipes[2].Name != "command" {
		t.Errorf("Expected command, got %s", recipes[2].Name)
	}

	arguments := recipes[2].Arguments

	if len(arguments) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(arguments))
	}

	if arguments[0].Name != "arg" {
		t.Errorf("Expected arg, got %s", arguments[0].Name)
	}

	if arguments[1].Name != "varg" {
		t.Errorf("Expected varg, got %s", arguments[1].Name)
	}

	if arguments[1].Variadic != true {
		t.Errorf("Expected varg to be variadic")
	}

	if arguments[1].Optional != true {
		t.Errorf("Expected varg to be optional")
	}

	arguments = recipes[3].Arguments

	if len(arguments) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(arguments))
	}

	if arguments[0].Name != "arg" {
		t.Errorf("Expected arg, got %s", arguments[0].Name)
	}

	if arguments[0].Default != "foo" {
		t.Errorf("Expected foo, got %s", arguments[0].Default)
	}

	if arguments[0].Optional != true {
		t.Errorf("Expected arg to be optional")
	}
}
