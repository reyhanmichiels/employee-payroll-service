package attendance

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

func (a *attendance) getSQL(ctx context.Context, param entity.AttendanceParam) (entity.Attendance, error) {
	attendance := entity.Attendance{}

	a.log.Debug(ctx, fmt.Sprintf("get attendance with body: %v", param))

	param.QueryOption.DisableLimit = true
	qb := query.NewSQLQueryBuilder(a.db, "param", "db", &param.QueryOption)
	queryExt, queryArgs, _, _, err := qb.Build(&param)
	if err != nil {
		return attendance, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	row, err := a.db.QueryRow(ctx, "rAttendance", readAttendance+queryExt, queryArgs...)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return attendance, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	if err := row.StructScan(&attendance); err != nil && errors.Is(err, sql.ErrNotFound) {
		return attendance, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, err.Error())
	} else if err != nil {
		return attendance, errors.NewWithCode(codes.CodeSQLRowScan, err.Error())
	}

	a.log.Debug(ctx, fmt.Sprintf("success get attendance with body: %v", param))

	return attendance, nil
}

func (a *attendance) getListSQL(ctx context.Context, param entity.AttendanceParam) ([]entity.Attendance, *entity.Pagination, error) {
	attendanceList := []entity.Attendance{}
	pg := entity.Pagination{}

	a.log.Debug(ctx, fmt.Sprintf("get attendance list with body: %v", param))

	qb := query.NewSQLQueryBuilder(a.db, "param", "db", &param.QueryOption)
	queryExt, queryArgs, countExt, countArgs, err := qb.Build(&param)
	if err != nil {
		return attendanceList, &pg, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	rows, err := a.db.Query(ctx, "rAttendanceList", readAttendance+queryExt, queryArgs...)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return attendanceList, &pg, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		attendance := entity.Attendance{}
		err := rows.StructScan(&attendance)
		if err != nil {
			a.log.Error(ctx, errors.NewWithCode(codes.CodeSQLRowScan, err.Error()))
			continue
		}

		attendanceList = append(attendanceList, attendance)
	}

	pg = entity.Pagination{
		CurrentPage:     param.PaginationParam.Page,
		CurrentElements: int64(len(attendanceList)),
		SortBy:          param.SortBy,
	}

	if !param.QueryOption.DisableLimit && len(attendanceList) > 0 {
		err := a.db.Get(ctx, "cAttendanceList", countAttendance+countExt, &pg.TotalElements, countArgs...)
		if err != nil {
			return attendanceList, &pg, errors.NewWithCode(codes.CodeSQLRead, err.Error())
		}
	}

	pg.ProcessPagination(param.Limit)

	a.log.Debug(ctx, fmt.Sprintf("success get attendance list with body: %v", param))

	return attendanceList, &pg, nil
}

func (a *attendance) createSQL(ctx context.Context, inputParam entity.AttendanceInputParam) (entity.Attendance, error) {
	attendance := entity.Attendance{}

	a.log.Debug(ctx, fmt.Sprintf("create attendance with body: %v", inputParam))

	stmt, err := a.db.PrepareNamed(ctx, "iNewAttendance", insertAttendance)
	if err != nil {
		return attendance, errors.NewWithCode(codes.CodeSQLPrepareStmt, err.Error())
	}
	defer stmt.Close()

	err = stmt.Get(&attendance, inputParam)

	pgErr := &pq.Error{}
	if err != nil && errors.As(err, &pgErr) && pgErr.Code == entity.PSQLUniqueConstraintCode {
		return attendance, errors.NewWithCode(codes.CodeSQLUniqueConstraint, err.Error())
	} else if err != nil {
		return attendance, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	a.log.Debug(ctx, fmt.Sprintf("success create attendance with body: %v", inputParam))

	return attendance, nil
}

func (a *attendance) updateSQL(ctx context.Context, updateParam entity.AttendanceUpdateParam, selectParam entity.AttendanceParam) error {
	a.log.Debug(ctx, fmt.Sprintf("update attendance with body: %v", updateParam))

	qb := query.NewSQLQueryBuilder(a.db, "param", "db", &selectParam.QueryOption)
	queryUpdate, args, err := qb.BuildUpdate(&updateParam, &selectParam)
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	res, err := a.db.Exec(ctx, "uAttendance", updateAttendance+queryUpdate, args...)
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
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no attendance updated")
	}

	a.log.Debug(ctx, fmt.Sprintf("success update attendance with body: %v", updateParam))

	return nil
}

func (a *attendance) countUserAttendanceSQL(ctx context.Context, attendancePeriodID int64) (entity.UserAttendanceCount, error) {
	userAttendanceCount := make(entity.UserAttendanceCount)

	rows, err := a.db.Query(ctx, "cUserAttendance", countUserAttendance, attendancePeriodID)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return userAttendanceCount, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var (
			userID          int64
			attendanceCount int64
		)

		err := rows.Scan(&userID, &attendanceCount)
		if err != nil {
			a.log.Error(ctx, errors.NewWithCode(codes.CodeSQLRowScan, err.Error()))
			continue
		}

		userAttendanceCount[userID] = attendanceCount
	}

	return userAttendanceCount, nil
}
