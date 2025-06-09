package attendance_period

import (
	"context"
	"fmt"

	"github.com/go-sql-driver/mysql"
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

	res, err := a.db.NamedExec(ctx, "iNewAttendancePeriod", insertAttendancePeriod, inputParam)

	mysqlErr := &mysql.MySQLError{}
	// 1062 is the error code for duplicate entry
	if err != nil && errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return attendancePeriod, errors.NewWithCode(codes.CodeSQLUniqueConstraint, err.Error())
	} else if err != nil {
		return attendancePeriod, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return attendancePeriod, errors.NewWithCode(codes.CodeSQLNoRowsAffected, err.Error())
	} else if rowCount < 1 {
		return attendancePeriod, errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no attendance period created")
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return attendancePeriod, errors.NewWithCode(codes.CodeSQLNoRowsAffected, err.Error())
	}

	a.log.Debug(ctx, fmt.Sprintf("success create attendance period with body: %v", inputParam))

	// Construct the returned attendance period object based on input parameters
	attendancePeriod = entity.AttendancePeriod{
		ID:           lastID,
		StartDate:    inputParam.StartDate,
		EndDate:      inputParam.EndDate,
		PeriodStatus: inputParam.PeriodStatus,
		Status:       1,
		CreatedAt:    inputParam.CreatedAt,
		CreatedBy:    inputParam.CreatedBy,
	}

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
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no attendance period updated")
	}

	a.log.Debug(ctx, fmt.Sprintf("success update attendance period with body: %v", updateParam))

	return nil
}
