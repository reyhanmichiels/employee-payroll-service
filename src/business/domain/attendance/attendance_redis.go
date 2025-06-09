package attendance

import (
	"context"
	"fmt"
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/business/entity"
)

const (
	getAttendanceByKey           = "employeePayroll:attendance:get:%s"
	getAttendanceByQueryKey      = "employeePayroll:attendance:get:q:%s"
	getAttendanceByPaginationKey = "employeePayroll:attendance:get:p:%s"
	deleteAttendanceKeysPattern  = "employeePayroll:attendance*"
)

func (a *attendance) upsertCache(ctx context.Context, key string, attendance entity.Attendance, ttl time.Duration) error {
	marshalledAttendance, err := a.json.Marshal(attendance)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	err = a.redis.SetEX(ctx, key, string(marshalledAttendance), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheSetSimpleKey, err.Error())
	}

	return nil
}

func (a *attendance) getCache(ctx context.Context, key string) (entity.Attendance, error) {
	attendance := entity.Attendance{}

	marshalledAttendance, err := a.redis.Get(ctx, key)
	if err != nil {
		return attendance, err
	}

	err = a.json.Unmarshal([]byte(marshalledAttendance), &attendance)
	if err != nil {
		return attendance, errors.NewWithCode(codes.CodeCacheUnmarshal, err.Error())
	}

	return attendance, nil
}

func (a *attendance) upsertCacheList(ctx context.Context, param entity.AttendanceParam, attendanceList []entity.Attendance, pg entity.Pagination, ttl time.Duration) error {
	keyValue, err := a.json.Marshal(param)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	// Set attendance list to cache
	marshalledAttendanceList, err := a.json.Marshal(attendanceList)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}
	err = a.redis.SetEX(ctx, fmt.Sprintf(getAttendanceByQueryKey, string(keyValue)), string(marshalledAttendanceList), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheSetSimpleKey, err.Error())
	}

	// Set pagination to cache
	marshalledPagination, err := a.json.Marshal(pg)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	err = a.redis.SetEX(ctx, fmt.Sprintf(getAttendanceByPaginationKey, string(keyValue)), string(marshalledPagination), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheSetSimpleKey, err.Error())
	}

	return nil
}

func (a *attendance) getCacheList(ctx context.Context, param entity.AttendanceParam) ([]entity.Attendance, entity.Pagination, error) {
	var (
		attendanceList = []entity.Attendance{}
		pg             = entity.Pagination{}
	)

	keyValue, err := a.json.Marshal(param)
	if err != nil {
		return attendanceList, pg, errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	// Get attendance list from redis
	marshalledAttendanceList, err := a.redis.Get(ctx, fmt.Sprintf(getAttendanceByQueryKey, string(keyValue)))
	if err != nil {
		return attendanceList, pg, err
	}

	err = a.json.Unmarshal([]byte(marshalledAttendanceList), &attendanceList)
	if err != nil {
		return attendanceList, pg, errors.NewWithCode(codes.CodeCacheUnmarshal, err.Error())
	}

	// Get pagination from redis
	marshalledPagination, err := a.redis.Get(ctx, fmt.Sprintf(getAttendanceByPaginationKey, string(keyValue)))
	if err != nil {
		return attendanceList, pg, err
	}

	err = a.json.Unmarshal([]byte(marshalledPagination), &pg)
	if err != nil {
		return attendanceList, pg, errors.NewWithCode(codes.CodeCacheUnmarshal, err.Error())
	}

	return attendanceList, pg, nil
}

func (a *attendance) deleteCache(ctx context.Context, key string) error {
	err := a.redis.Del(ctx, key)
	if err != nil {
		return err
	}

	return nil
}
