package handler

import (
	"bytes"
	"net/http"
	"strconv"

	"github.com/AVZotov/draft-survey/internal/handler/fields"
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

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit := fields.SurveyListLimit

	surveys, err := h.services.Survey.GetAll()
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	if offset == 0 {
		stats := service.CalcStats(surveys)
		var statsBuf bytes.Buffer
		if err = surveylist.Stats(stats).Render(r.Context(), &statsBuf); err != nil {
			h.logger.Error(op, err)
		} else {
			h.broker.Publish(sse.Event{Type: sse.EventSurveyStats, Data: statsBuf.String()})
		}
	}

	total := len(surveys)
	end := offset + limit
	if end > limit {
		end = total
	}

	h.logger.Debug(op, "offset", offset, "total", len(surveys), "end", end)
	if offset >= total {
		if err := surveylist.Rows(nil, 0, false).Render(r.Context(), w); err != nil {
			h.logger.Error(op, err)
		}
		return
	}
	page := surveys[offset:end]
	hasMore := end < total
	nextOffset := end

	dtos := make([]service.SurveyDTO, len(page))
	for i, s := range surveys {
		h.logger.Debug(op, "converting survey", s.ID)
		dtos[i] = service.ToSurveyDTO(s)
	}

	if err = surveylist.Rows(dtos, nextOffset, hasMore).Render(r.Context(), w); err != nil {
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
