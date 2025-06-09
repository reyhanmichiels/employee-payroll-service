package reimbursement

import (
	"context"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/query"
	"github.com/reyhanmichiels/go-pkg/v2/sql"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/business/entity"
)

func (r *reimbursement) getSQL(ctx context.Context, param entity.ReimbursementParam) (entity.Reimbursement, error) {
	reimbursement := entity.Reimbursement{}

	r.log.Debug(ctx, fmt.Sprintf("get reimbursement with body: %v", param))

	param.QueryOption.DisableLimit = true
	qb := query.NewSQLQueryBuilder(r.db, "param", "db", &param.QueryOption)
	queryExt, queryArgs, _, _, err := qb.Build(&param)
	if err != nil {
		return reimbursement, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	row, err := r.db.QueryRow(ctx, "rReimbursement", readReimbursement+queryExt, queryArgs...)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return reimbursement, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	if err := row.StructScan(&reimbursement); err != nil && errors.Is(err, sql.ErrNotFound) {
		return reimbursement, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, err.Error())
	} else if err != nil {
		return reimbursement, errors.NewWithCode(codes.CodeSQLRowScan, err.Error())
	}

	r.log.Debug(ctx, fmt.Sprintf("success get reimbursement with body: %v", param))

	return reimbursement, nil
}

func (r *reimbursement) getListSQL(ctx context.Context, param entity.ReimbursementParam) ([]entity.Reimbursement, *entity.Pagination, error) {
	reimbursementList := []entity.Reimbursement{}
	pg := entity.Pagination{}

	r.log.Debug(ctx, fmt.Sprintf("get reimbursement list with body: %v", param))

	qb := query.NewSQLQueryBuilder(r.db, "param", "db", &param.QueryOption)
	queryExt, queryArgs, countExt, countArgs, err := qb.Build(&param)
	if err != nil {
		return reimbursementList, &pg, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	rows, err := r.db.Query(ctx, "rReimbursementList", readReimbursement+queryExt, queryArgs...)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return reimbursementList, &pg, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		reimbursement := entity.Reimbursement{}
		err := rows.StructScan(&reimbursement)
		if err != nil {
			r.log.Error(ctx, errors.NewWithCode(codes.CodeSQLRowScan, err.Error()))
			continue
		}

		reimbursementList = append(reimbursementList, reimbursement)
	}

	pg = entity.Pagination{
		CurrentPage:     param.PaginationParam.Page,
		CurrentElements: int64(len(reimbursementList)),
		SortBy:          param.SortBy,
	}

	if !param.QueryOption.DisableLimit && len(reimbursementList) > 0 {
		err := r.db.Get(ctx, "cReimbursementList", countReimbursement+countExt, &pg.TotalElements, countArgs...)
		if err != nil {
			return reimbursementList, &pg, errors.NewWithCode(codes.CodeSQLRead, err.Error())
		}
	}

	pg.ProcessPagination(param.Limit)

	r.log.Debug(ctx, fmt.Sprintf("success get reimbursement list with body: %v", param))

	return reimbursementList, &pg, nil
}

func (r *reimbursement) createSQL(ctx context.Context, inputParam entity.ReimbursementInputParam) (entity.Reimbursement, error) {
	reimbursement := entity.Reimbursement{}

	r.log.Debug(ctx, fmt.Sprintf("create reimbursement with body: %v", inputParam))

	res, err := r.db.NamedExec(ctx, "iNewReimbursement", insertReimbursement, inputParam)

	mysqlErr := &mysql.MySQLError{}
	// 1062 is the error code for duplicate entry
	if err != nil && errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return reimbursement, errors.NewWithCode(codes.CodeSQLUniqueConstraint, err.Error())
	} else if err != nil {
		return reimbursement, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return reimbursement, errors.NewWithCode(codes.CodeSQLNoRowsAffected, err.Error())
	} else if rowCount < 1 {
		return reimbursement, errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no reimbursement created")
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return reimbursement, errors.NewWithCode(codes.CodeSQLNoRowsAffected, err.Error())
	}

	r.log.Debug(ctx, fmt.Sprintf("success create reimbursement with body: %v", inputParam))

	// Construct the returned reimbursement object based on input parameters
	reimbursement = entity.Reimbursement{
		ID:          lastID,
		UserID:      inputParam.UserID,
		Description: inputParam.Description,
		Amount:      inputParam.Amount,
		Status:      1,
		CreatedAt:   inputParam.CreatedAt,
		CreatedBy:   inputParam.CreatedBy,
	}

	return reimbursement, nil
}

func (r *reimbursement) updateSQL(ctx context.Context, updateParam entity.ReimbursementUpdateParam, selectParam entity.ReimbursementParam) error {
	r.log.Debug(ctx, fmt.Sprintf("update reimbursement with body: %v", updateParam))

	qb := query.NewSQLQueryBuilder(r.db, "param", "db", &selectParam.QueryOption)
	queryUpdate, args, err := qb.BuildUpdate(&updateParam, &selectParam)
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	res, err := r.db.Exec(ctx, "uReimbursement", updateReimbursement+queryUpdate, args...)
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
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no reimbursement updated")
	}

	r.log.Debug(ctx, fmt.Sprintf("success update reimbursement with body: %v", updateParam))

	return nil
}
