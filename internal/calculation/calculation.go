package calculation

import (
	"errors"
	"math"

	"github.com/AVZotov/draft-survey/internal/constants"
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/internal/vessel"
)

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

func MeanDrafts(m types.Marks) types.MeanDraft {
	return types.MeanDraft{
		DraftFwdMean: round3((markVal(m.FwdPort.Value) + markVal(m.FwdStarboard.Value)) / 2),
		DraftMidMean: round3((markVal(m.MidPort.Value) + markVal(m.MidStarboard.Value)) / 2),
		DraftAftMean: round3((markVal(m.AftPort.Value) + markVal(m.AftStarboard.Value)) / 2),
	}
}

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

func CalcMMC(draftsWKeel types.DraftsWKeel, v vessel.VesselData) float64 {
	if v.VesselType == vessel.VesselTypeMarine {
		return round3((draftsWKeel.FwdDraftWKeel + round3(6*draftsWKeel.MidDraftWKeel) + draftsWKeel.AftDraftWKeel) / 8)
	}

	if v.VesselType == vessel.VesselTypeRiver {
		return round3((draftsWKeel.FwdDraftWKeel + round3(4*draftsWKeel.MidDraftWKeel) + draftsWKeel.AftDraftWKeel) / 6)
	}

	if v.VesselType == vessel.VesselTypeBarge {
		return round3((round3(3*draftsWKeel.FwdDraftWKeel) + round3(14*draftsWKeel.MidDraftWKeel) + round3(3*draftsWKeel.AftDraftWKeel)) / 20)
	}

	return 0
}

func Interpolate(fact, lowerDraft, lowerValue, upperDraft, upperValue float64) float64 {
	result := round3(lowerValue + ((fact - lowerDraft) * (upperValue - lowerValue) / (upperDraft - lowerDraft)))
	return result
}

func CalcHydrostatics(mmc float64, hr []types.HydrostaticRow, v vessel.VesselData) types.Hydrostatics {
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

	if upper.LCFDirection == types.LCFDirectionFromAP || markVal(upper.LCF) > v.LBP*k3 {
		lowerLcf = (v.LBP / 2) - markVal(lower.LCF)
		upperLcf = (v.LBP / 2) - markVal(upper.LCF)
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

func CalcListCorrection(marks types.Marks, tpcListPort, tpcListStarboard float64) float64 {
	if markVal(marks.MidPort.Value) == markVal(marks.MidStarboard.Value) {
		return 0.0
	}
	return round3(6 * math.Abs(markVal(marks.MidPort.Value)-markVal(marks.MidStarboard.Value)) * math.Abs(tpcListPort-tpcListStarboard))
}

func CalcDensityCorrection(displacement float64, firstTrim float64, secondTrim float64, listCorrection float64, dockwaterDensity float64, tableDensity float64) float64 {
	displacementCorrected := round3(displacement + firstTrim + secondTrim + listCorrection)
	return round3(displacementCorrected * (dockwaterDensity - tableDensity) / tableDensity)
}

func CalcTotalDeductibles(bwt []types.Tank, fwt []types.Tank, d types.Deductibles) float64 {
	tbw := TotalTanksWeight(bwt)
	tfw := TotalTanksWeight(fwt)

	return round3(tbw + tfw + markVal(d.HFO) + markVal(d.MDO) + markVal(d.LubOil) + markVal(d.BilgeWater) + markVal(d.SewageWater) + markVal(d.Others))
}

func CalcNetDisplacement(displacement, firstTrim, secondTrim, listCorrection, densityCorrection, totalDeductibles float64) float64 {
	displCorrToDensity := round3(displacement + firstTrim + secondTrim + listCorrection + densityCorrection)
	return round3(displCorrToDensity - totalDeductibles)
}

func CalcCargoWeight(netDisplacementIni, netDisplacementFin float64) float64 {
	return round3(math.Abs(netDisplacementFin - netDisplacementIni))
}

func CalcConstant(netDisplacementIni float64, lightship float64) float64 {
	return round3(netDisplacementIni - lightship)
}

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

// Wrapper on true volume caclulations
func CalcBwTankVolume(trim, listDegrees float64, tank types.Tank) (float64, error) {
	if tank.Sounding == nil {
		return 0, errors.New("no sounding in measurements")
	}
	if !tank.Correction.IsValid() {
		return 0, errors.New("calibration data is incomplete or missing")
	}

	switch tank.Correction.TableType {
	case constants.VolumeCalibrationType2:
		return calcVolumeType2(*tank.Sounding, trim, listDegrees, tank.Correction), nil
	case constants.VolumeCalibrationType3:
		return calcVolumeType3(*tank.Sounding, trim, listDegrees, tank.Correction), nil
	default:
		return calcVolumeType1(*tank.Sounding, trim, listDegrees, tank.Correction), nil
	}
}
