package vessel

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
	Flag                 string           `json:"flag"`
	IMO                  string           `json:"imo"`
	BuiltYear            int              `json:"built_year"`
	Lightship            float64          `json:"lightship"` // вес порожнем, MT
	Breadth              float64          `json:"breadth"`   // ширина корпуса, м
	Depth                float64          `json:"depth"`     // высота корпуса, м
	LBP                  float64          `json:"lbp"`       // длина между перпендикулярами, м
	SummerDraft          float64          `json:"summer_draft"`
	SummerDWT            float64          `json:"summer_dwt"`
	SummerTPC            float64          `json:"summer_tpc"`
	SummerFreeboard      float64          `json:"summer_freeboard"`
	DistancePPFwd        float64          `json:"distance_pp_fwd"`
	PPFwdDirection       PPDirection      `json:"pp_fwd_direction"`
	DistancePPMid        float64          `json:"distance_pp_mid"`
	PPMidDirection       PPDirection      `json:"pp_mid_direction"`
	DistancePPAft        float64          `json:"distance_pp_aft"`
	PPAftDirection       PPDirection      `json:"pp_aft_direction"`
	KeelFwd              float64          `json:"keel_fwd"`
	KeelMid              float64          `json:"keel_mid"`
	KeelAft              float64          `json:"keel_aft"`
	VesselType           VesselType       `json:"vessel_type"`
	CorrectionMethod     CorrectionMethod `json:"correction_method"`
	IsLcfDetectionManual bool             `json:"is_lcf_detection_manual"`
	TableDensity         *float64         `json:"table_density"`
}
