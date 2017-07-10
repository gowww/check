package check

/*
import (
	"reflect"
	"testing"
)

func TestErrorsNotEmpty(t *testing.T) {
	errs := Errors{}
	if errs.NotEmpty() {
		t.Fail()
	}
}

func TestErrorsHas(t *testing.T) {
	errs := Errors{"one": nil}
	want := false
	got := errs.Has("two")
	if want != got {
		t.Errorf("False Errors.Has: want %t, got %t", want, got)
	}

	want = true
	got = errs.Has("one")
	if want != got {
		t.Errorf("True Errors.Has: want %t, got %t", want, got)
	}
}

func TestErrorsFirst(t *testing.T) {
	errs := Errors{"one": nil}
	want := ""
	got := errs.First("one")
	if want != got {
		t.Errorf("Empty Errors.First: want %s, got %s", want, got)
	}

	errs["one"] = []string{ErrRequired}
	want = ErrRequired
	got = errs.First("one")
	if want != got {
		t.Errorf("Existent Errors.First: want %s, got %s", want, got)
	}
}

func TestErrorsMerge(t *testing.T) {
	errs := Errors{"one": {ErrNotUnique}}
	errs.Merge(Errors{"one": {ErrNotAlpha}, "two": {ErrRequired}})
	want := Errors{"one": {ErrNotUnique, ErrNotAlpha}, "two": {ErrRequired}}
	if !reflect.DeepEqual(want, errs) {
		t.Errorf("Errors.Merge: want %s, got %s", want, errs)
	}
}

func TestErrorsJSON(t *testing.T) {
	errs := Errors{"one": {ErrNotUnique}}
	want := map[string]interface{}{"errors": errs}
	got := errs.JSON()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Errors.JSON: want %s, got %s", want, got)
	}
}
*/
