package types

type DraftResult struct {
	LBM                   float64       `json:"lbm"`
	LBMAftMid             float64       `json:"lbm_aft_mid"`
	LBMMidFwd             float64       `json:"lbm_mid_fwd"`
	MeanDraft             MeanDraft     `json:"mean_draft"`
	MeanFwdAft            float64       `json:"mean_fwd_aft"`
	PPCorrections         PPCorrections `json:"pp_corrections"`
	DraftsWKeel           DraftsWKeel   `json:"drafts_w_keel"`
	MMC                   float64       `json:"mmc"`
	ObservedTrim          float64       `json:"observed_trim"`
	TrueTrim              float64       `json:"true_trim"`
	ListMeters            float64       `json:"list_meters"`
	ListDegrees           float64       `json:"list_degrees"`
	Deflection            float64       `json:"deflection"`
	Hydrostatics          Hydrostatics  `json:"hydrostatics"`
	FirstTrimCorrection   float64       `json:"first_trim_correction"`
	SecondTrimCorrection  float64       `json:"second_trim_correction"`
	DeltaMTC              float64       `json:"delta_mtc"`
	ListCorrection        float64       `json:"list_correction"`
	TotalTrimCorrection   float64       `json:"total_trim_correction"`
	DensityCorrection     float64       `json:"density_correction"`
	DisplacementCorrected float64       `json:"displacement_corrected"`
	TotalDeductibles      float64       `json:"total_deductibles"`
	NetDisplacement       float64       `json:"net_displacement"`
	Constant              float64       `json:"constant"`
	CurrentDWT            float64       `json:"current_dwt"`
	TotalBwTanksWeight    float64       `json:"total_bw_tanks_weight"`
	TotalFwTanksWeight    float64       `json:"total_fw_tanks_weight"`
}
