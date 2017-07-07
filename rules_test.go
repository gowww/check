package check

import (
	"reflect"
	"testing"
)

type ruleCases []struct {
	v    string
	data map[string][]string
	rule Rule
	want []string
}

func testRule(t *testing.T, name string, cases ruleCases) {
	for _, c := range cases {
		errs := c.rule(c.v, c.data)
		if !reflect.DeepEqual(c.want, errs) {
			t.Errorf("%s(%q): want %v, got %v", name, c.v, c.want, errs)
		}
	}
}

func TestAlpha(t *testing.T) {
	testRule(t, "Alpha", ruleCases{
		{"a", nil, Alpha, nil},
		{"aa", nil, Alpha, nil},
		{"1", nil, Alpha, []string{ErrNotAlpha}},
		{"a.a", nil, Alpha, []string{ErrNotAlpha}},
		{"a@a", nil, Alpha, []string{ErrNotAlpha}},
	})
}

func TestEmail(t *testing.T) {
	testRule(t, "Email", ruleCases{
		{"a@a.aa", nil, Email, nil},
		{"a+a@a.aa", nil, Email, nil},
		{"a", nil, Email, []string{ErrNotEmail}},
		{"a@a.a", nil, Email, []string{ErrNotEmail}},
		{"a+a@a.a", nil, Email, []string{ErrNotEmail}},
		{"@a.a", nil, Email, []string{ErrNotEmail}},
		{"a@a.", nil, Email, []string{ErrNotEmail}},
		{"a@a", nil, Email, []string{ErrNotEmail}},
		{"a.a", nil, Email, []string{ErrNotEmail}},
	})
}

func TestInteger(t *testing.T) {
	testRule(t, "Email", ruleCases{
		{"1", nil, Integer, nil},
		{"123", nil, Integer, nil},
		{".", nil, Integer, []string{ErrNotInteger}},
		{". ", nil, Integer, []string{ErrNotInteger}},
		{"1 123", nil, Integer, []string{ErrNotInteger}},
		{"123.45", nil, Integer, []string{ErrNotInteger}},
		{"123,45", nil, Integer, []string{ErrNotInteger}},
		{"a123", nil, Integer, []string{ErrNotInteger}},
		{"123a", nil, Integer, []string{ErrNotInteger}},
		{"a", nil, Integer, []string{ErrNotInteger}},
	})
}

func TestLatitude(t *testing.T) {
	testRule(t, "Latitude", ruleCases{
		{"12.3", nil, Latitude, nil},
		{"+12.3", nil, Latitude, nil},
		{"-12.3", nil, Latitude, nil},
		{"200", nil, Latitude, []string{ErrNotLatitude}},
		{"-200", nil, Latitude, []string{ErrNotLatitude}},
		{"a", nil, Latitude, []string{ErrNotNumber}},
		{"a1", nil, Latitude, []string{ErrNotNumber}},
	})
}

func TestLongitude(t *testing.T) {
	testRule(t, "Longitude", ruleCases{
		{"78", nil, Longitude, nil},
		{"+78.9", nil, Longitude, nil},
		{"-78.9", nil, Longitude, nil},
		{"200", nil, Longitude, []string{ErrNotLongitude}},
		{"-200", nil, Longitude, []string{ErrNotLongitude}},
		{"a", nil, Longitude, []string{ErrNotNumber}},
		{"a1", nil, Longitude, []string{ErrNotNumber}},
	})
}

func TestMax(t *testing.T) {
	testRule(t, "Max", ruleCases{
		{"0", nil, Max(3), nil},
		{"1", nil, Max(3), nil},
		{"3", nil, Max(3), nil},
		{"-123.45", nil, Max(3), nil},
		{"5", nil, Max(3), []string{ErrMax + ":3"}},
		{"a", nil, Max(3), []string{ErrNotNumber}},
		{"a1", nil, Max(3), []string{ErrNotNumber}},
		{".", nil, Max(-1), []string{ErrNotNumber}},
	})
}

func TestMaxLen(t *testing.T) {
	testRule(t, "MaxLen", ruleCases{
		{"a", nil, MaxLen(3), nil},
		{"   ", nil, MaxLen(3), nil},
		{"aaaa", nil, MaxLen(3), []string{ErrMaxLen + ":3"}},
	})
}

func TestMin(t *testing.T) {
	testRule(t, "Min", ruleCases{
		{"3", nil, Min(3), nil},
		{"+123.45", nil, Min(3), nil},
		{"1", nil, Min(3), []string{ErrMin + ":3"}},
		{"a", nil, Min(3), []string{ErrNotNumber}},
		{"a1", nil, Min(3), []string{ErrNotNumber}},
		{".", nil, Min(3), []string{ErrNotNumber}},
	})
}

func TestMinLen(t *testing.T) {
	testRule(t, "MinLen", ruleCases{
		{"aaa", nil, MinLen(3), nil},
		{"    ", nil, MinLen(3), nil},
		{"a", nil, MinLen(3), []string{ErrMinLen + ":3"}},
	})
}

func TestNumber(t *testing.T) {
	testRule(t, "Number", ruleCases{
		{"1", nil, Number, nil},
		{"123", nil, Number, nil},
		{"-123.45", nil, Number, nil},
		{"a1", nil, Number, []string{ErrNotNumber}},
		{"a", nil, Number, []string{ErrNotNumber}},
		{".", nil, Number, []string{ErrNotNumber}},
	})
}

func TestPhone(t *testing.T) {
	testRule(t, "Phone", ruleCases{
		{"0012345678901", nil, Phone, nil},
		{"+12 (0) 345.67.89.01", nil, Phone, nil},
		{"00123", nil, Phone, []string{ErrNotPhone}},
		{"aaa", nil, Phone, []string{ErrNotPhone}},
		{"aaaaaaaaaa", nil, Phone, []string{ErrNotPhone}},
		{"aaa12345678901", nil, Phone, []string{ErrNotPhone}},
	})
}

func TestRange(t *testing.T) {
	testRule(t, "Range", ruleCases{
		{"5", nil, Range(3, 6), nil},
		{"1", nil, Range(1, 1), nil},
		{"2", nil, Range(3, 6), []string{ErrMin + ":3"}},
		{"0", nil, Range(1, 1), []string{ErrMin + ":1"}},
		{"2", nil, Range(1, 1), []string{ErrMax + ":1"}},
		{"a", nil, Range(3, 6), []string{ErrNotNumber}},
		{"a1", nil, Range(3, 6), []string{ErrNotNumber}},
		{".", nil, Range(0, 0), []string{ErrNotNumber}},
	})
}

func TestRangeLen(t *testing.T) {
	testRule(t, "RangeLen", ruleCases{
		{"a", nil, RangeLen(1, 1), nil},
		{"aaaaa", nil, RangeLen(3, 6), nil},
		{"     ", nil, RangeLen(3, 6), nil},
		{"a", nil, RangeLen(3, 6), []string{ErrMinLen + ":3"}},
		{"aaa", nil, RangeLen(1, 2), []string{ErrMaxLen + ":2"}},
	})
}

func TestSame(t *testing.T) {
	testRule(t, "Range", ruleCases{
		{"v", map[string][]string{"k": {"v"}}, Same("k"), nil},
		{"v", map[string][]string{"k": {"v"}}, Same("x"), []string{ErrNotSame + ":x"}},
		{"x", map[string][]string{"k": {"v"}}, Same("k"), []string{ErrNotSame + ":k"}},
		{"v", nil, Same("k"), []string{ErrNotSame + ":k"}},
	})
}
