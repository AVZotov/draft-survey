package types

type VesselType string

const (
	VesselTypeMarine VesselType = "marine"
	VesselTypeRiver  VesselType = "river"
	VesselTypeBarge  VesselType = "barge"
)

type CorrectionMethod string

const (
	CorrectionMethodFullLBP CorrectionMethod = "Full LBP"
	CorrectionMethodHalfLBP CorrectionMethod = "Half LBP"
)

type PPDirection string

const (
	PPDirectionForward PPDirection = "F"
	PPDirectionAft     PPDirection = "A"
)

type VesselData struct {
	Name                 string           `json:"name"`
	Flag                 Country          `json:"flag"`
	IMO                  string           `json:"imo"`
	BuiltYear            *int             `json:"built_year"`
	Lightship            *float64         `json:"lightship"`
	Breadth              *float64         `json:"breadth"`
	Depth                *float64         `json:"depth"`
	LBP                  *float64         `json:"lbp"`
	SummerDraft          *float64         `json:"summer_draft"`
	SummerDWT            *float64         `json:"summer_dwt"`
	SummerTPC            *float64         `json:"summer_tpc"`
	SummerFreeboard      *float64         `json:"summer_freeboard"`
	VesselType           VesselType       `json:"vessel_type"`
	CorrectionMethod     CorrectionMethod `json:"correction_method"`
	IsLcfDetectionManual bool             `json:"is_lcf_detection_manual"`
	TableDensity         *float64         `json:"table_density"` //oficial vessel documentation. With table corrections made on specified density
	HoldsTotal           *int             `json:"holds_total"`
}
