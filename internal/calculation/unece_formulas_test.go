package calculation

// unece_formulas_test.go
//
// UNECE 1992 Reference Implementation
// Source: Code of Uniform Standards and Procedures for the Performance of
//         DRAUGHT SURVEYS of Coal Cargoes (ECE/ENERGY/19, Ed. October 15, 1991)
//
// Rounding policy:
//   r3() is applied to every value that appears as a discrete reported field
//   in UNECE Form C (ECE/ENERGY/19). Intermediate calculations that are NOT
//   recorded in any Form field are computed at full float64 precision.
//   This mirrors the physical survey process: a surveyor records rounded
//   values on paper (Form C) and uses those recorded values in subsequent
//   calculations. Any legal dispute would reference the reported (rounded)
//   values — not intermediate float64 precision.
//
// Usage:
//   These functions are the immutable UNECE standard implementation.
//   They are used as expected values in beam_initial_test.go.
//   The calculation package functions are tested against these references.

import "math"

// r3 rounds to 3 decimal places — standard survey reporting precision.
// Applied only to values that appear as reported fields in UNECE Form C.
func r3(v float64) float64 {
	return math.Round(v*1000) / 1000
}

// ptr64 returns a pointer to a float64 value.
// Used for building Draft and VesselData structs with pointer fields.
func ptr64(v float64) *float64 {
	return &v
}

// =============================================================================
// UNECE Form C — Draught Readings & Mean Drafts
// Lines 129-131 (fwd), 136-138 (aft), 145-147 (mid)
// Formula: mean = (port + starboard) / 2
// Reported in Form C → r3 applied
// =============================================================================

// uneceMeanDrafts computes mean port/starboard drafts at each mark position.
// UNECE Form C lines 131 (fwd mean), 138 (aft mean), 147 (mid mean)
func uneceMeanDrafts(fwdP, fwdS, midP, midS, aftP, aftS float64) (meanF, meanM, meanA float64) {
	meanF = r3((fwdP + fwdS) / 2)
	meanM = r3((midP + midS) / 2)
	meanA = r3((aftP + aftS) / 2)
	return
}

// =============================================================================
// UNECE Form B / Form C — Length Between Marks
// Formula: LBM = LBP - dAft(signed) + dFwd(signed)
// Direction convention: F (forward of PP) = positive, A (aft of PP) = negative
// Reported in survey report → r3 applied
// =============================================================================

// uneceLBM computes Length Between Marks.
// dFwdSigned: positive if mark is forward of FP, negative if aft
// dAftSigned: positive if mark is forward of AP, negative if aft
func uneceLBM(lbp, dFwdSigned, dAftSigned float64) float64 {
	return r3(lbp - dAftSigned + dFwdSigned)
}

// =============================================================================
// UNECE Form C — Corrections from Marks to Perpendiculars
// Lines 133-134 (fwd), 140-141 (aft), 148-149 (mid)
// Formula: Corr = ±(T / D1) × D2
//   T  = trim between draught marks (meanA - meanF)
//   D1 = LBM (distance between draught marks)
//   D2 = distance from mark to perpendicular (signed)
// Sign: depends on slope of waterplane at mark location vs perpendicular
// Reported in Form C → r3 applied
// =============================================================================

// unecePPCorrections computes corrections from marks to perpendiculars.
// UNECE Form C lines 133-134 (fwd), 140-141 (aft), 148-149 (mid)
func unecePPCorrections(meanF, meanA, lbm, dFwdSigned, dMidSigned, dAftSigned float64) (fwdCorr, midCorr, aftCorr float64) {
	trim := meanA - meanF // observed trim — NOT rounded, used only internally
	fwdCorr = r3(dFwdSigned * trim / lbm)
	midCorr = r3(dMidSigned * trim / lbm)
	aftCorr = r3(dAftSigned * trim / lbm)
	return
}

// =============================================================================
// UNECE Form C — Observed Trim
// Line 155: Trim = Dap - Dfp (from mean mark readings, before PP correction)
// Reported in Form C → r3 applied
// =============================================================================

// uneceObservedTrim computes observed trim from mean mark readings.
func uneceObservedTrim(meanF, meanA float64) float64 {
	return r3(meanA - meanF)
}

// =============================================================================
// UNECE Form C — Corrected Drafts (at Perpendiculars, with Keel)
// Lines 134 (fwd at FP), 141 (aft at AP), corrected midship
// Formula: wKeel = mean + ppCorr - (keel_mm / 1000)
// Keel correction always negative (subtracts keel plate thickness)
// Reported in Form C → r3 applied
// =============================================================================

// uneceDraftsWKeel computes drafts corrected for PP position and keel thickness.
// keelFwd, keelMid, keelAft are in millimetres.
func uneceDraftsWKeel(meanF, meanM, meanA, fwdCorr, midCorr, aftCorr, keelFwd, keelMid, keelAft float64) (fwdWK, midWK, aftWK float64) {
	fwdWK = r3(meanF + fwdCorr - (keelFwd / 1000))
	midWK = r3(meanM + midCorr - (keelMid / 1000))
	aftWK = r3(meanA + aftCorr - (keelAft / 1000))
	return
}

// =============================================================================
// UNECE Form C — True Trim
// Line 155: Trim = Dap - Dfp (from corrected perpendicular drafts)
// Positive = trim by stern, Negative = trim by head
// Reported in Form C → r3 applied
// =============================================================================

// uneceTrueTrim computes true trim from corrected perpendicular drafts.
func uneceTrueTrim(fwdWK, aftWK float64) float64 {
	return r3(aftWK - fwdWK)
}

// =============================================================================
// Mean of Forward & Aft corrected drafts
// Formula: MeanF&A = (FwdWKeel + AftWKeel) / 2
// Reported in survey → r3 applied
// =============================================================================

// uneceMeanFwdAft computes mean of fwd and aft corrected drafts.
func uneceMeanFwdAft(fwdWK, aftWK float64) float64 {
	return r3((fwdWK + aftWK) / 2)
}

// =============================================================================
// UNECE Form C — Hull Deflection (Hogging / Sagging)
// Line 150: Deflection = (MidWKeel - MeanF&A) × 100  [cm]
// Positive = Hogging (mid higher than mean F&A)
// Negative = Sagging (mid lower than mean F&A)
// Reported in Form C → r3 applied
// =============================================================================

// uneceDeflection computes hull deflection in centimetres.
func uneceDeflection(fwdWK, midWK, aftWK float64) float64 {
	meanFA := uneceMeanFwdAft(fwdWK, aftWK)
	return r3((midWK - meanFA) * 100)
}

// =============================================================================
// UNECE Form C — Mean of Mean Corrected Draft (Quarter Mean / MMC)
// Line 151-154: Standard formula for marine vessels
// Formula: MMC = (FwdWKeel + 6×MidWKeel + AftWKeel) / 8
// This is the draft used to enter hydrostatic tables.
// Reported in Form C → r3 applied
// =============================================================================

// uneceMMC computes Mean of Mean Corrected draft for marine vessels.
// UNECE Form C line 151 — standard formula (6M/8)
func uneceMMC(fwdWK, midWK, aftWK float64) float64 {
	return r3((fwdWK + 6*midWK + aftWK) / 8)
}

// =============================================================================
// List in meters and degrees
// listMeters = MidPort - MidStarboard
// listDegrees = arctan(listMeters / Breadth) × (180/π)
// Reported in survey → r3 applied
// =============================================================================

// uneceListMeters computes transverse list in metres.
func uneceListMeters(midPort, midStarboard float64) float64 {
	return r3(midPort - midStarboard)
}

// uneceListDegrees computes transverse list angle in degrees.
func uneceListDegrees(listMeters, breadth float64) float64 {
	return r3(math.Atan2(listMeters, breadth) * 180 / math.Pi)
}

// =============================================================================
// UNECE Form C — Linear Interpolation for Hydrostatic Tables
// Used for Displacement, TPC, LCF lookup at MMC
// Formula: result = lower + (fact - lowerX) × (upper - lower) / (upperX - lowerX)
// =============================================================================

// uneceInterpolate performs linear interpolation between two table values.
// Reported hydrostatic values → r3 applied to result
func uneceInterpolate(fact, lowerX, lowerVal, upperX, upperVal float64) float64 {
	return r3(lowerVal + (fact-lowerX)*(upperVal-lowerVal)/(upperX-lowerX))
}

// =============================================================================
// UNECE Form C — Hydrostatic Values at MMC
// Displacement (line 160), TPC, LCF interpolated from vessel tables
// LCF sign convention (manual mode — UNECE standard):
//   Direction F (forward of midship) → negative value
//   Direction A (aft of midship)     → positive value
// Reported in Form C → r3 applied via uneceInterpolate
// =============================================================================

// uneceDisplacement interpolates displacement at MMC from hydrostatic table.
func uneceDisplacement(mmc, lowerDraft, lowerDispl, upperDraft, upperDispl float64) float64 {
	return uneceInterpolate(mmc, lowerDraft, lowerDispl, upperDraft, upperDispl)
}

// uneceTPC interpolates TPC at MMC from hydrostatic table.
func uneceTPC(mmc, lowerDraft, lowerTPC, upperDraft, upperTPC float64) float64 {
	return uneceInterpolate(mmc, lowerDraft, lowerTPC, upperDraft, upperTPC)
}

// uneceLCF interpolates LCF at MMC from hydrostatic table.
// lowerLCF and upperLCF must already have correct sign applied:
//   F direction → negative, A direction → positive
func uneceLCF(mmc, lowerDraft, lowerLCF, upperDraft, upperLCF float64) float64 {
	return uneceInterpolate(mmc, lowerDraft, lowerLCF, upperDraft, upperLCF)
}

// =============================================================================
// UNECE Form C — First Trim Correction (LCF Correction)
// Line 162: LCF Corr = ±(TPC × LCF × T × 100) / LBP
// Sign rule per UNECE table:
//   Trim Aft (+) AND LCF Fwd (-) → negative correction
//   Trim Aft (+) AND LCF Aft (+) → positive correction
//   Trim Fwd (-) AND LCF Fwd (-) → positive correction
//   Trim Fwd (-) AND LCF Aft (+) → negative correction
// Reported in Form C → r3 applied
// =============================================================================

// uneceFirstTrimCorrection computes 1st trim correction per UNECE line 162.
func uneceFirstTrimCorrection(trueTrim, tpc, lcf, lbp float64) float64 {
	value := math.Abs(trueTrim * tpc * lcf * 100 / lbp)
	if (trueTrim > 0 && lcf <= 0) || (trueTrim < 0 && lcf >= 0) {
		return r3(-value)
	}
	return r3(value)
}

// =============================================================================
// UNECE Form C — Delta MTC
// ΔMTC = MTC(+0.5m) - MTC(-0.5m)
// Input to Nemoto (2nd trim) correction
// Reported in Form C → r3 applied
// =============================================================================

// uneceDeltaMTC computes change in MTC between ±0.5m drafts.
func uneceDeltaMTC(mtcPlus, mtcMinus float64) float64 {
	return r3(mtcPlus - mtcMinus)
}

// =============================================================================
// UNECE Form C — Second Trim Correction (Nemoto Correction)
// Line 163: Nemoto Corr = 50 × ΔMTC × T² / LBP
// Always positive — magnitude correction for hull form non-linearity
// Reported in Form C → r3 applied
// =============================================================================

// uneceSecondTrimCorrection computes 2nd trim correction per UNECE line 163.
func uneceSecondTrimCorrection(trueTrim, deltaMTC, lbp float64) float64 {
	return r3(50 * math.Pow(trueTrim, 2) * deltaMTC / lbp)
}

// =============================================================================
// UNECE Form C — Density Correction
// Line 164: DensityCorr = DisplCorrToTrimList × (ρ_actual - ρ_table) / ρ_table
// Applied when actual water density differs from hydrostatic table density
// Reported in Form C → r3 applied
// =============================================================================

// uneceDensityCorrection computes correction for water density difference.
// displCorrToTrimList = Displacement + 1stTrim + 2ndTrim + ListCorr
func uneceDensityCorrection(displCorrToTrimList, actualDensity, tableDensity float64) float64 {
	return r3(displCorrToTrimList * (actualDensity - tableDensity) / tableDensity)
}

// =============================================================================
// UNECE Form C — Net Displacement (Displacement corrected to deductibles)
// Line 170: NetDispl = Displ + 1stTrim + 2ndTrim + ListCorr + DensityCorr - Deductibles
// Reported in Form C → r3 applied
// =============================================================================

// uneceNetDisplacement computes final net displacement.
func uneceNetDisplacement(displ, firstTrim, secondTrim, listCorr, densityCorr, deductibles float64) float64 {
	displCorrToDensity := r3(displ + firstTrim + secondTrim + listCorr + densityCorr)
	return r3(displCorrToDensity - deductibles)
}

// =============================================================================
// UNECE Form B — Ship's Constant
// Line 090: Constant = Corrected Displacement - Lightship Weight
// Reported in Form B → r3 applied
// =============================================================================

// uneceConstant computes ship's constant.
func uneceConstant(netDisplacement, lightship float64) float64 {
	return r3(netDisplacement - lightship)
}
