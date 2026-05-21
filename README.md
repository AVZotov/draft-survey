# Draft Survey Tool

> A professional cargo weight calculator for marine surveyors, implementing the UNECE 1992 Draft Survey Code.

**Status:** 🚧 Active Development — v0.4.0 | Core calculations complete and verified against UNECE 1992 standard

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
- Real-time calculation panel updated as you type
- SVG vessel diagram with dynamic trim and list visualization
- Sea condition logging (Wave / Ice with detailed condition selection)
- BW and FW tank management with calibration table modal
- Survey status alerts (constant deviation, data warnings)
- Surveyor profile with company details
- Graceful shutdown — safe database close on exit

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | [Go 1.25](https://go.dev/) + [Fiber v2.52](https://gofiber.io/) |
| Templates | [Templ v0.3.977](https://templ.guide/) |
| Interactivity | [HTMX 2.0.8](https://htmx.org/) + [Alpine.js 3.15.8](https://alpinejs.dev/) |
| DOM Morphing | [Idiomorph](https://github.com/bigskysoftware/idiomorph) |
| Fonts | [IBM Plex Sans & Mono](https://www.ibm.com/plex/) |
| Storage | [SQLite](https://www.sqlite.org/) via [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) (pure Go, no CGO) |

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

## Project Structure

```
cmd/server/         — application entry point
internal/
  calculation/      — UNECE 1992 math engine
  types/            — shared domain types
  storage/          — repository pattern (SQLite + embed.FS)
  handler/          — HTTP handlers (Fiber)
  format/           — value formatting helpers
  constants/        — form field names, header constants
data/
  dictionaries/     — static data (ports, countries, conditions)
docs/               — documentation and glossary
web/
  templates/        — page templates (Templ)
  widgets/          — reusable UI components
  components/       — shared UI elements
  static/           — CSS, JS, fonts
embed.go            — go:embed directives (self-contained build)
```

---

## Roadmap

See [docs/ROADMAP.md](docs/ROADMAP.md) for the full development plan.

**Coming next (v0.5.0):**
- PDF report generation (UNECE-compliant format)
- Structured logging
- Unified error handling
- Pre-built binaries for Windows, macOS, Linux

---

## Acknowledgements

This project is built with excellent open source tools. Grateful to their authors and contributors:

| Project | Version | License |
|---|---|---|
| [Go](https://go.dev/) | 1.25 | BSD-3-Clause |
| [Fiber](https://gofiber.io/) | v2.52.11 | MIT |
| [Templ](https://templ.guide/) | v0.3.977 | MIT |
| [HTMX](https://htmx.org/) | 2.0.8 | BSD Zero Clause |
| [Alpine.js](https://alpinejs.dev/) | 3.15.8 | MIT |
| [Idiomorph](https://github.com/bigskysoftware/idiomorph) | latest | BSD Zero Clause |
| [IBM Plex](https://www.ibm.com/plex/) | — | SIL Open Font License 1.1 |
| [Google UUID](https://github.com/google/uuid) | v1.6.0 | BSD-3-Clause |
| [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) | v1.50.1 | BSD-3-Clause |

All dependencies are compatible with commercial use.

---

## License

This project is licensed under the **MIT License** — see [LICENSE](LICENSE) for details.

---

## Standard Reference

All calculations implement the methods described in:

**UNECE Code for the Uniform Application of the Rules for the Measurement of Bulk Cargoes, January 1992**

Available from the United Nations Economic Commission for Europe.
