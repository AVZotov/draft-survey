package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/AVZotov/draft-survey/internal/calculation"
	"github.com/AVZotov/draft-survey/internal/format"
	"github.com/AVZotov/draft-survey/internal/handler/routes"
	"github.com/AVZotov/draft-survey/internal/service"
	"github.com/AVZotov/draft-survey/internal/sse"
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/templates/pages"
	"github.com/AVZotov/draft-survey/web/widgets/tanks"
	"github.com/AVZotov/draft-survey/web/widgets/tanks/corrections"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) getTanks(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.getTanks"

	id := chi.URLParam(r, "id")
	draftIndex, err := strconv.Atoi(chi.URLParam(r, "draftIndex"))
	if err != nil {
		h.logger.Error(op, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := h.services.Survey.GetPageData(id)
	if err != nil {
		h.logger.Error(op, err)
		h.respondPageError(w, r, routes.SurveyList())
		return
	}

	if data.Outcome != nil && data.Outcome.Redirect != "" {
		return
	}

	// Trim/list are informational only — nil is acceptable if the draft has no marks yet.
	var trim, list *float64
	var trimDir, listDir string
	if draftIndex < len(data.Survey.Drafts) {
		trim, list, trimDir, listDir = tankTrimList(data.Survey.Drafts[draftIndex], data.Survey.VesselData)
	}

	lp := web.TanksLayoutProps(data.User, h.appVersion)
	tp := web.TanksPageProps(*data.Survey, draftIndex, trim, list, trimDir, listDir)

	if err = pages.Tanks(lp, tp).Render(r.Context(), w); err != nil {
		h.logger.Error(op, err)
	}
}

// tankTrimList computes the current true trim/list for a draft, used both by
// the Tanks page itself and by the calibration modal (which displays them as
// read-only context for the surveyor).
func tankTrimList(draft types.Draft, vessel types.VesselData) (trim, list *float64, trimDir, listDir string) {
	dr := calculation.CalcDraft(draft, vessel)
	if dr.TrueTrim != 0 {
		trim = &dr.TrueTrim
		trimDir = format.TrimDirection(dr.TrueTrim)
	}
	if dr.ListDegrees != 0 {
		list = &dr.ListDegrees
		listDir = format.ListDirection(dr.ListDegrees)
	}
	return trim, list, trimDir, listDir
}

func (h *Handler) getTankCorrections(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.getTankCorrections"

	id := chi.URLParam(r, "id")
	tankID := chi.URLParam(r, "tankID")
	draftIndex, err := strconv.Atoi(chi.URLParam(r, "draftIndex"))
	if err != nil {
		h.logger.Error(op, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := h.services.Survey.GetPageData(id)
	if err != nil {
		h.logger.Error(op, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if draftIndex < 0 || draftIndex >= len(data.Survey.Drafts) {
		h.logger.Error(op, fmt.Errorf("draft index %d out of range", draftIndex))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idx := slices.IndexFunc(data.Survey.Drafts[draftIndex].BallastWaterTanks, func(t types.Tank) bool {
		return t.ID == tankID
	})
	if idx == -1 {
		h.logger.Error(op, fmt.Errorf("tank %s not found", tankID))
		w.WriteHeader(http.StatusNotFound)
		return
	}
	tank := data.Survey.Drafts[draftIndex].BallastWaterTanks[idx]

	trim, list, trimDir, listDir := tankTrimList(data.Survey.Drafts[draftIndex], data.Survey.VesselData)
	tp := web.TanksPageProps(*data.Survey, draftIndex, trim, list, trimDir, listDir)

	templ.Handler(corrections.ModalForm(tank, tp)).ServeHTTP(w, r)
}

func (h *Handler) addBWTank(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.addBWTank"

	id := chi.URLParam(r, "id")
	draftIndex, err := strconv.Atoi(chi.URLParam(r, "draftIndex"))
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	if err = r.ParseForm(); err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	survey, err := h.services.Survey.Get(id)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	if draftIndex < 0 || draftIndex >= len(survey.Drafts) {
		err := fmt.Errorf("draft index %d out of range", draftIndex)
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	tank := types.Tank{ID: uuid.New().String()}
	h.decoder.TankMeta(r, &tank)

	sr, _, err := h.services.Tank.AddBW(survey, draftIndex, tank)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}
	bwTanks := survey.Drafts[draftIndex].BallastWaterTanks
	added := bwTanks[len(bwTanks)-1]

	h.publishTankCalc(r.Context(), *survey, *sr, draftIndex)

	if err = tanks.TankItem(id, draftIndex, added, true, false, len(bwTanks) == 1, true).Render(r.Context(), w); err != nil {
		h.logger.Error(op, err)
	}
}

func (h *Handler) addFWTank(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.addFWTank"

	id := chi.URLParam(r, "id")
	draftIndex, err := strconv.Atoi(chi.URLParam(r, "draftIndex"))
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	if err = r.ParseForm(); err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	survey, err := h.services.Survey.Get(id)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	if draftIndex < 0 || draftIndex >= len(survey.Drafts) {
		err := fmt.Errorf("draft index %d out of range", draftIndex)
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	tank := types.Tank{ID: uuid.New().String()}
	h.decoder.TankMeta(r, &tank)

	sr, _, err := h.services.Tank.AddFW(survey, draftIndex, tank)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}
	fwTanks := survey.Drafts[draftIndex].FreshWaterTanks
	added := fwTanks[len(fwTanks)-1]

	h.publishTankCalc(r.Context(), *survey, *sr, draftIndex)

	if err = tanks.TankItem(id, draftIndex, added, false, false, len(fwTanks) == 1, true).Render(r.Context(), w); err != nil {
		h.logger.Error(op, err)
	}
}

func (h *Handler) updateBWTank(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.updateBWTank"

	id := chi.URLParam(r, "id")
	tankID := chi.URLParam(r, "tankID")
	draftIndex, err := strconv.Atoi(chi.URLParam(r, "draftIndex"))
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	if err = r.ParseForm(); err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	survey, err := h.services.Survey.Get(id)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	if draftIndex < 0 || draftIndex >= len(survey.Drafts) {
		err := fmt.Errorf("draft index %d out of range", draftIndex)
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	idx := slices.IndexFunc(survey.Drafts[draftIndex].BallastWaterTanks, func(t types.Tank) bool {
		return t.ID == tankID
	})
	if idx == -1 {
		err := fmt.Errorf("tank %s not found", tankID)
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}
	tank := survey.Drafts[draftIndex].BallastWaterTanks[idx]

	h.decoder.Tank(r, &tank)
	// The corrections modal also posts Type/Name (the row itself never does) —
	// parse them here too so an edit made from the modal isn't silently dropped.
	h.decoder.TankMeta(r, &tank)
	h.decoder.TankCorrections(r, &tank)

	sr, _, err := h.services.Tank.UpdateBW(survey, draftIndex, tank)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}
	bwTanks := survey.Drafts[draftIndex].BallastWaterTanks
	updated := bwTanks[idx]

	h.publishTankCalc(r.Context(), *survey, *sr, draftIndex)

	if err = tanks.TankItem(id, draftIndex, updated, true, false, idx == 0, idx == len(bwTanks)-1).Render(r.Context(), w); err != nil {
		h.logger.Error(op, err)
	}
}

func (h *Handler) updateFWTank(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.updateFWTank"

	id := chi.URLParam(r, "id")
	tankID := chi.URLParam(r, "tankID")
	draftIndex, err := strconv.Atoi(chi.URLParam(r, "draftIndex"))
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	if err = r.ParseForm(); err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	survey, err := h.services.Survey.Get(id)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	if draftIndex < 0 || draftIndex >= len(survey.Drafts) {
		err := fmt.Errorf("draft index %d out of range", draftIndex)
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	idx := slices.IndexFunc(survey.Drafts[draftIndex].FreshWaterTanks, func(t types.Tank) bool {
		return t.ID == tankID
	})
	if idx == -1 {
		err := fmt.Errorf("tank %s not found", tankID)
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}
	tank := survey.Drafts[draftIndex].FreshWaterTanks[idx]

	h.decoder.Tank(r, &tank)

	sr, _, err := h.services.Tank.UpdateFW(survey, draftIndex, tank)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}
	fwTanks := survey.Drafts[draftIndex].FreshWaterTanks
	updated := fwTanks[idx]

	h.publishTankCalc(r.Context(), *survey, *sr, draftIndex)

	if err = tanks.TankItem(id, draftIndex, updated, false, false, idx == 0, idx == len(fwTanks)-1).Render(r.Context(), w); err != nil {
		h.logger.Error(op, err)
	}
}

func (h *Handler) deleteBWTank(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.deleteBWTank"

	id := chi.URLParam(r, "id")
	tankID := chi.URLParam(r, "tankID")
	draftIndex, err := strconv.Atoi(chi.URLParam(r, "draftIndex"))
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	survey, err := h.services.Survey.Get(id)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	if draftIndex < 0 || draftIndex >= len(survey.Drafts) {
		err := fmt.Errorf("draft index %d out of range", draftIndex)
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	sr, outcome, err := h.services.Tank.DeleteBW(survey, draftIndex, tankID)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	h.publishTankCalc(r.Context(), *survey, *sr, draftIndex)

	h.respond(w, outcome)
}

func (h *Handler) deleteFWTank(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.deleteFWTank"

	id := chi.URLParam(r, "id")
	tankID := chi.URLParam(r, "tankID")
	draftIndex, err := strconv.Atoi(chi.URLParam(r, "draftIndex"))
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	survey, err := h.services.Survey.Get(id)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	if draftIndex < 0 || draftIndex >= len(survey.Drafts) {
		err := fmt.Errorf("draft index %d out of range", draftIndex)
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	sr, outcome, err := h.services.Tank.DeleteFW(survey, draftIndex, tankID)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	h.publishTankCalc(r.Context(), *survey, *sr, draftIndex)

	h.respond(w, outcome)
}

func (h *Handler) applyDensity(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.applyDensity"

	id := chi.URLParam(r, "id")
	draftIndex, err := strconv.Atoi(chi.URLParam(r, "draftIndex"))
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	if err = r.ParseForm(); err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	survey, err := h.services.Survey.Get(id)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	if draftIndex < 0 || draftIndex >= len(survey.Drafts) {
		err := fmt.Errorf("draft index %d out of range", draftIndex)
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	density, ok := h.decoder.DockwaterDensity(r)
	if !ok || density == nil {
		err := errors.New("dockwater density is required")
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	sr, outcome, err := h.services.Tank.ApplyDensity(survey, draftIndex, *density)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	h.publishTankCalc(r.Context(), *survey, *sr, draftIndex)

	h.respond(w, outcome)
}

// publishTankCalc renders the BW/FW block-header totals, the action-bar
// totals, and the full BW/FW tank rows (order + move-button disabled state
// can change on any mutation) for the given draft, and publishes them as one
// EventTankCalc SSE event. Mirrors publishDraftCalc in draft.go — bypasses
// Outcome (reserved for Toast/Alert) since these are pre-rendered read-only
// display fragments, not notifications.
func (h *Handler) publishTankCalc(ctx context.Context, survey types.Survey, sr types.SurveyResult, draftIndex int) {
	const op = "Handler.publishTankCalc"

	if draftIndex < 0 || draftIndex >= len(sr.DraftTotals) {
		return
	}
	dr := sr.DraftTotals[draftIndex].DraftResult
	draftType := survey.Drafts[draftIndex].Type
	bwWeight := format.WeightFormatted(dr.TotalBwTanksWeight)
	fwWeight := format.WeightFormatted(dr.TotalFwTanksWeight)

	var buf bytes.Buffer
	if err := tanks.TableFormHeader(tanks.BWTableHeaderID(draftIndex), "Ballast Water Tanks", bwWeight, draftType, true).Render(ctx, &buf); err != nil {
		h.logger.Error(op, err)
		return
	}
	if err := tanks.TableFormHeader(tanks.FWTableHeaderID(draftIndex), "Fresh Water Tanks", fwWeight, draftType, true).Render(ctx, &buf); err != nil {
		h.logger.Error(op, err)
		return
	}
	if err := tanks.ABTotals(bwWeight, fwWeight, draftIndex, true).Render(ctx, &buf); err != nil {
		h.logger.Error(op, err)
		return
	}
	if err := tanks.TankRows(survey.ID, draftIndex, survey.Drafts[draftIndex].BallastWaterTanks, true, true).Render(ctx, &buf); err != nil {
		h.logger.Error(op, err)
		return
	}
	if err := tanks.TankRows(survey.ID, draftIndex, survey.Drafts[draftIndex].FreshWaterTanks, false, true).Render(ctx, &buf); err != nil {
		h.logger.Error(op, err)
		return
	}

	h.broker.Publish(sse.Event{Type: sse.EventTankCalc, Data: buf.String()})
}

func (h *Handler) moveBWTankUp(w http.ResponseWriter, r *http.Request) {
	h.moveTank(w, r, "Handler.moveBWTankUp", h.services.Tank.MoveBWTankUp)
}

func (h *Handler) moveBWTankDown(w http.ResponseWriter, r *http.Request) {
	h.moveTank(w, r, "Handler.moveBWTankDown", h.services.Tank.MoveBWTankDown)
}

func (h *Handler) moveFWTankUp(w http.ResponseWriter, r *http.Request) {
	h.moveTank(w, r, "Handler.moveFWTankUp", h.services.Tank.MoveFWTankUp)
}

func (h *Handler) moveFWTankDown(w http.ResponseWriter, r *http.Request) {
	h.moveTank(w, r, "Handler.moveFWTankDown", h.services.Tank.MoveFWTankDown)
}

// moveTank is the shared body for the four move handlers above — they only
// differ in which TankService method swaps the tank.
func (h *Handler) moveTank(
	w http.ResponseWriter, r *http.Request, op string,
	move func(survey *types.Survey, draftIndex int, tankID string) (*types.SurveyResult, *service.Outcome, error),
) {
	id := chi.URLParam(r, "id")
	tankID := chi.URLParam(r, "tankID")
	draftIndex, err := strconv.Atoi(chi.URLParam(r, "draftIndex"))
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	survey, err := h.services.Survey.Get(id)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	if draftIndex < 0 || draftIndex >= len(survey.Drafts) {
		err := fmt.Errorf("draft index %d out of range", draftIndex)
		h.logger.Error(op, err)
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	sr, outcome, err := move(survey, draftIndex, tankID)
	if err != nil {
		h.logger.Error(op, err)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	h.publishTankCalc(r.Context(), *survey, *sr, draftIndex)

	h.respond(w, outcome)
}
