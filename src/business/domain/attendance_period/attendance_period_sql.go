package attendance_period

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

func (a *attendancePeriod) getSQL(ctx context.Context, param entity.AttendancePeriodParam) (entity.AttendancePeriod, error) {
	attendancePeriod := entity.AttendancePeriod{}

	a.log.Debug(ctx, fmt.Sprintf("get attendance period with body: %v", param))

	param.QueryOption.DisableLimit = true
	qb := query.NewSQLQueryBuilder(a.db, "param", "db", &param.QueryOption)
	queryExt, queryArgs, _, _, err := qb.Build(&param)
	if err != nil {
		return attendancePeriod, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	row, err := a.db.QueryRow(ctx, "rAttendancePeriod", readAttendancePeriod+queryExt, queryArgs...)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return attendancePeriod, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	if err := row.StructScan(&attendancePeriod); err != nil && errors.Is(err, sql.ErrNotFound) {
		return attendancePeriod, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, err.Error())
	} else if err != nil {
		return attendancePeriod, errors.NewWithCode(codes.CodeSQLRowScan, err.Error())
	}

	a.log.Debug(ctx, fmt.Sprintf("success get attendance period with body: %v", param))

	return attendancePeriod, nil
}

func (a *attendancePeriod) getListSQL(ctx context.Context, param entity.AttendancePeriodParam) ([]entity.AttendancePeriod, *entity.Pagination, error) {
	attendancePeriodList := []entity.AttendancePeriod{}
	pg := entity.Pagination{}

	a.log.Debug(ctx, fmt.Sprintf("get attendance period list with body: %v", param))

	qb := query.NewSQLQueryBuilder(a.db, "param", "db", &param.QueryOption)
	queryExt, queryArgs, countExt, countArgs, err := qb.Build(&param)
	if err != nil {
		return attendancePeriodList, &pg, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	rows, err := a.db.Query(ctx, "rAttendancePeriodList", readAttendancePeriod+queryExt, queryArgs...)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return attendancePeriodList, &pg, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		attendancePeriod := entity.AttendancePeriod{}
		err := rows.StructScan(&attendancePeriod)
		if err != nil {
			a.log.Error(ctx, errors.NewWithCode(codes.CodeSQLRowScan, err.Error()))
			continue
		}

		attendancePeriodList = append(attendancePeriodList, attendancePeriod)
	}

	pg = entity.Pagination{
		CurrentPage:     param.PaginationParam.Page,
		CurrentElements: int64(len(attendancePeriodList)),
		SortBy:          param.SortBy,
	}

	if !param.QueryOption.DisableLimit && len(attendancePeriodList) > 0 {
		err := a.db.Get(ctx, "cAttendancePeriodList", countAttendancePeriod+countExt, &pg.TotalElements, countArgs...)
		if err != nil {
			return attendancePeriodList, &pg, errors.NewWithCode(codes.CodeSQLRead, err.Error())
		}
	}

	pg.ProcessPagination(param.Limit)

	a.log.Debug(ctx, fmt.Sprintf("success get attendance period list with body: %v", param))

	return attendancePeriodList, &pg, nil
}

func (a *attendancePeriod) createSQL(ctx context.Context, inputParam entity.AttendancePeriodInputParam) (entity.AttendancePeriod, error) {
	attendancePeriod := entity.AttendancePeriod{}

	a.log.Debug(ctx, fmt.Sprintf("create attendance period with body: %v", inputParam))

	stmt, err := a.db.PrepareNamed(ctx, "iNewAttendancePeriod", insertAttendancePeriod)
	if err != nil {
		return attendancePeriod, errors.NewWithCode(codes.CodeSQLPrepareStmt, err.Error())
	}
	defer stmt.Close()

	err = stmt.Get(&attendancePeriod, inputParam)

	pgErr := &pq.Error{}
	if err != nil && errors.As(err, &pgErr) {
		switch pgErr.Code {
		case entity.PSQLExclusionConstraintCode:
			return attendancePeriod, errors.NewWithCode(codes.CodeSQLUniqueConstraint, err.Error())
		case entity.PSQLUniqueConstraintCode:
			return attendancePeriod, errors.NewWithCode(codes.CodeSQLUniqueConstraint, err.Error())
		default:
			return attendancePeriod, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
		}
	} else if err != nil {
		return attendancePeriod, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	a.log.Debug(ctx, fmt.Sprintf("success create attendance period with body: %v", inputParam))

	return attendancePeriod, nil
}

func (a *attendancePeriod) updateSQL(ctx context.Context, updateParam entity.AttendancePeriodUpdateParam, selectParam entity.AttendancePeriodParam) error {
	a.log.Debug(ctx, fmt.Sprintf("update attendance period with body: %v", updateParam))

	qb := query.NewSQLQueryBuilder(a.db, "param", "db", &selectParam.QueryOption)
	queryUpdate, args, err := qb.BuildUpdate(&updateParam, &selectParam)
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	res, err := a.db.Exec(ctx, "uAttendancePeriod", updateAttendancePeriod+queryUpdate, args...)

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
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no attendance period updated")
	}

	a.log.Debug(ctx, fmt.Sprintf("success update attendance period with body: %v", updateParam))

	return nil
}
