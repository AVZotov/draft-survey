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
		formattedErrors[i] = fmt.Sprintf("%s: %s", toLabel(err.Field), err.Message)
	}

	return strings.Join(formattedErrors, ", ")
}

// toLabel converts a Go struct field name to a human-readable label for
// user-facing validation messages.
func toLabel(field string) string {
	switch field {
	case "FirstName":
		return "First name"
	case "LastName":
		return "Last name"
	case "Email":
		return "Email"
	default:
		return field
	}
}
