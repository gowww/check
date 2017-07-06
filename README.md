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

1. Make a [Checker](https://godoc.org/github.com/gowww/check#Checker) with [rules](#rules) (separated by comma) for keys:

	```Go
	checker := check.Checker{
		"email": "required,email",
		"phone": "phone",
		"stars": "required,min:3",
	}
	```

2. Check data:

	- From a map:
	
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

3. Handle errors:

	```Go
	if errs.NotEmpty() {
		fmt.Println(errs)
	}
	```

	If errors must be JSON formatted (for an HTTP API response, by example), use [Errors.JSON](https://godoc.org/github.com/gowww/check#Errors.JSON):

	```Go
	if errs.NotEmpty() {
		errsjs, _ := json.Marshal(errs)
		w.Write(errsjs)
	}
	```

### Rules

Function    | Description                         | Usage        | Possible errors
------------|-------------------------------------|--------------|------------------------------
`alpha`     | Contains alpha characters only.     | `alpha`      | `notAlpha`
`email`     | Represents an email.                | `email`      | `notEmail`
`integer`   | Represents an integer.              | `integer`    | `notInteger`
`latitude`  | Represents a latitude.              | `latitude`   | `notLatitude`, `notNumber`
`longitude` | Represents a longitude.             | `longitude`  | `notLongitude`, `notNumber`
`max`       | Is below or equals max.             | `max:1`      | `max:1`, `notNumber`
`maxlen`    | Length is below or equals max.      | `maxlen:1`   | `maxLen:1`, `notNumber`
`min`       | Is over or equals min.              | `min:1`      | `min:1`, `notNumber`
`minlen`    | Length is over or equals min.       | `minlen:1`   | `minLen:1`, `notNumber`
`number`    | Represents a number.                | `number`     | `notNumber`
`phone`     | Represents a phone number.          | `phone`      | `notPhone`
`range`     | Represents a number inside a range. | `range:1:10` | `max:1`, `min:1`, `notNumber`
`required`  | Value is not empry.                 | `required`   | `required`
