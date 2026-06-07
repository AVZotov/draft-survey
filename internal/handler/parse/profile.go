package parse

import (
	"io"
	"net/http"

	"github.com/AVZotov/draft-survey/internal/handler/fields"
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/go-playground/form/v4"
)

func (d *Decoder) Profile(r *http.Request) (*types.User, form.DecodeErrors, error) {
	user := new(types.User)
	de, err := d.decode(r, user)
	return user, de, err
}

func (d *Decoder) ProfileSignature(r *http.Request) ([]byte, error) {
	file, _, err := r.FormFile(fields.FieldSignature)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}
