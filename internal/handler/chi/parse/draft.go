package parse

import (
	"net/http"
	"strconv"

	"github.com/AVZotov/draft-survey/internal/handler/chi/fields"
	"github.com/AVZotov/draft-survey/internal/types"
)

func DraftInfo(r *http.Request, d *types.Draft, index int) {
	d.LoadPort.PortManual = String(r, fields.WithIndex(fields.FieldLoadPort, index))
	d.DischargePort.PortManual = String(r, fields.WithIndex(fields.FieldDischargePort, index))
	d.ProductGrade = String(r, fields.WithIndex(fields.FieldProductGrade, index))
}

func SeaCondition(r *http.Request, d *types.Draft, index int) {
	seaType := String(r, fields.WithIndex(fields.FieldSeaType, index))
	if seaType == "" {
		return
	}
	d.SeaCondition.Type = types.SeaConditionType(seaType)

	condition := String(r, fields.WithIndex(fields.FieldSeaCondition, index))
	switch types.SeaConditionType(seaType) {
	case types.SeaConditionTypeWave:
		d.SeaCondition.Wave = types.WaveCondition(condition)
		d.SeaCondition.Ice = ""
	case types.SeaConditionTypeIce:
		d.SeaCondition.Ice = types.IceCondition(condition)
		d.SeaCondition.Wave = ""
	}
}

func HoldsBlock(r *http.Request, d *types.Draft, index int) {
	d.HoldsActive = []int{}
	values := r.PostForm[fields.WithIndex(fields.FieldHoldsActive, index)]
	for _, v := range values {
		i, err := strconv.Atoi(v)
		if err == nil {
			d.HoldsActive = append(d.HoldsActive, i)
		}
	}
}

func Marks(r *http.Request, d *types.Draft, index int) error {
	if err := Float(r, fields.WithIndex(fields.FieldFwdPort, index), &d.Marks.FwdPort.Value); err != nil {
		return err
	}
	if m := String(r, fields.WithIndex(fields.FieldFwdPortMethod, index)); m != "" {
		d.Marks.FwdPort.Method = types.ReadingMethod(m)
	}
	if err := Float(r, fields.WithIndex(fields.FieldFwdStbd, index), &d.Marks.FwdStarboard.Value); err != nil {
		return err
	}
	if m := String(r, fields.WithIndex(fields.FieldFwdStbdMethod, index)); m != "" {
		d.Marks.FwdStarboard.Method = types.ReadingMethod(m)
	}
	if err := Float(r, fields.WithIndex(fields.FieldMidPort, index), &d.Marks.MidPort.Value); err != nil {
		return err
	}
	if m := String(r, fields.WithIndex(fields.FieldMidPortMethod, index)); m != "" {
		d.Marks.MidPort.Method = types.ReadingMethod(m)
	}
	if err := Float(r, fields.WithIndex(fields.FieldMidStbd, index), &d.Marks.MidStarboard.Value); err != nil {
		return err
	}
	if m := String(r, fields.WithIndex(fields.FieldMidStbdMethod, index)); m != "" {
		d.Marks.MidStarboard.Method = types.ReadingMethod(m)
	}
	if err := Float(r, fields.WithIndex(fields.FieldAftPort, index), &d.Marks.AftPort.Value); err != nil {
		return err
	}
	if m := String(r, fields.WithIndex(fields.FieldAftPortMethod, index)); m != "" {
		d.Marks.AftPort.Method = types.ReadingMethod(m)
	}
	if err := Float(r, fields.WithIndex(fields.FieldAftStbd, index), &d.Marks.AftStarboard.Value); err != nil {
		return err
	}
	if m := String(r, fields.WithIndex(fields.FieldAftStbdMethod, index)); m != "" {
		d.Marks.AftStarboard.Method = types.ReadingMethod(m)
	}
	if err := Float(r, fields.WithIndex(fields.FieldTPCListPort, index), &d.TPCListPort); err != nil {
		return err
	}
	if err := Float(r, fields.WithIndex(fields.FieldTPCListStbd, index), &d.TPCListStarboard); err != nil {
		return err
	}
	return nil
}

func PPKeel(r *http.Request, d *types.Draft, index int) error {
	if err := Float(r, fields.WithIndex(fields.FieldDFwd, index), &d.DistancePPFwd); err != nil {
		return err
	}
	d.PPFwdDirection = String(r, fields.WithIndex(fields.FieldDFwdDir, index))
	if err := Float(r, fields.WithIndex(fields.FieldDMid, index), &d.DistancePPMid); err != nil {
		return err
	}
	d.PPMidDirection = String(r, fields.WithIndex(fields.FieldDMidDir, index))
	if err := Float(r, fields.WithIndex(fields.FieldDAft, index), &d.DistancePPAft); err != nil {
		return err
	}
	d.PPAftDirection = String(r, fields.WithIndex(fields.FieldDAftDir, index))
	if err := Float(r, fields.WithIndex(fields.FieldKeelFwd, index), &d.KeelFwd); err != nil {
		return err
	}
	if err := Float(r, fields.WithIndex(fields.FieldKeelMid, index), &d.KeelMid); err != nil {
		return err
	}
	if err := Float(r, fields.WithIndex(fields.FieldKeelAft, index), &d.KeelAft); err != nil {
		return err
	}
	return nil
}

func Hydrostatics(r *http.Request, d *types.Draft, index int) error {
	if err := Float(r, fields.WithIndex(fields.FieldUDraft, index), &d.HydrostaticRows[0].Draft); err != nil {
		return err
	}
	if err := Float(r, fields.WithIndex(fields.FieldUDisp, index), &d.HydrostaticRows[0].Displacement); err != nil {
		return err
	}
	if err := Float(r, fields.WithIndex(fields.FieldUTPC, index), &d.HydrostaticRows[0].TPC); err != nil {
		return err
	}
	if err := Float(r, fields.WithIndex(fields.FieldULCF, index), &d.HydrostaticRows[0].LCF); err != nil {
		return err
	}
	if dir := String(r, fields.WithIndex(fields.FieldULCFDir, index)); dir != "" {
		d.HydrostaticRows[0].LCFDirection = types.LCFDirection(dir)
	}
	if err := Float(r, fields.WithIndex(fields.FieldLDraft, index), &d.HydrostaticRows[1].Draft); err != nil {
		return err
	}
	if err := Float(r, fields.WithIndex(fields.FieldLDisp, index), &d.HydrostaticRows[1].Displacement); err != nil {
		return err
	}
	if err := Float(r, fields.WithIndex(fields.FieldLTPC, index), &d.HydrostaticRows[1].TPC); err != nil {
		return err
	}
	if err := Float(r, fields.WithIndex(fields.FieldLLCF, index), &d.HydrostaticRows[1].LCF); err != nil {
		return err
	}
	if dir := String(r, fields.WithIndex(fields.FieldLLCFDir, index)); dir != "" {
		d.HydrostaticRows[1].LCFDirection = types.LCFDirection(dir)
	}
	if err := Float(r, fields.WithIndex(fields.FieldPMTCDraft, index), &d.MTCRows[0].Draft); err != nil {
		return err
	}
	if err := Float(r, fields.WithIndex(fields.FieldPMTC, index), &d.MTCRows[0].MTC); err != nil {
		return err
	}
	if err := Float(r, fields.WithIndex(fields.FieldNMTCDraft, index), &d.MTCRows[1].Draft); err != nil {
		return err
	}
	if err := Float(r, fields.WithIndex(fields.FieldNMTC, index), &d.MTCRows[1].MTC); err != nil {
		return err
	}
	return nil
}

func Deductibles(r *http.Request, d *types.Draft, index int) error {
	if err := Float(r, fields.WithIndex(fields.FieldHFO, index), &d.Deductibles.HFO); err != nil {
		return err
	}
	if err := Float(r, fields.WithIndex(fields.FieldMDO, index), &d.Deductibles.MDO); err != nil {
		return err
	}
	if err := Float(r, fields.WithIndex(fields.FieldLubOil, index), &d.Deductibles.LubOil); err != nil {
		return err
	}
	if err := Float(r, fields.WithIndex(fields.FieldBilgeWater, index), &d.Deductibles.BilgeWater); err != nil {
		return err
	}
	if err := Float(r, fields.WithIndex(fields.FieldDockwaterDensity, index), &d.Density); err != nil {
		return err
	}
	if err := Float(r, fields.WithIndex(fields.FieldOthers, index), &d.Deductibles.Others); err != nil {
		return err
	}
	return nil
}
