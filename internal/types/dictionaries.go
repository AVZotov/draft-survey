package types

type Country struct {
	CountryCode string `json:"code"`
	Name        string `json:"name"`
}

type Port struct {
	Locode      string `json:"locode"`
	Name        string `json:"name"`
	CountryCode string `json:"country_code"`
	Coordinates string `json:"coordinates"`
}
