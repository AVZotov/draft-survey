package handler

import (
	"github.com/AVZotov/draft-survey/internal/handler/tadaptor"
	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/templates/pages"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) surveys(c *fiber.Ctx) error {
	user, err := h.userRepository.Get()
	if err != nil {
		return err
	}

	lp := web.SurveyListLayoutProps(user)
	slp, err := getSurveyListProps(h)
	if err != nil {
		return err
	}

	return tadaptor.Render(c, pages.SurveyList(lp, *slp))
}
