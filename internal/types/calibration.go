package types

import "github.com/AVZotov/draft-survey/internal/constants"

type CalibrationTableType string

const (
	CalibrationTypeVolumeByTrim       CalibrationTableType = constants.VolumeCalibrationType1
	CalibrationTypeSoundingCorrection CalibrationTableType = constants.VolumeCalibrationType2
	CalibrationTypeVolumeCorrection   CalibrationTableType = constants.VolumeCalibrationType3
)

type CalibrationRow struct {
	Sounding  *float64 `json:"calib_row_sounding"`
	VolumeLow *float64 `json:"calib_row_vol_sound_low"`
	VolumeUp  *float64 `json:"calib_row_vol_sound_up"`
}

type SoundingVolume struct {
	TableSounding *float64 `json:"calib_row_sounding"`
	TableVolume   *float64 `json:"calib_row_vol_sound"`
}

type VolumeCalibrationData struct {
	TableType         CalibrationTableType `json:"tank_calib_table_type"`
	TableTrimLow      *float64             `json:"calib_ttl"`
	TableTrimUpper    *float64             `json:"calib_ttu"`
	TableListLow      *float64             `json:"calib_tll"`
	TableListUpper    *float64             `json:"calib_tlu"`
	HasListCorrection bool                 `json:"calib_has_list_correction"`
	VolumeByTrimRow   []CalibrationRow     `json:"trim_rows"`
	ListRows          []CalibrationRow     `json:"list_rows"`
	VolumeRows        []SoundingVolume     `json:"volume_rows"` // Table 2 for Type2 and Type3
}

func (v VolumeCalibrationData) IsValid() bool {
	if v.TableType == "" || v.TableTrimLow == nil || v.TableTrimUpper == nil {
		return false
	}

	if v.HasListCorrection {
		if len(v.ListRows) < 2 || v.TableListLow == nil || v.TableListUpper == nil {
			return false
		}
		for _, r := range v.ListRows {
			if r.Sounding == nil || r.VolumeLow == nil || r.VolumeUp == nil {
				return false
			}
		}
	}

	switch v.TableType {
	case CalibrationTypeVolumeByTrim:
		if len(v.VolumeByTrimRow) < 2 {
			return false
		}
		for _, r := range v.VolumeByTrimRow {
			if r.Sounding == nil || r.VolumeLow == nil || r.VolumeUp == nil {
				return false
			}
		}
		return true
	case CalibrationTypeSoundingCorrection, CalibrationTypeVolumeCorrection:
		if len(v.VolumeRows) < 2 {
			return false
		}
		for _, r := range v.VolumeRows {
			if r.TableSounding == nil || r.TableVolume == nil {
				return false
			}
		}
		return true
	}
	return false
}
