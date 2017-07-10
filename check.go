// Package check provides request form checking.
package check

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

var errNoFileProvided = errors.New("check: no file provided")

// A Checker contains keys with their checking rules.
type Checker map[string][]Rule

// Check makes the check for a multipart.Form (values and files) and returns errors.
//
// Result is guaranteed to be non-nil.
func (c Checker) Check(form *multipart.Form) Errors {
	errs := make(Errors)
	for key, rules := range c {
		for _, rule := range rules {
			rule(errs, form, key)
		}
	}
	return errs
}

// CheckValues makes the check for a values map (key to multiple values) and returns errors.
//
// Result is guaranteed to be non-nil.
func (c Checker) CheckValues(values map[string][]string) Errors {
	return c.Check(&multipart.Form{Value: values})
}

// CheckFiles makes the check for a files map (key to multiple files) and returns errors.
//
// Result is guaranteed to be non-nil.
func (c Checker) CheckFiles(files map[string][]*multipart.FileHeader) Errors {
	return c.Check(&multipart.Form{File: files})
}

// CheckRequest makes the check for an HTTP request and returns errors.
//
// Request data can have multiple values with the same key (or field).
// In tis case, all values are checked and if one fails, the error is set for the whole key.
//
// Result is guaranteed to be non-nil.
func (c Checker) CheckRequest(r *http.Request) Errors {
	if r.Form == nil {
		r.ParseMultipartForm(32 << 20) // 32 MB
	}
	form := &multipart.Form{Value: r.Form}
	if r.MultipartForm != nil {
		if r.MultipartForm.Value != nil {
			for k, v := range r.MultipartForm.Value {
				form.Value[k] = append(form.Value[k], v...)
			}
		}
		form.File = r.MultipartForm.File
	}
	return c.Check(form)
}

func fileSize(file *multipart.FileHeader) (int64, error) {
	if file == nil {
		return 0, errNoFileProvided
	}
	// TODO: In next Go versions, use new Size attribute (https://go-review.googlesource.com/c/39223).
	f, err := file.Open()
	if err != nil {
		return 0, err
	}
	var size int64
	switch ft := f.(type) {
	case *os.File:
		fi, _ := ft.Stat()
		size = fi.Size()
	default:
		size, _ = ft.Seek(0, io.SeekEnd)
		f.Seek(0, io.SeekStart) // Reset reader.
	}
	return size, nil
}

func fileType(file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", errNoFileProvided
	}
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Seek(0, io.SeekStart) // Reset reader.
	fh := make([]byte, 512)
	if _, err = f.Read(fh); err != nil {
		return "", err
	}
	ct := http.DetectContentType(fh)
	if i := strings.IndexByte(ct, ';'); i != -1 {
		ct = ct[:i]
	}
	return ct, nil
}

func sliceContainsString(ss []string, s string) bool {
	for _, e := range ss {
		if s == e {
			return true
		}
	}
	return false
}

func stringsToInterfaces(ss []string) []interface{} {
	ii := make([]interface{}, len(ss))
	for i := 0; i < len(ss); i++ {
		ii[i] = ss[i]
	}
	return ii
}

func interfacesToStrings(ii []interface{}) []string {
	ss := make([]string, len(ii))
	for i := 0; i < len(ii); i++ {
		ss[i] = fmt.Sprint(ii[i])
	}
	return ss
}
