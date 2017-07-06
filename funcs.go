package check

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	reEmail = regexp.MustCompile(`^[^@]+@[^@]+\.[^@]{2,63}$`)
	rePhone = regexp.MustCompile(`^\+?(\d|\(|\)|\.|\s){9,20}$`)
)

// IsAlpha checks if v contains alpha characters only.
// If check pass, nil is returned.
func IsAlpha(v string) []string {
	for i := 0; i < len(v); i++ {
		if v[i] < 65 || v[i] > 90 && v[i] < 97 || v[i] > 122 {
			return []string{ErrNotAlpha}
		}
	}
	return nil
}

// IsEmail checks if v represents an email.
// If check pass, nil is returned.
func IsEmail(v string) []string {
	if reEmail.MatchString(v) {
		return nil
	}
	return []string{ErrNotEmail}
}

// IsInRange checks if v represents a number inside a range.
// If check pass, nil is returned.
func IsInRange(v string, min, max float64) []string {
	return append(IsMax(v, max), IsMin(v, min)...)
}

// IsInteger checks if v represents an integer.
// If check pass, nil is returned.
func IsInteger(v string) []string {
	if _, err := strconv.Atoi(v); err != nil {
		return []string{ErrNotInteger}
	}
	return nil
}

// IsLatitude checks if v represents a latitude.
// If check pass, nil is returned.
func IsLatitude(v string) []string {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return []string{ErrNotNumber}
	}
	if f < -90 || f > 90 {
		return []string{ErrNotLatitude}
	}
	return nil
}

// IsLongitude checks if v represents a longitude.
// If check pass, nil is returned.
func IsLongitude(v string) []string {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return []string{ErrNotNumber}
	}
	if f < -180 || f > 180 {
		return []string{ErrNotLongitude}
	}
	return nil
}

// IsMax checks if v is below or equals max.
// If check pass, nil is returned.
func IsMax(v string, max float64) []string {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return []string{ErrNotNumber}
	}
	if f > max {
		return []string{fmt.Sprintf("%s:%f", ErrMax, max)}
	}
	return nil
}

// IsMaxLen checks if v length is below or equals max.
// If check pass, nil is returned.
func IsMaxLen(v string, max int) []string {
	if len(v) > max {
		return []string{fmt.Sprintf("%s:%d", ErrMaxLen, max)}
	}
	return nil
}

// IsMin checks if v is over or equals min.
// If check pass, nil is returned.
func IsMin(v string, min float64) []string {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return []string{ErrNotNumber}
	}
	if f < min {
		return []string{fmt.Sprintf("%s:%f", ErrMin, min)}
	}
	return nil
}

// IsMinLen checks if v length is over or equals min.
// If check pass, nil is returned.
func IsMinLen(v string, min int) []string {
	if len(v) < min {
		return []string{fmt.Sprintf("%s:%d", ErrMinLen, min)}
	}
	return nil
}

// IsNumber checks if v represents a number.
// If check pass, nil is returned.
func IsNumber(v string) []string {
	_, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return []string{ErrNotNumber}
	}
	return nil
}

// IsPhone checks if v represents a phone number.
// If check pass, nil is returned.
func IsPhone(v string) []string {
	if rePhone.MatchString(v) {
		return nil
	}
	return []string{ErrNotPhone}
}
