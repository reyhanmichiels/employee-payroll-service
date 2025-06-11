package entity

import (
	"github.com/reyhanmichiels/go-pkg/v2/null"
	"github.com/reyhanmichiels/go-pkg/v2/query"
)

type Attendance struct {
	ID                 int64     `db:"id" json:"id"`
	AttendancePeriodID int64     `db:"fk_attendance_period_id" json:"attendancePeriodID"`
	UserID             int64     `db:"fk_user_id" json:"userID"`
	AttendanceDate     null.Date `db:"attendance_date" json:"attendanceDate"`

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

type AttendanceInputParam struct {
	AttendancePeriodID int64      `db:"fk_attendance_period_id" json:"attendancePeriodID"`
	UserID             int64      `db:"fk_user_id" json:"userID"`
	AttendanceDate     null.Date  `db:"attendance_date" json:"attendanceDate"`
	CreatedAt          null.Time  `db:"created_at" json:"-"`
	CreatedBy          null.Int64 `db:"created_by" json:"-"`
}

type AttendanceUpdateParam struct {
	AttendancePeriodID int64      `db:"fk_attendance_period_id" json:"attendancePeriodID"`
	AttendanceDate     null.Date  `db:"attendance_date" json:"attendanceDate"`
	Status             null.Int64 `db:"status" json:"status"`
	UpdatedAt          null.Time  `db:"updated_at" json:"-"`
	UpdatedBy          null.Int64 `db:"updated_by" json:"-"`
}

type AttendanceParam struct {
	ID          int64 `db:"id" param:"id" json:"id"`
	QueryOption query.Option
	BypassCache bool
	PaginationParam
}

type UserAttendanceCount map[int64]int64
