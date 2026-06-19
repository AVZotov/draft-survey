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
	Add(survey *types.Survey) error
	Start(survey *types.Survey, index int) error
	Finish(survey *types.Survey, index int) error
	Delete(survey *types.Survey) error
	CalcResults(survey *types.Survey) (*types.SurveyResult, error)
}

type Services struct {
	User       UserService
	Survey     SurveyService
	Draft      DraftService
	Dictionary DictionaryService
}
