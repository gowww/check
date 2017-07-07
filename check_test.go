package check

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
)

var (
	testAlpha = "PeWGHxQBZfiMVQNbWGzkegklTRMHVmuO"
	testInt   = "-19001231"
	testFloat = "-19001231.558877"
)

func TestCheckerCheck(t *testing.T) {
	chk := &Checker{
		"email": {Required, Email},
		"phone": {Phone},
		"stars": {Required, Range(3, 5)},
	}
	errs := chk.Check(map[string][]string{
		"name":  {"foobar"},
		"phone": {"0012345678901"},
		"stars": {"2"},
	})
	fmt.Println(errs)
}

func BenchmarkAlphaRegexp(b *testing.B) {
	reInt := regexp.MustCompile("^[a-zA-Z]+$")
	for i := 0; i < b.N; i++ {
		reInt.MatchString(testAlpha)
	}
}

func BenchmarkAlphaStrconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// for i := range testAlpha {
		for i := 0; i < len(testAlpha); i++ {
			if testAlpha[i] < 65 || testAlpha[i] > 90 && testAlpha[i] < 97 || testAlpha[i] > 122 {
			}
		}
	}
}

func BenchmarkIntRegexp(b *testing.B) {
	reInt := regexp.MustCompile("^(?:[-+]?(?:0|[1-9][0-9]*))$")
	for i := 0; i < b.N; i++ {
		reInt.MatchString(testInt)
	}
}

func BenchmarkIntStrconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.Atoi(testInt)
	}
}

func BenchmarkNumberRegexp(b *testing.B) {
	reInt := regexp.MustCompile("^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$")
	for i := 0; i < b.N; i++ {
		reInt.MatchString(testFloat)
	}
}

func BenchmarkNumberStrconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.ParseFloat(testFloat, 64)
	}
}
