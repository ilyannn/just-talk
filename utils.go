package main

import "strings"

//	 The name of the function to be called. Must be a-z, A-Z, 0-9, or contain underscores and dashes, with a maximum
//		length of 64.
func toValidName(input string) string {
	var sb strings.Builder
	for _, r := range input {
		switch {
		case (r >= '0' && r <= '9') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= 'a' && r <= 'z') ||
			r == '_' || r == '-':
			sb.WriteRune(r)
		default:
			sb.WriteRune('_')
		}
	}
	result := sb.String()
	if len(result) > 64 {
		result = result[:64]
	}
	return result
}
