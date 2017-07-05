package check

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	reEmail = regexp.MustCompile(`.+@.+\..+`)
	rePhone = regexp.MustCompile(`^\+?(\d|\(|\)|\.|\s){9,20}$`)
)

// IsEmail checks if v is an email.
func IsEmail(v string) []error {
	if reEmail.MatchString(v) {
		return nil
	}
	return []error{ErrNotEmail}
}

// IsInRange check if v is a number in a certain range.
func IsInRange(v string, min, max float64) []error {
	return append(IsMax(v, max), IsMin(v, min)...)
}

// IsInteger check if v is an integer.
func IsInteger(v string) []error {
	if _, err := strconv.Atoi(v); err != nil {
		return []error{ErrNotInteger}
	}
	return nil
}

// IsMax check if v is below or equals max.
func IsMax(v string, max float64) []error {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return []error{ErrNotNumber}
	}
	if f > max {
		return []error{fmt.Errorf("%v:%v", ErrMax, max)}
	}
	return nil
}

// IsMin check if v is over or equals min.
func IsMin(v string, min float64) []error {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return []error{ErrNotNumber}
	}
	if f < min {
		return []error{fmt.Errorf("%v:%v", ErrMin, min)}
	}
	return nil
}

// IsNumber check if v is a number.
func IsNumber(v string) []error {
	_, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return []error{ErrNotNumber}
	}
	return nil
}

// IsPhone checks if v is a phone number.
func IsPhone(v string) []error {
	if rePhone.MatchString(v) {
		return nil
	}
	return []error{ErrNotPhone}
}
