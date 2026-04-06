package types

import (
	"time"

	"github.com/AVZotov/draft-survey/internal/vessel"
)

type MeanDraft struct {
	DraftFwdMean float64 `json:"draft_fwd_mean"`
	DraftMidMean float64 `json:"draft_mid_mean"`
	DraftAftMean float64 `json:"draft_aft_mean"`
}

type PPCorrections struct {
	FwdCorrection float64 `json:"fwd_correction"`
	MidCorrection float64 `json:"mid_correction"`
	AftCorrection float64 `json:"aft_correction"`
}

type DraftsWKeel struct {
	FwdDraftWKeel float64 `json:"fwd_draft_w_keel"`
	MidDraftWKeel float64 `json:"mid_draft_w_keel"`
	AftDraftWKeel float64 `json:"aft_draft_w_keel"`
}

type Job struct {
	JobNumber string `json:"job_number"`
	DSNumber  int    `json:"ds_number"`
	Principal string `json:"principal"`
}

type CargoOperation struct {
	PlaceOfInspection string `json:"place_of_inspection"`
	Destination       string `json:"destination"`
	Operation         string `json:"operation"`
	Cargo             string `json:"cargo"`
	Packing           string `json:"packing"`
}

type SurveyStatus string

const (
	SurveyStatusDraft      SurveyStatus = "draft"
	SurveyStatusInProgress SurveyStatus = "in_progress"
	SurveyStatusComplete   SurveyStatus = "complete"
)

type DraftType string

const (
	DraftTypeInitial      DraftType = "initial"
	DraftTypeIntermediate DraftType = "intermediate"
	DraftTypeFinal        DraftType = "final"
)

type DraftStatus string

const (
	DraftStatusPending  DraftStatus = "pending"
	DraftStatusActive   DraftStatus = "active"
	DraftStatusComplete DraftStatus = "complete"
)

type Draft struct {
	Type              DraftType        `json:"type"`
	Status            DraftStatus      `json:"status"`
	SeaCondition      SeaCondition     `json:"sea_condition"`
	Marks             Marks            `json:"marks"`
	Deductibles       Deductibles      `json:"deductibles"`
	BallastWaterTanks []Tank           `json:"ballast_water_tanks"`
	FreshWaterTanks   []Tank           `json:"fresh_water_tanks"`
	Density           *float64         `json:"density"`
	MTCRows           []MTCRow         `json:"mtc_rows"`
	HydrostaticRows   []HydrostaticRow `json:"hydrostatic_rows"`
	TPCListPort       float64          `json:"tpc_list_port"`
	TPCListStarboard  float64          `json:"tpc_list_starboard"`
	DistancePPFwd     *float64         `json:"distance_pp_fwd"`
	PPFwdDirection    string           `json:"pp_fwd_direction"`
	DistancePPMid     *float64         `json:"distance_pp_mid"`
	PPMidDirection    string           `json:"pp_mid_direction"`
	DistancePPAft     *float64         `json:"distance_pp_aft"`
	PPAftDirection    string           `json:"pp_aft_direction"`
	KeelFwd           *float64         `json:"keel_fwd"`
	KeelMid           *float64         `json:"keel_mid"`
	KeelAft           *float64         `json:"keel_aft"`
	StartedAt         time.Time        `json:"started_at"`
	FinishedAt        time.Time        `json:"finished_at"`
	LoadingCommenced  time.Time        `json:"loading_commenced"`
}

type Survey struct {
	Surveyor         User              `json:"surveyor"`
	Status           SurveyStatus      `json:"status"`
	ID               string            `json:"id"`
	CreatedAt        time.Time         `json:"created_at"`
	Drafts           []Draft           `json:"drafts"`
	Job              Job               `json:"job"`
	CargoOperation   CargoOperation    `json:"cargo_operation"`
	CargoDeclared    *float64          `json:"cargo_declared"`
	ConstantDeclared *float64          `json:"constant_declared"`
	VesselData       vessel.VesselData `json:"vessel_data"`
	SeaCondition     SeaCondition      `json:"sea_condition"`
	Remarks          string            `json:"remarks"`
}
