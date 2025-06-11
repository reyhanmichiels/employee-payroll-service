package dto

import (
	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
)

type PubSubGeneratePayrollMessage struct {
	AttendancePeriod entity.AttendancePeriod `json:"attendancePeriod"`
	LoginUser        auth.User               `json:"loginUser"`
}
