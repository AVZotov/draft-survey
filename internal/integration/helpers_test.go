//go:build integration

// Package integration contains end-to-end HTTP tests against the real
// service/storage stack (real SQLite, real handlers, real chi router) behind
// httptest.NewServer. Kept in a separate package/build tag from the unit
// tests — see internal/CLAUDE.md and docs/TESTING_CHEATSHEET.md.
package integration

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	draftsurvey "github.com/AVZotov/draft-survey"
	"github.com/AVZotov/draft-survey/internal/handler"
	"github.com/AVZotov/draft-survey/internal/handler/fields"
	"github.com/AVZotov/draft-survey/internal/service"
	"github.com/AVZotov/draft-survey/internal/sse"
	"github.com/AVZotov/draft-survey/internal/storage"
	"github.com/AVZotov/draft-survey/internal/validation"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// nopLogger discards all log output. Each vessel fixture below sends on the
// order of a hundred HTTP requests; the real SlogLogger would otherwise
// flood `go test -v` with DEBUG lines for every one of them.
type nopLogger struct{}

func (nopLogger) Info(string, ...any)          {}
func (nopLogger) Warn(string, ...any)          {}
func (nopLogger) Error(string, error, ...any)  {}
func (nopLogger) Audit(string, string, ...any) {}
func (nopLogger) Debug(string, ...any)         {}

type testServer struct {
	server *httptest.Server
}

// newTestServer wires the same real stack cmd/server/main.go wires — real
// SQLite (in-memory), real broker, real service layer, real Handler, real
// chi router — behind an httptest.Server.
//
// Deviation from the literal DSN in the task ("file::memory:?cache=shared"):
// that DSN is not unique, and SQLite's shared-cache in-memory mode is keyed
// by the URI string *per process* — every test in this package would open
// the exact same database. Since `go test` runs every test function in one
// process, that would let tests bleed state into each other (most subtly:
// the seeded surveyor profile, since SQLiteUserStore always keys off a
// single fixed row id=1). Each call here gets a fresh, uniquely-named
// in-memory database instead.
func newTestServer(t *testing.T) *testServer {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", uuid.New().String())
	db, err := storage.NewDB(dsn, draftsurvey.Dictionaries)
	if err != nil {
		t.Fatalf("open in-memory db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	log := nopLogger{}
	validator := validation.New()

	userStore := storage.NewSQLiteUserStore(db)
	dictionaryStore := storage.NewDictionariesStore(db)
	surveyStore := storage.NewSQLiteSurveyStore(db)

	services := &service.Services{
		User:       service.NewUserService(userStore, log, validator),
		Survey:     service.NewSurveyService(surveyStore, userStore, log),
		Draft:      service.NewDraftService(surveyStore, userStore, log, validation.NewDraftValidator()),
		Tank:       service.NewTankService(surveyStore, log),
		Dictionary: service.NewDictionaryService(dictionaryStore),
	}

	broker := sse.NewBroker()
	h := handler.New(services, log, "test", broker)

	r := chi.NewRouter()
	if err := handler.SetupRoutesChi(r, h); err != nil {
		t.Fatalf("setup routes: %v", err)
	}

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	ts := &testServer{server: srv}
	// SurveyService.Create dereferences the current user unconditionally
	// (mirrors production: home.go redirects to /profile first if none
	// exists) — createSurvey() would panic without this.
	ts.createProfile(t)
	return ts
}

func (ts *testServer) url(path string) string {
	return ts.server.URL + path
}

func toValues(m map[string]string) url.Values {
	v := url.Values{}
	for k, val := range m {
		v.Set(k, val)
	}
	return v
}

func (ts *testServer) postForm(t *testing.T, path string, form url.Values) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, ts.url(path), strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("build POST %s: %v", path, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST %s: %v", path, err)
	}
	return resp
}

func (ts *testServer) putForm(t *testing.T, path string, form url.Values) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPut, ts.url(path), strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("build PUT %s: %v", path, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT %s: %v", path, err)
	}
	return resp
}

func requireOK(t *testing.T, resp *http.Response, op string) {
	t.Helper()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("%s: unexpected status %d: %s", op, resp.StatusCode, body)
	}
}

// createProfile seeds the single surveyor profile row every other operation
// depends on. Not in the task's literal helper list — added because
// createSurvey() cannot succeed without it (see newTestServer's doc comment).
func (ts *testServer) createProfile(t *testing.T) {
	t.Helper()
	requireOK(t, ts.postForm(t, "/api/v1/profile", url.Values{
		"first_name": {"Test"},
		"last_name":  {"Surveyor"},
	}), "createProfile")
}

// createSurvey follows the real newSurvey flow: GET /survey/new redirects
// (303) to /survey/{id}. The new survey's id is read off the Location header.
func (ts *testServer) createSurvey(t *testing.T) string {
	t.Helper()

	client := &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(ts.url("/survey/new"))
	if err != nil {
		t.Fatalf("createSurvey: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("createSurvey: expected 303, got %d", resp.StatusCode)
	}

	loc := resp.Header.Get("Location")
	id := strings.TrimPrefix(loc, "/survey/")
	if id == "" || id == loc {
		t.Fatalf("createSurvey: unexpected redirect location %q", loc)
	}
	return id
}

// putSurvey saves survey-level fields. Named to match the task's helper
// list, but flagging a factual correction: the registered route is
// POST /api/v1/survey/{id} (h.saveSurvey in routes.go), not PUT.
func (ts *testServer) putSurvey(t *testing.T, id string, fieldsMap map[string]string) {
	t.Helper()
	requireOK(t, ts.postForm(t, "/api/v1/survey/"+id, toValues(fieldsMap)), "putSurvey")
}

// putDraft saves one draft's fields. fieldsMap keys are the *unsuffixed*
// fields.FieldXxx constants — putDraft appends "-{index}" itself, matching
// how internal/handler/parse/draft.go expects every field to be named.
func (ts *testServer) putDraft(t *testing.T, id string, index int, fieldsMap map[string]string) {
	t.Helper()
	form := url.Values{}
	for k, v := range fieldsMap {
		form.Set(fmt.Sprintf("%s-%d", k, index), v)
	}
	requireOK(t, ts.putForm(t, fmt.Sprintf("/api/v1/survey/%s/draft/%d", id, index), form), "putDraft")
}

func (ts *testServer) startDraft(t *testing.T, id string, index int) {
	t.Helper()
	requireOK(t, ts.postForm(t, fmt.Sprintf("/api/v1/survey/%s/draft/%d/start", id, index), url.Values{}), "startDraft")
}

// applyDensity triggers the bulk "Apply Density" action (fills BW tank
// Density only where currently nil). Not in the task's literal helper list —
// used as a clean, single-publish way to force one more EventTankCalc
// broadcast after a batch of tank additions settle, since it doesn't itself
// do a two-step add-then-update HTTP round trip the way addBWTank/addFWTank
// do (see their doc comment on why capturing SSE around those is unsafe for
// getting a "final settled" total).
func (ts *testServer) applyDensity(t *testing.T, id string, draftIndex int, density string) {
	t.Helper()
	requireOK(t, ts.putForm(t, fmt.Sprintf("/api/v1/survey/%s/tanks/%d/density", id, draftIndex), url.Values{
		fields.FieldDockwaterDensity: {density},
	}), "applyDensity")
}

// addDraft appends a second draft to the survey. Not in the task's literal
// helper list — added because CargoFromPrev/CargoOnBoard (assertions Step 2
// asks for) are only meaningful with two drafts, and all three cheatsheet
// vessels are two-draft (initial + final) surveys.
func (ts *testServer) addDraft(t *testing.T, id string) {
	t.Helper()
	requireOK(t, ts.postForm(t, "/api/v1/survey/"+id+"/draft", url.Values{}), "addDraft")
}

func (ts *testServer) getDraftPage(t *testing.T, id string) string {
	t.Helper()
	resp, err := http.Get(ts.url("/survey/" + id + "/draft"))
	if err != nil {
		t.Fatalf("getDraftPage: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("getDraftPage: unexpected status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("getDraftPage: read body: %v", err)
	}
	return string(body)
}

// tankFixture is one row from the cheatsheet's tank tables. Empty string
// fields are omitted from the request (never entered), matching the
// *float64 "nil = not entered" semantics documented in internal/CLAUDE.md.
type tankFixture struct {
	Type, Name       string
	Sounding, Volume string
	Density          string // BW only — FW tanks never submit this field
}

var tankRowIDRe = regexp.MustCompile(`id="tank-row-([a-f0-9-]+)"`)

// addTank mirrors the real two-step UI flow: POST creates a bare row (type +
// name only — see AddRowForm/addBWTank/addFWTank), then a PUT fills in the
// measurement fields (see TankItem's row-level autosave block). Returns the
// server-generated tank ID.
//
// Caveat for callers using captureSSEEvent(sse.EventTankCalc, ...): this
// method makes two requests, and BOTH publish their own EventTankCalc (the
// POST's publish reflects the tank at zero volume/density, before the PUT
// applies its measurements). Wrapping a call to addTank/addBWTank/addFWTank
// directly in captureSSEEvent's trigger will catch whichever of the two
// publishes happens to be read first — not reliably the final, fully-measured
// one. Use applyDensity (or another single-publish action) as the trigger
// instead once every tank you care about has already been added.
func (ts *testServer) addTank(t *testing.T, surveyID string, draftIndex int, bw bool, tf tankFixture) string {
	t.Helper()

	kind := "fw"
	if bw {
		kind = "bw"
	}

	addResp := ts.postForm(t, fmt.Sprintf("/api/v1/survey/%s/tanks/%d/%s", surveyID, draftIndex, kind), url.Values{
		fields.FieldTankType: {tf.Type},
		fields.FieldTankName: {tf.Name},
	})
	defer addResp.Body.Close()
	body, err := io.ReadAll(addResp.Body)
	if err != nil {
		t.Fatalf("addTank: read body: %v", err)
	}
	if addResp.StatusCode != http.StatusOK {
		t.Fatalf("addTank: unexpected status %d: %s", addResp.StatusCode, body)
	}

	m := tankRowIDRe.FindSubmatch(body)
	if m == nil {
		t.Fatalf("addTank: could not find tank id in response: %s", body)
	}
	tankID := string(m[1])

	measure := url.Values{}
	if tf.Sounding != "" {
		measure.Set(fields.WithTankID(fields.FieldTankSounding, tankID), tf.Sounding)
	}
	if tf.Volume != "" {
		measure.Set(fields.WithTankID(fields.FieldTankVolume, tankID), tf.Volume)
	}
	if bw && tf.Density != "" {
		measure.Set(fields.WithTankID(fields.FieldTankDensity, tankID), tf.Density)
	}
	if len(measure) > 0 {
		requireOK(t, ts.putForm(t, fmt.Sprintf("/api/v1/survey/%s/tanks/%d/%s/%s", surveyID, draftIndex, kind, tankID), measure), "addTank(measure)")
	}

	return tankID
}

func (ts *testServer) addBWTank(t *testing.T, surveyID string, draftIndex int, tf tankFixture) string {
	return ts.addTank(t, surveyID, draftIndex, true, tf)
}

func (ts *testServer) addFWTank(t *testing.T, surveyID string, draftIndex int, tf tankFixture) string {
	return ts.addTank(t, surveyID, draftIndex, false, tf)
}

// captureSSEEvent opens a fresh /events connection, waits for the response
// headers (confirming the server has already called broker.Subscribe() —
// see the flush-on-subscribe fix in internal/handler/events.go), THEN calls
// trigger. Without that ordering, a publish fired a moment before the
// listener subscribes is silently dropped: the broker is non-blocking by
// design (internal/CLAUDE.md), so there is no built-in redelivery to save a
// naive "fire trigger, then start listening" test from missing the event.
func (ts *testServer) captureSSEEvent(t *testing.T, eventType sse.EventType, timeout time.Duration, trigger func()) string {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout+2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.url("/events"), nil)
	if err != nil {
		t.Fatalf("captureSSEEvent: build request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("captureSSEEvent: connect: %v", err)
	}
	defer resp.Body.Close()

	type result struct {
		data string
		err  error
	}
	resultCh := make(chan result, 1)

	go func() {
		scanner := bufio.NewScanner(resp.Body)
		scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
		var eventName string
		for scanner.Scan() {
			line := scanner.Text()
			switch {
			case strings.HasPrefix(line, "event: "):
				eventName = strings.TrimPrefix(line, "event: ")
			case strings.HasPrefix(line, "data: "):
				if eventName == string(eventType) {
					resultCh <- result{data: strings.TrimPrefix(line, "data: ")}
					return
				}
			}
		}
		resultCh <- result{err: scanner.Err()}
	}()

	trigger()

	select {
	case res := <-resultCh:
		if res.err != nil {
			t.Fatalf("captureSSEEvent: stream error: %v", res.err)
		}
		return res.data
	case <-time.After(timeout):
		t.Fatalf("captureSSEEvent: timed out waiting for event %q", eventType)
		return ""
	}
}

// --- HTML fragment parsing -------------------------------------------------
//
// Step 2 asks for parsing "by the same ids that calc_panel.templ and
// totals.templ render" — these helpers do exactly that, scoped per draft
// index so a single EventDraftCalc blob (which always carries fragments for
// every draft, see publishDraftCalc in draft.go) can't cross-match the wrong
// draft's numbers.

const tolerance = 0.001

func assertNear(t *testing.T, label string, got, want float64) {
	t.Helper()
	diff := got - want
	if diff < 0 {
		diff = -diff
	}
	if diff > tolerance {
		t.Errorf("%s: got %v, want %v (diff %v > tolerance %v)", label, got, want, diff, tolerance)
	}
}

func mustParseFloat(t *testing.T, raw string) float64 {
	t.Helper()
	clean := strings.ReplaceAll(strings.TrimSpace(raw), " ", "")
	v, err := strconv.ParseFloat(clean, 64)
	if err != nil {
		t.Fatalf("parse float %q: %v", raw, err)
	}
	return v
}

// extractByDivID returns the substring spanning a top-level <div id="...">
// through its matching closing </div>, walking nested <div> depth. Good
// enough for these generated fragments (calc-panel, draft-total-N), which
// nest only other divs/spans below the target id, never re-open a sibling
// div with the same id.
func extractByDivID(t *testing.T, html, id string) string {
	t.Helper()
	idx := strings.Index(html, `id="`+id+`"`)
	if idx == -1 {
		t.Fatalf("extractByDivID: id %q not found", id)
	}
	start := strings.LastIndex(html[:idx], "<div")
	if start == -1 {
		t.Fatalf("extractByDivID: no opening <div for id %q", id)
	}

	depth := 0
	i := start
	for i < len(html) {
		switch {
		case strings.HasPrefix(html[i:], "<div"):
			depth++
			i += len("<div")
		case strings.HasPrefix(html[i:], "</div>"):
			depth--
			i += len("</div>")
			if depth == 0 {
				return html[start:i]
			}
		default:
			i++
		}
	}
	t.Fatalf("extractByDivID: unbalanced <div> for id %q", id)
	return ""
}

// extractCPVal reads a calc-panel cell's value by its label
// (web/widgets/drafts/calc_panel.templ: <span class="cp-lbl">label</span>
// <span class="cp-val">value [direction-letter]</span>). The regex only
// captures the leading numeric run, dropping any trailing direction suffix.
func extractCPVal(t *testing.T, calcPanelHTML, label string) float64 {
	t.Helper()
	re := regexp.MustCompile(regexp.QuoteMeta(label) + `</span>\s*<span class="cp-val">\s*(-?[0-9. ]+)`)
	m := re.FindStringSubmatch(calcPanelHTML)
	if m == nil {
		t.Fatalf("extractCPVal: label %q not found in:\n%s", label, calcPanelHTML)
	}
	return mustParseFloat(t, m[1])
}

func extractMMC(t *testing.T, sseData string, draftIndex int) float64 {
	t.Helper()
	panelHTML := extractByDivID(t, sseData, fmt.Sprintf("calc-panel-%d", draftIndex))
	return extractCPVal(t, panelHTML, "MMC")
}

func extractTrueTrim(t *testing.T, sseData string, draftIndex int) float64 {
	t.Helper()
	panelHTML := extractByDivID(t, sseData, fmt.Sprintf("calc-panel-%d", draftIndex))
	return extractCPVal(t, panelHTML, "True Trim m")
}

// extractLBM reads hydrostatics.templ's `LBM {value}m` span by its
// draft-scoped id (constants.LBM + "-{index}").
func extractLBM(t *testing.T, sseData string, draftIndex int) float64 {
	t.Helper()
	re := regexp.MustCompile(fmt.Sprintf(`id="lbm-%d"[^>]*>\s*LBM\s*(-?[0-9.]+)\s*m`, draftIndex))
	m := re.FindStringSubmatch(sseData)
	if m == nil {
		t.Fatalf("extractLBM: lbm-%d not found in:\n%s", draftIndex, sseData)
	}
	return mustParseFloat(t, m[1])
}

// extractTotalsCellValue reads a Totals.templ secondary-grid cell
// (`.totals-cell-label` / `.totals-cell-value` pair) by its label.
func extractTotalsCellValue(t *testing.T, totalsHTML, label string) float64 {
	t.Helper()
	re := regexp.MustCompile(`totals-cell-label">` + regexp.QuoteMeta(label) + `</span>\s*<span class="totals-cell-value">\s*([0-9. -]+)`)
	m := re.FindStringSubmatch(totalsHTML)
	if m == nil {
		t.Fatalf("extractTotalsCellValue: label %q not found in:\n%s", label, totalsHTML)
	}
	return mustParseFloat(t, m[1])
}

func extractNetDisplacement(t *testing.T, sseData string, draftIndex int) float64 {
	t.Helper()
	totalsHTML := extractByDivID(t, sseData, fmt.Sprintf("draft-total-%d", draftIndex))
	return extractTotalsCellValue(t, totalsHTML, "Net Displ.")
}

// extractCargoOnBoard reads Totals.templ's `.totals-cargo-value` cell.
//
// Deviation flagged: the task asked to assert "CargoFromPrev (draft 1 vs
// draft 0)". Totals.templ only renders CargoFromPrev under
// `.totals-cargo-value` when the draft's Type is Intermediate — for a
// two-draft survey the second (last) draft is always Final, where the same
// slot instead renders sr.CargoOnBoard (survey-level). All three cheatsheet
// vessels are two-draft (initial + final) surveys, and CargoOnBoard's
// expected value is numerically the same quantity ("cargo added between the
// two drafts") — confirmed by the cheatsheet's own section 5 numbers
// matching draft 1's CargoFromPrev row exactly. So this asserts CargoOnBoard
// instead, since CargoFromPrev is never present in the rendered DOM for a
// Final-type draft.
func extractCargoOnBoard(t *testing.T, sseData string, draftIndex int) float64 {
	t.Helper()
	totalsHTML := extractByDivID(t, sseData, fmt.Sprintf("draft-total-%d", draftIndex))
	re := regexp.MustCompile(`totals-cargo-label">Cargo On Board</span>\s*<span class="totals-cargo-value">\s*([0-9. -]+)`)
	m := re.FindStringSubmatch(totalsHTML)
	if m == nil {
		t.Fatalf("extractCargoOnBoard: not found in:\n%s", totalsHTML)
	}
	return mustParseFloat(t, m[1])
}

// extractTotalDeductibles reads DeductiblesGrid's readonly "Total Deduct, MT"
// input. Deviation flagged: this field is rendered only on the full Draft
// Readings page (web/widgets/drafts/deductibles.templ), not in the
// EventDraftCalc SSE payload — publishDraftCalc in draft.go only re-renders
// CalcPanel/MMCRow/LBM/DeltaMtc/Totals, not DeductiblesGrid. So this parses
// a full GET /survey/{id}/draft page (via getDraftPage), not the SSE blob.
func extractTotalDeductibles(t *testing.T, draftPageHTML string, draftIndex int) float64 {
	t.Helper()
	gridHTML := extractByDivID(t, draftPageHTML, fmt.Sprintf("deduct-%d", draftIndex))
	re := regexp.MustCompile(`<label>Total Deduct, MT</label>\s*<input readonly value="([^"]*)"`)
	m := re.FindStringSubmatch(gridHTML)
	if m == nil {
		t.Fatalf("extractTotalDeductibles: not found in:\n%s", gridHTML)
	}
	return mustParseFloat(t, m[1])
}
