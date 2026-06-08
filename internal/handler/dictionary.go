package handler

import (
	"net/http"

	"github.com/AVZotov/draft-survey/internal/handler/fields"
	"github.com/AVZotov/draft-survey/web/components"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetSurveyCargoSelect(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.GetSurveyCargoSelect"

	id := chi.URLParam(r, "id")
	survey, err := h.services.Survey.Get(id)
	if err != nil {
		h.logger.Error(op, err)
		templ.Handler(components.StringSelect(nil, "", fields.FieldCargo)).ServeHTTP(w, r)
		return
	}

	types, err := h.services.Dictionary.GetCargoTypes()
	if err != nil {
		h.logger.Error(op, err)
		templ.Handler(components.StringSelect(nil, "", fields.FieldCargo)).ServeHTTP(w, r)
		return
	}

	selected := ""
	if survey != nil {
		selected = survey.CargoOperation.Cargo
	}

	templ.Handler(components.StringSelect(types, selected, fields.FieldCargo)).ServeHTTP(w, r)
}

func (h *Handler) GetSurveyPackingSelect(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.GetSurveyPackingSelect"

	id := chi.URLParam(r, "id")
	survey, err := h.services.Survey.Get(id)
	if err != nil {
		h.logger.Error(op, err)
		templ.Handler(components.StringSelect(nil, "", fields.FieldPacking)).ServeHTTP(w, r)
		return
	}

	items, err := h.services.Dictionary.GetPacking()
	if err != nil {
		h.logger.Error(op, err)
		templ.Handler(components.StringSelect(nil, "", fields.FieldPacking)).ServeHTTP(w, r)
		return
	}

	selected := ""
	if survey != nil {
		selected = survey.CargoOperation.Packing
	}

	templ.Handler(components.StringSelect(items, selected, fields.FieldPacking)).ServeHTTP(w, r)
}
