package generator

import "github.com/sethvargo/go-password/password"

func GeneratePassword(length int, digits int, symbols int, disableUpper bool, allowRepeats bool) string {
	if length <= 0 {
		length = 16
	}

	if digits <= 0 {
		digits = 0
	}

	if symbols <= 0 {
		symbols = 0
	}

	result, _ := password.Generate(length, digits, symbols, disableUpper, allowRepeats)

	return result
}

func GenerateDefault() string {
	result, _ := password.Generate(24, 5, 5, false, false)
	return result
}