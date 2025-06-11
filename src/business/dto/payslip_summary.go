package dto

type PayslipSummary struct {
	TotalEmployeeTakeHomePay float64          `json:"totalEmployeeTakeHomePay" example:"5000000.00"`
	TotalEmployee            int64            `json:"totalEmployee" example:"100"`
	EmployeePayouts          []EmployeePayout `json:"employeePayouts"`
}

type EmployeePayout struct {
	ID          int64   `json:"id" example:"1"`
	Name        string  `json:"name" example:"John Doe"`
	TakeHomePay float64 `json:"takeHomePay" example:"5000000.00"`
}
