package check_test

import "github.com/gowww/check"

func Example() {
	checker := check.Checker{
		"email": "required,email",
		"phone": "phone",
		"stars": "required,min:3",
	}

	errs := checker.Check(map[string][]string{
		"name":  {"foobar"},
		"phone": {"0012345678901"},
		"stars": {"2"},
	})

	if errs.NotEmpty() {
		// Handle errors.
	}
}
