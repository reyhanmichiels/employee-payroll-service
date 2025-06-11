package payslip_detail

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

func (p *payslipDetail) getSQL(ctx context.Context, param entity.PayslipDetailParam) (entity.PayslipDetail, error) {
	payslipDetail := entity.PayslipDetail{}

	p.log.Debug(ctx, fmt.Sprintf("get payslip detail with body: %v", param))

	param.QueryOption.DisableLimit = true
	qb := query.NewSQLQueryBuilder(p.db, "param", "db", &param.QueryOption)
	queryExt, queryArgs, _, _, err := qb.Build(&param)
	if err != nil {
		return payslipDetail, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	row, err := p.db.QueryRow(ctx, "rPayslipDetail", readPayslipDetail+queryExt, queryArgs...)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return payslipDetail, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	if err := row.StructScan(&payslipDetail); err != nil && errors.Is(err, sql.ErrNotFound) {
		return payslipDetail, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, err.Error())
	} else if err != nil {
		return payslipDetail, errors.NewWithCode(codes.CodeSQLRowScan, err.Error())
	}

	p.log.Debug(ctx, fmt.Sprintf("success get payslip detail with body: %v", param))

	return payslipDetail, nil
}

func (p *payslipDetail) getListSQL(ctx context.Context, param entity.PayslipDetailParam) ([]entity.PayslipDetail, *entity.Pagination, error) {
	payslipDetailList := []entity.PayslipDetail{}
	pg := entity.Pagination{}

	p.log.Debug(ctx, fmt.Sprintf("get payslip detail list with body: %v", param))

	qb := query.NewSQLQueryBuilder(p.db, "param", "db", &param.QueryOption)
	queryExt, queryArgs, countExt, countArgs, err := qb.Build(&param)
	if err != nil {
		return payslipDetailList, &pg, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	rows, err := p.db.Query(ctx, "rPayslipDetailList", readPayslipDetail+queryExt, queryArgs...)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return payslipDetailList, &pg, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		payslipDetail := entity.PayslipDetail{}
		err := rows.StructScan(&payslipDetail)
		if err != nil {
			p.log.Error(ctx, errors.NewWithCode(codes.CodeSQLRowScan, err.Error()))
			continue
		}

		payslipDetailList = append(payslipDetailList, payslipDetail)
	}

	pg = entity.Pagination{
		CurrentPage:     param.PaginationParam.Page,
		CurrentElements: int64(len(payslipDetailList)),
		SortBy:          param.SortBy,
	}

	if !param.QueryOption.DisableLimit && len(payslipDetailList) > 0 {
		err := p.db.Get(ctx, "cPayslipDetailList", countPayslipDetail+countExt, &pg.TotalElements, countArgs...)
		if err != nil {
			return payslipDetailList, &pg, errors.NewWithCode(codes.CodeSQLRead, err.Error())
		}
	}

	pg.ProcessPagination(param.Limit)

	p.log.Debug(ctx, fmt.Sprintf("success get payslip detail list with body: %v", param))

	return payslipDetailList, &pg, nil
}

func (p *payslipDetail) createSQL(ctx context.Context, inputParam entity.PayslipDetailInputParam) (entity.PayslipDetail, error) {
	payslipDetail := entity.PayslipDetail{}

	p.log.Debug(ctx, fmt.Sprintf("create payslip detail with body: %v", inputParam))

	stmt, err := p.db.PrepareNamed(ctx, "iNewPayslipDetail", insertPayslipDetail)
	if err != nil {
		return payslipDetail, errors.NewWithCode(codes.CodeSQLPrepareStmt, err.Error())
	}
	defer stmt.Close()

	err = stmt.Get(&payslipDetail, inputParam)

	pgErr := &pq.Error{}
	if err != nil && errors.As(err, &pgErr) && pgErr.Code == entity.PSQLUniqueConstraintCode {
		return payslipDetail, errors.NewWithCode(codes.CodeSQLUniqueConstraint, err.Error())
	} else if err != nil {
		return payslipDetail, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	p.log.Debug(ctx, fmt.Sprintf("success create payslip detail with body: %v", inputParam))

	return payslipDetail, nil
}

func (p *payslipDetail) createManySQL(ctx context.Context, inputParams []entity.PayslipDetailInputParam) error {
	p.log.Debug(ctx, fmt.Sprintf("create many payslip detail with body: %v", inputParams))

	res, err := p.db.NamedExec(ctx, "iManyPayslipDetail", insertManyPayslipDetail, inputParams)

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
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no payslip detail created")
	}

	p.log.Debug(ctx, fmt.Sprintf("success create many payslip detail with body: %v", inputParams))

	return nil
}

func (p *payslipDetail) updateSQL(ctx context.Context, updateParam entity.PayslipDetailUpdateParam, selectParam entity.PayslipDetailParam) error {
	p.log.Debug(ctx, fmt.Sprintf("update payslip detail with body: %v", updateParam))

	qb := query.NewSQLQueryBuilder(p.db, "param", "db", &selectParam.QueryOption)
	queryUpdate, args, err := qb.BuildUpdate(&updateParam, &selectParam)
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	res, err := p.db.Exec(ctx, "uPayslipDetail", updatePayslipDetail+queryUpdate, args...)
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
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no payslip detail updated")
	}

	p.log.Debug(ctx, fmt.Sprintf("success update payslip detail with body: %v", updateParam))

	return nil
}
