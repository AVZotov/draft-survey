package chi

import (
	"io/fs"
	"net/http"

	draftsurvey "github.com/AVZotov/draft-survey"
	"github.com/AVZotov/draft-survey/internal/handler/chi/routes"
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
	r.Get(routes.Survey("{id}"), nil)
	r.Get(routes.Draft("{id}"), nil)
	r.Get(routes.Results("{id}"), nil)
	r.Get(routes.SurveyList(), nil)
	r.Get(routes.Tanks("{id}", "{draftIndex}"), nil)
	r.Get(routes.TankCorrections("{id}", "{draftIndex}", "{tankID}"), nil)

	// Dictionaries
	r.Get(routes.DictionarySeaOptions(), nil)
	r.Get(routes.DictionaryPorts(), nil)
	r.Get(routes.DictionaryCountries(), nil)

	// HTMX fragments
	r.Get(routes.SurveyListRows(), nil)
	r.Get(routes.SurveyListStats(), nil)
	r.Get(routes.ResultsValidate("{id}"), nil)

	// API v1
	r.Route(
		"/api/v1", func(r chi.Router) {
			r.Post(routes.APICreateSurvey(), nil)
			r.Put(routes.APISurvey("{id}"), nil)
			r.Delete(routes.APIDeleteSurvey("{id}"), nil)

			r.Post(routes.APIDraft("{id}"), nil)
			r.Delete(routes.APIDraft("{id}"), nil)
			r.Put(routes.APIDraftUpdate("{id}", "{draftIndex}"), nil)
			r.Post(routes.APIDraftStatus("{id}", "{draftIndex}"), nil)

			r.Post(routes.APIProfile(), nil)

			r.Put(routes.APITanks("{id}", "{draftIndex}"), nil)
			r.Post(routes.APITank("{id}", "{draftIndex}", "{tankID}"), nil)
			r.Delete(routes.APITank("{id}", "{draftIndex}", "{tankID}"), nil)
			r.Put(routes.APITank("{id}", "{draftIndex}", "{tankID}"), nil)
		},
	)

	return nil
}
