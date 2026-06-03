package parse

import (
	"net/http"
	
	"github.com/AVZotov/draft-survey/internal/handler/fields"
	"github.com/AVZotov/draft-survey/internal/types"
)

func Tank(r *http.Request, tank *types.Tank) error {
	if tank == nil {
		return nil
	}
	
	tank.Name = String(r, fields.FieldTankName)
	tank.Type = String(r, fields.FieldTankType)
	
	if err := Float(r, fields.WithTankID(fields.FieldTankSounding, tank.ID), &tank.Sounding); err != nil {
		return err
	}
	if err := Float(r, fields.WithTankID(fields.FieldTankVolume, tank.ID), &tank.Volume); err != nil {
		return err
	}
	if !tank.IsFWTTank {
		if err := Float(r, fields.WithTankID(fields.FieldTankDensity, tank.ID), &tank.Density); err != nil {
			return err
		}
	}
	
	tableType := String(r, fields.FieldCalibType)
	if tableType == "" {
		return nil
	}
	tank.Correction.TableType = types.CalibrationTableType(tableType)
	tank.Correction.HasListCorrection = Bool(r, fields.FieldHasListCorrection)
	
	switch tank.Correction.TableType {
	case types.CalibrationTypeVolumeByTrim:
		if err := TrimRows(r, tank); err != nil {
			return err
		}
		if tank.Correction.HasListCorrection {
			if err := ListRows(r, tank); err != nil {
				return err
			}
		}
	case types.CalibrationTypeSoundingCorrection, types.CalibrationTypeVolumeCorrection:
		if err := TrimRows(r, tank); err != nil {
			return err
		}
		if err := CorrectionRows(r, tank); err != nil {
			return err
		}
		if tank.Correction.HasListCorrection {
			if err := ListRows(r, tank); err != nil {
				return err
			}
		}
	}
	return nil
}

func TrimRows(r *http.Request, tank *types.Tank) error {
	if len(tank.Correction.TrimRows) == 0 {
		tank.Correction.TrimRows = make([]types.CalibrationRow, 2)
	}
	if err := Float(r, fields.FieldTTL, &tank.Correction.TableTrimLow); err != nil {
		return err
	}
	if err := Float(r, fields.FieldTTU, &tank.Correction.TableTrimUpper); err != nil {
		return err
	}
	if err := Float(r, fields.FieldTrimTSLS, &tank.Correction.TrimRows[0].Sounding); err != nil {
		return err
	}
	if err := Float(r, fields.FieldTrimTSLVL, &tank.Correction.TrimRows[0].VolumeLow); err != nil {
		return err
	}
	if err := Float(r, fields.FieldTrimTSLVU, &tank.Correction.TrimRows[0].VolumeUp); err != nil {
		return err
	}
	if err := Float(r, fields.FieldTrimTSUS, &tank.Correction.TrimRows[1].Sounding); err != nil {
		return err
	}
	if err := Float(r, fields.FieldTrimTSUVL, &tank.Correction.TrimRows[1].VolumeLow); err != nil {
		return err
	}
	if err := Float(r, fields.FieldTrimTSUVU, &tank.Correction.TrimRows[1].VolumeUp); err != nil {
		return err
	}
	return nil
}

func ListRows(r *http.Request, tank *types.Tank) error {
	if len(tank.Correction.ListRows) == 0 {
		tank.Correction.ListRows = make([]types.CalibrationRow, 2)
	}
	if err := Float(r, fields.FieldTLL, &tank.Correction.TableListLow); err != nil {
		return err
	}
	if err := Float(r, fields.FieldTLU, &tank.Correction.TableListUpper); err != nil {
		return err
	}
	if err := Float(r, fields.FieldListTSLS, &tank.Correction.ListRows[0].Sounding); err != nil {
		return err
	}
	if err := Float(r, fields.FieldListTSLVL, &tank.Correction.ListRows[0].VolumeLow); err != nil {
		return err
	}
	if err := Float(r, fields.FieldListTSLVU, &tank.Correction.ListRows[0].VolumeUp); err != nil {
		return err
	}
	if err := Float(r, fields.FieldListTSUS, &tank.Correction.ListRows[1].Sounding); err != nil {
		return err
	}
	if err := Float(r, fields.FieldListTSUVL, &tank.Correction.ListRows[1].VolumeLow); err != nil {
		return err
	}
	if err := Float(r, fields.FieldListTSUVU, &tank.Correction.ListRows[1].VolumeUp); err != nil {
		return err
	}
	return nil
}

func CorrectionRows(r *http.Request, tank *types.Tank) error {
	if len(tank.Correction.CorrectionRows) == 0 {
		tank.Correction.CorrectionRows = make([]types.CorrectionRows, 2)
	}
	if err := Float(r, fields.FieldCorrSoundingLow, &tank.Correction.CorrectionRows[0].TableSounding); err != nil {
		return err
	}
	if err := Float(r, fields.FieldCorrVolumeLow, &tank.Correction.CorrectionRows[0].TableVolume); err != nil {
		return err
	}
	if err := Float(r, fields.FieldCorrSoundingUp, &tank.Correction.CorrectionRows[1].TableSounding); err != nil {
		return err
	}
	if err := Float(r, fields.FieldCorrVolumeUp, &tank.Correction.CorrectionRows[1].TableVolume); err != nil {
		return err
	}
	return nil
}
