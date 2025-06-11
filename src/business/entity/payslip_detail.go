package entity

import (
	"github.com/reyhanmichiels/go-pkg/v2/null"
	"github.com/reyhanmichiels/go-pkg/v2/query"
)

// PayslipItemType enum constants
const (
	PayslipItemTypeEarningBasePay  = "EARNING_BASE_PAY"
	PayslipItemTypeEarningOvertime = "EARNING_OVERTIME"
	PayslipItemTypeReimbursement   = "REIMBURSEMENT"
)

type PayslipDetail struct {
	ID          int64   `db:"id" json:"id"`
	PayslipID   int64   `db:"fk_payslip_id" json:"payslipID"`
	ItemType    string  `db:"item_type" json:"itemType"`
	Description string  `db:"description" json:"description"`
	Amount      float64 `db:"amount" json:"amount"`

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

type PayslipDetailInputParam struct {
	PayslipID   int64        `db:"fk_payslip_id" json:"payslipID"`
	ItemType    string       `db:"item_type" json:"itemType"`
	Description string       `db:"description" json:"description"`
	Amount      null.Float64 `db:"amount" json:"amount"`
	CreatedAt   null.Time    `db:"created_at" json:"-"`
	CreatedBy   null.Int64   `db:"created_by" json:"-"`
}

type PayslipDetailUpdateParam struct {
	PayslipID   null.Int64   `db:"fk_payslip_id" json:"payslipID"`
	ItemType    null.String  `db:"item_type" json:"itemType"`
	Description null.String  `db:"description" json:"description"`
	Amount      null.Float64 `db:"amount" json:"amount"`
	Status      null.Int64   `db:"status" json:"status"`
	UpdatedAt   null.Time    `db:"updated_at" json:"-"`
	UpdatedBy   null.Int64   `db:"updated_by" json:"-"`
}

type PayslipDetailParam struct {
	ID          int64 `db:"id" param:"id" json:"id"`
	QueryOption query.Option
	BypassCache bool
	PaginationParam
}
