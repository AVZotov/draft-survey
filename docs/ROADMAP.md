# Draft Survey Tool — Development Roadmap

**Project:** https://github.com/AVZotov/draft-survey  
**Status:** v0.5.0 — Chi migration complete; core math and web UI verified against UNECE 1992 and cross-checked against independent commercial software (7 integration tests, 5 vessels) 🚧  
**Last Updated:** 2026-08-21

---

## PHASE 0: Foundation ✅ COMPLETE

### Infrastructure
- [x] Repository created (lowercase naming)
- [x] Project structure with Go standard layout
- [x] Allowlist-style `.gitignore` for security
- [x] GitHub Actions CI (ubuntu-22.04, Go 1.25)
- [x] SSH authentication configured

### Decisions Made
- **Workflow:** GitHub Flow (main + feature branches)
- **Go Version:** 1.25 (minimum compatibility)
- **CI Runner:** ubuntu-22.04 (LTS, reproducible builds)
- **PDF Reports:** Generate on-demand only (no auto-save)
- **Storage:** SQLite via modernc.org/sqlite (pure Go, no CGO)
- **Temp Files:** Auto-deleted after final report generated

---

## PHASE 1: MVP — Open Source (Offline-First)

**Goal:** Working draft survey calculator with local storage

---

### 1.1 Core Mathematics Module ✅ COMPLETE
**Location:** `internal/calculation/`

- [x] UNECE 1992 formulas implemented:
  - Quarter Mean Draft (6M/8, 4M/6, 3/10f variants)
  - PP Corrections — Full LBP and Half LBP methods
  - Keel corrections
  - First Trim Correction (LCF/FTC)
  - Second Trim Correction (Nemoto)
  - List Correction — V1 (manual TPC) and V2 (Summer TPC interpolation)
  - LCF/LCA interpolation — auto (k3 coefficient) and manual (UNECE standard) modes
  - Density Correction (dockwater vs table density)
  - Net Displacement with full correction chain
  - Cargo Weight, Ship's Constant, Current DWT, True Constant
- [x] Tank volume calculation via calibration tables:
  - Type 1: Volume by Trim (most common)
  - Type 2: Sounding Correction
  - Type 3: Volume Correction
  - Optional list correction for all types
- [x] `CalcDraft` — single entry point returning complete `DraftResult`
- [x] `CalcSurvey` — multi-draft orchestrator with cargo and constant tracking
- [x] Guard against division by zero (LBP, interpolation, list correction)
- [x] Rounding policy aligned with UNECE Form C reporting fields
- [x] Golden tests against UNECE 1992 standard
- [x] Integration tests — 7 tests across 5 vessels (AGGELOS B, HUA YOU 2, NEWSUN VISION against UNECE 1992 golden data; POLAR STAR, BEAM against DSGear commercial-software reference data), zero discrepancies found (±0.001 tolerance)
- [x] Standalone `cmd/verify/` CLI tool — runs CalcDraft/CalcSurvey against raw JSON survey data for cross-checking outside the HTTP stack

---

### 1.2 Vessel Data Module ✅ COMPLETE
**Location:** `internal/types/` (merged into shared types)

- [x] Vessel information structure (name, IMO, flag, built year)
- [x] Hydrostatic table rows with LCF direction support
- [x] PP distances and keel corrections per draft
- [x] Vessel types: marine / river / barge
- [x] Correction methods: Full LBP / Half LBP
- [x] LCF detection modes: auto (k3) / manual (UNECE)
- [x] Summer draft, DWT, TPC, freeboard
- [x] Table density field for density correction
- [x] JSON serialization

---

### 1.3 User Profile ✅ COMPLETE
**Location:** `internal/storage/`

- [x] Surveyor profile (name, position, company, license, country)
- [x] Stored in SQLite database
- [x] Profile persistence between sessions

---

### 1.4 Storage Layer ✅ COMPLETE
**Location:** `internal/storage/`

- [x] Survey CRUD operations
- [x] SQLite database (one file — `data/draft-survey.db`)
- [x] Hybrid schema: metadata columns + JSON blob
- [x] UUID-based survey IDs
- [x] Survey list and search (`SurveyQueryRepository`)
- [x] Auto-save on every change
- [x] Schema migration mechanism with versioning
- [x] Graceful shutdown — safe DB close on SIGTERM/SIGINT

---

### 1.5 Chi Migration (v0.5.0) ✅ COMPLETE
**Location:** `internal/handler/`, `internal/sse/`, `internal/service/`

Full architectural rewrite, not just a router swap — see `docs/V050_MIGRATION_CONTRACT_EN.md` for the complete contract.

- [x] Chi replaces Fiber entirely (Fiber fully removed from go.mod)
- [x] Single unified SSE broker for every server→client notification (Toast/Alert/live calc-panel updates) — no `HX-Trigger` headers, no parallel mechanisms
- [x] Handler/Service/Repository layering enforced project-wide — handler as "low-paid secretary," zero business logic outside the service layer
- [x] `internal/handler/fields/` — form field constants (replaces `internal/constants/`)
- [x] `internal/handler/routes/` — URL builder functions (no inline `fmt.Sprintf`)
- [x] `internal/handler/parse/` — hand-written Decoder for true partial-update autosave (not `go-playground/form`, not `gorilla/schema`)
- [x] `Outcome` DTO — every service method returns `(*Outcome, error)`, including Update/Delete
- [x] Draft Readings, Tanks, Results pages — all wired and working
- [x] `internal/validation/playground/`, `internal/types/dto/`, `internal/constants/` — pending deletion once verified zero references (see migration contract Wave 4)

---

### 1.6 Report Generation 🔲 PLANNED
**Location:** `internal/report/` *(scaffolded, not yet implemented)*

- [ ] PDF library selection (gofpdf / unidoc)
- [ ] UNECE Form C compliant report template
- [ ] Vessel data, calculations
- [ ] Save to custom location

**Target:** v0.7.0 — see [#146](https://github.com/AVZotov/draft-survey/issues/146)

---

### 1.7 Web Interface (HTMX) ✅ COMPLETE
**Location:** `web/`, `internal/handler/`

- [x] Chi router + templ templates
- [x] Survey list with search and date filter
- [x] New survey form
- [x] Draft readings input with real-time calculation panel
- [x] SVG vessel diagram — dynamic trim and list visualization
- [x] Hydrostatic data entry
- [x] MTC data entry with MMC-based draft hints
- [x] Deductibles block — BW/FW tanks, fuel, others
- [x] Tank calibration modal — Type 1/2/3 with trim/list tables
- [x] Totals & Results block
- [x] Sea condition logging
- [x] Surveyor profile page
- [x] Final results page (full survey summary)
- [x] Unified numeric input formatting (`XInput` component, `x-mask`)
- [ ] Alerts block (trim by head, list > 0.5°, constant deviation, deflection) — moved to v0.7.0 soft validation, see [#148](https://github.com/AVZotov/draft-survey/issues/148) and the DSGear-derived validation rules [#149](https://github.com/AVZotov/draft-survey/issues/149)–[#156](https://github.com/AVZotov/draft-survey/issues/156)

---

### 1.8 Logging ✅ COMPLETE
- [x] Structured logging (slog), configurable via `LOG_LEVEL` env var
- [x] slog-based HTTP request middleware
- [ ] Log to file
- [ ] Log rotation

---

### 1.9 Error Handling 🚧 IN PROGRESS
- [x] `internal/errors/` sentinel errors (`ErrEmptyField`, `ErrInvalidFormat`)
- [x] User-friendly error display via Toast/Alert (SSE-delivered, not a separate mechanism)
- [ ] Broader sentinel error coverage (`storage.ErrNotFound` and similar)
- [ ] Unified error handling strategy across all handlers

---

### 1.10 Testing ✅ SUBSTANTIALLY COMPLETE
- [x] Golden tests against UNECE 1992 standard
- [x] Integration tests — 7 tests across 5 vessels (AGGELOS B, HUA YOU 2, NEWSUN VISION, POLAR STAR, BEAM), zero discrepancies vs. UNECE 1992 and DSGear reference data (±0.001 tolerance)
- [x] Unit tests — full calculation chain
- [x] Tank volume calibration tests (Type 1/2/3)
- [x] Standalone `cmd/verify/` CLI tool for cross-checking raw survey data outside the HTTP stack
- [ ] Storage layer tests (SQLite)
- [ ] E2E test: create survey → calculate → generate PDF

**Coverage:** >80% for `internal/calculation/`

---

### 1.11 Dictionaries ✅ COMPLETE
- [x] Ports list (JSON, embedded)
- [x] Country flags (JSON, embedded)
- [x] Sea conditions (wave / ice)
- [x] Tank types dictionary
- [x] All assets embedded in binary at build time

---

### 1.12 Distribution 🔲 PLANNED
- [ ] Cross-compilation (Windows / macOS / Linux)
- [ ] First-run setup
- [ ] Pre-built binaries via GitHub Releases
- [ ] Third-party license acknowledgment (THIRD_PARTY_NOTICES)

**Target:** relates to v0.6.0's one-click desktop launch, see [#145](https://github.com/AVZotov/draft-survey/issues/145)

---

## Known Issues / Open Questions

- [x] ~~List correction method when TPC upper == TPC lower~~ — resolved: `calcListCorrectionV2` implements Summer TPC interpolation for this case, matching DSGear's approach.
- [x] ~~TPC List Port / Stbd fields not yet in draft form UI~~ — resolved: wired into `web/widgets/drafts/marks.templ`.
- [x] ~~`Lightship` field — currently `float64`, should be `*float64`~~ — resolved: `types.VesselData.Lightship` is `*float64`.
- [ ] MTC draft hint fields — pre-filled vs readonly, awaiting surveyor feedback
- [ ] User signature storage — separate table needed to avoid loading binary data on every profile fetch

---

## PHASE 1.6: v0.6.0 — Frontend Redesign 🔲 PLANNED

**Goal:** Theming, localization, and desktop-friendly distribution for the surveyor audience.

- [ ] Dark/light theme toggle — CSS variables already use design tokens, add the toggle ([#143](https://github.com/AVZotov/draft-survey/issues/143))
- [ ] Russian/English UI locale — currently English-only ([#144](https://github.com/AVZotov/draft-survey/issues/144))
- [ ] One-click desktop launch (systray/icon) — launch by clicking an icon, not a terminal command ([#145](https://github.com/AVZotov/draft-survey/issues/145))

---

## PHASE 1.7: v0.7.0 — Advanced Math & Surveyor UX 🔲 PLANNED

**Goal:** Close the gap with commercial draft-survey software — reporting, soft validation, and hold-level cargo math.

- [ ] PDF report generation — UNECE Form C compliant ([#146](https://github.com/AVZotov/draft-survey/issues/146))
- [ ] In-app user-facing documentation page — Footer Documentation link currently shows a toast stub ([#147](https://github.com/AVZotov/draft-survey/issues/147))
- [ ] Soft validation of surveyor input — `DraftValidator` stub already exists ([#148](https://github.com/AVZotov/draft-survey/issues/148)), informed by 8 validation rules extracted from a surveyor's DSGear XLSM files ([#149](https://github.com/AVZotov/draft-survey/issues/149)–[#156](https://github.com/AVZotov/draft-survey/issues/156)): list angle > 0.5° requiring LOP, trim-by-head, Current DWT exceeding Summer DWT, hull-twist cross-check, Depth/Draft/Freeboard consistency, missing draft readings, instrument calibration expiry, tank listing consistency
- [ ] Partial cargo calculation by hold groups (see Phase 3.2)
- [ ] AI surveyor agent — surveyor UX automation, exact scope TBD

---

## PHASE 2: Commercial Version (Cloud-Sync)

**Goal:** Multi-user system with central repository

### 2.1 Authentication
- [ ] User login/password
- [ ] JWT with configurable expiry (temporary access for assistants)
- [ ] Session management

### 2.2 Backend Server
**Tech Stack:** Go + Chi + PostgreSQL
- [ ] REST API
- [ ] Database schema
- [ ] UUID-based IDs

### 2.3 Offline-First Sync
- [ ] Conflict resolution
- [ ] Delta sync
- [ ] Sync status indicator

### 2.4 User Roles
- [ ] Surveyor, Coordinator, Admin

### 2.5 Deployment
- [ ] Docker Compose
- [ ] Self-hosted instructions

---

## PHASE 3: Polish & Scale

### 3.1 Internationalization
- [ ] Russian
- [ ] Chinese
- [ ] English (current)

### 3.2 Advanced Features
- [ ] Historical data analysis
- [ ] Export to Excel
- [ ] Data Import/Export (JSON, MDB/Access legacy format)
- [ ] Partial cargo calculation by hold groups
- [ ] IsExtrapolated flag for tank volumes outside calibration range
- [ ] System tray / single-click shutdown for desktop mode

---

## Versioning

| Version | Status | Description |
|---------|--------|-------------|
| v0.3.0 | ✅ Released | Core math complete, web UI, modal strategy |
| v0.4.0 | ✅ Released | SQLite storage migration, bug fixes |
| v0.5.0 | 🚧 Current | Chi migration (Fiber removed), unified SSE broker, structured logging, Draft Readings/Tanks/Results pages complete, 7 integration tests across 5 vessels |
| v0.6.0 | 🔲 Planned | Dark/light theme, RU/EN locale, one-click desktop launch |
| v0.7.0 | 🔲 Planned | PDF reports, soft validation, hold-group calculations, AI surveyor agent |
| v1.0.0 | 🔲 Planned | Full MVP with distribution |

---

## Development Strategy

- `main` — stable releases only
- `feature/*` — feature development
- `fix/*` — bug fixes
- `tech/*` — technical debt
- `docs/*` — documentation updates
- PR to `main` when ready
- Semantic versioning (v0.1.0, v0.2.0, v1.0.0)
- GitHub Issues for all tasks and open questions

---

*This is a living document — updated with each significant change.*
