package service

import (
	"time"

	"github.com/AVZotov/draft-survey/internal/types"
)

type DictionaryService interface {
	GetPorts() (*[]types.Port, error)
	GetCountries() (*[]types.Country, error)
}

type UserService interface {
	Get() (*types.User, error)
	Save(user *types.User) (*Outcome, error)
	SaveSignature(signature []byte) (*Outcome, error)
	Delete() error
}

type SurveyFilter struct {
	Query string
	From  time.Time
	To    time.Time
}

type SurveyService interface {
	Create(user types.User) (*types.Survey, error)
	Get(id string) (*types.Survey, error)
	Update(survey *types.Survey) error
	Delete(id string) error
	Search(filter SurveyFilter) ([]*types.Survey, error)
}

type DraftService interface {
	Add(survey *types.Survey) error
	Start(survey *types.Survey, index int) error
	Finish(survey *types.Survey, index int) error
	Delete(survey *types.Survey) error
	Update(survey *types.Survey, index int) (*types.SurveyResult, error)
}

type Services struct {
	User       UserService
	Survey     SurveyService
	Draft      DraftService
	Dictionary DictionaryService
}
