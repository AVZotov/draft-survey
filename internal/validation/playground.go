package validation

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type PlaygroundValidator struct {
	validator *validator.Validate
}

func New() Validator {
	v := &PlaygroundValidator{
		validator: validator.New(),
	}
	err := v.validator.RegisterValidation("lbp", func(fl validator.FieldLevel) bool {
		if fl.Field().String() == "Half LBP" || fl.Field().String() == "Full LBP" {
			return true
		}
		return false
	})

	if err != nil {
		panic(err)
	}

	return v
}

func (pv *PlaygroundValidator) Validate(data any) FieldErrors {
	var fieldErrors FieldErrors
	err := pv.validator.Struct(data)
	if err != nil {
		if errs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			for _, fe := range errs {
				fieldErrors = append(fieldErrors, FieldError{
					Field:   fe.Field(),
					Message: messageFromTag(fe.Tag()),
				})
			}
		}
	}
	return fieldErrors
}

func messageFromTag(tag string) string {
	switch tag {
	case "required":
		return "required field is empty"
	default:
		return "invalid field value"
	}
}
