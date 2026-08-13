package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AVZotov/draft-survey/internal/calculation"
	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/templates/pages"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) getDraft(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.getDraft"

	id := chi.URLParam(r, "id")

	data, err := h.services.Survey.GetPageData(id)
	if err != nil {
		h.logger.Error(op, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if data.Outcome != nil && data.Outcome.Redirect != "" {
		return
	}

	sr := calculation.CalcSurvey(*data.Survey)

	lp := web.DraftLayoutProps(data.User, h.appVersion)
	dp := web.DraftsPageProps(*data.Survey)

	if err = pages.DraftReadings(lp, dp, sr).Render(r.Context(), w); err != nil {
		h.logger.Error(op, err)
	}
}

func (h *Handler) updateDraft(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.updateDraft"

	id := chi.URLParam(r, "id")
	indexStr := chi.URLParam(r, "index")

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	survey, err := h.services.Survey.Get(id)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	if index < 0 || index >= len(survey.Drafts) {
		err := fmt.Errorf("draft index %d out of range", index)
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	h.decoder.Draft(r, &survey.Drafts[index], index)

	outcome, err := h.services.Survey.Update(survey)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	h.respond(w, outcome)
}

func (h *Handler) startDraft(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.startDraft"

	id := chi.URLParam(r, "id")
	index, err := strconv.Atoi(chi.URLParam(r, "index"))
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	survey, err := h.services.Survey.Get(id)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	if index < 0 || index >= len(survey.Drafts) {
		err := fmt.Errorf("draft index %d out of range", index)
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	outcome, err := h.services.Draft.Start(survey, index)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	h.respond(w, outcome)
}

func (h *Handler) finishDraft(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.finishDraft"

	id := chi.URLParam(r, "id")
	index, err := strconv.Atoi(chi.URLParam(r, "index"))
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	survey, err := h.services.Survey.Get(id)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	if index < 0 || index >= len(survey.Drafts) {
		err := fmt.Errorf("draft index %d out of range", index)
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	outcome, err := h.services.Draft.Finish(survey, index)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	h.respond(w, outcome)
}

func (h *Handler) addDraft(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.addDraft"

	id := chi.URLParam(r, "id")

	survey, err := h.services.Survey.Get(id)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	outcome, err := h.services.Draft.Add(survey)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	h.respond(w, outcome)
}

func (h *Handler) deleteDraft(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.deleteDraft"

	id := chi.URLParam(r, "id")

	survey, err := h.services.Survey.Get(id)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	outcome, err := h.services.Draft.Delete(survey)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	h.respond(w, outcome)
}
