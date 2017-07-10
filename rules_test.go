package check

/*
import (
	"mime/multipart"
	"net/url"
	"reflect"
	"testing"
)

type valueCases []struct {
	v      string
	values url.Values
	rule   Rule
	want   []string
}

type fileCases []struct {
	f      string
	values map[string][]*multipart.FileHeader
	rule   Rule
	want   []string
}

func testValueRule(t *testing.T, name string, cases valueCases) {
	for _, c := range cases {
		errs := make(Errors)
		if c.values == nil {
			c.values = make(url.Values)
		}
		c.values[""] = []string{c.v}
		form := &multipart.Form{Value: c.values}
		c.rule(errs, form, "")
		if !reflect.DeepEqual(c.want, errs[""]) {
			t.Errorf("%s(%q): want %v, got %v", name, c.v, c.want, errs)
		}
	}
}

func TestAlpha(t *testing.T) {
	testValueRule(t, "Alpha", valueCases{
		{"a", nil, Alpha, nil},
		{"a", nil, Alpha, nil},
		{"aa", nil, Alpha, nil},
		{"1", nil, Alpha, []string{ErrNotAlpha}},
		{"a.a", nil, Alpha, []string{ErrNotAlpha}},
		{"a@a", nil, Alpha, []string{ErrNotAlpha}},
	})
}

func TestEmail(t *testing.T) {
	testValueRule(t, "Email", valueCases{
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
	testValueRule(t, "Email", valueCases{
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
	testValueRule(t, "Latitude", valueCases{
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
	testValueRule(t, "Longitude", valueCases{
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
	testValueRule(t, "Max", valueCases{
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
	testValueRule(t, "MaxLen", valueCases{
		{"a", nil, MaxLen(3), nil},
		{"   ", nil, MaxLen(3), nil},
		{"aaaa", nil, MaxLen(3), []string{ErrMaxLen + ":3"}},
	})
}

func TestMin(t *testing.T) {
	testValueRule(t, "Min", valueCases{
		{"3", nil, Min(3), nil},
		{"+123.45", nil, Min(3), nil},
		{"1", nil, Min(3), []string{ErrMin + ":3"}},
		{"a", nil, Min(3), []string{ErrNotNumber}},
		{"a1", nil, Min(3), []string{ErrNotNumber}},
		{".", nil, Min(3), []string{ErrNotNumber}},
	})
}

func TestMinLen(t *testing.T) {
	testValueRule(t, "MinLen", valueCases{
		{"aaa", nil, MinLen(3), nil},
		{"    ", nil, MinLen(3), nil},
		{"a", nil, MinLen(3), []string{ErrMinLen + ":3"}},
	})
}

func TestNumber(t *testing.T) {
	testValueRule(t, "Number", valueCases{
		{"1", nil, Number, nil},
		{"123", nil, Number, nil},
		{"-123.45", nil, Number, nil},
		{"a1", nil, Number, []string{ErrNotNumber}},
		{"a", nil, Number, []string{ErrNotNumber}},
		{".", nil, Number, []string{ErrNotNumber}},
	})
}

func TestPhone(t *testing.T) {
	testValueRule(t, "Phone", valueCases{
		{"0012345678901", nil, Phone, nil},
		{"+12 (0) 345.67.89.01", nil, Phone, nil},
		{"00123", nil, Phone, []string{ErrNotPhone}},
		{"aaa", nil, Phone, []string{ErrNotPhone}},
		{"aaaaaaaaaa", nil, Phone, []string{ErrNotPhone}},
		{"aaa12345678901", nil, Phone, []string{ErrNotPhone}},
	})
}

func TestRange(t *testing.T) {
	testValueRule(t, "Range", valueCases{
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
	testValueRule(t, "RangeLen", valueCases{
		{"a", nil, RangeLen(1, 1), nil},
		{"aaaaa", nil, RangeLen(3, 6), nil},
		{"     ", nil, RangeLen(3, 6), nil},
		{"a", nil, RangeLen(3, 6), []string{ErrMinLen + ":3"}},
		{"aaa", nil, RangeLen(1, 2), []string{ErrMaxLen + ":2"}},
	})
}

func TestSame(t *testing.T) {
	testValueRule(t, "Range", valueCases{
		{"v", url.Values{"k": {"v"}, "l": {"v"}}, Same("k", "l"), nil},
		{"v", url.Values{"k": {"v"}}, Same("x"), []string{ErrNotSame + ":x"}},
		{"x", url.Values{"k": {"v"}}, Same("k"), []string{ErrNotSame + ":k"}},
		{"v", url.Values{"k": {"v"}, "l": {"x"}}, Same("k", "l"), []string{ErrNotSame + ":k,l"}},
		{"v", nil, Same("k"), []string{ErrNotSame + ":k"}},
	})
}

func TestURL(t *testing.T) {
	testValueRule(t, "URL", valueCases{
		{"http://example.com", nil, URL, nil},
		{"https://example.com", nil, URL, nil},
		{"http://example.com/?foo=bar#baz=qux", nil, URL, nil},
		{"http://example.com/?q=%2F", nil, URL, nil},
		{"http://example.com?foo=bar", nil, URL, nil},
		{"http://example.com#com", nil, URL, nil},
		{"http://user:pass@www.example.com/", nil, URL, nil},
		{"http://user:pass@www.example.com/path/file", nil, URL, nil},
		{"2001:db8:85a3::ac1f:8001/index.html", nil, URL, nil},
		{"user:pass@[::1]:9093/a/b/c/?a=v#abc", nil, URL, nil},
		{"127.0.0.1", nil, URL, nil},
		{"localhost:3000", nil, URL, nil},
		{"foo://example.com", nil, URL, nil},
		{"mailto:foo@example.com", nil, URL, nil},
		{"www.ex.ample.com", nil, URL, nil},
		{"www.example.com", nil, URL, nil},
		{"www.ex---ample.com", nil, URL, nil},
		{"www.ex--am--ple.com", nil, URL, nil},
		{"w_w-w.example.com", nil, URL, nil},
		{"127.0.0.1", nil, URL, nil},
		{"example", nil, URL, nil},
		{"example.:9093/", nil, URL, nil},
		{"example:80", nil, URL, nil},
		{"example.COM", nil, URL, nil},
		{"example.中文网", nil, URL, nil},
		{"example.c", nil, URL, nil},
		{"example.com/", nil, URL, nil},
		{"example.com/?", nil, URL, nil},
		{"example.com/$-_.+!*\\'(),", nil, URL, nil},
		{"example.com/a-", nil, URL, nil},
		{"example.com#foo=bar", nil, URL, nil},
		{"example.com", nil, URL, nil},
		{"[::1]:8080", nil, URL, nil},
		{"example.c_o_m", nil, URL, []string{ErrNotURL}},
		{"www-.example.com", nil, URL, []string{ErrNotURL}},
		{"-www.example.com", nil, URL, []string{ErrNotURL}},
		{"_example", nil, URL, []string{ErrNotURL}},
		{"example_", nil, URL, []string{ErrNotURL}},
		{".example.com", nil, URL, []string{ErrNotURL}},
		{",example.com", nil, URL, []string{ErrNotURL}},
		{"?example.com", nil, URL, []string{ErrNotURL}},
		{"//example.com", nil, URL, []string{ErrNotURL}},
		{"://example.com", nil, URL, []string{ErrNotURL}},
		{"ex_ample.com", nil, URL, []string{ErrNotURL}},
		{"ex^ample.com", nil, URL, []string{ErrNotURL}},
		{"ex&ample.com", nil, URL, []string{ErrNotURL}},
		{"ex ample.com", nil, URL, []string{ErrNotURL}},
		{"ex&*ample.com", nil, URL, []string{ErrNotURL}},
		{"www.-example.com", nil, URL, []string{ErrNotURL}},
		{"irc://#channel@network", nil, URL, []string{ErrNotURL}},
		{"[::1]:70000", nil, URL, []string{ErrNotURL}},
		{"2001:db8::85a3::ac1f:8001", nil, URL, []string{ErrNotURL}},
		{"/abs/test/dir", nil, URL, []string{ErrNotURL}},
		{"./rel/test/dir", nil, URL, []string{ErrNotURL}},
		{".com", nil, URL, []string{ErrNotURL}},
		{".", nil, URL, []string{ErrNotURL}},
		{"", nil, URL, []string{ErrNotURL}},
	})
}
*/
