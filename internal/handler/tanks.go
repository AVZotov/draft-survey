package handler

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/AVZotov/draft-survey/internal/calculation"
	"github.com/AVZotov/draft-survey/internal/constants"
	"github.com/AVZotov/draft-survey/internal/format"
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

type props struct {
	surveyID   string
	draftIndex int
	survey     *types.Survey
	user       *types.User
	tankID     string
	tankIndex  int
	bwTanks    []types.Tank
	trim       *float64
	trimDir    string
	list       *float64
	listDir    string
}

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
	h.parseTank(c, &wt)

	wtType := c.Get(constants.HXTank)

	switch wtType {
	case constants.BWTank:
		p.survey.Drafts[p.draftIndex].BallastWaterTanks =
			append(p.survey.Drafts[p.draftIndex].BallastWaterTanks, wt)
	case constants.FWTank:
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

func (h *Handler) deleteBwTank(c *fiber.Ctx) error {
	p, err := getProps(h, c)
	if err != nil {
		return err
	}

	p.bwTanks = slices.Delete(p.bwTanks, p.tankIndex, p.tankIndex+1)

	p.survey.Drafts[p.draftIndex].BallastWaterTanks = p.bwTanks
	if err := h.surveyRepository.Save(p.survey); err != nil {
		return err
	}

	tanksProps := web.TanksPageProps(*p.survey, p.draftIndex, p.trim, p.list, p.trimDir, p.listDir)

	c.Status(http.StatusOK)
	return tadaptor.Render(c, tanks.BwTableHeaderForm(tanksProps, true))
}

func (h *Handler) updateBwTank(c *fiber.Ctx) error {
	p, err := getProps(h, c)
	if err != nil {
		return err
	}

	tank := p.bwTanks[p.tankIndex]
	h.parseTank(c, &tank)

	p.survey.Drafts[p.draftIndex].BallastWaterTanks[p.tankIndex] = tank

	if err := h.surveyRepository.Save(p.survey); err != nil {
		return err
	}

	tanksProps := web.TanksPageProps(*p.survey, p.draftIndex, p.trim, p.list, p.trimDir, p.listDir)

	c.Status(http.StatusOK)
	return tadaptor.Render(c, templ.Join(
		components.TankItem(p.survey.ID, p.draftIndex, tank),
		tanks.BwTableHeaderForm(tanksProps, true)))
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

func getProps(h *Handler, c *fiber.Ctx) (*props, error) {
	id := c.Params("id")
	survey, err := h.surveyRepository.Get(id)
	if err != nil {
		return nil, err
	}

	draftIndexStr := c.Params("draftIndex")
	draftIndex, err := strconv.Atoi(draftIndexStr)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepository.Get()
	if err != nil {
		return nil, err
	}

	var trueTrim, listDegrees *float64
	draft := survey.Drafts[draftIndex]
	if draft.Marks.FwdPort.Value != nil && draft.Marks.FwdPort.Method != "" &&
		draft.Marks.AftPort.Value != nil && draft.Marks.AftPort.Method != "" &&
		draft.Marks.MidPort.Value != nil && draft.Marks.MidPort.Method != "" &&
		draft.DistancePPFwd != nil && draft.DistancePPMid != nil && draft.DistancePPAft != nil &&
		survey.VesselData.LBP > 0 {
		results := calculation.CalcDraft(draft, survey.VesselData)
		trueTrim = &results.TrueTrim
		listDegrees = &results.ListDegrees
	}
	var trimDir, listDir string
	if trueTrim != nil && listDegrees != nil {
		trimDir = format.TrimDirection(*trueTrim)
		listDir = format.ListDirection(*listDegrees)
	}

	tankID := c.Params("tankID")

	if tankID == "" {
		return &props{
			surveyID:   id,
			draftIndex: draftIndex,
			survey:     survey,
			user:       user,
			trim:       trueTrim,
			trimDir:    trimDir,
			list:       listDegrees,
			listDir:    listDir,
		}, nil
	}

	bwTanks := survey.Drafts[draftIndex].BallastWaterTanks
	tankIndex := slices.IndexFunc(bwTanks, func(tank types.Tank) bool {
		return tank.ID == tankID
	})
	if tankIndex == -1 {
		return nil, errors.New("tank not found")
	}

	return &props{
		surveyID:   id,
		draftIndex: draftIndex,
		survey:     survey,
		user:       user,
		tankID:     tankID,
		tankIndex:  tankIndex,
		bwTanks:    bwTanks,
		trim:       trueTrim,
		trimDir:    trimDir,
		list:       listDegrees,
		listDir:    listDir,
	}, nil
}
