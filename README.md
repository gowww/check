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

1. Make a [Checker](https://godoc.org/github.com/gowww/check#Checker) with [rules](#rules) for keys:

	```Go
	checker := check.Checker{
		"email": {check.Required, check.Email},
		"phone": {check.Phone},
		"stars": {check.Required, check.Range(3, 5)},
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
		errsjs, _ := json.Marshal(errs.JSON())
		w.Write(errsjs)
	}
	```

### Rules

Function                                                        | Usage                               | Possible errors
----------------------------------------------------------------|-------------------------------------|------------------------------------
[Alpha](https://godoc.org/github.com/gowww/check#Alpha)         | `Alpha`                             | `notAlpha`
[Email](https://godoc.org/github.com/gowww/check#Email)         | `Email`                             | `notEmail`
[Integer](https://godoc.org/github.com/gowww/check#Integer)     | `Integer`                           | `notInteger`
[Latitude](https://godoc.org/github.com/gowww/check#Latitude)   | `Latitude`                          | `notLatitude`, `notNumber`
[Longitude](https://godoc.org/github.com/gowww/check#Longitude) | `Longitude`                         | `notLongitude`, `notNumber`
[Max](https://godoc.org/github.com/gowww/check#Max)             | `Max(1)`                            | `max:1`, `notNumber`
[MaxLen](https://godoc.org/github.com/gowww/check#MaxLen)       | `MaxLen(1)`                         | `maxLen:1`, `notNumber`
[Min](https://godoc.org/github.com/gowww/check#Min)             | `Min(1)`                            | `min:1`, `notNumber`
[MinLen](https://godoc.org/github.com/gowww/check#MinLen)       | `MinLen(1)`                         | `minLen:1`, `notNumber`
[Number](https://godoc.org/github.com/gowww/check#Number)       | `Number`                            | `notNumber`
[Phone](https://godoc.org/github.com/gowww/check#Phone)         | `Phone`                             | `notPhone`
[Range](https://godoc.org/github.com/gowww/check#Range)         | `Range(1, 5)`                       | `max:5`, `min:1`, `notNumber`
[RangeLen](https://godoc.org/github.com/gowww/check#RangeLen)   | `RangeLen(1, 5)`                    | `maxLen:5`, `minLen:1`, `notNumber`
[Required](https://godoc.org/github.com/gowww/check#Required)   | `Required`                          | `required`
[Same](https://godoc.org/github.com/gowww/check#Same)           | `Same("key1", "key2")`              | `notSame:key1,key2`
[Unique](https://godoc.org/github.com/gowww/check#Unique)       | `Unique(db, "users", "email", "?")` | `notUnique`
