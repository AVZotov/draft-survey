package types

type Country struct {
	CountryCode string `json:"code" schema:"code"`
	Name        string `json:"name" schema:"name"`
}

type Port struct {
	Locode      string `json:"locode"`
	Name        string `json:"name"`
	CountryCode string `json:"country_code"`
	Coordinates string `json:"coordinates"`
	PortManual  string `json:"port_manual"`
}
