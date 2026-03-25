package constants

// General fixed names for tanks parsing on backend
const (
	// TankType - from vessel's passport data
	TankType = "w_tank_type"
	// TankName - inputted manually from surveyor
	TankName = "w_tank_name"
	// TankSounding — measured by surveyor directly from tank
	TankSounding = "tank_sounding" // measured by ullage/sounding tape, m
	// TankDensity — measured by surveyor using water sample from tank
	TankDensity = "tank_density" // seawater density, t/m³ (BWT only)
	// TankVolume — calculated using tank calibration data or entered manually as exception
	TankVolume = "tank_volume" // m³
	// TankWeight - calculated value
	TankWeight = "tank_weight" //MT
	// HasListCorrection - bool to show if tank calibration data contains List corrections
	HasListCorrection = "has_list_correction"
)

// Fixed names used in names and ID generation for parsing calibration data passed
// from modal window
// Type 1 calibration approach (most common) with corrections on Trim and Volume
const (
	// VolumeCalibrationType1 - most common type of calibration tables at marine vessels
	VolumeCalibrationType1 = "standard_volume_by_trim"
	// TTL - table trim lower from vessel's data
	TTL = "table_trim_low"
	// TTU - table trim up from vessel's data
	TTU = "table_trim_up"
	// TLL - table list lower from vessel's data
	TLL = "table_list_low"
	// TLU - table list up from vessel's data
	TLU = "table_list_up"
	// TrimTableTSLS - lower sounding from vessel's trim table TTL column
	TrimTableTSLS = "trim_table_table_sounding_low_sounding"
	// TrimTableTSLVL - lower volume from vessel's trim table TTL column
	TrimTableTSLVL = "trim_table_table_sounding_low_volume_low"
	// TrimTableTSLVU - upper volume from vessel's trim table TTL column
	TrimTableTSLVU = "trim_table_table_sounding_low_volume_up"
	// TrimTableTSUS - upper sounding from vessel's trim table TTU column
	TrimTableTSUS = "trim_table_table_sounding_up_sounding"
	// TrimTableTSUVL - lower volume from vessel's trim table TTU column
	TrimTableTSUVL = "trim_table_table_sounding_up_volume_low"
	// TrimTableTSUVU - upper volume from vessel's trim table TTU column
	TrimTableTSUVU = "trim_table_table_sounding_up_volume_up"
	// ListTableTSLS - lower sounding from vessel's list table TTL column
	ListTableTSLS = "list_table_table_sounding_low_sounding"
	// ListTableTSLVL - lower volume from vessel's list table TTL column
	ListTableTSLVL = "list_table_table_sounding_low_volume_low"
	// ListTableTSLVU - upper volume from vessel's list table TTL column
	ListTableTSLVU = "list_table_table_sounding_low_volume_up"
	// ListTableTSUS - upper sounding from vessel's list table TTU column
	ListTableTSUS = "list_table_table_sounding_up_sounding"
	// ListTableTSUVL - lower volume from vessel's list table TTU column
	ListTableTSUVL = "list_table_table_sounding_up_volume_low"
	// ListTableTSUVU - upper volume from vessel's list table TTU column
	ListTableTSUVU = "list_table_table_sounding_up_volume_up"
)

// Fixed names used in names and ID generation for parsing calibration data passed
// from modal window
// ===== TYPE 2: Sounding Correction =====
// Table 1 — sounding correction by trim (same structure as Type 1 trim table)
// Reuse: TTL, TTU, TrimTableTSLS TrimTableTSLVL TrimTableTSLVU TrimTableTSUS TrimTableTSUVL TrimTableTSUVU
// Reuse: ListTableTSLS ListTableTSLVL ListTableTSLVU ListTableTSUS ListTableTSUVL ListTableTSUVU (if list correction exists)
// Table 2 — volume at corrected sounding (1D, no trim columns)
// ===== TYPE 3: Volume Correction =====
// Table 1 — volume correction by trim (same structure as Type 1 trim table)
// Reuse: TTL, TTU, TrimTableTSLS TrimTableTSLVL TrimTableTSLVU TrimTableTSUS TrimTableTSUVL TrimTableTSUVU
// Reuse: ListTableTSLS ListTableTSLVL ListTableTSLVU ListTableTSUS ListTableTSUVL ListTableTSUVU (if list correction exists)
// Table 2 — base volume at zero trim (1D, no trim columns)
// Both uses same field constants as table is identical
const (
	// VolumeCalibrationType2 - type of calibration data with 1D volume at corrected soundings
	VolumeCalibrationType2 = "sounding_correction"
	// VolumeCalibrationType3 - type of calibration data based on volume at zero trim (1D, no trim columns)
	VolumeCalibrationType3 = "volume_correction"
	// CorrTableSL - lower sounding from vessel's correction table
	CorrTableSL = "corr_table_sounding_low"
	// CorrTableVL - lower volume from vessel's correction table
	CorrTableVL = "corr_table_volume_low"
	// CorrTableSU - upper sounding from vessel's correction table
	CorrTableSU = "corr_table_sounding_up"
	// CorrTableVU - upper volume from vessel's correction table
	CorrTableVU = "corr_table_volume_up"
)
