package aspect

import (
	"fmt"
	"strings"
)

// Error holds meta and field-specific errors
type Error struct {
	Meta   []error
	Fields map[string]error
}

// AddMeta adds a meta error
func (err *Error) AddMeta(msg string, args ...interface{}) {
	err.Meta = append(err.Meta, fmt.Errorf(msg, args...))
}

// Error implements the built-in error interface
func (err Error) Error() string {
	n := len(err.Meta) + len(err.Fields)
	if n == 0 {
		return "aspect: no errors"
	}
	errors := []string{}
	if n == 1 {
		errors = append(errors, "aspect: 1 error:")
	} else {
		errors = append(errors, fmt.Sprintf("aspect: %d errors:", n))
	}

	for _, e := range err.Meta {
		errors = append(errors, fmt.Sprintf(" * %s", e.Error()))
	}
	for field, e := range err.Fields {
		errors = append(errors, fmt.Sprintf(" * %s: %s", field, e.Error()))
	}
	return strings.Join(errors, "\n")
}

// Exists returns true if there are any meta or field errors
func (err Error) Exists() bool {
	return len(err.Meta) > 0 || len(err.Fields) > 0
}

// SetField sets an error on a specific field
func (err *Error) SetField(field, msg string, args ...interface{}) {
	err.Fields[field] = fmt.Errorf(msg, args...)
}

// NewError creates a new Error
func NewError() *Error {
	return &Error{Fields: make(map[string]error)}
}
