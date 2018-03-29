package check

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gowww/i18n"
	"golang.org/x/text/language"
)

// An ErrorID defines a standard and translatable error to be used during a check.
type ErrorID struct {
	ID      string
	Locales map[language.Tag]string
}

// Error identifiers.
// The first locale in Locales map is used when no one matched.
var (
	ErrBadFileType = &ErrorID{ID: "badFileType", Locales: map[language.Tag]string{
		language.English: "Only these file types are accepted: %v.",
		language.French:  "Seul ces types de fichier sont acceptés: %v.",
	}}
	ErrIllogical = &ErrorID{ID: "illogical", Locales: map[language.Tag]string{
		language.English: "This value is illogical.",
		language.French:  "Cette valeur est illogique.",
	}}
	ErrInvalid = &ErrorID{ID: "invalid", Locales: map[language.Tag]string{
		language.English: "This value is invalid.",
		language.French:  "Cette valeur est invalide.",
	}}
	ErrMax = &ErrorID{ID: "max", Locales: map[language.Tag]string{
		language.English: "The maximal value is %v.",
		language.French:  "La valeur maximale est de %v",
	}}
	ErrMaxFileSize = &ErrorID{ID: "maxFileSize", Locales: map[language.Tag]string{
		language.English: "File size is over %v.",
		language.French:  "La taille du fichier dépasse %v.",
	}}
	ErrMaxLen = &ErrorID{ID: "maxLen", Locales: map[language.Tag]string{
		language.English: "The value exceeds %v characters.",
		language.French:  "La valeur dépasse %v caractères.",
	}}
	ErrMin = &ErrorID{ID: "min", Locales: map[language.Tag]string{
		language.English: "The minimal value is %v.",
		language.French:  "La valeur minimale est de %v",
	}}
	ErrMinFileSize = &ErrorID{ID: "minFileSize", Locales: map[language.Tag]string{
		language.English: "File size must be at least %v.",
		language.French:  "La taille du fichier doit être d'au moins %v.",
	}}
	ErrMinLen = &ErrorID{ID: "minLen", Locales: map[language.Tag]string{
		language.English: "The value must have more than %v characters.",
		language.French:  "La veleur doit comporter au moins %v caractères.",
	}}
	ErrNotAlpha = &ErrorID{ID: "notAlpha", Locales: map[language.Tag]string{
		language.English: "It's not a letters-only string.",
		language.French:  "Ce n'est pas une suite de lettres (uniquement).",
	}}
	ErrNotAlphanumeric = &ErrorID{ID: "notAlphanumeric", Locales: map[language.Tag]string{
		language.English: "It's not an alphanumeric-only string.",
		language.French:  "Ce n'est pas une suite alphanumérique (uniquement).",
	}}
	ErrNotEmail = &ErrorID{ID: "notEmail", Locales: map[language.Tag]string{
		language.English: "It's not an email.",
		language.French:  "Ce n'est pas un e-mail.",
	}}
	ErrNotFloat = &ErrorID{ID: "notFloat", Locales: map[language.Tag]string{
		language.English: "It's not a floating point number.",
		language.French:  "Ce n'est pas un nombre à virgule.",
	}}
	ErrNotImage = &ErrorID{ID: "notImage", Locales: map[language.Tag]string{
		language.English: "It's not an image.",
		language.French:  "Ce n'est pas une image.",
	}}
	ErrNotInteger = &ErrorID{ID: "notInteger", Locales: map[language.Tag]string{
		language.English: "It's not a integer number.",
		language.French:  "Ce n'est pas un nombre entier.",
	}}
	ErrNotLatitude = &ErrorID{ID: "notLatitude", Locales: map[language.Tag]string{
		language.English: "It's not a latitude.",
		language.French:  "Ce n'est pas une latitude.",
	}}
	ErrNotLongitude = &ErrorID{ID: "notLongitude", Locales: map[language.Tag]string{
		language.English: "It's not a longitude.",
		language.French:  "Ce n'est pas une longitude.",
	}}
	ErrNotNumber = &ErrorID{ID: "notNumber", Locales: map[language.Tag]string{
		language.English: "It's not a number.",
		language.French:  "Ce n'est pas un nombre.",
	}}
	ErrNotPhone = &ErrorID{ID: "notPhone", Locales: map[language.Tag]string{
		language.English: "It's not a phone number.",
		language.French:  "Ce n'est pas un numéro de téléphone.",
	}}
	ErrNotSame = &ErrorID{ID: "notSame", Locales: map[language.Tag]string{
		language.English: "The value must equals these fields: %v.",
		language.French:  "La valeur doit être identique aux champs suivants: %v.",
	}}
	ErrNotURL = &ErrorID{ID: "notURL", Locales: map[language.Tag]string{
		language.English: "It's not a web address.",
		language.French:  "Ce n'est pas une adresse web.",
	}}
	ErrNotUnique = &ErrorID{ID: "notUnique", Locales: map[language.Tag]string{
		language.English: "This value already exists.",
		language.French:  "Cette valeur existe déjà.",
	}}
	ErrRequired = &ErrorID{ID: "required", Locales: map[language.Tag]string{
		language.English: "A value is required.",
		language.French:  "Une valeur est requise.",
	}}
	ErrWrongPassword = &ErrorID{ID: "password", Locales: map[language.Tag]string{
		language.English: "The password is wrong.",
		language.French:  "Le mot de passe est incorrect.",
	}}
)

// An Error is a checking error from a rule, with rule's variables.
type Error struct {
	Error *ErrorID
	Args  []interface{}
}

// T returns an Error translation from an i18n.Translator (from key "error" + title case error ID, like "errorNotImage").
// If custom translation is not defined, the default ErrorID translation is used.
func (e *Error) T(t *i18n.Translator) string {
	if t == nil {
		return e.TDefault(language.English)
	}
	s := t.T("error"+strings.Title(e.Error.ID), e.Args...)
	if s == "" {
		return e.TDefault(t.Locale())
	}
	return s

}

// TDefault returns the default translation of Error.
// If Error has no translations, the raw string representation.
func (e *Error) TDefault(l language.Tag) string {
	if len(e.Error.Locales) == 0 {
		return e.String()
	}
	t := make([]language.Tag, 0, len(e.Error.Locales))
	for lt := range e.Error.Locales {
		t = append(t, lt)
	}
	l, _, _ = language.NewMatcher(t).Match(l)

	for i, arg := range e.Args {
		if ta, ok := arg.(i18n.Translatable); ok {
			e.Args[i] = ta.T(l) // Translate translatable arguments.
		}
	}
	return fmt.Sprintf(e.Error.Locales[l], e.Args...)

}

func (e *Error) String() string {
	if len(e.Args) == 0 {
		return e.Error.ID
	}
	return fmt.Sprintf("%s:%s", e.Error.ID, strings.Join(interfacesToStrings(e.Args), ","))
}

// Errors is a map of keys and their errors.
type Errors map[string][]*Error

// Add appends a failed validation Error to key.
func (e Errors) Add(key string, err *Error) {
	if err.Error == ErrRequired { // ErrRequired is always lonely.
		e[key] = []*Error{{Error: ErrRequired}}
		return
	}
	if checkErrors := e[key]; len(checkErrors) > 0 {
		for _, ce := range checkErrors {
			if ce.Error == err.Error || ce.Error == ErrRequired { // No duplicated errors and no other errors when ErrRequired exists.
				return
			}
		}
	}
	if e[key] == nil {
		e[key] = []*Error{err}
	} else {
		e[key] = append(e[key], err)
	}
}

// Empty tells if the errors map contains no keys.
func (e Errors) Empty() bool {
	return len(e) == 0
}

// NotEmpty tells if the errors map contains keys.
func (e Errors) NotEmpty() bool {
	return !e.Empty()
}

// Has tells if the errors map contains a key.
func (e Errors) Has(key string) bool {
	_, ok := e[key]
	return ok
}

// First returns the first error for key.
// If the key doesn't exist, nil.
func (e Errors) First(key string) *Error {
	v := e[key]
	if len(v) == 0 {
		return nil
	}
	return v[0]
}

func (e Errors) String() string {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "\t")
	if err := enc.Encode(e.StringMap()); err != nil {
		panic(err)
	}
	return buf.String()
}

// Merge merges two Errors maps.
func (e Errors) Merge(e2 Errors) {
	for k, errs := range e2 {
		for _, err := range errs {
			e.Add(k, err)
		}
	}
}

// StringMap returns the errors as a readable string map.
func (e Errors) StringMap() map[string][]string {
	m := make(map[string][]string, len(e))
	for k, errs := range e {
		m[k] = make([]string, 0, len(errs))
		for _, err := range errs {
			m[k] = append(m[k], err.String())
		}
	}
	return m
}

// JSON returns the errors map under the "errors" key, ready to be encoded.
func (e Errors) JSON() interface{} {
	return map[string]map[string][]string{"errors": e.StringMap()}
}

// TranslatedErrors is a map of keys and their translated errors.
type TranslatedErrors map[string][]string

// T returns a tranlated Errors map.
// If t is nil, built-in translations are used.
func (e Errors) T(t *i18n.Translator) TranslatedErrors {
	te := make(TranslatedErrors, len(e))
	for k, errs := range e {
		te[k] = make([]string, 0, len(errs))
		for _, err := range errs {
			te[k] = append(te[k], err.T(t))
		}
	}
	return te
}

// Has tells if the translated errors map contains a key.
func (e TranslatedErrors) Has(key string) bool {
	_, ok := e[key]
	return ok
}

// First returns the first translated error for key.
// If the key doesn't exist, an empty string.
func (e TranslatedErrors) First(key string) string {
	v := e[key]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}
