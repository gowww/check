package check

import (
	"mime/multipart"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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
		{"example", nil, URL, nil},
		{"example.com/", nil, URL, nil},
		{"example.com/?", nil, URL, nil},
		{"www.example.com", nil, URL, nil},
		{"example.com", nil, URL, nil},
		{"?example.com", nil, URL, []string{ErrNotURL}},
		{"//example.com", nil, URL, []string{ErrNotURL}},
		{"://example.com", nil, URL, []string{ErrNotURL}},
		{".", nil, URL, []string{ErrNotURL}},
		{"", nil, URL, []string{ErrNotURL}},
		// From https://github.com/asaskevich/govalidator/blob/master/validator_test.go#L572
		{"http://foo.bar#com", nil, URL, nil},
		{"http://foobar.com", nil, URL, nil},
		{"https://foobar.com", nil, URL, nil},
		{"foobar.com", nil, URL, nil},
		{"http://foobar.coffee/", nil, URL, nil},
		{"http://foobar.中文网/", nil, URL, nil},
		{"http://foobar.org/", nil, URL, nil},
		{"http://foobar.ORG", nil, URL, nil},
		{"http://foobar.org:8080/", nil, URL, nil},
		{"ftp://foobar.ru/", nil, URL, nil},
		{"ftp.foo.bar", nil, URL, nil},
		{"http://user:pass@www.foobar.com/", nil, URL, nil},
		{"http://user:pass@www.foobar.com/path/file", nil, URL, nil},
		{"http://127.0.0.1/", nil, URL, nil},
		{"http://duckduckgo.com/?q=%2F", nil, URL, nil},
		{"http://localhost:3000/", nil, URL, nil},
		{"http://foobar.com/?foo=bar#baz=qux", nil, URL, nil},
		{"http://foobar.com?foo=bar", nil, URL, nil},
		{"http://www.xn--froschgrn-x9a.net/", nil, URL, nil},
		{"http://foobar.com/a-", nil, URL, nil},
		{"http://foobar.پاکستان/", nil, URL, nil},
		{"xyz://foobar.com", nil, URL, nil},
		{"rtmp://foobar.com", nil, URL, nil},
		{"http://localhost:3000/", nil, URL, nil},
		{"http://foobar.com#baz=qux", nil, URL, nil},
		{"http://foobar.com/t$-_.+!*\\'(),", nil, URL, nil},
		{"http://www.foobar.com/~foobar", nil, URL, nil},
		{"http://www.foo---bar.com/", nil, URL, nil},
		{"http://r6---snnvoxuioq6.googlevideo.com", nil, URL, nil},
		{"mailto:someone@example.com", nil, URL, nil},
		{"irc://irc.server.org/channel", nil, URL, nil},
		{"http://foo.bar.org", nil, URL, nil},
		{"http://www.foo.bar.org", nil, URL, nil},
		{"http://www.foo.co.uk", nil, URL, nil},
		{"http://myservice.:9093/", nil, URL, nil},
		{"https://pbs.twimg.com/profile_images/560826135676588032/j8fWrmYY_normal.jpeg", nil, URL, nil},
		{"http://prometheus-alertmanager.service.q:9093", nil, URL, nil},
		{"https://www.logn-123-123.url.with.sigle.letter.d:12345/url/path/foo?bar=zzz#user", nil, URL, nil},
		{"http://me.example.com", nil, URL, nil},
		{"http://www.me.example.com", nil, URL, nil},
		{"https://farm6.static.flickr.com", nil, URL, nil},
		{"https://zh.wikipedia.org/wiki/Wikipedia:%E9%A6%96%E9%A1%B5", nil, URL, nil},
		{"google", nil, URL, nil},
		{"http://hyphenated-host-name.example.co.in", nil, URL, nil},
		{"http://www.domain-can-have-dashes.com", nil, URL, nil},
		{"http://m.abcd.com/test.html", nil, URL, nil},
		{"http://m.abcd.com/a/b/c/d/test.html?args=a&b=c", nil, URL, nil},
		{"http://[::1]:9093", nil, URL, nil},
		{"http://[2001:db8:a0b:12f0::1]/index.html", nil, URL, nil},
		{"http://[1200:0000:AB00:1234:0000:2552:7777:1313]", nil, URL, nil},
		{"http://user:pass@[::1]:9093/a/b/c/?a=v#abc", nil, URL, nil},
		{"https://127.0.0.1/a/b/c?a=v&c=11d", nil, URL, nil},
		{"https://foo_bar.example.com", nil, URL, nil},
		{"http://foo_bar.example.com", nil, URL, nil},
		{"http://foo_bar_fizz_buzz.example.com", nil, URL, nil},
		{"foo_bar.example.com", nil, URL, nil},
		{"foo_bar_fizz_buzz.example.com", nil, URL, nil},
		{"http://hello_world.example.com", nil, URL, nil},
		{"http://foobar.c_o_m", nil, URL, []string{ErrNotURL}},
		{"", nil, URL, []string{ErrNotURL}},
		{".com", nil, URL, []string{ErrNotURL}},
		{"http://www.foo_bar.com/", nil, URL, []string{ErrNotURL}},
		{"http://www.-foobar.com/", nil, URL, []string{ErrNotURL}},
		{"irc://#channel@network", nil, URL, []string{ErrNotURL}},
		{"/abs/test/dir", nil, URL, []string{ErrNotURL}},
		{"./rel/test/dir", nil, URL, []string{ErrNotURL}},
		{"http://foo^bar.org", nil, URL, []string{ErrNotURL}},
		{"http://foo&*bar.org", nil, URL, []string{ErrNotURL}},
		{"http://foo&bar.org", nil, URL, []string{ErrNotURL}},
		{"http://foo bar.org", nil, URL, []string{ErrNotURL}},
		{"foo", nil, URL, []string{ErrNotURL}},
		{"http://.foo.com", nil, URL, []string{ErrNotURL}},
		{"http://,foo.com", nil, URL, []string{ErrNotURL}},
		{",foo.com", nil, URL, []string{ErrNotURL}},
		{"http://cant-end-with-hyphen-.example.com", nil, URL, []string{ErrNotURL}},
		{"http://-cant-start-with-hyphen.example.com", nil, URL, []string{ErrNotURL}},
		{"http://[::1]:909388", nil, URL, []string{ErrNotURL}},
		{"1200::AB00:1234::2552:7777:1313", nil, URL, []string{ErrNotURL}},
		{"http://_cant_start_with_underescore", nil, URL, []string{ErrNotURL}},
		{"http://cant_end_with_underescore_", nil, URL, []string{ErrNotURL}},
	})
}

func BenchmarkAlphaRegexp(b *testing.B) {
	re := regexp.MustCompile("^[a-zA-Z]+$")
	for i := 0; i < b.N; i++ {
		re.MatchString(testAlpha)
	}
}

func BenchmarkAlphaStrconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(testAlpha); i++ {
			if testAlpha[i] < 65 || testAlpha[i] > 90 && testAlpha[i] < 97 || testAlpha[i] > 122 {
			}
		}
	}
}

func BenchmarkIntRegexp(b *testing.B) {
	re := regexp.MustCompile("^(?:[-+]?(?:0|[1-9][0-9]*))$")
	for i := 0; i < b.N; i++ {
		re.MatchString(testInt)
	}
}

func BenchmarkIntStrconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.Atoi(testInt)
	}
}

func BenchmarkNumberRegexp(b *testing.B) {
	re := regexp.MustCompile("^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$")
	for i := 0; i < b.N; i++ {
		re.MatchString(testFloat)
	}
}

func BenchmarkNumberStrconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.ParseFloat(testFloat, 64)
	}
}

func BenchmarkURLRegexp(b *testing.B) {
	reIP := `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	reURLSchema := `((ftp|tcp|udp|wss?|https?):\/\/)`
	reURLIP := `([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`
	reURLUsername := `(\S+(:\S*)?@)`
	reURLSubdomain := `((www\.)|([a-zA-Z0-9]([-\.][-\._a-zA-Z0-9]+)*))`
	reURLPort := `(:(\d{1,5}))`
	reURLPath := `((\/|\?|#)[^\s]*)`
	re := regexp.MustCompile(`^` + reURLSchema + `?` + reURLUsername + `?` + `((` + reURLIP + `|(\[` + reIP + `\])|(([a-zA-Z0-9]([a-zA-Z0-9-_]+)?[a-zA-Z0-9]([-\.][a-zA-Z0-9]+)*)|(` + reURLSubdomain + `?))?(([a-zA-Z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-zA-Z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-zA-Z\x{00a1}-\x{ffff}]{1,}))?))\.?` + reURLPort + `?` + reURLPath + `?$`)
	for i := 0; i < b.N; i++ {
		re.MatchString(testURL)
	}
}

func BenchmarkURLParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if len(testURL) < 4 {
		}
		testURL = strings.Replace(testURL, "127.0.0.1", "localhost", 1)
		if i := strings.IndexByte(testURL, '#'); i != -1 {
			testURL = testURL[:i]
		}
		if !strings.Contains(testURL, "://") {
			testURL = "http://" + testURL
		}
		u, err := url.ParseRequestURI(testURL)
		if err != nil || u.Host == "" || u.Host[0] == '-' || strings.Contains(u.Host, ".-") || strings.Contains(u.Host, "-.") {
		}
		parts := strings.Split(u.Host, ".")
		if parts[0] == "" || len(parts[len(parts)-1]) < 2 || len(parts[len(parts)-1]) > 63 {
		}
		var domain string
		if len(parts) > 2 {
			domain = strings.Join(parts[len(parts)-2:], ".")
		} else {
			domain = strings.Join(parts, ".")
		}
		if strings.ContainsAny(domain, "_,!&") {
		}
		if strings.Count(domain, "::") > 1 { // Only 1 substitution ("::") allowed in IPv6 address.
		}
		parts = strings.Split(domain, ":")
		port, err := strconv.Atoi(parts[len(parts)-1])
		if err == nil && (port < 1 || port > 65535) {
		}
	}
}
