package check

import (
	"reflect"
	"testing"
)

type ruleCases []struct {
	v    string
	rule Rule
	want []string
}

func testRule(t *testing.T, name string, cases ruleCases) {
	for _, c := range cases {
		errs := c.rule(c.v)
		if !reflect.DeepEqual(c.want, errs) {
			t.Errorf("%s(%q): want %v, got %v", name, c.v, c.want, errs)
		}
	}
}

func TestAlpha(t *testing.T) {
	testRule(t, "Alpha", ruleCases{
		{"a", Alpha, nil},
		{"aa", Alpha, nil},
		{"1", Alpha, []string{ErrNotAlpha}},
		{"a.a", Alpha, []string{ErrNotAlpha}},
		{"a@a", Alpha, []string{ErrNotAlpha}},
	})
}

func TestEmail(t *testing.T) {
	testRule(t, "Email", ruleCases{
		{"a@a.aa", Email, nil},
		{"a+a@a.aa", Email, nil},
		{"a", Email, []string{ErrNotEmail}},
		{"a@a.a", Email, []string{ErrNotEmail}},
		{"a+a@a.a", Email, []string{ErrNotEmail}},
		{"@a.a", Email, []string{ErrNotEmail}},
		{"a@a.", Email, []string{ErrNotEmail}},
		{"a@a", Email, []string{ErrNotEmail}},
		{"a.a", Email, []string{ErrNotEmail}},
	})
}

func TestInteger(t *testing.T) {
	testRule(t, "Email", ruleCases{
		{"1", Integer, nil},
		{"123", Integer, nil},
		{".", Integer, []string{ErrNotInteger}},
		{". ", Integer, []string{ErrNotInteger}},
		{"1 123", Integer, []string{ErrNotInteger}},
		{"123.45", Integer, []string{ErrNotInteger}},
		{"123,45", Integer, []string{ErrNotInteger}},
		{"a123", Integer, []string{ErrNotInteger}},
		{"123a", Integer, []string{ErrNotInteger}},
		{"a", Integer, []string{ErrNotInteger}},
	})
}

func TestLatitude(t *testing.T) {
	testRule(t, "Latitude", ruleCases{
		{"12.3", Latitude, nil},
		{"+12.3", Latitude, nil},
		{"-12.3", Latitude, nil},
		{"200", Latitude, []string{ErrNotLatitude}},
		{"-200", Latitude, []string{ErrNotLatitude}},
		{"a", Latitude, []string{ErrNotNumber}},
		{"a1", Latitude, []string{ErrNotNumber}},
	})
}

func TestLongitude(t *testing.T) {
	testRule(t, "Longitude", ruleCases{
		{"78", Longitude, nil},
		{"+78.9", Longitude, nil},
		{"-78.9", Longitude, nil},
		{"200", Longitude, []string{ErrNotLongitude}},
		{"-200", Longitude, []string{ErrNotLongitude}},
		{"a", Longitude, []string{ErrNotNumber}},
		{"a1", Longitude, []string{ErrNotNumber}},
	})
}

func TestMax(t *testing.T) {
	testRule(t, "Max", ruleCases{
		{"0", Max(3), nil},
		{"1", Max(3), nil},
		{"3", Max(3), nil},
		{"-123.45", Max(3), nil},
		{"5", Max(3), []string{ErrMax + ":3"}},
		{"a", Max(3), []string{ErrNotNumber}},
		{"a1", Max(3), []string{ErrNotNumber}},
		{".", Max(-1), []string{ErrNotNumber}},
	})
}

func TestMaxLen(t *testing.T) {
	testRule(t, "MaxLen", ruleCases{
		{"a", MaxLen(3), nil},
		{"   ", MaxLen(3), nil},
		{"aaaa", MaxLen(3), []string{ErrMaxLen + ":3"}},
	})
}

func TestMin(t *testing.T) {
	testRule(t, "Min", ruleCases{
		{"3", Min(3), nil},
		{"+123.45", Min(3), nil},
		{"1", Min(3), []string{ErrMin + ":3"}},
		{"a", Min(3), []string{ErrNotNumber}},
		{"a1", Min(3), []string{ErrNotNumber}},
		{".", Min(3), []string{ErrNotNumber}},
	})
}

func TestMinLen(t *testing.T) {
	testRule(t, "MinLen", ruleCases{
		{"aaa", MinLen(3), nil},
		{"    ", MinLen(3), nil},
		{"a", MinLen(3), []string{ErrMinLen + ":3"}},
	})
}

func TestNumber(t *testing.T) {
	testRule(t, "Number", ruleCases{
		{"1", Number, nil},
		{"123", Number, nil},
		{"-123.45", Number, nil},
		{"a1", Number, []string{ErrNotNumber}},
		{"a", Number, []string{ErrNotNumber}},
		{".", Number, []string{ErrNotNumber}},
	})
}

func TestPhone(t *testing.T) {
	testRule(t, "Phone", ruleCases{
		{"0012345678901", Phone, nil},
		{"+12 (0) 345.67.89.01", Phone, nil},
		{"00123", Phone, []string{ErrNotPhone}},
		{"aaa", Phone, []string{ErrNotPhone}},
		{"aaaaaaaaaa", Phone, []string{ErrNotPhone}},
		{"aaa12345678901", Phone, []string{ErrNotPhone}},
	})
}

func TestRange(t *testing.T) {
	testRule(t, "Range", ruleCases{
		{"5", Range(3, 6), nil},
		{"1", Range(1, 1), nil},
		{"2", Range(3, 6), []string{ErrMin + ":3"}},
		{"0", Range(1, 1), []string{ErrMin + ":1"}},
		{"2", Range(1, 1), []string{ErrMax + ":1"}},
		{"a", Range(3, 6), []string{ErrNotNumber}},
		{"a1", Range(3, 6), []string{ErrNotNumber}},
		{".", Range(0, 0), []string{ErrNotNumber}},
	})
}

func TestRangeLen(t *testing.T) {
	testRule(t, "RangeLen", ruleCases{
		{"a", RangeLen(1, 1), nil},
		{"aaaaa", RangeLen(3, 6), nil},
		{"     ", RangeLen(3, 6), nil},
		{"a", RangeLen(3, 6), []string{ErrMinLen + ":3"}},
		{"aaa", RangeLen(1, 2), []string{ErrMaxLen + ":2"}},
	})
}
