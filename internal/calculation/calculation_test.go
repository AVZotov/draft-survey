package calculation

import (
	"testing"

	"github.com/AVZotov/draft-survey/internal/types"
)

func fp(v float64) *float64 { return &v }

func getFreshWaterTank() types.Tank {
	return types.Tank{
		Name:      "test",
		Sounding:  fp(3.5),
		Volume:    fp(3.5),
		IsFWTTank: true,
	}
}

func getInitFreshWaterTanks() []types.Tank {
	return []types.Tank{
		{Name: "FW P", Sounding: fp(364.000), Volume: fp(364.000), IsFWTTank: true},
	}
}

func getBallastWaterTank() types.Tank {
	return types.Tank{
		Name:     "test",
		Sounding: fp(3.5),
		Volume:   fp(3.5),
		Density:  fp(1.025),
	}
}

func getInitBallastWaterTanks() []types.Tank {
	return []types.Tank{
		{Name: "FPT", Sounding: fp(10347.899), Volume: fp(10347.899), Density: fp(1.025)},
	}
}

func getInitDraftData() types.Draft {
	return types.Draft{
		TPCListPort:      fp(49.665),
		TPCListStarboard: fp(49.688),
		Density:          fp(1.023),
	}
}

func getMarks() types.Marks {
	return types.Marks{
		FwdPort:      types.Mark{Value: fp(3.41)},
		FwdStarboard: types.Mark{Value: fp(3.41)},
		MidPort:      types.Mark{Value: fp(4.51)},
		MidStarboard: types.Mark{Value: fp(4.54)},
		AftPort:      types.Mark{Value: fp(5.69)},
		AftStarboard: types.Mark{Value: fp(5.70)},
	}
}

func getVesselData() types.VesselData {
	return types.VesselData{
		LBP:        182.000,
		VesselType: types.VesselTypeMarine,
		Lightship:  8390.000,
	}
}

func getDraft() types.Draft {
	return types.Draft{
		DistancePPFwd:  fp(1.400),
		PPFwdDirection: "A",
		DistancePPMid:  fp(0.400),
		PPMidDirection: "A",
		DistancePPAft:  fp(9.950),
		PPAftDirection: "F",
		KeelFwd:        fp(0.000),
		KeelMid:        fp(0.000),
		KeelAft:        fp(0.000),
	}
}

func getInitHydrostaticRows() []types.HydrostaticRow {
	return []types.HydrostaticRow{
		{Draft: fp(4.54), Displacement: fp(21226), TPC: fp(49.7), LCF: fp(6.93), LCFDirection: types.LCFDirectionForward},
		{Draft: fp(4.55), Displacement: fp(21276), TPC: fp(49.7), LCF: fp(6.92), LCFDirection: types.LCFDirectionForward},
	}
}

func getInitMtcRows() []types.MTCRow {
	return []types.MTCRow{
		{Draft: fp(4.04), MTC: fp(529.4)},
		{Draft: fp(5.04), MTC: fp(548.0)},
	}
}

func getInitDeductibles() types.Deductibles {
	return types.Deductibles{
		HFO: fp(683.868),
		MDO: fp(89.130),
	}
}

// calcTestListCorrection is a helper that builds CalcListCorrection call
// for tests using getInitDraftData() with manually entered TPC list values.
func calcTestListCorrection(marks types.Marks, mmc float64, hydrostatics types.Hydrostatics) float64 {
	initDS := getInitDraftData()
	hr := getInitHydrostaticRows()
	vesselData := getVesselData()
	return CalcListCorrection(marks, initDS.TPCListPort, initDS.TPCListStarboard, hr, mmc, hydrostatics, vesselData)
}

func TestFreshWaterTank_GetWeight(t *testing.T) {
	const weight = 3.5
	tank := getFreshWaterTank()
	tankWeight := tank.CalcWeight()
	if tankWeight != weight {
		t.Errorf("Expected %f, got %f", weight, tankWeight)
	}
}

func TestBallastWaterTank_GetWeight(t *testing.T) {
	const weight = 3.587
	tank := getBallastWaterTank()
	tankWeight := round3(tank.CalcWeight())
	if tankWeight != weight {
		t.Errorf("Expected %f, got %f", weight, tankWeight)
	}
}

func TestTotalFreshWater(t *testing.T) {
	const weight = 364.000
	tanks := getInitFreshWaterTanks()
	totalWeight := TotalTanksWeight(tanks)
	if totalWeight != weight {
		t.Errorf("Expected %f, got %f", weight, totalWeight)
	}
}

func TestTotalBallastWater(t *testing.T) {
	const weight = 17.935
	tank := getBallastWaterTank()
	var ballastWaterTanks []types.Tank
	for i := 0; i < 5; i++ {
		ballastWaterTanks = append(ballastWaterTanks, tank)
	}
	totalWeight := round3(TotalTanksWeight(ballastWaterTanks))
	if totalWeight != weight {
		t.Errorf("Expected %f, got %f", weight, totalWeight)
	}
}

func TestMeanDrafts(t *testing.T) {
	draftFWDExpected := 3.410
	draftMIDExpected := 4.525
	draftAFTExpected := 5.695

	marks := getMarks()
	meanDrafts := MeanDrafts(marks)

	if meanDrafts.DraftFwdMean != draftFWDExpected {
		t.Errorf("Expected %f, got %f", draftFWDExpected, meanDrafts.DraftFwdMean)
	}
	if meanDrafts.DraftMidMean != draftMIDExpected {
		t.Errorf("Expected %f, got %f", draftMIDExpected, meanDrafts.DraftMidMean)
	}
	if meanDrafts.DraftAftMean != draftAFTExpected {
		t.Errorf("Expected %f, got %f", draftAFTExpected, meanDrafts.DraftAftMean)
	}
}

func TestCalcFullLBPPPCorrections(t *testing.T) {
	fwdCorrectionExpected := -0.019
	midCorrectionExpected := -0.005
	aftCorrectionExpected := 0.133

	meanDrafts := MeanDrafts(getMarks())
	vesselData := getVesselData()
	ppCorrections := CalcFullLBPPPCorrections(meanDrafts, getDraft(), vesselData.LBP)

	if fwdCorrectionExpected != ppCorrections.FwdCorrection {
		t.Errorf("Expected %f, got %f", fwdCorrectionExpected, ppCorrections.FwdCorrection)
	}
	if midCorrectionExpected != ppCorrections.MidCorrection {
		t.Errorf("Expected %f, got %f", midCorrectionExpected, ppCorrections.MidCorrection)
	}
	if aftCorrectionExpected != ppCorrections.AftCorrection {
		t.Errorf("Expected %f, got %f", aftCorrectionExpected, ppCorrections.AftCorrection)
	}
}

func TestCalcHalfLBPPPCorrections(t *testing.T) {
	fwdCorrectionExpected := -0.017
	midCorrectionExpected := -0.005
	aftCorrectionExpected := 0.144

	meanDrafts := MeanDrafts(getMarks())
	vesselData := getVesselData()
	ppCorrections := CalcHalfLBPPPCorrections(meanDrafts, getDraft(), vesselData.LBP)

	if fwdCorrectionExpected != ppCorrections.FwdCorrection {
		t.Errorf("Expected %f, got %f", fwdCorrectionExpected, ppCorrections.FwdCorrection)
	}
	if midCorrectionExpected != ppCorrections.MidCorrection {
		t.Errorf("Expected %f, got %f", midCorrectionExpected, ppCorrections.MidCorrection)
	}
	if aftCorrectionExpected != ppCorrections.AftCorrection {
		t.Errorf("Expected %f, got %f", aftCorrectionExpected, ppCorrections.AftCorrection)
	}
}

func TestCalcDraftsWKeel(t *testing.T) {
	FWDDraftsWKeelExpected := 3.391
	MIDDraftsWKeelExpected := 4.520
	AFTDraftsWKeelExpected := 5.828

	marks := getMarks()
	meanDraft := MeanDrafts(marks)
	vesselData := getVesselData()
	ppCorrections := CalcFullLBPPPCorrections(meanDraft, getDraft(), vesselData.LBP)
	draftsWKeel := CalcDraftsWKeel(meanDraft, ppCorrections, getDraft())

	if FWDDraftsWKeelExpected != draftsWKeel.FwdDraftWKeel {
		t.Errorf("Expected %f, got %f", FWDDraftsWKeelExpected, draftsWKeel.FwdDraftWKeel)
	}
	if MIDDraftsWKeelExpected != draftsWKeel.MidDraftWKeel {
		t.Errorf("Expected %f, got %f", MIDDraftsWKeelExpected, draftsWKeel.MidDraftWKeel)
	}
	if AFTDraftsWKeelExpected != draftsWKeel.AftDraftWKeel {
		t.Errorf("Expected %f, got %f", AFTDraftsWKeelExpected, draftsWKeel.AftDraftWKeel)
	}
}

func TestCalcMMC(t *testing.T) {
	mmcExpected := 4.542
	marks := getMarks()
	meanDraft := MeanDrafts(marks)
	vesselData := getVesselData()
	ppCorrections := CalcFullLBPPPCorrections(meanDraft, getDraft(), vesselData.LBP)
	draftsWKeel := CalcDraftsWKeel(meanDraft, ppCorrections, getDraft())
	mmc := CalcMMC(draftsWKeel, vesselData)

	if mmcExpected != mmc {
		t.Errorf("Expected %f, got %f", mmcExpected, mmc)
	}
}

func TestInterpolate(t *testing.T) {
	expected := 21236.000
	got := Interpolate(4.542, 4.540, 21226.000, 4.550, 21276.000)

	if expected != got {
		t.Errorf("Expected %f, got %f", expected, got)
	}
}

func TestCalcHydrostatics(t *testing.T) {
	displacementExpected := 21236.000
	tpcExpected := 49.700
	lcfExpected := -6.928
	marks := getMarks()
	meanDraft := MeanDrafts(marks)
	vesselData := getVesselData()
	ppCorrections := CalcFullLBPPPCorrections(meanDraft, getDraft(), vesselData.LBP)
	draftsWKeel := CalcDraftsWKeel(meanDraft, ppCorrections, getDraft())
	mmc := CalcMMC(draftsWKeel, vesselData)
	hr := getInitHydrostaticRows()
	hydrostatics := CalcHydrostatics(mmc, hr, vesselData)

	if displacementExpected != hydrostatics.Displacement {
		t.Errorf("Expected %f, got %f", displacementExpected, hydrostatics.Displacement)
	}
	if tpcExpected != hydrostatics.TPC {
		t.Errorf("Expected %f, got %f", tpcExpected, hydrostatics.TPC)
	}
	if lcfExpected != hydrostatics.LCF {
		t.Errorf("Expected %f, got %f", lcfExpected, hydrostatics.LCF)
	}
}

func TestCalcFirstTrimCorrection(t *testing.T) {
	firstTrimCorrectionExpected := -461.050
	marks := getMarks()
	meanDraft := MeanDrafts(marks)
	vesselData := getVesselData()
	ppCorrections := CalcFullLBPPPCorrections(meanDraft, getDraft(), vesselData.LBP)
	draftsWKeel := CalcDraftsWKeel(meanDraft, ppCorrections, getDraft())
	mmc := CalcMMC(draftsWKeel, vesselData)
	hr := getInitHydrostaticRows()
	hydrostatics := CalcHydrostatics(mmc, hr, vesselData)
	firstTrimCorrectionGot := CalcFirstTrimCorrection(draftsWKeel, hydrostatics.TPC, hydrostatics.LCF, vesselData.LBP)

	if firstTrimCorrectionExpected != firstTrimCorrectionGot {
		t.Errorf("Expected %f, got %f", firstTrimCorrectionExpected, firstTrimCorrectionGot)
	}
}

func TestCalcSecondTrimCorrection(t *testing.T) {
	secondTrimCorrectionExpected := 30.347
	marks := getMarks()
	meanDraft := MeanDrafts(marks)
	vesselData := getVesselData()
	ppCorrections := CalcFullLBPPPCorrections(meanDraft, getDraft(), vesselData.LBP)
	draftsWKeel := CalcDraftsWKeel(meanDraft, ppCorrections, getDraft())
	mtcRows := getInitMtcRows()
	secondTrimCorrectionGot := CalcSecondTrimCorrection(draftsWKeel, mtcRows, vesselData.LBP)

	if secondTrimCorrectionExpected != secondTrimCorrectionGot {
		t.Errorf("Expected %f, got %f", secondTrimCorrectionExpected, secondTrimCorrectionGot)
	}
}

func TestCalcListCorrection(t *testing.T) {
	listCorrectionExpected := 0.004
	marks := getMarks()
	meanDraft := MeanDrafts(marks)
	vesselData := getVesselData()
	ppCorrections := CalcFullLBPPPCorrections(meanDraft, getDraft(), vesselData.LBP)
	draftsWKeel := CalcDraftsWKeel(meanDraft, ppCorrections, getDraft())
	mmc := CalcMMC(draftsWKeel, vesselData)
	hr := getInitHydrostaticRows()
	hydrostatics := CalcHydrostatics(mmc, hr, vesselData)
	initDS := getInitDraftData()

	listCorrectionGot := CalcListCorrection(
		marks,
		initDS.TPCListPort,
		initDS.TPCListStarboard,
		hr,
		mmc,
		hydrostatics,
		vesselData,
	)
	if listCorrectionExpected != listCorrectionGot {
		t.Errorf("Expected %f, got %f", listCorrectionExpected, listCorrectionGot)
	}
}

func TestCalcDensityCorrection(t *testing.T) {
	densityCorrExpected := -40.596
	marks := getMarks()
	meanDraft := MeanDrafts(marks)
	vesselData := getVesselData()
	ppCorrections := CalcFullLBPPPCorrections(meanDraft, getDraft(), vesselData.LBP)
	draftsWKeel := CalcDraftsWKeel(meanDraft, ppCorrections, getDraft())
	mmc := CalcMMC(draftsWKeel, vesselData)
	hr := getInitHydrostaticRows()
	hydrostatics := CalcHydrostatics(mmc, hr, vesselData)
	mtcRows := getInitMtcRows()
	initDS := getInitDraftData()
	firstTrim := CalcFirstTrimCorrection(draftsWKeel, hydrostatics.TPC, hydrostatics.LCF, vesselData.LBP)
	secondTrim := CalcSecondTrimCorrection(draftsWKeel, mtcRows, vesselData.LBP)
	listCorrection := CalcListCorrection(marks, initDS.TPCListPort, initDS.TPCListStarboard, hr, mmc, hydrostatics, vesselData)
	densityCorrGot := CalcDensityCorrection(hydrostatics.Displacement, firstTrim, secondTrim, listCorrection, *initDS.Density, 1.025)

	if densityCorrExpected != densityCorrGot {
		t.Errorf("Expected %f, got %f", densityCorrExpected, densityCorrGot)
	}
}

func TestCalcTotalDeductibles(t *testing.T) {
	totalDeductiblesExpected := 11743.594
	bwt := getInitBallastWaterTanks()
	fwt := getInitFreshWaterTanks()
	d := getInitDeductibles()
	totalDeductiblesGot := CalcTotalDeductibles(bwt, fwt, d)

	if totalDeductiblesExpected != totalDeductiblesGot {
		t.Errorf("Expected %f, got %f", totalDeductiblesExpected, totalDeductiblesGot)
	}
}

func TestCalcNetDisplacement(t *testing.T) {
	netDisplacementExpected := 9021.111
	bwt := getInitBallastWaterTanks()
	fwt := getInitFreshWaterTanks()
	d := getInitDeductibles()
	totalDeductibles := CalcTotalDeductibles(bwt, fwt, d)
	marks := getMarks()
	meanDraft := MeanDrafts(marks)
	vesselData := getVesselData()
	ppCorrections := CalcFullLBPPPCorrections(meanDraft, getDraft(), vesselData.LBP)
	draftsWKeel := CalcDraftsWKeel(meanDraft, ppCorrections, getDraft())
	mmc := CalcMMC(draftsWKeel, vesselData)
	hr := getInitHydrostaticRows()
	hydrostatics := CalcHydrostatics(mmc, hr, vesselData)
	mtcRows := getInitMtcRows()
	initDS := getInitDraftData()
	firstTrim := CalcFirstTrimCorrection(draftsWKeel, hydrostatics.TPC, hydrostatics.LCF, vesselData.LBP)
	secondTrim := CalcSecondTrimCorrection(draftsWKeel, mtcRows, vesselData.LBP)
	listCorrection := CalcListCorrection(marks, initDS.TPCListPort, initDS.TPCListStarboard, hr, mmc, hydrostatics, vesselData)
	densityCorr := CalcDensityCorrection(hydrostatics.Displacement, firstTrim, secondTrim, listCorrection, *initDS.Density, 1.025)
	netDisplacementGot := CalcNetDisplacement(hydrostatics.Displacement, firstTrim, secondTrim, listCorrection, densityCorr, totalDeductibles)

	if netDisplacementExpected != netDisplacementGot {
		t.Errorf("Expected %f, got %f", netDisplacementExpected, netDisplacementGot)
	}
}

func TestCalcCargoWeight(t *testing.T) {
	netDisplacementIni := 9000.000
	netDisplacementFin := 49000.000
	cargoWeightExpected := 40000.000
	cargoWeightGot := CalcCargoWeight(netDisplacementIni, netDisplacementFin)
	if cargoWeightExpected != cargoWeightGot {
		t.Errorf("Expected %f, got %f", cargoWeightExpected, cargoWeightGot)
	}
}

func TestCalcConstant(t *testing.T) {
	constantExpected := 631.111
	bwt := getInitBallastWaterTanks()
	fwt := getInitFreshWaterTanks()
	d := getInitDeductibles()
	totalDeductibles := CalcTotalDeductibles(bwt, fwt, d)
	marks := getMarks()
	meanDraft := MeanDrafts(marks)
	vesselData := getVesselData()
	ppCorrections := CalcFullLBPPPCorrections(meanDraft, getDraft(), vesselData.LBP)
	draftsWKeel := CalcDraftsWKeel(meanDraft, ppCorrections, getDraft())
	mmc := CalcMMC(draftsWKeel, vesselData)
	hr := getInitHydrostaticRows()
	hydrostatics := CalcHydrostatics(mmc, hr, vesselData)
	mtcRows := getInitMtcRows()
	initDS := getInitDraftData()
	firstTrim := CalcFirstTrimCorrection(draftsWKeel, hydrostatics.TPC, hydrostatics.LCF, vesselData.LBP)
	secondTrim := CalcSecondTrimCorrection(draftsWKeel, mtcRows, vesselData.LBP)
	listCorrection := CalcListCorrection(marks, initDS.TPCListPort, initDS.TPCListStarboard, hr, mmc, hydrostatics, vesselData)
	densityCorr := CalcDensityCorrection(hydrostatics.Displacement, firstTrim, secondTrim, listCorrection, *initDS.Density, 1.025)
	netDisplacement := CalcNetDisplacement(hydrostatics.Displacement, firstTrim, secondTrim, listCorrection, densityCorr, totalDeductibles)
	constantGot := CalcConstant(netDisplacement, vesselData.Lightship)

	if constantExpected != constantGot {
		t.Errorf("Expected %f, got %f", constantExpected, constantGot)
	}
}

func TestCalcCurrentDWT(t *testing.T) {
	constantExpected := 12374.705
	marks := getMarks()
	meanDraft := MeanDrafts(marks)
	vesselData := getVesselData()
	ppCorrections := CalcFullLBPPPCorrections(meanDraft, getDraft(), vesselData.LBP)
	draftsWKeel := CalcDraftsWKeel(meanDraft, ppCorrections, getDraft())
	mmc := CalcMMC(draftsWKeel, vesselData)
	hr := getInitHydrostaticRows()
	hydrostatics := CalcHydrostatics(mmc, hr, vesselData)
	mtcRows := getInitMtcRows()
	initDS := getInitDraftData()
	firstTrim := CalcFirstTrimCorrection(draftsWKeel, hydrostatics.TPC, hydrostatics.LCF, vesselData.LBP)
	secondTrim := CalcSecondTrimCorrection(draftsWKeel, mtcRows, vesselData.LBP)
	listCorrection := CalcListCorrection(marks, initDS.TPCListPort, initDS.TPCListStarboard, hr, mmc, hydrostatics, vesselData)
	densityCorr := CalcDensityCorrection(hydrostatics.Displacement, firstTrim, secondTrim, listCorrection, *initDS.Density, 1.025)
	displCorrToDensity := round3(hydrostatics.Displacement + firstTrim + secondTrim + listCorrection + densityCorr)
	DWTGot := CalcCurrentDWT(displCorrToDensity, vesselData.Lightship)

	if constantExpected != DWTGot {
		t.Errorf("Expected %f, got %f", constantExpected, DWTGot)
	}
}
