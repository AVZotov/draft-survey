package service

import (
	"github.com/AVZotov/draft-survey/internal/storage"
	"github.com/AVZotov/draft-survey/internal/types"
)

type dictionaryService struct {
	repo storage.DictionariesRepository
}

func NewDictionaryService(repo storage.DictionariesRepository) DictionaryService {
	return &dictionaryService{repo: repo}
}

func (s *dictionaryService) GetCountries() (*[]types.Country, error) {
	return s.repo.GetCountries()
}

func (s *dictionaryService) GetPorts(countryCode string) (*[]types.Port, error) {
	return s.repo.GetPorts(countryCode)
}
