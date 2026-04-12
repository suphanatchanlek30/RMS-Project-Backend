package models

type CreateEmployeeRequest struct {
	EmployeeName string `json:"employeeName"`
	RoleID       int    `json:"roleId"`
	PhoneNumber  string `json:"phoneNumber"`
	Email        string `json:"email"`
	HireDate     string `json:"hireDate"`
	Password     string `json:"password"`
}

type Employee struct {
	EmployeeID     int    `json:"employeeId"`
	EmployeeName   string `json:"employeeName"`
	RoleID         int    `json:"roleId"`
	RoleName       string `json:"roleName"`
	PhoneNumber    string `json:"phoneNumber"`
	Email          string `json:"email"`
	HireDate       string `json:"hireDate"`
	EmployeeStatus bool   `json:"employeeStatus"`
}

type UpdateEmployeeRequest struct {
	EmployeeName string `json:"employeeName"`
	PhoneNumber  string `json:"phoneNumber"`
	RoleID       int    `json:"roleId"`
}

type UpdateEmployeeStatusRequest struct {
	EmployeeStatus bool `json:"employeeStatus"`
}