package attendance

import (
	"context"
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/null"
	"github.com/reyhanmichiels/go-pkg/v2/query"
	attendance_dom "github.com/reyhanmichies/employee-payroll-service/src/business/domain/attendance"
	"github.com/reyhanmichies/employee-payroll-service/src/business/domain/attendance_period"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
)

var Now = time.Now

type Interface interface {
	Create(ctx context.Context) error
}

type attendance struct {
	attendancePeriodDom attendance_period.Interface
	attendanceDom       attendance_dom.Interface
	auth                auth.Interface
}

type InitParam struct {
	AttendancePeriod attendance_period.Interface
	Attendance       attendance_dom.Interface
	Auth             auth.Interface
}

func Init(param InitParam) Interface {
	return &attendance{
		attendancePeriodDom: param.AttendancePeriod,
		attendanceDom:       param.Attendance,
		auth:                param.Auth,
	}
}

func (a *attendance) Create(ctx context.Context) error {
	loginUser, err := a.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return err
	}

	// check if current time is not weekday
	currentTime := null.TimeFrom(Now())
	if currentTime.Time.Weekday() == time.Saturday || currentTime.Time.Weekday() == time.Sunday {
		return errors.NewWithCode(codes.CodeBadRequest, "attendance cannot be submitted on weekend")
	}

	// find the current attendance period
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
			return errors.NewWithCode(codes.CodeNotFound, "no open attendance period found")
		default:
			return err
		}
	}

	// submit attendance
	_, err = a.attendanceDom.Create(
		ctx,
		entity.AttendanceInputParam{
			AttendancePeriodID: attendancePeriod.ID,
			UserID:             loginUser.ID,
			AttendanceDate:     null.DateFrom(currentTime.Time),
			CreatedAt:          currentTime,
			CreatedBy:          null.Int64From(loginUser.ID),
		},
	)
	if err != nil {
		switch errors.GetCode(err) {
		case codes.CodeSQLUniqueConstraint:
			return errors.NewWithCode(codes.CodeConflict, "attendance already submitted for today")
		default:
			return err
		}
	}

	return nil
}
