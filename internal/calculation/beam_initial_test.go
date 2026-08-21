package calculation

// beam_initial_test.go
//
// Golden tests for BEAM (IMO 9591741) — Initial Draft Survey
// Survey: 444572_001, Loading, Port Vanino
// Source: Real survey data verified by certified marine surveyor
//
// Test strategy:
//   expected = unece*() functions from unece_formulas_test.go
//   got      = CalcDraft(getInitialDraft(), getVessel()).Field
//
//   unece*() functions implement UNECE 1992 formulas with r3() applied
//   to all values that appear as reported fields in UNECE Form C.
//   This matches the physical survey process where surveyors record
//   rounded values and use those in subsequent calculations.
//
// To add a new vessel test case:
//   1. Create new constants block (e.g. vessel2*, initial2*)
//   2. Create getDraft2() and getVessel2() helpers
//   3. Call unece*() functions with the new constants
//   4. Write Test* functions following the same pattern

import (
	"testing"

	"github.com/AVZotov/draft-survey/internal/types"
)

// =============================================================================
// SECTION 1: VESSEL CONSTANTS — BEAM (IMO 9591741)
// Marine vessel, Full LBP method, Manual LCF detection (UNECE standard mode)
// =============================================================================

const (
	beamLBP          = 283.5   // Length Between Perpendiculars, m
	beamBreadth      = 45.0    // Extreme breadth, m
	beamLightship    = 26328.0 // Lightship weight, MT
	beamTableDensity = 1.025   // Hydrostatic table density, t/m³
	beamSummerDraft  = 18.322  // Summer load line draft, m
	beamSummerDWT    = 179100.3
	beamSummerTPC    = 122.4
)

// =============================================================================
// SECTION 2: INITIAL DRAFT — MARK READINGS
// All readings: Direct method
// =============================================================================

const (
	iniFwdPort      = 9.35 // Forward port mark, m
	iniFwdStarboard = 9.25 // Forward starboard mark, m
	iniMidPort      = 9.45 // Midship port mark, m
	iniMidStarboard = 9.08 // Midship starboard mark, m
	iniAftPort      = 9.52 // Aft port mark, m
	iniAftStarboard = 9.47 // Aft starboard mark, m
)

// =============================================================================
// SECTION 3: INITIAL DRAFT — PP DISTANCES & KEEL
// Sign convention: F (forward of perpendicular) = positive
//                 A (aft of perpendicular)      = negative
// =============================================================================

const (
	iniDFwdSigned = +3.42  // dFwd = 3.42m, direction F → positive
	iniDMidSigned = -0.83  // dMid = 0.83m, direction A → negative
	iniDAftSigned = +15.40 // dAft = 15.40m, direction F → positive
	iniKeelFwd    = 0.0    // Keel plate thickness at fwd mark, mm
	iniKeelMid    = 0.0    // Keel plate thickness at mid mark, mm
	iniKeelAft    = 0.0    // Keel plate thickness at aft mark, mm
)

// =============================================================================
// SECTION 4: INITIAL DRAFT — HYDROSTATIC DATA
// Two rows bracketing MMC, LCF direction F → negative sign in calculations
// =============================================================================

const (
	iniHydroUpperDraft = 9.30    // Upper hydrostatic row draft, m
	iniHydroUpperDispl = 98821.8 // Upper displacement, MT
	iniHydroUpperTPC   = 112.7   // Upper TPC, MT/cm
	iniHydroUpperLCF   = 10.618  // Upper LCF magnitude, m (direction F → negative)

	iniHydroLowerDraft = 9.25    // Lower hydrostatic row draft, m
	iniHydroLowerDispl = 98258.2 // Lower displacement, MT
	iniHydroLowerTPC   = 112.7   // Lower TPC, MT/cm
	iniHydroLowerLCF   = 10.669  // Lower LCF magnitude, m (direction F → negative)
)

// =============================================================================
// SECTION 5: INITIAL DRAFT — MTC DATA
// =============================================================================

const (
	iniMTCPlusDraft  = 9.8    // Draft for MTC +0.5m row, m
	iniMTCPlusValue  = 2076.1 // MTC at +0.5m draft, MT·m/cm
	iniMTCMinusDraft = 8.8    // Draft for MTC -0.5m row, m
	iniMTCMinusValue = 2023.1 // MTC at -0.5m draft, MT·m/cm
)

// =============================================================================
// SECTION 6: INITIAL DRAFT — DENSITY & DEDUCTIBLES
// List TPC port/stbd = 0 (not entered) → List correction = 0
// =============================================================================

const (
	iniDockwaterDensity = 1.020 // Measured dockwater density, t/m³

	// Deductibles — taken as reported constants (not calculated from tanks)
	iniBallastWater = 53099.555 // Ballast water total weight, MT
	iniFreshWater   = 179.930   // Fresh water total weight, MT
	iniHFO          = 863.380   // Heavy fuel oil, MT
	iniMDO          = 204.930   // Marine diesel oil, MT
	iniLubOil       = 62.100    // Lubricating oil, MT
	iniOthers       = 17000.0   // Other deductibles (declared), MT
	// Bilge water = 0, Sewage water = 0
)

// =============================================================================
// SECTION 7: HELPERS — build Draft and VesselData from constants
// =============================================================================

func getBeamInitialDraft() types.Draft {
	return types.Draft{
		Marks: types.Marks{
			FwdPort:      types.Mark{Value: ptr64(iniFwdPort)},
			FwdStarboard: types.Mark{Value: ptr64(iniFwdStarboard)},
			MidPort:      types.Mark{Value: ptr64(iniMidPort)},
			MidStarboard: types.Mark{Value: ptr64(iniMidStarboard)},
			AftPort:      types.Mark{Value: ptr64(iniAftPort)},
			AftStarboard: types.Mark{Value: ptr64(iniAftStarboard)},
		},
		DistancePPFwd:  ptr64(3.42),
		PPFwdDirection: "F",
		DistancePPMid:  ptr64(0.83),
		PPMidDirection: "A",
		DistancePPAft:  ptr64(15.40),
		PPAftDirection: "F",
		KeelFwd:        ptr64(iniKeelFwd),
		KeelMid:        ptr64(iniKeelMid),
		KeelAft:        ptr64(iniKeelAft),
		HydrostaticRows: []types.HydrostaticRow{
			{
				Draft:        ptr64(iniHydroUpperDraft),
				Displacement: ptr64(iniHydroUpperDispl),
				TPC:          ptr64(iniHydroUpperTPC),
				LCF:          ptr64(iniHydroUpperLCF),
				LCFDirection: types.LCFDirectionForward,
			},
			{
				Draft:        ptr64(iniHydroLowerDraft),
				Displacement: ptr64(iniHydroLowerDispl),
				TPC:          ptr64(iniHydroLowerTPC),
				LCF:          ptr64(iniHydroLowerLCF),
				LCFDirection: types.LCFDirectionForward,
			},
		},
		MTCRows: []types.MTCRow{
			{Draft: ptr64(iniMTCPlusDraft), MTC: ptr64(iniMTCPlusValue)},
			{Draft: ptr64(iniMTCMinusDraft), MTC: ptr64(iniMTCMinusValue)},
		},
		Density:          ptr64(iniDockwaterDensity),
		TPCListPort:      fp(0),
		TPCListStarboard: fp(0),
		Deductibles: types.Deductibles{
			HFO:    ptr64(iniHFO),
			MDO:    ptr64(iniMDO),
			LubOil: ptr64(iniLubOil),
			Others: ptr64(iniOthers),
		},
		BallastWaterTanks: []types.Tank{
			{
				IsFWTTank: false,
				Volume:    ptr64(iniBallastWater), // ballast water passed as MT
				Density:   ptr64(1.000),           // ballast water passed as MT
			},
		},
		FreshWaterTanks: []types.Tank{
			{
				IsFWTTank: true,
				Volume:    ptr64(iniFreshWater),
			},
		},
	}
}

func getBeamVessel() types.VesselData {
	return types.VesselData{
		LBP:                  ptr64(beamLBP),
		Breadth:              ptr64(beamBreadth),
		Lightship:            ptr64(beamLightship),
		TableDensity:         ptr64(beamTableDensity),
		SummerDraft:          ptr64(beamSummerDraft),
		SummerDWT:            ptr64(beamSummerDWT),
		SummerTPC:            ptr64(beamSummerTPC),
		VesselType:           types.VesselTypeMarine,
		CorrectionMethod:     types.CorrectionMethodFullLBP,
		IsLcfDetectionManual: true, // UNECE standard: trust entered LCF direction
	}
}

// =============================================================================
// SECTION 8: UNECE reference values computed from constants
// These are the expected values for each test.
// All intermediate values follow UNECE Form C rounding policy.
// =============================================================================

func beamIniMeanDrafts() (meanF, meanM, meanA float64) {
	return uneceMeanDrafts(iniFwdPort, iniFwdStarboard, iniMidPort, iniMidStarboard, iniAftPort, iniAftStarboard)
}

func beamIniLBM() float64 {
	return uneceLBM(beamLBP, iniDFwdSigned, iniDAftSigned)
}

func beamIniPPCorrections() (fwdCorr, midCorr, aftCorr float64) {
	meanF, _, meanA := beamIniMeanDrafts()
	lbm := beamIniLBM()
	return unecePPCorrections(meanF, meanA, lbm, iniDFwdSigned, iniDMidSigned, iniDAftSigned)
}

func beamIniDraftsWKeel() (fwdWK, midWK, aftWK float64) {
	meanF, meanM, meanA := beamIniMeanDrafts()
	fwdCorr, midCorr, aftCorr := beamIniPPCorrections()
	return uneceDraftsWKeel(meanF, meanM, meanA, fwdCorr, midCorr, aftCorr, iniKeelFwd, iniKeelMid, iniKeelAft)
}

func beamIniMMC() float64 {
	fwdWK, midWK, aftWK := beamIniDraftsWKeel()
	return uneceMMC(fwdWK, midWK, aftWK)
}

func beamIniDisplacement() float64 {
	mmc := beamIniMMC()
	return uneceDisplacement(mmc, iniHydroLowerDraft, iniHydroLowerDispl, iniHydroUpperDraft, iniHydroUpperDispl)
}

func beamIniTPC() float64 {
	mmc := beamIniMMC()
	return uneceTPC(mmc, iniHydroLowerDraft, iniHydroLowerTPC, iniHydroUpperDraft, iniHydroUpperTPC)
}

func beamIniLCF() float64 {
	mmc := beamIniMMC()
	// Direction F → negative sign applied before interpolation
	return uneceLCF(mmc, iniHydroLowerDraft, -iniHydroLowerLCF, iniHydroUpperDraft, -iniHydroUpperLCF)
}

func beamIniTrueTrim() float64 {
	fwdWK, _, aftWK := beamIniDraftsWKeel()
	return uneceTrueTrim(fwdWK, aftWK)
}

func beamIniFirstTrimCorrection() float64 {
	return uneceFirstTrimCorrection(beamIniTrueTrim(), beamIniTPC(), beamIniLCF(), beamLBP)
}

func beamIniSecondTrimCorrection() float64 {
	deltaMTC := uneceDeltaMTC(iniMTCPlusValue, iniMTCMinusValue)
	return uneceSecondTrimCorrection(beamIniTrueTrim(), deltaMTC, beamLBP)
}

func beamIniDensityCorrection() float64 {
	displ := beamIniDisplacement()
	firstTrim := beamIniFirstTrimCorrection()
	secondTrim := beamIniSecondTrimCorrection()
	displCorrToTrimList := r3(displ + firstTrim + secondTrim) // list = 0
	return uneceDensityCorrection(displCorrToTrimList, iniDockwaterDensity, beamTableDensity)
}

func beamIniTotalDeductibles() float64 {
	return r3(iniBallastWater + iniFreshWater + iniHFO + iniMDO + iniLubOil + iniOthers)
}

func beamIniNetDisplacement() float64 {
	return uneceNetDisplacement(
		beamIniDisplacement(),
		beamIniFirstTrimCorrection(),
		beamIniSecondTrimCorrection(),
		0, // list correction = 0 (TPC list not entered)
		beamIniDensityCorrection(),
		beamIniTotalDeductibles(),
	)
}

// =============================================================================
// SECTION 9: TESTS
// Pattern: expected = beamIni*(), got = CalcDraft(...).Field
// =============================================================================

func TestBeam_Initial_MeanDraft_Fwd(t *testing.T) {
	expected, _, _ := beamIniMeanDrafts()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).MeanDraft.DraftFwdMean
	if expected != got {
		t.Errorf("MeanDraft Fwd: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_MeanDraft_Mid(t *testing.T) {
	_, expected, _ := beamIniMeanDrafts()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).MeanDraft.DraftMidMean
	if expected != got {
		t.Errorf("MeanDraft Mid: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_MeanDraft_Aft(t *testing.T) {
	_, _, expected := beamIniMeanDrafts()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).MeanDraft.DraftAftMean
	if expected != got {
		t.Errorf("MeanDraft Aft: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_LBM(t *testing.T) {
	expected := beamIniLBM()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).LBM
	if expected != got {
		t.Errorf("LBM: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_ObservedTrim(t *testing.T) {
	meanF, _, meanA := beamIniMeanDrafts()
	expected := uneceObservedTrim(meanF, meanA)
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).ObservedTrim
	if expected != got {
		t.Errorf("ObservedTrim: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_PPCorrection_Fwd(t *testing.T) {
	expected, _, _ := beamIniPPCorrections()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).PPCorrections.FwdCorrection
	if expected != got {
		t.Errorf("PP Fwd: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_PPCorrection_Mid(t *testing.T) {
	_, expected, _ := beamIniPPCorrections()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).PPCorrections.MidCorrection
	if expected != got {
		t.Errorf("PP Mid: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_PPCorrection_Aft(t *testing.T) {
	_, _, expected := beamIniPPCorrections()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).PPCorrections.AftCorrection
	if expected != got {
		t.Errorf("PP Aft: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_DraftWKeel_Fwd(t *testing.T) {
	expected, _, _ := beamIniDraftsWKeel()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).DraftsWKeel.FwdDraftWKeel
	if expected != got {
		t.Errorf("DraftWKeel Fwd: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_DraftWKeel_Mid(t *testing.T) {
	_, expected, _ := beamIniDraftsWKeel()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).DraftsWKeel.MidDraftWKeel
	if expected != got {
		t.Errorf("DraftWKeel Mid: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_DraftWKeel_Aft(t *testing.T) {
	_, _, expected := beamIniDraftsWKeel()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).DraftsWKeel.AftDraftWKeel
	if expected != got {
		t.Errorf("DraftWKeel Aft: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_TrueTrim(t *testing.T) {
	expected := beamIniTrueTrim()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).TrueTrim
	if expected != got {
		t.Errorf("TrueTrim: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_MeanFwdAft(t *testing.T) {
	fwdWK, _, aftWK := beamIniDraftsWKeel()
	expected := uneceMeanFwdAft(fwdWK, aftWK)
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).MeanFwdAft
	if expected != got {
		t.Errorf("MeanFwdAft: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_Deflection(t *testing.T) {
	fwdWK, midWK, aftWK := beamIniDraftsWKeel()
	expected := uneceDeflection(fwdWK, midWK, aftWK)
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).Deflection
	if expected != got {
		t.Errorf("Deflection: expected %.3f cm, got %.3f cm", expected, got)
	}
}

func TestBeam_Initial_ListMeters(t *testing.T) {
	expected := uneceListMeters(iniMidPort, iniMidStarboard)
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).ListMeters
	if expected != got {
		t.Errorf("ListMeters: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_ListDegrees(t *testing.T) {
	listM := uneceListMeters(iniMidPort, iniMidStarboard)
	expected := uneceListDegrees(listM, beamBreadth)
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).ListDegrees
	if expected != got {
		t.Errorf("ListDegrees: expected %.3f°, got %.3f°", expected, got)
	}
}

func TestBeam_Initial_MMC(t *testing.T) {
	expected := beamIniMMC()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).MMC
	if expected != got {
		t.Errorf("MMC: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_Displacement(t *testing.T) {
	expected := beamIniDisplacement()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).Hydrostatics.Displacement
	if expected != got {
		t.Errorf("Displacement: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_TPC(t *testing.T) {
	expected := beamIniTPC()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).Hydrostatics.TPC
	if expected != got {
		t.Errorf("TPC: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_LCF(t *testing.T) {
	expected := beamIniLCF()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).Hydrostatics.LCF
	if expected != got {
		t.Errorf("LCF: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_DeltaMTC(t *testing.T) {
	expected := uneceDeltaMTC(iniMTCPlusValue, iniMTCMinusValue)
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).DeltaMTC
	if expected != got {
		t.Errorf("DeltaMTC: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_FirstTrimCorrection(t *testing.T) {
	expected := beamIniFirstTrimCorrection()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).FirstTrimCorrection
	if expected != got {
		t.Errorf("1st Trim Correction: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_SecondTrimCorrection(t *testing.T) {
	expected := beamIniSecondTrimCorrection()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).SecondTrimCorrection
	if expected != got {
		t.Errorf("2nd Trim Correction: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_DensityCorrection(t *testing.T) {
	expected := beamIniDensityCorrection()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).DensityCorrection
	if expected != got {
		t.Errorf("Density Correction: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_TotalDeductibles(t *testing.T) {
	expected := beamIniTotalDeductibles()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).TotalDeductibles
	if expected != got {
		t.Errorf("Total Deductibles: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_NetDisplacement(t *testing.T) {
	expected := beamIniNetDisplacement()
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).NetDisplacement
	if expected != got {
		t.Errorf("Net Displacement: expected %.3f, got %.3f", expected, got)
	}
}

func TestBeam_Initial_Constant(t *testing.T) {
	expected := uneceConstant(beamIniNetDisplacement(), beamLightship)
	got := CalcDraft(getBeamInitialDraft(), getBeamVessel()).Constant
	if expected != got {
		t.Errorf("Constant: expected %.3f, got %.3f", expected, got)
	}
}
