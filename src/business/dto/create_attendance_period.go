package dto

import (
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/null"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
)

type CreateAttendancePeriodParam struct {
	StartDate null.Date `json:"startDate" swaggertype:"string" example:"2025-06-09T00:00:00Z"`
	EndDate   null.Date `json:"endDate" swaggertype:"string" example:"2025-06-09T00:00:00Z"`
}

func (c *CreateAttendancePeriodParam) Validate() error {
	if !c.StartDate.Valid {
		return errors.NewWithCode(codes.CodeBadRequest, "startDate is required")
	}

	if !c.EndDate.Valid {
		return errors.NewWithCode(codes.CodeBadRequest, "endDate is required")
	}

	if c.StartDate.Time.After(c.EndDate.Time) {
		return errors.NewWithCode(codes.CodeBadRequest, "startDate cannot be after endDate")
	}

	if c.StartDate.Time.Before(time.Now()) {
		return errors.NewWithCode(codes.CodeBadRequest, "startDate cannot be before current date")
	}

	return nil
}

func (c *CreateAttendancePeriodParam) ToAttendancePeriodInputParam(currentTime null.Time, userID int64) entity.AttendancePeriodInputParam {
	return entity.AttendancePeriodInputParam{
		StartDate:    c.StartDate,
		EndDate:      c.EndDate,
		PeriodStatus: entity.PeriodStatusUpcoming,
		CreatedAt:    currentTime,
		CreatedBy:    null.Int64From(userID),
	}
}
