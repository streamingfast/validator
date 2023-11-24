package validator

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type ruleTestCase struct {
	name          string
	value         interface{}
	expectedError string
}

func TestDateTimeRule(t *testing.T) {
	tag := "date_time"
	rule := DateTimeRuleFactory(time.RFC3339)
	validator := func(field string, value interface{}) error {
		return rule(field, tag, "", value)
	}

	tests := []ruleTestCase{
		{"should be a string", true, "The test field must be a string"},
		{"should fail on valid layout", "2019-01-12 15:23:34", "The test field is not a valid date time string according to layout 2006-01-02T15:04:05Z07:00"},

		{"valid", "2019-01-12T15:23:34+00:00", ""},
	}

	runRuleTestCases(t, tag, tests, validator)
}

func TestHexRowRule(t *testing.T) {
	tag := "hex"
	validator := func(field string, value interface{}) error {
		return HexRule(field, tag, "", value)
	}

	deprecatedValidator := func(field string, value interface{}) error {
		return HexRowRule(field, tag, "", value)
	}

	tests := []ruleTestCase{
		{"should be a string", true, "The test field must be a string"},
		{"should contains something", "", "The test field must be a valid hexadecimal"},
		{"should contains a least two characters", "a", "The test field must be a valid hexadecimal"},
		{"should not contains invalid characters", "az", "The test field must be a valid hexadecimal"},
		{"should be a multple of 2", "ab01020", "The test field must be a valid hexadecimal"},

		{"valid", "ab", ""},
		{"valid", "1234567890abcdefABCDEF", ""},
	}

	runRuleTestCases(t, tag, tests, validator)
	runRuleTestCases(t, tag+"_deprecated", tests, deprecatedValidator)
}

func TestHexRowsRule(t *testing.T) {
	tag := "hex_slice"
	validator := func(field string, value interface{}) error {
		return HexSliceRule(field, tag, "", value)
	}

	deprecatedValidator := func(field string, value interface{}) error {
		return HexRowsRule(field, tag, "", value)
	}

	tests := []ruleTestCase{
		{"should be an array", "", "The test field must be a string array"},
		{"should have at least 1 row", []string{}, "The test field must have at least 1 element"},
		{"should fail on single error", []string{"a"}, "The test[0] field must be a valid hexadecimal"},
		{"should fail if any row error", []string{"ab", "zz"}, "The test[1] field must be a valid hexadecimal"},

		{"valid single row", []string{"ab"}, ""},
		{"valid multiple rows", []string{"ab", "de"}, ""},
	}

	runRuleTestCases(t, tag, tests, validator)
	runRuleTestCases(t, tag+"_deprecated", tests, deprecatedValidator)
}

func runRuleTestCases(t *testing.T, tag string, tests []ruleTestCase, validator func(field string, value interface{}) error) {
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s_%s", tag, test.name), func(t *testing.T) {
			err := validator("test", test.value)

			if test.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, errors.New(test.expectedError), err)
			}
		})
	}
}
