package types

type User struct {
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
	Position   string `json:"position"`
	Email      string `json:"email" validate:"omitempty,email"`
	Company    string `json:"company"`
	License    string `json:"license_no"`
	Country    string `json:"country"`
	EmployeeID string `json:"employee_id"`
}
