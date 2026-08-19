package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AVZotov/draft-survey/internal/calculation"
	"github.com/AVZotov/draft-survey/internal/handler/routes"
	"github.com/AVZotov/draft-survey/internal/sse"
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/templates/pages"
	"github.com/AVZotov/draft-survey/web/widgets/drafts"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) getDraft(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.getDraft"

	id := chi.URLParam(r, "id")

	data, err := h.services.Survey.GetPageData(id)
	if err != nil {
		h.logger.Error(op, err)
		h.respondPageError(w, r, routes.SurveyList())
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

	sr, outcome, err := h.services.Draft.Update(survey, index)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	h.publishDraftCalc(r.Context(), *survey, *sr)

	h.respond(w, outcome)
}

// publishDraftCalc renders the calc-panel/hydrostatics-result/totals
// fragments for every draft (SurveyResult.DraftTotals covers all drafts,
// since a single draft's change can shift cumulative cargo/constant figures
// on later drafts) and publishes them as one EventDraftCalc SSE event.
// Mirrors how getSurveyListRows/deleteSurvey render+publish EventSurveyStats
// directly in the handler, bypassing Outcome (which is reserved for
// Toast/Alert notifications).
func (h *Handler) publishDraftCalc(ctx context.Context, survey types.Survey, sr types.SurveyResult) {
	const op = "Handler.publishDraftCalc"

	var buf bytes.Buffer
	for i, dt := range sr.DraftTotals {
		if err := drafts.CalcPanel(survey, i, dt.DraftResult, true).Render(ctx, &buf); err != nil {
			h.logger.Error(op, err)
			return
		}
		if err := drafts.MMCRow(survey, i, dt.DraftResult, true).Render(ctx, &buf); err != nil {
			h.logger.Error(op, err)
			return
		}
		if err := drafts.LBM(dt.DraftResult, i, true).Render(ctx, &buf); err != nil {
			h.logger.Error(op, err)
			return
		}
		if err := drafts.DeltaMtc(survey, i, dt.DraftResult, true).Render(ctx, &buf); err != nil {
			h.logger.Error(op, err)
			return
		}
		if err := drafts.Totals(survey, i, sr).Render(ctx, &buf); err != nil {
			h.logger.Error(op, err)
			return
		}
	}

	h.broker.Publish(sse.Event{Type: sse.EventDraftCalc, Data: buf.String()})
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
