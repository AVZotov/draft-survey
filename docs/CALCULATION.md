# Calculation Module

**Location:** `internal/calculation/`  
**Standard:** UNECE 1992 Code for the Uniform Application of the Rules for the Measurement of Bulk Cargoes  
**Last Updated:** 2026-05-14

---

## Table of Contents

1. [Overview](#1-overview)
2. [Rounding Policy](#2-rounding-policy)
3. [Calculation Flow](#3-calculation-flow)
4. [Formulas](#4-formulas)
   - 4.1 [Mean Drafts](#41-mean-drafts)
   - 4.2 [PP Corrections](#42-pp-corrections)
   - 4.3 [Drafts with Keel Correction](#43-drafts-with-keel-correction)
   - 4.4 [Observed Trim & True Trim](#44-observed-trim--true-trim)
   - 4.5 [Mean F&A and Deflection](#45-mean-fa-and-deflection)
   - 4.6 [List](#46-list)
   - 4.7 [Quarter Mean (MMC)](#47-quarter-mean-mmc)
   - 4.8 [Interpolation](#48-interpolation)
   - 4.9 [Hydrostatics — LCF Detection Modes](#49-hydrostatics--lcf-detection-modes)
   - 4.10 [First Trim Correction](#410-first-trim-correction)
   - 4.11 [Delta MTC](#411-delta-mtc)
   - 4.12 [Second Trim Correction](#412-second-trim-correction)
   - 4.13 [List Correction](#413-list-correction)
   - 4.14 [Density Correction](#414-density-correction)
   - 4.15 [Net Displacement](#415-net-displacement)
   - 4.16 [Cargo Weight & Constants](#416-cargo-weight--constants)
5. [Tank Volume Calibration](#5-tank-volume-calibration)
6. [Entry Points](#6-entry-points)
7. [Types Reference](#7-types-reference)
8. [Test Coverage](#8-test-coverage)
9. [Known Deviations from UNECE 1992](#9-known-deviations-from-unece-1992)

---

## 1. Overview

The calculation module implements the complete draft survey mathematical chain per UNECE 1992 standard. It is the core of the application — all other modules depend on it.

**Key design decisions:**
- Single entry point `CalcDraft()` returns a complete `DraftResult` struct
- Multi-draft orchestration via `CalcSurvey()` returning `SurveyResult`
- All functions are pure — no side effects, no external dependencies
- `*float64` for surveyor-entered values (can be nil = not entered)
- `float64` for calculated values (always present)

---

## 2. Rounding Policy

**Rule:** `round3()` (3 decimal places) is applied to every value that appears as a reported field in **UNECE Form C**.

Intermediate calculations that are NOT recorded in Form C are computed at full `float64` precision without rounding.

This mirrors the physical survey process: a surveyor records rounded values on paper (Form C) and uses those recorded values in subsequent calculations. In any legal dispute, the reported (rounded) values are the reference.

```go
func round3(v float64) float64 {
    return math.Round(v*1000) / 1000
}
```

**Form C fields with round3 applied:**

| Field | Applied |
|-------|---------|
| meanF, meanM, meanA | ✅ |
| LBM | ✅ |
| PP corrections (fwd/mid/aft) | ✅ |
| DraftsWKeel (fwd/mid/aft) | ✅ |
| ObservedTrim, TrueTrim | ✅ |
| MeanF&A, Deflection | ✅ |
| ListMeters, ListDegrees | ✅ |
| MMC | ✅ |
| Displacement, TPC, LCF | ✅ |
| DeltaMTC | ✅ |
| 1st/2nd Trim Correction | ✅ |
| List Correction | ✅ |
| Density Correction | ✅ |
| Net Displacement | ✅ |
| Constant | ✅ |

---

## 3. Calculation Flow

```
Marks (6 readings: FWD/MID/AFT × Port/Starboard)
    │
    ▼
MeanDrafts ──────────────────────────── round3
(meanF, meanM, meanA)
    │
    ▼
LBM = LBP - dAft_signed + dFwd_signed ─ round3
    │
    ▼
PPCorrections ◄── LBP, PP distances, directions
(Full LBP or Half LBP method) ────────── round3
    │
    ▼
DraftsWKeel ◄── keel corrections (mm→m)
(fwdWK, midWK, aftWK) ─────────────────  round3
    │
    ├──── ObservedTrim = meanA - meanF ── round3
    ├──── TrueTrim = aftWK - fwdWK ───── round3
    ├──── MeanF&A = (fwdWK + aftWK)/2 ── round3
    ├──── Deflection = (midWK-MeanF&A)×100 round3
    ├──── ListMeters = midP - midS ────── round3
    ├──── ListDegrees = atan(list/B)×180/π round3
    │
    ▼
CalcMMC ─────────────────────────────── round3
(Quarter Mean: 6M/8, 4M/6, or 3/10f)
    │
    ├──── MTCPlusDraft = MMC + 0.5 ──── (hint for surveyor)
    ├──── MTCMinusDraft = MMC - 0.5 ─── (hint for surveyor)
    │
    ▼
CalcHydrostatics ◄── HydrostaticRows (2 rows)
(Displacement, TPC, LCF) ──────────── round3 via Interpolate
    │
    ├──── CalcFirstTrimCorrection ────── round3
    ├──── CalcDeltaMTC ──────────────── round3
    ├──── CalcSecondTrimCorrection ───── round3
    ├──── CalcListCorrection ─────────── round3
    │
    ▼
CalcDensityCorrection ─────────────── round3
    │
    ▼
CalcTotalDeductibles
    │
    ▼
CalcNetDisplacement ──────────────── round3
    │
    ├──── CalcConstant = NetDispl - Lightship
    └──── CalcCurrentDWT = DisplCorrToDensity - Lightship
```

---

## 4. Formulas

### 4.1 Mean Drafts

```
meanF = (FWD_port + FWD_starboard) / 2
meanM = (MID_port + MID_starboard) / 2
meanA = (AFT_port + AFT_starboard) / 2
```

*UNECE Form C lines 131, 138, 147*

---

### 4.2 PP Corrections

**Sign convention for distances:**
- Direction `F` (Forward of perpendicular) → **positive**
- Direction `A` (Aft of perpendicular) → **negative**

**Full LBP method** *(marine vessels, default)*:
```
LBM = LBP - dAft_signed + dFwd_signed

FWD_corr = dFwd_signed × (meanA - meanF) / LBM
MID_corr = dMid_signed × (meanA - meanF) / LBM
AFT_corr = dAft_signed × (meanA - meanF) / LBM
```

**Half LBP method** *(river vessels)*:
```
LBMmid-fwd = (LBP/2) - dMid_signed - dFwd_signed
LBMaft-mid = (LBP/2) - dAft_signed - dMid_signed

FWD_corr = dFwd_signed × (meanM - meanF) / LBMmid-fwd
MID_corr = dMid_signed × (meanM - meanF) / LBMmid-fwd
midWKeel  = meanM + MID_corr - (KeelMID / 1000)
AFT_corr  = dAft_signed × (meanA - midWKeel) / LBMaft-mid
```

*UNECE Form C lines 133-134, 140-141, 148-149*

---

### 4.3 Drafts with Keel Correction

```
FWD_wKeel = meanF + FWD_corr - (KeelFWD / 1000)
MID_wKeel = meanM + MID_corr - (KeelMID / 1000)
AFT_wKeel = meanA + AFT_corr - (KeelAFT / 1000)
```

*Keel values are in mm, converted to meters. Always negative correction.*

---

### 4.4 Observed Trim & True Trim

```
ObservedTrim = meanA - meanF
TrueTrim     = AFT_wKeel - FWD_wKeel
```

- Positive = trim by stern (normal)
- Negative = trim by head (unusual)

---

### 4.5 Mean F&A and Deflection

```
MeanF&A    = (FWD_wKeel + AFT_wKeel) / 2
Deflection = (MID_wKeel - MeanF&A) × 100   [cm]
```

- `Deflection > 0` → **Hogging**
- `Deflection < 0` → **Sagging**

---

### 4.6 List

```
ListMeters  = MID_port - MID_starboard
ListDegrees = arctan(ListMeters / Breadth) × (180 / π)
```

Alert thresholds:
- `|ListDegrees| >= 0.5°` → LOP recommended (UNECE/DSGear)
- `|ListDegrees| >= 1.0°` → beyond calibration table range

---

### 4.7 Quarter Mean (MMC)

```
Marine:  MMC = (FWD + 6 × MID + AFT) / 8     (6M/8)
River:   MMC = (FWD + 4 × MID + AFT) / 6     (4M/6)
Barge:   MMC = (3×FWD + 14×MID + 3×AFT) / 20 (3/10f)
```

**MTC Draft Hints** (shown to surveyor when all marks entered):
```
MTCPlusDraft  = round3(MMC + 0.5)
MTCMinusDraft = round3(MMC - 0.5)
```

---

### 4.8 Interpolation

```
value = lowerValue + (fact - lowerX) × (upperValue - lowerValue) / (upperX - lowerX)
```

---

### 4.9 Hydrostatics — LCF Detection Modes

**Auto mode** (`IsLcfDetectionManual = false`, default):
```
k3 = 0.045
IF upper.LCF > LBP × k3 THEN
    LCF is from AP → xf = (LBP/2) - LCF
ELSE
    Direction F → negative, Direction A → positive
```

**Manual mode** (`IsLcfDetectionManual = true`, UNECE standard):
- Trust direction entered by surveyor
- `F` → negative, `A` → positive, `AP` → convert from AP

*Manual mode used in golden tests for strict UNECE compliance.*

---

### 4.10 First Trim Correction

```
FTC = ABS(TrueTrim × TPC × LCF × 100 / LBP)
```

Sign rule:

| TrueTrim | LCF | Sign |
|----------|-----|------|
| > 0 (Aft) | < 0 (Fwd) | Negative |
| < 0 (Fwd) | > 0 (Aft) | Negative |
| > 0 (Aft) | > 0 (Aft) | Positive |
| < 0 (Fwd) | < 0 (Fwd) | Positive |

---

### 4.11 Delta MTC

```
ΔMTC = MTC(MMC + 0.5m) - MTC(MMC - 0.5m)
```

---

### 4.12 Second Trim Correction

```
STC = 50 × TrueTrim² × ΔMTC / LBP
```

*Always positive. UNECE Form C line 163.*

---

### 4.13 List Correction

```
ListCorr = 6 × |MID_port - MID_starboard| × |TPC_port - TPC_starboard|
```

**Method selection:**

| Condition | Method |
|-----------|--------|
| `TPCListPort` and `TPCListStarboard` entered | **V1** — direct formula |
| Both nil + hydrostatic TPC upper == TPC lower | **V2** — Summer TPC interpolation |
| Both nil + TPC differ | `0` |

**V2 formula:**
```
TPC_port = hydroTPC + (midPort - MMC) × (SummerTPC - hydroTPC) / (SummerDraft - MMC)
TPC_stbd = hydroTPC + (midStbd - MMC) × (SummerTPC - hydroTPC) / (SummerDraft - MMC)
ListCorr = 6 × |midPort - midStbd| × |TPC_port - TPC_stbd|
```

---

### 4.14 Density Correction

```
DisplCorrToTrimList = Displacement + FTC + STC + ListCorr
DensityCorr = DisplCorrToTrimList × (ρ_actual - ρ_table) / ρ_table
```

---

### 4.15 Net Displacement

```
DisplCorrToDensity = Displacement + FTC + STC + ListCorr + DensityCorr
NetDisplacement    = DisplCorrToDensity - TotalDeductibles
```

---

### 4.16 Cargo Weight & Constants

```
CargoWeight  = |NetDispl_final - NetDispl_initial|
Constant     = NetDisplacement - Lightship
TrueConstant = Constant - CumulativeCargo
CurrentDWT   = DisplCorrToDensity - Lightship
```

**TrueConstant** should be stable across all drafts. Significant deviation indicates survey error.

---

## 5. Tank Volume Calibration

**Notation:**
- `AS` — actual sounding (m)
- `TT` — true trim (m)
- `TTL/TTU` — table trim lower/upper (m)
- `TLL/TLU` — table list lower/upper (degrees)
- `TSL/TSU` — table sounding lower/upper (m)

**Base bilinear interpolation:**
```
AB1 = A1 + (TT - TTL) × (B1 - A1) / (TTU - TTL)
AB2 = A2 + (TT - TTL) × (B2 - A2) / (TTU - TTL)
AB  = AB1 + (AS - TSL) × (AB2 - AB1) / (TSU - TSL)
```

If trim = 0: 1D interpolation by sounding only.

### Type 1 — Volume by Trim *(most common)*
Table contains volume (m³) for each (sounding × trim).
```
Volume = bilinear interpolation by trim then sounding
If list table: Volume += bilinear interpolation by list then sounding
```

### Type 2 — Sounding Correction *(rare)*
```
SoundingCorr = interpolate from Table 1 [mm] → round to integer mm
If list: SoundingCorr += list correction [mm]
CorrectedSounding = ActualSounding + SoundingCorr / 1000
Volume = 1D interpolation from Table 2 at CorrectedSounding
```

### Type 3 — Volume Correction *(rare)*
```
BaseVolume = 1D interpolation from Table 2 at ActualSounding
TrimCorr   = interpolate from Table 1 by trim [m³]
If list: ListCorr = interpolate from Table 1 by list [m³]
Volume = BaseVolume + TrimCorr + ListCorr
```

---

## 6. Entry Points

### 6.1 CalcDraft

```go
func CalcDraft(draft types.Draft, v types.VesselData) types.DraftResult
```

**DraftResult key fields:**

| Field | Description |
|-------|-------------|
| `LBM` | Length between marks |
| `MeanDraft` | FWD/MID/AFT mean drafts |
| `MMC` | Quarter mean draft |
| `MTCPlusDraft` / `MTCMinusDraft` | MTC lookup hints |
| `TrueTrim` / `ObservedTrim` | Trim values |
| `ListMeters` / `ListDegrees` | List values |
| `Deflection` | Hogging/sagging in cm |
| `Hydrostatics` | Displacement, TPC, LCF |
| `DeltaMTC` | MTC+ - MTC- |
| `FirstTrimCorrection` / `SecondTrimCorrection` / `ListCorrection` | Corrections |
| `DensityCorrection` | Density correction |
| `DisplacementCorrected` | After all corrections |
| `NetDisplacement` | After deductibles |
| `Constant` / `TrueConstant` (via CalcSurvey) | Ship's constant |
| `CurrentDWT` | Current deadweight |

### 6.2 CalcSurvey

```go
func CalcSurvey(s types.Survey) types.SurveyResult
```

**SurveyResult:**

| Field | Description |
|-------|-------------|
| `DraftTotals[]` | Per-draft: DraftResult + CargoFromPrev + ConstantDiff + TrueConstant |
| `CargoOnBoard` | Final - Initial net displacement |
| `CargoDiffDeclared` | CargoOnBoard - CargoDeclared |

---

## 7. Types Reference

| Type | Description |
|------|-------------|
| `Marks` | 6 raw draft readings |
| `VesselData` | Vessel particulars incl. LBP, breadth, lightship, summer data |
| `HydrostaticRow` | Draft, Displacement, TPC, LCF, LCFDirection |
| `MTCRow` | Draft, MTC value |
| `Tank` | BWT or FWT with optional calibration data |
| `VolumeCalibrationData` | Calibration table (Type 1/2/3, trim/list/volume rows) |
| `CalibrationRow` | 2D row: Sounding, VolumeLow, VolumeUp |
| `SoundingVolume` | 1D row: TableSounding, TableVolume |
| `DraftResult` | Complete single-draft calculation output |
| `SurveyResult` | Complete multi-draft survey output |

---

## 8. Test Coverage

| File | Vessel | Tests |
|------|--------|-------|
| `beam_initial_test.go` | BEAM IMO 9591741 (LBP=283.5m) | 27 golden tests against UNECE 1992 |
| `unece_formulas_test.go` | — | UNECE reference implementation (immutable) |
| `polar_star_test.go` | POLAR_STAR (LBP=183m) | TrimNoList + TrimList scenarios |
| `calculation_test.go` | DSGear reference (LBP=182m) | Unit tests full chain |
| `tank_volume_test.go` | — | Type 1/2/3 volume calibration |

**Coverage:** >80% for `internal/calculation/`

---

## 9. Known Deviations from UNECE 1992

| Item | UNECE | Implementation | Reason |
|------|-------|----------------|--------|
| LCF auto-detection | Not specified | k3=0.045 coefficient | DSGear UX extension |
| List correction V2 | Not specified | Summer TPC interpolation | When hydrostatic TPC rows are equal |
| Intermediate rounding | Not specified | round3 on all Form C fields | Mirrors physical paper survey |

---

*See also: `docs/glossary_en.md` and `docs/glossary_ru.md` for term definitions.*