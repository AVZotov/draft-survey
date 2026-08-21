package types

type Tank struct {
	IsFWTTank  bool                  `json:"is_fresh_water_tank"`
	Type       string                `json:"tank_type"`
	Name       string                `json:"tank_name"`
	ID         string                `json:"tank_id"`
	Sounding   *float64              `json:"tank_sounding"`
	Density    *float64              `json:"tank_density"`
	Volume     *float64              `json:"tank_volume"`
	VolumeCalc *float64              `json:"volume_calc"`
	Correction VolumeCalibrationData `json:"correction"`
}

// EffectiveVolume returns the surveyor's manual override (Volume) when present,
// otherwise the calibration-calculated volume (VolumeCalc).
func (t Tank) EffectiveVolume() *float64 {
	if t.Volume != nil {
		return t.Volume
	}
	return t.VolumeCalc
}

func (t Tank) CalcWeight() float64 {
	vol := t.EffectiveVolume()
	if vol == nil {
		return 0.000
	}
	if t.IsFWTTank {
		return *vol
	}

	if t.Density == nil {
		return 0.000
	}
	return *vol * *t.Density
}

type Deductibles struct {
	HFO         *float64 `json:"hfo"`
	MDO         *float64 `json:"mdo"`
	LubOil      *float64 `json:"lub_oil"`
	BilgeWater  *float64 `json:"bilge_water"`
	SewageWater *float64 `json:"sewage_water"`
	Others      *float64 `json:"others"`
}
