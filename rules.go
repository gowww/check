package check

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	reEmail = regexp.MustCompile(`^[^@]+@[^@]+\.[^@]{2,63}$`)
	rePhone = regexp.MustCompile(`^\+?(\d|\(|\)|\.|\s){9,20}$`)
)

// A Rule is a checking function to use inside a Checker.
// It receives a value to check and the whole data map for relative checks.
// It returns error identifiers or nil if check pass.
type Rule func(string, map[string][]string) []string

// Alpha rule checks if v contains alpha characters only.
func Alpha(v string, _ map[string][]string) []string {
	for i := 0; i < len(v); i++ {
		if v[i] < 65 || v[i] > 90 && v[i] < 97 || v[i] > 122 {
			return []string{ErrNotAlpha}
		}
	}
	return nil
}

// Email rule checks if v represents an email.
func Email(v string, _ map[string][]string) []string {
	if !reEmail.MatchString(v) {
		return []string{ErrNotEmail}
	}
	return nil
}

// Integer rule checks if v represents an integer.
func Integer(v string, _ map[string][]string) []string {
	if v == "." {
		return []string{ErrNotInteger}
	}
	if _, err := strconv.Atoi(v); err != nil {
		return []string{ErrNotInteger}
	}
	return nil
}

// Latitude rule checks if v represents a latitude.
func Latitude(v string, _ map[string][]string) []string {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return []string{ErrNotNumber}
	}
	if f < -90 || f > 90 {
		return []string{ErrNotLatitude}
	}
	return nil
}

// Longitude rule checks if v represents a longitude.
func Longitude(v string, _ map[string][]string) []string {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return []string{ErrNotNumber}
	}
	if f < -180 || f > 180 {
		return []string{ErrNotLongitude}
	}
	return nil
}

// Max rule checks if v is below or equals max.
func Max(max float64) Rule {
	return func(v string, _ map[string][]string) []string {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return []string{ErrNotNumber}
		}
		if f > max {
			return []string{fmt.Sprintf("%s:%g", ErrMax, max)}
		}
		return nil
	}
}

// MaxLen rule checks if v length is below or equals max.
func MaxLen(max int) Rule {
	return func(v string, _ map[string][]string) []string {
		if len(v) > max {
			return []string{fmt.Sprintf("%s:%d", ErrMaxLen, max)}
		}
		return nil
	}
}

// Min rule checks if v is over or equals min.
func Min(min float64) Rule {
	return func(v string, _ map[string][]string) []string {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return []string{ErrNotNumber}
		}
		if f < min {
			return []string{fmt.Sprintf("%s:%g", ErrMin, min)}
		}
		return nil
	}
}

// MinLen rule checks if v length is over or equals min.
func MinLen(min int) Rule {
	return func(v string, _ map[string][]string) []string {
		if len(v) < min {
			return []string{fmt.Sprintf("%s:%d", ErrMinLen, min)}
		}
		return nil
	}
}

// Number rule checks if v represents a number.
func Number(v string, _ map[string][]string) []string {
	_, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return []string{ErrNotNumber}
	}
	return nil
}

// Phone rule checks if v represents a phone number.
func Phone(v string, _ map[string][]string) []string {
	if !rePhone.MatchString(v) {
		return []string{ErrNotPhone}
	}
	return nil
}

// Range rule checks if v represents a number inside a range.
func Range(min, max float64) Rule {
	return func(v string, _ map[string][]string) []string {
		if errs := Max(max)(v, nil); errs != nil {
			return errs
		}
		return Min(min)(v, nil)
	}
}

// RangeLen rule checks if v length is between or equal min and max.
func RangeLen(min, max int) Rule {
	return func(v string, _ map[string][]string) []string {
		if errs := MinLen(min)(v, nil); errs != nil {
			return errs
		}
		return MaxLen(max)(v, nil)
	}
}

// Required rule checks that v is not empty.
// v is not trimmed so a single space can pass the check.
func Required(v string, _ map[string][]string) []string {
	if v == "" {
		return []string{ErrRequired}
	}
	return nil
}

// Same rule checks that v is the same as another key value.
func Same(keys ...string) Rule {
	return func(v string, data map[string][]string) []string {
	KeysLoop:
		for _, key := range keys {
			vv := data[key]
			if len(vv) == 0 {
				return []string{ErrNotSame + ":" + strings.Join(keys, ",")}
			}
			for _, v2 := range vv {
				if v == v2 {
					continue KeysLoop
				}
			}
			return []string{ErrNotSame + ":" + strings.Join(keys, ",")}
		}
		return nil
	}
}

// Unique rule checks if v is unique in database.
// The placeholder ("?", "$1" or other) must be provided as it depends on the SQL driver.
func Unique(db *sql.DB, table, column, placeholder string) Rule {
	if db == nil {
		panic(`check: no database provided for "unique" rule`)
	}
	return func(v string, _ map[string][]string) []string {
		var n int
		if err := db.QueryRow("SELECT COUNT() FROM "+table+" WHERE "+column+" = "+placeholder, v).Scan(&n); err != nil {
			panic(err)
		}
		if n > 0 {
			return []string{ErrNotUnique}
		}
		return nil
	}
}
