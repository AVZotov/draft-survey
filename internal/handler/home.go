package handler

import (
	"errors"
	"net/http"
	"os"

	"github.com/AVZotov/draft-survey/internal/handler/tadaptor"
	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/templates/pages"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) home(c *fiber.Ctx) error {
	user, err := h.userRepository.Get()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	props := web.DashboardLayoutProps(user, h.appVersion)
	c.Status(http.StatusOK)
	return tadaptor.Render(c, pages.Dashboard(props))
}
