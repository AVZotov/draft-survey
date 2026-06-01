package chi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AVZotov/draft-survey/internal/service"
	"github.com/AVZotov/draft-survey/internal/sse"
	"github.com/AVZotov/draft-survey/web/components"
)

func (h *Handler) Events(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	ch := h.broker.Subscribe()
	defer h.broker.Unsubscribe(ch)

	for {
		select {
		case event := <-ch:
			data := h.renderEvent(r.Context(), event)
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event.Type, data)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (h *Handler) renderEvent(ctx context.Context, event sse.Event) string {
	const op = "Handler.renderEvent"

	var n service.Notification
	if err := json.Unmarshal([]byte(event.Data), &n); err != nil {
		h.logger.Error(op, err)
		return event.Data
	}

	var buf bytes.Buffer
	switch event.Type {
	case sse.EventToast:
		if err := components.Toast(
			components.ToastKind(n.Kind),
			n.Header,
			n.Message,
		).Render(ctx, &buf); err != nil {
			h.logger.Error(op, err)
			return ""
		}
	case sse.EventAlert:
		// TODO: alert component
		return event.Data
	}
	return buf.String()
}
