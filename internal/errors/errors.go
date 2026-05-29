package errors

import "errors"

var (
	ErrEmptyField    = errors.New("empty field")
	ErrInvalidFormat = errors.New("invalid format")
)
