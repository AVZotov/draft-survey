package parse

import (
	"net/http"

	"github.com/AVZotov/draft-survey/internal/handler/fields"
	"github.com/AVZotov/draft-survey/internal/types"
)

// Tank updates only the measurement fields present in the request.
// existing must not be nil — caller ensures this.
func (d *Decoder) Tank(r *http.Request, existing *types.Tank) *types.Tank {
	if v, ok := d.pFloat(r, fields.WithTankID(fields.FieldTankSounding, existing.ID)); ok {
		existing.Sounding = v
	}
	if v, ok := d.pFloat(r, fields.WithTankID(fields.FieldTankVolume, existing.ID)); ok {
		existing.Volume = v
	}
	if v, ok := d.pFloat(r, fields.WithTankID(fields.FieldTankDensity, existing.ID)); ok {
		existing.Density = v
	}

	return existing
}

// TankMeta updates the type and name fields present in the request.
// Used when adding a new tank, where no tank ID exists yet to suffix the field names with.
// existing must not be nil — caller ensures this.
func (d *Decoder) TankMeta(r *http.Request, existing *types.Tank) *types.Tank {
	if r.PostForm.Has(fields.FieldTankType) {
		existing.Type = d.str(r, fields.FieldTankType)
	}
	if r.PostForm.Has(fields.FieldTankName) {
		existing.Name = d.str(r, fields.FieldTankName)
	}

	return existing
}

// TankCorrections updates existing.Correction with the calibration modal fields present
// in the request.
// existing must not be nil — caller ensures this.
func (d *Decoder) TankCorrections(r *http.Request, existing *types.Tank) *types.Tank {
	if r.PostForm.Has(fields.FieldCalibType) {
		existing.Correction.TableType = types.CalibrationTableType(d.str(r, fields.FieldCalibType))
	}
	if v, ok := d.pFloat(r, fields.FieldTTL); ok {
		existing.Correction.TableTrimLow = v
	}
	if v, ok := d.pFloat(r, fields.FieldTTU); ok {
		existing.Correction.TableTrimUpper = v
	}

	d.tankTrimRows(r, existing)

	existing.Correction.HasListCorrection = d.boolean(r, fields.FieldHasListCorrection)
	if v, ok := d.pFloat(r, fields.FieldTLL); ok {
		existing.Correction.TableListLow = v
	}
	if v, ok := d.pFloat(r, fields.FieldTLU); ok {
		existing.Correction.TableListUpper = v
	}

	d.tankListRows(r, existing)
	d.tankCorrectionRows(r, existing)

	return existing
}

func (d *Decoder) tankTrimRows(r *http.Request, existing *types.Tank) {
	if len(existing.Correction.TrimRows) == 0 {
		existing.Correction.TrimRows = make([]types.CalibrationRow, 2)
	}

	if v, ok := d.pFloat(r, fields.FieldTrimTSLS); ok {
		existing.Correction.TrimRows[0].Sounding = v
	}
	if v, ok := d.pFloat(r, fields.FieldTrimTSLVL); ok {
		existing.Correction.TrimRows[0].VolumeLow = v
	}
	if v, ok := d.pFloat(r, fields.FieldTrimTSLVU); ok {
		existing.Correction.TrimRows[0].VolumeUp = v
	}
	if v, ok := d.pFloat(r, fields.FieldTrimTSUS); ok {
		existing.Correction.TrimRows[1].Sounding = v
	}
	if v, ok := d.pFloat(r, fields.FieldTrimTSUVL); ok {
		existing.Correction.TrimRows[1].VolumeLow = v
	}
	if v, ok := d.pFloat(r, fields.FieldTrimTSUVU); ok {
		existing.Correction.TrimRows[1].VolumeUp = v
	}
}

func (d *Decoder) tankListRows(r *http.Request, existing *types.Tank) {
	if len(existing.Correction.ListRows) == 0 {
		existing.Correction.ListRows = make([]types.CalibrationRow, 2)
	}

	if v, ok := d.pFloat(r, fields.FieldListTSLS); ok {
		existing.Correction.ListRows[0].Sounding = v
	}
	if v, ok := d.pFloat(r, fields.FieldListTSLVL); ok {
		existing.Correction.ListRows[0].VolumeLow = v
	}
	if v, ok := d.pFloat(r, fields.FieldListTSLVU); ok {
		existing.Correction.ListRows[0].VolumeUp = v
	}
	if v, ok := d.pFloat(r, fields.FieldListTSUS); ok {
		existing.Correction.ListRows[1].Sounding = v
	}
	if v, ok := d.pFloat(r, fields.FieldListTSUVL); ok {
		existing.Correction.ListRows[1].VolumeLow = v
	}
	if v, ok := d.pFloat(r, fields.FieldListTSUVU); ok {
		existing.Correction.ListRows[1].VolumeUp = v
	}
}

func (d *Decoder) tankCorrectionRows(r *http.Request, existing *types.Tank) {
	if len(existing.Correction.CorrectionRows) == 0 {
		existing.Correction.CorrectionRows = make([]types.CorrectionRows, 2)
	}

	if v, ok := d.pFloat(r, fields.FieldCorrSoundingLow); ok {
		existing.Correction.CorrectionRows[0].TableSounding = v
	}
	if v, ok := d.pFloat(r, fields.FieldCorrVolumeLow); ok {
		existing.Correction.CorrectionRows[0].TableVolume = v
	}
	if v, ok := d.pFloat(r, fields.FieldCorrSoundingUp); ok {
		existing.Correction.CorrectionRows[1].TableSounding = v
	}
	if v, ok := d.pFloat(r, fields.FieldCorrVolumeUp); ok {
		existing.Correction.CorrectionRows[1].TableVolume = v
	}
}

// DockwaterDensity parses the single density field from the Apply Density form.
// bool reports whether the field was present in the request.
func (d *Decoder) DockwaterDensity(r *http.Request) (*float64, bool) {
	return d.pFloat(r, fields.FieldDockwaterDensity)
}
