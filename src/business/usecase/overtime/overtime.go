package overtime

import (
	"context"
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/null"
	overtimeDom "github.com/reyhanmichies/employee-payroll-service/src/business/domain/overtime"
	"github.com/reyhanmichies/employee-payroll-service/src/business/dto"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
)

var Now = time.Now

type Interface interface {
	Create(ctx context.Context, inputParam dto.CreateOvertimeParam) (entity.Overtime, error)
}

type overtime struct {
	auth        auth.Interface
	overtimeDom overtimeDom.Interface
}

type InitParam struct {
	Auth        auth.Interface
	OvertimeDom overtimeDom.Interface
}

func Init(param InitParam) Interface {
	return &overtime{
		auth:        param.Auth,
		overtimeDom: param.OvertimeDom,
	}
}

func (o *overtime) Create(
	ctx context.Context,
	inputParam dto.CreateOvertimeParam,
) (entity.Overtime, error) {
	loginUser, err := o.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return entity.Overtime{}, err
	}

	currentTime := null.TimeFrom(Now())

	if err := inputParam.Validate(currentTime.Time); err != nil {
		return entity.Overtime{}, err
	}

	overtimeInputParam := inputParam.ToOvertimeInputParam(currentTime, loginUser.ID)
	overtimeInputParam.MockApprovalData(currentTime) // TODO: remove this line on production, this is only used for testing purpose
	overtime, err := o.overtimeDom.Create(ctx, overtimeInputParam)
	if err != nil {
		switch errors.GetCode(err) {
		case codes.CodeSQLUniqueConstraint:
			return entity.Overtime{}, errors.NewWithCode(codes.CodeConflict, "overtime already submitted for this date")
		default:
			return entity.Overtime{}, err
		}
	}

	return overtime, nil
}
