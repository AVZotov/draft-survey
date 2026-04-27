package calculation

import (
	"github.com/AVZotov/draft-survey/internal/types"
)

func CalcSurvey(s types.Survey) types.SurveyResult {
	draftTotals := make([]types.DraftTotal, 0, len(s.Drafts))
	var previousDR types.DraftResult
	var cumulativeCargo float64

	for i, draft := range s.Drafts {
		dr := CalcDraft(draft, s.VesselData)
		dt := draftTotal(s, dr)

		if i > 0 {
			dt.CargoFromPrev = round3(dr.NetDisplacement - previousDR.NetDisplacement)
			cumulativeCargo = round3(cumulativeCargo + dt.CargoFromPrev)
		}

		dt.TrueConstant = round3(dr.Constant - cumulativeCargo)
		draftTotals = append(draftTotals, dt)
		previousDR = dr
	}

	var cargoOnBoard float64
	var cargoDiffDeclared float64

	if len(draftTotals) >= 2 {
		first := draftTotals[0].DraftResult.NetDisplacement
		last := draftTotals[len(draftTotals)-1].DraftResult.NetDisplacement
		cargoOnBoard = round3(last - first)

		if s.CargoDeclared != nil {
			cargoDiffDeclared = round3(cargoOnBoard - *s.CargoDeclared)
		}
	}

	return types.SurveyResult{
		DraftTotals:       draftTotals,
		CargoOnBoard:      cargoOnBoard,
		CargoDiffDeclared: cargoDiffDeclared,
	}
}

func draftTotal(s types.Survey, dr types.DraftResult) types.DraftTotal {
	var dt types.DraftTotal
	dt.DraftResult = dr
	if s.ConstantDeclared != nil {
		dt.ConstantDiff = dr.Constant - *s.ConstantDeclared
	}
	return dt
}
