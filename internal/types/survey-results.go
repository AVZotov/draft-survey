package types

type DraftTotal struct {
	DraftResult   DraftResult `json:"draft_result"`
	CargoFromPrev float64     `json:"cargo_from_prev"` // груз с предыдущего драфта
	ConstantDiff  float64     `json:"constant_diff"`   // разница с заявленной константой
	TrueConstant  float64     `json:"true_constant"`   // always calculate constant substract cargo loaded
}

type SurveyResult struct {
	DraftTotals       []DraftTotal `json:"draft_totals"`
	CargoOnBoard      float64      `json:"cargo_on_board"`      // финал - начало
	CargoDiffDeclared float64      `json:"cargo_diff_declared"` // только финал vs заявленный
}
