package chi

import (
	"net/http"

	"github.com/AVZotov/draft-survey/internal/handler/chi/parse"
	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/templates/pages"
)

func (h *Handler) profile(w http.ResponseWriter, r *http.Request) {
	const op = "profile"
	user, err := h.services.User.Get()
	if err != nil {
		h.logger.Error(op, err)
		//TODO: Implement error handling. Avoid handling no user error
	}

	props := web.ProfileLayoutProps(user, h.appVersion)
	err = pages.Profile(props).Render(r.Context(), w)
	if err != nil {
		h.logger.Error(op, err)
		//TODO: Implement error handling.
	}
}

func (h *Handler) createProfile(w http.ResponseWriter, r *http.Request) {
	const op = "createProfile"
	user, err := parse.Profile(r)
	if err = h.services.User.Save(user); err != nil {
		h.logger.Error(op, err)
		//TODO: Implement error handling.
	}

	sign, err := parse.ProfileSignature(r)
	if err != nil {
		h.logger.Error(op, err)
		//TODO: Implement error handling.
	}

	if err = h.services.User.SaveSignature(sign); err != nil {
		h.logger.Error(op, err)
		//TODO: Implement error handling. Avoid handle no signature error
	}
}
