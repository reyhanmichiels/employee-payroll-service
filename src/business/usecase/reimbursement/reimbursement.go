package reimbursement

import (
	"context"
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichiels/go-pkg/v2/null"
	reimbursementDom "github.com/reyhanmichies/employee-payroll-service/src/business/domain/reimbursement"
	"github.com/reyhanmichies/employee-payroll-service/src/business/dto"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
)

var Now = time.Now

type Interface interface {
	Create(ctx context.Context, inputParam dto.CreateReimbursementParam) (entity.Reimbursement, error)
}

type reimbursement struct {
	auth             auth.Interface
	reimbursementDom reimbursementDom.Interface
}

type InitParam struct {
	Auth          auth.Interface
	Reimbursement reimbursementDom.Interface
}

func Init(param InitParam) Interface {
	return &reimbursement{
		auth:             param.Auth,
		reimbursementDom: param.Reimbursement,
	}
}

func (r *reimbursement) Create(
	ctx context.Context,
	inputParam dto.CreateReimbursementParam,
) (entity.Reimbursement, error) {
	loginUser, err := r.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return entity.Reimbursement{}, err
	}

	currentTime := null.TimeFrom(Now())

	if err := inputParam.Validate(currentTime.Time); err != nil {
		return entity.Reimbursement{}, err
	}

	reimbursementInputParam := inputParam.ToReimbursementInputParam(currentTime, loginUser.ID)
	reimbursementInputParam.MockApprovalData(currentTime) // TODO: remove this line on production, this is only used for testing purpose
	reimbursement, err := r.reimbursementDom.Create(ctx, reimbursementInputParam)
	if err != nil {
		return entity.Reimbursement{}, err
	}

	return reimbursement, nil
}
