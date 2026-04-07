package constants

// HX Custom headers (keys)
const (
	HXTankType    = "X-Tank-Type"
	HXSeaType     = "X-Sea-Type"
	HXCalcContext = "X-Calc-Context"
	HXSurveyID    = "X-Survey-ID"
	HXDraftIndex  = "X-Draft-Index"
	HXDraftBlock  = "X-Draft-Block"
	HXTankID      = "X-Tank-ID"
	HXTankItem    = "X-Tank-Item"
)

// Tanks values (values)
const (
	FWTank = "fwt"
	BWTank = "bwt"
)

// Context values (values)
const (
	UpdtTanksWeight    = "updt_tanks_weight"
	UpdtTankVolume     = "updt_tank_volume"
	UpdtDraftCalcPanel = "updt-draft-calc-panel"
)

// Draft Block names (values)
const (
	SeaConditionBlock = "sea-condition-block"
	MarksBlock        = "marks-block"
	PPKeelBlock       = "ppkeel-block"
	HydrostaticsBlock = "hydrostatics-block"
	DeductiblesBlock  = "deductibles-block"
)
