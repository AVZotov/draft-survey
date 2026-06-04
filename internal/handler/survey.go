package handler

import (
	"net/http"

	"github.com/AVZotov/draft-survey/internal/handler/routes"
	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/templates/pages"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) newSurvey(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.newSurvey"

	survey, err := h.services.Survey.Create()
	if err != nil {
		h.logger.Error(op, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, routes.Survey(survey.ID), http.StatusSeeOther)
}

func (h *Handler) getSurvey(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.getSurvey"

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

	lp := web.SurveyLayoutProps(data.User, h.appVersion)
	sp := web.SurveyPageProps(*data.Survey)

	if err = pages.NewSurvey(lp, sp).Render(r.Context(), w); err != nil {
		h.logger.Error(op, err)
	}
}
