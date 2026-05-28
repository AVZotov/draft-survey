package parse

import (
	"io"
	"net/http"

	"github.com/AVZotov/draft-survey/internal/handler/chi/fields"
	"github.com/AVZotov/draft-survey/internal/types"
)

func Profile(r *http.Request) (*types.User, error) {
	user := &types.User{
		FirstName:  r.FormValue(fields.FieldFirstName),
		LastName:   r.FormValue(fields.FieldLastName),
		Position:   r.FormValue(fields.FieldPosition),
		Email:      r.FormValue(fields.FieldEmail),
		Company:    r.FormValue(fields.FieldCompany),
		License:    r.FormValue(fields.FieldLicense),
		Country:    r.FormValue(fields.FieldCountry),
		EmployeeID: r.FormValue(fields.FieldEmployeeID),
	}
	return user, nil
}

func ProfileSignature(r *http.Request) ([]byte, error) {
	file, _, err := r.FormFile(fields.FieldSignature)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}
