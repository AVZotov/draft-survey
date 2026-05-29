package parse

import (
	"net/http"
	"strconv"
	"strings"

	apperrors "github.com/AVZotov/draft-survey/internal/errors"
)

func Float(r *http.Request, name string, dest **float64) error {
	v := r.FormValue(name)
	v = strings.ReplaceAll(v, " ", "")
	if v == "" {
		*dest = nil
		return nil
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return apperrors.ErrInvalidFormat
	}
	*dest = &f
	return nil
}

func Int(r *http.Request, name string, dest **int) error {
	v := r.FormValue(name)
	if v == "" {
		*dest = nil
		return nil
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return apperrors.ErrInvalidFormat
	}
	*dest = &i
	return nil
}

func String(r *http.Request, name string) string {
	return r.FormValue(name)
}

func Bool(r *http.Request, name string) bool {
	v := r.FormValue(name)
	if v == "" {
		return false
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return b
}
