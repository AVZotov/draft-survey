package handler

import (
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

func (h *Handler) calculate(survey *types.Survey, c *fiber.Ctx, main templ.Component, oobComponents ...templ.Component) error {
	hxCtx := (c.Locals(constants.HXCalcContext)).(string)
	if hxCtx == "" {
		return ErrUndefinedHXContext
	}

	switch hxCtx {
	case constants.UpdtTanksWeight:
		component, err := updtTanksWeight(survey, c, main, oobComponents...)
		if err != nil {
			return err
		}
		c.Status(http.StatusOK)
		return tadaptor.Render(c, component)
	default:
		return ErrUndefinedHXContext
	}
}

func updtTanksWeight(survey *types.Survey, c *fiber.Ctx, main templ.Component, oobComponets ...templ.Component) (templ.Component, error) {
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
	var components []templ.Component

	if main != nil {
		components = append(components, main)
	}

	switch wtType {
	case constants.BWTank:
		components = append(components,
			tanks.TableFormHeader(constants.BwTableHeaderID, bwTitle, bwTWeight, draft.Type, true),
			tanks.ABTotals(bwTWeight, fwTWeight, true),
		)

	case constants.FWTank:
		components = append(components,
			tanks.TableFormHeader(constants.FwTableHeaderID, fwTitle, fwTWeight, draft.Type, true),
			tanks.ABTotals(bwTWeight, fwTWeight, true),
		)
	default:
		return nil, ErrUndefinedTankType
	}

	components = append(components, oobComponets...)
	return templ.Join(components...), nil
}
