package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AVZotov/draft-survey/internal/handler/tadaptor"
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/templates/pages"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) newSurvey(c *fiber.Ctx) error {
	id := uuid.New().String()
	user, err := h.userRepository.Get()
	if err != nil {
		return err
	}
	survey := &types.Survey{
		Status:    types.SurveyStatusDraft,
		ID:        id,
		CreatedAt: time.Now(),
		Surveyor:  *user,
	}

	if err := h.surveyRepository.Save(survey); err != nil {
		return err
	}

	slp := web.SurveyLayoutProps(user)
	spp := web.SurveyPageProps(*survey)

	return tadaptor.Render(c, pages.NewSurvey(slp, spp))
}

func (h *Handler) saveSurvey(c *fiber.Ctx) error {
	p, err := getSurveyProps(h, c)
	if err != nil {
		return err
	}

	parseSurveyPage(c, p.survey)

	if err = h.surveyRepository.Save(p.survey); err != nil {
		return err
	}

	url := fmt.Sprintf("/survey/%s", p.surveyID)
	c.Set("HX-Replace-Url", url)

	return c.SendStatus(http.StatusOK)
}

func (h *Handler) getSurvey(c *fiber.Ctx) error {
	id := c.Params("id")
	survey, err := h.surveyRepository.Get(id)
	if err != nil {
		return err
	}

	user, err := h.userRepository.Get()
	if err != nil {
		return err
	}

	slp := web.SurveyLayoutProps(user)
	spp := web.SurveyPageProps(*survey)

	return tadaptor.Render(c, pages.NewSurvey(slp, spp))
}
