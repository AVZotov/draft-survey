package constants

const (
	// Marks
	FwdPort         = "fwd_port"
	FwdPortMarkRead = "fwd_port_mr"
	MidPort         = "mid_port"
	MidPortMarkRead = "mid_port_mr"
	AftPort         = "aft_port"
	AftPortMarkRead = "aft_port_mr"
	FwdStbd         = "fwd_stbd"
	FwdStbdMarkRead = "fwd_stbd_mr"
	MidStbd         = "mid_stbd"
	MidStbdMarkRead = "mid_stbd_mr"
	AftStbd         = "aft_stbd"
	AftStbdMarkRead = "aft_stbd_mr"

	// PP distances
	DFwd    = "dFwd"
	DFwdDir = "dFwdDir"
	DMid    = "dMid"
	DMidDir = "dMidDir"
	DAft    = "dAft"
	DAftDir = "dAftDir"

	// Keel
	KeelFwd = "Fk"
	KeelMid = "Mk"
	KeelAft = "Ak"

	// Hydrostatics
	UDraft  = "uDraft"
	UDisp   = "uDisp"
	UTpc    = "uTpc"
	ULcfLca = "uLcfLca"
	ULcfDir = "uLcfDir"
	LDraft  = "lDraft"
	LDisp   = "lDisp"
	LTpc    = "lTpc"
	LLcfLca = "lLcfLca"
	LLcfDir = "lLcfDir"

	// MTC
	PMtcDraft = "pMtcDraft"
	PMtc      = "pMtc"
	NMtcDraft = "nMtcDraft"
	NMtc      = "nMtc"

	// Deductibles
	HFO           = "hfo"
	MDO           = "mdo"
	LubOil        = "lub"
	BilgeWater    = "bilW"
	Others        = "others"
	ConstDeclared = "constDeclared"
	CargoDeclared = "cargoDeclared"

	// Density
	TableDensity     = "tDens"
	DockwaterDensity = "dwDens"

	// Sea condition
	SeaType          = "sea_type"
	SeaCondition     = "sea_condition"
	SeaConditionWave = "sea_wave"
	SeaConditionIce  = "sea_ice"
)
