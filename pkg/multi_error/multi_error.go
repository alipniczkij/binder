package multi_error

import "strings"

type MultiError struct {
	errors []string
}

func New() *MultiError {
	return &MultiError{
		errors: make([]string, 0),
	}
}

func (e *MultiError) Error() string {
	return strings.Join(e.errors, "\n")
}

func (e *MultiError) IsEmpty() bool {
	return len(e.errors) == 0
}

func (e *MultiError) Append(err error) {
	e.errors = append(e.errors, err.Error())
}
