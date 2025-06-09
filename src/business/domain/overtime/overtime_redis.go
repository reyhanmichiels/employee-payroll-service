package overtime

import (
	"context"
	"fmt"
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/business/entity"
)

const (
	getOvertimeByKey           = "employeePayroll:overtime:get:%s"
	getOvertimeByQueryKey      = "employeePayroll:overtime:get:q:%s"
	getOvertimeByPaginationKey = "employeePayroll:overtime:get:p:%s"
	deleteOvertimeKeysPattern  = "employeePayroll:overtime*"
)

func (o *overtime) upsertCache(ctx context.Context, key string, overtime entity.Overtime, ttl time.Duration) error {
	marshalledOvertime, err := o.json.Marshal(overtime)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	err = o.redis.SetEX(ctx, key, string(marshalledOvertime), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheSetSimpleKey, err.Error())
	}

	return nil
}

func (o *overtime) getCache(ctx context.Context, key string) (entity.Overtime, error) {
	overtime := entity.Overtime{}

	marshalledOvertime, err := o.redis.Get(ctx, key)
	if err != nil {
		return overtime, err
	}

	err = o.json.Unmarshal([]byte(marshalledOvertime), &overtime)
	if err != nil {
		return overtime, errors.NewWithCode(codes.CodeCacheUnmarshal, err.Error())
	}

	return overtime, nil
}

func (o *overtime) upsertCacheList(ctx context.Context, param entity.OvertimeParam, overtimeList []entity.Overtime, pg entity.Pagination, ttl time.Duration) error {
	keyValue, err := o.json.Marshal(param)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	// Set overtime list to cache
	marshalledOvertimeList, err := o.json.Marshal(overtimeList)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}
	err = o.redis.SetEX(ctx, fmt.Sprintf(getOvertimeByQueryKey, string(keyValue)), string(marshalledOvertimeList), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheSetSimpleKey, err.Error())
	}

	// Set pagination to cache
	marshalledPagination, err := o.json.Marshal(pg)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	err = o.redis.SetEX(ctx, fmt.Sprintf(getOvertimeByPaginationKey, string(keyValue)), string(marshalledPagination), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheSetSimpleKey, err.Error())
	}

	return nil
}

func (o *overtime) getCacheList(ctx context.Context, param entity.OvertimeParam) ([]entity.Overtime, entity.Pagination, error) {
	var (
		overtimeList = []entity.Overtime{}
		pg           = entity.Pagination{}
	)

	keyValue, err := o.json.Marshal(param)
	if err != nil {
		return overtimeList, pg, errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	// Get overtime list from redis
	marshalledOvertimeList, err := o.redis.Get(ctx, fmt.Sprintf(getOvertimeByQueryKey, string(keyValue)))
	if err != nil {
		return overtimeList, pg, err
	}

	err = o.json.Unmarshal([]byte(marshalledOvertimeList), &overtimeList)
	if err != nil {
		return overtimeList, pg, errors.NewWithCode(codes.CodeCacheUnmarshal, err.Error())
	}

	// Get pagination from redis
	marshalledPagination, err := o.redis.Get(ctx, fmt.Sprintf(getOvertimeByPaginationKey, string(keyValue)))
	if err != nil {
		return overtimeList, pg, err
	}

	err = o.json.Unmarshal([]byte(marshalledPagination), &pg)
	if err != nil {
		return overtimeList, pg, errors.NewWithCode(codes.CodeCacheUnmarshal, err.Error())
	}

	return overtimeList, pg, nil
}

func (o *overtime) deleteCache(ctx context.Context, key string) error {
	err := o.redis.Del(ctx, key)
	if err != nil {
		return err
	}

	return nil
}
