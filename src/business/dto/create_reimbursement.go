package dto

import (
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/null"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
)

type CreateReimbursementParam struct {
	Description       string    `json:"description" example:"Reimbursement for office supplies"`
	Amount            float64   `json:"amount" example:"150000.00"`
	ReimbursementDate null.Date `json:"reimbursementDate" swaggertype:"string" example:"2025-06-19T00:00:00Z"`
}

func (c *CreateReimbursementParam) Validate(currentTime time.Time) error {
	if c.Description == "" {
		return errors.NewWithCode(codes.CodeBadRequest, "description is required")
	}

	if c.Amount <= 0 {
		return errors.NewWithCode(codes.CodeBadRequest, "amount must be greater than zero")
	}

	if !c.ReimbursementDate.Valid {
		return errors.NewWithCode(codes.CodeBadRequest, "reimbursementDate is required")
	}

	if currentTime.Before(c.ReimbursementDate.Time) {
		return errors.NewWithCode(codes.CodeBadRequest, "reimbursementDate cannot be in the future")
	}

	return nil
}

func (c *CreateReimbursementParam) ToReimbursementInputParam(currentTime null.Time, userID int64) entity.ReimbursementInputParam {
	return entity.ReimbursementInputParam{
		UserID:            userID,
		Description:       c.Description,
		Amount:            c.Amount,
		ReimbursementDate: c.ReimbursementDate,
		CreatedAt:         currentTime,
		CreatedBy:         null.Int64From(userID),
	}
}
