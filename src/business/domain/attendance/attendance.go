package attendance

import (
	"context"
	"fmt"

	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/log"
	"github.com/reyhanmichiels/go-pkg/v2/parser"
	"github.com/reyhanmichiels/go-pkg/v2/redis"
	"github.com/reyhanmichiels/go-pkg/v2/sql"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
)

type Interface interface {
	Get(ctx context.Context, param entity.AttendanceParam) (entity.Attendance, error)
	GetList(ctx context.Context, param entity.AttendanceParam) ([]entity.Attendance, *entity.Pagination, error)
	CountUserAttendance(ctx context.Context, attendancePeriodID int64) (entity.UserAttendanceCount, error)
	Create(ctx context.Context, param entity.AttendanceInputParam) (entity.Attendance, error)
	Update(ctx context.Context, updateParam entity.AttendanceUpdateParam, selectParam entity.AttendanceParam) error
}

type attendance struct {
	db    sql.Interface
	log   log.Interface
	redis redis.Interface
	json  parser.JSONInterface
}

type InitParam struct {
	Db    sql.Interface
	Log   log.Interface
	Redis redis.Interface
	Json  parser.JSONInterface
}

func Init(param InitParam) Interface {
	return &attendance{
		db:    param.Db,
		log:   param.Log,
		redis: param.Redis,
		json:  param.Json,
	}
}

func (a *attendance) Get(ctx context.Context, param entity.AttendanceParam) (entity.Attendance, error) {
	attendance := entity.Attendance{}

	marshalledParam, err := a.json.Marshal(param)
	if err != nil {
		return attendance, err
	}

	if !param.BypassCache {
		attendance, err = a.getCache(ctx, fmt.Sprintf(getAttendanceByKey, string(marshalledParam)))
		switch {
		case errors.Is(err, redis.Nil):
			a.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedisNil, err.Error()))
		case err != nil:
			a.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
		default:
			return attendance, nil
		}
	}

	attendance, err = a.getSQL(ctx, param)
	if err != nil {
		return attendance, err
	}

	err = a.upsertCache(ctx, fmt.Sprintf(getAttendanceByKey, string(marshalledParam)), attendance, a.redis.GetDefaultTTL(ctx))
	if err != nil {
		a.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return attendance, nil
}

func (a *attendance) GetList(ctx context.Context, param entity.AttendanceParam) ([]entity.Attendance, *entity.Pagination, error) {
	if !param.BypassCache {
		attendanceList, pg, err := a.getCacheList(ctx, param)
		switch {
		case errors.Is(err, redis.Nil):
			a.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedisNil, err.Error()))
		case err != nil:
			a.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
		default:
			return attendanceList, &pg, nil
		}
	}

	attendanceList, pg, err := a.getListSQL(ctx, param)
	if err != nil {
		return attendanceList, pg, err
	}

	err = a.upsertCacheList(ctx, param, attendanceList, *pg, a.redis.GetDefaultTTL(ctx))
	if err != nil {
		a.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return attendanceList, pg, nil
}

func (a *attendance) Create(ctx context.Context, param entity.AttendanceInputParam) (entity.Attendance, error) {
	attendance, err := a.createSQL(ctx, param)
	if err != nil {
		return attendance, err
	}

	err = a.deleteCache(ctx, deleteAttendanceKeysPattern)
	if err != nil {
		a.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return attendance, nil
}

func (a *attendance) Update(ctx context.Context, updateParam entity.AttendanceUpdateParam, selectParam entity.AttendanceParam) error {
	err := a.updateSQL(ctx, updateParam, selectParam)
	if err != nil {
		return err
	}

	err = a.deleteCache(ctx, deleteAttendanceKeysPattern)
	if err != nil {
		a.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return nil
}

func (a *attendance) CountUserAttendance(ctx context.Context, attendancePeriodID int64) (entity.UserAttendanceCount, error) {
	userAttendancePeriod, err := a.countUserAttendanceSQL(ctx, attendancePeriodID)
	if err != nil {
		return userAttendancePeriod, err
	}

	return userAttendancePeriod, nil
}
