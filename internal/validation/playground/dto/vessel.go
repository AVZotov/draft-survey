package dto

import "github.com/AVZotov/draft-survey/internal/types"

type VesselData struct {
	Name             string   `validate:"required"`
	IMO              string   `validate:"len=7"`
	Lightship        float64  `validate:"required"`
	Breadth          float64  `validate:"required"`
	Depth            float64  `validate:"required"`
	LBP              float64  `validate:"required"`
	SummerDraft      float64  `validate:"required"`
	SummerDWT        float64  `validate:"required"`
	SummerTPC        float64  `validate:"required"`
	SummerFreeboard  float64  `validate:"required"`
	VesselType       string   `validate:"oneof=marine river barge"`
	CorrectionMethod string   `validate:"lbp"`
	TableDensity     *float64 `validate:"gte=1.000"`
}

func NewVesselDto(v *types.VesselData) *VesselData {
	if v == nil {
		return nil
	}
	return &VesselData{
		Name:             v.Name,
		IMO:              v.IMO,
		Lightship:        v.Lightship,
		Breadth:          v.Breadth,
		Depth:            v.Depth,
		LBP:              v.LBP,
		SummerDraft:      v.SummerDraft,
		SummerDWT:        v.SummerDWT,
		SummerTPC:        v.SummerTPC,
		SummerFreeboard:  v.SummerFreeboard,
		VesselType:       string(v.VesselType),
		CorrectionMethod: string(v.CorrectionMethod),
		TableDensity:     v.TableDensity,
	}
}
