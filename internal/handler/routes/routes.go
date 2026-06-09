package routes

// Pages
func Home() string            { return "/" }
func Profile() string         { return "/profile" }
func NewSurvey() string       { return "/survey/new" }
func Survey(id string) string { return "/survey/" + id }
func Draft(id string) string  { return "/survey/" + id + "draft" }
func SurveyList() string      { return "/survey-list" }

// SSE
func Events() string { return "/events" }

// API v1 - dictionaries
func APISurveyCargoSelect(id string) string   { return "/api/v1/survey/" + id + "/cargo-select" }
func APISurveyPackingSelect(id string) string { return "/api/v1/survey/" + id + "/packing-select" }

// API v1 — profile
func APIProfile() string              { return "/api/v1/profile" }
func APIProfileCountrySelect() string { return "/api/v1/profile/country-select" }

// API v1 — survey
func APISurvey(id string) string              { return "/api/v1/survey/" + id }
func APISurveyCountrySelect(id string) string { return "/api/v1/survey/" + id + "/country-select" }
func APISurveyFlagSelect(id string) string    { return "/api/v1/survey/" + id + "/flag-select" }

// API v1 - survey-list
func SurveyListRows() string           { return "/survey-list/rows" }
func APIDeleteSurvey(id string) string { return "/api/v1/survey/" + id }
