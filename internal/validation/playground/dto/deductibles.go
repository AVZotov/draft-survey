package dto

import "github.com/AVZotov/draft-survey/internal/types"

type Deductibles struct {
	HFO        *float64 `validate:"required"`
	MDO        *float64 `validate:"required"`
	LubOil     *float64 `validate:"required"`
	BilgeWater *float64 `validate:"required"`
	Others     *float64 `validate:"required"`
}

func NewDeductiblesDTO(d *types.Deductibles) *Deductibles {
	if d == nil {
		return nil
	}
	return &Deductibles{
		HFO:        d.HFO,
		MDO:        d.MDO,
		LubOil:     d.LubOil,
		BilgeWater: d.BilgeWater,
		Others:     d.Others,
	}
}
