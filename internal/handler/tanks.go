package handler

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/AVZotov/draft-survey/internal/calculation"
	"github.com/AVZotov/draft-survey/internal/constants"
	"github.com/AVZotov/draft-survey/internal/handler/tadaptor"
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/components"
	"github.com/AVZotov/draft-survey/web/templates/pages"
	"github.com/AVZotov/draft-survey/web/widgets/tanks"
	"github.com/AVZotov/draft-survey/web/widgets/tanks/corrections"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) tanks(c *fiber.Ctx) error {
	p, err := getProps(h, c)
	if err != nil {
		return err
	}
	results := calculation.CalcDraft(p.survey.Drafts[p.draftIndex], p.survey.VesselData)
	tanksLayoutProps := web.TanksLayoutProps(p.user)
	tanksProps := web.TanksPageProps(*p.survey, p.draftIndex, p.trim, p.list, p.trimDir, p.listDir)

	return tadaptor.Render(c, pages.Tanks(tanksLayoutProps, tanksProps, results))
}

func (h *Handler) newTank(c *fiber.Ctx) error {
	tankID := uuid.New().String()
	p, err := getProps(h, c)
	if err != nil {
		return err
	}

	wt := types.Tank{
		ID: tankID,
	}
	parseTank(c, &wt)

	wtType := c.Get(constants.HXTankType)

	switch wtType {
	case constants.BWTank:
		p.survey.Drafts[p.draftIndex].BallastWaterTanks =
			append(p.survey.Drafts[p.draftIndex].BallastWaterTanks, wt)
	case constants.FWTank:
		wt.IsFWTTank = true
		p.survey.Drafts[p.draftIndex].FreshWaterTanks =
			append(p.survey.Drafts[p.draftIndex].FreshWaterTanks, wt)
	default:
		return errors.New("undefined tank type")
	}

	if err = h.surveyRepository.Save(p.survey); err != nil {
		return err
	}

	return tadaptor.Render(c, templ.Join(
		components.TankItem(p.survey.ID, p.draftIndex, wt),
		tanks.AddRowForm(p.survey.ID, p.draftIndex, true)))
}

func (h *Handler) deleteTank(c *fiber.Ctx) error {
	p, err := getProps(h, c)
	if err != nil {
		return err
	}

	wtType := c.Get(constants.HXTankType)
	switch wtType {
	case constants.BWTank:
		p.bwTanks = slices.Delete(p.bwTanks, p.tankIndex, p.tankIndex+1)
		p.survey.Drafts[p.draftIndex].BallastWaterTanks = p.bwTanks
	case constants.FWTank:
		p.bwTanks = slices.Delete(p.bwTanks, p.tankIndex, p.tankIndex+1)
		p.survey.Drafts[p.draftIndex].BallastWaterTanks = p.bwTanks
	default:
		return errors.New("undefined tank type")
	}

	if err := h.surveyRepository.Save(p.survey); err != nil {
		return err
	}

	c.Locals(constants.HXCalcContext, constants.UpdtTanksWeight)
	c.Locals(constants.HXDraftIndex, strconv.Itoa(p.draftIndex))

	return h.calculate(p.survey, c, nil)
}

func (h *Handler) updateBwTank(c *fiber.Ctx) error {
	p, err := getProps(h, c)
	if err != nil {
		return err
	}

	tank := p.bwTanks[p.tankIndex]
	parseTank(c, &tank)
	p.survey.Drafts[p.draftIndex].BallastWaterTanks[p.tankIndex] = tank

	if err := h.surveyRepository.Save(p.survey); err != nil {
		return err
	}

	c.Locals(constants.HXCalcContext, constants.UpdtTanksWeight)
	c.Locals(constants.HXDraftIndex, strconv.Itoa(p.draftIndex))

	return h.calculate(p.survey, c,
		components.TankItem(p.survey.ID, p.draftIndex, tank))
}

func (h *Handler) tanksCorrections(c *fiber.Ctx) error {
	p, err := getProps(h, c)
	if err != nil {
		return err
	}

	tank := p.bwTanks[p.tankIndex]
	tanksProps := web.TanksPageProps(*p.survey, p.draftIndex, p.trim, p.list, p.trimDir, p.listDir)

	if c.Get("HX-Request") != "true" {
		return c.Redirect(fmt.Sprintf("/survey/%s/tanks/%d", p.surveyID, p.draftIndex))
	}
	c.Status(http.StatusOK)
	return tadaptor.Render(c, corrections.ModalForm(tank, tanksProps))
}
