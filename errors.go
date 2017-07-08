package check

// Error identifiers.
const (
	ErrBadFileType  = "badFileType"
	ErrIllogical    = "illogical"
	ErrInvalid      = "invalid"
	ErrMax          = "max"
	ErrMaxFileSize  = "maxFileSize"
	ErrMaxLen       = "maxLen"
	ErrMin          = "min"
	ErrMinFileSize  = "minFileSize"
	ErrMinLen       = "minLen"
	ErrNotAlpha     = "notAlpha"
	ErrNotEmail     = "notEmail"
	ErrNotFloat     = "notFloat"
	ErrNotImage     = "notImage"
	ErrNotInteger   = "notInteger"
	ErrNotLatitude  = "notLatitude"
	ErrNotLongitude = "notLongitude"
	ErrNotNumber    = "notNumber"
	ErrNotPhone     = "notPhone"
	ErrNotSame      = "notSame"
	ErrNotURL       = "notURL"
	ErrNotUnique    = "notUnique"
	ErrRequired     = "required"
)

// Errors is a map of keys and their errors.
type Errors map[string][]string

// Add appends a failed validation Error to key.
func (e Errors) Add(key string, err string) {
	if err == ErrRequired { // ErrRequired is always lonely.
		e[key] = []string{ErrRequired}
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
		e[key] = []string{err}
	} else {
		e[key] = append(e[key], err)
	}
}

// Empty tells if the errors map is empty.
func (e Errors) Empty() bool {
	return len(e) == 0
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

// First returns the first error for key.
// If the key doesn't exist, an empty string is returned.
func (e Errors) First(key string) (err string) {
	v := e[key]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

// Merge merges 2 error maps.
func (e Errors) Merge(e2 Errors) {
	for k, errs := range e2 {
		for _, err := range errs {
			e.Add(k, err)
		}
	}
}

// JSON returns the errors map under the "errors" key, ready to be encoded.
func (e Errors) JSON() interface{} {
	return map[string]interface{}{"errors": e}
}
