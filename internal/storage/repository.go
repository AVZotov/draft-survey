package storage

import (
	"github.com/AVZotov/draft-survey/internal/types"
)

type SurveyRepository interface {
	Save(survey *types.Survey) error
	Get(id string) (*types.Survey, error)
	GetAll() ([]*types.Survey, error)
	Delete(id string) error
}

type UserRepository interface {
	Save(user *types.User) error
	SaveSignature(data []byte, ext string) error
	Get() (*types.User, error)
	Delete() error
}

type DictionariesRepository interface {
	GetPorts() (*[]types.Port, error)
	GetCountries() (*[]types.Country, error)
}
