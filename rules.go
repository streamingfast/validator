package validator

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/eoscanada/eos-go"
	"github.com/streamingfast/opaque"
)

type Rule func(field string, rule string, message string, value interface{}) error

func EOSBlockNumRule(field string, rule string, message string, value interface{}) error {
	val, ok := value.(string)
	if !ok {
		return fmt.Errorf("The %s field must be a string", field)
	}

	_, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return fmt.Errorf("The %s field must be a valid EOS block num", field)
	}

	return nil
}

func EOSNameRule(field string, rule string, message string, value interface{}) error {
	checkName := func(field string, name string) error {
		if !IsValidName(name) {
			return fmt.Errorf("The %s field must be a valid EOS name", field)
		}

		return nil
	}

	switch v := value.(type) {
	case string:
		return checkName(field, v)
	case eos.Name, eos.PermissionName, eos.ActionName, eos.AccountName, eos.TableName:
		return checkName(field, fmt.Sprintf("%s", v))
	default:
		return fmt.Errorf("The %s field is not a known type for an EOS name", field)
	}
}

func EOSExtendedNameRule(field string, rule string, message string, value interface{}) error {
	checkName := func(field string, name string) error {
		if !IsValidExtendedName(name) {
			return fmt.Errorf("The %s field must be a valid EOS name", field)
		}

		return nil
	}

	switch v := value.(type) {
	case string:
		return checkName(field, v)
	case eos.Symbol:
		return checkName(field, v.String())
	case eos.SymbolCode:
		return checkName(field, v.String())
	case eos.Name, eos.PermissionName, eos.ActionName, eos.AccountName, eos.TableName:
		return checkName(field, fmt.Sprintf("%s", v))
	default:
		return fmt.Errorf("The %s field is not a known type for an EOS name", field)
	}
}

func EOSNamesListRuleFactory(sep string, maxCount int) Rule {
	return StringListRuleFactory(sep, maxCount, EOSNameRule)
}

func EOSExtendedNamesListRuleFactory(sep string, maxCount int) Rule {
	return StringListRuleFactory(sep, maxCount, EOSExtendedNameRule)
}

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

func EOSTrxIDRule(field string, rule string, message string, value interface{}) error {
	err := HexRowRule(field, rule, message, value)
	if err != nil {
		return err
	}

	val := value.(string)
	if len(val) != 64 {
		return fmt.Errorf("The %s field must have exactly 64 characters", field)
	}

	return nil
}

func CursorRule(field string, rule string, message string, value interface{}) error {
	val, ok := value.(string)
	if !ok {
		return fmt.Errorf("The %s field must be a string", field)
	}

	if val == "" {
		return nil
	}

	_, err := opaque.FromOpaque(val)
	if err != nil {
		return fmt.Errorf("The %s field is not a valid cursor", field)
	}

	return nil
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
