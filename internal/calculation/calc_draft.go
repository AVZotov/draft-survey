package calculation

import (
	"math"

	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/internal/vessel"
)

func CalcDraft(draft types.Draft, v vessel.VesselData) DraftResult {
	var meanDraft types.MeanDraft
	var meanFwdAft float64
	var ppCorrections types.PPCorrections
	var draftsWKeel types.DraftsWKeel
	var mmc float64
	var observedTrim float64
	var trueTrim float64
	var listMeters float64
	var listDegrees float64
	var deflection float64
	var hydrostatics types.Hydrostatics
	var firstTrimCorrection float64
	var secondTrimCorrection float64
	var deltaMTC float64
	var listCorrection float64
	var totalTrimCorrection float64
	var densityCorrection float64
	var displacementCorrected float64
	var totalDeductibles float64
	var netDisplacement float64
	var constant float64
	var currentDWT float64
	var totalBwTanksWeight float64
	var totalFwTanksWeight float64

	meanDraft = MeanDrafts(draft.Marks)
	if v.CorrectionMethod == vessel.CorrectionMethodFullLBP {
		ppCorrections = CalcFullLBPPPCorrections(meanDraft, draft, v.LBP)
	}
	if v.CorrectionMethod == vessel.CorrectionMethodHalfLBP {
		ppCorrections = CalcHalfLBPPPCorrections(meanDraft, draft, v.LBP)
	}
	meanFwdAft = round3((draftsWKeel.FwdDraftWKeel + draftsWKeel.AftDraftWKeel) / 2)
	draftsWKeel = CalcDraftsWKeel(meanDraft, ppCorrections, draft)
	mmc = CalcMMC(draftsWKeel, v)
	observedTrim = round3(meanDraft.DraftAftMean - meanDraft.DraftFwdMean)
	trueTrim = round3(draftsWKeel.AftDraftWKeel - draftsWKeel.FwdDraftWKeel)
	listMeters = round3(markVal(draft.Marks.MidPort.Value) - markVal(draft.Marks.MidStarboard.Value))
	listDegrees = round3(math.Atan2(listMeters, v.Breadth) * 180 / math.Pi)
	deflection = round3((draftsWKeel.MidDraftWKeel - (draftsWKeel.FwdDraftWKeel+draftsWKeel.AftDraftWKeel)/2) * 100)
	if len(draft.HydrostaticRows) >= 2 {
		hydrostatics = CalcHydrostatics(mmc, draft.HydrostaticRows, v)
		firstTrimCorrection = CalcFirstTrimCorrection(draftsWKeel, hydrostatics.TPC, hydrostatics.LCF, v.LBP)
	}
	if len(draft.MTCRows) >= 2 {
		secondTrimCorrection = CalcSecondTrimCorrection(draftsWKeel, draft.MTCRows, v.LBP)
		deltaMTC = CalcDeltaMTC(draft.MTCRows)
	}

	listCorrection = CalcListCorrection(draft.Marks, draft.TPCListPort, draft.TPCListStarboard)
	totalTrimCorrection = round3(firstTrimCorrection + secondTrimCorrection + listCorrection)
	if draft.Density != nil {
		densityCorrection = CalcDensityCorrection(
			hydrostatics.Displacement,
			firstTrimCorrection,
			secondTrimCorrection,
			listCorrection,
			*draft.Density,
			*v.TableDensity,
		)
	}
	totalDeductibles = CalcTotalDeductibles(draft.BallastWaterTanks, draft.FreshWaterTanks, draft.Deductibles)
	displacementCorrected = round3(hydrostatics.Displacement + firstTrimCorrection + secondTrimCorrection + listCorrection + densityCorrection)
	netDisplacement = CalcNetDisplacement(
		hydrostatics.Displacement,
		firstTrimCorrection,
		secondTrimCorrection,
		listCorrection,
		densityCorrection,
		totalDeductibles,
	)
	constant = CalcConstant(netDisplacement, v.Lightship)
	currentDWT = CalcCurrentDWT(displacementCorrected, v.Lightship)
	totalBwTanksWeight = TotalTanksWeight(draft.BallastWaterTanks)
	totalFwTanksWeight = TotalTanksWeight(draft.FreshWaterTanks)

	return DraftResult{
		MeanDraft:             meanDraft,
		MeanFwdAft:            meanFwdAft,
		PPCorrections:         ppCorrections,
		DraftsWKeel:           draftsWKeel,
		MMC:                   mmc,
		ObservedTrim:          observedTrim,
		TrueTrim:              trueTrim,
		ListMeters:            listMeters,
		ListDegrees:           listDegrees,
		Deflection:            deflection,
		Hydrostatics:          hydrostatics,
		FirstTrimCorrection:   firstTrimCorrection,
		SecondTrimCorrection:  secondTrimCorrection,
		DeltaMTC:              deltaMTC,
		ListCorrection:        listCorrection,
		TotalTrimCorrection:   totalTrimCorrection,
		DensityCorrection:     densityCorrection,
		DisplacementCorrected: displacementCorrected,
		TotalDeductibles:      totalDeductibles,
		NetDisplacement:       netDisplacement,
		Constant:              constant,
		CurrentDWT:            currentDWT,
		TotalBwTanksWeight:    totalBwTanksWeight,
		TotalFwTanksWeight:    totalFwTanksWeight,
	}
}
