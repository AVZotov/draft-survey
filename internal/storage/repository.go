package storage

import "github.com/AVZotov/draft-survey/internal/types"

type SurveyRepository interface {
	Save(survey *types.Survey) error
	Get(id string) (*types.Survey, error)
	GetPage(limit, offset int) ([]*types.Survey, error)
	GetStats() (types.SurveyStats, error)
	Delete(id string) error
}

type UserRepository interface {
	Save(user *types.User) error
	SaveSignature(data []byte) error
	Get() (*types.User, error)
	Delete() error
}

type DictionariesRepository interface {
	GetCountries() (*[]types.Country, error)
	GetPorts(countryCode string) (*[]types.Port, error)
	GetCargoTypes() ([]string, error)
	GetPacking() ([]string, error)
}
