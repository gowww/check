package check

import (
	"reflect"
	"testing"
)

func TestEmail(t *testing.T) {
	cases := []struct {
		v    string
		want []string
	}{
		{"a@example.com", nil},
		{"a+a@example.com", nil},
		{"a@a.a", nil},
		{"a@a.a", []string{ErrNotEmail}},
		{"a+a@a.a", []string{ErrNotEmail}},
		{"@a.a", []string{ErrNotEmail}},
		{"a@a.", []string{ErrNotEmail}},
		{"a@a", []string{ErrNotEmail}},
		{"a.a", []string{ErrNotEmail}},
		{"foobar", []string{ErrNotEmail}},
	}
	checker := Checker{"email": {Email}}
	for _, c := range cases {
		errs := checker.Check(map[string][]string{"email": {c.v}})
		got := errs["email"]
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Email(%v): want %v, got %v", c.v, c.want, got)
		}
	}
}

func TestInteger(t *testing.T) {
	cases := []struct {
		v    string
		want []string
	}{
		{"123", nil},
		{"123.45", []string{ErrNotInteger}},
		{"123,45", []string{ErrNotInteger}},
		{"a123", []string{ErrNotInteger}},
		{"123a", []string{ErrNotInteger}},
		{"1 123", []string{ErrNotInteger}},
		{"foobar", []string{ErrNotInteger}},
	}
	checker := Checker{"integer": {Integer}}
	for _, c := range cases {
		errs := checker.Check(map[string][]string{"integer": {c.v}})
		got := errs["integer"]
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Integer(%v): want %v, got %v", c.v, c.want, got)
		}
	}
}

func TestPhone(t *testing.T) {
	cases := []struct {
		v    string
		want []string
	}{
		{"0012345678901", nil},
		{"+12 (0) 345.67.89.01", nil},
		{"00123", []string{ErrNotPhone}},
		{"foobar", []string{ErrNotPhone}},
	}
	checker := Checker{"phone": {Phone}}
	for _, c := range cases {
		errs := checker.Check(map[string][]string{"phone": {c.v}})
		got := errs["phone"]
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Phone(%v): want %v, got %v", c.v, c.want, got)
		}
	}
}
