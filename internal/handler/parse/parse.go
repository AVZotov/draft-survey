package parse

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/form/v4"
)

type Decoder struct {
	d *form.Decoder
}

func NewDecoder() *Decoder {
	d := form.NewDecoder()
	d.SetTagName("json")
	d.RegisterCustomTypeFunc(parsePFloat, new(float64))
	d.RegisterCustomTypeFunc(parseFloat, float64(0.0))
	d.RegisterCustomTypeFunc(parsePInt, new(int))
	d.RegisterCustomTypeFunc(parseInt, int(0))
	return &Decoder{
		d: d,
	}
}

func (d *Decoder) decode(r *http.Request, dst any) (form.DecodeErrors, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	err := d.d.Decode(dst, r.PostForm)
	if err != nil {
		var decodeErrors form.DecodeErrors
		if errors.As(err, &decodeErrors) {
			return decodeErrors, nil
		}
		return nil, err
	}
	return nil, nil
}

func parsePFloat(str []string) (any, error) {
	v := str[0]
	v = strings.ReplaceAll(v, " ", "")
	v = strings.ReplaceAll(v, ",", ".")
	if v == "" {
		return (*float64)(nil), nil
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return (*float64)(nil), err
	}
	return &f, nil
}

func parseFloat(str []string) (any, error) {
	var f float64
	v := str[0]
	v = strings.ReplaceAll(v, " ", "")
	v = strings.ReplaceAll(v, ",", ".")
	if v == "" {
		return f, nil
	}

	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return f, err
	}
	return f, nil
}

func parsePInt(str []string) (any, error) {
	v := str[0]
	v = strings.ReplaceAll(v, " ", "")
	v = strings.ReplaceAll(v, ",", ".")
	if v == "" {
		return (*int)(nil), nil
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return (*int)(nil), err
	}
	return &i, nil
}

func parseInt(str []string) (any, error) {
	var i int
	v := str[0]
	v = strings.ReplaceAll(v, " ", "")
	v = strings.ReplaceAll(v, ",", ".")
	if v == "" {
		return i, nil
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return i, err
	}
	return i, nil
}
