package handler

import (
	"io/fs"
	"net/http"

	"github.com/AVZotov/draft-survey"
	"github.com/go-chi/chi/v5"
)

func SetupRoutesChi(r chi.Router, h *HandlerChi) error {
	staticFS, err := fs.Sub(draft_survey.StaticFiles, "web/static")
	if err != nil {
		return err
	}

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	//TODO: This comment as WIP route per route with handlers migrations to Chi
	r.Get("/", nil)
	r.Group(
		func(r chi.Router) {
			//TODO: Add deny all direct requests except of main page
			r.Get("/profile", nil)
			r.Get("/dictionary/sea-options", nil)
			r.Get("/dictionary/ports", nil)
			r.Get("/dictionary/countries", nil)
			r.Get("/survey/new", nil)
			r.Get("/survey/{id}", nil)
			r.Get("/survey/{id}/draft", nil)
			r.Get("/survey/{id}/tanks/{draftIndex}", nil)
			r.Get("/survey/{id}/tanks/{draftIndex}/bw-tank/{tankID}/corrections", nil)
			r.Get("/survey/{id}/results", nil)
			r.Get("/survey/{id}/results/validate", nil)
			r.Get("/survey-list", nil)
			r.Get("/survey-list/rows", nil)
			r.Get("/survey-list/stats", nil)

			r.Route(
				"/api/v1", func(r chi.Router) {
					r.Post("/profile", nil)
					r.Put("/survey/{id}", nil)
					r.Post("/survey/{id}", nil)
					r.Post("/survey/{id}/draft/{index}/start", nil)
					r.Post("/survey/{id}/draft/{index}/finish", nil)
					r.Post("/survey/{id}/draft", nil)
					r.Delete("/survey/{id}/draft", nil)
					r.Put("/survey/{id}/draft/{draftIndex}", nil)
					r.Put("/survey/{id}/tanks/{draftIndex}", nil)
					//TODO: Change bw-tank to tank or w-tank as more relevant route
					r.Post("/survey/{id}/tanks/{draftIndex}/bw-tank", nil)
					r.Delete("/survey/{id}/tanks/{draftIndex}/bw-tank/{tankID}", nil)
					r.Put("/survey/{id}/tanks/{draftIndex}/bw-tank/{tankID}", nil)
					r.Delete("/survey-list/rows/{id}", nil)
				},
			)
		},
	)
	return nil
}
