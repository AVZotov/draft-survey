package service

import (
	"github.com/AVZotov/draft-survey/internal/types"
)

// SurveyPageData bundles everything a page handler needs to render a survey
// page in one call: the survey itself, the current user (for layout/signature
// display), and an Outcome if the lookup produced a redirect/notification
// instead of a survey (e.g. survey not found).
type SurveyPageData struct {
	Survey  *types.Survey
	User    *types.User
	Outcome *Outcome
}

// DictionaryService serves read-only reference data (countries, ports, cargo
// types, packing) used to populate select dropdowns across survey pages.
type DictionaryService interface {
	// GetCountries returns every country in the dictionary.
	GetCountries() (*[]types.Country, error)
	// GetPorts returns the ports belonging to the given country code.
	GetPorts(countryCode string) (*[]types.Port, error)
	// GetCargoTypes returns the list of known cargo type names.
	GetCargoTypes() ([]string, error)
	// GetPacking returns the list of known packing type names.
	GetPacking() ([]string, error)
}

// UserService manages the single surveyor profile this offline-first
// application operates as.
type UserService interface {
	// Get returns the current user profile, or nil if none has been set up yet.
	Get() (*types.User, error)
	// Save persists changes to the user profile.
	Save(user *types.User) (*Outcome, error)
	// SaveSignature persists the surveyor's signature image bytes.
	SaveSignature(signature []byte) (*Outcome, error)
	// Delete removes the user profile.
	Delete() error
}

// SurveyService manages survey lifecycle: creation, retrieval, listing,
// search, autosave updates, and deletion.
type SurveyService interface {
	// Create starts a new survey with a single pending initial draft,
	// attributed to the current user profile if one exists.
	Create() (*types.Survey, error)
	// Get returns the survey with the given ID.
	Get(id string) (*types.Survey, error)
	// GetPageData resolves a survey ID into everything a page handler needs
	// to render (survey, user, outcome) in one call.
	GetPageData(id string) (*SurveyPageData, error)
	// GetPage returns a page of surveys for the survey list, most recent first.
	GetPage(limit, offset int) ([]*types.Survey, error)
	// GetStats returns aggregate survey counts for the dashboard.
	GetStats() (types.SurveyStats, error)
	// Search returns surveys matching the given filter.
	Search(filter types.SurveyFilter) ([]*types.Survey, error)
	// Update persists changes to an existing survey (autosave).
	Update(survey *types.Survey) (*Outcome, error)
	// Delete removes the survey with the given ID.
	Delete(id string) (*Outcome, error)
}

// DraftService manages the drafts within a survey: adding, starting,
// finishing, deleting, and recalculating them on autosave.
type DraftService interface {
	// Add appends a new pending final draft to the survey, copying tank
	// structure (without measurements) from the previous draft if any exist,
	// and re-derives every draft's Type based on its new position.
	Add(survey *types.Survey) (*Outcome, error)
	// Start marks the draft at index as active and records the surveyor and
	// start time.
	Start(survey *types.Survey, index int) (*Outcome, error)
	// Finish marks the draft at index as complete and records the finish time.
	Finish(survey *types.Survey, index int) (*Outcome, error)
	// Delete removes the last draft in the survey (at least one intermediate
	// draft must remain) and re-derives every remaining draft's Type.
	Delete(survey *types.Survey) (*Outcome, error)
	// Update recalculates the survey after an autosave, appends any
	// validator-produced audit events, and persists the result. The
	// returned SurveyResult is calculation data only (never persisted) —
	// the handler uses it to render the EventDraftCalc SSE payload, since
	// the service layer never touches templ/HTML.
	Update(survey *types.Survey, index int) (*types.SurveyResult, *Outcome, error)
}

// TankService manages the ballast water and fresh water tanks within a
// single draft: add, update, delete, reorder, and calibration-driven volume
// recalculation.
//
// TankService mirrors DraftService.Update's (*SurveyResult, *Outcome, error)
// signature across every method: every mutation shifts the draft's BW/FW
// totals, and the handler needs SurveyResult to render the EventTankCalc SSE
// payload, same rationale as internal/CLAUDE.md documents for DraftService.Update.
type TankService interface {
	// AddBW appends a new ballast water tank to the given draft.
	AddBW(survey *types.Survey, draftIndex int, tank types.Tank) (*types.SurveyResult, *Outcome, error)
	// AddFW appends a new fresh water tank to the given draft.
	AddFW(survey *types.Survey, draftIndex int, tank types.Tank) (*types.SurveyResult, *Outcome, error)
	// UpdateBW persists changes to an existing ballast water tank (autosave).
	UpdateBW(survey *types.Survey, draftIndex int, tank types.Tank) (*types.SurveyResult, *Outcome, error)
	// UpdateFW persists changes to an existing fresh water tank (autosave).
	UpdateFW(survey *types.Survey, draftIndex int, tank types.Tank) (*types.SurveyResult, *Outcome, error)
	// DeleteBW removes a ballast water tank from the given draft.
	DeleteBW(survey *types.Survey, draftIndex int, tankID string) (*types.SurveyResult, *Outcome, error)
	// DeleteFW removes a fresh water tank from the given draft.
	DeleteFW(survey *types.Survey, draftIndex int, tankID string) (*types.SurveyResult, *Outcome, error)
	// ApplyDensity sets the dockwater density used for BW tank calibration
	// lookups on the given draft.
	ApplyDensity(survey *types.Survey, draftIndex int, density float64) (*types.SurveyResult, *Outcome, error)
	// MoveBWTankUp swaps a ballast water tank with the one before it in
	// display order.
	MoveBWTankUp(survey *types.Survey, draftIndex int, tankID string) (*types.SurveyResult, *Outcome, error)
	// MoveBWTankDown swaps a ballast water tank with the one after it in
	// display order.
	MoveBWTankDown(survey *types.Survey, draftIndex int, tankID string) (*types.SurveyResult, *Outcome, error)
	// MoveFWTankUp swaps a fresh water tank with the one before it in
	// display order.
	MoveFWTankUp(survey *types.Survey, draftIndex int, tankID string) (*types.SurveyResult, *Outcome, error)
	// MoveFWTankDown swaps a fresh water tank with the one after it in
	// display order.
	MoveFWTankDown(survey *types.Survey, draftIndex int, tankID string) (*types.SurveyResult, *Outcome, error)
	// CopyFromPrevious copies tank structure (without measurements) from the
	// preceding draft into the given draft.
	CopyFromPrevious(survey *types.Survey, draftIndex int) (*types.SurveyResult, *Outcome, error)
}

// Services aggregates every domain service the handler layer depends on.
type Services struct {
	User       UserService
	Survey     SurveyService
	Draft      DraftService
	Tank       TankService
	Dictionary DictionaryService
}
