## EOS Canada Validator Library

This repository contains all common stuff around validation handling across our
various services

### Philosophy

We started using [govalidator](https://github.com/thedevsaddam/govalidator) for our validation
needs.

This library is based on rules which are plain string identifier backed by a corresponding
registered function that validates the input. When you when to validate a request query parameters
or JSON payload, you can define a set of rules which is simply a `map[string][]string{}` where the
key is the field/parameter to validate while the value is the set of rules identifier you would like
to validate the field against.

Most of the rules accepts string for query parameters and typed value for JSON payload. The [govalidator](https://github.com/thedevsaddam/govalidator) has a bunch of predefined rules also already integrated. See
[govalidator pre-defined rules set](https://github.com/thedevsaddam/govalidator#validation-rules).

The package aims at providing a [quick API](./api.go) to validate either the query parameters of a request
`validator.ValidateQueryParams` or the JSON body payload received `validator.ValidateJSONBody`.

The package also includes a bunch of [predefined rules](./rules.go) useful across all our projects, most of
them related to blockchain ecosystem.

In the rules set, a function ending with `Factory` means it's used to create a `validator.Rule` function
which is parameterized by the parameters received by the factory function. For example, the
`EOSNamesRulesFactory` is used to create a `Rule` function exploding a list of `eos.Name` based on
a separator `sep` and with a maximum count of `maxCount`.

```
eosNamesListRule := validator.EOSNamesRulesFactory("|", 10)
```

### Usage

To efficiently used pre-defined rules inside validator as a string, you must register them through
a central location (`func init()` in a `validators.go` in the package is probably the most common
place).

```
func init() {
    govalidator.AddCustomRule("eos.blockNum", validator.EOSBlockNumRule)
    govalidator.AddCustomRule("eos.name", validator.EOSNameRule)
    govalidator.AddCustomRule("eos.accountsList", validator.EOSNamesListRuleFactory("|", 10))
}
```

You can then pass them as string in your rules set when validating query parameters:

```
errors := validator.ValidateQueryParams(r, govalidator.MapData{
    "account":   []string{"required", "eos.name"},
    "block_num": []string{"eos.blockNum"},
})
```

#### Validate Query Parameters

Simply pass your request and receives back `url.Values` object which is a simple
`map[string][]string` where the key is the name of the offending field and the
value is an array of all the error messages for this field.

```
errors := validator.ValidateQueryParams(r, validator.Rules{
    "account":   []string{"required", "eos.name"},
    "block_num": []string{"eos.blockNum"},
})
```

#### Validate JSON Body Payload

Similar to `validator.ValidateQueryParams` but you pass and extra parameters
which will be the object into which the JSON payload will be deserialized in.

If the deserialization is successful, the `data` object will be populated with
the JSON data. This holds even if the validation later failed, you still have
access to the full deserialized object.

Similar to how one use `json.UnmarshalJSON`, you must pass a pointer to
the `validator.ValidateJSONBody`.

You will receives back `url.Values` object which is a simple
`map[string][]string` where the key is the name of the offending field and the
value is an array of all the error messages for this field.

```
type Request struct {
    Account eos.Name `json:"account"`
}

request := Request{}
errors := validator.ValidateJSONBody(r, &request, validator.Rules{
    "account":   []string{"required", "eos.name"},
})
```

If the deserialization of the JSON failed altogether, validation rules are
not checked. However, to denote such deserialization error, you will get
an errors map containing a single entry whose key will be named `_error` and
the values will be a single element array containing the message why
the JSON deserialization failed.

#### Validate Struct

Can be used to validate any kind of structure. The field names are determined
via the `json` tag if present of the field name (case sensitive) if the `json`
tag is not present. It's also possible to use `validator.TagIdentifierOption`
to specify an alternative tag to use instead of `json`.

Simply pass your struct and receives back `url.Values` object which is a simple
`map[string][]string` where the key is the name of the offending field and the
value is an array of all the error messages for this field.

```
type test struct {
    accountName string `json:account`
    blockNum    uint32 `json:block_num`
}

data := &test{
    accountName: "eosio",
    blockNum: 0,
}

errors := validator.ValidateStruct(data, validator.Rules{
    "account":   []string{"required", "eos.name"},
    "block_num": []string{"eos.blockNum"},
})
```

#### Reference

For now, not much reference documentation exists. You are invited to read the
following files to get a clearer understanding of what you can do.

- [api.go](./api.go)
- [rules.go](./rules.go)

You are specially invited to check out the unit tests which expresses
most of the usage that can be made out of this library:

- [api_test.go](./api_test.go)
- [rules_test.go](./rules_test.go)
