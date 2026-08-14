package service

import (
	"time"

	"github.com/AVZotov/draft-survey/internal/calculation"
	"github.com/AVZotov/draft-survey/internal/handler/routes"
	"github.com/AVZotov/draft-survey/internal/logger"
	"github.com/AVZotov/draft-survey/internal/storage"
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/internal/validation"
)

type draftService struct {
	surveyRepo storage.SurveyRepository
	usersRepo  storage.UserRepository
	logger     logger.Logger
	validator  validation.DraftValidator
}

func NewDraftService(
	surveyRepo storage.SurveyRepository, userRepo storage.UserRepository, log logger.Logger,
	validator validation.DraftValidator,
) DraftService {
	return &draftService{
		surveyRepo: surveyRepo,
		usersRepo:  userRepo,
		logger:     log,
		validator:  validator,
	}
}

func (s *draftService) Add(survey *types.Survey) (*Outcome, error) {
	const op = "DraftService.Add"

	survey.Drafts = append(survey.Drafts, types.Draft{
		Type:            types.DraftTypeFinal,
		Status:          types.DraftStatusPending,
		MTCRows:         make([]types.MTCRow, 2),
		HydrostaticRows: make([]types.HydrostaticRow, 2),
	})

	updateDraftType(survey)

	if err := s.surveyRepo.Save(survey); err != nil {
		s.logger.Error(op, err)
		return nil, err
	}

	return &Outcome{Redirect: routes.Draft(survey.ID)}, nil
}

func (s *draftService) Start(survey *types.Survey, index int) (*Outcome, error) {
	const op = "DraftService.Start"

	user, err := s.usersRepo.Get()
	if err != nil {
		s.logger.Error(op, err)
		return nil, err
	}

	survey.Drafts[index].Status = types.DraftStatusActive
	survey.Drafts[index].StartedAt = time.Now()
	survey.Drafts[index].Surveyor = *user

	if err = s.surveyRepo.Save(survey); err != nil {
		s.logger.Error(op, err)
		return nil, err
	}

	return &Outcome{Redirect: routes.Draft(survey.ID)}, nil
}

func (s *draftService) Finish(survey *types.Survey, index int) (*Outcome, error) {
	const op = "DraftService.Finish"

	survey.Drafts[index].Status = types.DraftStatusComplete
	survey.Drafts[index].FinishedAt = time.Now()

	if err := s.surveyRepo.Save(survey); err != nil {
		s.logger.Error(op, err)
		return nil, err
	}

	return &Outcome{Redirect: routes.Draft(survey.ID)}, nil
}

func (s *draftService) Delete(survey *types.Survey) (*Outcome, error) {
	const op = "DraftService.Delete"

	if len(survey.Drafts) > 2 {
		survey.Drafts = survey.Drafts[:len(survey.Drafts)-1]
	}

	updateDraftType(survey)

	if err := s.surveyRepo.Save(survey); err != nil {
		s.logger.Error(op, err)
		return nil, err
	}

	return &Outcome{
		Redirect: routes.Draft(survey.ID),
		Toast: &Notification{
			Kind:    KindSuccess,
			Header:  "Deleted",
			Message: "Draft removed successfully",
		},
	}, nil
}

func (s *draftService) Update(survey *types.Survey, index int) (*types.SurveyResult, *Outcome, error) {
	const op = "DraftService.Update"

	sr := calculation.CalcSurvey(*survey)

	if index >= 0 && index < len(survey.Drafts) && index < len(sr.DraftTotals) {
		if events := s.validator.Validate(survey.Drafts[index], sr.DraftTotals[index].DraftResult); len(events) > 0 {
			survey.Audit = append(survey.Audit, events...)
		}
	}

	if err := s.surveyRepo.Save(survey); err != nil {
		s.logger.Error(op, err)
		return nil, nil, err
	}

	return &sr, nil, nil
}

// updateDraftType re-derives Type for every draft based on its position in the
// slice. Ported as-is from v0.4.0 — pure domain logic, index-only, no side effects.
func updateDraftType(s *types.Survey) {
	if len(s.Drafts) == 1 {
		return
	}

	if len(s.Drafts) == 2 {
		s.Drafts[1].Type = types.DraftTypeFinal
		return
	}

	for i := len(s.Drafts) - 1; i > 0; i-- {
		if i == len(s.Drafts)-1 {
			s.Drafts[i].Type = types.DraftTypeFinal
			continue
		}
		s.Drafts[i].Type = types.DraftTypeIntermediate
	}
}
