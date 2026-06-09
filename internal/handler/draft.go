package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) getDraft(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.getDraft"
	id := chi.URLParam(r, "id")
	h.logger.Debug(op, "draft id: ", id)
	w.WriteHeader(http.StatusOK)
}
