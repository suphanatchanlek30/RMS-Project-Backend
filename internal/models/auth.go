package models

type RoleResponse struct {
	RoleID   int    `json:"roleId"`
	RoleName string `json:"roleName"`
}

type EmployeeAuth struct {
	EmployeeID     int    `json:"employeeId"`
	EmployeeName   string `json:"employeeName"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phoneNumber,omitempty"`
	EmployeeStatus bool   `json:"employeeStatus"`
	PasswordHash   string `json:"-"`
	RoleID         int    `json:"-"`
	RoleName       string `json:"-"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginEmployeeResponse struct {
	EmployeeID   int          `json:"employeeId"`
	EmployeeName string       `json:"employeeName"`
	Email        string       `json:"email"`
	Role         RoleResponse `json:"role"`
}

type LoginResponseData struct {
	AccessToken string                `json:"accessToken"`
	TokenType   string                `json:"tokenType"`
	ExpiresIn   int                   `json:"expiresIn"`
	Employee    LoginEmployeeResponse `json:"employee"`
}

type MeResponseData struct {
	EmployeeID     int          `json:"employeeId"`
	EmployeeName   string       `json:"employeeName"`
	Email          string       `json:"email"`
	PhoneNumber    string       `json:"phoneNumber"`
	EmployeeStatus bool         `json:"employeeStatus"`
	Role           RoleResponse `json:"role"`
}
