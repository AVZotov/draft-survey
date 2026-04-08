# Calculation Module — Phase 1.1

**Location:** `internal/calculation/`  
**Standard:** UNECE 1992 Draught Survey Code + DSGear_Comments  

---

## Calculation Flow

```
Marks (6 readings)
    │
    ▼
MeanDrafts
(FWD/MID/AFT mean)
    │
    ▼
PPCorrections  ◄─── Vessel (LBP, dFWD/dMID/dAFT, directions)
(Full LBP or Half LBP)
    │
    ▼
DraftsWKeel
(+ keel correction)
    │
    ├──────────────────────┐
    ▼                      ▼
CalcMMC              TrueTrim = AFT - FWD
(Quarter Mean)
    │
    ▼
CalcHydrostatics ◄─── HydrostaticRows (2 rows bracketing MMC)
(Displacement, TPC, LCF)
    │
    ├─────────────────────────────────┐
    ▼                                 ▼
CalcFirstTrimCorrection       CalcSecondTrimCorrection
(TrueTrim × TPC × LCF / LBP)  (TrueTrim² × ΔMTC / LBP)
    │                                 │
    └──────────────┬──────────────────┘
                   │
                   ▼
           CalcListCorrection
           (6 × ΔMid × ΔTPC)
                   │
                   ▼
         CalcDensityCorrection
         (Displacement × (ρ - 1.025) / 1.025)
                   │
                   ▼
          CalcNetDisplacement
          (Displacement + all corrections - Deductibles)
                   │
                   ▼
           CalcCargoWeight
           |NetDispl_fin - NetDispl_ini|
```

---

## Formulas

### 1. Mean Drafts
```
meanF = (FWD_port + FWD_starboard) / 2
meanM = (MID_port + MID_starboard) / 2
meanA = (AFT_port + AFT_starboard) / 2
```

### 2. PP Corrections

**Full LBP method:**
```
LBM = LBP - dAft_signed + dFwd_signed
FWD_corr = dFwd_signed × (meanA - meanF) / LBM
MID_corr = dMid_signed × (meanA - meanF) / LBM
AFT_corr = dAft_signed × (meanA - meanF) / LBM
```

**Half LBP method** *(river vessels only)*:
```
LBMmid-fwd = (LBP/2) - dMid_signed - dFwd_signed
LBMaft-mid = (LBP/2) - dAft_signed - dMid_signed

FWD_corr = dFwd_signed × (meanM - meanF) / LBMmid-fwd
MID_corr = dMid_signed × (meanM - meanF) / LBMmid-fwd
midWKeel  = meanM + MID_corr - (KeelMID / 1000)
AFT_corr  = dAft_signed × (meanA - midWKeel) / LBMaft-mid
```

Sign convention for distances:
- Direction `A` (Aft of perpendicular) → negative
- Direction `F` (Forward of perpendicular) → positive

### 3. Drafts with Keel Correction
```
FWD_wKeel = meanF + FWD_corr - (KeelFWD / 1000)
MID_wKeel = meanM + MID_corr - (KeelMID / 1000)
AFT_wKeel = meanA + AFT_corr - (KeelAFT / 1000)
```
*Keel values are in mm, converted to meters.*

### 4. Quarter Mean (MMC)
```
Marine:  MMC = (FWD + 6×MID + AFT) / 8
River:   MMC = (FWD + 4×MID + AFT) / 6
Barge:   MMC = (3×FWD + 14×MID + 3×AFT) / 20
```

### 5. Interpolation
```
value = lower_value + (fact - lower_draft) × (upper_value - lower_value) / (upper_draft - lower_draft)
```
Used for: Displacement, TPC, LCF from hydrostatic table rows.

### 6. Hydrostatics (LCF formats)

Two formats exist in real vessel hydrostatic tables:

**Format 1 — from Midship (F/A direction):**
- Small values, e.g. `6.93 F`
- Condition: `LCF <= LBP × k3` where `k3 = 0.045`
- Sign: `F` → negative, `A` → positive

**Format 2 — from AP (no direction):**
- Large values, e.g. `98.457`
- Condition: `LCF > LBP × k3`
- Convert: `xf = (LBP/2) - LCF`

Both formats produce the same signed `LCF` result for subsequent calculations.

### 7. First Trim Correction
```
trueTrim = AFT_wKeel - FWD_wKeel
FTC = ABS(trueTrim × TPC × LCF × 100 / LBP)
```
Sign rule:
- `trueTrim < 0` AND `LCF >= 0` → negative
- `trueTrim > 0` AND `LCF <= 0` → negative
- otherwise → positive

### 8. Second Trim Correction
```
ΔMTC = MTC(draft+0.5) - MTC(draft-0.5)
STC  = 50 × trueTrim² × ΔMTC / LBP
```
*Always positive.*

### 9. List Correction
```
ListCorr = 6 × |MID_port - MID_starboard| × |TPC_port - TPC_starboard|
```
*Zero if MID_port == MID_starboard.*

### 10. Density Correction
```
Disp_trimmed = Displacement + FTC + STC + ListCorr
DensityCorr  = Disp_trimmed × (ρ_actual - 1.025) / 1.025
```
*Negative when ρ_actual < 1.025 (fresh/brackish water).*

### 11. Net Displacement
```
Disp_density = Displacement + FTC + STC + ListCorr + DensityCorr
NetDispl     = Disp_density - TotalDeductibles
```

### 12. Cargo Weight
```
CargoWeight = |NetDispl_final - NetDispl_initial|
```

---

## Types

| Type | Description |
|------|-------------|
| `Marks` | 6 raw draft readings (FWD/MID/AFT × Port/Starboard) |
| `MeanDraft` | Averaged FWD/MID/AFT drafts |
| `PPCorrections` | Corrections to perpendiculars |
| `DraftsWKeel` | Drafts corrected for PP + keel |
| `Hydrostatics` | Interpolated Displacement, TPC, LCF |
| `HydrostaticRow` | Single row from vessel's hydrostatic table |
| `MTCRow` | Single MTC value at a given draft |
| `Vessel` | Vessel particulars (LBP, PP distances, keel, type) |
| `VesselType` | `marine` / `river` / `barge` |
| `LCFDirection` | `F` / `A` / `AP` |
| `PPDirection` | `F` / `A` |
| `CorrectionMethod` | `Full LBP` / `Half LBP` |


---

## Test Coverage

| File | Type | Data source |
|------|------|-------------|
| `calculation_test.go` | Unit tests | DSGear.xlsm (LBP=182m vessel) |
| `polar_star_test.go` | Integration tests | POLAR_STAR scenarios (LBP=183m vessel) |

**POLAR_STAR scenarios covered:**
- `TrimnoList` — trim only, LCF from AP format
- `TrimList` — trim + list, full correction chain
