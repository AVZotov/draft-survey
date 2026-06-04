package handler

import (
	"net/http"
	"net/url"
	"time"

	"github.com/AVZotov/draft-survey/internal/handler/fields"
	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/templates/pages"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	const op = "Hadler.home"
	user, err := h.services.User.Get()
	if err != nil {
		h.logger.Error(op, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		value := url.QueryEscape(`{"header":"Profile Required","message":"Please set up your profile to continue"}`)
		cookie := &http.Cookie{
			Name:     fields.CookieFlashAlert,
			Value:    value,
			Path:     "/",
			Expires:  time.Now().Add(5 * time.Second),
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	props := web.DashboardLayoutProps(user, h.appVersion)
	err = pages.Dashboard(props).Render(r.Context(), w)
	if err != nil {
		h.logger.Error(op, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
