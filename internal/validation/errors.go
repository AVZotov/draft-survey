package validation

import (
	"fmt"
	"strings"
)

type FieldError struct {
	Field   string
	Message string
}

type FieldErrors []FieldError

func (errs FieldErrors) Error() string {
	formattedErrors := make([]string, len(errs))
	for i, err := range errs {
		formattedErrors[i] = fmt.Sprintf("%s: %s", err.Field, err.Message)
	}

	return strings.Join(formattedErrors, ", ")
}
