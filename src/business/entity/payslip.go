package entity

import (
	"github.com/reyhanmichiels/go-pkg/v2/null"
	"github.com/reyhanmichiels/go-pkg/v2/query"
)

type Payslip struct {
	ID                     int64   `db:"id" json:"id"`
	UserID                 int64   `db:"fk_user_id" json:"userID"`
	AttendancePeriodID     int64   `db:"fk_attendance_period_id" json:"attendancePeriodID"`
	BasePayComponent       float64 `db:"base_pay_component" json:"basePayComponent"`
	OvertimeComponent      float64 `db:"overtime_component" json:"overtimeComponent"`
	ReimbursementComponent float64 `db:"reimbursement_component" json:"reimbursementComponent"`
	TotalTakeHomePay       float64 `db:"total_take_home_pay" json:"totalTakeHomePay"`

	// Utility Column
	Status    int64       `db:"status" json:"status"`
	Flag      int64       `db:"flag" json:"flag,omitempty"`
	Meta      null.String `db:"meta" json:"meta,omitempty" swaggertype:"string"`
	CreatedAt null.Time   `db:"created_at" json:"createdAt" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	CreatedBy null.Int64  `db:"created_by" json:"createdBy" swaggertype:"string"`
	UpdatedAt null.Time   `db:"updated_at" json:"updatedAt" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	UpdatedBy null.String `db:"updated_by" json:"updatedBy" swaggertype:"string"`
	DeletedAt null.Time   `db:"deleted_at" json:"deletedAt,omitempty" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	DeletedBy null.String `db:"deleted_by" json:"deletedBy,omitempty" swaggertype:"string"`
}

type PayslipInputParam struct {
	UserID                 int64        `db:"fk_user_id" json:"userID"`
	AttendancePeriodID     int64        `db:"fk_attendance_period_id" json:"attendancePeriodID"`
	BasePayComponent       null.Float64 `db:"base_pay_component" json:"basePayComponent"`
	OvertimeComponent      null.Float64 `db:"overtime_component" json:"overtimeComponent"`
	ReimbursementComponent null.Float64 `db:"reimbursement_component" json:"reimbursementComponent"`
	TotalTakeHomePay       null.Float64 `db:"total_take_home_pay" json:"totalTakeHomePay"`
	CreatedAt              null.Time    `db:"created_at" json:"-"`
	CreatedBy              null.Int64   `db:"created_by" json:"-"`
}

type PayslipUpdateParam struct {
	UserID                 null.Int64   `db:"fk_user_id" json:"userID"`
	AttendancePeriodID     null.Int64   `db:"fk_attendance_period_id" json:"attendancePeriodID"`
	BasePayComponent       null.Float64 `db:"base_pay_component" json:"basePayComponent"`
	OvertimeComponent      null.Float64 `db:"overtime_component" json:"overtimeComponent"`
	ReimbursementComponent null.Float64 `db:"reimbursement_component" json:"reimbursementComponent"`
	TotalTakeHomePay       null.Float64 `db:"total_take_home_pay" json:"totalTakeHomePay"`
	Status                 null.Int64   `db:"status" json:"status"`
	UpdatedAt              null.Time    `db:"updated_at" json:"-"`
	UpdatedBy              null.Int64   `db:"updated_by" json:"-"`
}

type PayslipParam struct {
	ID                 int64 `db:"id" param:"id" json:"id"`
	AttendancePeriodID int64 `db:"fk_attendance_period_id" param:"attendance_period_id"`
	UserID             int64 `db:"fk_user_id" param:"user_id" `
	QueryOption        query.Option
	BypassCache        bool
	PaginationParam
}
