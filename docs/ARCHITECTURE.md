# Draft Survey — Architecture & Design Decisions

> **This document records the architecture and design decisions behind Draft Survey.**
> Read this alongside `internal/CLAUDE.md` (backend patterns) and `web/CLAUDE.md` (frontend patterns).
> All three documents together define the full architecture. No deviations without explicit instruction.

---

## What changed in v0.5.0 and why

This was not a router swap. It was a full architectural rewrite of a working v0.4.0 application:

1. **Chi replaces Fiber** — Fiber is incompatible with `net/http` stdlib and has been fully removed.
2. **Single unified SSE approach** — v0.4.0 used a mix of `HX-Trigger` headers, custom events, and other mechanisms. v0.5.0 uses one delivery path for everything: SSE broker + EventType + JS listener. No exceptions per page.
   When `Outcome` contains both `Redirect` and `Toast`/`Alert`, the notification is written to a flash cookie (`CookieFlashToast` / `CookieFlashAlert`) instead of publishing immediately to the SSE broker. The `/events` endpoint reads these cookies on reconnect and publishes them after the redirect completes. This prevents the race where `HX-Redirect` tears down the SSE connection before the broker publish reaches a listener.
3. **Consistent, predictable patterns** — v0.4.0 became hard to maintain months after release. v0.5.0 enforces the same pattern on every page so any page can be understood without reading all the others.
4. **Handler as "low-paid secretary"** — parses and delegates only. Zero business logic in handlers.
5. **localhost is not a bottleneck** — all data is sent per block, not minimized. Optimization for network transfer is not a goal. Clarity and consistency are.
6. **Calculated results are never persisted** — always recomputed on the fly. No stale state to track.

---

## Migration status

| Area | Status |
|---|---|
| Chi router (Fiber fully removed) | Done |
| SSE broker (internal/sse/) | Done |
| Logger interface + SlogLogger | Done — configurable via LOG_LEVEL env var. slog-based HTTP middleware. |
| fields/fields.go — form field constants | Done |
| routes/routes.go — URL builder functions | Done |
| parse/parse.go — Decoder + private helpers | Done |
| parse/survey.go — reference parser implementation | Done |
| respond.go — respond() + respondError() | Done |
| Profile page | Done |
| Dashboard / Home | Done |
| New Survey page | Done |
| Survey List page | Done |
| Draft Readings page | Done |
| Tanks page | Not started |
| Results page | Not started |
| DraftService implementation | Done |
| parse/draft.go | Done |
| parse/tanks.go | Not started |
| Delete internal/constants/ | Pending — after all pages migrated |
| Delete internal/types/dto/ | Pending — after all pages migrated |
| Delete internal/validation/playground/ | Pending — after all pages migrated |
| Squash migrations 1-7 into one | Pending — v0.5.0 release task |

---

## Reference implementations (read these first)

Before implementing anything new, read these files as the ground truth for patterns:

| File | What it demonstrates |
|---|---|
| internal/handler/survey.go | Page handler + API handler + select handlers |
| internal/handler/parse/survey.go | Full parser with block methods and field constants |
| internal/handler/respond.go | respond() and respondError() |
| internal/handler/routes/routes.go | URL builder function naming and format |
| internal/handler/routes.go | SetupRoutesChi — how routes are registered |
| web/data.go | LayoutProps + PageProps constructors |
| web/widgets/survey/ | Templ component structure for a complete page |
| web/widgets/survey_list/ | Pagination, SSE stats, search, delete patterns |

---

## Current blocker — Draft Readings

**Current state:** getDraft in internal/handler/draft.go is a stub (logs id, writes 200, returns). pages.DraftReadings is never called. CalcSurvey is called nowhere outside its own package.

**The panic:** block.templ:13 does sr.DraftTotals[draftIndex].DraftResult unconditionally. Passing types.SurveyResult{} (zero value, nil DraftTotals) causes index-out-of-range panic because at least one Draft always exists in a Survey.

**The fix** — see internal/CLAUDE.md section "Draft Readings — current blocker and fix" for the exact implementation.

**Why CalcSurvey is called directly in the handler (not through DraftService):** The initial page load is read-only. DraftService.Update() is for autosave flows (field change + save + recalculate + SSE publish). A read-only page render does not save anything, so calling CalcSurvey directly is correct and intentional.

---

## Implementation order — remaining pages

### Wave 1: Draft Readings (current)

Step-by-step — do not skip steps or merge them:

1. Implement getDraft handler — page renders without panic
2. Implement parse/draft.go — Decoder.Draft(r, existing, index) with block methods
3. Add PUT routes for draft blocks in routes/routes.go + SetupRoutesChi
4. Implement updateDraft handler (autosave one draft block)
5. Implement DraftService.Update() — save + CalcSurvey + broker.Publish
6. Add EventDraftCalc SSE event type + JS listener for calc panel updates
7. Add port/country modal (native dialog element) for Load Port / Discharge Port entry
8. Test with real vessel data

### Wave 2: Tanks

- Separate page, separate route
- Draft Readings displays Total BW + Total FW as read-only aggregates computed from Tank page data
- Follows the same handler + parser + service + SSE pattern as Draft Readings
- Do not start until Wave 1 is complete and tested

### Wave 3: Results

- Read-only page: load Survey + CalcSurvey + render
- No autosave, no parser required
- Do not start until Wave 2 is complete

### Wave 4: Cleanup (separate task, after all pages)

Run this before cleanup to find remaining references:

  grep -rn "internal/constants|types/dto|validation/playground" . --include="*.go"

Delete in this order:
1. internal/validation/playground/ — no v0.5.0 consumers
2. internal/types/dto/ — replaced by direct types + service.Outcome
3. internal/constants/ — replaced by internal/handler/fields/
4. Squash SQLite migrations 1-7 into one v0.5.0 migration

Do not delete any of these prematurely — verify zero references first.

---

## Domain model — key facts for implementation

### Survey to Draft relationship
- Survey.Drafts []Draft — index is slice position, no stored index field in Draft struct
- Draft index exists only in URL/context (draftIndex param), never persisted
- On delete: slice is recompacted (append(drafts[:i], drafts[i+1:]...)), no gaps ever
- First draft (index 0) cannot be deleted — always DraftTypeInitial
- Draft initialization invariant: every new Draft (whether created by Survey.Create() or
  DraftService.Add()) must be initialized with:
  - HydrostaticRows: make([]types.HydrostaticRow, 2)
  - MTCRows: make([]types.MTCRow, 2)
  - Type and Status set explicitly (never rely on zero values)
  Leaving these as nil causes index-out-of-range panics in hydrostatics.templ which unconditionally
  accesses rows[0] and rows[1].
- Survey.Audit must always be initialized as []AuditEvent{} (not nil) in
  SurveyService.Create() — nil slice serializes as null in JSON, breaking
  deserialization of existing records.

### Draft types and lifecycle
- DraftType: initial / intermediate / final
- DraftStatus: pending / active / complete
- Type is set explicitly on creation — not inferred from slice position

### Draft Readings page — input scope

| Block | Fields | Notes |
|---|---|---|
| Marks | 6x (Value *float64 + Method ReadingMethod) | 12 inputs total |
| PP Corrections | DistancePPFwd/Mid/Aft *float64 + 3 direction strings | |
| Keel | KeelFwd/Mid/Aft *float64 | values in mm |
| Hydrostatics | 2 rows x (Draft, Displacement, TPC, LCF, LCFDirection) | always exactly 2 rows |
| MTC | 2 rows x (Draft, MTC) | always exactly 2 rows, independent of Hydrostatics |
| TPC List | TPCListPort/Starboard *float64 | optional |
| Density | *float64 | |
| Deductibles | HFO/MDO/LubOil/BilgeWater/SewageWater/Others *float64 | direct weights, no sub-calculation |
| Sea Condition | Type toggle (wave/ice) + enum select | |
| Tanks | Total BW + Total FW | read-only on this page, computed on Tanks page |

### Calculation facts
- calculation.CalcSurvey(survey) SurveyResult — single entry point for full recalculation
- SurveyResult.DraftTotals length always equals len(survey.Drafts) — guaranteed by CalcSurvey
- Results are never stored in the database
- The calculation/ package is complete and fully tested — do not reimplement math in handlers or services

---

## Absolute rules — NEVER do these

- NEVER use internal/constants/ — use internal/handler/fields/fields.go
- NEVER use fmt.Sprintf for URL construction in templ or handlers — use routes.X()
- NEVER use HX-Trigger response headers for notifications — use SSE broker via Outcome
- NEVER use gorilla/schema or go-playground/form for parsing — use hand-written Decoder
- NEVER add form tags — autosave blocks only, no Save buttons
- NEVER persist calculated results — always recompute on the fly
- NEVER put business logic in handlers — service layer only
- NEVER add a new notification mechanism — Toast and Alert via SSE are the only two
- NEVER use r.Route() groups for route registration — direct registration only (caused double-prefix bug)
- NEVER start a new page implementation before the previous wave is complete and tested
- NEVER expect respond() to write a response body — it only sets headers and status. Autosave blocks must use hx-swap='none'. Pre-rendered SSE payloads are published directly via broker.Publish(), not through respond().
