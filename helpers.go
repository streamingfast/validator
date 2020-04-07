package validator

import (
	"regexp"
	"strings"
)

var symbolRegexp = regexp.MustCompile(`^[0-9],[A-Z]{1,7}$`)
var symbolCodeRegexp = regexp.MustCompile(`^[A-Z]{1,7}$`)
var nameRegexp = regexp.MustCompile(`^[\.a-z1-5]{0,13}$`)

func ExplodeNames(input string, sep string) (names []string) {
	rawNames := strings.Split(input, sep)
	for _, rawName := range rawNames {
		account := strings.TrimSpace(rawName)
		if account == "" {
			continue
		}

		names = append(names, rawName)
	}

	return
}

// FIXME: Use eso-go IsValidName once merged, not perfect Regex for now, 13 characters if present is restricted to a different subset
func IsValidName(input string) bool {
	// An empty string name means a uint64 transformed name with a 0 value
	if input == "" {
		return true
	}

	return nameRegexp.MatchString(input)
}

func IsValidExtendedName(input string) bool {
	// An empty string name means a uint64 transformed name with a 0 value
	if input == "" {
		return true
	}

	return nameRegexp.MatchString(input) || symbolCodeRegexp.MatchString(input) || symbolRegexp.MatchString(input)
}
