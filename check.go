// Package check provides form validation utilities.
package check

import (
	"net/http"
	"strconv"
	"strings"
)

// A Checker contains keys with their checking rules.
type Checker map[string]string

// Check makes the check for data (key to multiple values) and returns errors.
func (c Checker) Check(data map[string][]string) Errors {
	errs := make(Errors)
	for key, srules := range c {
		rules := strings.Split(srules, ",")
		values, ok := data[key]
		if !ok {
			// Check rules contains require.
			for _, rule := range rules {
				if rule == "required" {
					errs.Add(key, ErrRequired)
					break
				}
			}
			continue
		}
		for _, v := range values {
			for _, rule := range rules {
				if rule == "required" && v == "" {
					errs.Add(key, ErrRequired)
				} else if rule == "email" {
					errs.Add(key, IsEmail(v)...)
				} else if rule == "integer" {
					errs.Add(key, IsInteger(v)...)
				} else if rule == "number" {
					errs.Add(key, IsNumber(v)...)
				} else if rule == "phone" {
					errs.Add(key, IsPhone(v)...)
				} else { // Rules with arguments.
					ruleParts := strings.Split(rule, ":")
					rule = ruleParts[0]
					if rule == "max" {
						max, err := strconv.ParseFloat(ruleParts[1], 64)
						if err != nil {
							panic(err)
						}
						errs.Add(key, IsMax(v, max)...)
					} else if rule == "min" {
						min, err := strconv.ParseFloat(ruleParts[1], 64)
						if err != nil {
							panic(err)
						}
						errs.Add(key, IsMin(v, min)...)
					} else if rule == "range" {
						min, err := strconv.ParseFloat(ruleParts[1], 64)
						if err != nil {
							panic(err)
						}
						max, err := strconv.ParseFloat(ruleParts[2], 64)
						if err != nil {
							panic(err)
						}
						errs.Add(key, IsInRange(v, min, max)...)
					}
				}
			}
		}
	}
	return errs
}

// CheckRequest makes the check for an HTTP request and returns errors.
func (c Checker) CheckRequest(r *http.Request) Errors {
	if r.Form == nil {
		r.ParseMultipartForm(32 << 20) // 32 MB
	}
	return c.Check(r.Form)
}
