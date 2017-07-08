package check

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	reEmail = regexp.MustCompile(`^[^@]+@[^@]+\.[^@]{2,63}$`)
	rePhone = regexp.MustCompile(`^\+?(\d|\(|\)|\.|\s){9,20}$`)
)

// A Rule is a checking function to be used inside a Checker.
// It receives the errors map to add encountered errors, the whole form for relative checks, and the specific key to check.
type Rule func(errs Errors, form *multipart.Form, key string)

// Alpha rule checks that value contains alpha characters only.
func Alpha(errs Errors, form *multipart.Form, key string) {
	if form == nil && form.Value == nil {
		return
	}
	for _, v := range form.Value[key] {
		for i := 0; i < len(v); i++ {
			if v[i] < 65 || v[i] > 90 && v[i] < 97 || v[i] > 122 {
				errs.Add(key, ErrNotAlpha)
				return
			}
		}
	}
}

// Email rule checks that value represents an email.
func Email(errs Errors, form *multipart.Form, key string) {
	if form == nil && form.Value == nil {
		return
	}
	for _, v := range form.Value[key] {
		if !reEmail.MatchString(v) {
			errs.Add(key, ErrNotEmail)
			return
		}
	}
}

// FileType rule checks that file is one of given MIME types.
func FileType(types ...string) Rule {
	return func(errs Errors, form *multipart.Form, key string) {
		if form == nil && form.File == nil {
			return
		}
		for _, file := range form.File[key] {
			ct, err := fileType(file)
			if err != nil {
				continue
			}
			if !sliceContainsString(types, ct) {
				errs.Add(key, ErrBadFileType+":"+strings.Join(types, ","))
				return
			}
		}
	}
}

// Image rule checks that file is GIF, JPEG or PNG.
func Image(errs Errors, form *multipart.Form, key string) {
	if form == nil && form.File == nil {
		return
	}
	for _, file := range form.File[key] {
		ct, err := fileType(file)
		if err != nil {
			continue
		}
		if !sliceContainsString([]string{"image/gif", "image/jpeg", "image/png"}, ct) {
			errs.Add(key, ErrNotImage)
			return
		}
	}
}

// Integer rule checks that value represents an integer.
func Integer(errs Errors, form *multipart.Form, key string) {
	if form == nil && form.Value == nil {
		return
	}
	for _, v := range form.Value[key] {
		if v == "." {
			errs.Add(key, ErrNotInteger)
			return
		}
		if _, err := strconv.Atoi(v); err != nil {
			errs.Add(key, ErrNotInteger)
			return
		}
	}
}

// Latitude rule checks that value represents a latitude.
func Latitude(errs Errors, form *multipart.Form, key string) {
	if form == nil && form.Value == nil {
		return
	}
	for _, v := range form.Value[key] {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			errs.Add(key, ErrNotNumber)
			return
		}
		if f < -90 || f > 90 {
			errs.Add(key, ErrNotLatitude)
			return
		}
	}
}

// Longitude rule checks that value represents a longitude.
func Longitude(errs Errors, form *multipart.Form, key string) {
	if form == nil && form.Value == nil {
		return
	}
	for _, v := range form.Value[key] {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			errs.Add(key, ErrNotNumber)
			return
		}
		if f < -180 || f > 180 {
			errs.Add(key, ErrNotLongitude)
			return
		}
	}
}

// Max rule checks that value is below or equals max.
func Max(max float64) Rule {
	return func(errs Errors, form *multipart.Form, key string) {
		if form == nil && form.Value == nil {
			return
		}
		for _, v := range form.Value[key] {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				errs.Add(key, ErrNotNumber)
				return
			}
			if f > max {
				errs.Add(key, fmt.Sprintf("%s:%g", ErrMax, max))
				return
			}
		}
	}
}

// MaxFileSize rule checks if file has max or less bytes.
func MaxFileSize(max int64) Rule {
	return func(errs Errors, form *multipart.Form, key string) {
		if form == nil && form.File == nil {
			return
		}
		for _, file := range form.File[key] {
			size, err := fileSize(file)
			if err != nil {
				continue
			}
			if size > max {
				errs.Add(key, fmt.Sprintf("%s:%d", ErrMaxFileSize, max))
				return
			}
		}
	}
}

// MaxLen rule checks that value length is below or equals max.
func MaxLen(max int) Rule {
	return func(errs Errors, form *multipart.Form, key string) {
		if form == nil && form.Value == nil {
			return
		}
		for _, v := range form.Value[key] {
			if len(v) > max {
				errs.Add(key, fmt.Sprintf("%s:%d", ErrMaxLen, max))
				return
			}
		}
	}
}

// Min rule checks that value is over or equals min.
func Min(min float64) Rule {
	return func(errs Errors, form *multipart.Form, key string) {
		if form == nil && form.Value == nil {
			return
		}
		for _, v := range form.Value[key] {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				errs.Add(key, ErrNotNumber)
				return
			}
			if f < min {
				errs.Add(key, fmt.Sprintf("%s:%g", ErrMin, min))
				return
			}
		}
	}
}

// MinFileSize rule checks if file has min or more bytes.
func MinFileSize(min int64) Rule {
	return func(errs Errors, form *multipart.Form, key string) {
		if form == nil && form.File == nil {
			return
		}
		for _, file := range form.File[key] {
			size, err := fileSize(file)
			if err != nil {
				continue
			}
			if size < min {
				errs.Add(key, fmt.Sprintf("%s:%d", ErrMinFileSize, min))
				return
			}
		}
	}
}

// MinLen rule checks that value length is over or equals min.
func MinLen(min int) Rule {
	return func(errs Errors, form *multipart.Form, key string) {
		if form == nil && form.Value == nil {
			return
		}
		for _, v := range form.Value[key] {
			if len(v) < min {
				errs.Add(key, fmt.Sprintf("%s:%d", ErrMinLen, min))
				return
			}
		}
	}
}

// Number rule checks that value represents a number.
func Number(errs Errors, form *multipart.Form, key string) {
	if form == nil && form.Value == nil {
		return
	}
	for _, v := range form.Value[key] {
		_, err := strconv.ParseFloat(v, 64)
		if err != nil {
			errs.Add(key, ErrNotNumber)
			return
		}
	}
}

// Phone rule checks that value represents a phone number.
func Phone(errs Errors, form *multipart.Form, key string) {
	if form == nil && form.Value == nil {
		return
	}
	for _, v := range form.Value[key] {
		if !rePhone.MatchString(v) {
			errs.Add(key, ErrNotPhone)
			return
		}
	}
}

// Range rule checks that value represents a number inside a range.
func Range(min, max float64) Rule {
	return func(errs Errors, form *multipart.Form, key string) {
		if form == nil && form.Value == nil {
			return
		}
		for _, v := range form.Value[key] {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				errs.Add(key, ErrNotNumber)
				return
			}
			if f < min {
				errs.Add(key, fmt.Sprintf("%s:%g", ErrMin, min))
				return
			}
			if f > max {
				errs.Add(key, fmt.Sprintf("%s:%g", ErrMax, max))
				return
			}
		}
	}
}

// RangeFileSize rule checks if file length is inside a range.
func RangeFileSize(min, max int64) Rule {
	return func(errs Errors, form *multipart.Form, key string) {
		if form == nil && form.File == nil {
			return
		}
		for _, file := range form.File[key] {
			size, err := fileSize(file)
			if err != nil {
				continue
			}
			if size > max {
				errs.Add(key, fmt.Sprintf("%s:%d", ErrMaxFileSize, max))
				return
			}
			if size < min {
				errs.Add(key, fmt.Sprintf("%s:%d", ErrMinFileSize, min))
				return
			}
		}
	}
}

// RangeLen rule checks that value length is inside a range.
func RangeLen(min, max int) Rule {
	return func(errs Errors, form *multipart.Form, key string) {
		if form == nil && form.Value == nil {
			return
		}
		for _, v := range form.Value[key] {
			if len(v) > max {
				errs.Add(key, fmt.Sprintf("%s:%d", ErrMaxLen, max))
				return
			}
			if len(v) < min {
				errs.Add(key, fmt.Sprintf("%s:%d", ErrMinLen, min))
				return
			}
		}
	}
}

// Required rule checks that value or file exists and is not empty.
// A value is not trimmed so a single space can pass the check.
func Required(errs Errors, form *multipart.Form, key string) {
	if form == nil {
		errs.Add(key, ErrRequired)
		return
	}
	if form.Value != nil {
		for _, v := range form.Value[key] {
			if v != "" {
				return
			}
		}
	}
	if form.File != nil {
		for _, v := range form.File[key] {
			if v == nil {
				return
			}
		}
	}
	errs.Add(key, ErrRequired)
	return
}

// Same rule checks that value deeply equals another key value.
func Same(keys ...string) Rule {
	return func(errs Errors, form *multipart.Form, key string) {
		if form == nil && form.Value == nil {
			return
		}
		for _, k := range keys {
			if !reflect.DeepEqual(form.Value[key], form.Value[k]) {
				errs.Add(key, ErrNotSame+":"+strings.Join(keys, ","))
				return
			}
		}
	}
}

// Unique rule checks that value is unique in database.
// The placeholder ("?", "$1" or other) must be provided as it depends on the SQL driver.
func Unique(db *sql.DB, table, column, placeholder string) Rule {
	if db == nil {
		panic(`check: no database provided for "unique" rule`)
	}
	return func(errs Errors, form *multipart.Form, key string) {
		for _, v := range form.Value[key] {
			var n int
			if err := db.QueryRow("SELECT COUNT() FROM "+table+" WHERE "+column+" = "+placeholder, v).Scan(&n); err != nil {
				panic(err)
			}
			if n > 0 {
				errs.Add(key, ErrNotUnique)
				return
			}
		}
	}
}

// URL rule checks that value represents an URL.
func URL(errs Errors, form *multipart.Form, key string) {
	if form == nil && form.Value == nil {
		return
	}
	for _, v := range form.Value[key] {
		if len(v) < 4 {
			errs.Add(key, ErrNotURL)
			return
		}
		v = strings.Replace(v, "127.0.0.1", "localhost", 1)
		if i := strings.IndexByte(v, '#'); i != -1 {
			v = v[:i]
		}
		if !strings.Contains(v, "://") {
			v = "http://" + v
		}
		u, err := url.ParseRequestURI(v)
		if err != nil || u.Host == "" || u.Host[0] == '-' || strings.Contains(u.Host, ".-") || strings.Contains(u.Host, "-.") {
			errs.Add(key, ErrNotURL)
			return
		}
		parts := strings.Split(u.Host, ".")
		if parts[0] == "" || len(parts[len(parts)-1]) < 2 || len(parts[len(parts)-1]) > 63 {
			errs.Add(key, ErrNotURL)
			return
		}
		var domain string
		if len(parts) > 2 {
			domain = strings.Join(parts[len(parts)-2:], ".")
		} else {
			domain = strings.Join(parts, ".")
		}
		if strings.ContainsAny(domain, "_,!&") {
			errs.Add(key, ErrNotURL)
			return
		}
		if strings.Count(domain, "::") > 1 { // Only 1 substitution ("::") allowed in IPv6 address.
			errs.Add(key, ErrNotURL)
			return
		}
		parts = strings.Split(domain, ":")
		port, err := strconv.Atoi(parts[len(parts)-1])
		if err == nil && (port < 1 || port > 65535) {
			errs.Add(key, ErrNotURL)
			return
		}
	}
}
