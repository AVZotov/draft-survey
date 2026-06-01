package chi

import (
	"net/http"

	"github.com/AVZotov/draft-survey/internal/handler/chi/parse"
	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/templates/pages"
)

func (h *Handler) profile(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.profile"

	user, err := h.services.User.Get()
	if err != nil {
		h.logger.Error(op, err)
	}

	props := web.ProfileLayoutProps(user, h.appVersion)
	if err = pages.Profile(props).Render(r.Context(), w); err != nil {
		h.logger.Error(op, err)
	}
}

func (h *Handler) createProfile(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.createProfile"

	user, err := parse.Profile(r)
	if err != nil {
		h.logger.Error(op, err)
		h.respond(w, nil)
		return
	}

	outcome, err := h.services.User.Save(user)
	if err != nil {
		h.logger.Error(op, err)
		h.respond(w, nil)
		return
	}

	if sign, err := parse.ProfileSignature(r); err == nil && len(sign) > 0 {
		signOutcome, err := h.services.User.SaveSignature(sign)
		if err != nil {
			h.logger.Error(op, err)
		} else {
			outcome = signOutcome
		}
	}

	h.respond(w, outcome)
}
