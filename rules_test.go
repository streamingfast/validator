package validator

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/eoscanada/eos-go"
	"github.com/stretchr/testify/assert"
)

type ruleTestCase struct {
	name          string
	value         interface{}
	expectedError string
}

func TestEOSBlockNumRule(t *testing.T) {
	tag := "eos_block_num"
	validator := func(field string, value interface{}) error {
		return EOSBlockNumRule(field, tag, "", value)
	}

	tests := []ruleTestCase{
		{"should be a string", true, "The test field must be a string"},
		{"should not contains invalid characters", "!", "The test field must be a valid EOS block num"},

		{"valid block num", "10", ""},
	}

	runRuleTestCases(t, tag, tests, validator)
}

func TestEOSNameRule(t *testing.T) {
	tag := "eos_name"
	validator := func(field string, value interface{}) error {
		return EOSNameRule(field, tag, "", value)
	}

	tests := []ruleTestCase{
		{"should be a string", true, "The test field is not a known type for an EOS name"},
		{"should not contains invalid characters", "6", "The test field must be a valid EOS name"},
		{"should not be longer than 13", "abcdefghigklma", "The test field must be a valid EOS name"},

		{"valid empty", "", ""},
		{"valid single", "e", ""},
		{"valid limit", "5", ""},
		{"valid with dots and 13 chars", "eosio.tokenfl", ""},
		{"valid eos.Name", eos.Name("eosio"), ""},
		{"valid eos.PermissionName", eos.PermissionName("eosio"), ""},
		{"valid eos.ActionName", eos.ActionName("eosio"), ""},
		{"valid eos.AccountName", eos.AccountName("eosio"), ""},
		{"valid eos.TableName", eos.TableName("eosio"), ""},
	}

	runRuleTestCases(t, tag, tests, validator)
}

func TestEOSExtendedNameRule(t *testing.T) {
	tag := "eos_extended_name"
	validator := func(field string, value interface{}) error {
		return EOSExtendedNameRule(field, tag, "", value)
	}

	tests := []ruleTestCase{
		{"should be a string", true, "The test field is not a known type for an EOS name"},
		{"should not contains invalid characters", "6", "The test field must be a valid EOS name"},
		{"should not be longer than 13", "abcdefghigklma", "The test field must be a valid EOS name"},

		{"valid empty", "", ""},
		{"valid single", "e", ""},
		{"valid limit", "5", ""},
		{"valid with dots and 13 chars", "eosio.tokenfl", ""},
		{"valid with whem symbol", "4,EOS", ""},
		{"valid with whem symbol code", "EOS", ""},

		{"valid eos.Name", eos.Name("eosio"), ""},
		{"valid eos.PermissionName", eos.PermissionName("eosio"), ""},
		{"valid eos.ActionName", eos.ActionName("eosio"), ""},
		{"valid eos.AccountName", eos.AccountName("eosio"), ""},
		{"valid eos.TableName", eos.TableName("eosio"), ""},
	}

	runRuleTestCases(t, tag, tests, validator)
}

func TestEOSNamesListRule(t *testing.T) {
	tag := "eos_names_list"
	rule := EOSNamesListRuleFactory("|", 2)
	validator := func(field string, value interface{}) error {
		return rule(field, tag, "", value)
	}

	tests := []ruleTestCase{
		{"should be a string", true, "The test field must be a string"},
		{"should have at least 1 element", "", "The test field must have at least 1 element"},
		{"should have at max macCount element", "eos|eos|eos", "The test field must have at most 2 elements"},
		{"should fail on single error", "6", "The test[0] field must be a valid EOS name"},
		{"should fail if any element error", "ab|6", "The test[1] field must be a valid EOS name"},

		{"valid single", "ab", ""},
		{"valid multiple", "ded|eos", ""},
	}

	runRuleTestCases(t, tag, tests, validator)
}

func TestEOSExtendedNamesListRule(t *testing.T) {
	tag := "eos_extended_names_list"
	rule := EOSExtendedNamesListRuleFactory("|", 3)
	validator := func(field string, value interface{}) error {
		return rule(field, tag, "", value)
	}

	tests := []ruleTestCase{
		{"should be a string", true, "The test field must be a string"},
		{"should have at least 1 element", "", "The test field must have at least 1 element"},
		{"should have at max macCount element", "eos|eos|eos|eos", "The test field must have at most 3 elements"},
		{"should fail on single error", "6", "The test[0] field must be a valid EOS name"},
		{"should fail if any element error", "ab|6", "The test[1] field must be a valid EOS name"},

		{"valid single", "ab", ""},
		{"valid multiple", "ded|eos", ""},
		{"valid multiple symbol", "ded|eos|EOS", ""},
		{"valid multiple symbol code", "ded|4,EOS", ""},
		{"valid multiple mixed", "ded|EOS|4,EOS", ""},
	}

	runRuleTestCases(t, tag, tests, validator)
}

func TestEOSTrxIDRule(t *testing.T) {
	tag := "eos_trx_id"
	validator := func(field string, value interface{}) error {
		return EOSTrxIDRule(field, tag, "", value)
	}

	tests := []ruleTestCase{
		{"should be a string", true, "The test field must be a string"},
		{"should contains something", "", "The test field must be a valid hexadecimal"},
		{"should contains a least two characters", "a", "The test field must be a valid hexadecimal"},
		{"should not contains invalid characters", "az", "The test field must be a valid hexadecimal"},
		{"should be a multple of 2", "ab01020", "The test field must be a valid hexadecimal"},
		{"should be long enough", "d8fe02221408fbcc221d1207c1b8cc67e0d9b3ca1c6005a36ea10428dd7fd1", "The test field must have exactly 64 characters"},

		{"valid", "d8fe02221408fbcc221d1207c1b8cc67e0d9b3ca1c6005a36ea10428dd7fd148", ""},
		{"valid", "D8FE02221408FBCC221D1207C1B8CC67E0D9B3CA1C6005A36EA10428DD7FD148", ""},
	}

	runRuleTestCases(t, tag, tests, validator)
}

func TestCursorRule(t *testing.T) {
	tag := "cursor"
	validator := func(field string, value interface{}) error {
		return CursorRule(field, tag, "", value)
	}

	tests := []ruleTestCase{
		{"happy path", "Wf_IQ72XbdmObmHnniHTKPazJ8IwBwxqBl3tfhdIh4z19XLF2p6hU2N9PUzZla_yjhLjTQis29jKHC9_ocZY7dDuyr9g73JpQS8pxYjp-eflePPybA==", ""},
		{"empty cursor", "", ""},
		{"invalid characters in cursor", "-----==", "The test field is not a valid cursor"},
		{"invalid cursor", "abc", "The test field is not a valid cursor"},
	}

	runRuleTestCases(t, tag, tests, validator)
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
