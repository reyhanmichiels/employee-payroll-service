package attendance_period

import (
	"context"
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/log"
	"github.com/reyhanmichiels/go-pkg/v2/null"
	"github.com/reyhanmichiels/go-pkg/v2/parser"
	"github.com/reyhanmichiels/go-pkg/v2/query"
	"github.com/reyhanmichiels/go-pkg/v2/sql"
	"github.com/reyhanmichies/employee-payroll-service/src/business/domain/attendance"
	"github.com/reyhanmichies/employee-payroll-service/src/business/domain/attendance_period"
	"github.com/reyhanmichies/employee-payroll-service/src/business/domain/overtime"
	"github.com/reyhanmichies/employee-payroll-service/src/business/domain/payslip"
	"github.com/reyhanmichies/employee-payroll-service/src/business/domain/payslip_detail"
	"github.com/reyhanmichies/employee-payroll-service/src/business/domain/reimbursement"
	"github.com/reyhanmichies/employee-payroll-service/src/business/domain/transactor"
	"github.com/reyhanmichies/employee-payroll-service/src/business/domain/user"
	"github.com/reyhanmichies/employee-payroll-service/src/business/dto"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
	"github.com/reyhanmichies/employee-payroll-service/src/handler/pubsub/publisher"
)

var Now = time.Now

type Interface interface {
	Create(ctx context.Context, inputParam dto.CreateAttendancePeriodParam) (entity.AttendancePeriod, error)
	GetCurrentAttendancePeriod(ctx context.Context) (entity.AttendancePeriod, error)
	GeneratePayroll(ctx context.Context, attendancePeriodID int64) error

	PubSubGeneratePayroll(ctx context.Context, message entity.PubSubMessage) error
}

type attendancePeriod struct {
	auth                auth.Interface
	attendancePeriodDom attendance_period.Interface
	publisher           publisher.Interface
	transactor          transactor.Interface
	json                parser.JSONInterface
	log                 log.Interface
	payslipDom          payslip.Interface
	payslipDetailDom    payslip_detail.Interface
	userDom             user.Interface
	overtimeDom         overtime.Interface
	reimbursementDom    reimbursement.Interface
	attendanceDom       attendance.Interface
}

type InitParam struct {
	Auth             auth.Interface
	AttendancePeriod attendance_period.Interface
	Publisher        publisher.Interface
	Transactor       transactor.Interface
	Json             parser.JSONInterface
	Log              log.Interface
	Payslip          payslip.Interface
	PayslipDetail    payslip_detail.Interface
	User             user.Interface
	Overtime         overtime.Interface
	Reimbursement    reimbursement.Interface
	Attendance       attendance.Interface
}

func Init(param InitParam) Interface {
	return &attendancePeriod{
		auth:                param.Auth,
		attendancePeriodDom: param.AttendancePeriod,
		publisher:           param.Publisher,
		transactor:          param.Transactor,
		json:                param.Json,
		log:                 param.Log,
		payslipDom:          param.Payslip,
		payslipDetailDom:    param.PayslipDetail,
		userDom:             param.User,
		overtimeDom:         param.Overtime,
		reimbursementDom:    param.Reimbursement,
		attendanceDom:       param.Attendance,
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

func (a *attendancePeriod) GeneratePayroll(
	ctx context.Context,
	attendancePeriodID int64,
) error {
	loginUser, err := a.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return err
	}

	attendancePeriodParam := entity.AttendancePeriodParam{
		ID: attendancePeriodID,
		QueryOption: query.Option{
			IsActive: true,
		},
	}

	attendancePeriod, err := a.attendancePeriodDom.Get(
		ctx,
		attendancePeriodParam,
	)
	if err != nil {
		switch errors.GetCode(err) {
		case codes.CodeSQLRecordDoesNotExist:
			return errors.NewWithCode(codes.CodeNotFound, "attendance period not found")
		default:
			return err
		}
	}

	switch attendancePeriod.PeriodStatus {
	case entity.PeriodStatusProcessed:
		return errors.NewWithCode(codes.CodeBadRequest, "attendance period has already been processed")
	case entity.PeriodStatusProcessing:
		return errors.NewWithCode(codes.CodeConflict, "attendance period is currently being processed")
	case entity.PeriodStatusOpen:
		return errors.NewWithCode(codes.CodeBadRequest, "attendance period is still open and cannot be processed")
	case entity.PeriodStatusUpcoming:
		return errors.NewWithCode(codes.CodeBadRequest, "attendance period is upcoming and cannot be processed")
	}

	return a.transactor.Execute(ctx, "txGeneratePayroll", sql.TxOptions{}, func(ctx context.Context) error {
		err = a.attendancePeriodDom.Update(
			ctx,
			entity.AttendancePeriodUpdateParam{
				PeriodStatus: entity.PeriodStatusProcessed,
				UpdatedAt:    null.TimeFrom(Now()),
				UpdatedBy:    null.Int64From(loginUser.ID),
			},
			attendancePeriodParam,
		)
		if err != nil {
			return err
		}

		err = a.publisher.Publish(
			ctx,
			entity.ExchangePayrollEvent,
			entity.RoutingKeyPayrollCalculate,
			dto.PubSubGeneratePayrollMessage{
				AttendancePeriod: attendancePeriod,
				LoginUser:        loginUser,
			},
		)
		if err != nil {
			return err
		}

		return nil
	})
}
