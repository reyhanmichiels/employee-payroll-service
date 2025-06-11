package usecase

import (
	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichiels/go-pkg/v2/hash"
	"github.com/reyhanmichiels/go-pkg/v2/log"
	"github.com/reyhanmichiels/go-pkg/v2/parser"
	"github.com/reyhanmichies/employee-payroll-service/src/business/domain"
	"github.com/reyhanmichies/employee-payroll-service/src/business/usecase/attendance"
	"github.com/reyhanmichies/employee-payroll-service/src/business/usecase/attendance_period"
	"github.com/reyhanmichies/employee-payroll-service/src/business/usecase/overtime"
	"github.com/reyhanmichies/employee-payroll-service/src/business/usecase/reimbursement"
	"github.com/reyhanmichies/employee-payroll-service/src/business/usecase/user"
	"github.com/reyhanmichies/employee-payroll-service/src/handler/pubsub/publisher"
)

type Usecases struct {
	User             user.Interface
	AttendancePeriod attendance_period.Interface
	Attendance       attendance.Interface
	Overtime         overtime.Interface
	Reimbursement    reimbursement.Interface
}

type InitParam struct {
	Dom       *domain.Domains
	Json      parser.JSONInterface
	Log       log.Interface
	Hash      hash.Interface
	Auth      auth.Interface
	Publisher publisher.Interface
}

func Init(param InitParam) *Usecases {
	return &Usecases{
		User:             user.Init(user.InitParam{UserDomain: param.Dom.User, Auth: param.Auth, Hash: param.Hash}),
		AttendancePeriod: attendance_period.Init(attendance_period.InitParam{Auth: param.Auth, AttendancePeriod: param.Dom.AttendancePeriod, Publisher: param.Publisher, Transactor: param.Dom.Transactor, Json: param.Json, Log: param.Log, Payslip: param.Dom.Payslip, PayslipDetail: param.Dom.PayslipDetail, User: param.Dom.User, Overtime: param.Dom.Overtime, Reimbursement: param.Dom.Reimbursement, Attendance: param.Dom.Attendance}),
		Attendance:       attendance.Init(attendance.InitParam{Auth: param.Auth, AttendancePeriod: param.Dom.AttendancePeriod, Attendance: param.Dom.Attendance}),
		Overtime:         overtime.Init(overtime.InitParam{Auth: param.Auth, OvertimeDom: param.Dom.Overtime}),
		Reimbursement:    reimbursement.Init(reimbursement.InitParam{Auth: param.Auth, Reimbursement: param.Dom.Reimbursement}),
	}
}
