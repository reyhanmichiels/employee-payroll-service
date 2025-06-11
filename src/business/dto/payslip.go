package dto

import "github.com/reyhanmichiels/go-pkg/v2/null"

type Payslip struct {
	StartDate          null.Date       `json:"startDate" example:"2025-06-01"`
	EndDate            null.Date       `json:"endDate" example:"2025-06-30"`
	BasePayComponent   float64         `json:"basePayComponent" example:"5000000.00"`
	OvertimeComponent  float64         `json:"overtimeComponent" example:"200000.00"`
	ReimburseComponent float64         `json:"reimburseComponent" example:"150000.00"`
	TotalTakeHomePay   float64         `json:"totalTakeHomePay" example:"5350000.00"`
	Details            []PayslipDetail `json:"details"`
}

type PayslipDetail struct {
	Type        string  `json:"type" example:"EARNING_OVERTIME"`
	Description string  `json:"description" example:"Overtime: 2 hours on 2025-06-15"`
	Amount      float64 `json:"amount" example:"200000.00"`
}
