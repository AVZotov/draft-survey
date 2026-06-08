package service

import (
	"time"

	"github.com/AVZotov/draft-survey/internal/logger"
	"github.com/AVZotov/draft-survey/internal/storage"
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/google/uuid"
)

type surveyService struct {
	surveyRepo storage.SurveyRepository
	usersRepo  storage.UserRepository
	logger     logger.Logger
}

func NewSurveyService(surveyRepo storage.SurveyRepository, userRepo storage.UserRepository, log logger.Logger) SurveyService {
	return &surveyService{
		surveyRepo: surveyRepo,
		usersRepo:  userRepo,
		logger:     log,
	}
}

func (s *surveyService) Create() (*types.Survey, error) {
	const op = "SurveyService.Create"

	user, err := s.usersRepo.Get()
	if err != nil {
		s.logger.Error(op, err)
		return nil, err
	}

	survey := &types.Survey{
		Status:    types.SurveyStatusDraft,
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		Surveyor:  *user,
		VesselData: types.VesselData{
			IsLcfDetectionManual: true,
		},
	}

	if err = s.surveyRepo.Save(survey); err != nil {
		s.logger.Error(op, err)
		return nil, err
	}

	return survey, nil
}

func (s *surveyService) Get(id string) (*types.Survey, error) {
	const op = "SurveyService.Get"

	survey, err := s.surveyRepo.Get(id)
	if err != nil {
		s.logger.Error(op, err)
		return nil, err
	}

	return survey, nil
}

func (s *surveyService) GetPageData(id string) (*SurveyPageData, error) {
	const op = "SurveyService.GetPageData"

	survey, err := s.surveyRepo.Get(id)
	if err != nil {
		s.logger.Error(op, err)
		return nil, err
	}

	user, err := s.usersRepo.Get()
	if err != nil {
		s.logger.Error(op, err)
		return nil, err
	}

	return &SurveyPageData{
		Survey: survey,
		User:   user,
	}, nil
}

func (s *surveyService) Update(survey *types.Survey) (*Outcome, error) {
	const op = "SurveyService.Update"

	if err := s.surveyRepo.Save(survey); err != nil {
		s.logger.Error(op, err)
		return nil, err
	}

	return nil, nil
}

func (s *surveyService) Delete(id string) error {
	const op = "SurveyService.Delete"

	if err := s.surveyRepo.Delete(id); err != nil {
		s.logger.Error(op, err)
		return err
	}

	return nil
}
