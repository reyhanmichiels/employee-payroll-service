package reimbursement

import (
	"context"
	"fmt"
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/business/entity"
)

const (
	getReimbursementByKey           = "employeePayroll:reimbursement:get:%s"
	getReimbursementByQueryKey      = "employeePayroll:reimbursement:get:q:%s"
	getReimbursementByPaginationKey = "employeePayroll:reimbursement:get:p:%s"
	deleteReimbursementKeysPattern  = "employeePayroll:reimbursement*"
)

func (r *reimbursement) upsertCache(ctx context.Context, key string, reimbursement entity.Reimbursement, ttl time.Duration) error {
	marshalledReimbursement, err := r.json.Marshal(reimbursement)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	err = r.redis.SetEX(ctx, key, string(marshalledReimbursement), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheSetSimpleKey, err.Error())
	}

	return nil
}

func (r *reimbursement) getCache(ctx context.Context, key string) (entity.Reimbursement, error) {
	reimbursement := entity.Reimbursement{}

	marshalledReimbursement, err := r.redis.Get(ctx, key)
	if err != nil {
		return reimbursement, err
	}

	err = r.json.Unmarshal([]byte(marshalledReimbursement), &reimbursement)
	if err != nil {
		return reimbursement, errors.NewWithCode(codes.CodeCacheUnmarshal, err.Error())
	}

	return reimbursement, nil
}

func (r *reimbursement) upsertCacheList(ctx context.Context, param entity.ReimbursementParam, reimbursementList []entity.Reimbursement, pg entity.Pagination, ttl time.Duration) error {
	keyValue, err := r.json.Marshal(param)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	// Set reimbursement list to cache
	marshalledReimbursementList, err := r.json.Marshal(reimbursementList)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}
	err = r.redis.SetEX(ctx, fmt.Sprintf(getReimbursementByQueryKey, string(keyValue)), string(marshalledReimbursementList), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheSetSimpleKey, err.Error())
	}

	// Set pagination to cache
	marshalledPagination, err := r.json.Marshal(pg)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	err = r.redis.SetEX(ctx, fmt.Sprintf(getReimbursementByPaginationKey, string(keyValue)), string(marshalledPagination), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheSetSimpleKey, err.Error())
	}

	return nil
}

func (r *reimbursement) getCacheList(ctx context.Context, param entity.ReimbursementParam) ([]entity.Reimbursement, entity.Pagination, error) {
	var (
		reimbursementList = []entity.Reimbursement{}
		pg                = entity.Pagination{}
	)

	keyValue, err := r.json.Marshal(param)
	if err != nil {
		return reimbursementList, pg, errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	// Get reimbursement list from redis
	marshalledReimbursementList, err := r.redis.Get(ctx, fmt.Sprintf(getReimbursementByQueryKey, string(keyValue)))
	if err != nil {
		return reimbursementList, pg, err
	}

	err = r.json.Unmarshal([]byte(marshalledReimbursementList), &reimbursementList)
	if err != nil {
		return reimbursementList, pg, errors.NewWithCode(codes.CodeCacheUnmarshal, err.Error())
	}

	// Get pagination from redis
	marshalledPagination, err := r.redis.Get(ctx, fmt.Sprintf(getReimbursementByPaginationKey, string(keyValue)))
	if err != nil {
		return reimbursementList, pg, err
	}

	err = r.json.Unmarshal([]byte(marshalledPagination), &pg)
	if err != nil {
		return reimbursementList, pg, errors.NewWithCode(codes.CodeCacheUnmarshal, err.Error())
	}

	return reimbursementList, pg, nil
}

func (r *reimbursement) deleteCache(ctx context.Context, key string) error {
	err := r.redis.Del(ctx, key)
	if err != nil {
		return err
	}

	return nil
}
