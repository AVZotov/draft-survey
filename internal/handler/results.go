package handler

import (
	"encoding/json"

	"github.com/AVZotov/draft-survey/internal/calculation"
	"github.com/AVZotov/draft-survey/internal/handler/tadaptor"
	"github.com/AVZotov/draft-survey/internal/types/dto"
	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/components"
	"github.com/AVZotov/draft-survey/web/templates/pages"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) results(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userRepository.Get()
	if err != nil {
		return err
	}

	survey, err := h.surveyRepository.Get(id)
	if err != nil {
		return err
	}

	surveyResults := calculation.CalcSurvey(*survey)

	alpineData := dto.AlpineDTO{
		Survey:        *survey,
		SurveyResults: surveyResults,
	}

	alpineJSON, err := json.Marshal(alpineData)
	if err != nil {
		return err
	}

	props := components.ResultProps{
		Survey:    *survey,
		Alpine:    string(alpineJSON),
		Lastindex: len(surveyResults.DraftTotals) - 1,
	}

	lp := web.ResultsLayoutProps(user)

	return tadaptor.Render(c, pages.Results(lp, props))
}
