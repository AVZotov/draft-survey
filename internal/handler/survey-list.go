package handler

import (
	"net/http"
	"slices"

	"github.com/AVZotov/draft-survey/internal/constants"
	"github.com/AVZotov/draft-survey/internal/handler/tadaptor"
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/templates/pages"
	surveylist "github.com/AVZotov/draft-survey/web/widgets/survey-list"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) surveys(c *fiber.Ctx) error {
	user, err := h.userRepository.Get()
	if err != nil {
		return err
	}
	surveys, err := h.surveyRepository.GetAll()
	if err != nil {
		return err
	}
	slices.SortFunc(surveys, func(a, b *types.Survey) int {
		return b.CreatedAt.Compare(a.CreatedAt)
	})
	props := getSurveyStats(surveys)
	lp := web.SurveyListLayoutProps(user)

	return tadaptor.Render(c, pages.SurveyList(lp, props))
}

func (h *Handler) surveyRows(c *fiber.Ctx) error {
	q := c.Query("q")
	offset := c.QueryInt("offset", 0)
	limit := c.QueryInt("limit", constants.SurveyListLimit)
	from := c.Query("from")
	to := c.Query("to")

	surveys, err := h.surveyRepository.GetAll()
	if err != nil {
		return err
	}

	slices.SortFunc(surveys, func(a, b *types.Survey) int {
		return b.CreatedAt.Compare(a.CreatedAt)
	})

	surveys = filterSurveys(surveys, q, from, to)
	rows := getSurveyRows(surveys, offset, limit)

	nextOffset := offset + limit
	hasMore := nextOffset < len(surveys)

	return tadaptor.Render(c, surveylist.Rows(rows, nextOffset, hasMore))
}

func (h *Handler) surveyStats(c *fiber.Ctx) error {
	surveys, err := h.surveyRepository.GetAll()
	if err != nil {
		return err
	}
	stats := getSurveyStats(surveys)

	return tadaptor.Render(c, surveylist.Stats(stats))
}

func (h *Handler) deleteSurveyRow(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.surveyRepository.Delete(id); err != nil {
		return err
	}

	c.Set("HX-Trigger", "survey-deleted")

	return c.SendStatus(http.StatusOK)
}
