package parse

import (
	"fmt"
	"net/http"

	"github.com/AVZotov/draft-survey/internal/handler/fields"
	"github.com/AVZotov/draft-survey/internal/types"
)

// Draft updates only the fields present in the request.
// existing must not be nil — caller ensures this via services.Survey.Get()
func (d *Decoder) Draft(r *http.Request, existing *types.Draft, index int) *types.Draft {
	if existing == nil {
		existing = &types.Draft{}
	}

	existing.Marks = d.marksBlock(r, existing.Marks, index)
	existing.Deductibles = d.deductiblesBlock(r, existing.Deductibles, index)
	existing.SeaCondition = d.seaConditionBlock(r, existing.SeaCondition, index)
	d.ppKeelBlock(r, existing, index)
	d.hydrostaticsBlock(r, existing, index)
	d.mtcBlock(r, existing, index)
	d.draftMetaBlock(r, existing, index)

	return existing
}

func (d *Decoder) marksBlock(r *http.Request, existing types.Marks, index int) types.Marks {
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldFwdPort, index)); ok {
		existing.FwdPort.Value = v
	}
	if r.PostForm.Has(fmt.Sprintf("%s-%d", fields.FieldFwdPortMethod, index)) {
		existing.FwdPort.Method = types.ReadingMethod(d.str(r, fmt.Sprintf("%s-%d", fields.FieldFwdPortMethod, index)))
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldFwdStbd, index)); ok {
		existing.FwdStarboard.Value = v
	}
	if r.PostForm.Has(fmt.Sprintf("%s-%d", fields.FieldFwdStbdMethod, index)) {
		existing.FwdStarboard.Method = types.ReadingMethod(d.str(r, fmt.Sprintf("%s-%d", fields.FieldFwdStbdMethod, index)))
	}

	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldMidPort, index)); ok {
		existing.MidPort.Value = v
	}
	if r.PostForm.Has(fmt.Sprintf("%s-%d", fields.FieldMidPortMethod, index)) {
		existing.MidPort.Method = types.ReadingMethod(d.str(r, fmt.Sprintf("%s-%d", fields.FieldMidPortMethod, index)))
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldMidStbd, index)); ok {
		existing.MidStarboard.Value = v
	}
	if r.PostForm.Has(fmt.Sprintf("%s-%d", fields.FieldMidStbdMethod, index)) {
		existing.MidStarboard.Method = types.ReadingMethod(d.str(r, fmt.Sprintf("%s-%d", fields.FieldMidStbdMethod, index)))
	}

	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldAftPort, index)); ok {
		existing.AftPort.Value = v
	}
	if r.PostForm.Has(fmt.Sprintf("%s-%d", fields.FieldAftPortMethod, index)) {
		existing.AftPort.Method = types.ReadingMethod(d.str(r, fmt.Sprintf("%s-%d", fields.FieldAftPortMethod, index)))
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldAftStbd, index)); ok {
		existing.AftStarboard.Value = v
	}
	if r.PostForm.Has(fmt.Sprintf("%s-%d", fields.FieldAftStbdMethod, index)) {
		existing.AftStarboard.Method = types.ReadingMethod(d.str(r, fmt.Sprintf("%s-%d", fields.FieldAftStbdMethod, index)))
	}

	return existing
}

func (d *Decoder) ppKeelBlock(r *http.Request, existing *types.Draft, index int) {
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldDFwd, index)); ok {
		existing.DistancePPFwd = v
	}
	if r.PostForm.Has(fmt.Sprintf("%s-%d", fields.FieldDFwdDir, index)) {
		existing.PPFwdDirection = d.str(r, fmt.Sprintf("%s-%d", fields.FieldDFwdDir, index))
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldDMid, index)); ok {
		existing.DistancePPMid = v
	}
	if r.PostForm.Has(fmt.Sprintf("%s-%d", fields.FieldDMidDir, index)) {
		existing.PPMidDirection = d.str(r, fmt.Sprintf("%s-%d", fields.FieldDMidDir, index))
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldDAft, index)); ok {
		existing.DistancePPAft = v
	}
	if r.PostForm.Has(fmt.Sprintf("%s-%d", fields.FieldDAftDir, index)) {
		existing.PPAftDirection = d.str(r, fmt.Sprintf("%s-%d", fields.FieldDAftDir, index))
	}

	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldKeelFwd, index)); ok {
		existing.KeelFwd = v
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldKeelMid, index)); ok {
		existing.KeelMid = v
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldKeelAft, index)); ok {
		existing.KeelAft = v
	}
}

// hydrostaticsBlock parses the two HydrostaticRows.
// The frontend (web/widgets/drafts/hydrostatics.templ) uses distinct field
// constants for the upper (row 0) and lower (row 1) rows rather than a
// shared field name with a row-index suffix.
func (d *Decoder) hydrostaticsBlock(r *http.Request, existing *types.Draft, index int) {
	if len(existing.HydrostaticRows) < 2 {
		existing.HydrostaticRows = make([]types.HydrostaticRow, 2)
	}

	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldUDraft, index)); ok {
		existing.HydrostaticRows[0].Draft = v
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldUDisp, index)); ok {
		existing.HydrostaticRows[0].Displacement = v
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldUTPC, index)); ok {
		existing.HydrostaticRows[0].TPC = v
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldULCF, index)); ok {
		existing.HydrostaticRows[0].LCF = v
	}
	if r.PostForm.Has(fmt.Sprintf("%s-%d", fields.FieldULCFDir, index)) {
		existing.HydrostaticRows[0].LCFDirection = types.LCFDirection(d.str(r, fmt.Sprintf("%s-%d", fields.FieldULCFDir, index)))
	}

	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldLDraft, index)); ok {
		existing.HydrostaticRows[1].Draft = v
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldLDisp, index)); ok {
		existing.HydrostaticRows[1].Displacement = v
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldLTPC, index)); ok {
		existing.HydrostaticRows[1].TPC = v
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldLLCF, index)); ok {
		existing.HydrostaticRows[1].LCF = v
	}
	if r.PostForm.Has(fmt.Sprintf("%s-%d", fields.FieldLLCFDir, index)) {
		existing.HydrostaticRows[1].LCFDirection = types.LCFDirection(d.str(r, fmt.Sprintf("%s-%d", fields.FieldLLCFDir, index)))
	}
}

// mtcBlock parses the two MTCRows (+0.5m / -0.5m), matching the field
// constants used in web/widgets/drafts/hydrostatics.templ.
func (d *Decoder) mtcBlock(r *http.Request, existing *types.Draft, index int) {
	if len(existing.MTCRows) < 2 {
		existing.MTCRows = make([]types.MTCRow, 2)
	}

	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldPMTCDraft, index)); ok {
		existing.MTCRows[0].Draft = v
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldPMTC, index)); ok {
		existing.MTCRows[0].MTC = v
	}

	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldNMTCDraft, index)); ok {
		existing.MTCRows[1].Draft = v
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldNMTC, index)); ok {
		existing.MTCRows[1].MTC = v
	}
}

func (d *Decoder) deductiblesBlock(r *http.Request, existing types.Deductibles, index int) types.Deductibles {
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldHFO, index)); ok {
		existing.HFO = v
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldMDO, index)); ok {
		existing.MDO = v
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldLubOil, index)); ok {
		existing.LubOil = v
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldBilgeWater, index)); ok {
		existing.BilgeWater = v
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldSewageWater, index)); ok {
		existing.SewageWater = v
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldOthers, index)); ok {
		existing.Others = v
	}

	return existing
}

// seaConditionBlock parses the sea condition radio (Type) and the single
// condition select (web/components/sea_condition.templ posts one value under
// fields.FieldSeaCondition, whose meaning depends on the active Type).
func (d *Decoder) seaConditionBlock(r *http.Request, existing types.SeaCondition, index int) types.SeaCondition {
	typeKey := fmt.Sprintf("%s-%d", fields.FieldSeaType, index)
	if r.PostForm.Has(typeKey) {
		existing.Type = types.SeaConditionType(d.str(r, typeKey))
	}

	condKey := fmt.Sprintf("%s-%d", fields.FieldSeaCondition, index)
	if r.PostForm.Has(condKey) {
		val := d.str(r, condKey)
		switch existing.Type {
		case types.SeaConditionTypeWave:
			existing.Wave = types.WaveCondition(val)
		case types.SeaConditionTypeIce:
			existing.Ice = types.IceCondition(val)
		}
	}

	return existing
}

// draftMetaBlock parses fields that don't belong to a larger sub-struct.
// Density maps to fields.FieldDockwaterDensity — see
// web/widgets/drafts/deductibles.templ ("DockW Density" input).
func (d *Decoder) draftMetaBlock(r *http.Request, existing *types.Draft, index int) {
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldDockwaterDensity, index)); ok {
		existing.Density = v
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldTPCListPort, index)); ok {
		existing.TPCListPort = v
	}
	if v, ok := d.pFloat(r, fmt.Sprintf("%s-%d", fields.FieldTPCListStbd, index)); ok {
		existing.TPCListStarboard = v
	}
	if r.PostForm.Has(fmt.Sprintf("%s-%d", fields.FieldProductGrade, index)) {
		existing.ProductGrade = d.str(r, fmt.Sprintf("%s-%d", fields.FieldProductGrade, index))
	}
}
