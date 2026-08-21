package parse

import (
	"io"
	"net/http"

	"github.com/AVZotov/draft-survey/internal/handler/fields"
	"github.com/AVZotov/draft-survey/internal/types"
)

func (d *Decoder) Profile(r *http.Request, existing *types.User) *types.User {
	user := existing
	if user == nil {
		user = &types.User{}
	}

	if r.PostForm.Has(fields.FieldFirstName) {
		user.FirstName = d.str(r, fields.FieldFirstName)
	}
	if r.PostForm.Has(fields.FieldLastName) {
		user.LastName = d.str(r, fields.FieldLastName)
	}
	if r.PostForm.Has(fields.FieldPosition) {
		user.Position = d.str(r, fields.FieldPosition)
	}
	if r.PostForm.Has(fields.FieldEmail) {
		user.Email = d.str(r, fields.FieldEmail)
	}
	if r.PostForm.Has(fields.FieldCompany) {
		user.Company = d.str(r, fields.FieldCompany)
	}
	if r.PostForm.Has(fields.FieldLicense) {
		user.License = d.str(r, fields.FieldLicense)
	}
	if r.PostForm.Has(fields.FieldCountryCode) {
		user.CountryCode = d.str(r, fields.FieldCountryCode)
	}
	if r.PostForm.Has(fields.FieldEmployeeID) {
		user.EmployeeID = d.str(r, fields.FieldEmployeeID)
	}

	return user
}

func (d *Decoder) ProfileSignature(r *http.Request) ([]byte, error) {
	file, _, err := r.FormFile(fields.FieldSignature)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}
