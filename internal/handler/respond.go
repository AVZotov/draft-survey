package handler

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/AVZotov/draft-survey/internal/handler/fields"
	"github.com/AVZotov/draft-survey/internal/service"
	"github.com/AVZotov/draft-survey/internal/sse"
)

// respond publishes Outcome.Toast/Outcome.Alert based on whether a redirect
// accompanies them. HX-Redirect makes htmx perform a full page navigation, which
// tears down the current SSE connection before a broker.Publish on this request
// could ever reach a live listener — the same redirect-then-notify race the flash
// cookie in home.go already solves for Alert. When Redirect is set, both
// notification kinds are staged as flash cookies instead, and the new page's
// /events connection publishes them on connect.
func (h *Handler) respond(w http.ResponseWriter, outcome *service.Outcome) {
	if outcome == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	if outcome.Redirect != "" {
		if outcome.Toast != nil {
			h.setFlashCookie(w, fields.CookieFlashToast, outcome.Toast)
		}
		if outcome.Alert != nil {
			h.setFlashCookie(w, fields.CookieFlashAlert, outcome.Alert)
		}
		w.Header().Set("HX-Redirect", outcome.Redirect)
		w.WriteHeader(http.StatusOK)
		return
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

func (h *Handler) setFlashCookie(w http.ResponseWriter, name string, n *service.Notification) {
	data, err := json.Marshal(n)
	if err != nil {
		h.logger.Error("respond: marshal flash cookie", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(string(data)),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   2,
	})
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
