package types

// DraftTotal pairs one draft's DraftResult with survey-level figures that
// depend on its position among the other drafts: cargo moved since the
// previous draft, the discrepancy against the surveyor-declared constant,
// and the running true constant after cumulative cargo is backed out.
type DraftTotal struct {
	DraftResult   DraftResult `json:"draft_result"`
	CargoFromPrev float64     `json:"cargo_from_prev"` // cargo moved since the previous draft
	ConstantDiff  float64     `json:"constant_diff"`   // difference from the surveyor-declared constant
	TrueConstant  float64     `json:"true_constant"`   // constant with cumulative cargo backed out
}

// SurveyResult is the full recalculation of a survey, produced by
// calculation.CalcSurvey. DraftTotals always has one entry per survey.Drafts
// entry, in the same order. Never persisted; always recomputed on the fly.
type SurveyResult struct {
	DraftTotals       []DraftTotal `json:"draft_totals"`
	CargoOnBoard      float64      `json:"cargo_on_board"`      // last draft's net displacement minus first draft's
	CargoDiffDeclared float64      `json:"cargo_diff_declared"` // cargo on board minus surveyor-declared cargo
}
