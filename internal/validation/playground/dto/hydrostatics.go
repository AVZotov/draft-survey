package dto

import "github.com/AVZotov/draft-survey/internal/types"

type HydrostaticRow struct {
	Draft        *float64 `validate:"required"`
	Displacement *float64 `validate:"required"`
	TPC          *float64 `validate:"required"`
	LCF          *float64 `validate:"required"`
	LCFDirection string   `validate:"oneof=F A AP"`
}

type MTCRow struct {
	Draft *float64 `validate:"required"`
	MTC   *float64 `validate:"required"`
}

func NewHydrostaticRowDTO(hRow *types.HydrostaticRow) *HydrostaticRow {
	if hRow == nil {
		return nil
	}
	return &HydrostaticRow{
		Draft:        hRow.Draft,
		Displacement: hRow.Displacement,
		TPC:          hRow.TPC,
		LCF:          hRow.LCF,
		LCFDirection: string(hRow.LCFDirection),
	}
}

func NewMTCRowDTO(mtc *types.MTCRow) *MTCRow {
	if mtc == nil {
		return nil
	}
	return &MTCRow{
		Draft: mtc.Draft,
		MTC:   mtc.MTC,
	}
}
