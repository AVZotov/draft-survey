package routes

import "fmt"

// Pages

func Home() string             { return "/" }
func Profile() string          { return "/profile" }
func Survey(id string) string  { return "/survey/" + id }
func Draft(id string) string   { return "/survey/" + id + "/draft" }
func Results(id string) string { return "/survey/" + id + "/results" }
func SurveyList() string       { return "/survey-list" }

func Tanks(id string, draftIndex string) string {
	return fmt.Sprintf("/survey/%s/tanks/%s", id, draftIndex)
}
func TankCorrections(id string, draftIndex string, tankID string) string {
	return fmt.Sprintf("/survey/%s/tanks/%s/tank/%s/corrections", id, draftIndex, tankID)
}

// API v1

func APICreateSurvey() string          { return "/api/v1/survey" }
func APISurvey(id string) string       { return "/api/v1/survey/" + id }
func APIDeleteSurvey(id string) string { return "/api/v1/survey/" + id }
func APIProfile() string               { return "/api/v1/profile" }
func APIDraft(id string) string        { return "/api/v1/survey/" + id + "/draft" }
func APIDraftUpdate(id string, draftIndex string) string {
	return fmt.Sprintf("/api/v1/survey/%s/draft/%s", id, draftIndex)
}
func APIDraftStatus(id string, draftIndex string) string {
	return fmt.Sprintf("/api/v1/survey/%s/draft/%s/status", id, draftIndex)
}

func APITanks(id string, draftIndex string) string {
	return fmt.Sprintf("/api/v1/survey/%s/tanks/%s", id, draftIndex)
}
func APITank(id string, draftIndex string, tankID string) string {
	return fmt.Sprintf("/api/v1/survey/%s/tanks/%s/tank/%s", id, draftIndex, tankID)
}

// HTMX fragments

func SurveyListRows() string  { return "/survey-list/rows" }
func SurveyListStats() string { return "/survey-list/stats" }
func ResultsValidate(id string) string {
	return "/survey/" + id + "/results/validate"
}

// Dictionaries

func DictionarySeaOptions() string { return "/dictionary/sea-options" }
func DictionaryPorts() string      { return "/dictionary/ports" }
func DictionaryCountries() string  { return "/dictionary/countries" }

// SSE Routes

func Events() string { return "/events" }
