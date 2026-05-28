package chi

import (
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/templates/pages"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) profile(w http.ResponseWriter, r *http.Request) {
	const op = "profile"

	user, err := h.userRepository.Get()
	if err != nil {
		h.logger.Error(op, err)
		//TODO: Implement error handling. Avoid handling no user error
	}

	props := web.ProfileLayoutProps(user, h.appVersion)
	err = pages.Dashboard(props).Render(r.Context(), w)
	if err != nil {
		h.logger.Error(op, err)
		//TODO: Implement error handling.
	}
}

func (h *Handler) createProfile(c *fiber.Ctx) error {
	user := &types.User{
		FirstName:  c.FormValue("first_name"),
		LastName:   c.FormValue("last_name"),
		Position:   c.FormValue("position"),
		Email:      c.FormValue("email"),
		Company:    c.FormValue("company"),
		License:    c.FormValue("license_no"),
		Country:    c.FormValue("country"),
		EmployeeID: c.FormValue("employee_id"),
	}
	if err := h.userRepository.Save(user); err != nil {
		return err
	}

	if file, err := c.FormFile("signature"); err == nil {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer func(src multipart.File) {
			cerr := src.Close()
			if cerr != nil && err == nil {
				err = cerr
			}
		}(src)
		data, err := io.ReadAll(src)
		if err != nil {
			return err
		}
		ext := filepath.Ext(file.Filename)
		if err := h.userRepository.SaveSignature(data, ext); err != nil {
			return err
		}
	}

	c.Set("HX-Redirect", "/")
	return c.SendStatus(http.StatusOK)
}
