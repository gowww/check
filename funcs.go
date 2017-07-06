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

// IsAlpha checks if v contains alpha characters only.
func IsAlpha(v string) []error {
	for i := 0; i < len(v); i++ {
		if v[i] < 65 || v[i] > 90 && v[i] < 97 || v[i] > 122 {
			return []error{ErrNotAlpha}
		}
	}
	return nil
}

// IsEmail checks if v represents an email.
func IsEmail(v string) []error {
	if reEmail.MatchString(v) {
		return nil
	}
	return []error{ErrNotEmail}
}

// IsInRange check if v represents a number inside a range.
func IsInRange(v string, min, max float64) []error {
	return append(IsMax(v, max), IsMin(v, min)...)
}

// IsInteger check if v represents an integer.
func IsInteger(v string) []error {
	if _, err := strconv.Atoi(v); err != nil {
		return []error{ErrNotInteger}
	}
	return nil
}

// IsLatitude check if v represents a latitude.
func IsLatitude(v string) []error {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return []error{ErrNotNumber}
	}
	if f < -90 || f > 90 {
		return []error{ErrNotLatitude}
	}
	return nil
}

// IsLongitude check if v represents a longitude.
func IsLongitude(v string) []error {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return []error{ErrNotNumber}
	}
	if f < -180 || f > 180 {
		return []error{ErrNotLongitude}
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

// IsMaxLen check if v length is below or equals max.
func IsMaxLen(v string, max int) []error {
	if len(v) > max {
		return []error{fmt.Errorf("%v:%v", ErrMaxLen, max)}
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

// IsMinLen check if v length is over or equals min.
func IsMinLen(v string, min int) []error {
	if len(v) < min {
		return []error{fmt.Errorf("%v:%v", ErrMinLen, min)}
	}
	return nil
}

// IsNumber check if v represents a number.
func IsNumber(v string) []error {
	_, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return []error{ErrNotNumber}
	}
	return nil
}

// IsPhone checks if v represents a phone number.
func IsPhone(v string) []error {
	if rePhone.MatchString(v) {
		return nil
	}
	return []error{ErrNotPhone}
}
