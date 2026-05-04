package validation

type Validator interface {
	Validate(data any) FieldErrors
}
