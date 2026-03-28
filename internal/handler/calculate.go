package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/AVZotov/draft-survey/internal/calculation"
	"github.com/AVZotov/draft-survey/internal/constants"
	"github.com/AVZotov/draft-survey/internal/format"
	"github.com/AVZotov/draft-survey/internal/handler/tadaptor"
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/web/widgets/tanks"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) calculate(c *fiber.Ctx) error {
	hxCtx := (c.Locals(constants.HXCalcContext)).(string)
	if hxCtx == "" {
		return templ.Error{Err: ErrEmptyField}
	}

	survey, err := h.surveyRepository.Get((c.Locals(constants.HXSurveyID)).(string))
	if err != nil {
		return err
	}

	switch hxCtx {
	case constants.UpdtTanksWeight:
		component, err := updtTanksWeight(survey, c)
		if err != nil {
			return err
		}
		c.Status(http.StatusOK)
		return tadaptor.Render(c, component)
	default:
		return errors.New("undefined HX context")
	}
}

func updtTanksWeight(survey *types.Survey, c *fiber.Ctx) (templ.Component, error) {
	draftIndexStr := (c.Locals(constants.HXDraftIndex)).(string)
	draftIndex, err := strconv.Atoi(draftIndexStr)
	if err != nil {
		return nil, err
	}
	draft := survey.Drafts[draftIndex]
	results := calculation.CalcDraft(draft, survey.VesselData)
	bwTitle := "Ballast Water Tanks"
	fwTitle := "Fresh Water Tanks"
	bwTWeight := format.WeightFormatted(results.TotalBwTanksWeight)
	fwTWeight := format.WeightFormatted(results.TotalFwTanksWeight)

	wtType := c.Get(constants.HXTankType)

	switch wtType {
	case constants.BWTank:
		return templ.Join(tanks.TableFormHeader(constants.BwTableHeaderID, bwTitle, bwTWeight, draft.Type, true)), nil
	case constants.FWTank:
		return templ.Join(tanks.TableFormHeader(constants.FwTableHeaderID, fwTitle, fwTWeight, draft.Type, true)), nil
	default:
		return nil, errors.New("undefined tank type")
	}
}
