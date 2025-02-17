package main

import (
	"strings"
)

func parseJustOutput(raw string) []Recipe {
	var recipes []Recipe
	lines := strings.Split(raw, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.Contains(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "#", 2)
		left := strings.Fields(parts[0])
		description := ""
		if len(parts) > 1 {
			description = strings.TrimSpace(parts[1])
		}

		if len(left) == 0 {
			continue
		}

		name := left[0]
		argTokens := left[1:]

		var arguments []Argument
		for _, arg := range argTokens {
			optional := false
			variadic := false
			defaultValue := ""

			if strings.HasPrefix(arg, "+") {
				variadic = true
				arg = arg[1:]
			} else if strings.HasPrefix(arg, "*") {
				variadic = true
				optional = true
				arg = arg[1:]
			}

			if strings.Contains(arg, "=") {
				parts := strings.SplitN(arg, "=", 2)
				arg = parts[0]
				defaultValue = parts[1]
				optional = true
			}

			arguments = append(arguments, Argument{
				Name:     arg,
				Optional: optional,
				Variadic: variadic,
				Default:  defaultValue,
			})
		}

		recipes = append(recipes, Recipe{
			Name:        name,
			Description: description,
			Arguments:   arguments,
		})
	}

	return recipes
}
