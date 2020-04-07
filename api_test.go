package validator

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	eos "github.com/eoscanada/eos-go"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
	"github.com/thedevsaddam/govalidator"
)

func init() {
	govalidator.AddCustomRule("eos.blockNum", EOSBlockNumRule)
	govalidator.AddCustomRule("eos.name", EOSNameRule)
}

func TestValidateQueryParams(t *testing.T) {
	singleRules := map[string][]string{
		"block_num": []string{"eos.blockNum"},
	}

	tests := []struct {
		name   string
		query  string
		rules  Rules
		errors url.Values
	}{
		{"block_num valid", "block_num=1", singleRules, url.Values{}},
		{"block_num not valid", "block_num=a", singleRules, url.Values{
			"block_num": []string{"The block_num field must be a valid EOS block num"},
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
		"account": []string{"eos.name"},
	}

	type payload struct {
		Account eos.Name `json:"account"`
	}

	tests := []struct {
		name         string
		rules        Rules
		expectedData payload
		errors       url.Values
	}{
		{"account valid", singleRules, payload{Account: "eos"}, url.Values{}},
		{"account not valid", singleRules, payload{Account: "6"}, url.Values{
			"account": []string{"The account field must be a valid EOS name"},
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
		"account": []string{"eos.name"},
	}

	type payload struct {
		Account eos.Name `schema:"account"`
	}

	tests := []struct {
		name         string
		rules        Rules
		expectedData payload
		errors       url.Values
	}{
		{"account valid", singleRules, payload{Account: "eos"}, url.Values{}},
		{"account not valid", singleRules, payload{Account: "6"}, url.Values{
			"account": []string{"The account field must be a valid EOS name"},
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
		"account": []string{"eos.name"},
	}

	type payload struct {
		Account eos.Name `json:"account"`
	}

	tests := []struct {
		name         string
		body         string
		rules        Rules
		expectedData payload
		errors       url.Values
	}{
		{"account valid", `{"account":"eos"}`, singleRules, payload{Account: "eos"}, url.Values{}},
		{"account not valid", `{"account":"6"}`, singleRules, payload{Account: "6"}, url.Values{
			"account": []string{"The account field must be a valid EOS name"},
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
