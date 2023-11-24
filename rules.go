package validator

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Rule func(field string, rule string, message string, value interface{}) error

func StringListRuleFactory(sep string, maxCount int, elementRule Rule) Rule {
	return func(field string, rule string, message string, value interface{}) error {
		rawNames, ok := value.(string)
		if !ok {
			return fmt.Errorf("The %s field must be a string", field)
		}

		names := ExplodeNames(rawNames, sep)
		nameCount := len(names)
		if nameCount <= 0 {
			return fmt.Errorf("The %s field must have at least 1 element", field)
		}

		if nameCount > maxCount {
			return fmt.Errorf("The %s field must have at most %d elements", field, maxCount)
		}

		for i, name := range names {
			err := elementRule(fmt.Sprintf("%s[%d]", field, i), rule, message, name)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func DateTimeRuleFactory(layout string) Rule {
	return func(field string, rule string, message string, value interface{}) error {
		val, ok := value.(string)
		if !ok {
			return fmt.Errorf("The %s field must be a string", field)
		}

		_, err := time.Parse(layout, val)
		if err != nil {
			return fmt.Errorf("The %s field is not a valid date time string according to layout %s", field, layout)
		}

		return nil
	}
}

// Deprecated: Use `HexRule` instead
var HexRowRule = HexRule

func HexRule(field string, rule string, message string, value interface{}) error {
	hexRow, ok := value.(string)
	if !ok {
		return fmt.Errorf("The %s field must be a string", field)
	}

	match, _ := regexp.MatchString("^[A-Fa-f0-9]+$", hexRow)
	if !match {
		return fmt.Errorf("The %s field must be a valid hexadecimal", field)
	}

	if len(hexRow)%2 != 0 {
		return fmt.Errorf("The %s field must be a valid hexadecimal", field)
	}

	return nil
}

// Deprecated: Use `HexRowsRule` instead
var HexRowsRule = HexSliceRule

func HexSliceRule(field string, rule string, message string, value interface{}) error {
	hexRows, ok := value.([]string)
	if !ok {
		return fmt.Errorf("The %s field must be a string array", field)
	}

	if len(hexRows) <= 0 {
		return fmt.Errorf("The %s field must have at least 1 element", field)
	}

	for i, hexData := range hexRows {
		err := HexRowRule(fmt.Sprintf("%s[%d]", field, i), rule, message, hexData)
		if err != nil {
			return err
		}
	}

	return nil
}

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
