package fields

import "fmt"

// WithIndex helper function
// for consistent naming and parsing of fields received from HTML requests
func WithIndex(field string, index int) string {
	return fmt.Sprintf("%s-%d", field, index)
}

// WithTankID helper function
// for consistent tanks naming and parsing received from HTML requests
func WithTankID(field string, tankID string) string {
	return fmt.Sprintf("%s-%s", field, tankID)
}

// Profile form fields
const (
	FieldFirstName  = "first_name"
	FieldLastName   = "last_name"
	FieldPosition   = "position"
	FieldEmail      = "email"
	FieldCompany    = "company"
	FieldLicense    = "license_no"
	FieldCountry    = "country"
	FieldEmployeeID = "employee_id"
	FieldSignature  = "signature"
)

// Survey Block
const (
	FieldJobNumber       = "job_number"
	FieldDSNumber        = "ds_number"
	FieldCargoOperation  = "cargo_operation"
	FieldCargo           = "cargo"
	FieldPacking         = "packing"
	FieldClient          = "client"
	FieldCargoDeclared   = "cargo_declared"
	FieldRemarks         = "remarks"
	FieldLoadCountryCode = "load_country_code"
	FieldLoadCountryName = "load_country_name"
)

// Vessel Block
const (
	FieldVesselName      = "vessel_name"
	FieldIMO             = "imo"
	FieldBuiltYear       = "built_year"
	FieldConstDeclared   = "const_declared"
	FieldTableDensity    = "table_density"
	FieldLightship       = "lightship"
	FieldLBP             = "lbp"
	FieldBreadth         = "breadth"
	FieldDepth           = "depth"
	FieldSummerDraft     = "summer_draft"
	FieldSummerDWT       = "summer_dwt"
	FieldSummerTPC       = "summer_tpc"
	FieldSummerFreeboard = "summer_freeboard"
	FieldFlagCountryCode = "flag_country_code"
	FieldFlagCountryName = "flag_country_name"
	FieldHoldsTotal      = "holds_total"
)

// Settings Block
const (
	FieldMMCMethod    = "mmc_method"
	FieldCorrMethod   = "corr_method"
	FieldLCFDetection = "lcf_detection"
)

// Draft Info Block
const (
	FieldLoadPort      = "load_port"
	FieldDischargePort = "discharge_port"
	FieldProductGrade  = "product_grade"
)

// Sea Condition Block
const (
	FieldSeaType      = "sea_type"
	FieldSeaCondition = "sea_condition"
)

// FieldHoldsActive Holds Block
const (
	FieldHoldsActive = "holds_active"
)

// Marks Block
const (
	FieldFwdPort       = "fwd_port"
	FieldFwdPortMethod = "fwd_port_method"
	FieldFwdStbd       = "fwd_stbd"
	FieldFwdStbdMethod = "fwd_stbd_method"
	FieldMidPort       = "mid_port"
	FieldMidPortMethod = "mid_port_method"
	FieldMidStbd       = "mid_stbd"
	FieldMidStbdMethod = "mid_stbd_method"
	FieldAftPort       = "aft_port"
	FieldAftPortMethod = "aft_port_method"
	FieldAftStbd       = "aft_stbd"
	FieldAftStbdMethod = "aft_stbd_method"
	FieldTPCListPort   = "tpc_list_port"
	FieldTPCListStbd   = "tpc_list_stbd"
)

// PP Keel Block
const (
	FieldDFwd    = "d_fwd"
	FieldDFwdDir = "d_fwd_dir"
	FieldDMid    = "d_mid"
	FieldDMidDir = "d_mid_dir"
	FieldDAft    = "d_aft"
	FieldDAftDir = "d_aft_dir"
	FieldKeelFwd = "keel_fwd"
	FieldKeelMid = "keel_mid"
	FieldKeelAft = "keel_aft"
)

// Hydrostatics Block
const (
	FieldUDraft    = "u_draft"
	FieldUDisp     = "u_disp"
	FieldUTPC      = "u_tpc"
	FieldULCF      = "u_lcf"
	FieldULCFDir   = "u_lcf_dir"
	FieldLDraft    = "l_draft"
	FieldLDisp     = "l_disp"
	FieldLTPC      = "l_tpc"
	FieldLLCF      = "l_lcf"
	FieldLLCFDir   = "l_lcf_dir"
	FieldPMTCDraft = "p_mtc_draft"
	FieldPMTC      = "p_mtc"
	FieldNMTCDraft = "n_mtc_draft"
	FieldNMTC      = "n_mtc"
)

// Deductibles Block
const (
	FieldHFO              = "hfo"
	FieldMDO              = "mdo"
	FieldLubOil           = "lub_oil"
	FieldBilgeWater       = "bilge_water"
	FieldDockwaterDensity = "dockwater_density"
	FieldOthers           = "others"
)

// Tank fields
const (
	FieldTankName          = "tank_name"
	FieldTankType          = "tank_type"
	FieldTankSounding      = "tank_sounding"
	FieldTankVolume        = "tank_volume"
	FieldTankDensity       = "tank_density"
	FieldCalibType         = "calib_type"
	FieldHasListCorrection = "has_list_correction"
	FieldTTL               = "ttl"
	FieldTTU               = "ttu"
	FieldTLL               = "tll"
	FieldTLU               = "tlu"
	FieldTrimTSLS          = "trim_table_sounding_low"
	FieldTrimTSLVL         = "trim_table_sounding_low_vol_low"
	FieldTrimTSLVU         = "trim_table_sounding_low_vol_up"
	FieldTrimTSUS          = "trim_table_sounding_up"
	FieldTrimTSUVL         = "trim_table_sounding_up_vol_low"
	FieldTrimTSUVU         = "trim_table_sounding_up_vol_up"
	FieldListTSLS          = "list_table_sounding_low"
	FieldListTSLVL         = "list_table_sounding_low_vol_low"
	FieldListTSLVU         = "list_table_sounding_low_vol_up"
	FieldListTSUS          = "list_table_sounding_up"
	FieldListTSUVL         = "list_table_sounding_up_vol_low"
	FieldListTSUVU         = "list_table_sounding_up_vol_up"
	FieldCorrSoundingLow   = "corr_sounding_low"
	FieldCorrVolumeLow     = "corr_volume_low"
	FieldCorrSoundingUp    = "corr_sounding_up"
	FieldCorrVolumeUp      = "corr_volume_up"
)
