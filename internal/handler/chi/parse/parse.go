package parse

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

var ErrEmptyField = errors.New("empty field")

func Float(r *http.Request, name string) (*float64, error) {
	v := r.FormValue(name)
	if v == "" {
		return nil, ErrEmptyField
	}
	v = strings.ReplaceAll(v, " ", "") // x-mask support
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func Int(r *http.Request, name string) (*int, error) {
	v := r.FormValue(name)
	if v == "" {
		return nil, ErrEmptyField
	}
	f, err := strconv.Atoi(v)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func String(r *http.Request, name string) (string, error) {
	v := r.FormValue(name)
	if v == "" {
		return "", ErrEmptyField
	}
	return v, nil
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
