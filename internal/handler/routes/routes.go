package routes

// Pages
func Home() string            { return "/" }
func Profile() string         { return "/profile" }
func NewSurvey() string       { return "/survey/new" }
func Survey(id string) string { return "/survey/" + id }

// SSE
func Events() string { return "/events" }

// API v1 — profile
func APIProfile() string              { return "/api/v1/profile" }
func APIProfileCountrySelect() string { return "/api/v1/profile/country-select" }

// API v1 — survey
func APISurvey(id string) string { return "/api/v1/survey/" + id }
