# Draft Survey Tool

> A professional cargo weight calculator for marine surveyors, implementing the UNECE 1992 Draft Survey Code.

**Status:** 🚧 Active Development — v0.5.0 | Chi migration complete, core calculations verified against UNECE 1992 standard and cross-checked against independent commercial software

---

## Overview

Draft Survey Tool is an offline-first desktop application that helps marine surveyors calculate cargo weight using the draft survey method. All calculations strictly follow the **UNECE Code for the Uniform Application of the Rules for the Measurement of Bulk Cargoes (1992)**.

The mathematics is fully implemented, tested against real-world golden data, and verified to match the standard. The application runs locally — no internet connection required, no data leaves your machine.

---

## Screenshots

![Survey List](/.github/images/survey-list.png)

![Draft Readings](/.github/images/draft-readings.png)

![Calculation Results](/.github/images/results.png)

---

## Features

**Calculations (UNECE 1992 compliant)**
- Quarter Mean Draft (QMD) with PP and keel corrections
- Full LBP and Half LBP correction methods
- First and Second Trim Correction
- List Correction
- LCF/LCA interpolation — auto and manual direction detection
- Density Correction with FWA/DWA
- Deductibles: fuel, lubricants, ballast water, fresh water, others
- Tank volume calculation via calibration tables — 3 table types with optional list correction
- Multi-draft surveys: Initial → Intermediate(s) → Final
- Cargo on board with discrepancy from Shipper/Receiver declaration

**Application**
- Offline-first — all data stored locally in SQLite database
- Self-contained binary — dictionaries and assets embedded at build time
- Survey list with search by vessel name / IMO and date range filter
- Real-time calculation panel updated as you type, delivered over a single SSE broker (toast/alert notifications, live calc panel updates)
- SVG vessel diagram with dynamic trim and list visualization
- Sea condition logging (Wave / Ice with detailed condition selection)
- BW and FW tank management with calibration table modal
- Survey status alerts (constant deviation, data warnings)
- Surveyor profile with company details
- Auto-selects the next free port if the configured one is busy
- Graceful shutdown — safe database close on exit

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | [Go 1.25](https://go.dev/) + [Chi v5](https://github.com/go-chi/chi) |
| Templates | [Templ v0.3.1020](https://templ.guide/) |
| Interactivity | [HTMX 2.0.8](https://htmx.org/) + [Alpine.js 3.15.8](https://alpinejs.dev/) |
| DOM Morphing | [Idiomorph](https://github.com/bigskysoftware/idiomorph) |
| Fonts | [IBM Plex Sans & Mono](https://www.ibm.com/plex/) |
| Storage | [SQLite](https://www.sqlite.org/) via [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) (pure Go, no CGO) |
| Validation | [go-playground/validator](https://github.com/go-playground/validator) |

---

## Getting Started

### Prerequisites

- [Go 1.25+](https://go.dev/dl/)
- [Templ CLI](https://templ.guide/quick-start/installation) — for template generation
- [Air](https://github.com/air-verse/air) — for hot reload during development

### Run in Development Mode

```bash
# Clone the repository
git clone git@github.com:AVZotov/draft-survey.git
cd draft-survey

# Install dependencies
go mod download

# Run with hot reload
air

# Run with debug configuration
air -c .air.debug.toml
```

The application will be available at `http://localhost:3399`

### Data Directory

On first run, the application creates the following directory structure:

```
data/
  draft-survey.db    — SQLite database (surveys, user profile)
```

Dictionaries (ports, countries) are embedded in the binary and require no external files.

---

## Configuration

All configuration is via environment variables (see `.env.example`):

| Variable | Default | Description |
|---|---|---|
| `PORT` | `:3399` | Server listen address. If busy, the app automatically scans upward for the next free port and logs the one actually used. |
| `DB_PATH` | `./data/draft-survey.db` | SQLite database file location |
| `LOG_LEVEL` | `info` | `debug` \| `info` \| `warn` \| `error` |
| `VERSION` | `dev` | Version string shown in the footer |

---

## Pages

| Page | Route | Purpose |
|---|---|---|
| Dashboard | `/` | Landing page, entry point |
| Profile | `/profile` | Surveyor and company details |
| New Survey | `/survey/new` | Vessel, job, and cargo setup for a new survey |
| Survey List | `/survey-list` | Search and browse existing surveys |
| Draft Readings | `/survey/{id}/draft` | Marks, PP/keel corrections, hydrostatics, MTC, deductibles, live calc panel |
| Tanks | `/survey/{id}/tanks/{draftIndex}` | BW/FW tank soundings and calibration-table volume lookups |
| Results | `/survey/{id}/results` | Full survey summary — all drafts, cargo on board, constant |

---

## Architecture

Handler → Service → Repository, strictly one-directional (handlers never contain business logic, services never touch HTTP or templ, repositories never touch business rules — see `internal/CLAUDE.md`):

- **Handler** — parses the request, calls the service, calls `respond()`/`respondError()`. Nothing else.
- **Service** — business logic, validation, talks to the repository, returns `(*Outcome, error)` — including for `Update`/`Delete`, not just `Create`. A `nil` Outcome is a valid "silent success."
- **Repository** — CRUD only, raw SQL, no business rules.
- **SSE Broker** (`internal/sse/`) — the single delivery path for every server→client notification (toast, alert, live calc-panel updates). Non-blocking publish; a slow/dead client never blocks the publisher. No `HX-Trigger` headers, no parallel notification mechanism.
- **Decoder** (`internal/handler/parse/`) — a hand-written form parser (not `go-playground/form` or `gorilla/schema`) that updates only the fields actually present in the request, enabling true partial-update autosave without overwriting untouched fields.
- **Calculation engine** (`internal/calculation/`) — pure functions, read-only, never persisted; always recomputed on the fly from stored survey data.

---

## Project Structure

```
cmd/
  server/           — application entry point
  verify/           — standalone CLI tool: runs CalcDraft/CalcSurvey against a JSON file for verification
internal/
  calculation/      — UNECE 1992 math engine
  types/            — shared domain types
  storage/          — repository pattern (SQLite + embed.FS)
  handler/          — HTTP handlers (Chi)
    fields/         — form field name constants
    parse/          — hand-written form Decoder
    routes/         — URL builder functions
  service/          — business logic, Outcome-based responses
  sse/              — SSE broker and event types
  config/           — environment-variable configuration
  logger/           — structured logging (slog)
  format/           — value formatting helpers
  validation/       — go-playground/validator wrapper + DraftValidator
  integration/      — end-to-end HTTP tests (build tag: integration)
data/
  dictionaries/     — static data (ports, countries, conditions)
docs/               — documentation and glossary
web/
  templates/pages/  — page templates (Templ)
  layouts/          — page shell, HTMX/Alpine script includes
  widgets/          — reusable UI components
  components/       — shared UI elements
  ids/              — DOM id constants for OOB swaps and hx-include targeting
  static/           — CSS, JS, fonts
embed.go            — go:embed directives (self-contained build)
```

---

## Integration Tests

`internal/integration/` runs the full stack — real HTTP handlers behind `httptest.Server`, real (in-memory) SQLite, real service layer — driving each survey through the same requests the browser would send, then asserting the results against verified reference data.

```bash
go test -tags=integration ./internal/integration/...
```

**7 tests across 5 vessels**, each cross-checked against an independent source of truth:

| Vessel | Reference | Covers |
|---|---|---|
| AGGELOS B | UNECE 1992 golden data | Baseline draft calculation, tank calibration |
| HUA YOU 2 | UNECE 1992 golden data | Multi-tank ballast/fresh-water deductibles |
| NEWSUN VISION | UNECE 1992 golden data | Ice sea condition, LubOil/BilgeWater/Others deductibles |
| POLAR STAR | DSGear (commercial software) | LCF auto-detection / from-AP heuristic, nonzero list correction |
| BEAM | DSGear (commercial software) | Equal Port/Stbd TPC (list correction forced to 0), large-scale `Others` deductible |

Zero discrepancies (tolerance ±0.001) were found between this app's calculation engine and either reference across all five vessels — see `cmd/verify/` for the standalone verification tool used to cross-check raw survey data outside the HTTP stack.

---

## Roadmap

See [docs/ROADMAP.md](docs/ROADMAP.md) for the full development plan.

**Coming next:**
- **v0.6.0** — dark/light theme toggle, Russian/English UI locale, one-click desktop launch
- **v0.7.0** — PDF report generation, in-app documentation, soft validation of surveyor input, hold-group calculations

---

## Acknowledgements

This project is built with excellent open source tools. Grateful to their authors and contributors:

| Project | Version | License |
|---|---|---|
| [Go](https://go.dev/) | 1.25 | BSD-3-Clause |
| [Chi](https://github.com/go-chi/chi) | v5.3.0 | MIT |
| [Templ](https://templ.guide/) | v0.3.1020 | MIT |
| [HTMX](https://htmx.org/) | 2.0.8 | BSD Zero Clause |
| [Alpine.js](https://alpinejs.dev/) | 3.15.8 | MIT |
| [Idiomorph](https://github.com/bigskysoftware/idiomorph) | latest | BSD Zero Clause |
| [IBM Plex](https://www.ibm.com/plex/) | — | SIL Open Font License 1.1 |
| [Google UUID](https://github.com/google/uuid) | v1.6.0 | BSD-3-Clause |
| [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) | v1.51.0 | BSD-3-Clause |
| [go-playground/validator](https://github.com/go-playground/validator) | v10.30.3 | MIT |

All dependencies are compatible with commercial use.

---

## License

This project is licensed under the **MIT License** — see [LICENSE](LICENSE) for details.

---

## Standard Reference

All calculations implement the methods described in:

**UNECE Code for the Uniform Application of the Rules for the Measurement of Bulk Cargoes, January 1992**

Available from the United Nations Economic Commission for Europe.
