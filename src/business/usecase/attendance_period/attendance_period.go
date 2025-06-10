package attendance_period

import (
	"context"
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/null"
	"github.com/reyhanmichiels/go-pkg/v2/query"
	"github.com/reyhanmichies/employee-payroll-service/src/business/domain/attendance_period"
	"github.com/reyhanmichies/employee-payroll-service/src/business/dto"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
)

var Now = time.Now

type Interface interface {
	Create(ctx context.Context, inputParam dto.CreateAttendancePeriodParam) (entity.AttendancePeriod, error)
	GetCurrentAttendancePeriod(ctx context.Context) (entity.AttendancePeriod, error)
}

type attendancePeriod struct {
	auth                auth.Interface
	attendancePeriodDom attendance_period.Interface
}

type InitParam struct {
	Auth             auth.Interface
	AttendancePeriod attendance_period.Interface
}

func Init(param InitParam) Interface {
	return &attendancePeriod{
		auth:                param.Auth,
		attendancePeriodDom: param.AttendancePeriod,
	}
}

func (a *attendancePeriod) Create(
	ctx context.Context,
	inputParam dto.CreateAttendancePeriodParam,
) (
	entity.AttendancePeriod,
	error,
) {
	loginUser, err := a.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return entity.AttendancePeriod{}, err
	}

	if err := inputParam.Validate(); err != nil {
		return entity.AttendancePeriod{}, err
	}

	currentTime := null.TimeFrom(Now())

	attendancePeriod, err := a.attendancePeriodDom.Create(
		ctx,
		inputParam.ToAttendancePeriodInputParam(
			currentTime,
			loginUser.ID,
		),
	)
	if err != nil {
		switch errors.GetCode(err) {
		case codes.CodeSQLUniqueConstraint:
			return entity.AttendancePeriod{}, errors.NewWithCode(codes.CodeConflict, "attendance period already exists for the given date range")
		default:
			return entity.AttendancePeriod{}, err
		}
	}

	return attendancePeriod, nil
}

func (a *attendancePeriod) GetCurrentAttendancePeriod(ctx context.Context) (entity.AttendancePeriod, error) {
	attendancePeriod, err := a.attendancePeriodDom.Get(
		ctx,
		entity.AttendancePeriodParam{
			PeriodStatus: entity.PeriodStatusOpen,
			QueryOption: query.Option{
				IsActive: true,
			},
		},
	)
	if err != nil {
		switch errors.GetCode(err) {
		case codes.CodeSQLRecordDoesNotExist:
			return entity.AttendancePeriod{}, errors.NewWithCode(codes.CodeNotFound, "no open attendance period found")
		default:
			return entity.AttendancePeriod{}, err
		}
	}

	return attendancePeriod, nil
}
