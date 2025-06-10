package overtime

import (
	"context"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/query"
	"github.com/reyhanmichiels/go-pkg/v2/sql"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
)

func (o *overtime) getSQL(ctx context.Context, param entity.OvertimeParam) (entity.Overtime, error) {
	overtime := entity.Overtime{}

	o.log.Debug(ctx, fmt.Sprintf("get overtime with body: %v", param))

	param.QueryOption.DisableLimit = true
	qb := query.NewSQLQueryBuilder(o.db, "param", "db", &param.QueryOption)
	queryExt, queryArgs, _, _, err := qb.Build(&param)
	if err != nil {
		return overtime, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	row, err := o.db.QueryRow(ctx, "rOvertime", readOvertime+queryExt, queryArgs...)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return overtime, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	if err := row.StructScan(&overtime); err != nil && errors.Is(err, sql.ErrNotFound) {
		return overtime, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, err.Error())
	} else if err != nil {
		return overtime, errors.NewWithCode(codes.CodeSQLRowScan, err.Error())
	}

	o.log.Debug(ctx, fmt.Sprintf("success get overtime with body: %v", param))

	return overtime, nil
}

func (o *overtime) getListSQL(ctx context.Context, param entity.OvertimeParam) ([]entity.Overtime, *entity.Pagination, error) {
	overtimeList := []entity.Overtime{}
	pg := entity.Pagination{}

	o.log.Debug(ctx, fmt.Sprintf("get overtime list with body: %v", param))

	qb := query.NewSQLQueryBuilder(o.db, "param", "db", &param.QueryOption)
	queryExt, queryArgs, countExt, countArgs, err := qb.Build(&param)
	if err != nil {
		return overtimeList, &pg, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	rows, err := o.db.Query(ctx, "rOvertimeList", readOvertime+queryExt, queryArgs...)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return overtimeList, &pg, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		overtime := entity.Overtime{}
		err := rows.StructScan(&overtime)
		if err != nil {
			o.log.Error(ctx, errors.NewWithCode(codes.CodeSQLRowScan, err.Error()))
			continue
		}

		overtimeList = append(overtimeList, overtime)
	}

	pg = entity.Pagination{
		CurrentPage:     param.PaginationParam.Page,
		CurrentElements: int64(len(overtimeList)),
		SortBy:          param.SortBy,
	}

	if !param.QueryOption.DisableLimit && len(overtimeList) > 0 {
		err := o.db.Get(ctx, "cOvertimeList", countOvertime+countExt, &pg.TotalElements, countArgs...)
		if err != nil {
			return overtimeList, &pg, errors.NewWithCode(codes.CodeSQLRead, err.Error())
		}
	}

	pg.ProcessPagination(param.Limit)

	o.log.Debug(ctx, fmt.Sprintf("success get overtime list with body: %v", param))

	return overtimeList, &pg, nil
}

func (o *overtime) createSQL(ctx context.Context, inputParam entity.OvertimeInputParam) (entity.Overtime, error) {
	overtime := entity.Overtime{}

	o.log.Debug(ctx, fmt.Sprintf("create overtime with body: %v", inputParam))

	stmt, err := o.db.PrepareNamed(ctx, "iNewOvertime", insertOvertime)
	if err != nil {
		return overtime, errors.NewWithCode(codes.CodeSQLPrepareStmt, err.Error())
	}
	defer stmt.Close()

	err = stmt.Get(&overtime, inputParam)

	pgErr := &pq.Error{}
	if err != nil && errors.As(err, &pgErr) && pgErr.Code == entity.PSQLUniqueConstraintCode {
		return overtime, errors.NewWithCode(codes.CodeSQLUniqueConstraint, err.Error())
	} else if err != nil {
		return overtime, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	o.log.Debug(ctx, fmt.Sprintf("success create overtime with body: %v", inputParam))

	return overtime, nil
}

func (o *overtime) updateSQL(ctx context.Context, updateParam entity.OvertimeUpdateParam, selectParam entity.OvertimeParam) error {
	o.log.Debug(ctx, fmt.Sprintf("update overtime with body: %v", updateParam))

	qb := query.NewSQLQueryBuilder(o.db, "param", "db", &selectParam.QueryOption)
	queryUpdate, args, err := qb.BuildUpdate(&updateParam, &selectParam)
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	res, err := o.db.Exec(ctx, "uOvertime", updateOvertime+queryUpdate, args...)
	mysqlErr := &mysql.MySQLError{}
	// 1062 is the error code for duplicate entry
	if err != nil && errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return errors.NewWithCode(codes.CodeSQLUniqueConstraint, err.Error())
	} else if err != nil {
		return errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, err.Error())
	} else if rowCount < 1 {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no overtime updated")
	}

	o.log.Debug(ctx, fmt.Sprintf("success update overtime with body: %v", updateParam))

	return nil
}
