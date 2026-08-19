package handler

import (
	"io/fs"
	"net/http"

	draftsurvey "github.com/AVZotov/draft-survey"
	"github.com/AVZotov/draft-survey/internal/handler/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutesChi(r chi.Router, h *Handler) error {
	r.Use(middleware.RequestID)
	r.Use(h.requestLogger)

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
	r.Get(routes.Draft("{id}"), h.getDraft)
	r.Get(routes.TanksPattern("{id}"), h.getTanks)
	r.Get(routes.Results("{id}"), h.getResults)
	r.Get(routes.SurveyList(), h.getSurveyList)

	// SSE
	r.Get(routes.Events(), h.Events)

	//API v1 - dictionaries
	r.Get(routes.APISurveyCargoSelect("{id}"), h.GetSurveyCargoSelect)
	r.Get(routes.APISurveyPackingSelect("{id}"), h.GetSurveyPackingSelect)
	r.Get(routes.APICargoTypesSelect(), h.GetCargoTypesSelect)
	r.Get(routes.APIDictionarySeaOptions(), h.GetSeaOptions)

	// API v1 — profile
	r.Post(routes.APIProfile(), h.createProfile)
	r.Get(routes.APIProfileCountrySelect(), h.GetProfileCountrySelect)

	// API v1 — survey
	r.Post(routes.APISurvey("{id}"), h.saveSurvey)
	r.Get(routes.APISurveyCountrySelect("{id}"), h.GetSurveyCountrySelect)
	r.Get(routes.APISurveyFlagSelect("{id}"), h.GetSurveyFlagSelect)

	// API v1 — draft
	r.Put(routes.APIDraftUpdatePattern("{id}"), h.updateDraft)
	r.Post(routes.APIDraftStartPattern("{id}"), h.startDraft)
	r.Post(routes.APIDraftFinishPattern("{id}"), h.finishDraft)
	r.Post(routes.APIDraftAdd("{id}"), h.addDraft)
	r.Delete(routes.APIDraftDelete("{id}"), h.deleteDraft)

	// API v1 — tanks
	r.Get(routes.APITankBWCorrectionsPattern("{id}"), h.getTankCorrections)
	r.Post(routes.APITankBWPattern("{id}"), h.addBWTank)
	r.Put(routes.APITankBWIDPattern("{id}"), h.updateBWTank)
	r.Delete(routes.APITankBWIDPattern("{id}"), h.deleteBWTank)
	r.Post(routes.APITankFWPattern("{id}"), h.addFWTank)
	r.Put(routes.APITankFWIDPattern("{id}"), h.updateFWTank)
	r.Delete(routes.APITankFWIDPattern("{id}"), h.deleteFWTank)
	r.Put(routes.APITankApplyDensityPattern("{id}"), h.applyDensity)
	r.Put(routes.APITankBWMoveUpPattern("{id}"), h.moveBWTankUp)
	r.Put(routes.APITankBWMoveDownPattern("{id}"), h.moveBWTankDown)
	r.Put(routes.APITankFWMoveUpPattern("{id}"), h.moveFWTankUp)
	r.Put(routes.APITankFWMoveDownPattern("{id}"), h.moveFWTankDown)

	// API v1 — results
	r.Get(routes.APIResultsDraft("{id}"), h.getResultsDraft)

	// API v1 - survey-list
	r.Get(routes.SurveyListRows(), h.getSurveyListRows)
	r.Delete(routes.APIDeleteSurvey("{id}"), h.deleteSurvey)
	r.Get(routes.SurveySearch(), h.searchSurveys)

	return nil
}
