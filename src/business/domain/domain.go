package domain

import (
	"github.com/reyhanmichiels/go-pkg/v2/log"
	"github.com/reyhanmichiels/go-pkg/v2/parser"
	"github.com/reyhanmichiels/go-pkg/v2/redis"
	"github.com/reyhanmichiels/go-pkg/v2/sql"
	"github.com/reyhanmichies/employee-payroll-service/src/business/domain/attendance"
	"github.com/reyhanmichies/employee-payroll-service/src/business/domain/attendance_period"
	"github.com/reyhanmichies/employee-payroll-service/src/business/domain/overtime"
	"github.com/reyhanmichies/employee-payroll-service/src/business/domain/payslip"
	"github.com/reyhanmichies/employee-payroll-service/src/business/domain/payslip_detail"
	"github.com/reyhanmichies/employee-payroll-service/src/business/domain/reimbursement"
	"github.com/reyhanmichies/employee-payroll-service/src/business/domain/transactor"
	"github.com/reyhanmichies/employee-payroll-service/src/business/domain/user"
)

type Domains struct {
	User             user.Interface
	Attendance       attendance.Interface
	AttendancePeriod attendance_period.Interface
	Overtime         overtime.Interface
	Reimbursement    reimbursement.Interface
	Transactor       transactor.Interface
	Payslip          payslip.Interface
	PayslipDetail    payslip_detail.Interface
}

type InitParam struct {
	Log   log.Interface
	Db    sql.Interface
	Redis redis.Interface
	Json  parser.JSONInterface
	// TODO: add audit
}

func Init(param InitParam) *Domains {
	return &Domains{
		User:             user.Init(user.InitParam{Db: param.Db, Log: param.Log, Redis: param.Redis, Json: param.Json}),
		Attendance:       attendance.Init(attendance.InitParam{Db: param.Db, Log: param.Log, Redis: param.Redis, Json: param.Json}),
		AttendancePeriod: attendance_period.Init(attendance_period.InitParam{Db: param.Db, Log: param.Log, Redis: param.Redis, Json: param.Json}),
		Overtime:         overtime.Init(overtime.InitParam{Db: param.Db, Log: param.Log, Redis: param.Redis, Json: param.Json}),
		Reimbursement:    reimbursement.Init(reimbursement.InitParam{Db: param.Db, Log: param.Log, Redis: param.Redis, Json: param.Json}),
		Transactor:       transactor.Init(param.Db),
		Payslip:          payslip.Init(payslip.InitParam{Db: param.Db, Log: param.Log, Redis: param.Redis, Json: param.Json}),
		PayslipDetail:    payslip_detail.Init(payslip_detail.InitParam{Db: param.Db, Log: param.Log, Redis: param.Redis, Json: param.Json}),
	}
}
