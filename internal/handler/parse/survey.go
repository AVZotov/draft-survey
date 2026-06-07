package parse

// func SurveyBlock(r *http.Request, s *types.Survey) error {
// 	s.Job.JobNumber = String(r, fields.FieldJobNumber)
// 	//s.CargoOperation.Operation = String(r, fields.FieldCargoOperation)
// 	s.CargoOperation.Cargo = String(r, fields.FieldCargo)
// 	s.CargoOperation.Packing = String(r, fields.FieldPacking)
// 	s.Job.Principal = String(r, fields.FieldClient)
// 	s.Remarks = String(r, fields.FieldRemarks)
// 	s.Country.CountryCode = String(r, fields.FieldLoadCountryCode)
// 	s.Country.Name = String(r, fields.FieldLoadCountryName)

// 	if err := Int(r, fields.FieldDSNumber, &s.Job.DSNumber); err != nil {
// 		return err
// 	}
// 	if err := Float(r, fields.FieldCargoDeclared, &s.CargoDeclared); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func VesselBlock(r *http.Request, s *types.Survey) error {
// 	s.VesselData.Name = String(r, fields.FieldVesselName)
// 	s.VesselData.IMO = String(r, fields.FieldIMO)
// 	s.VesselData.Flag.CountryCode = String(r, fields.FieldFlagCountryCode)
// 	s.VesselData.Flag.Name = String(r, fields.FieldFlagCountryName)

// 	if err := Int(r, fields.FieldBuiltYear, &s.VesselData.BuiltYear); err != nil {
// 		return err
// 	}
// 	if err := Int(r, fields.FieldHoldsTotal, &s.VesselData.HoldsTotal); err != nil {
// 		return err
// 	}
// 	if err := Float(r, fields.FieldConstDeclared, &s.ConstantDeclared); err != nil {
// 		return err
// 	}
// 	if err := Float(r, fields.FieldTableDensity, &s.VesselData.TableDensity); err != nil {
// 		return err
// 	}
// 	if err := Float(r, fields.FieldLightship, &s.VesselData.Lightship); err != nil {
// 		return err
// 	}
// 	if err := Float(r, fields.FieldLBP, &s.VesselData.LBP); err != nil {
// 		return err
// 	}
// 	if err := Float(r, fields.FieldBreadth, &s.VesselData.Breadth); err != nil {
// 		return err
// 	}
// 	if err := Float(r, fields.FieldDepth, &s.VesselData.Depth); err != nil {
// 		return err
// 	}
// 	if err := Float(r, fields.FieldSummerDraft, &s.VesselData.SummerDraft); err != nil {
// 		return err
// 	}
// 	if err := Float(r, fields.FieldSummerDWT, &s.VesselData.SummerDWT); err != nil {
// 		return err
// 	}
// 	if err := Float(r, fields.FieldSummerTPC, &s.VesselData.SummerTPC); err != nil {
// 		return err
// 	}
// 	if err := Float(r, fields.FieldSummerFreeboard, &s.VesselData.SummerFreeboard); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func SettingsBlock(r *http.Request, s *types.Survey) {
// 	if method := String(r, fields.FieldMMCMethod); method != "" {
// 		s.VesselData.VesselType = types.VesselType(method)
// 	}
// 	if corr := String(r, fields.FieldCorrMethod); corr != "" {
// 		s.VesselData.CorrectionMethod = types.CorrectionMethod(corr)
// 	}
// 	if lcf := String(r, fields.FieldLCFDetection); lcf != "" {
// 		s.VesselData.IsLcfDetectionManual = lcf == "manual"
// 	}
// }
