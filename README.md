# [![gowww](https://avatars.githubusercontent.com/u/18078923?s=20)](https://github.com/gowww) check [![GoDoc](https://godoc.org/github.com/gowww/check?status.svg)](https://godoc.org/github.com/gowww/check) [![Build](https://travis-ci.org/gowww/check.svg?branch=master)](https://travis-ci.org/gowww/check) [![Coverage](https://coveralls.io/repos/github/gowww/check/badge.svg?branch=master)](https://coveralls.io/github/gowww/check?branch=master) [![Go Report](https://goreportcard.com/badge/github.com/gowww/check)](https://goreportcard.com/report/github.com/gowww/check)

Package [check](https://godoc.org/github.com/gowww/check) provides form validation utilities.

## Installing

1. Get package:

	```Shell
	go get -u github.com/gowww/check
	```

2. Import it in your code:

	```Go
	import "github.com/gowww/check"
	```

## Usage

1. Make a [Checker](https://godoc.org/github.com/gowww/check#Checker) with rules for keys:

	```Go
	checker := check.Checker{
		"email": "required,email",
		"phone": "phone",
		"stars": "required,min:3",
	}
	```

2. Check you data:

	- Using a map:
	
		```Go
		errs := checker.Check(map[string][]string{
			"name":  {"foobar"},
			"phone": {"0012345678901"},
			"stars": {"2"},
		})
		```

	- From an [http.Request](https://golang.org/pkg/net/http/#Request):
	
		```Go
		errs := checker.CheckRequest(r)
		```

3. Use errors like you want:

	```Go
	if errs.NotEmpty() {
		// Handle errors.
	}
	```
