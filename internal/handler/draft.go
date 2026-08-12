package handler

import (
	"net/http"

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
