package chi

import (
	"net/http"

	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/templates/pages"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	const op = "home"
	user, err := h.services.User.Get()
	if err != nil {
		h.logger.Error(op, err)
		//TODO: Implement error handling. Avoid handling no user error
	}
	if user == nil {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	props := web.DashboardLayoutProps(user, h.appVersion)
	err = pages.Dashboard(props).Render(r.Context(), w)
	if err != nil {
		h.logger.Error(op, err)
		//TODO: Implement error handling.
	}
}
