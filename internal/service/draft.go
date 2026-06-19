package service

import "github.com/AVZotov/draft-survey/internal/types"

type NoopDraftService struct{}

func (s *NoopDraftService) Add(survey *types.Survey) error               { return nil }
func (s *NoopDraftService) Start(survey *types.Survey, index int) error  { return nil }
func (s *NoopDraftService) Finish(survey *types.Survey, index int) error { return nil }
func (s *NoopDraftService) Delete(survey *types.Survey) error            { return nil }
func (s *NoopDraftService) CalcResults(survey *types.Survey) (*types.SurveyResult, error) {
	return nil, nil
}
