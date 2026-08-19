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

func NewSurveyService(
	surveyRepo storage.SurveyRepository, userRepo storage.UserRepository, log logger.Logger,
) SurveyService {
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

	initialDraft := types.NewDraft()
	initialDraft.Surveyor = *user
	initialDraft.Type = types.DraftTypeInitial
	initialDraft.Status = types.DraftStatusPending

	survey := &types.Survey{
		Status:    types.SurveyStatusDraft,
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		Surveyor:  *user,
		Audit:     []types.AuditEvent{},
		VesselData: types.VesselData{
			IsLcfDetectionManual: true,
			// Matches the MMC-Method radio group's own default: form.templ
			// pre-checks "Standard/marine" whenever VesselType isn't
			// river/barge, so a surveyor who never touches that radio should
			// persist the same value the UI already shows as selected —
			// see calculation.vesselTypeOrDefault for the calc-layer half
			// of this fix (covers surveys created before this change).
			VesselType: types.VesselTypeMarine,
		},
		Drafts: []types.Draft{initialDraft},
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

func (s *surveyService) GetPage(limit, offset int) ([]*types.Survey, error) {
	const op = "SurveyService.GetPage"
	
	surveys, err := s.surveyRepo.GetPage(limit, offset)
	if err != nil {
		s.logger.Error(op, err)
		return nil, err
	}
	
	return surveys, nil
}

func (s *surveyService) GetStats() (types.SurveyStats, error) {
	const op = "SurveyService.GetStats"
	
	stats, err := s.surveyRepo.GetStats()
	if err != nil {
		s.logger.Error(op, err)
		return types.SurveyStats{}, err
	}
	
	return stats, nil
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

func (s *surveyService) Delete(id string) (*Outcome, error) {
	const op = "SurveyService.Delete"
	
	if err := s.surveyRepo.Delete(id); err != nil {
		s.logger.Error(op, err)
		return nil, err
	}
	
	return &Outcome{
		Toast: &Notification{
			Kind:    KindSuccess,
			Header:  "Deleted",
			Message: "Survey removed successfully",
		},
	}, nil
}

func (s *surveyService) Search(filter types.SurveyFilter) ([]*types.Survey, error) {
	const op = "SurveyService.Search"
	
	surveys, err := s.surveyRepo.Search(filter)
	if err != nil {
		s.logger.Error(op, err)
		return nil, err
	}
	
	return surveys, nil
}
