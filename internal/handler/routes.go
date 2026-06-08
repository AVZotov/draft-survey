package handler

import (
	"io/fs"
	"net/http"

	draftsurvey "github.com/AVZotov/draft-survey"
	"github.com/AVZotov/draft-survey/internal/handler/routes"
	"github.com/go-chi/chi/v5"
)

func SetupRoutesChi(r chi.Router, h *Handler) error {
	staticFS, err := fs.Sub(draftsurvey.StaticFiles, "web/static")
	if err != nil {
		return err
	}

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Pages
	r.Get(routes.Home(), h.home)
	r.Get(routes.Profile(), h.profile)
	r.Get(routes.NewSurvey(), h.newSurvey)
	r.Get(routes.Survey("{id}"), h.getSurvey)

	// SSE
	r.Get(routes.Events(), h.Events)

	// API v1 — profile
	r.Post(routes.APIProfile(), h.createProfile)
	r.Get(routes.APIProfileCountrySelect(), h.GetProfileCountrySelect)

	// API v1 — survey
	r.Post(routes.APISurvey("{id}"), h.saveSurvey)
	r.Get(routes.APISurveyCountrySelect("{id}"), h.GetSurveyCountrySelect)
	r.Get(routes.APISurveyFlagSelect("{id}"), h.GetSurveyFlagSelect)

	return nil
}
