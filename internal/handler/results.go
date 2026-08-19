package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AVZotov/draft-survey/internal/calculation"
	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/templates/pages"
	"github.com/AVZotov/draft-survey/web/widgets/results"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) getResults(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.getResults"

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

	// CalcSurvey guarantees DraftTotals length == len(survey.Drafts)
	// Direct call is correct here: initial page load is read-only, no save occurs.
	sr := calculation.CalcSurvey(*data.Survey)

	lp := web.ResultsLayoutProps(data.User, h.appVersion)
	rp := web.ResultsPageProps(*data.Survey, sr, 0, len(sr.DraftTotals)-1)

	if err = pages.Results(lp, rp).Render(r.Context(), w); err != nil {
		h.logger.Error(op, err)
	}
}

// getResultsDraft renders a single comparison column (Block) for the
// draft selector's hx-get — used to swap either the "first" or "last"
// column without reloading the page.
func (h *Handler) getResultsDraft(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.getResultsDraft"

	id := chi.URLParam(r, "id")

	index, err := strconv.Atoi(r.URL.Query().Get("index"))
	if err != nil {
		h.logger.Error(op, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	position := r.URL.Query().Get("position")
	if position != "first" && position != "last" {
		h.logger.Error(op, fmt.Errorf("invalid position %q", position))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	isFirst := position == "first"

	data, err := h.services.Survey.GetPageData(id)
	if err != nil {
		h.logger.Error(op, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if index < 0 || index >= len(data.Survey.Drafts) {
		h.logger.Error(op, fmt.Errorf("draft index %d out of range", index))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sr := calculation.CalcSurvey(*data.Survey)

	rp := web.ResultsPageProps(*data.Survey, sr, index, index)

	templ.Handler(results.Block(rp, index, isFirst)).ServeHTTP(w, r)
}
