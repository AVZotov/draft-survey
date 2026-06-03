package types

type User struct {
	FirstName   string `json:"first_name" schema:"first_name" validate:"required"`
	LastName    string `json:"last_name" schema:"last_name" validate:"required"`
	Position    string `json:"position" schema:"position"`
	Email       string `json:"email" schema:"email" validate:"omitempty,email"`
	Company     string `json:"company" schema:"company"`
	License     string `json:"license_no" schema:"license_no"`
	CountryCode string `json:"country_code" schema:"country_code"`
	EmployeeID  string `json:"employee_id" schema:"employee_id"`
}
