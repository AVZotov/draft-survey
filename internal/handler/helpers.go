package handler

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/AVZotov/draft-survey/internal/calculation"
	"github.com/AVZotov/draft-survey/internal/constants"
	"github.com/AVZotov/draft-survey/internal/format"
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/gofiber/fiber/v2"
)

type tankProps struct {
	surveyID   string
	draftIndex int
	survey     *types.Survey
	user       *types.User
	tankID     string
	tankIndex  int
	bwTanks    []types.Tank
	fwTanks    []types.Tank
	trim       *float64
	trimDir    string
	list       *float64
	listDir    string
}

type draftProps struct {
	user       *types.User
	survey     *types.Survey
	surveyID   string
	draftIndex int
}

type surveyProps struct {
	user     *types.User
	survey   *types.Survey
	surveyID string
}

var ErrEmptyField = errors.New("empty field")
var ErrUndefinedTankType = errors.New("undefined tank type")
var ErrUndefinedHXContext = errors.New("undefined HX context")
var ErrNilPointer = errors.New("nil poiner error")

func parseFloat(c *fiber.Ctx, name string) (*float64, error) {

	v := c.FormValue(name)
	if v == "" {
		return nil, ErrEmptyField
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil, ErrEmptyField
	}
	return &f, nil
}

func parseInt(c *fiber.Ctx, field string) (*int, error) {
	v := c.FormValue(field)
	if v == "" {
		return nil, ErrEmptyField
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return nil, ErrEmptyField
	}
	return &i, nil
}

func parseString(c *fiber.Ctx, field string) (string, error) {
	v := c.FormValue(field)
	if v == "" {
		return "", ErrEmptyField
	}
	return v, nil
}

// Deprecated To be deleted and substituted with granulated parsings
func (h *Handler) parseDraft(c *fiber.Ctx, survey *types.Survey) {
	tableDensity, err := parseFloat(c, constants.TableDensity)
	if err == nil {
		survey.VesselData.TableDensity = tableDensity
	}
	for i := range survey.Drafts {
		//Getting draft marks
		fwdPort, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.FwdPort, i))
		if err == nil {
			survey.Drafts[i].Marks.FwdPort.Value = fwdPort
		}
		fwdPortMethod, err := parseString(c, fmt.Sprintf("%s-d%d", constants.FwdPortMarkRead, i))
		if err == nil {
			survey.Drafts[i].Marks.FwdPort.Method = types.ReadingMethod(fwdPortMethod)
		}
		midPort, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.MidPort, i))
		if err == nil {
			survey.Drafts[i].Marks.MidPort.Value = midPort
		}
		midPortMethod, err := parseString(c, fmt.Sprintf("%s-d%d", constants.MidPortMarkRead, i))
		if err == nil {
			survey.Drafts[i].Marks.MidPort.Method = types.ReadingMethod(midPortMethod)
		}
		aftPort, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.AftPort, i))
		if err == nil {
			survey.Drafts[i].Marks.AftPort.Value = aftPort
		}
		aftPortMethod, err := parseString(c, fmt.Sprintf("%s-d%d", constants.AftPortMarkRead, i))
		if err == nil {
			survey.Drafts[i].Marks.AftPort.Method = types.ReadingMethod(aftPortMethod)
		}
		fwdStbd, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.FwdStbd, i))
		if err == nil {
			survey.Drafts[i].Marks.FwdStarboard.Value = fwdStbd
		}
		fwdStbdMethod, err := parseString(c, fmt.Sprintf("%s-d%d", constants.FwdStbdMarkRead, i))
		if err == nil {
			survey.Drafts[i].Marks.FwdStarboard.Method = types.ReadingMethod(fwdStbdMethod)
		}
		midStbd, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.MidStbd, i))
		if err == nil {
			survey.Drafts[i].Marks.MidStarboard.Value = midStbd
		}
		midStbdMethod, err := parseString(c, fmt.Sprintf("%s-d%d", constants.MidStbdMarkRead, i))
		if err == nil {
			survey.Drafts[i].Marks.MidStarboard.Method = types.ReadingMethod(midStbdMethod)
		}
		aftStbd, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.AftStbd, i))
		if err == nil {
			survey.Drafts[i].Marks.AftStarboard.Value = aftStbd
		}
		aftStbdMethod, err := parseString(c, fmt.Sprintf("%s-d%d", constants.AftStbdMarkRead, i))
		if err == nil {
			survey.Drafts[i].Marks.AftStarboard.Method = types.ReadingMethod(aftStbdMethod)
		}

		//Getting sea condition
		seaType, err := parseString(c, fmt.Sprintf("%s-d%d", constants.SeaType, i))
		if err == nil {
			survey.Drafts[i].SeaCondition.Type = types.SeaConditionType(seaType)
		}
		waveCondition, err := parseString(c, fmt.Sprintf("%s-d%d", constants.SeaConditionWave, i))
		if err == nil {
			survey.Drafts[i].SeaCondition.Wave = types.WaveCondition(waveCondition)
		}
		iceCondition, err := parseString(c, fmt.Sprintf("%s-d%d", constants.SeaConditionIce, i))
		if err == nil {
			survey.Drafts[i].SeaCondition.Ice = types.IceCondition(iceCondition)
		}

		//Getting deductibles
		hfo, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.HFO, i))
		if err == nil {
			survey.Drafts[i].Deductibles.HFO = hfo

		}
		mdo, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.MDO, i))
		if err == nil {
			survey.Drafts[i].Deductibles.MDO = mdo

		}
		lubOil, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.LubOil, i))
		if err == nil {
			survey.Drafts[i].Deductibles.LubOil = lubOil

		}
		bilgeWater, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.BilgeWater, i))
		if err == nil {
			survey.Drafts[i].Deductibles.BilgeWater = bilgeWater

		}
		others, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.Others, i))
		if err == nil {
			survey.Drafts[i].Deductibles.Others = others
		}

		//Getting water density
		dockWaterDensity, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.DockwaterDensity, i))
		if err == nil {
			survey.Drafts[i].Density = dockWaterDensity
		}

		//Getting vessel's passport data
		distancePPFwd, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.DFwd, i))
		if err == nil {
			survey.Drafts[i].DistancePPFwd = distancePPFwd

		}
		dirPPFwd, err := parseString(c, fmt.Sprintf("%s-d%d", constants.DFwdDir, i))
		if err == nil {
			survey.Drafts[i].PPFwdDirection = dirPPFwd
		}
		distancePPMid, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.DMid, i))
		if err == nil {
			survey.Drafts[i].DistancePPMid = distancePPMid
		}
		dirPPMid, err := parseString(c, fmt.Sprintf("%s-d%d", constants.DMidDir, i))
		if err == nil {
			survey.Drafts[i].PPMidDirection = dirPPMid
		}
		distancePPAft, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.DAft, i))
		if err == nil {
			survey.Drafts[i].DistancePPAft = distancePPAft
		}
		dirPPAft, err := parseString(c, fmt.Sprintf("%s-d%d", constants.DAftDir, i))
		if err == nil {
			survey.Drafts[i].PPAftDirection = dirPPAft
		}
		keelFwd, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.KeelFwd, i))
		if err == nil {
			survey.Drafts[i].KeelFwd = keelFwd
		}
		keelMid, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.KeelMid, i))
		if err == nil {
			survey.Drafts[i].KeelMid = keelMid
		}
		keelAft, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.KeelAft, i))
		if err == nil {
			survey.Drafts[i].KeelAft = keelAft
		}

		//Getting hydrostatics data
		if len(survey.Drafts[i].HydrostaticRows) == 0 {
			survey.Drafts[i].HydrostaticRows = make([]types.HydrostaticRow, 2)
		}
		if len(survey.Drafts[i].MTCRows) == 0 {
			survey.Drafts[i].MTCRows = make([]types.MTCRow, 2)
		}
		uDraft, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.UDraft, i))
		if err == nil {
			survey.Drafts[i].HydrostaticRows[0].Draft = uDraft
		}
		uDisp, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.UDisp, i))
		if err == nil {
			survey.Drafts[i].HydrostaticRows[0].Displacement = uDisp
		}
		uTpc, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.UTpc, i))
		if err == nil {
			survey.Drafts[i].HydrostaticRows[0].TPC = uTpc
		}
		uLcfLca, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.ULcfLca, i))
		if err == nil {
			survey.Drafts[i].HydrostaticRows[0].LCF = uLcfLca
		}
		uLcfDir, err := parseString(c, fmt.Sprintf("%s-d%d", constants.ULcfDir, i))
		if err == nil {
			survey.Drafts[i].HydrostaticRows[0].LCFDirection = types.LCFDirection(uLcfDir)
		}
		lDraft, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.LDraft, i))
		if err == nil {
			survey.Drafts[i].HydrostaticRows[1].Draft = lDraft
		}
		lDisp, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.LDisp, i))
		if err == nil {
			survey.Drafts[i].HydrostaticRows[1].Displacement = lDisp
		}
		lTpc, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.LTpc, i))
		if err == nil {
			survey.Drafts[i].HydrostaticRows[1].TPC = lTpc
		}
		lLcfLca, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.LLcfLca, i))
		if err == nil {
			survey.Drafts[i].HydrostaticRows[1].LCF = lLcfLca
		}
		lLcfDir, err := parseString(c, fmt.Sprintf("%s-d%d", constants.LLcfDir, i))
		if err == nil {
			survey.Drafts[i].HydrostaticRows[1].LCFDirection = types.LCFDirection(lLcfDir)
		}
		pMtcDraft, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.PMtcDraft, i))
		if err == nil {
			survey.Drafts[i].MTCRows[0].Draft = pMtcDraft
		}
		pMtc, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.PMtc, i))
		if err == nil {
			survey.Drafts[i].MTCRows[0].MTC = pMtc
		}
		nMtcDraft, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.NMtcDraft, i))
		if err == nil {
			survey.Drafts[i].MTCRows[1].Draft = nMtcDraft
		}
		nMtc, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.NMtc, i))
		if err == nil {
			survey.Drafts[i].MTCRows[1].MTC = nMtc
		}
	}
}

func getTankProps(h *Handler, c *fiber.Ctx) (*tankProps, error) {
	id := c.Params("id")
	survey, err := h.surveyRepository.Get(id)
	if err != nil {
		return nil, err
	}

	draftIndexStr := c.Params("draftIndex")
	draftIndex, err := strconv.Atoi(draftIndexStr)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepository.Get()
	if err != nil {
		return nil, err
	}

	var trueTrim, listDegrees *float64
	draft := survey.Drafts[draftIndex]
	if draft.Marks.FwdPort.Value != nil && draft.Marks.FwdPort.Method != "" &&
		draft.Marks.AftPort.Value != nil && draft.Marks.AftPort.Method != "" &&
		draft.Marks.MidPort.Value != nil && draft.Marks.MidPort.Method != "" &&
		draft.DistancePPFwd != nil && draft.DistancePPMid != nil && draft.DistancePPAft != nil &&
		survey.VesselData.LBP > 0 {
		results := calculation.CalcDraft(draft, survey.VesselData)
		trueTrim = &results.TrueTrim
		listDegrees = &results.ListDegrees
	}
	var trimDir, listDir string
	if trueTrim != nil && listDegrees != nil {
		trimDir = format.TrimDirection(*trueTrim)
		listDir = format.ListDirection(*listDegrees)
	}
	bwTanks := survey.Drafts[draftIndex].BallastWaterTanks
	fwTanks := survey.Drafts[draftIndex].FreshWaterTanks

	tankID := c.Params("tankID")

	if tankID == "" {
		return &tankProps{
			surveyID:   id,
			draftIndex: draftIndex,
			survey:     survey,
			user:       user,
			bwTanks:    bwTanks,
			fwTanks:    fwTanks,
			trim:       trueTrim,
			trimDir:    trimDir,
			list:       listDegrees,
			listDir:    listDir,
		}, nil
	}

	wtType := c.Get(constants.HXTankType)
	var tankIndex int

	switch wtType {
	case constants.BWTank:
		tankIndex = slices.IndexFunc(bwTanks, func(tank types.Tank) bool {
			return tank.ID == tankID
		})
		if tankIndex == -1 {
			return nil, errors.New("tank not found")
		}
	case constants.FWTank:
		tankIndex = slices.IndexFunc(fwTanks, func(tank types.Tank) bool {
			return tank.ID == tankID
		})
		if tankIndex == -1 {
			return nil, errors.New("tank not found")
		}
	default:
		return nil, errors.New("undefined tank type")
	}

	return &tankProps{
		surveyID:   id,
		draftIndex: draftIndex,
		survey:     survey,
		user:       user,
		tankID:     tankID,
		tankIndex:  tankIndex,
		bwTanks:    bwTanks,
		fwTanks:    fwTanks,
		trim:       trueTrim,
		trimDir:    trimDir,
		list:       listDegrees,
		listDir:    listDir,
	}, nil
}

func parseTank(c *fiber.Ctx, tank *types.Tank) {
	if tank == nil {
		return
	}

	tName, err := parseString(c, constants.TankName)
	if err == nil {
		tank.Name = tName
	}

	tType, err := parseString(c, constants.TankType)
	if err == nil {
		tank.Type = tType
	}

	sounding, err := parseFloat(c, fmt.Sprintf("%s-%s", constants.TankSounding, tank.ID))
	if err == nil {
		tank.Sounding = sounding
	}

	volume, err := parseFloat(c, fmt.Sprintf("%s-%s", constants.TankVolume, tank.ID))
	if err == nil {
		tank.Volume = volume
	}

	if !tank.IsFWTTank {
		density, err := parseFloat(c, fmt.Sprintf("%s-%s", constants.TankDensity, tank.ID))
		if err == nil {
			tank.Density = density
		}
	}

	//Corrections Parsing

	tableType, err := parseString(c, "calib_type")
	if err != nil {
		return
	}
	tank.Correction.TableType = types.CalibrationTableType(tableType)

	hasListCorrection, _ := parseString(c, constants.HasListCorrection)
	tank.Correction.HasListCorrection = hasListCorrection != ""

	switch tank.Correction.TableType {
	case types.CalibrationTypeVolumeByTrim:
		parseTrimRows(c, tank)
		if tank.Correction.HasListCorrection {
			parseListRows(c, tank)
		}
		return
	case types.CalibrationTypeSoundingCorrection, types.CalibrationTypeVolumeCorrection:
		parseTrimRows(c, tank)
		parseCorrectionRows(c, tank)
		if tank.Correction.HasListCorrection {
			parseListRows(c, tank)
		}
	}
}

func parseTrimRows(c *fiber.Ctx, tank *types.Tank) {
	if tank == nil {
		return
	}

	ttl, err := parseFloat(c, constants.TTL)
	if err == nil {
		tank.Correction.TableTrimLow = ttl
	}
	ttu, err := parseFloat(c, constants.TTU)
	if err == nil {
		tank.Correction.TableTrimUpper = ttu
	}

	if len(tank.Correction.TrimRows) == 0 {
		tank.Correction.TrimRows = make([]types.CalibrationRow, 2)
	}
	LowSounding, err := parseFloat(c, constants.TrimTableTSLS)
	if err == nil {
		tank.Correction.TrimRows[0].Sounding = LowSounding
	}
	TSLVolumeLow, err := parseFloat(c, constants.TrimTableTSLVL)
	if err == nil {
		tank.Correction.TrimRows[0].VolumeLow = TSLVolumeLow
	}
	TSLVolumeUp, err := parseFloat(c, constants.TrimTableTSLVU)
	if err == nil {
		tank.Correction.TrimRows[0].VolumeUp = TSLVolumeUp
	}
	UpSounding, err := parseFloat(c, constants.TrimTableTSUS)
	if err == nil {
		tank.Correction.TrimRows[1].Sounding = UpSounding
	}
	TSUVolumeLow, err := parseFloat(c, constants.TrimTableTSUVL)
	if err == nil {
		tank.Correction.TrimRows[1].VolumeLow = TSUVolumeLow
	}
	TSUVolumeUp, err := parseFloat(c, constants.TrimTableTSUVU)
	if err == nil {
		tank.Correction.TrimRows[1].VolumeUp = TSUVolumeUp
	}
}

func parseListRows(c *fiber.Ctx, tank *types.Tank) {
	if tank == nil {
		return
	}

	tll, err := parseFloat(c, constants.TLL)
	if err == nil {
		tank.Correction.TableListLow = tll
	}
	tlu, err := parseFloat(c, constants.TLU)
	if err == nil {
		tank.Correction.TableListUpper = tlu
	}

	if len(tank.Correction.ListRows) == 0 {
		tank.Correction.ListRows = make([]types.CalibrationRow, 2)
	}
	LowSounding, err := parseFloat(c, constants.ListTableTSLS)
	if err == nil {
		tank.Correction.ListRows[0].Sounding = LowSounding
	}
	TSLVolumeLow, err := parseFloat(c, constants.ListTableTSLVL)
	if err == nil {
		tank.Correction.ListRows[0].VolumeLow = TSLVolumeLow
	}
	TSLVolumeUp, err := parseFloat(c, constants.ListTableTSLVU)
	if err == nil {
		tank.Correction.ListRows[0].VolumeUp = TSLVolumeUp
	}
	UpSounding, err := parseFloat(c, constants.ListTableTSUS)
	if err == nil {
		tank.Correction.ListRows[1].Sounding = UpSounding
	}
	TSUVolumeLow, err := parseFloat(c, constants.ListTableTSUVL)
	if err == nil {
		tank.Correction.ListRows[1].VolumeLow = TSUVolumeLow
	}
	TSUVolumeUp, err := parseFloat(c, constants.ListTableTSUVU)
	if err == nil {
		tank.Correction.ListRows[1].VolumeUp = TSUVolumeUp
	}
}

func parseCorrectionRows(c *fiber.Ctx, tank *types.Tank) {
	if tank == nil {
		return
	}

	if len(tank.Correction.CorrectionRows) == 0 {
		tank.Correction.CorrectionRows = make([]types.CorrectionRows, 2)
	}
	soundingLow, err := parseFloat(c, constants.CorrTableSL)
	if err == nil {
		tank.Correction.CorrectionRows[0].TableSounding = soundingLow
	}
	volumeLow, err := parseFloat(c, constants.CorrTableVL)
	if err == nil {
		tank.Correction.CorrectionRows[0].TableVolume = volumeLow
	}
	soundingUp, err := parseFloat(c, constants.CorrTableSU)
	if err == nil {
		tank.Correction.CorrectionRows[1].TableSounding = soundingUp
	}
	volumeUp, err := parseFloat(c, constants.CorrTableVU)
	if err == nil {
		tank.Correction.CorrectionRows[1].TableVolume = volumeUp
	}
}

func getDraftProps(h *Handler, c *fiber.Ctx) (*draftProps, error) {
	id := c.Params("id")
	survey, err := h.surveyRepository.Get(id)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepository.Get()
	if err != nil {
		return nil, err
	}

	draftIndexStr := c.Params("draftIndex")
	if draftIndexStr == "" {
		return &draftProps{
			user:     user,
			survey:   survey,
			surveyID: id,
		}, nil
	}

	draftIndex, err := strconv.Atoi(draftIndexStr)
	if err != nil {
		return nil, err
	}

	return &draftProps{
		user:       user,
		survey:     survey,
		surveyID:   id,
		draftIndex: draftIndex,
	}, nil

}

func parseDraftBlocks(c *fiber.Ctx, d *types.Draft, draftIndex int) {
	blockCtx := c.Get(constants.HXDraftBlock)

	switch blockCtx {
	case constants.SeaConditionBlock:
		parseSeaCondition(c, d, draftIndex)
	case constants.MarksBlock:
		parseMarks(c, d, draftIndex)
	case constants.PPKeelBlock:
		parsePPKeel(c, d, draftIndex)
	case constants.HydrostaticsBlock:
		parseHydrostatic(c, d, draftIndex)
	case constants.DeductiblesBlock:
		parseDeductibles(c, d, draftIndex)
	}
}

func parseSeaCondition(c *fiber.Ctx, d *types.Draft, draftIndex int) {
	if d == nil {
		return
	}

	seaType, err := parseString(c, fmt.Sprintf("%s-%d", constants.SeaType, draftIndex))
	if err == nil {
		d.SeaCondition.Type = types.SeaConditionType(seaType)
	}

	switch types.SeaConditionType(seaType) {
	case types.SeaConditionTypeWave:
		condition, err := parseString(c, fmt.Sprintf("%s-%d", constants.SeaCondition, draftIndex))
		if err == nil {
			d.SeaCondition.Wave = types.WaveCondition(condition)
			d.SeaCondition.Ice = ""
		}
	case types.SeaConditionTypeIce:
		condition, err := parseString(c, fmt.Sprintf("%s-%d", constants.SeaCondition, draftIndex))
		if err == nil {
			d.SeaCondition.Ice = types.IceCondition(condition)
			d.SeaCondition.Wave = ""
		}
	}
}

func parseMarks(c *fiber.Ctx, d *types.Draft, draftIndex int) {
	if d == nil {
		return
	}
	fwdPort, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.FwdPort, draftIndex))
	if err == nil {
		d.Marks.FwdPort.Value = fwdPort
	}
	fwdPortMark, err := parseString(c, fmt.Sprintf("%s-%d", constants.FwdPortMarkRead, draftIndex))
	if err == nil {
		d.Marks.FwdPort.Method = types.ReadingMethod(fwdPortMark)
	}
	midPort, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.MidPort, draftIndex))
	if err == nil {
		d.Marks.MidPort.Value = midPort
	}
	midPortMark, err := parseString(c, fmt.Sprintf("%s-%d", constants.MidPortMarkRead, draftIndex))
	if err == nil {
		d.Marks.MidPort.Method = types.ReadingMethod(midPortMark)
	}
	aftPort, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.AftPort, draftIndex))
	if err == nil {
		d.Marks.AftPort.Value = aftPort
	}
	aftPortMark, err := parseString(c, fmt.Sprintf("%s-%d", constants.AftPortMarkRead, draftIndex))
	if err == nil {
		d.Marks.AftPort.Method = types.ReadingMethod(aftPortMark)
	}
	fwdStbd, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.FwdStbd, draftIndex))
	if err == nil {
		d.Marks.FwdStarboard.Value = fwdStbd
	}
	fwdStbdMark, err := parseString(c, fmt.Sprintf("%s-%d", constants.FwdStbdMarkRead, draftIndex))
	if err == nil {
		d.Marks.FwdStarboard.Method = types.ReadingMethod(fwdStbdMark)
	}
	midStbd, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.MidStbd, draftIndex))
	if err == nil {
		d.Marks.MidStarboard.Value = midStbd
	}
	midStbdMark, err := parseString(c, fmt.Sprintf("%s-%d", constants.MidStbdMarkRead, draftIndex))
	if err == nil {
		d.Marks.MidStarboard.Method = types.ReadingMethod(midStbdMark)
	}
	aftStbd, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.AftStbd, draftIndex))
	if err == nil {
		d.Marks.AftStarboard.Value = aftStbd
	}
	aftStbdMark, err := parseString(c, fmt.Sprintf("%s-%d", constants.AftStbdMarkRead, draftIndex))
	if err == nil {
		d.Marks.AftStarboard.Method = types.ReadingMethod(aftStbdMark)
	}
}

func parsePPKeel(c *fiber.Ctx, d *types.Draft, draftIndex int) {
	if d == nil {
		return
	}
	dFwd, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.DFwd, draftIndex))
	if err == nil {
		d.DistancePPFwd = dFwd
	}
	dFwdDir, err := parseString(c, fmt.Sprintf("%s-%d", constants.DFwdDir, draftIndex))
	if err == nil {
		d.PPFwdDirection = dFwdDir
	}
	dMid, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.DMid, draftIndex))
	if err == nil {
		d.DistancePPMid = dMid
	}
	dMidDir, err := parseString(c, fmt.Sprintf("%s-%d", constants.DMidDir, draftIndex))
	if err == nil {
		d.PPMidDirection = dMidDir
	}
	dAft, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.DAft, draftIndex))
	if err == nil {
		d.DistancePPAft = dAft
	}
	dAftDir, err := parseString(c, fmt.Sprintf("%s-%d", constants.DAftDir, draftIndex))
	if err == nil {
		d.PPAftDirection = dAftDir
	}
	keelFwd, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.KeelFwd, draftIndex))
	if err == nil {
		d.KeelFwd = keelFwd
	}
	keelMid, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.KeelMid, draftIndex))
	if err == nil {
		d.KeelMid = keelMid
	}
	keelAft, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.KeelAft, draftIndex))
	if err == nil {
		d.KeelAft = keelAft
	}
}

func parseHydrostatic(c *fiber.Ctx, d *types.Draft, draftIndex int) {
	if d == nil {
		return
	}
	uDraft, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.UDraft, draftIndex))
	if err == nil {
		d.HydrostaticRows[0].Draft = uDraft
	}
	uDisp, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.UDisp, draftIndex))
	if err == nil {
		d.HydrostaticRows[0].Displacement = uDisp
	}
	uTpc, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.UTpc, draftIndex))
	if err == nil {
		d.HydrostaticRows[0].TPC = uTpc
	}
	uLcfLca, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.ULcfLca, draftIndex))
	if err == nil {
		d.HydrostaticRows[0].LCF = uLcfLca
	}
	uLcfLcaDir, err := parseString(c, fmt.Sprintf("%s-%d", constants.ULcfDir, draftIndex))
	if err == nil {
		d.HydrostaticRows[0].LCFDirection = types.LCFDirection(uLcfLcaDir)
	}
	lDraft, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.LDraft, draftIndex))
	if err == nil {
		d.HydrostaticRows[1].Draft = lDraft
	}
	lDisp, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.LDisp, draftIndex))
	if err == nil {
		d.HydrostaticRows[1].Displacement = lDisp
	}
	lTpc, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.LTpc, draftIndex))
	if err == nil {
		d.HydrostaticRows[1].TPC = lTpc
	}
	lLcfLca, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.LLcfLca, draftIndex))
	if err == nil {
		d.HydrostaticRows[1].LCF = lLcfLca
	}
	lLcfLcaDir, err := parseString(c, fmt.Sprintf("%s-%d", constants.LLcfDir, draftIndex))
	if err == nil {
		d.HydrostaticRows[1].LCFDirection = types.LCFDirection(lLcfLcaDir)
	}
	pMtcDraft, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.PMtcDraft, draftIndex))
	if err == nil {
		d.MTCRows[0].Draft = pMtcDraft
	}
	pMtc, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.PMtc, draftIndex))
	if err == nil {
		d.MTCRows[0].MTC = pMtc
	}
	nMtcDraft, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.NMtcDraft, draftIndex))
	if err == nil {
		d.MTCRows[1].Draft = nMtcDraft
	}
	nMtc, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.NMtc, draftIndex))
	if err == nil {
		d.MTCRows[1].MTC = nMtc
	}
}

func parseDeductibles(c *fiber.Ctx, d *types.Draft, draftIndex int) {
	if d == nil {
		return
	}
	hfo, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.HFO, draftIndex))
	if err == nil {
		d.Deductibles.HFO = hfo
	}
	mdo, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.MDO, draftIndex))
	if err == nil {
		d.Deductibles.MDO = mdo
	}
	lub, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.LubOil, draftIndex))
	if err == nil {
		d.Deductibles.LubOil = lub
	}
	bglW, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.BilgeWater, draftIndex))
	if err == nil {
		d.Deductibles.BilgeWater = bglW
	}
	dwDen, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.DockwaterDensity, draftIndex))
	if err == nil {
		d.Density = dwDen
	}
	others, err := parseFloat(c, fmt.Sprintf("%s-%d", constants.Others, draftIndex))
	if err == nil {
		d.Deductibles.Others = others
	}
}

func getSurveyProps(h *Handler, c *fiber.Ctx) (*surveyProps, error) {
	id := c.Params("id")
	survey, err := h.surveyRepository.Get(id)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepository.Get()
	if err != nil {
		return nil, err
	}

	return &surveyProps{
		user:     user,
		survey:   survey,
		surveyID: id,
	}, nil
}

func parseSurveyPage(c *fiber.Ctx, s *types.Survey) {
	blockCtx := c.Get(constants.HXSurveyBlock)

	switch blockCtx {
	case constants.SurveyBlock:
		parseSurveyBlock(c, s)
	case constants.VesselBlock:
		parseVesselBlock(c, s)
	case constants.SettingsBlock:
		parseSettingBlock(c, s)
	}
}

func parseSurveyBlock(c *fiber.Ctx, s *types.Survey) {
	jobNumber, err := parseString(c, constants.JobNumber)
	if err == nil {
		s.Job.JobNumber = jobNumber
	}

	dsNumber, err := parseInt(c, constants.DSNumber)
	if err == nil {
		s.Job.DSNumber = *dsNumber
	}

	operation, err := parseString(c, constants.CargoOperation)
	if err == nil {
		s.CargoOperation.Operation = operation
	}

	cargo, err := parseString(c, constants.Cargo)
	if err == nil {
		s.CargoOperation.Cargo = cargo
	}

	packing, err := parseString(c, constants.Packing)
	if err == nil {
		s.CargoOperation.Packing = packing
	}

	principal, err := parseString(c, constants.Client)
	if err == nil {
		s.Job.Principal = principal
	}

	cargoDeclared, err := parseFloat(c, constants.CargoDeclared)
	if err == nil {
		s.CargoDeclared = cargoDeclared
	}

	remarks, err := parseString(c, constants.Remarks)
	if err == nil {
		s.Remarks = remarks
	}

	countryCode, err := parseString(c, constants.CountryCode)
	if err == nil {
		s.Country.CountryCode = countryCode
	}

	countryName, err := parseString(c, constants.CountryName)
	if err == nil {
		s.Country.Name = countryName
	}
}

func parseVesselBlock(c *fiber.Ctx, s *types.Survey) {
	vesselName, err := parseString(c, constants.VesselName)
	if err == nil {
		s.VesselData.Name = vesselName
	}

	imo, err := parseString(c, constants.IMO)
	if err == nil {
		s.VesselData.IMO = imo
	}

	builtYear, err := parseInt(c, constants.Built)
	if err == nil {
		s.VesselData.BuiltYear = *builtYear
	}

	tableDensity, err := parseFloat(c, constants.TableDensity)
	if err == nil {
		s.VesselData.TableDensity = tableDensity
	}

	lightship, err := parseFloat(c, constants.Lightship)
	if err == nil {
		s.VesselData.Lightship = *lightship
	}

	lbp, err := parseFloat(c, constants.LBP)
	if err == nil {
		s.VesselData.LBP = *lbp
	}

	breadth, err := parseFloat(c, constants.Breadth)
	if err == nil {
		s.VesselData.Breadth = *breadth
	}

	depth, err := parseFloat(c, constants.Depth)
	if err == nil {
		s.VesselData.Depth = *depth
	}

	summerDraft, err := parseFloat(c, constants.SummerDraft)
	if err == nil {
		s.VesselData.SummerDraft = *summerDraft
	}

	summerDWT, err := parseFloat(c, constants.SummerDWT)
	if err == nil {
		s.VesselData.SummerDWT = *summerDWT
	}

	summerTPC, err := parseFloat(c, constants.SummerTPC)
	if err == nil {
		s.VesselData.SummerTPC = *summerTPC
	}

	summerFreeboard, err := parseFloat(c, constants.SummerFreeboard)
	if err == nil {
		s.VesselData.SummerFreeboard = *summerFreeboard
	}

	countryCode, err := parseString(c, constants.CountryCode)
	if err == nil {
		s.VesselData.Flag.CountryCode = countryCode
	}

	countryName, err := parseString(c, constants.CountryName)
	if err == nil {
		s.VesselData.Flag.Name = countryName
	}

	holdsTotal, err := parseInt(c, constants.HoldsTotal)
	if err == nil {
		s.VesselData.HoldsTotal = *holdsTotal
	}
}

func parseSettingBlock(c *fiber.Ctx, s *types.Survey) {
	mmcMethod, err := parseString(c, constants.MMCMethod)
	if err == nil {
		s.VesselData.VesselType = types.VesselType(mmcMethod)
	}

	corrMethod, err := parseString(c, constants.CorrMethod)
	if err == nil {
		s.VesselData.CorrectionMethod = types.CorrectionMethod(corrMethod)
	}

	lcfDetection, err := parseString(c, constants.LCFDetection)
	if err == nil {
		s.VesselData.IsLcfDetectionManual = lcfDetection == "manual"
	}
}
