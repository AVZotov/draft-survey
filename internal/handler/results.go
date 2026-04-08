package handler

import (
	"github.com/AVZotov/draft-survey/internal/calculation"
	"github.com/AVZotov/draft-survey/internal/handler/tadaptor"
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/web"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) surveyResults(c *fiber.Ctx) error {
	id := c.Params("id")
	survey, err := h.surveyRepository.Get(id)
	if err != nil {
		return err
	}

	user, err := h.userRepository.Get()
	if err != nil {
		return err
	}

	results := make([]types.DraftResult, len(survey.Drafts))
	for i, draft := range survey.Drafts {
		if draft.Status == types.DraftStatusComplete ||
			draft.Status == types.DraftStatusActive {
			results[i] = calculation.CalcDraft(draft, survey.VesselData)
		}
	}
	props := web.ResultsPageProps(user, survey, &results)
	return tadaptor.Render(c, web.Results(props))
}
