package handler

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/AVZotov/draft-survey/internal/constants"
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/gofiber/fiber/v2"
)

var ErrEmptyField = errors.New("empty field")

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
		constDeclared, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.ConstDeclared, i))
		if err == nil {
			survey.Drafts[i].ConstantDeclared = constDeclared
		}
		cargoDeclared, err := parseFloat(c, fmt.Sprintf("%s-d%d", constants.CargoDeclared, i))
		if err == nil {
			survey.Drafts[i].CargoDeclared = cargoDeclared
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

func (h *Handler) parseBwTank(c *fiber.Ctx, tank *types.Tank) {
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

	density, err := parseFloat(c, fmt.Sprintf("%s-%s", constants.TankDensity, tank.ID))
	if err == nil {
		tank.Density = density
	}
}
