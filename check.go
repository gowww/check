// Package check provides form validation utilities.
package check

import (
	"fmt"
	"net/http"
	"strings"
)

func errRuleFormat(rule string, args []string) error {
	if len(args) > 0 {
		rule += ":" + strings.Join(args, ":")
	}
	return fmt.Errorf("check: cannot parse rule %q", rule)
}

// Rules in a map of checking rules for keys.
// type Rules map[string][]Rule

// A Checker contains keys with their checking rules.
type Checker map[string][]Rule

// Check makes the check for data (key to multiple values) and returns errors.
func (c Checker) Check(data map[string][]string) Errors {
	errs := make(Errors)
	for k, rules := range c {
		if vv, ok := data[k]; ok {
			for _, v := range vv {
				for _, rule := range rules {
					rule(errs, k, v)
				}
			}
			continue
		}
		// No data for key: see if it's required by checker.
		for _, rule := range rules {
			rule(errs, k, "")
		}
		if kerrs := errs[k]; len(kerrs) == 1 && kerrs[0] == ErrRequired {
			continue
		}
		delete(errs, k) // Checks has been made for key and no "required" error at the end: remove other potential errors.
	}
	return errs
}

// CheckRequest makes the check for an HTTP request and returns errors.
//
// Request data can have multiple values with the same key (or field).
// In tis case, all values are checked and if one fails, the error is set for the whole key.
func (c Checker) CheckRequest(r *http.Request) Errors {
	if r.Form == nil {
		r.ParseMultipartForm(32 << 20) // 32 MB
	}
	return c.Check(r.Form)
}
