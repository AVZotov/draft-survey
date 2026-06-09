package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AVZotov/draft-survey/internal/service"
	"github.com/AVZotov/draft-survey/internal/sse"
)

func (h *Handler) respond(w http.ResponseWriter, outcome *service.Outcome) {
	if outcome == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	if outcome.Redirect != "" {
		w.Header().Set("HX-Redirect", outcome.Redirect)
	}

	if outcome.Toast != nil {
		data, err := json.Marshal(outcome.Toast)
		if err != nil {
			h.logger.Error("respond: marshal toast", err)
		} else {
			h.broker.Publish(sse.Event{Type: sse.EventToast, Data: string(data)})
		}
	}

	if outcome.Alert != nil {
		data, err := json.Marshal(outcome.Alert)
		if err != nil {
			h.logger.Error("respond: marshal alert", err)
		} else {
			h.broker.Publish(sse.Event{Type: sse.EventAlert, Data: string(data)})
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) respondError(w http.ResponseWriter, status int, err error) {
	message := "Something went wrong. Please try again."
	if status == http.StatusBadRequest && err != nil {
		message = err.Error()
	}

	h.broker.Publish(sse.Event{
		Type: sse.EventAlert,
		Data: mustJSON(service.Notification{
			Kind:    service.KindError,
			Header:  "Error",
			Message: message,
		}),
	})
	w.WriteHeader(status)
}

func mustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return `{"header":"Error","message":"Something went wrong"}`
	}
	return string(b)
}
