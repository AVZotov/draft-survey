package types

type DraftResult struct {
	LBM                   float64
	LBMAftMid             float64
	LBMMidFwd             float64
	MeanDraft             MeanDraft
	MeanFwdAft            float64
	PPCorrections         PPCorrections
	DraftsWKeel           DraftsWKeel
	MMC                   float64
	ObservedTrim          float64
	TrueTrim              float64
	ListMeters            float64
	ListDegrees           float64
	Deflection            float64
	Hydrostatics          Hydrostatics
	FirstTrimCorrection   float64
	SecondTrimCorrection  float64
	DeltaMTC              float64
	ListCorrection        float64
	TotalTrimCorrection   float64
	DensityCorrection     float64
	DisplacementCorrected float64
	TotalDeductibles      float64
	NetDisplacement       float64
	Constant              float64
	CurrentDWT            float64
	TotalBwTanksWeight    float64
	TotalFwTanksWeight    float64
}
