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

// MoveBWTankUp swaps the tank with its predecessor in the list. A no-op
// (still saved/recalculated) if the tank is already first.
func (s *tankService) MoveBWTankUp(survey *types.Survey, draftIndex int, tankID string) (*types.SurveyResult, *Outcome, error) {
	const op = "TankService.MoveBWTankUp"
	return s.moveTank(op, survey, draftIndex, tankID, true, true)
}

// MoveBWTankDown swaps the tank with its successor in the list. A no-op
// (still saved/recalculated) if the tank is already last.
func (s *tankService) MoveBWTankDown(survey *types.Survey, draftIndex int, tankID string) (*types.SurveyResult, *Outcome, error) {
	const op = "TankService.MoveBWTankDown"
	return s.moveTank(op, survey, draftIndex, tankID, true, false)
}

// MoveFWTankUp swaps the tank with its predecessor in the list. A no-op
// (still saved/recalculated) if the tank is already first.
func (s *tankService) MoveFWTankUp(survey *types.Survey, draftIndex int, tankID string) (*types.SurveyResult, *Outcome, error) {
	const op = "TankService.MoveFWTankUp"
	return s.moveTank(op, survey, draftIndex, tankID, false, true)
}

// MoveFWTankDown swaps the tank with its successor in the list. A no-op
// (still saved/recalculated) if the tank is already last.
func (s *tankService) MoveFWTankDown(survey *types.Survey, draftIndex int, tankID string) (*types.SurveyResult, *Outcome, error) {
	const op = "TankService.MoveFWTankDown"
	return s.moveTank(op, survey, draftIndex, tankID, false, false)
}

// moveTank finds tankID in the BW or FW list of the given draft and swaps it
// with the adjacent tank in the requested direction. Bounds are respected:
// the first tank cannot move up, the last cannot move down — both are silent
// no-ops rather than errors, since the button is expected to be disabled in
// that state and this is purely a display-order operation.
func (s *tankService) moveTank(op string, survey *types.Survey, draftIndex int, tankID string, isBW, up bool) (*types.SurveyResult, *Outcome, error) {
	var list *[]types.Tank
	if isBW {
		list = &survey.Drafts[draftIndex].BallastWaterTanks
	} else {
		list = &survey.Drafts[draftIndex].FreshWaterTanks
	}

	idx := slices.IndexFunc(*list, func(t types.Tank) bool {
		return t.ID == tankID
	})
	if idx == -1 {
		s.logger.Error(op, ErrTankNotFound)
		return nil, nil, ErrTankNotFound
	}

	other := idx - 1
	if !up {
		other = idx + 1
	}
	if other >= 0 && other < len(*list) {
		(*list)[idx], (*list)[other] = (*list)[other], (*list)[idx]
	}

	sr, err := s.saveAndCalc(op, survey)
	if err != nil {
		return nil, nil, err
	}
	return sr, nil, nil
}

// CopyFromPrevious replaces draftIndex's BW/FW tanks with structural copies
// (identity, naming, calibration table) of draftIndex-1's tanks, resetting
// per-draft measurements — same reset shape as DraftService.Add uses when a
// new draft is created from the previous one.
func (s *tankService) CopyFromPrevious(survey *types.Survey, draftIndex int) (*types.SurveyResult, *Outcome, error) {
	const op = "TankService.CopyFromPrevious"

	prev := survey.Drafts[draftIndex-1]
	survey.Drafts[draftIndex].BallastWaterTanks = copyTanksWithoutMeasurements(prev.BallastWaterTanks)
	survey.Drafts[draftIndex].FreshWaterTanks = copyTanksWithoutMeasurements(prev.FreshWaterTanks)

	sr, err := s.saveAndCalc(op, survey)
	if err != nil {
		return nil, nil, err
	}
	return sr, &Outcome{
		Toast: &Notification{
			Kind:    KindSuccess,
			Header:  "Copied",
			Message: "Tanks copied from previous draft",
		},
	}, nil
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
