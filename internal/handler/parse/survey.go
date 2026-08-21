package parse

import (
	"net/http"

	"github.com/AVZotov/draft-survey/internal/handler/fields"
	"github.com/AVZotov/draft-survey/internal/types"
)

func (d *Decoder) Survey(r *http.Request, existing *types.Survey) *types.Survey {
	survey := existing
	if survey == nil {
		survey = &types.Survey{}
	}

	if r.PostForm.Has(fields.FieldJobNumber) {
		survey.Job.JobNumber = d.str(r, fields.FieldJobNumber)
	}
	if v, ok := d.pInt(r, fields.FieldDSNumber); ok {
		survey.Job.DSNumber = v
	}
	if r.PostForm.Has(fields.FieldClient) {
		survey.Job.Principal = d.str(r, fields.FieldClient)
	}
	if r.PostForm.Has(fields.FieldCargo) {
		survey.CargoOperation.Cargo = d.str(r, fields.FieldCargo)
	}
	if r.PostForm.Has(fields.FieldCargoOperation) {
		survey.CargoOperation.Operation = types.Operation(d.str(r, fields.FieldCargoOperation))
	}
	if r.PostForm.Has(fields.FieldDescription) {
		survey.CargoOperation.Description = d.str(r, fields.FieldDescription)
	}
	if r.PostForm.Has(fields.FieldPacking) {
		survey.CargoOperation.Packing = d.str(r, fields.FieldPacking)
	}
	if r.PostForm.Has(fields.FieldRemarks) {
		survey.Remarks = d.str(r, fields.FieldRemarks)
	}
	if v, ok := d.pFloat(r, fields.FieldCargoDeclared); ok {
		survey.CargoDeclared = v
	}
	if r.PostForm.Has(fields.FieldCountryCode) {
		survey.Country.CountryCode = d.str(r, fields.FieldCountryCode)
	}

	if r.PostForm.Has(fields.FieldVesselName) {
		survey.VesselData.Name = d.str(r, fields.FieldVesselName)
	}
	if r.PostForm.Has(fields.FieldIMO) {
		survey.VesselData.IMO = d.str(r, fields.FieldIMO)
	}
	if v, ok := d.pInt(r, fields.FieldBuiltYear); ok {
		survey.VesselData.BuiltYear = v
	}
	if v, ok := d.pInt(r, fields.FieldHoldsTotal); ok {
		survey.VesselData.HoldsTotal = v
	}
	if v, ok := d.pFloat(r, fields.FieldConstDeclared); ok {
		survey.ConstantDeclared = v
	}
	if v, ok := d.pFloat(r, fields.FieldTableDensity); ok {
		survey.VesselData.TableDensity = v
	}
	if v, ok := d.pFloat(r, fields.FieldLightship); ok {
		survey.VesselData.Lightship = v
	}
	if v, ok := d.pFloat(r, fields.FieldLBP); ok {
		survey.VesselData.LBP = v
	}
	if v, ok := d.pFloat(r, fields.FieldBreadth); ok {
		survey.VesselData.Breadth = v
	}
	if v, ok := d.pFloat(r, fields.FieldDepth); ok {
		survey.VesselData.Depth = v
	}
	if v, ok := d.pFloat(r, fields.FieldSummerDraft); ok {
		survey.VesselData.SummerDraft = v
	}
	if v, ok := d.pFloat(r, fields.FieldSummerDWT); ok {
		survey.VesselData.SummerDWT = v
	}
	if v, ok := d.pFloat(r, fields.FieldSummerTPC); ok {
		survey.VesselData.SummerTPC = v
	}
	if v, ok := d.pFloat(r, fields.FieldSummerFreeboard); ok {
		survey.VesselData.SummerFreeboard = v
	}
	if r.PostForm.Has(fields.FieldFlagCountryCode) {
		survey.VesselData.Flag.CountryCode = d.str(r, fields.FieldFlagCountryCode)
	}

	if r.PostForm.Has(fields.FieldMMCMethod) {
		survey.VesselData.VesselType = types.VesselType(d.str(r, fields.FieldMMCMethod))
	}
	if r.PostForm.Has(fields.FieldCorrMethod) {
		survey.VesselData.CorrectionMethod = types.CorrectionMethod(d.str(r, fields.FieldCorrMethod))
	}
	if r.PostForm.Has(fields.FieldLCFDetection) {
		survey.VesselData.IsLcfDetectionManual = d.str(r, fields.FieldLCFDetection) == "manual"
	}

	return survey
}
