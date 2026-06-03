package parse

import (
	"io"
	"net/http"
	
	"github.com/AVZotov/draft-survey/internal/handler/fields"
	"github.com/AVZotov/draft-survey/internal/types"
)

func Profile(r *http.Request) (*types.User, error) {
	var user types.User
	if err := Decode(r, &user); err != nil {
		return nil, err
	}
	user.CountryCode = r.FormValue(fields.FieldCountryCode)
	return &user, nil
}

func ProfileSignature(r *http.Request) ([]byte, error) {
	file, _, err := r.FormFile(fields.FieldSignature)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	return io.ReadAll(file)
}
