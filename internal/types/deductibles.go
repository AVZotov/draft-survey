package types

type Tank struct {
	IsFWTTank  bool                  `json:"is_fresh_water_tank"`
	Type       string                `json:"tank_type"`
	Name       string                `json:"tank_name"`
	ID         string                `json:"tank_id"`
	Sounding   *float64              `json:"tank_sounding"`
	Density    *float64              `json:"tank_density"`
	Volume     *float64              `json:"tank_volume"`
	Correction VolumeCalibrationData `json:"correction"`
}

func (t Tank) CalcWeight() float64 {
	if t.Volume == nil {
		return 0.000
	}
	if t.IsFWTTank {
		return *t.Volume
	}

	if t.Density == nil {
		return 0.000
	}
	return *t.Volume * *t.Density
}

type OtherDeductibles struct {
	Others     *float64 `json:"others"`
	OthersName string   `json:"others_name"`
}
type Deductibles struct {
	HFO         *float64 `json:"hfo"`
	MDO         *float64 `json:"mdo"`
	LubOil      *float64 `json:"lub_oil"`
	BilgeWater  *float64 `json:"bilge_water"`
	SewageWater *float64 `json:"sewage_water"`
	OtherDeductibles
}
