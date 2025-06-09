package attendance_period

import (
	"context"
	"fmt"
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
)

const (
	getAttendancePeriodByKey           = "employeePayroll:attendanceperiod:get:%s"
	getAttendancePeriodByQueryKey      = "employeePayroll:attendanceperiod:get:q:%s"
	getAttendancePeriodByPaginationKey = "employeePayroll:attendanceperiod:get:p:%s"
	deleteAttendancePeriodKeysPattern  = "employeePayroll:attendanceperiod*"
)

func (a *attendancePeriod) upsertCache(ctx context.Context, key string, attendancePeriod entity.AttendancePeriod, ttl time.Duration) error {
	marshalledAttendancePeriod, err := a.json.Marshal(attendancePeriod)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	err = a.redis.SetEX(ctx, key, string(marshalledAttendancePeriod), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheSetSimpleKey, err.Error())
	}

	return nil
}

func (a *attendancePeriod) getCache(ctx context.Context, key string) (entity.AttendancePeriod, error) {
	attendancePeriod := entity.AttendancePeriod{}

	marshalledAttendancePeriod, err := a.redis.Get(ctx, key)
	if err != nil {
		return attendancePeriod, err
	}

	err = a.json.Unmarshal([]byte(marshalledAttendancePeriod), &attendancePeriod)
	if err != nil {
		return attendancePeriod, errors.NewWithCode(codes.CodeCacheUnmarshal, err.Error())
	}

	return attendancePeriod, nil
}

func (a *attendancePeriod) upsertCacheList(ctx context.Context, param entity.AttendancePeriodParam, attendancePeriodList []entity.AttendancePeriod, pg entity.Pagination, ttl time.Duration) error {
	keyValue, err := a.json.Marshal(param)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	// Set attendance period list to cache
	marshalledAttendancePeriodList, err := a.json.Marshal(attendancePeriodList)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}
	err = a.redis.SetEX(ctx, fmt.Sprintf(getAttendancePeriodByQueryKey, string(keyValue)), string(marshalledAttendancePeriodList), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheSetSimpleKey, err.Error())
	}

	// Set pagination to cache
	marshalledPagination, err := a.json.Marshal(pg)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	err = a.redis.SetEX(ctx, fmt.Sprintf(getAttendancePeriodByPaginationKey, string(keyValue)), string(marshalledPagination), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheSetSimpleKey, err.Error())
	}

	return nil
}

func (a *attendancePeriod) getCacheList(ctx context.Context, param entity.AttendancePeriodParam) ([]entity.AttendancePeriod, entity.Pagination, error) {
	var (
		attendancePeriodList = []entity.AttendancePeriod{}
		pg                   = entity.Pagination{}
	)

	keyValue, err := a.json.Marshal(param)
	if err != nil {
		return attendancePeriodList, pg, errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	// Get attendance period list from redis
	marshalledAttendancePeriodList, err := a.redis.Get(ctx, fmt.Sprintf(getAttendancePeriodByQueryKey, string(keyValue)))
	if err != nil {
		return attendancePeriodList, pg, err
	}

	err = a.json.Unmarshal([]byte(marshalledAttendancePeriodList), &attendancePeriodList)
	if err != nil {
		return attendancePeriodList, pg, errors.NewWithCode(codes.CodeCacheUnmarshal, err.Error())
	}

	// Get pagination from redis
	marshalledPagination, err := a.redis.Get(ctx, fmt.Sprintf(getAttendancePeriodByPaginationKey, string(keyValue)))
	if err != nil {
		return attendancePeriodList, pg, err
	}

	err = a.json.Unmarshal([]byte(marshalledPagination), &pg)
	if err != nil {
		return attendancePeriodList, pg, errors.NewWithCode(codes.CodeCacheUnmarshal, err.Error())
	}

	return attendancePeriodList, pg, nil
}

func (a *attendancePeriod) deleteCache(ctx context.Context, key string) error {
	err := a.redis.Del(ctx, key)
	if err != nil {
		return err
	}

	return nil
}
