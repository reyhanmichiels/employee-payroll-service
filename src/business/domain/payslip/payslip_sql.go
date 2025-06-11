package payslip

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/query"
	"github.com/reyhanmichiels/go-pkg/v2/sql"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
)

func (p *payslip) getSQL(ctx context.Context, param entity.PayslipParam) (entity.Payslip, error) {
	payslip := entity.Payslip{}

	p.log.Debug(ctx, fmt.Sprintf("get payslip with body: %v", param))

	param.QueryOption.DisableLimit = true
	qb := query.NewSQLQueryBuilder(p.db, "param", "db", &param.QueryOption)
	queryExt, queryArgs, _, _, err := qb.Build(&param)
	if err != nil {
		return payslip, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	row, err := p.db.QueryRow(ctx, "rPayslip", readPayslip+queryExt, queryArgs...)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return payslip, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	if err := row.StructScan(&payslip); err != nil && errors.Is(err, sql.ErrNotFound) {
		return payslip, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, err.Error())
	} else if err != nil {
		return payslip, errors.NewWithCode(codes.CodeSQLRowScan, err.Error())
	}

	p.log.Debug(ctx, fmt.Sprintf("success get payslip with body: %v", param))

	return payslip, nil
}

func (p *payslip) getListSQL(ctx context.Context, param entity.PayslipParam) ([]entity.Payslip, *entity.Pagination, error) {
	payslipList := []entity.Payslip{}
	pg := entity.Pagination{}

	p.log.Debug(ctx, fmt.Sprintf("get payslip list with body: %v", param))

	qb := query.NewSQLQueryBuilder(p.db, "param", "db", &param.QueryOption)
	queryExt, queryArgs, countExt, countArgs, err := qb.Build(&param)
	if err != nil {
		return payslipList, &pg, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	rows, err := p.db.Query(ctx, "rPayslipList", readPayslip+queryExt, queryArgs...)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return payslipList, &pg, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		payslip := entity.Payslip{}
		err := rows.StructScan(&payslip)
		if err != nil {
			p.log.Error(ctx, errors.NewWithCode(codes.CodeSQLRowScan, err.Error()))
			continue
		}

		payslipList = append(payslipList, payslip)
	}

	pg = entity.Pagination{
		CurrentPage:     param.PaginationParam.Page,
		CurrentElements: int64(len(payslipList)),
		SortBy:          param.SortBy,
	}

	if !param.QueryOption.DisableLimit && len(payslipList) > 0 {
		err := p.db.Get(ctx, "cPayslipList", countPayslip+countExt, &pg.TotalElements, countArgs...)
		if err != nil {
			return payslipList, &pg, errors.NewWithCode(codes.CodeSQLRead, err.Error())
		}
	}

	pg.ProcessPagination(param.Limit)

	p.log.Debug(ctx, fmt.Sprintf("success get payslip list with body: %v", param))

	return payslipList, &pg, nil
}

func (p *payslip) createSQL(ctx context.Context, inputParam entity.PayslipInputParam) (entity.Payslip, error) {
	payslip := entity.Payslip{}

	p.log.Debug(ctx, fmt.Sprintf("create payslip with body: %v", inputParam))

	stmt, err := p.db.PrepareNamed(ctx, "iNewPayslip", insertPayslip)
	if err != nil {
		return payslip, errors.NewWithCode(codes.CodeSQLPrepareStmt, err.Error())
	}
	defer stmt.Close()

	err = stmt.Get(&payslip, inputParam)

	pgErr := &pq.Error{}
	if err != nil && errors.As(err, &pgErr) && pgErr.Code == entity.PSQLUniqueConstraintCode {
		return payslip, errors.NewWithCode(codes.CodeSQLUniqueConstraint, err.Error())
	} else if err != nil {
		return payslip, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	p.log.Debug(ctx, fmt.Sprintf("success create payslip with body: %v", inputParam))

	return payslip, nil
}

func (p *payslip) createManySQL(ctx context.Context, inputParams []entity.PayslipInputParam) error {
	p.log.Debug(ctx, fmt.Sprintf("create many payslip with body: %v", inputParams))

	res, err := p.db.NamedExec(ctx, "iManyPayslip", insertManyPayslip, inputParams)

	pgErr := &pq.Error{}
	if err != nil && errors.As(err, &pgErr) && pgErr.Code == entity.PSQLUniqueConstraintCode {
		return errors.NewWithCode(codes.CodeSQLUniqueConstraint, err.Error())
	} else if err != nil {
		return errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, err.Error())
	} else if rowCount < int64(len(inputParams)) {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no payslip created")
	}

	p.log.Debug(ctx, fmt.Sprintf("success create many payslip with body: %v", inputParams))

	return nil
}

func (p *payslip) updateSQL(ctx context.Context, updateParam entity.PayslipUpdateParam, selectParam entity.PayslipParam) error {
	p.log.Debug(ctx, fmt.Sprintf("update payslip with body: %v", updateParam))

	qb := query.NewSQLQueryBuilder(p.db, "param", "db", &selectParam.QueryOption)
	queryUpdate, args, err := qb.BuildUpdate(&updateParam, &selectParam)
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	res, err := p.db.Exec(ctx, "uPayslip", updatePayslip+queryUpdate, args...)
	pgErr := &pq.Error{}
	if err != nil && errors.As(err, &pgErr) && pgErr.Code == entity.PSQLUniqueConstraintCode {
		return errors.NewWithCode(codes.CodeSQLUniqueConstraint, err.Error())
	} else if err != nil {
		return errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, err.Error())
	} else if rowCount < 1 {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no payslip updated")
	}

	p.log.Debug(ctx, fmt.Sprintf("success update payslip with body: %v", updateParam))

	return nil
}
