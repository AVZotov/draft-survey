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

	newDraft := types.Draft{
		Type:            types.DraftTypeFinal,
		Status:          types.DraftStatusPending,
		MTCRows:         make([]types.MTCRow, 2),
		HydrostaticRows: make([]types.HydrostaticRow, 2),
		PPFwdDirection:  "A",
		PPMidDirection:  "A",
		PPAftDirection:  "F",
	}

	if len(survey.Drafts) > 0 {
		prev := survey.Drafts[len(survey.Drafts)-1]
		newDraft.BallastWaterTanks = copyTanksWithoutMeasurements(prev.BallastWaterTanks)
		newDraft.FreshWaterTanks = copyTanksWithoutMeasurements(prev.FreshWaterTanks)
	}

	survey.Drafts = append(survey.Drafts, newDraft)

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

	if index >= 0 && index < len(survey.Drafts) {
		mmc := calculation.CalcDraft(survey.Drafts[index], survey.VesselData).MMC
		autoPopulateMTCDrafts(&survey.Drafts[index], mmc)
	}

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

// copyTanksWithoutMeasurements carries a draft's tank structure (identity,
// naming, calibration table) into the next draft while resetting per-draft
// measurements — the surveyor re-enters soundings on each draft, but the
// tank list itself does not change between drafts. FW density is always 1.0
// (handled by Tank.CalcWeight), so Density is only reset for BW tanks.
func copyTanksWithoutMeasurements(src []types.Tank) []types.Tank {
	if len(src) == 0 {
		return nil
	}

	dst := make([]types.Tank, len(src))
	for i, t := range src {
		dst[i] = types.Tank{
			ID:         t.ID,
			Type:       t.Type,
			Name:       t.Name,
			IsFWTTank:  t.IsFWTTank,
			Correction: t.Correction,
		}
	}
	return dst
}

// autoPopulateMTCDrafts fills MTCRows[0] (MTC+) and MTCRows[1] (MTC-) from the
// Upper/Lower hydrostatic rows and the current MMC, using the surveyor's
// formula, whenever MTC+ has not been manually entered.
func autoPopulateMTCDrafts(draft *types.Draft, mmc float64) {
	if len(draft.MTCRows) < 2 || len(draft.HydrostaticRows) < 2 {
		return
	}
	if draft.MTCRows[0].Draft != nil {
		return
	}

	upper := draft.HydrostaticRows[0].Draft
	lower := draft.HydrostaticRows[1].Draft
	if upper == nil || lower == nil {
		return
	}
	draftUpper := *upper
	draftLower := *lower

	var mtcPlus, mtcMinus float64
	if ((draftLower+0.5)+(draftUpper+0.5))/2 <= mmc+0.5 {
		mtcPlus = draftLower + 0.5
	} else {
		mtcPlus = draftUpper + 0.5
	}
	if ((draftLower-0.5)+(draftUpper-0.5))/2 <= mmc-0.5 {
		mtcMinus = draftLower - 0.5
	} else {
		mtcMinus = draftUpper - 0.5
	}

	draft.MTCRows[0].Draft = &mtcPlus
	draft.MTCRows[1].Draft = &mtcMinus
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
