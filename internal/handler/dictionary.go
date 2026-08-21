package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AVZotov/draft-survey/internal/handler/fields"
	"github.com/AVZotov/draft-survey/internal/types"
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

func (h *Handler) GetCargoTypesSelect(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.GetCargoTypesSelect"

	types, err := h.services.Dictionary.GetCargoTypes()
	if err != nil {
		h.logger.Error(op, err)
		templ.Handler(components.StringSelect(nil, "", fields.FieldCargo)).ServeHTTP(w, r)
		return
	}

	templ.Handler(components.StringSelect(types, "", fields.FieldCargo)).ServeHTTP(w, r)
}

func (h *Handler) GetSeaOptions(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.GetSeaOptions"

	seaType := r.URL.Query().Get("type")

	var scp components.SeaConditionProps
	if indexStr := r.URL.Query().Get("index"); indexStr != "" {
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			h.logger.Error(op, err)
			h.respondError(w, http.StatusBadRequest, err)
			return
		}
		scp.DraftIndex = &index
	}

	switch types.SeaConditionType(seaType) {
	case types.SeaConditionTypeWave:
		scp.SeaCondition.Type = types.SeaConditionTypeWave
	case types.SeaConditionTypeIce:
		scp.SeaCondition.Type = types.SeaConditionTypeIce
	default:
		h.logger.Error(op, fmt.Errorf("invalid sea type %q", seaType))
		h.respondError(w, http.StatusBadRequest, fmt.Errorf("invalid sea type %q", seaType))
		return
	}

	templ.Handler(components.SeaSelect(scp)).ServeHTTP(w, r)
}
