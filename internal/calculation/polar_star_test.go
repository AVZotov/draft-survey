package calculation

import (
	"testing"

	"github.com/AVZotov/draft-survey/internal/types"
)

// POLAR STAR - TrimnoList tests

func getPolarStarTrimNoListVessel() types.VesselData {
	return types.VesselData{
		LBP:        fp(183.000),
		VesselType: types.VesselTypeMarine,
	}
}

func getPolarStarDraft() types.Draft {
	return types.Draft{
		DistancePPFwd:  fp(4.800),
		PPFwdDirection: "A",
		DistancePPMid:  fp(0.500),
		PPMidDirection: "A",
		DistancePPAft:  fp(1.200),
		PPAftDirection: "A",
		KeelFwd:        fp(0),
		KeelMid:        fp(0),
		KeelAft:        fp(0),
	}
}

func getPolarStarTrimNoListMarks() types.Marks {
	return types.Marks{
		FwdPort: types.Mark{Value: fp(3.33)}, FwdStarboard: types.Mark{Value: fp(3.33)},
		MidPort: types.Mark{Value: fp(4.64)}, MidStarboard: types.Mark{Value: fp(4.64)},
		AftPort: types.Mark{Value: fp(6.12)}, AftStarboard: types.Mark{Value: fp(6.12)},
	}
}

func getPolarStarTrimNoListHydrostaticRows() []types.HydrostaticRow {
	return []types.HydrostaticRow{
		{Draft: fp(4.617), Displacement: fp(19182.7), TPC: fp(45.2), LCF: fp(98.457), LCFDirection: types.LCFDirectionFromAP},
		{Draft: fp(4.667), Displacement: fp(19409.0), TPC: fp(45.3), LCF: fp(98.405), LCFDirection: types.LCFDirectionFromAP},
	}
}

func getPolarStarTrimNoListMTCRows() []types.MTCRow {
	return []types.MTCRow{
		{Draft: fp(4.167), MTC: fp(500.2)},
		{Draft: fp(5.167), MTC: fp(526.9)},
	}
}

func TestPolarStar_TrimNoList_MeanDrafts(t *testing.T) {
	marks := getPolarStarTrimNoListMarks()
	got := MeanDrafts(marks)

	if got.DraftFwdMean != 3.330 {
		t.Errorf("meanF: expected 3.330, got %f", got.DraftFwdMean)
	}
	if got.DraftMidMean != 4.640 {
		t.Errorf("meanM: expected 4.640, got %f", got.DraftMidMean)
	}
	if got.DraftAftMean != 6.120 {
		t.Errorf("meanA: expected 6.120, got %f", got.DraftAftMean)
	}
}

func TestPolarStar_TrimnoList_PPCorrections(t *testing.T) {
	marks := getPolarStarTrimNoListMarks()
	vesselData := getPolarStarTrimNoListVessel()
	got := CalcFullLBPPPCorrections(MeanDrafts(marks), getPolarStarDraft(), *vesselData.LBP)

	if got.FwdCorrection != -0.075 {
		t.Errorf("FWD corr: expected -0.075, got %f", got.FwdCorrection)
	}
	if got.MidCorrection != -0.008 {
		t.Errorf("MID corr: expected -0.008, got %f", got.MidCorrection)
	}
	if got.AftCorrection != -0.019 {
		t.Errorf("AFT corr: expected -0.019, got %f", got.AftCorrection)
	}
}

func TestPolarStar_TrimNoList_DraftsWKeel(t *testing.T) {
	marks := getPolarStarTrimNoListMarks()
	vesselData := getPolarStarTrimNoListVessel()
	meanDraft := MeanDrafts(marks)
	ppCorrections := CalcFullLBPPPCorrections(meanDraft, getPolarStarDraft(), *vesselData.LBP)
	got := CalcDraftsWKeel(meanDraft, ppCorrections, getPolarStarDraft())

	if got.FwdDraftWKeel != 3.255 {
		t.Errorf("FWD wKeel: expected 3.255, got %f", got.FwdDraftWKeel)
	}
	if got.MidDraftWKeel != 4.632 {
		t.Errorf("MID wKeel: expected 4.632, got %f", got.MidDraftWKeel)
	}
	if got.AftDraftWKeel != 6.101 {
		t.Errorf("AFT wKeel: expected 6.101, got %f", got.AftDraftWKeel)
	}
}

func TestPolarStar_TrimNoList_MMC(t *testing.T) {
	marks := getPolarStarTrimNoListMarks()
	vesselData := getPolarStarTrimNoListVessel()
	meanDraft := MeanDrafts(marks)
	ppCorrections := CalcFullLBPPPCorrections(meanDraft, getPolarStarDraft(), *vesselData.LBP)
	draftsWKeel := CalcDraftsWKeel(meanDraft, ppCorrections, getPolarStarDraft())
	got := CalcMMC(draftsWKeel, vesselData)

	if got != 4.644 {
		t.Errorf("MMC: expected 4.644, got %f", got)
	}
}

func TestPolarStar_TrimNoListCalcHydrostatics(t *testing.T) {
	displacementExpected := 19304.902
	tpcExpected := 45.254
	lcfExpected := -6.929
	marks := getPolarStarTrimNoListMarks()
	vesselData := getPolarStarTrimNoListVessel()
	meanDraft := MeanDrafts(marks)
	ppCorrections := CalcFullLBPPPCorrections(meanDraft, getPolarStarDraft(), *vesselData.LBP)
	draftsWKeel := CalcDraftsWKeel(meanDraft, ppCorrections, getPolarStarDraft())
	mmc := CalcMMC(draftsWKeel, vesselData)
	hr := getPolarStarTrimNoListHydrostaticRows()

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

func TestPolarStar_TrimNoList_FirstTrimCorrection(t *testing.T) {
	marks := getPolarStarTrimNoListMarks()
	vesselData := getPolarStarTrimNoListVessel()
	meanDraft := MeanDrafts(marks)
	ppCorrections := CalcFullLBPPPCorrections(meanDraft, getPolarStarDraft(), *vesselData.LBP)
	draftsWKeel := CalcDraftsWKeel(meanDraft, ppCorrections, getPolarStarDraft())
	mmc := CalcMMC(draftsWKeel, vesselData)
	hr := getPolarStarTrimNoListHydrostaticRows()
	hydrostatics := CalcHydrostatics(mmc, hr, vesselData)
	got := CalcFirstTrimCorrection(draftsWKeel, hydrostatics.TPC, hydrostatics.LCF, *vesselData.LBP)

	if got != -487.653 {
		t.Errorf("1st trim: expected -487.653, got %f", got)
	}
}

func TestPolarStar_TrimNoList_SecondTrimCorrection(t *testing.T) {
	marks := getPolarStarTrimNoListMarks()
	vesselData := getPolarStarTrimNoListVessel()
	meanDraft := MeanDrafts(marks)
	ppCorrections := CalcFullLBPPPCorrections(meanDraft, getPolarStarDraft(), *vesselData.LBP)
	draftsWKeel := CalcDraftsWKeel(meanDraft, ppCorrections, getPolarStarDraft())
	mtcRows := getPolarStarTrimNoListMTCRows()
	got := CalcSecondTrimCorrection(draftsWKeel, mtcRows, *vesselData.LBP)

	if got != 59.088 {
		t.Errorf("2nd trim: expected 59.088, got %f", got)
	}
}

func TestPolarStar_TrimNoList_ListCorrection(t *testing.T) {
	marks := getPolarStarTrimNoListMarks()
	vesselData := getPolarStarTrimNoListVessel()
	hr := getPolarStarTrimNoListHydrostaticRows()
	mmc := 4.644
	hydrostatics := CalcHydrostatics(mmc, hr, vesselData)
	got := CalcListCorrection(marks, nil, nil, hr, mmc, hydrostatics, vesselData)
	if got != 0 {
		t.Errorf("List corr: expected 0, got %f", got)
	}
}

func TestPolarStar_TrimNoList_DensityCorrection(t *testing.T) {
	marks := getPolarStarTrimNoListMarks()
	vesselData := getPolarStarTrimNoListVessel()
	meanDraft := MeanDrafts(marks)
	ppCorrections := CalcFullLBPPPCorrections(meanDraft, getPolarStarDraft(), *vesselData.LBP)
	draftsWKeel := CalcDraftsWKeel(meanDraft, ppCorrections, getPolarStarDraft())
	mmc := CalcMMC(draftsWKeel, vesselData)
	hr := getPolarStarTrimNoListHydrostaticRows()
	hydrostatics := CalcHydrostatics(mmc, hr, vesselData)
	mtcRows := getPolarStarTrimNoListMTCRows()
	firstTrim := CalcFirstTrimCorrection(draftsWKeel, hydrostatics.TPC, hydrostatics.LCF, *vesselData.LBP)
	secondTrim := CalcSecondTrimCorrection(draftsWKeel, mtcRows, *vesselData.LBP)
	listCorrection := CalcListCorrection(marks, nil, nil, hr, mmc, hydrostatics, vesselData)
	got := CalcDensityCorrection(hydrostatics.Displacement, firstTrim, secondTrim, listCorrection, 1.017, 1.025)
	if got != -147.328 {
		t.Errorf("Density corr: expected -147.328, got %f", got)
	}
}

// POLAR STAR - TrimList tests

func getPolarStarTrimListVessel() types.VesselData {
	return types.VesselData{
		LBP:        fp(183.000),
		VesselType: types.VesselTypeMarine,
	}
}

func getPolarStarTrimListMarks() types.Marks {
	return types.Marks{
		FwdPort: types.Mark{Value: fp(3.39)}, FwdStarboard: types.Mark{Value: fp(3.36)},
		MidPort: types.Mark{Value: fp(4.64)}, MidStarboard: types.Mark{Value: fp(4.54)},
		AftPort: types.Mark{Value: fp(6.12)}, AftStarboard: types.Mark{Value: fp(6.12)},
	}
}

func getPolarStarTrimListHydrostaticRows() []types.HydrostaticRow {
	return []types.HydrostaticRow{
		{Draft: fp(4.567), Displacement: fp(18956.7), TPC: fp(45.2), LCF: fp(98.509), LCFDirection: types.LCFDirectionFromAP},
		{Draft: fp(4.617), Displacement: fp(19182.7), TPC: fp(45.2), LCF: fp(98.457), LCFDirection: types.LCFDirectionFromAP},
	}
}

func getPolarStarTrimListMTCRows() []types.MTCRow {
	return []types.MTCRow{
		{Draft: fp(4.117), MTC: fp(498.8)},
		{Draft: fp(5.117), MTC: fp(525.7)},
	}
}

func TestPolarStar_TrimList_MeanDrafts(t *testing.T) {
	got := MeanDrafts(getPolarStarTrimListMarks())
	if got.DraftFwdMean != 3.375 {
		t.Errorf("meanF: expected 3.375, got %f", got.DraftFwdMean)
	}
	if got.DraftMidMean != 4.590 {
		t.Errorf("meanM: expected 4.590, got %f", got.DraftMidMean)
	}
	if got.DraftAftMean != 6.120 {
		t.Errorf("meanA: expected 6.120, got %f", got.DraftAftMean)
	}
}

func TestPolarStar_TrimList_PPCorrections(t *testing.T) {
	vesselData := getPolarStarTrimListVessel()
	got := CalcFullLBPPPCorrections(MeanDrafts(getPolarStarTrimListMarks()), getPolarStarDraft(), *vesselData.LBP)
	if got.FwdCorrection != -0.073 {
		t.Errorf("FWD: expected -0.073, got %f", got.FwdCorrection)
	}
	if got.MidCorrection != -0.008 {
		t.Errorf("MID: expected -0.008, got %f", got.MidCorrection)
	}
	if got.AftCorrection != -0.018 {
		t.Errorf("AFT: expected -0.018, got %f", got.AftCorrection)
	}
}

func TestPolarStar_TrimList_DraftsWKeel(t *testing.T) {
	vesselData := getPolarStarTrimListVessel()
	meanDraft := MeanDrafts(getPolarStarTrimListMarks())
	got := CalcDraftsWKeel(meanDraft, CalcFullLBPPPCorrections(meanDraft, getPolarStarDraft(), *vesselData.LBP), getPolarStarDraft())
	if got.FwdDraftWKeel != 3.302 {
		t.Errorf("FWD wKeel: expected 3.302, got %f", got.FwdDraftWKeel)
	}
	if got.MidDraftWKeel != 4.582 {
		t.Errorf("MID wKeel: expected 4.582, got %f", got.MidDraftWKeel)
	}
	if got.AftDraftWKeel != 6.102 {
		t.Errorf("AFT wKeel: expected 6.102, got %f", got.AftDraftWKeel)
	}
}

func TestPolarStar_TrimList_MMC(t *testing.T) {
	vesselData := getPolarStarTrimListVessel()
	meanDraft := MeanDrafts(getPolarStarTrimListMarks())
	draftsWKeel := CalcDraftsWKeel(meanDraft, CalcFullLBPPPCorrections(meanDraft, getPolarStarDraft(), *vesselData.LBP), getPolarStarDraft())
	got := CalcMMC(draftsWKeel, vesselData)
	if got != 4.612 {
		t.Errorf("MMC: expected 4.612, got %f", got)
	}
}

func TestPolarStar_TrimList_Hydrostatics(t *testing.T) {
	vesselData := getPolarStarTrimListVessel()
	meanDraft := MeanDrafts(getPolarStarTrimListMarks())
	draftsWKeel := CalcDraftsWKeel(meanDraft, CalcFullLBPPPCorrections(meanDraft, getPolarStarDraft(), *vesselData.LBP), getPolarStarDraft())
	mmc := CalcMMC(draftsWKeel, vesselData)
	hydrostatics := CalcHydrostatics(mmc, getPolarStarTrimListHydrostaticRows(), vesselData)
	if hydrostatics.Displacement != 19160.1 {
		t.Errorf("Displacement: expected 19160.1, got %f", hydrostatics.Displacement)
	}
	if hydrostatics.TPC != 45.2 {
		t.Errorf("TPC: expected 45.2, got %f", hydrostatics.TPC)
	}
	if hydrostatics.LCF != -6.962 {
		t.Errorf("LCF: expected -6.962, got %f", hydrostatics.LCF)
	}
}

func TestPolarStar_TrimList_FirstTrimCorrection(t *testing.T) {
	vesselData := getPolarStarTrimListVessel()
	meanDraft := MeanDrafts(getPolarStarTrimListMarks())
	draftsWKeel := CalcDraftsWKeel(meanDraft, CalcFullLBPPPCorrections(meanDraft, getPolarStarDraft(), *vesselData.LBP), getPolarStarDraft())
	mmc := CalcMMC(draftsWKeel, vesselData)
	hydrostatics := CalcHydrostatics(mmc, getPolarStarTrimListHydrostaticRows(), vesselData)
	got := CalcFirstTrimCorrection(draftsWKeel, hydrostatics.TPC, hydrostatics.LCF, *vesselData.LBP)
	if got != -481.481 {
		t.Errorf("1st trim: expected -481.481, got %f", got)
	}
}

func TestPolarStar_TrimList_SecondTrimCorrection(t *testing.T) {
	vesselData := getPolarStarTrimListVessel()
	meanDraft := MeanDrafts(getPolarStarTrimListMarks())
	draftsWKeel := CalcDraftsWKeel(meanDraft, CalcFullLBPPPCorrections(meanDraft, getPolarStarDraft(), *vesselData.LBP), getPolarStarDraft())
	got := CalcSecondTrimCorrection(draftsWKeel, getPolarStarTrimListMTCRows(), *vesselData.LBP)
	if got != 57.622 {
		t.Errorf("2nd trim: expected 57.622, got %f", got)
	}
}

func TestPolarStar_TrimList_ListCorrection(t *testing.T) {
	marks := getPolarStarTrimListMarks()
	vesselData := getPolarStarTrimListVessel()
	hr := getPolarStarTrimListHydrostaticRows()
	mmc := 4.612
	hydrostatics := CalcHydrostatics(mmc, hr, vesselData)
	tpcPort := 45.212
	tpcStbd := 45.129
	got := CalcListCorrection(marks, &tpcPort, &tpcStbd, hr, mmc, hydrostatics, vesselData)
	if got != 0.05 {
		t.Errorf("List corr: expected 0.050, got %f", got)
	}
}

func TestPolarStar_TrimList_DensityCorrection(t *testing.T) {
	vesselData := getPolarStarTrimListVessel()
	marks := getPolarStarTrimListMarks()
	meanDraft := MeanDrafts(marks)
	draftsWKeel := CalcDraftsWKeel(meanDraft, CalcFullLBPPPCorrections(meanDraft, getPolarStarDraft(), *vesselData.LBP), getPolarStarDraft())
	mmc := CalcMMC(draftsWKeel, vesselData)
	hr := getPolarStarTrimListHydrostaticRows()
	hydrostatics := CalcHydrostatics(mmc, hr, vesselData)
	firstTrim := CalcFirstTrimCorrection(draftsWKeel, hydrostatics.TPC, hydrostatics.LCF, *vesselData.LBP)
	secondTrim := CalcSecondTrimCorrection(draftsWKeel, getPolarStarTrimListMTCRows(), *vesselData.LBP)
	tpcPort := 45.212
	tpcStbd := 45.129
	listCorr := CalcListCorrection(marks, &tpcPort, &tpcStbd, hr, mmc, hydrostatics, vesselData)
	got := CalcDensityCorrection(hydrostatics.Displacement, firstTrim, secondTrim, listCorr, 1.017, 1.025)
	if got != -146.234 {
		t.Errorf("Density corr: expected -146.234, got %f", got)
	}
}
