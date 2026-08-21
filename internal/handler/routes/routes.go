package routes

import "strconv"

// Pages
func Home() string            { return "/" }
func Profile() string         { return "/profile" }
func NewSurvey() string       { return "/survey/new" }
func Survey(id string) string { return "/survey/" + id }
func Draft(id string) string  { return "/survey/" + id + "/draft" }
func SurveyList() string      { return "/survey-list" }

// SSE
func Events() string { return "/events" }

// API v1 - dictionaries
func APISurveyCargoSelect(id string) string   { return "/api/v1/survey/" + id + "/cargo-select" }
func APISurveyPackingSelect(id string) string { return "/api/v1/survey/" + id + "/packing-select" }

// API v1 - dictionariess No preselected items for filter
func APICargoTypesSelect() string     { return "/api/v1/dictionary/cargo-select" }
func APIDictionarySeaOptions() string { return "/api/v1/dictionary/sea-options" }

// API v1 — profile
func APIProfile() string              { return "/api/v1/profile" }
func APIProfileCountrySelect() string { return "/api/v1/profile/country-select" }

// API v1 — survey
func APISurvey(id string) string              { return "/api/v1/survey/" + id }
func APISurveyCountrySelect(id string) string { return "/api/v1/survey/" + id + "/country-select" }
func APISurveyFlagSelect(id string) string    { return "/api/v1/survey/" + id + "/flag-select" }

// API v1 — draft
func APIDraftUpdate(id string, index int) string {
	return "/api/v1/survey/" + id + "/draft/" + strconv.Itoa(index)
}
func APIDraftUpdatePattern(id string) string { return "/api/v1/survey/" + id + "/draft/{index}" }
func APIDraftStart(id string, index int) string {
	return "/api/v1/survey/" + id + "/draft/" + strconv.Itoa(index) + "/start"
}
func APIDraftStartPattern(id string) string { return "/api/v1/survey/" + id + "/draft/{index}/start" }
func APIDraftFinish(id string, index int) string {
	return "/api/v1/survey/" + id + "/draft/" + strconv.Itoa(index) + "/finish"
}
func APIDraftFinishPattern(id string) string { return "/api/v1/survey/" + id + "/draft/{index}/finish" }
func APIDraftAdd(id string) string           { return "/api/v1/survey/" + id + "/draft" }
func APIDraftDelete(id string) string        { return "/api/v1/survey/" + id + "/draft" }

// Tanks page
func Tanks(id string, draftIndex int) string {
	return "/survey/" + id + "/tanks/" + strconv.Itoa(draftIndex)
}
func TanksPattern(id string) string { return "/survey/" + id + "/tanks/{draftIndex}" }

// API v1 — tanks: ballast water
func APITankBWAdd(id string, draftIndex int) string {
	return "/api/v1/survey/" + id + "/tanks/" + strconv.Itoa(draftIndex) + "/bw"
}
func APITankBWPattern(id string) string {
	return "/api/v1/survey/" + id + "/tanks/{draftIndex}/bw"
}
func APITankBWUpdate(id string, draftIndex int, tankID string) string {
	return "/api/v1/survey/" + id + "/tanks/" + strconv.Itoa(draftIndex) + "/bw/" + tankID
}
func APITankBWDelete(id string, draftIndex int, tankID string) string {
	return APITankBWUpdate(id, draftIndex, tankID)
}
func APITankBWIDPattern(id string) string {
	return "/api/v1/survey/" + id + "/tanks/{draftIndex}/bw/{tankID}"
}
func APITankBWCorrections(id string, draftIndex int, tankID string) string {
	return "/survey/" + id + "/tanks/" + strconv.Itoa(draftIndex) + "/bw/" + tankID + "/corrections"
}
func APITankBWCorrectionsPattern(id string) string {
	return "/survey/" + id + "/tanks/{draftIndex}/bw/{tankID}/corrections"
}
func APITankBWMoveUp(id string, draftIndex int, tankID string) string {
	return "/api/v1/survey/" + id + "/tanks/" + strconv.Itoa(draftIndex) + "/bw/" + tankID + "/move-up"
}
func APITankBWMoveUpPattern(id string) string {
	return "/api/v1/survey/" + id + "/tanks/{draftIndex}/bw/{tankID}/move-up"
}
func APITankBWMoveDown(id string, draftIndex int, tankID string) string {
	return "/api/v1/survey/" + id + "/tanks/" + strconv.Itoa(draftIndex) + "/bw/" + tankID + "/move-down"
}
func APITankBWMoveDownPattern(id string) string {
	return "/api/v1/survey/" + id + "/tanks/{draftIndex}/bw/{tankID}/move-down"
}

// API v1 — tanks: fresh water
func APITankFWAdd(id string, draftIndex int) string {
	return "/api/v1/survey/" + id + "/tanks/" + strconv.Itoa(draftIndex) + "/fw"
}
func APITankFWPattern(id string) string {
	return "/api/v1/survey/" + id + "/tanks/{draftIndex}/fw"
}
func APITankFWUpdate(id string, draftIndex int, tankID string) string {
	return "/api/v1/survey/" + id + "/tanks/" + strconv.Itoa(draftIndex) + "/fw/" + tankID
}
func APITankFWDelete(id string, draftIndex int, tankID string) string {
	return APITankFWUpdate(id, draftIndex, tankID)
}
func APITankFWIDPattern(id string) string {
	return "/api/v1/survey/" + id + "/tanks/{draftIndex}/fw/{tankID}"
}
func APITankFWMoveUp(id string, draftIndex int, tankID string) string {
	return "/api/v1/survey/" + id + "/tanks/" + strconv.Itoa(draftIndex) + "/fw/" + tankID + "/move-up"
}
func APITankFWMoveUpPattern(id string) string {
	return "/api/v1/survey/" + id + "/tanks/{draftIndex}/fw/{tankID}/move-up"
}
func APITankFWMoveDown(id string, draftIndex int, tankID string) string {
	return "/api/v1/survey/" + id + "/tanks/" + strconv.Itoa(draftIndex) + "/fw/" + tankID + "/move-down"
}
func APITankFWMoveDownPattern(id string) string {
	return "/api/v1/survey/" + id + "/tanks/{draftIndex}/fw/{tankID}/move-down"
}

// API v1 — tanks: bulk density apply
func APITankApplyDensity(id string, draftIndex int) string {
	return "/api/v1/survey/" + id + "/tanks/" + strconv.Itoa(draftIndex) + "/density"
}
func APITankApplyDensityPattern(id string) string {
	return "/api/v1/survey/" + id + "/tanks/{draftIndex}/density"
}

// API v1 — tanks: copy from previous draft
func APITankCopyFromPrevious(id string, draftIndex int) string {
	return "/api/v1/survey/" + id + "/tanks/" + strconv.Itoa(draftIndex) + "/copy-from-previous"
}
func APITankCopyFromPreviousPattern(id string) string {
	return "/api/v1/survey/" + id + "/tanks/{draftIndex}/copy-from-previous"
}

// Results page
func Results(id string) string { return "/survey/" + id + "/results" }

// API v1 — results
func APIResultsDraft(id string) string { return "/api/v1/survey/" + id + "/results/draft" }

// API v1 - survey-list
func SurveyListRows() string           { return "/survey-list/rows" }
func APIDeleteSurvey(id string) string { return "/api/v1/survey/" + id }
func SurveySearch() string             { return "/survey-list/search" }

func SurveyListPage(offset, limit int) string {
	return SurveyListRows() + "?offset=" + strconv.Itoa(offset) + "&limit=" + strconv.Itoa(limit)
}
