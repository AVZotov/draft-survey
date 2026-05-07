package handler

import (
	"net/http"

	draft_survey "github.com/AVZotov/draft-survey"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func SetupRoutes(app *fiber.App, h *Handler) {
	app.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(draft_survey.StaticFiles),
		PathPrefix: "web/static",
		Browse:     false,
	}))

	app.Get("/", h.home)
	app.Get("/profile", h.profile)
	app.Get("/dictionary/sea-options", h.seaOptions)
	app.Get("/dictionary/ports", h.ports)
	app.Get("/dictionary/countries", h.countries)
	app.Get("/survey/new", h.newSurvey)
	app.Get("/survey/:id", h.getSurvey)
	app.Get("/survey/:id/draft", h.draft)
	app.Get("/survey/:id/tanks/:draftIndex", h.tanks)
	app.Get("/survey/:id/tanks/:draftIndex/bw-tank/:tankID/corrections", h.tanksCorrections)
	app.Get("/survey/:id/results", h.results)
	app.Get("/survey/:id/results/validate", h.validate)
	app.Get("/survey-list", h.surveys)
	app.Get("/survey-list/rows", h.surveyRows)
	app.Get("/survey-list/stats", h.surveyStats)

	api := app.Group("/api/v1")
	api.Post("/profile", h.createProfile)
	api.Put("/survey/:id", h.saveSurvey)
	api.Post("/survey/:id", h.saveSurveyAndNavigate)
	api.Post("/survey/:id/draft/:index/start", h.startDraft)
	api.Post("/survey/:id/draft/:index/finish", h.finishDraft)
	api.Post("/survey/:id/draft/add-intermediate", h.addIntermediateDraft)
	api.Post("/survey/:id/draft/add-final", h.addFinalDraft)
	api.Put("/survey/:id/draft/:draftIndex", h.updateDraft)
	api.Put("/survey/:id/tanks/:draftIndex", h.updateTanks)
	api.Post("/survey/:id/tanks/:draftIndex/bw-tank", h.newTank)
	api.Delete("/survey/:id/tanks/:draftIndex/bw-tank/:tankID", h.deleteTank)
	api.Put("/survey/:id/tanks/:draftIndex/bw-tank/:tankID", h.updateTank)
	api.Delete("/survey-list/rows/:id", h.deleteSurveyRow)
}
