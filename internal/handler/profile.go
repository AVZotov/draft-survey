package handler

import (
	"net/http"
	"strings"

	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/components"
	"github.com/AVZotov/draft-survey/web/templates/pages"
	"github.com/a-h/templ"
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

	user, decerr, err := h.decoder.Profile(r)
	if decerr != nil {
		h.logger.Warn(op, "decode errors", "errors", decerr.Error())
	}
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
	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		if sign, err := h.decoder.ProfileSignature(r); err == nil && len(sign) > 0 {
			signOutcome, err := h.services.User.SaveSignature(sign)
			if err != nil {
				h.logger.Error(op, err)
			} else {
				outcome = signOutcome
			}
		}
	}
	h.respond(w, outcome)
}

func (h *Handler) GetProfileCountrySelect(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.GetProfileCountrySelect"

	user, err := h.services.User.Get()
	if err != nil {
		h.logger.Error(op, err)
		templ.Handler(components.CountrySelect(nil, "")).ServeHTTP(w, r)
		return
	}

	countries, err := h.services.Dictionary.GetCountries()
	if err != nil {
		h.logger.Error(op, err)
		templ.Handler(components.CountrySelect(nil, "")).ServeHTTP(w, r)
		return
	}

	selected := ""
	if user != nil {
		selected = user.CountryCode
	}

	templ.Handler(components.CountrySelect(*countries, selected)).ServeHTTP(w, r)
}
