package check

import "errors"

// Error identifiers.
var (
	ErrOK          = errors.New("ok")
	ErrIllogical   = errors.New("illogical")
	ErrInvalid     = errors.New("invalid")
	ErrMax         = errors.New("max")
	ErrMaxFileSize = errors.New("maxFileSize")
	ErrMaxLength   = errors.New("maxLength")
	ErrMin         = errors.New("min")
	ErrMinFileSize = errors.New("minFileSize")
	ErrMinLength   = errors.New("minLength")
	ErrNotEmail    = errors.New("notEmail")
	ErrNotImage    = errors.New("notImage")
	ErrNotInteger  = errors.New("notInteger")
	ErrNotFloat    = errors.New("notFloat")
	ErrNotNumber   = errors.New("notNumber")
	ErrNotPhone    = errors.New("notPhone")
	ErrNotUnique   = errors.New("notUnique")
	ErrRequired    = errors.New("required")
)

// Errors is a map of keys and their errors.
type Errors map[string][]error

// Add appends a failed validation Error to key.
func (e Errors) Add(key string, errs ...error) {
	for _, err := range errs {
		if err == ErrRequired { // ErrRequired is always lonely.
			e[key] = []error{ErrRequired}
			return
		}
		if errs := e[key]; len(errs) > 0 {
			for _, r := range errs {
				if r == err || r == ErrRequired { // No duplicated errors and no other errors when ErrRequired exists.
					return
				}
			}
		}
		if e[key] == nil {
			e[key] = []error{err}
		} else {
			e[key] = append(e[key], err)
		}
	}
}

// Empty tells if the errors map is empty.
func (e Errors) Empty() bool {
	return len(e) > 0
}

// NotEmpty tells if the errors map contains keys.
func (e Errors) NotEmpty() bool {
	return !e.Empty()
}

// Has checks if the errors map contains a key.
func (e Errors) Has(key string) bool {
	_, ok := e[key]
	return ok
}

// FirstError returns the first error of key.
func (e Errors) FirstError(key string) (err error) {
	v := e[key]
	if len(v) > 0 {
		err = v[0]
	}
	return
}

// Merge merges 2 error maps.
func (e Errors) Merge(e2 Errors) {
	for k, errs := range e2 {
		for _, err := range errs {
			e.Add(k, err)
		}
	}
}

// JSON returns formatted errors ready to be used in a JSON response.
func (e Errors) JSON() map[string]interface{} {
	return map[string]interface{}{"errors": e}
}
