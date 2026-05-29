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

func derefFloat(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

func NewVesselDto(v *types.VesselData) *VesselData {
	if v == nil {
		return nil
	}
	return &VesselData{
		Name:             v.Name,
		IMO:              v.IMO,
		Lightship:        derefFloat(v.Lightship),
		Breadth:          derefFloat(v.Breadth),
		Depth:            derefFloat(v.Depth),
		LBP:              derefFloat(v.LBP),
		SummerDraft:      derefFloat(v.SummerDraft),
		SummerDWT:        derefFloat(v.SummerDWT),
		SummerTPC:        derefFloat(v.SummerTPC),
		SummerFreeboard:  derefFloat(v.SummerFreeboard),
		VesselType:       string(v.VesselType),
		CorrectionMethod: string(v.CorrectionMethod),
		TableDensity:     v.TableDensity,
	}
}
