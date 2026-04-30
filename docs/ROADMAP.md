# Draft Survey Tool — Development Roadmap

**Project:** https://github.com/AVZotov/draft-survey
**Current Version:** v0.1.0-rc
**Last Updated:** 2026-04-30

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
- **Backup Strategy:** Each survey = separate JSON file
- **Temp Files:** Auto-deleted after final report generated

---

## PHASE 1: MVP — Open Source (Offline-First)

**Goal:** Working draft survey calculator with local storage
**Status:** 🚧 Near Complete — PDF report remaining

---

### 1.1 Core Mathematics Module ✅ COMPLETE
**Location:** `internal/calculation/`

- [x] UNECE 1992 formulas: Quarter Mean Draft, FTC, STC, Density Correction, Displacement
- [x] PP & Keel corrections (Full LBP and Half LBP methods)
- [x] List correction
- [x] LCF/LCA interpolation (auto and manual direction detection)
- [x] Second trim correction via MTC rows
- [x] Tank volume calculation — 3 calibration table types:
  - Type 1: Volume by Trim (bilinear interpolation)
  - Type 2: Sounding Correction by Trim
  - Type 3: Volume Correction by Trim
- [x] List correction for all tank calibration types (optional)
- [x] `CalcSurvey` — orchestrates multi-draft survey calculation
- [x] Unit tests with hardcoded golden data (POLAR STAR, BEAM vessels)
- [ ] Input range validation logic

**Deliverable:** ✅ Calculation engine passing all test cases

---

### 1.2 Vessel Data Module ✅ COMPLETE
**Location:** `internal/vessel/`, `internal/types/`

- [x] VesselData structure (name, IMO, flag, dimensions, summer marks)
- [x] Vessel geometry (LBP, PP corrections, Keel thickness)
- [x] Shared domain types (Marks, Deductibles, Hydrostatics, SeaCondition, Survey)
- [x] SeaCondition types (Wave and Ice conditions with dictionaries)
- [x] MMC calculation methods: Standard (6/8), River (4/6), Barge (3+14+3/20)
- [x] Correction methods: Full LBP, Half LBP

**Decisions Made:**
- `internal/types/` — shared domain types used across all packages
- `internal/vessel/` — vessel passport data (VesselData)

**Deliverable:** ✅ Complete vessel data structures with domain types

---

### 1.3 User Profile ✅ COMPLETE
**Location:** `internal/types/`, `internal/storage/`

- [x] Surveyor profile (name, position, company, employee ID)
- [x] Local storage — no authentication
- [x] Profile persistence in `user.json`
- [x] `UserRepository` interface + `UserStore` implementation
- [x] Unit tests (SaveAndGet, GetWithNoUser, Delete)
- [x] `Survey.Surveyor` field — surveyor embedded in survey record

**Decisions Made:**
- Logout = `UserStore.Delete()` removes `user.json` from disk
- `Survey.Surveyor *User` — pointer allows nil (survey can exist without profile)

**Deliverable:** ✅ Surveyor can set up profile once, persists between restarts

---

### 1.4 Storage Layer ✅ COMPLETE
**Location:** `internal/storage/`

- [x] Survey CRUD (Create, Read, Update, Delete)
- [x] List all surveys with DTO projection
- [x] Search and filter surveys (by name, IMO, date range)
- [x] Pagination support (offset + limit)
- [x] `SurveyRepository` interface + `SurveyStore` implementation
- [x] `UserRepository` interface + `UserStore` implementation
- [x] `DictionariesRepository` interface + `DictionariesStore` with `embed.FS`
- [x] Dictionaries embedded in binary via `go:embed` (countries, ports)

**Decisions Made:**
- JSON files — simple, no dependencies, human-readable
- UUID-based filenames: `{uuid}.json`
- Dictionaries embedded in binary — no external files needed at runtime
- Future: MySQL as optional cloud-sync backend (Phase 2)

**Deliverable:** ✅ Surveys persist between app restarts, binary is self-contained

---

### 1.5 Report Generation 🔴 NEXT PRIORITY
**Location:** `internal/report/`

- [ ] PDF generation library selection (gofpdf / unidoc / chromedp)
- [ ] UNECE-compliant report template
- [ ] Include: vessel data, all draft readings, calculation steps, cargo result
- [ ] Surveyor signature block
- [ ] Save to user-chosen location

**Deliverable:** Generate professional PDF report suitable for official use

---

### 1.6 Web Interface ✅ COMPLETE
**Location:** `web/`

**Tech Stack Decisions:**
- Fiber v2 + HTMX + Templ (type-safe HTML templates)
- Alpine.js for lightweight client-side interactivity
- IBM Plex Sans/Mono fonts, custom CSS design tokens (no Tailwind)
- CSS split by page: `styles.css` (global) + page-specific files
- Repository pattern for all storage access
- Manual form parsing via `c.FormValue()` (granular, block-based saving)
- On-the-fly save with "save before navigate" pattern for data safety

**Pages Implemented:**
- [x] Surveyor profile setup (`/profile`)
- [x] Dashboard (`/`)
- [x] Survey list with search, date filter, stats (`/survey-list`)
- [x] New survey form — vessel data, cargo, job info (`/survey/new`)
- [x] Survey edit (`/survey/:id`)
- [x] Draft readings (`/survey/:id/draft`) with:
  - [x] Initial / Intermediate / Final draft blocks
  - [x] SVG vessel diagram (dynamic trim/list visualization)
  - [x] Sea condition toggle (Wave/Ice with condition dropdowns)
  - [x] PP & Keel distance inputs
  - [x] Real-time calculation panel (MMC, trim, list, deflection)
  - [x] Hydrostatic data input with interpolated MMC row
  - [x] MTC rows for second trim correction
  - [x] Deductibles section (HFO, MDO, LubOil, BilgeWater, Others)
  - [x] Start/Finish timestamps with action bar state machine
- [x] BW/FW Tanks page (`/survey/:id/tanks/:draftIndex`) with:
  - [x] Add / Remove tanks (BW and FW separately)
  - [x] Global density apply to tanks with nil density
  - [x] Tank calibration modal (3 table types, optional list correction)
  - [x] Real-time total weight calculation
  - [x] Trim/list correction warning banner
- [x] Calculation results page (`/survey/:id/results`) with:
  - [x] Draft selector (compare any two drafts)
  - [x] Step-by-step calculation display (collapsible sections)
  - [x] Cargo on board with discrepancy from S/P
  - [x] Survey status alerts (constant %, data warnings)

**Deliverable:** ✅ Working UI at `localhost:3399`

---

### 1.7 Logging 🔴 TODO
**Location:** `internal/logger/`

- [ ] Structured logging (`slog` — standard library, no dependencies)
- [ ] Log to file (`data/app.log`)
- [ ] Log levels (DEBUG/INFO/WARN/ERROR)
- [ ] Request logging middleware for Fiber

**Deliverable:** All operations logged for debugging and audit

---

### 1.8 Error Handling 🔴 TODO
**Location:** `internal/errors/`

- [ ] Custom error types for domain validation
- [ ] User-friendly error display in UI (banners, inline messages)
- [ ] HTTP error pages (404, 500)

**Deliverable:** Consistent, user-friendly error handling throughout app

---

### 1.9 Testing 🟡 PARTIAL
**Location:** `internal/calculation/`

- [x] Unit tests for calculation module (golden data — POLAR STAR, BEAM)
- [x] Tank volume calculation tests (all 3 calibration types)
- [x] Storage unit tests (SaveAndGet, Delete, GetWithNoUser)
- [ ] Integration tests for HTTP handlers
- [ ] E2E test: create survey → draft readings → tanks → calculate → PDF

**Target Coverage:** >80% for `internal/calculation/`

---

### 1.10 Dictionaries ✅ COMPLETE
**Location:** `data/dictionaries/`, `embed.go`

- [x] Countries list (JSON, embedded in binary)
- [x] Ports list (JSON, embedded in binary)
- [x] Wave conditions dictionary (Go)
- [x] Ice conditions dictionary (Go)
- [x] Cargo types dictionary (Go)
- [x] Packing types dictionary (Go)
- [x] Tank type names dictionary (Go)
- [x] Glossary EN + RU (Markdown, embedded in binary)

**Deliverable:** ✅ All dropdowns populated, binary is self-contained

---

### 1.11 Audit & Metadata 🟡 IN PROGRESS
**Location:** `internal/types/audit.go`, `internal/types/metadata.go`

- [x] `SurveyMetadata` struct — schema version, app version
- [x] `AuditEvent` struct — timestamp, event type, message, user
- [ ] Audit event recording (suspicious data, overrides, warnings)
- [ ] Telegram notification for critical audit events

**Deliverable:** Full audit trail embedded in each survey record

---

### 1.12 Distribution 🔴 TODO

- [ ] `go build` with embedded assets — single binary
- [ ] Makefile with cross-compilation (Windows / macOS / Linux)
- [ ] First-run setup (auto-create `data/` directories)
- [ ] GitHub Releases with pre-built binaries

---

## PHASE 2: Commercial Version (Cloud-Sync)

**Goal:** Multi-user system with central repository

### 2.1 Authentication
- [ ] User login/password
- [ ] JWT with configurable expiry (temporary access for assistants)
- [ ] Session management

### 2.2 Backend & Database
**Tech Stack:** Go + Fiber + MySQL

- [ ] MySQL as cloud storage backend (parallel to JSON for offline)
- [ ] Database schema (users, surveys, sync state)
- [ ] REST API for sync

### 2.3 Offline-First Sync
- [ ] Conflict resolution strategy
- [ ] Delta sync (only changed records)
- [ ] Sync status indicator in UI

### 2.4 User Roles
- [ ] Surveyor (create surveys)
- [ ] Assistant (temporary access via time-limited JWT)
- [ ] Coordinator (manage ports/dictionaries)
- [ ] Admin (user management)

### 2.5 Dictionary Management
- [ ] CRUD for ports and countries (coordinators only)
- [ ] Push updated dictionaries to clients on sync

### 2.6 Deployment
- [ ] Docker Compose (app + MySQL)
- [ ] Self-hosted instructions
- [ ] Security hardening

---

## PHASE 3: Polish & Scale

### 3.1 Internationalization (i18n)
- [ ] Russian interface
- [ ] English interface (current)

### 3.2 Performance
- [ ] Survey index file for fast list loading (avoid loading all JSON files)
- [ ] Caching strategy for dictionaries

### 3.3 Advanced Features
- [ ] Multiple surveyors per survey (primary + assistant)
- [ ] Digital signatures (e-signing PDFs)
- [ ] Export to Excel
- [ ] Historical data analysis
- [ ] gRPC between microservices (if cloud version scales to multiple services)

---

## Development Strategy

### Branching
- `main` — stable releases only
- `feature/*` — work on specific features
- PR to `main` when feature complete

### Versioning
- Semantic versioning: v0.1.0, v0.2.0, v1.0.0
- Git tags for releases
- `SurveyMetadata.SchemaVersion` — tracks breaking changes to survey structure

### CI/CD
- GitHub Actions on every push to `feature/*`
- Automated tests (`go test ./...`)
- Build artifacts on merge to `main`

---

## Immediate Next Steps (v0.1.0 → v0.2.0)

1. **PDF Report** (Phase 1.5) — core deliverable for real usage
2. **Audit event recording** (Phase 1.11) — log suspicious data automatically
3. **Logging middleware** (Phase 1.7) — structured request logging
4. **Distribution** (Phase 1.12) — Makefile + GitHub Release with binaries
5. **go:embed for static assets** — embed CSS/JS into binary (full self-contained build)

---

## Notes

- This is a living document — update with each significant change
- `SchemaVersion` in `SurveyMetadata` must be bumped on breaking struct changes
- Keep Phase 1 focus before expanding to Phase 2
- Open source version stays JSON-based; MySQL is Phase 2 cloud feature