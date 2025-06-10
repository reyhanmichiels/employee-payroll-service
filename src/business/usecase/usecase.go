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
	"github.com/reyhanmichies/employee-payroll-service/src/business/usecase/user"
)

type Usecases struct {
	User             user.Interface
	AttendancePeriod attendance_period.Interface
	Attendance       attendance.Interface
	Overtime         overtime.Interface
}

type InitParam struct {
	Dom  *domain.Domains
	Json parser.JSONInterface
	Log  log.Interface
	Hash hash.Interface
	Auth auth.Interface
}

func Init(param InitParam) *Usecases {
	return &Usecases{
		User:             user.Init(user.InitParam{UserDomain: param.Dom.User, Auth: param.Auth, Hash: param.Hash}),
		AttendancePeriod: attendance_period.Init(attendance_period.InitParam{Auth: param.Auth, AttendancePeriod: param.Dom.AttendancePeriod}),
		Attendance:       attendance.Init(attendance.InitParam{Auth: param.Auth, AttendancePeriod: param.Dom.AttendancePeriod, Attendance: param.Dom.Attendance}),
		Overtime:         overtime.Init(overtime.InitParam{Auth: param.Auth, OvertimeDom: param.Dom.Overtime}),
	}
}
