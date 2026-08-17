package service

import (
	"github.com/AVZotov/draft-survey/internal/types"
)

type SurveyPageData struct {
	Survey  *types.Survey
	User    *types.User
	Outcome *Outcome
}

type DictionaryService interface {
	GetCountries() (*[]types.Country, error)
	GetPorts(countryCode string) (*[]types.Port, error)
	GetCargoTypes() ([]string, error)
	GetPacking() ([]string, error)
}

type UserService interface {
	Get() (*types.User, error)
	Save(user *types.User) (*Outcome, error)
	SaveSignature(signature []byte) (*Outcome, error)
	Delete() error
}

type SurveyService interface {
	Create() (*types.Survey, error)
	Get(id string) (*types.Survey, error)
	GetPageData(id string) (*SurveyPageData, error)
	GetPage(limit, offset int) ([]*types.Survey, error)
	GetStats() (types.SurveyStats, error)
	Search(filter types.SurveyFilter) ([]*types.Survey, error)
	Update(survey *types.Survey) (*Outcome, error)
	Delete(id string) (*Outcome, error)
}

type DraftService interface {
	Add(survey *types.Survey) (*Outcome, error)
	Start(survey *types.Survey, index int) (*Outcome, error)
	Finish(survey *types.Survey, index int) (*Outcome, error)
	Delete(survey *types.Survey) (*Outcome, error)
	// Update recalculates the survey after an autosave, appends any
	// validator-produced audit events, and persists the result. The
	// returned SurveyResult is calculation data only (never persisted) —
	// the handler uses it to render the EventDraftCalc SSE payload, since
	// the service layer never touches templ/HTML.
	Update(survey *types.Survey, index int) (*types.SurveyResult, *Outcome, error)
}

// TankService mirrors DraftService.Update's (*SurveyResult, *Outcome, error)
// signature across every method: every mutation shifts the draft's BW/FW
// totals, and the handler needs SurveyResult to render the EventTankCalc SSE
// payload, same rationale as internal/CLAUDE.md documents for DraftService.Update.
type TankService interface {
	AddBW(survey *types.Survey, draftIndex int, tank types.Tank) (*types.SurveyResult, *Outcome, error)
	AddFW(survey *types.Survey, draftIndex int, tank types.Tank) (*types.SurveyResult, *Outcome, error)
	UpdateBW(survey *types.Survey, draftIndex int, tank types.Tank) (*types.SurveyResult, *Outcome, error)
	UpdateFW(survey *types.Survey, draftIndex int, tank types.Tank) (*types.SurveyResult, *Outcome, error)
	DeleteBW(survey *types.Survey, draftIndex int, tankID string) (*types.SurveyResult, *Outcome, error)
	DeleteFW(survey *types.Survey, draftIndex int, tankID string) (*types.SurveyResult, *Outcome, error)
	ApplyDensity(survey *types.Survey, draftIndex int, density float64) (*types.SurveyResult, *Outcome, error)
}

type Services struct {
	User       UserService
	Survey     SurveyService
	Draft      DraftService
	Tank       TankService
	Dictionary DictionaryService
}
