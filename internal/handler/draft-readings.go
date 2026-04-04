package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AVZotov/draft-survey/internal/handler/tadaptor"
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/templates/pages"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) draft(c *fiber.Ctx) error {
	p, err := getDraftProps(h, c)
	if err != nil {
		return err
	}

	p.survey.Status = types.SurveyStatusInProgress

	if len(p.survey.Drafts) == 0 {
		drafts := []types.Draft{
			{
				Type:   types.DraftTypeInitial,
				Status: types.DraftStatusPending,
			},
		}
		p.survey.Drafts = drafts
	}

	if err = h.surveyRepository.Save(p.survey); err != nil {
		return err
	}

	layoutProps := web.DraftLayoutProps(p.user)
	pageProps := web.DraftsPageProps(*p.survey)

	return tadaptor.Render(c, pages.DraftReadings(layoutProps, pageProps))
}

func (h *Handler) startDraft(c *fiber.Ctx) error {
	id := c.Params("id")
	index, err := strconv.Atoi(c.Params("index"))
	if err != nil {
		return err
	}

	survey, err := h.surveyRepository.Get(id)
	if err != nil {
		return err
	}

	h.parseDraft(c, survey)

	survey.Drafts[index].Status = types.DraftStatusActive
	survey.Drafts[index].StartedAt = time.Now()

	if err = h.surveyRepository.Save(survey); err != nil {
		return err
	}

	c.Set("HX-Redirect", "/survey/"+id+"/draft")
	return c.SendStatus(http.StatusOK)
}

func (h *Handler) finishDraft(c *fiber.Ctx) error {
	id := c.Params("id")
	index, err := strconv.Atoi(c.Params("index"))
	if err != nil {
		return err
	}

	survey, err := h.surveyRepository.Get(id)
	if err != nil {
		return err
	}

	h.parseDraft(c, survey)

	survey.Drafts[index].Status = types.DraftStatusComplete
	survey.Drafts[index].FinishedAt = time.Now()

	if err = h.surveyRepository.Save(survey); err != nil {
		return err
	}

	c.Set("HX-Redirect", "/survey/"+id+"/draft")
	return c.SendStatus(http.StatusOK)
}

func (h *Handler) addIntermediateDraft(c *fiber.Ctx) error {
	id := c.Params("id")
	survey, err := h.surveyRepository.Get(id)
	if err != nil {
		return err
	}

	survey.Drafts = append(survey.Drafts, types.Draft{
		Type:   types.DraftTypeIntermediate,
		Status: types.DraftStatusPending,
	})

	if err = h.surveyRepository.Save(survey); err != nil {
		return err
	}

	c.Set("HX-Redirect", "/survey/"+id+"/draft")
	return c.SendStatus(http.StatusOK)
}

func (h *Handler) addFinalDraft(c *fiber.Ctx) error {
	id := c.Params("id")
	survey, err := h.surveyRepository.Get(id)
	if err != nil {
		return err
	}

	survey.Drafts = append(survey.Drafts, types.Draft{
		Type:   types.DraftTypeFinal,
		Status: types.DraftStatusPending,
	})

	if err = h.surveyRepository.Save(survey); err != nil {
		return err
	}

	c.Set("HX-Redirect", "/survey/"+id+"/draft")
	return c.SendStatus(http.StatusOK)
}

func (h *Handler) saveDrafts(c *fiber.Ctx) error {
	redirect := c.Query("redirect")
	id := c.Params("id")
	survey, err := h.surveyRepository.Get(id)
	if err != nil {
		return err
	}

	h.parseDraft(c, survey)

	if err = h.surveyRepository.Save(survey); err != nil {
		return err
	}

	if redirect != "" {
		c.Set("HX-Redirect", redirect)
	}
	return c.SendStatus(http.StatusOK)
}

func (h *Handler) updateDraft(c *fiber.Ctx) error {
	p, err := getDraftProps(h, c)
	if err != nil {
		return err
	}
	parseDraftBlocks(c, &p.survey.Drafts[p.draftIndex], p.draftIndex)

	if err := h.surveyRepository.Save(p.survey); err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}
