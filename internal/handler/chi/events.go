package chi

import (
	"fmt"
	"net/http"
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
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event.Type, event.Data)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
