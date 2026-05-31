package validation

import (
	"fmt"
	"strings"
)

type FieldError struct {
	Field     string
	Namespace string
	Message   string
}

type FieldErrors []FieldError

func (errs FieldErrors) Error() string {
	formattedErrors := make([]string, len(errs))
	for i, err := range errs {
		formattedErrors[i] = fmt.Sprintf("%s: %s", err.Namespace, err.Message)
	}

	return strings.Join(formattedErrors, ", ")
}
