package handler

import (
	"bytes"
	"net/http"

	"github.com/AVZotov/draft-survey/internal/service"
	"github.com/AVZotov/draft-survey/internal/sse"
	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/templates/pages"
	surveylist "github.com/AVZotov/draft-survey/web/widgets/survey_list"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) getSurveyList(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.getSurveyList"

	user, err := h.services.User.Get()
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	lp := web.SurveyListLayoutProps(user, h.appVersion)

	if err = pages.SurveyList(lp).Render(r.Context(), w); err != nil {
		h.logger.Error(op, err)
	}
}

func (h *Handler) getSurveyListRows(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.getSurveyListRows"

	surveys, err := h.services.Survey.GetAll()
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	stats := service.CalcStats(surveys)

	var statsBuf bytes.Buffer
	if err = surveylist.Stats(stats).Render(r.Context(), &statsBuf); err != nil {
		h.logger.Error(op, err)
	} else {
		h.broker.Publish(sse.Event{Type: sse.EventSurveyStats, Data: statsBuf.String()})
	}

	dtos := make([]service.SurveyDTO, len(surveys))
	for i, s := range surveys {
		dtos[i] = service.ToSurveyDTO(s)
	}

	if err = surveylist.Rows(dtos).Render(r.Context(), w); err != nil {
		h.logger.Error(op, err)
	}
}

func (h *Handler) deleteSurvey(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.deleteSurvey"

	id := chi.URLParam(r, "id")

	outcome, err := h.services.Survey.Delete(id)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	surveys, err := h.services.Survey.GetAll()
	if err != nil {
		h.logger.Error(op, err)
	} else {
		stats := service.CalcStats(surveys)
		var buf bytes.Buffer
		if err = surveylist.Stats(stats).Render(r.Context(), &buf); err != nil {
			h.logger.Error(op, err)
		} else {
			h.broker.Publish(sse.Event{Type: sse.EventSurveyStats, Data: buf.String()})
		}
	}

	h.respond(w, outcome)
}
