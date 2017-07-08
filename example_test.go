package check_test

import (
	"fmt"

	"github.com/gowww/check"
)

func Example() {
	checker := check.Checker{
		"email":   {check.Required, check.Email},
		"phone":   {check.Phone},
		"picture": {check.MaxFileSize(5000)},
		"stars":   {check.Required, check.Range(3, 5)},
	}

	errs := checker.CheckValues(map[string][]string{
		"name":  {"foobar"},
		"phone": {"0012345678901"},
		"stars": {"2"},
	})

	if errs.NotEmpty() {
		fmt.Println(errs)
	}
}
