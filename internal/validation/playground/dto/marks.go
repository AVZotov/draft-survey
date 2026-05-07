package dto

import (
	"github.com/AVZotov/draft-survey/internal/types"
)

type Mark struct {
	Value  *float64 `validate:"required"`
	Method string   `validate:"oneof=direct camera"`
}

type Marks struct {
	FwdPort      Mark
	FwdStarboard Mark
	MidPort      Mark
	MidStarboard Mark
	AftPort      Mark
	AftStarboard Mark
}

func NewMarksDto(marks *types.Marks) *Marks {
	if marks == nil {
		return nil
	}
	return &Marks{
		FwdPort: Mark{
			Value:  marks.FwdPort.Value,
			Method: string(marks.FwdPort.Method),
		},
		FwdStarboard: Mark{
			Value:  marks.FwdStarboard.Value,
			Method: string(marks.FwdStarboard.Method),
		},
		MidPort: Mark{
			Value:  marks.MidPort.Value,
			Method: string(marks.MidPort.Method),
		},
		MidStarboard: Mark{
			Value:  marks.MidStarboard.Value,
			Method: string(marks.MidStarboard.Method),
		},
		AftPort: Mark{
			Value:  marks.AftPort.Value,
			Method: string(marks.AftPort.Method),
		},
		AftStarboard: Mark{
			Value:  marks.AftStarboard.Value,
			Method: string(marks.AftStarboard.Method),
		},
	}
}
