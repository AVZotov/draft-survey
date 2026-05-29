package service

import "github.com/AVZotov/draft-survey/internal/types"

type NoopSurveyService struct{}

func (s *NoopSurveyService) Create(user types.User) (*types.Survey, error)       { return nil, nil }
func (s *NoopSurveyService) Get(id string) (*types.Survey, error)                { return nil, nil }
func (s *NoopSurveyService) Update(survey *types.Survey) error                   { return nil }
func (s *NoopSurveyService) Delete(id string) error                              { return nil }
func (s *NoopSurveyService) Search(filter SurveyFilter) ([]*types.Survey, error) { return nil, nil }
