//go:build integration

package integration

import (
	"regexp"
	"testing"
	"time"

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
