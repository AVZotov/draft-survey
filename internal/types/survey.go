package types

import (
	"time"
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
	DSNumber  *int   `json:"ds_number"`
	Principal string `json:"client"`
}

// PortInfo holds full port data — country, port from UN/LOCODE dict,
// and optional manual override (prevails over dict selection)
type PortInfo struct {
	Country    Country `json:"country"`
	Port       Port    `json:"port"`
	PortManual string  `json:"port_manual"`
}

// DisplayName returns manual port if set, otherwise port name from dict
func (p PortInfo) DisplayName() string {
	if p.PortManual != "" {
		return p.PortManual
	}
	if p.Port.Name != "" {
		return p.Port.Name
	}
	return ""
}

type Operation string

const (
	OperationLoading   Operation = "loading"
	OperationDischarge Operation = "discharge"
)

type CargoOperation struct {
	Operation   Operation `json:"cargo_operation"`
	Cargo       string    `json:"cargo"`
	Description string    `json:"description"`
	Packing     string    `json:"packing"`
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
	Surveyor          User             `json:"surveyor"`
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
	TPCListPort       *float64         `json:"tpc_list_port"`
	TPCListStarboard  *float64         `json:"tpc_list_stbd"`
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
	LoadPort          Port             `json:"load_port"`
	DischargePort     Port             `json:"discharge_port"`
	ProductGrade      string           `json:"product_grade"`
	HoldsActive       []int            `json:"holds_active"`
}

type Survey struct {
	Surveyor         User           `json:"surveyor"`
	Status           SurveyStatus   `json:"status"`
	ID               string         `json:"id"`
	CreatedAt        time.Time      `json:"created_at"`
	Drafts           []Draft        `json:"drafts"`
	Job              Job            `json:"job"`
	CargoOperation   CargoOperation `json:"cargo_operation"`
	CargoDeclared    *float64       `json:"cargo_declared"`
	ConstantDeclared *float64       `json:"const_declared"`
	VesselData       VesselData     `json:"vessel_data"`
	Remarks          string         `json:"remarks"`
	Country          Country        `json:"country"`
	Audit            AuditEvent     `json:"audit"`
	MetaData         MetaData       `json:"metadata"`
}
