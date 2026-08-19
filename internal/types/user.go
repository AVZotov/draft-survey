package types

type User struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Position    string `json:"position"`
	Email       string `json:"email" validate:"omitempty,email"`
	Company     string `json:"company"`
	License     string `json:"license_no"`
	CountryCode string `json:"country_code"`
	EmployeeID  string `json:"employee_id"`
}
