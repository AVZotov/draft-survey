package parse

import (
	"net/http"
	"strconv"
	"strings"
)

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

// str returns the form value for key, or an empty string if not present
func (d *Decoder) str(r *http.Request, key string) string {
	return r.FormValue(key)
}

// pFloat parses key as *float64.
// bool reports whether the field was present in the request.
// An empty value yields (nil, true); an unparsable value yields (nil, true) as well.
func (d *Decoder) pFloat(r *http.Request, key string) (*float64, bool) {
	if !r.PostForm.Has(key) {
		return nil, false
	}
	v := normalizeNumber(r.PostForm.Get(key))
	if v == "" {
		return nil, true
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil, true
	}
	return &f, true
}

// float parses key as float64.
// bool reports whether the field was present in the request.
func (d *Decoder) float(r *http.Request, key string) (float64, bool) {
	if !r.PostForm.Has(key) {
		return 0, false
	}
	v := normalizeNumber(r.PostForm.Get(key))
	if v == "" {
		return 0, true
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, true
	}
	return f, true
}

// pInt parses key as *int.
// bool reports whether the field was present in the request.
func (d *Decoder) pInt(r *http.Request, key string) (*int, bool) {
	if !r.PostForm.Has(key) {
		return nil, false
	}
	v := normalizeNumber(r.PostForm.Get(key))
	if v == "" {
		return nil, true
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return nil, true
	}
	return &i, true
}

// int parses key as int.
// bool reports whether the field was present in the request.
func (d *Decoder) int(r *http.Request, key string) (int, bool) {
	if !r.PostForm.Has(key) {
		return 0, false
	}
	v := normalizeNumber(r.PostForm.Get(key))
	if v == "" {
		return 0, true
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, true
	}
	return i, true
}

// boolean reports a checkbox state: absent fields are false,
// values "on", "true" or "1" are true.
func (d *Decoder) boolean(r *http.Request, key string) bool {
	switch r.FormValue(key) {
	case "on", "true", "1":
		return true
	default:
		return false
	}
}

func normalizeNumber(v string) string {
	v = strings.ReplaceAll(v, " ", "")
	v = strings.ReplaceAll(v, ",", ".")
	return v
}
