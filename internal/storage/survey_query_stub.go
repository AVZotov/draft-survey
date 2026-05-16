package storage

import "github.com/AVZotov/draft-survey/internal/types"

var _ SurveyQueryRepository = (*SurveyQueryStub)(nil)

type SurveyQueryStub struct {
}

func NewSurveyQueryStub() *SurveyQueryStub {
	return &SurveyQueryStub{}
}

func (s *SurveyQueryStub) Search(filter SurveyFilter) ([]*types.Survey, error) {
	panic("implement me")
}
