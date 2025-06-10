package entity

import (
	"github.com/reyhanmichiels/go-pkg/v2/null"
	"github.com/reyhanmichiels/go-pkg/v2/query"
)

// PeriodStatus constants represent the possible statuses of an attendance period.
const (
	// PeriodStatusUpcoming indicates that the attendance period is scheduled to start in the future.
	PeriodStatusUpcoming = "UPCOMING"

	// PeriodStatusOpen indicates that the attendance period is currently active and open.
	PeriodStatusOpen = "OPEN"

	// PeriodStatusClosed indicates that the attendance period has ended and is no longer active.
	PeriodStatusClosed = "CLOSED"

	// PeriodStatusProcessing indicates that the attendance period is being calculated to create a payslip.
	PeriodStatusProcessing = "PROCESSING"

	// PeriodStatusProcessed indicates that the attendance period has been fully processed.
	PeriodStatusProcessed = "PROCESSED"
)

type AttendancePeriod struct {
	ID           int64     `db:"id" json:"id"`
	StartDate    null.Date `db:"start_date" json:"startDate"`
	EndDate      null.Date `db:"end_date" json:"endDate"`
	PeriodStatus string    `db:"period_status" json:"periodStatus"`

	// Utility Column
	Status    int64       `db:"status" json:"status"`
	Flag      int64       `db:"flag" json:"flag,omitempty"`
	Meta      null.String `db:"meta" json:"meta,omitempty" swaggertype:"string"`
	CreatedAt null.Time   `db:"created_at" json:"createdAt" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	CreatedBy null.Int64  `db:"created_by" json:"createdBy" swaggertype:"integer"`
	UpdatedAt null.Time   `db:"updated_at" json:"updatedAt" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	UpdatedBy null.Int64  `db:"updated_by" json:"updatedBy" swaggertype:"integer"`
	DeletedAt null.Time   `db:"deleted_at" json:"deletedAt,omitempty" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	DeletedBy null.Int64  `db:"deleted_by" json:"deletedBy,omitempty" swaggertype:"integer"`
}

type AttendancePeriodInputParam struct {
	StartDate    null.Date  `db:"start_date" json:"startDate"`
	EndDate      null.Date  `db:"end_date" json:"endDate"`
	PeriodStatus string     `db:"period_status" json:"periodStatus"`
	CreatedAt    null.Time  `db:"created_at" json:"-"`
	CreatedBy    null.Int64 `db:"created_by" json:"-"`
}

type AttendancePeriodUpdateParam struct {
	StartDate    null.Date  `db:"start_date" json:"startDate"`
	EndDate      null.Date  `db:"end_date" json:"endDate"`
	PeriodStatus string     `db:"period_status" json:"periodStatus"`
	Status       null.Int64 `db:"status" json:"status"`
	UpdatedAt    null.Time  `db:"updated_at" json:"-"`
	UpdatedBy    null.Int64 `db:"updated_by" json:"-"`
}

type AttendancePeriodParam struct {
	ID           int64  `db:"id" param:"id" json:"id"`
	PeriodStatus string `db:"period_status" param:"period_status" json:"periodStatus"`
	QueryOption  query.Option
	BypassCache  bool
	PaginationParam
}
