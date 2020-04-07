package validator

import (
	"net/http"
	"net/url"

	"github.com/thedevsaddam/govalidator"
)

// Rules represents the set of validation to perform for each path
// of the struct.
type Rules map[string][]string

// Option represents an option that can be set on the validator
// to alter it's behavior
type Option interface {
	set(options *govalidator.Options)
}

// MessagesOption represents the set of messages to return explicitely
// for each validation, defaults to `nil` when undefined (predefined messages are used).
type MessagesOption map[string][]string

// TagIdentifierOption represents the tag to use when inspecting
// a struct to determine the field name, defaults to `json` when undefined.
type TagIdentifierOption string

func (o MessagesOption) set(options *govalidator.Options) {
	options.Messages = govalidator.MapData(o)
}

func (o TagIdentifierOption) set(options *govalidator.Options) {
	options.TagIdentifier = string(o)
}

func ValidateQueryParams(r *http.Request, rules Rules, options ...Option) url.Values {
	return newValidator(r, nil, rules, options).Validate()
}

func ValidateJSONBody(r *http.Request, data interface{}, rules Rules, options ...Option) url.Values {
	return newValidator(r, data, rules, options).ValidateJSON()
}

func ValidateStruct(data interface{}, rules Rules, options ...Option) url.Values {
	return newValidator(nil, data, rules, options).ValidateStruct()
}

func newValidator(r *http.Request, data interface{}, rules Rules, options []Option) *govalidator.Validator {
	opts := govalidator.Options{
		Request: r,
		Rules:   govalidator.MapData(rules),
		Data:    data,
	}

	for _, option := range options {
		option.set(&opts)
	}

	return govalidator.New(opts)
}
