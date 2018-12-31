# [![gowww](https://avatars.githubusercontent.com/u/18078923?s=20)](https://github.com/gowww) check [![GoDoc](https://godoc.org/github.com/gowww/check?status.svg)](https://godoc.org/github.com/gowww/check) [![Build](https://travis-ci.org/gowww/check.svg?branch=master)](https://travis-ci.org/gowww/check) [![Coverage](https://coveralls.io/repos/github/gowww/check/badge.svg?branch=master)](https://coveralls.io/github/gowww/check?branch=master) [![Go Report](https://goreportcard.com/badge/github.com/gowww/check)](https://goreportcard.com/report/github.com/gowww/check) ![Status Testing](https://img.shields.io/badge/status-testing-orange.svg)

Package [check](https://godoc.org/github.com/gowww/check) provides request form checking.

- [Installing](#installing)
- [Usage](#usage)
	- [JSON](#json)
	- [Internationalization](#internationalization)
	- [Rules](#rules)

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
	userChecker := check.Checker{
		"email":   {check.Required, check.Email, check.Unique(db, "users", "email", "?")},
		"phone":   {check.Phone},
		"picture": {check.MaxFileSize(5000000), check.Image},
	}
	```

	The rules order is significant so for example, it's smarter to check the format of a value before its uniqueness, avoiding some useless database requests.

2. Check data:

	- From a values map, with [Checker.CheckValues](https://godoc.org/github.com/gowww/check#Checker.CheckValues):

		```Go
		errs := userChecker.CheckValues(map[string][]string{
			"name":  {"foobar"},
			"phone": {"0012345678901"},
		})
		```

	- From an [http.Request](https://golang.org/pkg/net/http/#Request), with [Checker.CheckRequest](https://godoc.org/github.com/gowww/check#Checker.CheckRequest):

		```Go
		errs := userChecker.CheckRequest(r)
		```

3. Handle errors:

	```Go
	if errs.NotEmpty() {
		fmt.Println(errs)
	}
	```

### JSON

Use [Errors.JSON](https://godoc.org/github.com/gowww/check#Errors.JSON) to get errors in a map under `errors` key, ready to be JSON formatted (as an HTTP API response, for example):

```Go
if errs.NotEmpty() {
	errsjs, _ := json.Marshal(errs.JSON())
	w.Write(errsjs)
}
```

### Internationalization

Internationalization is handled by [gowww/i18n](https://godoc.org/github.com/gowww/i18n) and there are [built-in translations](https://godoc.org/github.com/gowww/check#pkg-variables) for all errors.

Use [Errors.T](https://godoc.org/github.com/gowww/check#Errors.T) with an [i18n.Translator](https://godoc.org/github.com/gowww/i18n#Translator) (usually stored in the request context) to get translated errors:

```Go
if errs.NotEmpty() {
	transErrs := errs.T(i18n.RequestTranslator(r))
	fmt.Println(transErrs)
}
```

You can provide custom translations for each error type under keys like "`error` + ErrorID":

```Go
var locales = i18n.Locales{
	language.English: {
		"hello": "Hello!",

		"errorMaxFileSize": "File too big (%v max.)",
		"errorRequired":    "Required field",
	},
}
```

If the [i18n.Translator](https://godoc.org/github.com/gowww/i18n#Translator) is `nil` or a custom translation is not found, the built-in translation of error is used.

### Rules

Function                                                            | Usage                               | Possible errors
--------------------------------------------------------------------|-------------------------------------|------------------------------------
[Alpha](https://godoc.org/github.com/gowww/check#Alpha)             | `Alpha`                             | `notAlpha`
[Email](https://godoc.org/github.com/gowww/check#Email)             | `Email`                             | `notEmail`
[FileType](https://godoc.org/github.com/gowww/check#FileType)       | `FileType("text/plain")`            | `badFileType:text/plain`
[Image](https://godoc.org/github.com/gowww/check#Image)             | `Image`                             | `notImage`
[Integer](https://godoc.org/github.com/gowww/check#Integer)         | `Integer`                           | `notInteger`
[Latitude](https://godoc.org/github.com/gowww/check#Latitude)       | `Latitude`                          | `notLatitude`, `notNumber`
[Longitude](https://godoc.org/github.com/gowww/check#Longitude)     | `Longitude`                         | `notLongitude`, `notNumber`
[Max](https://godoc.org/github.com/gowww/check#Max)                 | `Max(1)`                            | `max:1`, `notNumber`
[MaxFileSize](https://godoc.org/github.com/gowww/check#MaxFileSize) | `MaxFileSize(5000000)`              | `maxFileSize:5000000`
[MaxLen](https://godoc.org/github.com/gowww/check#MaxLen)           | `MaxLen(1)`                         | `maxLen:1`, `notNumber`
[Min](https://godoc.org/github.com/gowww/check#Min)                 | `Min(1)`                            | `min:1`, `notNumber`
[MinFileSize](https://godoc.org/github.com/gowww/check#MinFileSize) | `MinFileSize(10)`                   | `minFileSize:10`
[MinLen](https://godoc.org/github.com/gowww/check#MinLen)           | `MinLen(1)`                         | `minLen:1`, `notNumber`
[Number](https://godoc.org/github.com/gowww/check#Number)           | `Number`                            | `notNumber`
[Phone](https://godoc.org/github.com/gowww/check#Phone)             | `Phone`                             | `notPhone`
[Range](https://godoc.org/github.com/gowww/check#Range)             | `Range(1, 5)`                       | `max:5`, `min:1`, `notNumber`
[RangeLen](https://godoc.org/github.com/gowww/check#RangeLen)       | `RangeLen(1, 5)`                    | `maxLen:5`, `minLen:1`
[Required](https://godoc.org/github.com/gowww/check#Required)       | `Required`                          | `required`
[Same](https://godoc.org/github.com/gowww/check#Same)               | `Same("key1", "key2")`              | `notSame:key1,key2`
[Unique](https://godoc.org/github.com/gowww/check#Unique)           | `Unique(db, "users", "email", "?")` | `notUnique`
[URL](https://godoc.org/github.com/gowww/check#URL)                 | `URL`                               | `notURL`
