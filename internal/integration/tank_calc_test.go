//go:build integration

package integration

import (
	"regexp"
	"testing"
	"time"

	"github.com/AVZotov/draft-survey/internal/handler/fields"
	"github.com/AVZotov/draft-survey/internal/sse"
	"github.com/AVZotov/draft-survey/web/widgets/tanks"
)

// TestAggelosBTankCalc adds every AGGELOS B draft-0 BW tank from
// docs/TESTING_CHEATSHEET.md (manual-volume path — none of them carry a
// calibration table) and asserts the resulting TotalBwTanksWeight.
func TestAggelosBTankCalc(t *testing.T) {
	ts := newTestServer(t)

	id := ts.createSurvey(t)
	ts.startDraft(t, id, 0)

	bwTanks := aggelosBFixture().draft0BW
	for _, tf := range bwTanks {
		ts.addBWTank(t, id, 0, tf)
	}

	// Every AGGELOS B BW tank already carries an explicit density (manual-
	// volume path — see the cheatsheet's dataset note), so this fills
	// nothing; it exists purely to force one more, single-publish
	// EventTankCalc broadcast to observe now that every tank has settled.
	// (addBWTank itself makes two requests that each publish their own
	// EventTankCalc — see its doc comment for why capturing SSE around it
	// directly isn't reliable here.)
	sseData := ts.captureSSEEvent(t, sse.EventTankCalc, 5*time.Second, func() {
		ts.applyDensity(t, id, 0, "1.025")
	})

	const wantTotalBW = 9237.7930
	assertNear(t, "AggelosB draft0 TotalBwTanksWeight", extractBWTotal(t, sseData, 0), wantTotalBW)
}

// TestBeamTankCalc covers a gap flagged in the raw-data verification
// report: BEAM's Deductibles.Others is 17000 — two orders of magnitude
// larger than the largest existing fixture's Others value (170,
// NewsunVision). Adds BEAM's draft-0 BW/FW tanks plus its declared
// HFO/MDO/LubOil/Others, then verifies CalcTotalDeductibles sums all of it
// correctly at that scale. Uses beamFixture() from draft_calc_test.go for
// the tank data, same as TestAggelosBTankCalc reuses aggelosBFixture().
func TestBeamTankCalc(t *testing.T) {
	ts := newTestServer(t)

	id := ts.createSurvey(t)
	ts.startDraft(t, id, 0)

	vf := beamFixture()
	ts.putDraft(t, id, 0, map[string]string{
		fields.FieldHFO:    vf.draft0[fields.FieldHFO],
		fields.FieldMDO:    vf.draft0[fields.FieldMDO],
		fields.FieldLubOil: vf.draft0[fields.FieldLubOil],
		fields.FieldOthers: vf.draft0[fields.FieldOthers],
	})
	for _, tf := range vf.draft0BW {
		ts.addBWTank(t, id, 0, tf)
	}
	for _, tf := range vf.draft0FW {
		ts.addFWTank(t, id, 0, tf)
	}

	// TotalDeductibles is not in the SSE payload (see extractTotalDeductibles's
	// doc comment on draft_calc_test.go's runDraftCalcTest) — read it off the
	// full Draft Readings page instead.
	pageHTML := ts.getDraftPage(t, id)
	const wantTotalDeductibles = 71409.891 // 53099.551 (BW) + 179.93 (FW) + 863.38 + 204.93 + 62.1 + 17000
	assertNear(t, "Beam draft0 TotalDeductibles", extractTotalDeductibles(t, pageHTML, 0), wantTotalDeductibles)
}

var totalChipValRe = regexp.MustCompile(`total-chip-val">\s*([0-9. -]+)</span>`)

// extractBWTotal reads the BW block-header total-weight chip
// (web/widgets/tanks/table_form_header.templ). Uses the production
// tanks.BWTableHeaderID helper directly rather than reconstructing the id
// format, so this can't silently drift out of sync with the real id.
func extractBWTotal(t *testing.T, sseData string, draftIndex int) float64 {
	t.Helper()
	headerHTML := extractByDivID(t, sseData, tanks.BWTableHeaderID(draftIndex))
	m := totalChipValRe.FindStringSubmatch(headerHTML)
	if m == nil {
		t.Fatalf("extractBWTotal: total-chip-val not found in:\n%s", headerHTML)
	}
	return mustParseFloat(t, m[1])
}
