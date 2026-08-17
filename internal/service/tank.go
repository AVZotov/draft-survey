package service

import (
	"errors"
	"slices"

	"github.com/AVZotov/draft-survey/internal/calculation"
	"github.com/AVZotov/draft-survey/internal/logger"
	"github.com/AVZotov/draft-survey/internal/storage"
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/google/uuid"
)

var ErrTankNotFound = errors.New("tank not found")

type tankService struct {
	surveyRepo storage.SurveyRepository
	logger     logger.Logger
}

func NewTankService(surveyRepo storage.SurveyRepository, log logger.Logger) TankService {
	return &tankService{surveyRepo: surveyRepo, logger: log}
}

func (s *tankService) AddBW(survey *types.Survey, draftIndex int, tank types.Tank) (*types.SurveyResult, *Outcome, error) {
	const op = "TankService.AddBW"

	if tank.ID == "" {
		tank.ID = uuid.New().String()
	}
	s.recalcVolume(survey, draftIndex, &tank)

	survey.Drafts[draftIndex].BallastWaterTanks = append(survey.Drafts[draftIndex].BallastWaterTanks, tank)

	sr, err := s.saveAndCalc(op, survey)
	if err != nil {
		return nil, nil, err
	}
	return sr, nil, nil
}

func (s *tankService) AddFW(survey *types.Survey, draftIndex int, tank types.Tank) (*types.SurveyResult, *Outcome, error) {
	const op = "TankService.AddFW"

	if tank.ID == "" {
		tank.ID = uuid.New().String()
	}
	tank.IsFWTTank = true

	survey.Drafts[draftIndex].FreshWaterTanks = append(survey.Drafts[draftIndex].FreshWaterTanks, tank)

	sr, err := s.saveAndCalc(op, survey)
	if err != nil {
		return nil, nil, err
	}
	return sr, nil, nil
}

func (s *tankService) UpdateBW(survey *types.Survey, draftIndex int, tank types.Tank) (*types.SurveyResult, *Outcome, error) {
	const op = "TankService.UpdateBW"

	idx := slices.IndexFunc(survey.Drafts[draftIndex].BallastWaterTanks, func(t types.Tank) bool {
		return t.ID == tank.ID
	})
	if idx == -1 {
		s.logger.Error(op, ErrTankNotFound)
		return nil, nil, ErrTankNotFound
	}

	s.recalcVolume(survey, draftIndex, &tank)
	survey.Drafts[draftIndex].BallastWaterTanks[idx] = tank

	sr, err := s.saveAndCalc(op, survey)
	if err != nil {
		return nil, nil, err
	}
	return sr, nil, nil
}

func (s *tankService) UpdateFW(survey *types.Survey, draftIndex int, tank types.Tank) (*types.SurveyResult, *Outcome, error) {
	const op = "TankService.UpdateFW"

	idx := slices.IndexFunc(survey.Drafts[draftIndex].FreshWaterTanks, func(t types.Tank) bool {
		return t.ID == tank.ID
	})
	if idx == -1 {
		s.logger.Error(op, ErrTankNotFound)
		return nil, nil, ErrTankNotFound
	}

	tank.IsFWTTank = true
	survey.Drafts[draftIndex].FreshWaterTanks[idx] = tank

	sr, err := s.saveAndCalc(op, survey)
	if err != nil {
		return nil, nil, err
	}
	return sr, nil, nil
}

func (s *tankService) DeleteBW(survey *types.Survey, draftIndex int, tankID string) (*types.SurveyResult, *Outcome, error) {
	const op = "TankService.DeleteBW"

	idx := slices.IndexFunc(survey.Drafts[draftIndex].BallastWaterTanks, func(t types.Tank) bool {
		return t.ID == tankID
	})
	if idx == -1 {
		s.logger.Error(op, ErrTankNotFound)
		return nil, nil, ErrTankNotFound
	}
	survey.Drafts[draftIndex].BallastWaterTanks = slices.Delete(survey.Drafts[draftIndex].BallastWaterTanks, idx, idx+1)

	sr, err := s.saveAndCalc(op, survey)
	if err != nil {
		return nil, nil, err
	}
	return sr, deletedTankOutcome(), nil
}

func (s *tankService) DeleteFW(survey *types.Survey, draftIndex int, tankID string) (*types.SurveyResult, *Outcome, error) {
	const op = "TankService.DeleteFW"

	idx := slices.IndexFunc(survey.Drafts[draftIndex].FreshWaterTanks, func(t types.Tank) bool {
		return t.ID == tankID
	})
	if idx == -1 {
		s.logger.Error(op, ErrTankNotFound)
		return nil, nil, ErrTankNotFound
	}
	survey.Drafts[draftIndex].FreshWaterTanks = slices.Delete(survey.Drafts[draftIndex].FreshWaterTanks, idx, idx+1)

	sr, err := s.saveAndCalc(op, survey)
	if err != nil {
		return nil, nil, err
	}
	return sr, deletedTankOutcome(), nil
}

// ApplyDensity fills Density only on BW tanks where it is currently nil —
// never overrides a value the surveyor already entered per-tank.
func (s *tankService) ApplyDensity(survey *types.Survey, draftIndex int, density float64) (*types.SurveyResult, *Outcome, error) {
	const op = "TankService.ApplyDensity"

	for i := range survey.Drafts[draftIndex].BallastWaterTanks {
		if survey.Drafts[draftIndex].BallastWaterTanks[i].Density == nil {
			d := density
			survey.Drafts[draftIndex].BallastWaterTanks[i].Density = &d
		}
	}

	sr, err := s.saveAndCalc(op, survey)
	if err != nil {
		return nil, nil, err
	}
	return sr, nil, nil
}

// recalcVolume recomputes tank.VolumeCalc from the current calibration data,
// applied only when the calibration table is complete enough to trust
// (Correction.IsValid()). Never called for FW tanks — they have no
// calibration table; their volume is always the surveyor's direct entry.
func (s *tankService) recalcVolume(survey *types.Survey, draftIndex int, tank *types.Tank) {
	if !tank.Correction.IsValid() {
		return
	}

	dr := calculation.CalcDraft(survey.Drafts[draftIndex], survey.VesselData)
	vol, err := calculation.CalcBwTankVolume(dr.TrueTrim, dr.ListDegrees, *tank)
	if err != nil {
		return
	}
	tank.VolumeCalc = &vol
}

func (s *tankService) saveAndCalc(op string, survey *types.Survey) (*types.SurveyResult, error) {
	sr := calculation.CalcSurvey(*survey)

	if err := s.surveyRepo.Save(survey); err != nil {
		s.logger.Error(op, err)
		return nil, err
	}
	return &sr, nil
}

func deletedTankOutcome() *Outcome {
	return &Outcome{
		Toast: &Notification{
			Kind:    KindSuccess,
			Header:  "Deleted",
			Message: "Tank removed successfully",
		},
	}
}
