package check

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
)

var (
	reEmail = regexp.MustCompile(`^[^@]+@[^@]+\.[^@]{2,63}$`)
	rePhone = regexp.MustCompile(`^\+?(\d|\(|\)|\.|\s){9,20}$`)
)

// A Rule is a checking function to use inside a Checker.
type Rule func(Errors, string, string)

// Alpha checks if v contains alpha characters only.
func Alpha(errs Errors, k, v string) {
	for i := 0; i < len(v); i++ {
		if v[i] < 65 || v[i] > 90 && v[i] < 97 || v[i] > 122 {
			errs.Add(k, ErrNotAlpha)
			return
		}
	}
}

// Email checks if v represents an email.
func Email(errs Errors, k, v string) {
	if !reEmail.MatchString(v) {
		errs.Add(k, ErrNotEmail)
	}
}

// Integer checks if v represents an integer.
func Integer(errs Errors, k, v string) {
	if v == "." {
		errs.Add(k, ErrNotInteger)
		return
	}
	if _, err := strconv.Atoi(v); err != nil {
		errs.Add(k, ErrNotInteger)
	}
}

// Latitude checks if v represents a latitude.
func Latitude(errs Errors, k, v string) {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		errs.Add(k, ErrNotNumber)
		return
	}
	if f < -90 || f > 90 {
		errs.Add(k, ErrNotLatitude)
	}
}

// Longitude checks if v represents a longitude.
func Longitude(errs Errors, k, v string) {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		errs.Add(k, ErrNotNumber)
		return
	}
	if f < -180 || f > 180 {
		errs.Add(k, ErrNotLongitude)
	}
}

// Max checks if v is below or equals max.
func Max(max float64) Rule {
	return func(errs Errors, k string, v string) {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			errs.Add(k, ErrNotNumber)
			return
		}
		if f > max {
			errs.Add(k, fmt.Sprintf("%s:%g", ErrMax, max))
		}
	}
}

// MaxLen checks if v length is below or equals max.
func MaxLen(max int) Rule {
	return func(errs Errors, k string, v string) {
		if len(v) > max {
			errs.Add(k, fmt.Sprintf("%s:%d", ErrMaxLen, max))
		}
	}
}

// Min checks if v is over or equals min.
func Min(min float64) Rule {
	return func(errs Errors, k string, v string) {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			errs.Add(k, ErrNotNumber)
			return
		}
		if f < min {
			errs.Add(k, fmt.Sprintf("%s:%g", ErrMin, min))
		}
	}
}

// MinLen checks if v length is over or equals min.
func MinLen(min int) Rule {
	return func(errs Errors, k string, v string) {
		if len(v) < min {
			errs.Add(k, fmt.Sprintf("%s:%d", ErrMinLen, min))
		}
	}
}

// Number checks if v represents a number.
func Number(errs Errors, k, v string) {
	_, err := strconv.ParseFloat(v, 64)
	if err != nil {
		errs.Add(k, ErrNotNumber)
	}
}

// Phone checks if v represents a phone number.
func Phone(errs Errors, k, v string) {
	if !rePhone.MatchString(v) {
		errs.Add(k, ErrNotPhone)
	}
}

// Range checks if v represents a number inside a range.
func Range(min, max float64) Rule {
	return func(errs Errors, k string, v string) {
		Min(min)(errs, k, v)
		Max(max)(errs, k, v)
	}
}

// RangeLen checks if v length is between or equal min and max.
func RangeLen(min, max int) Rule {
	return func(errs Errors, k string, v string) {
		MinLen(min)(errs, k, v)
		MaxLen(max)(errs, k, v)
	}
}

// Required checks that v is not empty.
// V is not trimmed so a space is a value.
func Required(errs Errors, k, v string) {
	if v == "" {
		errs.Add(k, ErrRequired)
	}
}

// Unique checks if v is unique in database.
// The placeholder must be provided as it depends on the SQL driver.
func Unique(db *sql.DB, table, column, placeholder string) Rule {
	if db == nil {
		panic(`check: no database provided for "unique" rule`)
	}
	return func(errs Errors, k string, v string) {
		var n int
		if err := db.QueryRow("SELECT COUNT() FROM "+table+" WHERE "+column+" = "+placeholder, v).Scan(&n); err != nil {
			panic(err)
		}
		if n > 0 {
			errs.Add(k, ErrNotUnique)
		}
	}
}
