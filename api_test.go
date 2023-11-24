package validator

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedevsaddam/govalidator"
)

func init() {
	govalidator.AddCustomRule("date_time", DateTimeRuleFactory("2006-01-02"))
}

func TestValidateQueryParams(t *testing.T) {
	singleRules := map[string][]string{
		"block_num": {"date_time"},
	}

	tests := []struct {
		name   string
		query  string
		rules  Rules
		errors url.Values
	}{
		{"block_num valid", "block_num=2017-10-30", singleRules, url.Values{}},
		{"block_num not valid", "block_num=2017", singleRules, url.Values{
			"block_num": []string{"The block_num field is not a valid date time string according to layout 2006-01-02"},
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest("GET", "/?"+test.query, nil)
			require.NoError(t, err)

			errors := ValidateQueryParams(request, test.rules)
			assert.Equal(t, test.errors, errors)
		})
	}
}

func TestValidateStruct(t *testing.T) {
	singleRules := map[string][]string{
		"account": {"date_time"},
	}

	type payload struct {
		Account string `json:"account"`
	}

	tests := []struct {
		name         string
		rules        Rules
		expectedData payload
		errors       url.Values
	}{
		{"account valid", singleRules, payload{Account: "2017-02-10"}, url.Values{}},
		{"account not valid", singleRules, payload{Account: "6"}, url.Values{
			"account": []string{"The account field is not a valid date time string according to layout 2006-01-02"},
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errors := ValidateStruct(&test.expectedData, test.rules)
			assert.Equal(t, test.errors, errors)
		})
	}
}

func TestValidateStruct_CustomTagIdentifier(t *testing.T) {
	singleRules := map[string][]string{
		"account": {"date_time"},
	}

	type payload struct {
		Account string `schema:"account"`
	}

	tests := []struct {
		name         string
		rules        Rules
		expectedData payload
		errors       url.Values
	}{
		{"account valid", singleRules, payload{Account: "2017-02-10"}, url.Values{}},
		{"account not valid", singleRules, payload{Account: "6"}, url.Values{
			"account": []string{"The account field is not a valid date time string according to layout 2006-01-02"},
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errors := ValidateStruct(&test.expectedData, test.rules, TagIdentifierOption("schema"))
			assert.Equal(t, test.errors, errors)
		})
	}
}

func TestValidateJSONBody(t *testing.T) {
	singleRules := map[string][]string{
		"account": {"date_time"},
	}

	type payload struct {
		Account string `json:"account"`
	}

	tests := []struct {
		name         string
		body         string
		rules        Rules
		expectedData payload
		errors       url.Values
	}{
		{"account valid", `{"account":"2017-02-10"}`, singleRules, payload{Account: "2017-02-10"}, url.Values{}},
		{"account not valid", `{"account":"6"}`, singleRules, payload{Account: "6"}, url.Values{
			"account": []string{"The account field is not a valid date time string according to layout 2006-01-02"},
		}},
		{"account invalid JSON", `{"account":"6"`, singleRules, payload{}, url.Values{
			"_error": []string{"unexpected EOF"},
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest("GET", "/", strings.NewReader(test.body))
			require.NoError(t, err)

			data := payload{}
			errors := ValidateJSONBody(request, &data, test.rules)
			assert.Equal(t, test.errors, errors)
			assert.Equal(t, test.expectedData, data)
		})
	}
}
