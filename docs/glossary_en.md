# Draft Survey Glossary — English
*A practical reference for surveyors*

---

## Table of Contents

1. [The Vessel](#1-the-vessel)
2. [Draft Marks & Readings](#2-draft-marks--readings)
3. [Corrections to Marks](#3-corrections-to-marks)
4. [Mean Drafts & Key Calculations](#4-mean-drafts--key-calculations)
5. [Hydrostatic Data](#5-hydrostatic-data)
6. [Trim & List Corrections](#6-trim--list-corrections)
7. [Density Correction](#7-density-correction)
8. [Deductibles](#8-deductibles)
9. [Final Results](#9-final-results)
10. [Tank Volume Calculations](#10-tank-volume-calculations)

---

## 1. The Vessel

### LBP — Length Between Perpendiculars
**What it is:** The distance between the forward perpendicular (FP) and the aft perpendicular (AP) of the vessel. These are vertical lines at the bow and stern of the ship at the design waterline.

**Formula:** Given in vessel documentation (metres)

**Why it matters:** Almost every correction formula uses LBP as the base length. Without it, trim corrections cannot be calculated.

---

### LBM — Length Between Marks
**What it is:** The actual distance between the forward and aft draft marks on the hull. Draft marks are rarely located exactly at the perpendiculars — they are usually offset by a few metres.

**Formula (Full LBP method):**
```
LBM = LBP - dAft + dFwd
```
Where `dAft` and `dFwd` are signed distances from the marks to the perpendiculars.

**Why it matters:** Using LBP instead of LBM when marks are offset introduces error into PP corrections. The bigger the offset, the bigger the error.

---

### Lightship
**What it is:** The weight of the vessel itself — hull, machinery, equipment — with no cargo, fuel, ballast, crew, or stores on board. A fixed value from the vessel's stability booklet.

**Formula:** Given in vessel documentation (MT)

**Why it matters:** Lightship is subtracted from displacement to find the total weight of everything loaded on the vessel (DWT).

---

### Summer Draft / Summer DWT
**What it is:** The maximum allowable draft and deadweight at the summer load line. Marked on the hull as "S".

**Why it matters:** If the calculated MMC (corrected for density) exceeds the Summer Draft, the vessel is overloaded. The program issues a warning.

---

### Table Density
**What it is:** The water density used by the shipbuilder when constructing the hydrostatic tables. Usually 1.025 t/m³ (salt water) but can differ.

**Why it matters:** If the actual water density differs from the table density, a correction must be applied to the displacement. Without knowing the table density, the correction cannot be calculated correctly.

---

### Breadth
**What it is:** The maximum beam (width) of the vessel at the midship section (metres).

**Formula:** Given in vessel documentation

**Why it matters:** Used to calculate the list angle in degrees from the difference in port and starboard midship drafts.

---

## 2. Draft Marks & Readings

### Marks — Draft Readings
**What it is:** The six draft readings taken by the surveyor at the waterline on both sides of the vessel: Forward Port, Forward Starboard, Midship Port, Midship Starboard, Aft Port, Aft Starboard.

**Why it matters:** These are the raw measurements — everything else is calculated from them. Reading marks accurately is the most important skill in draft surveying.

---

### Mark Read (Direct / By Waterline)
**What it is:** How the draft was read. "Direct" means the surveyor read the mark directly from a boat or ladder. "By Waterline" means read from the vessel's side using the waterline level.

**Why it matters:** Recorded in the survey report to document method used. Affects credibility of the survey in disputed cases.

---

### dFwd / dMid / dAft — Distance to Perpendicular
**What it is:** The distance (in metres) between the draft mark and the nearest perpendicular (bow or stern). Direction is indicated as "F" (forward of perpendicular) or "A" (aft of perpendicular).

**Why it matters:** Used to calculate PP corrections — adjusting the observed draft at the mark to what it would be at the perpendicular.

---

### Keel (Fk / Mk / Ak)
**What it is:** The thickness of the keel plate at the forward, midship, and aft marks, measured in millimetres. In most surveys this is zero or the same at all three points.

**Formula:**
```
Keel correction = -keel_thickness / 1000  (converts mm to m, always negative)
```

**Why it matters:** The draft marks measure to the outside of the keel. The keel plate has thickness — subtracting it gives the true water depth to the bottom of the vessel.

---

## 3. Corrections to Marks

### PP Corrections (Perpendicular Corrections)
**What it is:** Adjustments applied to the mean drafts at each mark position to calculate what the draft would be at the actual perpendiculars (bow and stern). Needed because draft marks are rarely at the exact perpendicular positions.

**Formula (Full LBP method):**
```
FWD correction = dFwd × trim / LBM
MID correction = dMid × trim / LBM
AFT correction = dAft × trim / LBM
```

**Formula (Half LBP method — river vessels):**
```
FWD correction = dFwd × (meanM - meanF) / LBMmid-fwd
MID correction = dMid × (meanM - meanF) / LBMmid-fwd
AFT correction = dAft × (meanA - midWKeel) / LBMaft-mid
```

**Why it matters:** Without PP corrections, the calculated trim is based on mark positions, not the true perpendiculars. This leads to errors in displacement corrections, especially on vessels with large mark offsets.

---

### Drafts w/Keel (Corrected Drafts)
**What it is:** The draft at each position after applying both PP corrections and keel corrections.

**Formula:**
```
FWD wKeel = meanF + FWD PP correction - (Fk / 1000)
MID wKeel = meanM + MID PP correction - (Mk / 1000)
AFT wKeel = meanA + AFT PP correction - (Ak / 1000)
```

**Why it matters:** These are the "true" drafts used in all subsequent calculations.

---

## 4. Mean Drafts & Key Calculations

### Mean Drafts (meanF, meanM, meanA)
**What it is:** The arithmetic average of port and starboard readings at each mark position.

**Formula:**
```
meanF = (Fwd Port + Fwd Starboard) / 2
meanM = (Mid Port + Mid Starboard) / 2
meanA = (Aft Port + Aft Starboard) / 2
```

**Why it matters:** Averaging port and starboard eliminates the effect of list on the draft reading.

---

### Observed Trim
**What it is:** The difference in mean drafts between aft and forward, before any corrections.

**Formula:**
```
Observed Trim = meanA - meanF
```

**Why it matters:** Quick visual indicator of how the vessel is sitting. Positive = trim by stern (normal). Negative = trim by head (unusual, may affect calculations).

---

### True Trim
**What it is:** The difference between aft and forward corrected drafts (after PP and keel corrections).

**Formula:**
```
True Trim = AFT wKeel - FWD wKeel
```

**Why it matters:** This is the trim value used in all correction formulas. More accurate than observed trim because it accounts for mark positions and keel.

---

### List (listMeters / listDegrees)
**What it is:** The transverse inclination of the vessel — how much it is leaning to one side.

**Formula:**
```
listMeters  = Mid Port - Mid Starboard
listDegrees = arctan(listMeters / Breadth) × (180 / π)
```

**Why it matters:** List affects the accuracy of draft readings and tank measurements. A list greater than 0.5° requires a Letter of Protest to be issued per UNECE standards.

---

### Deflection (Hogging / Sagging)
**What it is:** The bending of the hull. If the midship is higher than the average of bow and stern — the vessel is hogging (arched like a bridge). If lower — sagging (drooping like a hammock).

**Formula:**
```
Deflection = (MID wKeel - (FWD wKeel + AFT wKeel) / 2) × 100  [cm]
```

- `Deflection > 0` → **Hogging**
- `Deflection < 0` → **Sagging**
- `Deflection = 0` → No deflection

**Why it matters:** Severe deflection affects the accuracy of the MMC calculation and indicates structural stress on the vessel.

---

### Mean F&A
**What it is:** The average of the forward and aft corrected drafts.

**Formula:**
```
Mean F&A = (FWD wKeel + AFT wKeel) / 2
```

**Why it matters:** Intermediate value used to calculate MM and as a reference for deflection.

---

### MM (Mean of Means)
**What it is:** The average of the midship corrected draft and Mean F&A.

**Formula:**
```
MM = (MID wKeel + Mean F&A) / 2
```

**Why it matters:** Intermediate step toward MMC. Provides a quick estimate of the vessel's average draft.

---

### MMC — Mean of Mean Corrected Draft (Quarter Mean)
**What it is:** The single representative draft of the vessel, calculated as a weighted average of the three corrected drafts. This is the draft used to enter the hydrostatic tables.

**Formula (Standard — marine vessels):**
```
MMC = (FWD wKeel + 6 × MID wKeel + AFT wKeel) / 8
```

**Formula (River vessels):**
```
MMC = (FWD wKeel + 4 × MID wKeel + AFT wKeel) / 6
```

**Formula (Barges):**
```
MMC = (3 × FWD wKeel + 14 × MID wKeel + 3 × AFT wKeel) / 20
```

**Why it matters:** MMC is the most important single number in a draft survey. It is used to look up displacement, TPC, and LCF from the vessel's hydrostatic tables. All final results depend on getting MMC right.

---

## 5. Hydrostatic Data

### Displacement
**What it is:** The total weight of water displaced by the vessel — equal to the total weight of the vessel and everything on board. Obtained from the hydrostatic tables by interpolation at MMC.

**Unit:** Metric Tonnes (MT)

**Why it matters:** This is the starting point for calculating cargo weight. The difference in displacement between initial and final surveys equals the cargo loaded or discharged.

---

### TPC — Tonnes Per Centimetre
**What it is:** How many metric tonnes of weight must be added to (or removed from) the vessel to change its mean draft by one centimetre.

**Unit:** MT/cm

**Why it matters:** Used in the First Trim Correction formula. Also useful for quick sanity checks — if displacement changes by X MT, the draft should change by approximately X/TPC centimetres.

---

### LCF / LCA — Longitudinal Centre of Flotation
**What it is:** The fore-and-aft position of the centre of the waterplane area. When weight is added or removed, the vessel pivots around this point.

**Unit:** metres, measured from midship (positive = aft, negative = forward) or from AP

**Formula (when given from AP):**
```
LCF from midship = (LBP / 2) - LCA
```

**Why it matters:** Used in the First Trim Correction. If LCF is forward of midship and the vessel trims by stern, the displacement correction is negative (the forward body is finer). Getting the sign right is critical.

---

### MTC — Moment to Change Trim
**What it is:** The trimming moment (in tonne-metres) required to change the vessel's trim by one centimetre.

**Unit:** MT·m/cm

**Why it matters:** Used to calculate the Second Trim Correction (ΔMTCmethod). The change in MTC between drafts ±0.5m from MMC indicates how non-linear the trim correction is.

---

### ΔMTC — Delta MTC
**What it is:** The difference in MTC values between the draft 0.5m above MMC and the draft 0.5m below MMC.

**Formula:**
```
ΔMTC = MTC(MMC + 0.5m) - MTC(MMC - 0.5m)
```

**Why it matters:** Used exclusively in the Second Trim Correction formula. A larger ΔMTC means the hull form changes significantly with draft — the second correction becomes more important.

---

## 6. Trim & List Corrections

### First Trim Correction
**What it is:** An adjustment to displacement for the fact that the vessel is not on an even keel. The waterplane is not symmetrical fore and aft, so a trimmed vessel displaces a different amount than the hydrostatic tables suggest at that mean draft.

**Formula:**
```
1st Trim = ±|True Trim × TPC × LCF × 100 / LBP|

Sign rule:
  If (trim < 0 AND LCF ≥ 0) OR (trim > 0 AND LCF ≤ 0) → negative
  Otherwise → positive
```

**Why it matters:** This is usually the largest correction in the calculation. On a vessel with 2–3 metres of trim, this can be hundreds of tonnes.

---

### Second Trim Correction
**What it is:** A refinement to the First Trim Correction that accounts for the non-linear change in the hull form with trim (the Nemoto correction).

**Formula:**
```
2nd Trim = 50 × True Trim² × ΔMTC / LBP
```

**Why it matters:** For moderate trim (1–3m), this correction is typically small (a few tonnes). For extreme trim it can become significant. Always positive.

---

### List Correction
**What it is:** An adjustment for the asymmetry introduced by the vessel's list. Port and starboard sides have slightly different waterplane areas when the vessel is heeled.

**Formula:**
```
List Correction = 6 × |Mid Port - Mid Starboard| × |TPC Port - TPC Starboard|
```

**Why it matters:** Usually small unless the list is significant. If the vessel has a large list, this correction can be several tonnes. The program warns when list exceeds 0.5°.

---

## 7. Density Correction

### Density Correction
**What it is:** An adjustment for the difference between the actual water density and the density used to build the hydrostatic tables.

**Formula:**
```
Displacement corrected to density = (Displacement + 1st Trim + 2nd Trim + List correction) × (actual density - table density) / table density
```

**Why it matters:** Fresh river water (1.000) vs salt sea water (1.025) is a 2.5% difference — on a large vessel this can be hundreds or even thousands of tonnes. Ignoring this would make the survey completely wrong.

---

### Dockwater Density
**What it is:** The actual density of the water in which the vessel is floating at the time of the survey. Measured by the surveyor using a hydrometer.

**Unit:** t/m³ (typical range: 1.000 – 1.025)

**Why it matters:** The raw material for the density correction. Must be measured carefully, usually at multiple depths and locations around the vessel.

---

### FWA — Fresh Water Allowance
**What it is:** The change in draft (in mm) when a vessel moves from salt water (1.025) to fresh water (1.000).

**Formula:**
```
FWA = Summer DWT / (4 × Summer TPC)  [mm]
```

**Why it matters:** Used to calculate DWA and as a quick reference for how sensitive this vessel is to density changes.

---

### DWA — Dock Water Allowance
**What it is:** The adjustment to the summer load line for the actual water density at the port.

**Formula:**
```
DWA = FWA × (table density - actual density) / (table density - 1.000)  [mm]
```

**Why it matters:** Tells the surveyor and captain how much the vessel can be loaded beyond or below the summer mark given the actual water density.

---

## 8. Deductibles

### Deductibles (Total Deductibles)
**What it is:** The total weight of everything on the vessel that is NOT cargo — ballast water, fresh water, fuel oil (HFO), diesel (MDO), lube oil, bilge water, sewage, and other consumables.

**Why it matters:** To find the cargo weight, deductibles are subtracted from the net displacement. Errors in deductibles directly cause errors in the final cargo figure. This is a common area of dispute between surveyors.

---

### Ballast Water
**What it is:** Sea water pumped into dedicated ballast tanks to stabilise the vessel when it is not fully loaded.

**Weight formula:**
```
Weight = Volume × Density
```

**Why it matters:** Usually the largest deductible. Ballast tanks must be carefully measured (sounded) and the volume calculated from calibration tables.

---

### Fresh Water
**What it is:** Drinking and utility water stored in dedicated tanks for crew use.

**Weight formula:**
```
Weight = Volume × 1.000
```

**Why it matters:** Fresh water density is always 1.000, making calculation simple. The volume is obtained by sounding.

---

### Sounding
**What it is:** The depth of liquid in a tank, measured by lowering a graduated tape through a sounding pipe. The reading is used to look up the volume in the vessel's calibration tables.

**Unit:** metres (or cm, depending on the table)

**Why it matters:** The raw measurement for all tank calculations. Reading accuracy directly affects the weight calculation. Trim and list of the vessel affect the sounding reading — hence the calibration corrections.

---

## 9. Final Results

### Net Displacement
**What it is:** The total displacement corrected for trim, list, and density — representing the true total weight of the vessel and all its contents.

**Formula:**
```
Net Displacement = Displacement + 1st Trim + 2nd Trim + List correction + Density correction - Total Deductibles
```

**Why it matters:** This is the base from which cargo weight is calculated.

---

### Constant (Ship's Constant)
**What it is:** The difference between net displacement and the sum of lightship plus all known deductibles. Represents everything on board that is not accounted for — dirt in bilges, crew and luggage, stores, rust scale, residual water, etc.

**Formula:**
```
Constant = Net Displacement - Lightship
```
(Net Displacement already has Total Deductibles subtracted out — see the Net Displacement formula above. Subtracting deductibles a second time here is a common mistake.)

**Why it matters:** One of the most important numbers in draft surveying. A stable, reasonable constant confirms the survey is correct. A constant that differs significantly from the declared value raises questions about data accuracy or honesty.

---

### Constant Declared
**What it is:** The ship's constant as declared by the captain at the start of the survey. Based on historical surveys of this vessel.

**Why it matters:** Compared against the calculated constant. A large discrepancy (typically >50 MT) should trigger investigation. Captains sometimes inflate the declared constant to allow more cargo to appear "heavier" than it is.

---

### Cargo Weight
**What it is:** The weight of cargo loaded or discharged during the operation.

**Formula:**
```
Cargo = Final Net Displacement - Initial Net Displacement
```

The app does not force this to a positive number: for a loading operation it comes out positive (cargo added), and for a discharge operation it comes out negative (cargo removed) — the sign itself tells you the direction of the operation.

**Why it matters:** This is the ultimate purpose of the draft survey — to determine how much cargo changed hands. It is used for billing, customs declarations, and cargo claims.

---

### Current DWT — Deadweight
**What it is:** Everything on the vessel except the vessel itself. Cargo + ballast + fuel + stores + crew.

**Formula:**
```
Current DWT = Displacement corrected to density - Lightship
```

**Why it matters:** Compared against the vessel's Summer DWT to check if the vessel is overloaded. If Current DWT approaches Summer DWT, the vessel is near its load limit.

---

## 10. Tank Volume Calculations

### Calibration Table
**What it is:** A table provided by the shipbuilder giving the volume of liquid in each tank for a given sounding, at various trim conditions.

**Why it matters:** Without calibration tables, only a rough estimate of tank volumes is possible. Accurate tank weights require proper interpolation from these tables.

---

### Calibration Table Type 1 — Standard (Volume by Trim)
**What it is:** A two-dimensional table with sounding on one axis and trim on the other. Each cell contains the actual volume in m³.

**How to use:** Find the two sounding rows and two trim columns bracketing the actual values. Interpolate first by trim, then by sounding (bilinear interpolation).

**Why it matters:** Most common type — found on the majority of vessels surveyed.

---

### Calibration Table Type 2 — Sounding Correction
**What it is:** Two tables. Table 1 gives a correction (in mm) to apply to the sounding. Table 2 gives the volume at the corrected sounding (at zero trim).

**How to use:**
1. Interpolate sounding correction from Table 1
2. Apply correction: corrected sounding = actual sounding + correction/1000
3. Look up volume from Table 2 at corrected sounding

**Why it matters:** Less common but important to recognise. Incorrectly treating it as Type 1 gives completely wrong volumes.

---

### Calibration Table Type 3 — Volume Correction
**What it is:** Two tables. Table 2 gives base volume at zero trim. Table 1 gives a volume correction (in m³) for the actual trim.

**How to use:**
1. Look up base volume from Table 2 at actual sounding
2. Interpolate volume correction from Table 1
3. Final volume = base volume + trim correction + list correction

**Why it matters:** Rare but encountered on some older vessels and barges.

---

### TTL / TTU — Table Trim Lower / Upper
**What it is:** The two trim values from the column headers of the calibration table that bracket the vessel's actual trim.

**Unit:** metres

**Why it matters:** Required for bilinear interpolation. The actual trim must fall between TTL and TTU, otherwise extrapolation is used (less reliable).

---

### TLL / TLU — Table List Lower / Upper
**What it is:** The two list values from the column headers of the list correction table that bracket the vessel's actual list angle.

**Unit:** degrees (°)

**Why it matters:** Same as TTL/TTU but for the list correction table. Note that TLL is typically negative (opposite side heel).

---

### TSL / TSU — Table Sounding Lower / Upper
**What it is:** The two sounding rows from the calibration table that bracket the vessel's actual sounding.

**Unit:** metres

**Why it matters:** The actual sounding must fall between TSL and TSU for interpolation to be valid.

---

*Last updated: 2026 | Draft Survey Tool — github.com/AVZotov/draft-survey*
