package calculation

import (
	"testing"

	"github.com/AVZotov/draft-survey/internal/types"
)

func ptr(v float64) *float64 {
	return &v
}

func getVolumeCalibrationType1Tank() types.Tank {
	return types.Tank{
		Sounding: ptr(1.990),
		Correction: types.VolumeCalibrationData{
			TableType:         types.CalibrationTypeVolumeByTrim,
			TableTrimLow:      ptr(2.000),
			TableTrimUpper:    ptr(3.000),
			HasListCorrection: false,
			TrimRows: []types.CalibrationRow{
				{Sounding: ptr(1.950), VolumeLow: ptr(808.780), VolumeUp: ptr(782.280)},
				{Sounding: ptr(2.000), VolumeLow: ptr(824.850), VolumeUp: ptr(800.410)},
			},
		},
	}
}

func getVolumeCalibrationType2Tank() types.Tank {
	return types.Tank{
		Sounding: ptr(2.460),
		Correction: types.VolumeCalibrationData{
			TableType:         types.CalibrationTypeSoundingCorrection,
			TableTrimLow:      ptr(2.500),
			TableTrimUpper:    ptr(3.000),
			HasListCorrection: false,
			TrimRows: []types.CalibrationRow{
				{Sounding: ptr(2.400), VolumeLow: ptr(-205), VolumeUp: ptr(-246)},
				{Sounding: ptr(2.500), VolumeLow: ptr(-207), VolumeUp: ptr(-248)},
			},
			CorrectionRows: []types.CorrectionRows{
				{TableSounding: ptr(2.220), TableVolume: ptr(565.540)},
				{TableSounding: ptr(2.230), TableVolume: ptr(556.370)},
			},
		},
	}
}

func getVolumeCalibrationType3Tank() types.Tank {
	return types.Tank{
		Sounding: ptr(1.980),
		Correction: types.VolumeCalibrationData{
			TableType:         types.CalibrationTypeVolumeCorrection,
			TableTrimLow:      ptr(2.000),
			TableTrimUpper:    ptr(3.000),
			HasListCorrection: false,
			TrimRows: []types.CalibrationRow{
				{Sounding: ptr(1.900), VolumeLow: ptr(-56.000), VolumeUp: ptr(-89.200)},
				{Sounding: ptr(2.000), VolumeLow: ptr(-58.630), VolumeUp: ptr(-97.150)},
			},
			CorrectionRows: []types.CorrectionRows{
				{TableSounding: ptr(1.900), TableVolume: ptr(892.060)},
				{TableSounding: ptr(2.000), TableVolume: ptr(899.250)},
			},
		},
	}
}

func TestVolumeCalibrationType1(t *testing.T) {
	const volume = 801.754
	const trim = 2.800

	tank := getVolumeCalibrationType1Tank()
	tankVolume, err := CalcBwTankVolume(trim, 0, tank)
	if err != nil {
		t.Errorf("function error %s", err.Error())
	}
	if tankVolume != volume {
		t.Errorf("Expected %f, got %f", volume, tankVolume)
	}
}

func TestVolumeCalibrationType2(t *testing.T) {
	const volume = 557.287
	const trim = 2.800

	tank := getVolumeCalibrationType2Tank()
	tankVolume, err := CalcBwTankVolume(trim, 0, tank)
	if err != nil {
		t.Errorf("function error %s", err.Error())
	}
	if tankVolume != volume {
		t.Errorf("Expected %f, got %f", volume, tankVolume)
	}
}

func TestVolumeCalibrationType3(t *testing.T) {
	const volume = 809.743
	const trim = 2.800

	tank := getVolumeCalibrationType3Tank()
	tankVolume, err := CalcBwTankVolume(trim, 0, tank)
	if err != nil {
		t.Errorf("function error %s", err.Error())
	}
	if tankVolume != volume {
		t.Errorf("Expected %f, got %f", volume, tankVolume)
	}
}
