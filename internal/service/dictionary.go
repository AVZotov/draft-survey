package service

import "github.com/AVZotov/draft-survey/internal/types"

type NoopDictionaryService struct{}

func (s *NoopDictionaryService) GetPorts() (*[]types.Port, error)        { return nil, nil }
func (s *NoopDictionaryService) GetCountries() (*[]types.Country, error) { return nil, nil }
