package dto

import (
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/null"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
)

type CreateOvertimeParam struct {
	OvertimeDate null.Date    `json:"overtimeDate" swaggertype:"string" example:"2025-06-09T00:00:00Z"`
	OvertimeHour null.Float64 `json:"overtimeHour" swaggertype:"number" example:"1.5"`
}

func (c *CreateOvertimeParam) Validate(currentTime time.Time) error {
	if !c.OvertimeDate.Valid {
		return errors.NewWithCode(codes.CodeBadRequest, "overtimeDate is required")
	}

	if !c.OvertimeHour.Valid || c.OvertimeHour.Float64 <= 0 {
		return errors.NewWithCode(codes.CodeBadRequest, "overtimeHour is required and must be a positive number")
	}

	if c.OvertimeHour.Float64 > 3 {
		return errors.NewWithCode(codes.CodeBadRequest, "overtimeHour cannot exceed 3 hours")
	}

	if c.OvertimeDate.Time.After(currentTime) {
		return errors.NewWithCode(codes.CodeBadRequest, "overtimeDate cannot be in the future")
	}

	if c.OvertimeDate.Time.Format(time.DateOnly) == currentTime.Format(time.DateOnly) && currentTime.Hour() < 17 {
		return errors.NewWithCode(codes.CodeBadRequest, "overtimeDate is today, but the current time must be after 5 PM")
	}

	return nil
}

func (c *CreateOvertimeParam) ToOvertimeInputParam(currentTime null.Time, userID int64) entity.OvertimeInputParam {
	return entity.OvertimeInputParam{
		UserID:       userID,
		OvertimeDate: c.OvertimeDate,
		OvertimeHour: c.OvertimeHour,
		CreatedAt:    currentTime,
		CreatedBy:    null.Int64From(userID),
	}
}
