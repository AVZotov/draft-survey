package types

type DraftTotal struct {
	DraftResult   DraftResult
	CargoFromPrev float64 // груз с предыдущего драфта
	ConstantDiff  float64 // разница с заявленной константой
}

type SurveyResult struct {
	DraftTotals       []DraftTotal
	CargoOnBoard      float64 // финал - начало
	CargoDiffDeclared float64 // только финал vs заявленный
}
