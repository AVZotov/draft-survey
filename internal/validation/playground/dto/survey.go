package dto

import "github.com/AVZotov/draft-survey/internal/types"

type CargoOperation struct {
	Operation string `validate:"oneof=loading discharge"`
	Cargo     string `validate:"required"`
}

type Draft struct {
	Type             string           `validate:"oneof=initial intermediate final"`
	Status           string           `validate:"oneof=pending complete active"`
	SeaCondition     *SeaCondition    `validate:"required"`
	Marks            *Marks           `validate:"required"`
	Deductibles      *Deductibles     `validate:"required"`
	Density          *float64         `validate:"gte=1.000"`
	MTCRows          []MTCRow         `validate:"len=2,dive"`
	HydrostaticRows  []HydrostaticRow `validate:"len=2,dive"`
	TPCListPort      *float64         `validate:"required"`
	TPCListStarboard *float64         `validate:"required"`
	DistancePPFwd    *float64         `validate:"required"`
	PPFwdDirection   string           `validate:"oneof=A F"`
	DistancePPMid    *float64         `validate:"required"`
	PPMidDirection   string           `validate:"oneof=A F"`
	DistancePPAft    *float64         `validate:"required"`
	PPAftDirection   string           `validate:"oneof=A F"`
}

type Survey struct {
	Status           string          `validate:"oneof=draft in_progress complete"`
	ID               string          `validate:"uuid"`
	Drafts           []Draft         `validate:"min=2,dive"`
	CargoOperation   *CargoOperation `validate:"required"`
	CargoDeclared    *float64        `validate:"required"`
	ConstantDeclared *float64        `validate:"required"`
	VesselData       *VesselData     `validate:"required"`
}

func NewCargoOperationDTO(co *types.CargoOperation) *CargoOperation {
	if co == nil {
		return nil
	}
	return &CargoOperation{
		Operation: co.Operation,
		Cargo:     co.Cargo,
	}
}

func NewDraftDto(d *types.Draft) *Draft {
	if d == nil {
		return nil
	}
	var mtcRows = make([]MTCRow, 0)
	var hRows = make([]HydrostaticRow, 0)
	for _, mtc := range d.MTCRows {
		mtcRows = append(mtcRows, *NewMTCRowDTO(&mtc))
	}
	for _, h := range d.HydrostaticRows {
		hRows = append(hRows, *NewHydrostaticRowDTO(&h))
	}

	return &Draft{
		Type:             string(d.Type),
		Status:           string(d.Status),
		SeaCondition:     NewSeaConditionDTO(&d.SeaCondition),
		Marks:            NewMarksDto(&d.Marks),
		Deductibles:      NewDeductiblesDTO(&d.Deductibles),
		Density:          d.Density,
		MTCRows:          mtcRows,
		HydrostaticRows:  hRows,
		TPCListPort:      d.TPCListPort,
		TPCListStarboard: d.TPCListStarboard,
		DistancePPFwd:    d.DistancePPFwd,
		PPFwdDirection:   d.PPFwdDirection,
		DistancePPMid:    d.DistancePPMid,
		PPMidDirection:   d.PPMidDirection,
		DistancePPAft:    d.DistancePPAft,
		PPAftDirection:   d.PPAftDirection,
	}
}

func NewSurveyDTO(s *types.Survey) *Survey {
	if s == nil {
		return nil
	}
	var drafts = make([]Draft, 0)
	for _, d := range s.Drafts {
		drafts = append(drafts, *NewDraftDto(&d))
	}
	return &Survey{
		Status:           string(s.Status),
		ID:               s.ID,
		Drafts:           drafts,
		CargoOperation:   NewCargoOperationDTO(&s.CargoOperation),
		CargoDeclared:    s.CargoDeclared,
		ConstantDeclared: s.ConstantDeclared,
		VesselData:       NewVesselDto(&s.VesselData),
	}
}
