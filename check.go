// Package check provides form validation utilities.
package check

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Checking rules.
const (
	RuleAlpha     = "alpha"
	RuleEmail     = "email"
	RuleInteger   = "integer"
	RuleLatitude  = "latitude"
	RuleLongitude = "longitude"
	RuleMax       = "max"
	RuleMin       = "min"
	RuleNumber    = "number"
	RulePhone     = "phone"
	RuleRange     = "range"
	RuleRequired  = "required"
)

func errRuleFormat(rule string, args []string) error {
	return fmt.Errorf("check: cannot parse rule %q", rule+":"+strings.Join(args, ":"))
}

// A Checker contains keys with their checking rules.
type Checker map[string]string

// Check makes the check for data (key to multiple values) and returns errors.
func (c Checker) Check(data map[string][]string) Errors {
	errs := make(Errors)
	for k, srules := range c {
		rules := strings.Split(srules, ",")
		values, ok := data[k]
		if !ok {
			// Check rules contains require.
			for _, rule := range rules {
				if rule == "required" {
					errs.Add(k, ErrRequired)
					break
				}
			}
			continue
		}
		for _, v := range values {
			for _, rule := range rules {
				ruleCheck(errs, rule, k, v)
			}
		}
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

func ruleCheck(errs Errors, rule, k, v string) {
	if rule == RuleRequired {
		if v == "" {
			errs.Add(k, ErrRequired)
		}
	} else if rule == RuleAlpha {
		errs.Add(k, IsAlpha(v)...)
	} else if rule == RuleEmail {
		errs.Add(k, IsEmail(v)...)
	} else if rule == RuleInteger {
		errs.Add(k, IsInteger(v)...)
	} else if rule == RuleLatitude {
		errs.Add(k, IsLatitude(v)...)
	} else if rule == RuleLongitude {
		errs.Add(k, IsLongitude(v)...)
	} else if rule == RuleNumber {
		errs.Add(k, IsNumber(v)...)
	} else if rule == RulePhone {
		errs.Add(k, IsPhone(v)...)
	} else { // Rules with arguments.
		ruleParts := strings.Split(rule, ":")
		if len(ruleParts) < 2 {
			panic(errRuleFormat(rule, nil))
		}
		rule = ruleParts[0]
		args := ruleParts[1:]
		if rule == RuleMax {
			errs.Add(k, IsMax(v, parseRuleFloat64(rule, args))...)
		} else if rule == RuleMin {
			errs.Add(k, IsMin(v, parseRuleFloat64(rule, args))...)
		} else if rule == RuleRange {
			min, max := parseRuleFloat64Float64(rule, args)
			errs.Add(k, IsInRange(v, min, max)...)
		}
	}
}

func parseRuleFloat64(rule string, args []string) float64 {
	if len(args) != 1 {
		panic(errRuleFormat(rule, args))
	}
	f, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		panic(errRuleFormat(rule, args))
	}
	return f
}

func parseRuleFloat64Float64(rule string, args []string) (float64, float64) {
	if len(args) != 2 {
		panic(errRuleFormat(rule, args))
	}
	f1, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		panic(errRuleFormat(rule, args))
	}
	f2, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		panic(errRuleFormat(rule, args))
	}
	return f1, f2
}
