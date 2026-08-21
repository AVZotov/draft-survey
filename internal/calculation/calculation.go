package calculation

import (
	"errors"
	"math"

	"github.com/AVZotov/draft-survey/internal/types"
)

// TotalTanksWeight sums the calculated weight (via Tank.CalcWeight) of every
// tank in the slice, rounding each tank's weight to 3 decimals before summing.
func TotalTanksWeight(tanks []types.Tank) float64 {
	var total float64
	for _, t := range tanks {
		total += round3(t.CalcWeight())
	}
	return total
}

func markVal(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

// vesselTypeOrDefault treats an unset VesselType as marine — the same
// fallback web/widgets/survey/form.templ's MMC-Method radio group already
// displays (`checked?={ VesselType != "river" && VesselType != "barge" }`
// pre-selects "Standard/marine" for the Go zero value too). A radio that's
// already visually checked on page load never fires htmx's change trigger
// unless the surveyor picks a different option, so VesselType can reach here
// as "" on a fully-entered survey. Without this, CalcMMC's exact-match
// switch falls through to its 0 fallback, and CalcDraft's LBM gate (which
// also checks VesselType == marine) skips assignment entirely — both read
// as "nothing was calculated" even though every input was entered.
func vesselTypeOrDefault(v types.VesselType) types.VesselType {
	if v != types.VesselTypeRiver && v != types.VesselTypeBarge {
		return types.VesselTypeMarine
	}
	return v
}

// correctionMethodOrDefault treats an unset CorrectionMethod as Full LBP —
// the same fallback web/widgets/survey/form.templ's method radio group
// already displays (`checked?={ CorrectionMethod != "Half LBP" }` pre-selects
// "Full LBP" for the Go zero value too). A radio that's already visually
// checked on page load never fires htmx's change trigger unless the
// surveyor picks a different option, so CorrectionMethod can reach here as
// "" on a fully-entered survey. Without this, CalcDraft's two exact-match
// CorrectionMethod checks both fail and PPCorrections is silently left at
// its zero value — same class of bug as vesselTypeOrDefault above.
func correctionMethodOrDefault(v types.CorrectionMethod) types.CorrectionMethod {
	if v == types.CorrectionMethodHalfLBP {
		return types.CorrectionMethodHalfLBP
	}
	return types.CorrectionMethodFullLBP
}

// MeanDrafts averages each pair of port/starboard draft marks (forward,
// midship, aft) into a MeanDraft. A nil mark reads as 0.
func MeanDrafts(m types.Marks) types.MeanDraft {
	return types.MeanDraft{
		DraftFwdMean: round3((markVal(m.FwdPort.Value) + markVal(m.FwdStarboard.Value)) / 2),
		DraftMidMean: round3((markVal(m.MidPort.Value) + markVal(m.MidStarboard.Value)) / 2),
		DraftAftMean: round3((markVal(m.AftPort.Value) + markVal(m.AftStarboard.Value)) / 2),
	}
}

// CalcFullLBPLBM computes the length between marks (LBM) for the full-LBP
// correction method: LBP adjusted by the forward and aft perpendicular
// distances, each signed by its recorded direction ("A" = aft, negated).
func CalcFullLBPLBM(draft types.Draft, lbp float64) float64 {
	var dFwdDir, dAftDir float64

	if dFwdDir = markVal(draft.DistancePPFwd); draft.PPFwdDirection == "A" {
		dFwdDir *= -1
	}
	if dAftDir = markVal(draft.DistancePPAft); draft.PPAftDirection == "A" {
		dAftDir *= -1
	}

	return round3(lbp - dAftDir + dFwdDir)
}

// CalcHalfLBPLBM computes the two half-length spans (aft-to-midship and
// midship-to-forward) used by the half-LBP correction method, each derived
// from half of LBP adjusted by the relevant perpendicular distances.
func CalcHalfLBPLBM(draft types.Draft, lbp float64) (lbmAftMid float64, lbmMidFwd float64) {
	var dFwdDir, dMidDir, dAftDir float64

	if dFwdDir = markVal(draft.DistancePPFwd); draft.PPFwdDirection == "A" {
		dFwdDir *= -1
	}
	if dMidDir = markVal(draft.DistancePPMid); draft.PPMidDirection == "A" {
		dMidDir *= -1
	}
	if dAftDir = markVal(draft.DistancePPAft); draft.PPAftDirection == "A" {
		dAftDir *= -1
	}

	lbmMidFwd = round3((lbp / 2) - dMidDir - dFwdDir)
	lbmAftMid = round3((lbp / 2) - dAftDir - dMidDir)

	return lbmAftMid, lbmMidFwd
}

// CalcFullLBPPPCorrections computes the forward/midship/aft perpendicular-
// point corrections for the full-LBP method, proportional to the observed
// trim (aft mean minus forward mean) and each mark's distance from LBM.
func CalcFullLBPPPCorrections(m types.MeanDraft, draft types.Draft, lbp float64) types.PPCorrections {
	trim := m.DraftAftMean - m.DraftFwdMean
	var dFwdDir, dMidDir, dAftDir float64

	if dFwdDir = markVal(draft.DistancePPFwd); draft.PPFwdDirection == "A" {
		dFwdDir *= -1
	}
	if dMidDir = markVal(draft.DistancePPMid); draft.PPMidDirection == "A" {
		dMidDir *= -1
	}
	if dAftDir = markVal(draft.DistancePPAft); draft.PPAftDirection == "A" {
		dAftDir *= -1
	}
	lbm := round3(lbp - dAftDir + dFwdDir)
	return types.PPCorrections{
		FwdCorrection: round3(dFwdDir * trim / lbm),
		MidCorrection: round3(dMidDir * trim / lbm),
		AftCorrection: round3(dAftDir * trim / lbm),
	}
}

// CalcHalfLBPPPCorrections computes the forward/midship/aft perpendicular-
// point corrections for the half-LBP method. Unlike CalcFullLBPPPCorrections,
// the aft correction depends on the midship correction already applied
// (midship-with-keel is computed first), since the method treats the vessel
// as two independent half-lengths rather than one full span.
func CalcHalfLBPPPCorrections(m types.MeanDraft, draft types.Draft, lbp float64) types.PPCorrections {
	var dFwdDir, dMidDir, dAftDir float64

	if dFwdDir = markVal(draft.DistancePPFwd); draft.PPFwdDirection == "A" {
		dFwdDir *= -1
	}
	if dMidDir = markVal(draft.DistancePPMid); draft.PPMidDirection == "A" {
		dMidDir *= -1
	}
	if dAftDir = markVal(draft.DistancePPAft); draft.PPAftDirection == "A" {
		dAftDir *= -1
	}

	lbmMidFwd := round3((lbp / 2) - dMidDir - dFwdDir)
	lbmAftMid := round3((lbp / 2) - dAftDir - dMidDir)

	fwdCorr := round3(dFwdDir * (m.DraftMidMean - m.DraftFwdMean) / lbmMidFwd)
	midCorr := round3(dMidDir * (m.DraftMidMean - m.DraftFwdMean) / lbmMidFwd)
	midWKeel := round3(m.DraftMidMean + midCorr - (markVal(draft.KeelMid) / 1000))
	aftCorr := round3(dAftDir * (m.DraftAftMean - midWKeel) / lbmAftMid)

	return types.PPCorrections{
		FwdCorrection: fwdCorr,
		MidCorrection: midCorr,
		AftCorrection: aftCorr,
	}
}

// CalcDraftsWKeel applies the PP corrections and keel corrections (keel
// readings in mm, negated and converted to metres) to each mean draft,
// producing the forward/midship/aft drafts with keel.
func CalcDraftsWKeel(
	meanDraft types.MeanDraft, ppCorrections types.PPCorrections, draft types.Draft) types.DraftsWKeel {
	keelCorrectionFwd := -1 * markVal(draft.KeelFwd) / 1000
	keelCorrectionMid := -1 * markVal(draft.KeelMid) / 1000
	keelCorrectionAft := -1 * markVal(draft.KeelAft) / 1000

	return types.DraftsWKeel{
		FwdDraftWKeel: round3(meanDraft.DraftFwdMean + ppCorrections.FwdCorrection + keelCorrectionFwd),
		MidDraftWKeel: round3(meanDraft.DraftMidMean + ppCorrections.MidCorrection + keelCorrectionMid),
		AftDraftWKeel: round3(meanDraft.DraftAftMean + ppCorrections.AftCorrection + keelCorrectionAft),
	}
}

// CalcMMC computes the Mean of Mean Constant draft from the forward/midship/
// aft drafts-with-keel, using the weighting formula for the vessel's type
// (marine: 1-6-1/8, river: 1-4-1/6, barge: 3-14-3/20). An unset VesselType
// is treated as marine (see vesselTypeOrDefault); an unrecognized type
// returns 0.
func CalcMMC(draftsWKeel types.DraftsWKeel, v types.VesselData) float64 {
	switch vesselTypeOrDefault(v.VesselType) {
	case types.VesselTypeMarine:
		return round3((draftsWKeel.FwdDraftWKeel + round3(6*draftsWKeel.MidDraftWKeel) + draftsWKeel.AftDraftWKeel) / 8)
	case types.VesselTypeRiver:
		return round3((draftsWKeel.FwdDraftWKeel + round3(4*draftsWKeel.MidDraftWKeel) + draftsWKeel.AftDraftWKeel) / 6)
	case types.VesselTypeBarge:
		return round3((round3(3*draftsWKeel.FwdDraftWKeel) + round3(14*draftsWKeel.MidDraftWKeel) + round3(3*draftsWKeel.AftDraftWKeel)) / 20)
	default:
		return 0
	}
}

// Interpolate performs linear interpolation of a value at fact between the
// (lowerDraft, lowerValue) and (upperDraft, upperValue) reference points.
// Returns 0 if the two reference drafts are equal (would divide by zero).
func Interpolate(fact, lowerDraft, lowerValue, upperDraft, upperValue float64) float64 {
	if upperDraft == lowerDraft {
		return 0
	}

	result := round3(lowerValue + ((fact - lowerDraft) * (upperValue - lowerValue) / (upperDraft - lowerDraft)))
	return result
}

// CalcHydrostatics interpolates displacement, TPC, and LCF at MMC between the
// two entered hydrostatic table rows (order-independent — rows are sorted by
// draft internally). LCF direction is either trusted as entered by the
// surveyor (UNECE manual mode) or auto-detected from AP by magnitude against
// LBP*0.045 (DSGear k3 auto mode), per v.IsLcfDetectionManual.
func CalcHydrostatics(mmc float64, hr []types.HydrostaticRow, v types.VesselData) types.Hydrostatics {
	var lower, upper types.HydrostaticRow
	if markVal(hr[0].Draft) < markVal(hr[1].Draft) {
		lower = hr[0]
		upper = hr[1]
	} else {
		lower = hr[1]
		upper = hr[0]
	}
	displacement := Interpolate(mmc, markVal(lower.Draft), markVal(lower.Displacement), markVal(upper.Draft), markVal(upper.Displacement))
	tpc := Interpolate(mmc, markVal(lower.Draft), markVal(lower.TPC), markVal(upper.Draft), markVal(upper.TPC))

	const k3 = 0.045

	lowerLcf := markVal(lower.LCF)
	upperLcf := markVal(upper.LCF)

	// Auto mode (DSGear k3): detect LCA from AP by magnitude
	// Manual mode (UNECE): trust direction entered by surveyor
	if !v.IsLcfDetectionManual && (upper.LCFDirection == types.LCFDirectionFromAP || markVal(upper.LCF) > markVal(v.LBP)*k3) {
		lowerLcf = (markVal(v.LBP) / 2) - markVal(lower.LCF)
		upperLcf = (markVal(v.LBP) / 2) - markVal(upper.LCF)
	} else {
		if lower.LCFDirection == types.LCFDirectionForward {
			lowerLcf *= -1
		}
		if upper.LCFDirection == types.LCFDirectionForward {
			upperLcf *= -1
		}
	}

	lcf := Interpolate(mmc, markVal(lower.Draft), lowerLcf, markVal(upper.Draft), upperLcf)

	return types.Hydrostatics{
		Displacement: displacement,
		TPC:          tpc,
		LCF:          lcf,
	}
}

// CalcFirstTrimCorrection computes the first (layer) trim correction from
// true trim, hydrostatic TPC, and LCF at MMC, applying the UNECE sign rule:
// negative when trim and LCF point the same direction, positive otherwise.
func CalcFirstTrimCorrection(dwk types.DraftsWKeel, tpc float64, lcf float64, lbp float64) float64 {
	trueTrim := dwk.AftDraftWKeel - dwk.FwdDraftWKeel
	var firstTrimCorrection float64

	if trueTrim < 0 && lcf >= 0 || trueTrim > 0 && lcf <= 0 {
		firstTrimCorrection = -1 * (math.Abs(trueTrim * tpc * lcf * 100 / lbp))
	} else {
		firstTrimCorrection = math.Abs(trueTrim * tpc * lcf * 100 / lbp)
	}

	return round3(firstTrimCorrection)
}

// CalcDeltaMTC is supporting function to get value for frontend
func CalcDeltaMTC(mtcRows []types.MTCRow) float64 {
	var lowerMtcRow, upperMtcRow types.MTCRow
	if markVal(mtcRows[0].Draft) < markVal(mtcRows[1].Draft) {
		lowerMtcRow = mtcRows[0]
		upperMtcRow = mtcRows[1]
	} else {
		lowerMtcRow = mtcRows[1]
		upperMtcRow = mtcRows[0]
	}

	return round3(markVal(upperMtcRow.MTC) - markVal(lowerMtcRow.MTC))
}

// CalcSecondTrimCorrection computes the second (Nemoto's) trim correction
// from true trim and the delta MTC between the two entered MTC rows
// (order-independent — rows are sorted by draft internally).
func CalcSecondTrimCorrection(dwk types.DraftsWKeel, mtcRows []types.MTCRow, lbp float64) float64 {
	var lowerMtcRow, upperMtcRow types.MTCRow
	if markVal(mtcRows[0].Draft) < markVal(mtcRows[1].Draft) {
		lowerMtcRow = mtcRows[0]
		upperMtcRow = mtcRows[1]
	} else {
		lowerMtcRow = mtcRows[1]
		upperMtcRow = mtcRows[0]
	}

	deltaMtc := markVal(upperMtcRow.MTC) - markVal(lowerMtcRow.MTC)
	trueTrim := dwk.AftDraftWKeel - dwk.FwdDraftWKeel

	return round3(50 * math.Pow(trueTrim, 2) * deltaMtc / lbp)
}

// CalcListCorrection calculates list correction.
// Automatically selects calculation method based on available data:
//
//	V1: when TPCListPort and TPCListStarboard are entered manually by surveyor
//	V2: when hydrostatic TPC values are equal — uses Summer TPC interpolation
//	0:  when TPC values differ and no manual TPC list entered
func CalcListCorrection(marks types.Marks, tpcListPort, tpcListStarboard *float64, hydroRows []types.HydrostaticRow, mmc float64, hydrostatics types.Hydrostatics, v types.VesselData) float64 {
	if markVal(marks.MidPort.Value) == markVal(marks.MidStarboard.Value) {
		return 0.0
	}

	if tpcListPort != nil && tpcListStarboard != nil {
		return calcListCorrectionV1(marks, *tpcListPort, *tpcListStarboard)
	}

	if len(hydroRows) >= 2 && markVal(hydroRows[0].TPC) == markVal(hydroRows[1].TPC) {
		return calcListCorrectionV2(marks, mmc, hydrostatics.TPC, v)
	}

	return 0
}

// Industry-standard list correction formula (not part of UNECE 1992 Code).
// UNECE 1992 contains no list correction formula.
// This is an app-specific addition following common survey practice.
func calcListCorrectionV1(marks types.Marks, tpcListPort, tpcListStarboard float64) float64 {
	if markVal(marks.MidPort.Value) == markVal(marks.MidStarboard.Value) {
		return 0.0
	}
	return round3(6 * math.Abs(markVal(marks.MidPort.Value)-markVal(marks.MidStarboard.Value)) * math.Abs(tpcListPort-tpcListStarboard))
}

// CalcListCorrectionV2 calculates list correction using interpolation
// between MMC hydrostatic TPC and Summer TPC as the upper reference point.
// Used when hydrostatic table TPC values are identical for upper and lower rows
// (standard interpolation would yield zero correction).
// Industry-standard list correction formula (not part of UNECE 1992 Code).
// UNECE 1992 contains no list correction formula.
// This is an app-specific addition following common survey practice.
func calcListCorrectionV2(marks types.Marks, mmc float64, hydroTPC float64, v types.VesselData) float64 {
	if markVal(marks.MidPort.Value) == markVal(marks.MidStarboard.Value) {
		return 0.0
	}

	if v.SummerDraft == nil || *v.SummerDraft <= mmc {
		return 0.0
	}

	midPort := markVal(marks.MidPort.Value)
	midStbd := markVal(marks.MidStarboard.Value)

	// Interpolate TPC at each midship draft
	// using (MMC, hydroTPC) as lower point and (SummerDraft, SummerTPC) as upper point
	tpcPort := round3(hydroTPC + (midPort-mmc)*(markVal(v.SummerTPC)-hydroTPC)/(*v.SummerDraft-mmc))
	tpcStbd := round3(hydroTPC + (midStbd-mmc)*(markVal(v.SummerTPC)-hydroTPC)/(*v.SummerDraft-mmc))

	return round3(6 * math.Abs(midPort-midStbd) * math.Abs(tpcPort-tpcStbd))
}

// CalcDensityCorrection corrects the trim/list-adjusted displacement for the
// difference between actual dockwater density and the hydrostatic table's
// reference density, proportional to that displacement.
func CalcDensityCorrection(displacement float64, firstTrim float64, secondTrim float64, listCorrection float64, dockwaterDensity float64, tableDensity float64) float64 {
	displacementCorrected := round3(displacement + firstTrim + secondTrim + listCorrection)
	return round3(displacementCorrected * (dockwaterDensity - tableDensity) / tableDensity)
}

// CalcTotalDeductibles sums everything subtracted from displacement to reach
// net cargo-relevant displacement: total ballast and fresh water tank weight
// plus the surveyor-entered deductible weights (fuel, lube oil, bilge/sewage
// water, other).
func CalcTotalDeductibles(bwt []types.Tank, fwt []types.Tank, d types.Deductibles) float64 {
	tbw := TotalTanksWeight(bwt)
	tfw := TotalTanksWeight(fwt)

	return round3(tbw + tfw + markVal(d.HFO) + markVal(d.MDO) + markVal(d.LubOil) + markVal(d.BilgeWater) + markVal(d.SewageWater) + markVal(d.Others))
}

// CalcNetDisplacement applies all corrections (trim, list, density) to the
// hydrostatic displacement, then subtracts total deductibles to yield the
// net displacement used for cargo and constant calculations.
func CalcNetDisplacement(displacement, firstTrim, secondTrim, listCorrection, densityCorrection, totalDeductibles float64) float64 {
	displCorrToDensity := round3(displacement + firstTrim + secondTrim + listCorrection + densityCorrection)
	return round3(displCorrToDensity - totalDeductibles)
}

// CalcCargoWeight returns the absolute difference in net displacement
// between the initial and final draft readings.
func CalcCargoWeight(netDisplacementIni, netDisplacementFin float64) float64 {
	return round3(math.Abs(netDisplacementFin - netDisplacementIni))
}

// CalcConstant returns the survey constant: net displacement minus lightship
// weight.
func CalcConstant(netDisplacement float64, lightship float64) float64 {
	return round3(netDisplacement - lightship)
}

// CalcCurrentDWT returns the current deadweight: trim/list/density-corrected
// displacement minus lightship weight.
func CalcCurrentDWT(displCorrToDensity float64, lightship float64) float64 {
	return round3(displCorrToDensity - lightship)
}

/*
CalcVolumeByRows performs bilinear interpolation of tank volume.
Used as a base for all calibration table types (trim or list correction).
First interpolates by trim/list angle, then by sounding.
Parameters:

	sounding — actual tank sounding (measured by tape), m
	trimOrList — actual trim (meanA - meanF) or list value, m
	rows — two calibration rows bracketing the actual sounding
	tableLow — lower trim/list value from calibration table (TTL)
	tableHigh — upper trim/list value from calibration table (TTU)

If trimOrList == 0, interpolation is performed by sounding only (1D).
*/
func calcVolumeBy1DRows(sounding float64, rows []types.CorrectionRows) float64 {
	// Sort rows by sounding
	var lower, upper types.CorrectionRows
	if markVal(rows[0].TableSounding) < markVal(rows[1].TableSounding) {
		lower = rows[0]
		upper = rows[1]
	} else {
		lower = rows[1]
		upper = rows[0]
	}
	return Interpolate(sounding,
		markVal(lower.TableSounding), markVal(lower.TableVolume),
		markVal(upper.TableSounding), markVal(upper.TableVolume))
}

func calcVolumeByRows(sounding, trimOrList float64, rows []types.CalibrationRow, tableLow, tableHigh float64) float64 {
	var lower, upper types.CalibrationRow
	if markVal(rows[0].Sounding) < markVal(rows[1].Sounding) {
		lower = rows[0]
		upper = rows[1]
	} else {
		lower = rows[1]
		upper = rows[0]
	}

	if trimOrList == 0 {
		return Interpolate(sounding,
			markVal(lower.Sounding), markVal(lower.VolumeLow),
			markVal(upper.Sounding), markVal(upper.VolumeLow))
	}

	ab1 := Interpolate(trimOrList,
		tableLow, markVal(lower.VolumeLow),
		tableHigh, markVal(lower.VolumeUp))

	ab2 := Interpolate(trimOrList,
		tableLow, markVal(upper.VolumeLow),
		tableHigh, markVal(upper.VolumeUp))

	return Interpolate(sounding, markVal(lower.Sounding), ab1, markVal(upper.Sounding), ab2)
}

/*
	 CalcVolumeType1 calculates tank volume using Type 1 calibration table.
	 Type 1 (most common): table contains Volume directly for each
	 (sounding × trim) combination. List correction is optional.

	 Parameters:

		sounding — actual tank sounding, m
		trim     — vessel true trim (meanA - meanF), m
		list     — observed vessel list (midPort - midStbd), m
		data     — tank calibration data
*/
func calcVolumeType1(sounding, trim, list float64, data types.VolumeCalibrationData) float64 {
	ttl := markVal(data.TableTrimLow)
	ttu := markVal(data.TableTrimUpper)
	tll := markVal(data.TableListLow)
	tlu := markVal(data.TableListUpper)

	volume := calcVolumeByRows(sounding, trim, data.TrimRows, ttl, ttu)

	if len(data.ListRows) > 0 {
		volume += calcVolumeByRows(sounding, list, data.ListRows, tll, tlu)
	}

	return round3(volume)
}

/*
		CalcVolumeType2 calculates tank volume using Type 2 calibration table.

	 Type 2 (rare): two separate tables —

		Table 1: sounding correction by trim/list
		Table 2: volume at corrected sounding (trim = 0)

	 Workflow:
	  1. Interpolate sounding correction from Table 1 (trim rows)
	  2. Apply list correction from Table 1 (list rows) if available
	  3. Correct actual sounding by found corrections
	  4. Interpolate final volume from Table 2 at corrected sounding

	 Parameters:

		sounding — actual tank sounding, m
		trim     — observed vessel trim, m
		list     — observed vessel list, m
		data     — Table 1 calibration data (sounding corrections)
		volumeRows — Table 2 rows (volume at corrected sounding, trim = 0)
*/
func calcVolumeType2(sounding, trim, list float64, data types.VolumeCalibrationData) float64 {
	ttl := markVal(data.TableTrimLow)
	ttu := markVal(data.TableTrimUpper)
	tll := markVal(data.TableListLow)
	tlu := markVal(data.TableListUpper)

	// Step 1: sounding correction by trim
	soundingCorr := math.Round(calcVolumeByRows(sounding, trim, data.TrimRows, ttl, ttu))

	// Step 2: list correction to sounding
	if len(data.ListRows) > 0 {
		soundingCorr += math.Round(calcVolumeByRows(sounding, list, data.ListRows, tll, tlu))
	}

	// Step 3: corrected sounding
	correctedSounding := sounding + (soundingCorr / 1000)

	// Step 4: volume at corrected sounding (1D)
	return calcVolumeBy1DRows(correctedSounding, data.CorrectionRows)
}

/*
CalcVolumeType3 calculates tank volume using Type 3 calibration table.
Type 3 (rare): two separate tables —

	Table 1: volume correction by trim/list
	Table 2: base volume at actual sounding, trim = 0

Workflow:
 1. Get base volume from Table 2 at actual sounding (no trim)
 2. Interpolate volume correction from Table 1 by trim
 3. Add list correction from Table 1 if available
 4. Sum: base volume + trim correction + list correction

Parameters:

	sounding   — actual tank sounding, m
	trim       — observed vessel trim, m
	list       — observed vessel list, m
	data       — Table 1 calibration data (volume corrections)
	volumeRows — Table 2 rows (base volume at zero trim)
*/
func calcVolumeType3(sounding, trim, list float64, data types.VolumeCalibrationData) float64 {
	ttl := markVal(data.TableTrimLow)
	ttu := markVal(data.TableTrimUpper)
	tll := markVal(data.TableListLow)
	tlu := markVal(data.TableListUpper)

	// Step 1: base volume at zero trim (1D)
	baseVolume := calcVolumeBy1DRows(sounding, data.CorrectionRows)

	// Step 2: volume correction by trim
	trimCorr := calcVolumeByRows(sounding, trim, data.TrimRows, ttl, ttu)

	// Step 3: volume correction by list
	listCorr := 0.0
	if len(data.ListRows) > 0 {
		listCorr = calcVolumeByRows(sounding, list, data.ListRows, tll, tlu)
	}

	return round3(baseVolume + trimCorr + listCorr)
}

// CalcBwTankVolume computes a ballast water tank's volume at the given trim
// and list from its calibration table, dispatching to the Type 1/2/3
// calibration-table workflow per tank.Correction.TableType. Returns an error
// if the tank has no sounding entered or its calibration data is incomplete.
func CalcBwTankVolume(trim, listDegrees float64, tank types.Tank) (float64, error) {
	if tank.Sounding == nil {
		return 0, errors.New("no sounding in measurements")
	}
	if !tank.Correction.IsValid() {
		return 0, errors.New("calibration data is incomplete or missing")
	}

	switch tank.Correction.TableType {
	case types.CalibrationTypeSoundingCorrection:
		return calcVolumeType2(*tank.Sounding, trim, listDegrees, tank.Correction), nil
	case types.CalibrationTypeVolumeCorrection:
		return calcVolumeType3(*tank.Sounding, trim, listDegrees, tank.Correction), nil
	default:
		return calcVolumeType1(*tank.Sounding, trim, listDegrees, tank.Correction), nil
	}
}
