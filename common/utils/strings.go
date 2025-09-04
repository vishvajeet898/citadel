package utils

import "strings"

func StringsEqualIgnoreCase(s1, s2 string) bool {
	return strings.EqualFold(s1, s2)
}

func ConvertStringToCamelCase(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		words[i] = CapitalizeFirstLetter(word)
	}
	return strings.Join(words, " ")
}

func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	return strings.ToUpper(s[:1]) + s[1:]
}
