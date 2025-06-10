package entity

import (
	"github.com/reyhanmichiels/go-pkg/v2/null"
	"github.com/reyhanmichiels/go-pkg/v2/query"
)

type Reimbursement struct {
	ID                int64      `db:"id" json:"id"`
	UserID            int64      `db:"fk_user_id" json:"userID"`
	Description       string     `db:"description" json:"description"`
	Amount            float64    `db:"amount" json:"amount"`
	ReimbursementDate null.Date  `db:"reimbursement_date" json:"reimbursementDate" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	ApprovedDate      null.Date  `db:"approved_date" json:"approvedDate"`
	ApprovedBy        null.Int64 `db:"approved_by" json:"approvedBy"`

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

type ReimbursementInputParam struct {
	UserID            int64      `db:"fk_user_id" json:"userID"`
	Description       string     `db:"description" json:"description"`
	Amount            float64    `db:"amount" json:"amount"`
	ReimbursementDate null.Date  `db:"reimbursement_date" json:"reimbursementDate" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	ApprovedDate      null.Date  `db:"approved_date" json:"approvedDate"`
	ApprovedBy        null.Int64 `db:"approved_by" json:"approvedBy"`
	CreatedAt         null.Time  `db:"created_at" json:"-"`
	CreatedBy         null.Int64 `db:"created_by" json:"-"`
}

type ReimbursementUpdateParam struct {
	Description  string     `db:"description" json:"description"`
	Amount       float64    `db:"amount" json:"amount"`
	ApprovedDate null.Date  `db:"approved_date" json:"approvedDate"`
	ApprovedBy   null.Int64 `db:"approved_by" json:"approvedBy"`
	Status       null.Int64 `db:"status" json:"status"`
	UpdatedAt    null.Time  `db:"updated_at" json:"-"`
	UpdatedBy    null.Int64 `db:"updated_by" json:"-"`
}

type ReimbursementParam struct {
	ID          int64 `db:"id" param:"id" json:"id"`
	QueryOption query.Option
	BypassCache bool
	PaginationParam
}

func (r *ReimbursementInputParam) MockApprovalData(currentTime null.Time) {
	r.ApprovedDate = null.DateFrom(currentTime.Time)
	r.ApprovedBy = null.Int64From(1)
}
