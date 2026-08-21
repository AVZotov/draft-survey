//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/AVZotov/draft-survey/internal/handler/fields"
	"github.com/AVZotov/draft-survey/internal/sse"
)

// draftExpected holds the subset of docs/TESTING_CHEATSHEET.md's "Expected
// Results" table that Step 2 asks to assert.
type draftExpected struct {
	MMC              float64
	LBM              float64
	TrueTrim         float64
	NetDisplacement  float64
	TotalDeductibles float64
}

// vesselFixture bundles one cheatsheet vessel end to end: vessel/job data,
// both drafts' fields, both drafts' tanks, and the expected results for
// each draft plus the survey-level CargoOnBoard for the final draft.
type vesselFixture struct {
	name string

	vessel map[string]string
	job    map[string]string

	draft0       map[string]string
	draft0BW     []tankFixture
	draft0FW     []tankFixture
	draft0Want   draftExpected
	draft1       map[string]string
	draft1BW     []tankFixture
	draft1FW     []tankFixture
	draft1Want   draftExpected
	cargoOnBoard float64
}

// runDraftCalcTest drives one vessel fixture through the real HTTP pipeline
// and asserts the Step 2 fields.
//
// Deviation flagged: the task's 4-step outline for this test only mentions
// entering draft 0's fields, but the assertions it asks for (NetDisplacement,
// TotalDeductibles, and especially CargoFromPrev "draft 1 vs draft 0") are
// mathematically dependent on BOTH drafts' full tank data — TotalDeductibles
// = TotalBwTanksWeight + TotalFwTanksWeight + fuel/water deductibles (see
// CalcTotalDeductibles), and CargoFromPrev/CargoOnBoard is a delta of
// NetDisplacement between drafts. So this drives both drafts (marks, tanks,
// everything) and adds the second draft via addDraft(), rather than stopping
// at draft 0.
func runDraftCalcTest(t *testing.T, vf vesselFixture) {
	ts := newTestServer(t)

	id := ts.createSurvey(t)
	ts.putSurvey(t, id, mergeMaps(vf.vessel, vf.job))

	// --- Draft 0 --------------------------------------------------------
	ts.startDraft(t, id, 0)
	ts.putDraft(t, id, 0, vf.draft0)
	for _, tf := range vf.draft0BW {
		ts.addBWTank(t, id, 0, tf)
	}
	for _, tf := range vf.draft0FW {
		ts.addFWTank(t, id, 0, tf)
	}

	// --- Draft 1 ----------------------------------------------------------
	ts.addDraft(t, id)
	ts.startDraft(t, id, 1)
	ts.putDraft(t, id, 1, vf.draft1)
	for _, tf := range vf.draft1BW {
		ts.addBWTank(t, id, 1, tf)
	}
	for _, tf := range vf.draft1FW {
		ts.addFWTank(t, id, 1, tf)
	}

	// Tanks are added via /tanks/... endpoints, which publish EventTankCalc,
	// not EventDraftCalc — so no EventDraftCalc broadcast yet reflects the
	// final BW/FW totals for either draft. publishDraftCalc always
	// recalculates from the currently-saved survey (both drafts' tanks are
	// already persisted at this point), so the cheapest way to get one is to
	// trigger one more (otherwise no-op) draft save.
	sseData := ts.captureSSEEvent(t, sse.EventDraftCalc, 5*time.Second, func() {
		ts.putDraft(t, id, 1, map[string]string{fields.FieldDockwaterDensity: vf.draft1[fields.FieldDockwaterDensity]})
	})

	assertNear(t, vf.name+" draft0 MMC", extractMMC(t, sseData, 0), vf.draft0Want.MMC)
	assertNear(t, vf.name+" draft0 LBM", extractLBM(t, sseData, 0), vf.draft0Want.LBM)
	assertNear(t, vf.name+" draft0 TrueTrim", extractTrueTrim(t, sseData, 0), vf.draft0Want.TrueTrim)
	assertNear(t, vf.name+" draft0 NetDisplacement", extractNetDisplacement(t, sseData, 0), vf.draft0Want.NetDisplacement)

	assertNear(t, vf.name+" draft1 MMC", extractMMC(t, sseData, 1), vf.draft1Want.MMC)
	assertNear(t, vf.name+" draft1 LBM", extractLBM(t, sseData, 1), vf.draft1Want.LBM)
	assertNear(t, vf.name+" draft1 TrueTrim", extractTrueTrim(t, sseData, 1), vf.draft1Want.TrueTrim)
	assertNear(t, vf.name+" draft1 NetDisplacement", extractNetDisplacement(t, sseData, 1), vf.draft1Want.NetDisplacement)

	// CargoFromPrev (draft 1 vs draft 0) — see extractCargoOnBoard's doc
	// comment in helpers_test.go for why this reads CargoOnBoard instead.
	assertNear(t, vf.name+" CargoOnBoard (== draft1 CargoFromPrev)", extractCargoOnBoard(t, sseData, 1), vf.cargoOnBoard)

	// TotalDeductibles is not in the SSE payload (see extractTotalDeductibles's
	// doc comment) — read it off the full Draft Readings page instead.
	pageHTML := ts.getDraftPage(t, id)
	assertNear(t, vf.name+" draft0 TotalDeductibles", extractTotalDeductibles(t, pageHTML, 0), vf.draft0Want.TotalDeductibles)
	assertNear(t, vf.name+" draft1 TotalDeductibles", extractTotalDeductibles(t, pageHTML, 1), vf.draft1Want.TotalDeductibles)
}

func mergeMaps(maps ...map[string]string) map[string]string {
	out := map[string]string{}
	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

// --- AGGELOS B ---------------------------------------------------------

func aggelosBFixture() vesselFixture {
	vessel := map[string]string{
		fields.FieldVesselName:      "AGGELOS B",
		fields.FieldIMO:             "9577434",
		fields.FieldFlagCountryCode: "MT",
		fields.FieldBuiltYear:       "0",
		fields.FieldLBP:             "189",
		fields.FieldBreadth:         "32.26",
		fields.FieldDepth:           "18.6",
		fields.FieldSummerDraft:     "13.02",
		fields.FieldSummerDWT:       "58479.8",
		fields.FieldSummerFreeboard: "5.621",
		fields.FieldSummerTPC:       "59",
		fields.FieldLightship:       "11440",
		fields.FieldHoldsTotal:      "5",
		fields.FieldTableDensity:    "1.025",
		fields.FieldMMCMethod:       "marine",
		fields.FieldCorrMethod:      "Full LBP",
		fields.FieldLCFDetection:    "manual",
	}
	job := map[string]string{
		fields.FieldJobNumber:     "423432_001",
		fields.FieldDSNumber:      "1",
		fields.FieldClient:        "Our principal",
		fields.FieldCargoDeclared: "53000",
		fields.FieldConstDeclared: "350",
	}

	draft0 := map[string]string{
		fields.FieldFwdPort: "3.38", fields.FieldFwdPortMethod: "direct",
		fields.FieldFwdStbd: "3.39", fields.FieldFwdStbdMethod: "direct",
		fields.FieldMidPort: "4.48", fields.FieldMidPortMethod: "direct",
		fields.FieldMidStbd: "4.53", fields.FieldMidStbdMethod: "direct",
		fields.FieldAftPort: "5.53", fields.FieldAftPortMethod: "direct",
		fields.FieldAftStbd: "5.53", fields.FieldAftStbdMethod: "direct",

		fields.FieldDFwd: "2.2", fields.FieldDFwdDir: "A",
		fields.FieldDMidDir: "A",
		fields.FieldDAft:    "8.5", fields.FieldDAftDir: "F",

		fields.FieldUDraft: "4.5", fields.FieldUDisp: "22341.3", fields.FieldUTPC: "52.4", fields.FieldULCF: "7.954", fields.FieldULCFDir: "F",
		fields.FieldLDraft: "4.6", fields.FieldLDisp: "22865.8", fields.FieldLTPC: "52.5", fields.FieldLLCF: "7.896", fields.FieldLLCFDir: "F",

		fields.FieldPMTCDraft: "5", fields.FieldPMTC: "616.9",
		fields.FieldNMTCDraft: "4", fields.FieldNMTC: "597.6",

		fields.FieldDockwaterDensity: "1.024",
		fields.FieldTPCListPort:      "52.2", fields.FieldTPCListStbd: "52.6",

		fields.FieldSeaType: "wave", fields.FieldSeaCondition: "under 0.1m",

		fields.FieldHFO: "657.17", fields.FieldMDO: "28.7",
	}
	draft0BW := []tankFixture{
		{Type: "FPT", Sounding: "11.2", Volume: "1374.73", Density: "1.025"},
		{Type: "WBT", Name: "1P", Sounding: "2.23", Volume: "646.81", Density: "1.025"},
		{Type: "WBT", Name: "1S", Sounding: "2.15", Volume: "639.31", Density: "1.025"},
		{Type: "WBT", Name: "2P", Sounding: "2.08", Volume: "826.55", Density: "1.025"},
		{Type: "WBT", Name: "2S", Sounding: "2.1", Volume: "829.23", Density: "1.025"},
		{Type: "WBT", Name: "3P", Sounding: "2.15", Volume: "837.34", Density: "1.025"},
		{Type: "WBT", Name: "3S", Sounding: "2.1", Volume: "831.45", Density: "1.025"},
		{Type: "WBT", Name: "4P", Sounding: "2.25", Volume: "847.82", Density: "1.025"},
		{Type: "WBT", Name: "4S", Sounding: "2.18", Volume: "839.51", Density: "1.025"},
		{Type: "WBT", Name: "5P", Sounding: "2.15", Volume: "671.22", Density: "1.025"},
		{Type: "WBT", Name: "5S", Sounding: "2.13", Volume: "668.21", Density: "1.025"},
		{Type: "APT", Sounding: "0", Volume: "0.3", Density: "1.025"},
	}
	draft0FW := []tankFixture{
		{Type: "FPT", Name: "FWT P", Volume: "73"},
		{Type: "FPT", Name: "FWT S", Volume: "77"},
	}
	draft0Want := draftExpected{MMC: 4.5030, LBM: 178.3000, TrueTrim: 2.2730, NetDisplacement: 11787.3690, TotalDeductibles: 10073.6630}

	draft1 := map[string]string{
		fields.FieldFwdPort: "12.12", fields.FieldFwdPortMethod: "direct",
		fields.FieldFwdStbd: "12.12", fields.FieldFwdStbdMethod: "direct",
		fields.FieldMidPort: "12.34", fields.FieldMidPortMethod: "direct",
		fields.FieldMidStbd: "12.4", fields.FieldMidStbdMethod: "direct",
		fields.FieldAftPort: "12.39", fields.FieldAftPortMethod: "direct",
		fields.FieldAftStbd: "12.41", fields.FieldAftStbdMethod: "direct",

		fields.FieldDFwd: "2.2", fields.FieldDFwdDir: "A",
		fields.FieldDMidDir: "A",
		fields.FieldDAft:    "8.5", fields.FieldDAftDir: "F",

		fields.FieldUDraft: "12.3", fields.FieldUDisp: "65685.2", fields.FieldUTPC: "58.6", fields.FieldULCF: "1.42", fields.FieldULCFDir: "A",
		fields.FieldLDraft: "12.4", fields.FieldLDisp: "66271.8", fields.FieldLTPC: "58.7", fields.FieldLLCF: "1.47", fields.FieldLLCFDir: "A",

		fields.FieldPMTCDraft: "12.8", fields.FieldPMTC: "840.6",
		fields.FieldNMTCDraft: "11.8", fields.FieldNMTC: "820.6",

		fields.FieldDockwaterDensity: "1.024",

		fields.FieldSeaType: "wave", fields.FieldSeaCondition: "0.1-0.5m",

		fields.FieldHFO: "652.09", fields.FieldMDO: "28.67",
	}
	draft1BW := []tankFixture{
		{Type: "FPT", Sounding: "0.29", Volume: "11.66", Density: "1.025"},
		{Type: "WBT", Name: "1P", Sounding: "0.09", Volume: "22.46", Density: "1.025"},
		{Type: "WBT", Name: "1S", Sounding: "0", Volume: "0", Density: "1.025"},
		{Type: "WBT", Name: "2P", Sounding: "0.2", Volume: "81.42", Density: "1.025"},
		{Type: "WBT", Name: "2S", Sounding: "0", Volume: "0", Density: "1.025"},
		{Type: "WBT", Name: "3P", Sounding: "0.09", Volume: "32.151", Density: "1.025"},
		{Type: "WBT", Name: "3S", Sounding: "0", Volume: "0", Density: "1.025"},
		{Type: "WBT", Name: "4P", Sounding: "0.08", Volume: "28.26", Density: "1.025"},
		{Type: "WBT", Name: "4S", Sounding: "0.06", Volume: "22.6", Density: "1.025"},
		{Type: "WBT", Name: "5P", Sounding: "0.18", Volume: "63.47", Density: "1.025"},
		{Type: "WBT", Name: "5S", Sounding: "0.03", Volume: "15.46", Density: "1.025"},
		{Type: "APT", Sounding: "0", Volume: "0.23", Density: "1.025"},
	}
	draft1FW := []tankFixture{
		{Type: "FPT", Name: "FWT P", Volume: "70"},
		{Type: "FPT", Name: "FWT S", Volume: "70"},
	}
	draft1Want := draftExpected{MMC: 12.3440, LBM: 178.3000, TrueTrim: 0.2960, NetDisplacement: 64787.2480, TotalDeductibles: 1105.4160}

	return vesselFixture{
		name:         "AggelosB",
		vessel:       vessel,
		job:          job,
		draft0:       draft0,
		draft0BW:     draft0BW,
		draft0FW:     draft0FW,
		draft0Want:   draft0Want,
		draft1:       draft1,
		draft1BW:     draft1BW,
		draft1FW:     draft1FW,
		draft1Want:   draft1Want,
		cargoOnBoard: 52999.8790,
	}
}

func TestAggelosBDraftCalc(t *testing.T) {
	runDraftCalcTest(t, aggelosBFixture())
}

// --- HUA YOU 2 -----------------------------------------------------------

func huaYou2Fixture() vesselFixture {
	vessel := map[string]string{
		fields.FieldVesselName:      "HUA YOU 2",
		fields.FieldIMO:             "8670095",
		fields.FieldFlagCountryCode: "PA",
		fields.FieldBuiltYear:       "2010",
		fields.FieldLBP:             "178",
		fields.FieldBreadth:         "28.8",
		fields.FieldDepth:           "15.2",
		fields.FieldSummerDraft:     "10.6",
		fields.FieldSummerDWT:       "36988",
		fields.FieldSummerFreeboard: "4.62",
		fields.FieldSummerTPC:       "47.9",
		fields.FieldLightship:       "9044",
		fields.FieldHoldsTotal:      "5",
		fields.FieldTableDensity:    "1.025",
		fields.FieldMMCMethod:       "marine",
		fields.FieldCorrMethod:      "Full LBP",
		fields.FieldLCFDetection:    "manual",
	}
	job := map[string]string{
		fields.FieldJobNumber:     "449258_004",
		fields.FieldDSNumber:      "1",
		fields.FieldClient:        "Rusaltrans OOO",
		fields.FieldCargoDeclared: "30467",
		fields.FieldConstDeclared: "300",
	}

	draft0 := map[string]string{
		fields.FieldFwdPort: "9.51", fields.FieldFwdPortMethod: "direct",
		fields.FieldFwdStbd: "9.55", fields.FieldFwdStbdMethod: "direct",
		fields.FieldMidPort: "9.53", fields.FieldMidPortMethod: "direct",
		fields.FieldMidStbd: "9.65", fields.FieldMidStbdMethod: "direct",
		fields.FieldAftPort: "9.58", fields.FieldAftPortMethod: "direct",
		fields.FieldAftStbd: "9.6", fields.FieldAftStbdMethod: "direct",

		fields.FieldDFwd: "0.6", fields.FieldDFwdDir: "A",
		fields.FieldDMid: "0.605", fields.FieldDMidDir: "A",
		fields.FieldDAftDir: "A",

		fields.FieldUDraft: "9.58", fields.FieldUDisp: "41105.4", fields.FieldUTPC: "47", fields.FieldULCF: "1.162", fields.FieldULCFDir: "A",
		fields.FieldLDraft: "9.6", fields.FieldLDisp: "41199.4", fields.FieldLTPC: "47", fields.FieldLLCF: "1.198", fields.FieldLLCFDir: "A",

		fields.FieldPMTCDraft: "10.08", fields.FieldPMTC: "593.6",
		fields.FieldNMTCDraft: "9.08", fields.FieldNMTC: "561",

		fields.FieldDockwaterDensity: "1.02",

		fields.FieldSeaType: "wave", fields.FieldSeaCondition: "under 0.1m",

		fields.FieldHFO: "141.36", fields.FieldMDO: "47.46",
	}
	draft0BW := []tankFixture{
		{Type: "FPT", Sounding: "0.9", Volume: "27.755", Density: "1.025"},
		{Type: "FPT", Name: "P", Sounding: "0.32", Volume: "4.315", Density: "1.025"},
		{Type: "FPT", Name: "S", Sounding: "0.38", Volume: "4.891", Density: "1.025"},
		{Type: "DBT", Name: "1P", Sounding: "0", Volume: "0", Density: "1.025"},
		{Type: "DBT", Name: "1S", Sounding: "0.32", Volume: "73.815", Density: "1.025"},
		{Type: "DBT", Name: "2P", Sounding: "0.04", Volume: "11.965", Density: "1.025"},
		{Type: "DBT", Name: "3P", Sounding: "0.03", Volume: "8.908", Density: "1.025"},
		{Type: "DBT", Name: "3S", Sounding: "0.02", Volume: "5.939", Density: "1.025"},
		{Type: "DBT", Name: "4P", Sounding: "0.61", Volume: "202.034", Density: "1.025"},
		{Type: "DBT", Name: "4S", Sounding: "0.14", Volume: "43.909", Density: "1.025"},
		{Type: "DBT", Name: "5P", Sounding: "0.51", Volume: "108.925", Density: "1.025"},
		{Type: "DBT", Name: "5S", Sounding: "0", Volume: "0", Density: "1.025"},
		{Type: "TST", Name: "1P", Sounding: "0", Volume: "0", Density: "1.025"},
		{Type: "TST", Name: "1S", Sounding: "0.08", Volume: "1.103", Density: "1.025"},
		{Type: "TST", Name: "2P", Sounding: "0.03", Volume: "1.014", Density: "1.025"},
		{Type: "TST", Name: "2S", Sounding: "0.06", Volume: "2.029", Density: "1.025"},
		{Type: "TST", Name: "3P", Sounding: "0.03", Volume: "1.151", Density: "1.025"},
		{Type: "TST", Name: "3S", Sounding: "0", Volume: "0.554", Density: "1.025"},
		{Type: "TST", Name: "4P", Sounding: "0.09", Volume: "1.702", Density: "1.025"},
		{Type: "TST", Name: "4S", Sounding: "0.04", Volume: "1.397", Density: "1.025"},
		{Type: "APT", Name: "P", Sounding: "0", Volume: "3.395", Density: "1.025"},
		{Type: "APT", Name: "S", Sounding: "1.15", Volume: "59.894", Density: "1.025"},
		{Type: "DBT", Name: "2S", Sounding: "0.05", Volume: "11.966", Density: "1.025"},
	}
	draft0FW := []tankFixture{
		{Type: "FPT", Name: "FW P", Volume: "65"},
		{Type: "FPT", Name: "FW S", Volume: "68.831"},
	}
	draft0Want := draftExpected{MMC: 9.5830, LBM: 177.4000, TrueTrim: 0.0600, NetDisplacement: 40007.1370, TotalDeductibles: 913.7290}

	draft1 := map[string]string{
		fields.FieldFwdPort: "2.6", fields.FieldFwdPortMethod: "direct",
		fields.FieldFwdStbd: "2.58", fields.FieldFwdStbdMethod: "direct",
		fields.FieldMidPort: "4.33", fields.FieldMidPortMethod: "direct",
		fields.FieldMidStbd: "4.23", fields.FieldMidStbdMethod: "direct",
		fields.FieldAftPort: "6", fields.FieldAftPortMethod: "direct",
		fields.FieldAftStbd: "5.98", fields.FieldAftStbdMethod: "direct",

		fields.FieldDFwd: "0.6", fields.FieldDFwdDir: "A",
		fields.FieldDMid: "0.605", fields.FieldDMidDir: "A",
		fields.FieldDAftDir: "A",

		fields.FieldUDraft: "4.26", fields.FieldUDisp: "17142.3", fields.FieldUTPC: "43.2", fields.FieldULCF: "5.67", fields.FieldULCFDir: "F",
		fields.FieldLDraft: "4.28", fields.FieldLDisp: "17228.7", fields.FieldLTPC: "43.3", fields.FieldLLCF: "5.665", fields.FieldLLCFDir: "F",

		fields.FieldPMTCDraft: "4.78", fields.FieldPMTC: "473.7",
		fields.FieldNMTCDraft: "3.78", fields.FieldNMTC: "453.3",

		fields.FieldDockwaterDensity: "1.022",

		fields.FieldSeaType: "wave", fields.FieldSeaCondition: "under 0.1m",

		fields.FieldHFO: "131", fields.FieldMDO: "44",
	}
	draft1BW := []tankFixture{
		{Type: "FPT", Sounding: "1.2", Volume: "41", Density: "1.0225"},
		{Type: "FPT", Name: "P", Sounding: "0.24", Volume: "3.5", Density: "1.0225"},
		{Type: "FPT", Name: "S", Sounding: "2", Volume: "3", Density: "1.0225"},
		{Type: "DBT", Name: "1P", Sounding: "7.7", Volume: "670", Density: "1.0225"},
		{Type: "DBT", Name: "1S", Sounding: "4.94", Volume: "638", Density: "1.0225"},
		{Type: "DBT", Name: "2P", Sounding: "11.64", Volume: "769", Density: "1.0225"},
		{Type: "DBT", Name: "2S", Sounding: "6.52", Volume: "769", Density: "1.0225"},
		{Type: "DBT", Name: "3P", Sounding: "4.22", Volume: "709", Density: "1.0225"},
		{Type: "DBT", Name: "3S", Sounding: "3.11", Volume: "643", Density: "1.0225"},
		{Type: "DBT", Name: "4P", Sounding: "15", Volume: "703", Density: "1.0225"},
		{Type: "DBT", Name: "4S", Sounding: "9.7", Volume: "710", Density: "1.0225"},
		{Type: "DBT", Name: "5P", Sounding: "4.57", Volume: "640", Density: "1.0225"},
		{Type: "DBT", Name: "5S", Sounding: "4.98", Volume: "664", Density: "1.0225"},
		{Type: "TST", Name: "1P", Sounding: "0", Volume: "0", Density: "1.0225"},
		{Type: "TST", Name: "1S", Sounding: "0.07", Volume: "0.4", Density: "1.0225"},
		{Type: "TST", Name: "2P", Sounding: "0.14", Volume: "1.5", Density: "1.0225"},
		{Type: "TST", Name: "2S", Sounding: "0.07", Volume: "0.5", Density: "1.0225"},
		{Type: "TST", Name: "3P", Sounding: "0.04", Volume: "0.2", Density: "1.0225"},
		{Type: "TST", Name: "3S", Sounding: "0.04", Volume: "0.2", Density: "1.0225"},
		{Type: "TST", Name: "4P", Sounding: "0.06", Volume: "3.7", Density: "1.0225"},
		{Type: "TST", Name: "4S", Sounding: "0.04", Volume: "3.1", Density: "1.0225"},
		{Type: "APT", Name: "P", Sounding: "0", Volume: "0", Density: "1.0225"},
		{Type: "APT", Name: "S", Sounding: "1.31", Volume: "70.5", Density: "1.0225"},
	}
	draft1FW := []tankFixture{
		{Type: "FPT", Name: "FW P", Volume: "65"},
		{Type: "FPT", Name: "FW S", Volume: "65"},
	}
	draft1Want := draftExpected{MMC: 4.2720, LBM: 177.4000, TrueTrim: 3.4110, NetDisplacement: 9235.8180, TotalDeductibles: 7506.0630}

	return vesselFixture{
		name:         "HuaYou2",
		vessel:       vessel,
		job:          job,
		draft0:       draft0,
		draft0BW:     draft0BW,
		draft0FW:     draft0FW,
		draft0Want:   draft0Want,
		draft1:       draft1,
		draft1BW:     draft1BW,
		draft1FW:     draft1FW,
		draft1Want:   draft1Want,
		cargoOnBoard: -30771.3190,
	}
}

func TestHuaYou2DraftCalc(t *testing.T) {
	runDraftCalcTest(t, huaYou2Fixture())
}

// --- NEWSUN VISION ---------------------------------------------------------

func newsunVisionFixture() vesselFixture {
	vessel := map[string]string{
		fields.FieldVesselName:      "NEWSUN VISION",
		fields.FieldIMO:             "9300867",
		fields.FieldFlagCountryCode: "VN",
		fields.FieldBuiltYear:       "2007",
		fields.FieldLBP:             "170",
		fields.FieldBreadth:         "28",
		fields.FieldDepth:           "14.4",
		fields.FieldSummerDraft:     "9.812",
		fields.FieldSummerDWT:       "30548.6",
		fields.FieldSummerFreeboard: "4.588",
		fields.FieldSummerTPC:       "43.368",
		fields.FieldLightship:       "8009.04",
		fields.FieldHoldsTotal:      "5",
		fields.FieldTableDensity:    "1.025",
		fields.FieldMMCMethod:       "marine",
		fields.FieldCorrMethod:      "Full LBP",
		fields.FieldLCFDetection:    "manual",
	}
	job := map[string]string{
		fields.FieldJobNumber:     "436397_005",
		fields.FieldDSNumber:      "1",
		fields.FieldClient:        "Rusaltrans OOO",
		fields.FieldCargoDeclared: "29402",
		fields.FieldConstDeclared: "300",
	}

	draft0 := map[string]string{
		fields.FieldFwdPort: "9.66", fields.FieldFwdPortMethod: "direct",
		fields.FieldFwdStbd: "9.61", fields.FieldFwdStbdMethod: "direct",
		fields.FieldMidPort: "10.16", fields.FieldMidPortMethod: "direct",
		fields.FieldMidStbd: "9.57", fields.FieldMidStbdMethod: "direct",
		fields.FieldAftPort: "10.47", fields.FieldAftPortMethod: "direct",
		fields.FieldAftStbd: "9.99", fields.FieldAftStbdMethod: "direct",

		fields.FieldDFwd: "1.6", fields.FieldDFwdDir: "A",
		fields.FieldDMid: "0.2", fields.FieldDMidDir: "A",
		fields.FieldDAft: "8.8", fields.FieldDAftDir: "F",

		fields.FieldUDraft: "9.824", fields.FieldUDisp: "38592.3", fields.FieldUTPC: "43.5", fields.FieldULCF: "1.933", fields.FieldULCFDir: "A",
		fields.FieldLDraft: "9.924", fields.FieldLDisp: "39027.9", fields.FieldLTPC: "43.6", fields.FieldLLCF: "2.058", fields.FieldLLCFDir: "A",

		fields.FieldPMTCDraft: "10.424", fields.FieldPMTC: "526.3",
		fields.FieldNMTCDraft: "9.424", fields.FieldNMTC: "499.2",

		fields.FieldDockwaterDensity: "1.023",

		fields.FieldSeaType: "ice", fields.FieldSeaCondition: "0.1-0.15m around",

		fields.FieldHFO: "153.379", fields.FieldMDO: "32.162", fields.FieldLubOil: "0", fields.FieldBilgeWater: "0", fields.FieldOthers: "170",
	}
	draft0BW := []tankFixture{
		{Type: "FPT", Sounding: "11", Volume: "358.256", Density: "1.025"},
	}
	draft0FW := []tankFixture{
		{Type: "FPT", Volume: "257.18"},
	}
	draft0Want := draftExpected{MMC: 9.8850, LBM: 159.6000, TrueTrim: 0.6340, NetDisplacement: 37838.0350, TotalDeductibles: 979.9330}

	draft1 := map[string]string{
		fields.FieldFwdPort: "2.17", fields.FieldFwdPortMethod: "direct",
		fields.FieldFwdStbd: "2.17", fields.FieldFwdStbdMethod: "direct",
		fields.FieldMidPort: "3.69", fields.FieldMidPortMethod: "direct",
		fields.FieldMidStbd: "3.72", fields.FieldMidStbdMethod: "direct",
		fields.FieldAftPort: "5.28", fields.FieldAftPortMethod: "direct",
		fields.FieldAftStbd: "5.28", fields.FieldAftStbdMethod: "direct",

		fields.FieldDFwd: "1.6", fields.FieldDFwdDir: "A",
		fields.FieldDMid: "0.2", fields.FieldDMidDir: "A",
		fields.FieldDAft: "8.8", fields.FieldDAftDir: "F",

		fields.FieldUDraft: "3.724", fields.FieldUDisp: "13580.3", fields.FieldUTPC: "39.2", fields.FieldULCF: "5.944", fields.FieldULCFDir: "F",
		fields.FieldLDraft: "3.824", fields.FieldLDisp: "13972.6", fields.FieldLTPC: "39.2", fields.FieldLLCF: "5.921", fields.FieldLLCFDir: "F",

		fields.FieldPMTCDraft: "4.224", fields.FieldPMTC: "388.7",
		fields.FieldNMTCDraft: "3.224", fields.FieldNMTC: "374",

		fields.FieldDockwaterDensity: "1.0235",

		fields.FieldSeaType: "ice", fields.FieldSeaCondition: "0.15-0.2m around",

		fields.FieldHFO: "235.601", fields.FieldMDO: "31.342", fields.FieldLubOil: "0", fields.FieldBilgeWater: "0", fields.FieldOthers: "170",
	}
	draft1BW := []tankFixture{
		{Type: "FPT", Sounding: "22", Volume: "3933.96", Density: "1.023"},
	}
	draft1FW := []tankFixture{
		{Type: "FPT", Volume: "225.38"},
	}
	draft1Want := draftExpected{MMC: 3.7250, LBM: 159.6000, TrueTrim: 3.3120, NetDisplacement: 8471.6570, TotalDeductibles: 4686.7640}

	return vesselFixture{
		name:         "NewsunVision",
		vessel:       vessel,
		job:          job,
		draft0:       draft0,
		draft0BW:     draft0BW,
		draft0FW:     draft0FW,
		draft0Want:   draft0Want,
		draft1:       draft1,
		draft1BW:     draft1BW,
		draft1FW:     draft1FW,
		draft1Want:   draft1Want,
		cargoOnBoard: -29366.3780,
	}
}

func TestNewsunVisionDraftCalc(t *testing.T) {
	runDraftCalcTest(t, newsunVisionFixture())
}

// --- POLAR STAR ----------------------------------------------------------
//
// Data source: cmd/verify/testdata/polar_star.json, built and verified
// against DSGear (independent commercial draft-survey software) in the raw
// data audit — see the verification report. Covers two gaps that fixture
// exposed: DSGear's own hydrostatic table rows for this vessel give LCF as
// distance-from-AP (~98.5m, far larger than LBP/2=91.5m allows for a
// from-midship value) rather than from-midship, so this is the first
// fixture to exercise CalcHydrostatics' k3 auto-detection/from-AP branch
// (IsLcfDetectionManual=false — every existing fixture uses "manual" and
// direction-labelled, from-midship LCF data). It also has a nonzero List
// Correction on its initial draft (0.05), unlike every existing fixture's
// draft0/draft1 pairs which never isolated that specific value for
// assertion.
//
// polarStarFixture() reuses vesselFixture/tankFixture/draftExpected so its
// data can be driven with the same helpers as the other three vessels, but
// TestPolarStarDraftCalc does not call the shared runDraftCalcTest: that
// harness's draftExpected only covers MMC/LBM/TrueTrim/NetDisplacement/
// TotalDeductibles/CargoOnBoard, not the FirstTrimCorrection/
// SecondTrimCorrection/ListCorrection/TrueConstant breakdown this fixture
// needs to assert — extending the shared struct would require backfilling
// those fields for AggelosB/HuaYou2/NewsunVision too, well beyond this
// task's scope. Driving the HTTP flow a second time through runDraftCalcTest
// just to get the base fields would also mean standing up two separate
// in-memory servers for one vessel; instead this asserts the full superset
// in a single pass, reusing the same testServer helpers and totals.templ
// label lookups (extractByDivID + extractTotalsCellValue) already used by
// extractNetDisplacement/extractCargoOnBoard.
func polarStarFixture() vesselFixture {
	vessel := map[string]string{
		fields.FieldVesselName:      "POLAR STAR",
		fields.FieldIMO:             "9471666",
		fields.FieldLBP:             "183",
		fields.FieldBreadth:         "28.5",
		fields.FieldDepth:           "15.1",
		fields.FieldSummerDraft:     "10.4",
		fields.FieldSummerDWT:       "37390.4",
		fields.FieldSummerFreeboard: "4.7",
		fields.FieldSummerTPC:       "50",
		fields.FieldLightship:       "9522.3",
		fields.FieldTableDensity:    "1.025",
		fields.FieldMMCMethod:       "marine",
		fields.FieldCorrMethod:      "Full LBP",
		fields.FieldLCFDetection:    "auto",
	}
	job := map[string]string{
		fields.FieldJobNumber:     "345333_033",
		fields.FieldDSNumber:      "1",
		fields.FieldClient:        "Our principal",
		fields.FieldCargoDeclared: "21626",
		fields.FieldConstDeclared: "342",
	}

	draft0 := map[string]string{
		fields.FieldFwdPort: "3.39", fields.FieldFwdPortMethod: "direct",
		fields.FieldFwdStbd: "3.36", fields.FieldFwdStbdMethod: "direct",
		fields.FieldMidPort: "4.64", fields.FieldMidPortMethod: "direct",
		fields.FieldMidStbd: "4.54", fields.FieldMidStbdMethod: "direct",
		fields.FieldAftPort: "6.12", fields.FieldAftPortMethod: "direct",
		fields.FieldAftStbd: "6.12", fields.FieldAftStbdMethod: "direct",

		fields.FieldDFwd: "4.8", fields.FieldDFwdDir: "A",
		fields.FieldDMid: "0.5", fields.FieldDMidDir: "A",
		fields.FieldDAft: "1.2", fields.FieldDAftDir: "A",

		// LCF given from AP (~98.5m) — see the doc comment above. Direction
		// "F" here is DSGear's raw label on the from-AP value, not a
		// from-midship forward/aft flag; CalcHydrostatics' k3 branch
		// re-derives the true from-midship sign via LBP/2 - LCF.
		fields.FieldUDraft: "4.567", fields.FieldUDisp: "18956.7", fields.FieldUTPC: "45.2", fields.FieldULCF: "98.509", fields.FieldULCFDir: "F",
		fields.FieldLDraft: "4.617", fields.FieldLDisp: "19182.7", fields.FieldLTPC: "45.2", fields.FieldLLCF: "98.457", fields.FieldLLCFDir: "F",

		fields.FieldPMTCDraft: "5.1", fields.FieldPMTC: "525.7",
		fields.FieldNMTCDraft: "4.1", fields.FieldNMTC: "498.8",

		fields.FieldDockwaterDensity: "1.017",
		fields.FieldTPCListPort:      "45.212", fields.FieldTPCListStbd: "45.129",

		fields.FieldSeaType: "wave", fields.FieldSeaCondition: "under 0.1m",

		fields.FieldHFO: "49.16", fields.FieldMDO: "87.3", fields.FieldLubOil: "34.767",
	}
	draft0BW := []tankFixture{
		{Type: "BW-TOTAL", Volume: "8266.233", Density: "1"},
	}
	draft0FW := []tankFixture{
		{Type: "FW-TOTAL", Volume: "287.8"},
	}
	draft0Want := draftExpected{MMC: 4.612, LBM: 179.4, TrueTrim: 2.800, NetDisplacement: 9864.797, TotalDeductibles: 8725.26}

	draft1 := map[string]string{
		fields.FieldFwdPort: "7.54", fields.FieldFwdPortMethod: "direct",
		fields.FieldFwdStbd: "7.53", fields.FieldFwdStbdMethod: "direct",
		fields.FieldMidPort: "7.75", fields.FieldMidPortMethod: "direct",
		fields.FieldMidStbd: "7.68", fields.FieldMidStbdMethod: "direct",
		fields.FieldAftPort: "8.01", fields.FieldAftPortMethod: "direct",
		fields.FieldAftStbd: "7.91", fields.FieldAftStbdMethod: "direct",

		fields.FieldDFwd: "4.8", fields.FieldDFwdDir: "A",
		fields.FieldDMid: "0.5", fields.FieldDMidDir: "A",
		fields.FieldDAft: "9.4", fields.FieldDAftDir: "F",

		fields.FieldUDraft: "7.717", fields.FieldUDisp: "33642.2", fields.FieldUTPC: "48", fields.FieldULCF: "93.421", fields.FieldULCFDir: "F",
		fields.FieldLDraft: "7.767", fields.FieldLDisp: "33882.3", fields.FieldLTPC: "48", fields.FieldLLCF: "93.334", fields.FieldLLCFDir: "F",

		fields.FieldPMTCDraft: "8.2", fields.FieldPMTC: "622.7",
		fields.FieldNMTCDraft: "7.2", fields.FieldNMTC: "590.5",

		fields.FieldDockwaterDensity: "1.021",
		fields.FieldTPCListPort:      "48", fields.FieldTPCListStbd: "48",

		fields.FieldSeaType: "wave", fields.FieldSeaCondition: "under 0.1m",

		fields.FieldHFO: "44.49", fields.FieldMDO: "83.47", fields.FieldLubOil: "34.767",
	}
	draft1BW := []tankFixture{
		{Type: "BW-TOTAL", Volume: "1589.092", Density: "1"},
	}
	draft1FW := []tankFixture{
		{Type: "FW-TOTAL", Volume: "280.7"},
	}
	draft1Want := draftExpected{MMC: 7.724, LBM: 168.8, TrueTrim: 0.461, NetDisplacement: 31490.747, TotalDeductibles: 2032.519}

	return vesselFixture{
		name:         "PolarStar",
		vessel:       vessel,
		job:          job,
		draft0:       draft0,
		draft0BW:     draft0BW,
		draft0FW:     draft0FW,
		draft0Want:   draft0Want,
		draft1:       draft1,
		draft1BW:     draft1BW,
		draft1FW:     draft1FW,
		draft1Want:   draft1Want,
		cargoOnBoard: 21625.950,
	}
}

func TestPolarStarDraftCalc(t *testing.T) {
	vf := polarStarFixture()
	ts := newTestServer(t)

	id := ts.createSurvey(t)
	ts.putSurvey(t, id, mergeMaps(vf.vessel, vf.job))

	ts.startDraft(t, id, 0)
	ts.putDraft(t, id, 0, vf.draft0)
	for _, tf := range vf.draft0BW {
		ts.addBWTank(t, id, 0, tf)
	}
	for _, tf := range vf.draft0FW {
		ts.addFWTank(t, id, 0, tf)
	}

	ts.addDraft(t, id)
	ts.startDraft(t, id, 1)
	ts.putDraft(t, id, 1, vf.draft1)
	for _, tf := range vf.draft1BW {
		ts.addBWTank(t, id, 1, tf)
	}
	for _, tf := range vf.draft1FW {
		ts.addFWTank(t, id, 1, tf)
	}

	sseData := ts.captureSSEEvent(t, sse.EventDraftCalc, 5*time.Second, func() {
		ts.putDraft(t, id, 1, map[string]string{fields.FieldDockwaterDensity: vf.draft1[fields.FieldDockwaterDensity]})
	})

	assertNear(t, "PolarStar draft0 MMC", extractMMC(t, sseData, 0), vf.draft0Want.MMC)
	assertNear(t, "PolarStar draft0 LBM", extractLBM(t, sseData, 0), vf.draft0Want.LBM)
	assertNear(t, "PolarStar draft0 TrueTrim", extractTrueTrim(t, sseData, 0), vf.draft0Want.TrueTrim)
	assertNear(t, "PolarStar draft0 NetDisplacement", extractNetDisplacement(t, sseData, 0), vf.draft0Want.NetDisplacement)

	assertNear(t, "PolarStar draft1 MMC", extractMMC(t, sseData, 1), vf.draft1Want.MMC)
	assertNear(t, "PolarStar draft1 LBM", extractLBM(t, sseData, 1), vf.draft1Want.LBM)
	assertNear(t, "PolarStar draft1 TrueTrim", extractTrueTrim(t, sseData, 1), vf.draft1Want.TrueTrim)
	assertNear(t, "PolarStar draft1 NetDisplacement", extractNetDisplacement(t, sseData, 1), vf.draft1Want.NetDisplacement)

	assertNear(t, "PolarStar CargoOnBoard (== draft1 CargoFromPrev)", extractCargoOnBoard(t, sseData, 1), vf.cargoOnBoard)

	// Trim/list/constant breakdown — not covered by draftExpected, asserted
	// directly against totals.templ's labelled cells via the same
	// extractByDivID + extractTotalsCellValue helpers extractNetDisplacement
	// and extractCargoOnBoard already use.
	totals0 := extractByDivID(t, sseData, "draft-total-0")
	assertNear(t, "PolarStar draft0 FirstTrimCorrection", extractTotalsCellValue(t, totals0, "1st Trim Corr."), -481.481)
	assertNear(t, "PolarStar draft0 SecondTrimCorrection", extractTotalsCellValue(t, totals0, "2nd Trim Corr."), 57.622)
	assertNear(t, "PolarStar draft0 ListCorrection", extractTotalsCellValue(t, totals0, "List Corr."), 0.05)
	assertNear(t, "PolarStar draft0 TrueConstant", extractTotalsCellValue(t, totals0, "Const Calculated"), 342.497)

	totals1 := extractByDivID(t, sseData, "draft-total-1")
	assertNear(t, "PolarStar draft1 FirstTrimCorrection", extractTotalsCellValue(t, totals1, "1st Trim Corr."), -23.083)
	assertNear(t, "PolarStar draft1 SecondTrimCorrection", extractTotalsCellValue(t, totals1, "2nd Trim Corr."), 1.87)
	assertNear(t, "PolarStar draft1 ListCorrection", extractTotalsCellValue(t, totals1, "List Corr."), 0.0)
	assertNear(t, "PolarStar draft1 TrueConstant", extractTotalsCellValue(t, totals1, "Const Calculated"), 342.497)

	// TotalDeductibles is not in the SSE payload — read off the full Draft
	// Readings page, same as runDraftCalcTest.
	pageHTML := ts.getDraftPage(t, id)
	assertNear(t, "PolarStar draft0 TotalDeductibles", extractTotalDeductibles(t, pageHTML, 0), vf.draft0Want.TotalDeductibles)
	assertNear(t, "PolarStar draft1 TotalDeductibles", extractTotalDeductibles(t, pageHTML, 1), vf.draft1Want.TotalDeductibles)
}

// --- BEAM ------------------------------------------------------------------
//
// Data source: cmd/verify/testdata/beam.json, built and verified against
// DSGear in the same raw data audit. Covers two more gaps: Port and
// Starboard list-TPC are equal on both drafts (112.7/112.7, 121.8/121.8),
// which forces List Correction to exactly 0 via calcListCorrectionV1's
// |tpcPort-tpcStbd| term — no existing fixture (all three have distinct
// Port/Stbd TPC) exercises that path. Deductibles.Others is 17000, two
// orders of magnitude larger than the largest existing fixture value (170,
// NewsunVision).
//
// Uses "manual" LCF detection (unlike PolarStarFixture): BEAM's hydrostatic
// LCF values (~10.6m, ~0.5m) are already from-midship and well under the k3
// threshold, so auto vs manual detection is not the scenario this fixture
// targets — keeping it on "manual" isolates PolarStar as the sole test of
// the auto-detection branch.
func beamFixture() vesselFixture {
	vessel := map[string]string{
		fields.FieldVesselName:      "BEAM",
		fields.FieldIMO:             "9591741",
		fields.FieldLBP:             "283.5",
		fields.FieldBreadth:         "45.0",
		fields.FieldDepth:           "24.8",
		fields.FieldSummerDraft:     "18.322",
		fields.FieldSummerDWT:       "179100.3",
		fields.FieldSummerFreeboard: "6.478",
		fields.FieldSummerTPC:       "122.4",
		fields.FieldLightship:       "26328.0",
		fields.FieldTableDensity:    "1.025",
		fields.FieldMMCMethod:       "marine",
		fields.FieldCorrMethod:      "Full LBP",
		fields.FieldLCFDetection:    "manual",
	}
	job := map[string]string{
		fields.FieldJobNumber:     "444572_001",
		fields.FieldDSNumber:      "1",
		fields.FieldClient:        "Our principal",
		fields.FieldCargoDeclared: "146700",
		fields.FieldConstDeclared: "500",
	}

	draft0 := map[string]string{
		fields.FieldFwdPort: "9.35", fields.FieldFwdPortMethod: "direct",
		fields.FieldFwdStbd: "9.25", fields.FieldFwdStbdMethod: "direct",
		fields.FieldMidPort: "9.45", fields.FieldMidPortMethod: "direct",
		fields.FieldMidStbd: "9.08", fields.FieldMidStbdMethod: "direct",
		fields.FieldAftPort: "9.52", fields.FieldAftPortMethod: "direct",
		fields.FieldAftStbd: "9.47", fields.FieldAftStbdMethod: "direct",

		fields.FieldDFwd: "3.42", fields.FieldDFwdDir: "A",
		fields.FieldDMid: "0.83", fields.FieldDMidDir: "A",
		fields.FieldDAft: "15.4", fields.FieldDAftDir: "F",

		fields.FieldUDraft: "9.25", fields.FieldUDisp: "98258.2", fields.FieldUTPC: "112.7", fields.FieldULCF: "10.669", fields.FieldULCFDir: "F",
		fields.FieldLDraft: "9.3", fields.FieldLDisp: "98821.8", fields.FieldLTPC: "112.7", fields.FieldLLCF: "10.618", fields.FieldLLCFDir: "F",

		fields.FieldPMTCDraft: "9.798", fields.FieldPMTC: "2076.1",
		fields.FieldNMTCDraft: "8.798", fields.FieldNMTC: "2023.1",

		fields.FieldDockwaterDensity: "1.02",
		fields.FieldTPCListPort:      "112.7", fields.FieldTPCListStbd: "112.7",

		fields.FieldSeaType: "ice", fields.FieldSeaCondition: "under 0.05m around",

		fields.FieldHFO: "863.38", fields.FieldMDO: "204.93", fields.FieldLubOil: "62.1", fields.FieldOthers: "17000",
	}
	draft0BW := []tankFixture{
		{Type: "BW-TOTAL", Volume: "53099.551", Density: "1"},
	}
	draft0FW := []tankFixture{
		{Type: "FW-TOTAL", Volume: "179.93"},
	}
	draft0Want := draftExpected{MMC: 9.298, LBM: 264.68, TrueTrim: 0.209, NetDisplacement: 26820.019, TotalDeductibles: 71409.891}

	draft1 := map[string]string{
		fields.FieldFwdPort: "17.1", fields.FieldFwdPortMethod: "direct",
		fields.FieldFwdStbd: "17.1", fields.FieldFwdStbdMethod: "direct",
		fields.FieldMidPort: "17.3", fields.FieldMidPortMethod: "direct",
		fields.FieldMidStbd: "17.3", fields.FieldMidStbdMethod: "direct",
		fields.FieldAftPort: "17.22", fields.FieldAftPortMethod: "direct",
		fields.FieldAftStbd: "17.22", fields.FieldAftStbdMethod: "direct",

		fields.FieldDFwd: "3.42", fields.FieldDFwdDir: "A",
		fields.FieldDMid: "0.83", fields.FieldDMidDir: "A",
		fields.FieldDAft: "15.4", fields.FieldDAftDir: "F",

		fields.FieldUDraft: "17.25", fields.FieldUDisp: "192341.0", fields.FieldUTPC: "121.7", fields.FieldULCF: "0.518", fields.FieldULCFDir: "A",
		fields.FieldLDraft: "17.3", fields.FieldLDisp: "192949.8", fields.FieldLTPC: "121.8", fields.FieldLLCF: "0.551", fields.FieldLLCFDir: "A",

		fields.FieldPMTCDraft: "17.75", fields.FieldPMTC: "2572.3",
		fields.FieldNMTCDraft: "16.75", fields.FieldNMTC: "2533.7",

		fields.FieldDockwaterDensity: "1.023",
		fields.FieldTPCListPort:      "121.8", fields.FieldTPCListStbd: "121.8",

		fields.FieldSeaType: "ice", fields.FieldSeaCondition: "under 0.05m around",

		fields.FieldHFO: "856.55", fields.FieldMDO: "204.93", fields.FieldLubOil: "62.1", fields.FieldOthers: "17000",
	}
	draft1BW := []tankFixture{
		{Type: "BW-TOTAL", Volume: "330.234", Density: "1"},
	}
	draft1FW := []tankFixture{
		{Type: "FW-TOTAL", Volume: "177.37"},
	}
	draft1Want := draftExpected{MMC: 17.266, LBM: 264.68, TrueTrim: 0.129, NetDisplacement: 173531.989, TotalDeductibles: 18631.184}

	return vesselFixture{
		name:         "Beam",
		vessel:       vessel,
		job:          job,
		draft0:       draft0,
		draft0BW:     draft0BW,
		draft0FW:     draft0FW,
		draft0Want:   draft0Want,
		draft1:       draft1,
		draft1BW:     draft1BW,
		draft1FW:     draft1FW,
		draft1Want:   draft1Want,
		cargoOnBoard: 146711.970,
	}
}

func TestBeamDraftCalc(t *testing.T) {
	vf := beamFixture()
	ts := newTestServer(t)

	id := ts.createSurvey(t)
	ts.putSurvey(t, id, mergeMaps(vf.vessel, vf.job))

	ts.startDraft(t, id, 0)
	ts.putDraft(t, id, 0, vf.draft0)
	for _, tf := range vf.draft0BW {
		ts.addBWTank(t, id, 0, tf)
	}
	for _, tf := range vf.draft0FW {
		ts.addFWTank(t, id, 0, tf)
	}

	ts.addDraft(t, id)
	ts.startDraft(t, id, 1)
	ts.putDraft(t, id, 1, vf.draft1)
	for _, tf := range vf.draft1BW {
		ts.addBWTank(t, id, 1, tf)
	}
	for _, tf := range vf.draft1FW {
		ts.addFWTank(t, id, 1, tf)
	}

	sseData := ts.captureSSEEvent(t, sse.EventDraftCalc, 5*time.Second, func() {
		ts.putDraft(t, id, 1, map[string]string{fields.FieldDockwaterDensity: vf.draft1[fields.FieldDockwaterDensity]})
	})

	assertNear(t, "Beam draft0 MMC", extractMMC(t, sseData, 0), vf.draft0Want.MMC)
	assertNear(t, "Beam draft0 LBM", extractLBM(t, sseData, 0), vf.draft0Want.LBM)
	assertNear(t, "Beam draft0 TrueTrim", extractTrueTrim(t, sseData, 0), vf.draft0Want.TrueTrim)
	assertNear(t, "Beam draft0 NetDisplacement", extractNetDisplacement(t, sseData, 0), vf.draft0Want.NetDisplacement)

	assertNear(t, "Beam draft1 MMC", extractMMC(t, sseData, 1), vf.draft1Want.MMC)
	assertNear(t, "Beam draft1 LBM", extractLBM(t, sseData, 1), vf.draft1Want.LBM)
	assertNear(t, "Beam draft1 TrueTrim", extractTrueTrim(t, sseData, 1), vf.draft1Want.TrueTrim)
	assertNear(t, "Beam draft1 NetDisplacement", extractNetDisplacement(t, sseData, 1), vf.draft1Want.NetDisplacement)

	assertNear(t, "Beam CargoOnBoard (== draft1 CargoFromPrev)", extractCargoOnBoard(t, sseData, 1), vf.cargoOnBoard)

	totals0 := extractByDivID(t, sseData, "draft-total-0")
	assertNear(t, "Beam draft0 FirstTrimCorrection", extractTotalsCellValue(t, totals0, "1st Trim Corr."), -88.235)
	assertNear(t, "Beam draft0 SecondTrimCorrection", extractTotalsCellValue(t, totals0, "2nd Trim Corr."), 0.408)
	assertNear(t, "Beam draft0 ListCorrection", extractTotalsCellValue(t, totals0, "List Corr."), 0.0)
	assertNear(t, "Beam draft0 TrueConstant", extractTotalsCellValue(t, totals0, "Const Calculated"), 492.019)

	totals1 := extractByDivID(t, sseData, "draft-total-1")
	assertNear(t, "Beam draft1 FirstTrimCorrection", extractTotalsCellValue(t, totals1, "1st Trim Corr."), 2.93)
	assertNear(t, "Beam draft1 SecondTrimCorrection", extractTotalsCellValue(t, totals1, "2nd Trim Corr."), 0.113)
	assertNear(t, "Beam draft1 ListCorrection", extractTotalsCellValue(t, totals1, "List Corr."), 0.0)
	assertNear(t, "Beam draft1 TrueConstant", extractTotalsCellValue(t, totals1, "Const Calculated"), 492.019)

	pageHTML := ts.getDraftPage(t, id)
	assertNear(t, "Beam draft0 TotalDeductibles", extractTotalDeductibles(t, pageHTML, 0), vf.draft0Want.TotalDeductibles)
	assertNear(t, "Beam draft1 TotalDeductibles", extractTotalDeductibles(t, pageHTML, 1), vf.draft1Want.TotalDeductibles)
}
