package main

type Argument struct {
	Name     string
	Optional bool
	Variadic bool
	Default  string
}

type Recipe struct {
	Name        string
	Description string
	Arguments   []Argument
}
