package handler

import (
	"net/http"

	"github.com/AVZotov/draft-survey/internal/handler/fields"
	"github.com/AVZotov/draft-survey/internal/handler/routes"
	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/components"
	"github.com/AVZotov/draft-survey/web/templates/pages"
	"github.com/a-h/templ"
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

func (h *Handler) saveSurvey(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.saveSurvey"

	id := chi.URLParam(r, "id")

	if err := r.ParseForm(); err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	h.logger.Debug(op, "form values", r.PostForm)

	existing, err := h.services.Survey.Get(id)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	survey := h.decoder.Survey(r, existing)
	h.logger.Debug(op, "vessel name after parse", survey.VesselData.Name)

	outcome, err := h.services.Survey.Update(survey)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	h.respond(w, outcome)
}

func (h *Handler) GetSurveyCountrySelect(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.GetSurveyCountrySelect"

	id := chi.URLParam(r, "id")
	survey, err := h.services.Survey.Get(id)
	if err != nil {
		h.logger.Error(op, err)
		templ.Handler(components.CountrySelect(nil, "", fields.FieldCountryCode)).ServeHTTP(w, r)
		return
	}

	countries, err := h.services.Dictionary.GetCountries()
	if err != nil {
		h.logger.Error(op, err)
		templ.Handler(components.CountrySelect(nil, "", fields.FieldCountryCode)).ServeHTTP(w, r)
		return
	}

	selected := ""
	if survey != nil {
		selected = survey.Country.CountryCode
	}

	templ.Handler(components.CountrySelect(*countries, selected, fields.FieldCountryCode)).ServeHTTP(w, r)
}

func (h *Handler) GetSurveyFlagSelect(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.GetSurveyFlagSelect"

	id := chi.URLParam(r, "id")
	survey, err := h.services.Survey.Get(id)
	if err != nil {
		h.logger.Error(op, err)
		templ.Handler(components.CountrySelect(nil, "", fields.FieldFlagCountryCode)).ServeHTTP(w, r)
		return
	}

	countries, err := h.services.Dictionary.GetCountries()
	if err != nil {
		h.logger.Error(op, err)
		templ.Handler(components.CountrySelect(nil, "", fields.FieldFlagCountryCode)).ServeHTTP(w, r)
		return
	}

	selected := ""
	if survey != nil {
		selected = survey.VesselData.Flag.CountryCode
	}

	templ.Handler(components.CountrySelect(*countries, selected, fields.FieldFlagCountryCode)).ServeHTTP(w, r)
}
