# Draft Survey Tool — Development Roadmap

**Project:** https://github.com/AVZotov/draft-survey  
**Status:** Phase 1.6 In Progress 🚧  
**Last Updated:** 2026-05-14

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
- **Backup Strategy:** Each survey = separate JSON file (max 10-20/month)
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
- [x] Rounding policy aligned with UNECE Form C reporting fields
- [x] Golden tests against UNECE 1992 standard (BEAM IMO 9591741)
- [x] Integration tests — POLAR_STAR (TrimNoList, TrimList scenarios)
- [x] Unit tests — DSGear reference data

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
- [x] Local JSON storage
- [x] Profile persistence between sessions

---

### 1.4 Storage Layer ✅ COMPLETE
**Location:** `internal/storage/`

- [x] Survey CRUD operations
- [x] JSON files (one file per survey)
- [x] UUID-based survey IDs
- [x] Survey list and search
- [x] Auto-save on every change

---

### 1.5 Report Generation 🔲 PLANNED
**Location:** `internal/report/` *(not yet created)*

- [ ] PDF library selection (gofpdf / unidoc)
- [ ] UNECE-compliant report template
- [ ] Vessel data, calculations
- [ ] Save to custom location

**Target:** v0.4.0

---

### 1.6 Web Interface (HTMX) 🚧 IN PROGRESS
**Location:** `web/`, `internal/handler/`

- [x] Fiber server
- [x] Templ templates
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
- [ ] Alerts block (trim by head, list > 0.5°, constant deviation, deflection)
- [x] Final results page (full survey summary)
- [x] Form validation improvements

---

### 1.7 Logging 🔲 PLANNED
- [ ] Structured logging (slog)
- [ ] Log to file
- [ ] Log rotation

**Target:** v0.4.0

---

### 1.8 Error Handling 🚧 IN PROGRESS
- [x] Custom error types for handlers
- [ ] Error messages dictionary
- [ ] User-friendly error display improvements

---

### 1.9 Testing ✅ SUBSTANTIALLY COMPLETE
- [x] Golden tests against UNECE 1992 (BEAM initial draft — 27 tests)
- [x] Integration tests POLAR_STAR (TrimNoList + TrimList)
- [x] Unit tests — full calculation chain
- [x] Tank volume calibration tests (Type 1/2/3)
- [x] Storage layer tests
- [ ] E2E test: create survey → calculate → generate PDF

**Coverage:** >80% for `internal/calculation/`

---

### 1.10 Dictionaries ✅ COMPLETE
- [x] Ports list (JSON, embedded)
- [x] Country flags (JSON, embedded)
- [x] Sea conditions (wave / ice)
- [x] Tank types dictionary
- [x] All assets embedded in binary at build time

---

### 1.11 Distribution 🔲 PLANNED
- [ ] Cross-compilation (Windows / macOS / Linux)
- [ ] First-run setup
- [ ] Pre-built binaries

**Target:** v0.4.0

---

## Known Issues / Open Questions

- [ ] List correction method when TPC upper == TPC lower — DSGear uses Summer TPC interpolation (~0.881 MT), UNECE strict gives 0. Awaiting surveyor clarification. See Issue #TBD
- [ ] TPC List Port / Stbd fields not yet in draft form UI — using V2 auto-calculation as interim
- [ ] MTC draft hint fields — pre-filled vs readonly, awaiting surveyor feedback

---

## PHASE 2: Commercial Version (Cloud-Sync)

**Goal:** Multi-user system with central repository

### 2.1 Authentication
- [ ] User login/password
- [ ] Long-lived tokens (2 months offline support)

### 2.2 Backend Server
**Tech Stack:** Go + Fiber + PostgreSQL
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
- [ ] PDF report generation
- [ ] Historical data analysis
- [ ] Export to Excel
- [ ] Partial cargo calculation by hold groups (see Issue #TBD)
- [ ] IsExtrapolated flag for tank volumes outside calibration range

---

## Versioning

| Version | Status | Description |
|---------|--------|-------------|
| v0.3.0-rc | 🚧 Current | Core math complete, web UI in progress |
| v0.4.0 | 🔲 Planned | PDF reports, logging, pre-built binaries |
| v1.0.0 | 🔲 Planned | Full MVP with distribution |

---

## Development Strategy

- `main` — stable releases only
- `feature/*` — feature development
- `fix/*` — bug fixes
- `docs/*` — documentation updates
- PR to `main` when ready
- Semantic versioning (v0.1.0, v0.2.0, v1.0.0)
- GitHub Issues for all tasks and open questions

---

*This is a living document — updated with each significant change.*