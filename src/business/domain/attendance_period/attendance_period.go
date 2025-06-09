package attendance_period

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
	Get(ctx context.Context, param entity.AttendancePeriodParam) (entity.AttendancePeriod, error)
	GetList(ctx context.Context, param entity.AttendancePeriodParam) ([]entity.AttendancePeriod, *entity.Pagination, error)
	Create(ctx context.Context, param entity.AttendancePeriodInputParam) (entity.AttendancePeriod, error)
	Update(ctx context.Context, updateParam entity.AttendancePeriodUpdateParam, selectParam entity.AttendancePeriodParam) error
}

type attendancePeriod struct {
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
	return &attendancePeriod{
		db:    param.Db,
		log:   param.Log,
		redis: param.Redis,
		json:  param.Json,
	}
}

func (a *attendancePeriod) Get(ctx context.Context, param entity.AttendancePeriodParam) (entity.AttendancePeriod, error) {
	attendancePeriod := entity.AttendancePeriod{}

	marshalledParam, err := a.json.Marshal(param)
	if err != nil {
		return attendancePeriod, err
	}

	if !param.BypassCache {
		attendancePeriod, err = a.getCache(ctx, fmt.Sprintf(getAttendancePeriodByKey, string(marshalledParam)))
		switch {
		case errors.Is(err, redis.Nil):
			a.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedisNil, err.Error()))
		case err != nil:
			a.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
		default:
			return attendancePeriod, nil
		}
	}

	attendancePeriod, err = a.getSQL(ctx, param)
	if err != nil {
		return attendancePeriod, err
	}

	err = a.upsertCache(ctx, fmt.Sprintf(getAttendancePeriodByKey, string(marshalledParam)), attendancePeriod, a.redis.GetDefaultTTL(ctx))
	if err != nil {
		a.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return attendancePeriod, nil
}

func (a *attendancePeriod) GetList(ctx context.Context, param entity.AttendancePeriodParam) ([]entity.AttendancePeriod, *entity.Pagination, error) {
	if !param.BypassCache {
		attendancePeriodList, pg, err := a.getCacheList(ctx, param)
		switch {
		case errors.Is(err, redis.Nil):
			a.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedisNil, err.Error()))
		case err != nil:
			a.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
		default:
			return attendancePeriodList, &pg, nil
		}
	}

	attendancePeriodList, pg, err := a.getListSQL(ctx, param)
	if err != nil {
		return attendancePeriodList, pg, err
	}

	err = a.upsertCacheList(ctx, param, attendancePeriodList, *pg, a.redis.GetDefaultTTL(ctx))
	if err != nil {
		a.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return attendancePeriodList, pg, nil
}

func (a *attendancePeriod) Create(ctx context.Context, param entity.AttendancePeriodInputParam) (entity.AttendancePeriod, error) {
	attendancePeriod, err := a.createSQL(ctx, param)
	if err != nil {
		return attendancePeriod, err
	}

	err = a.deleteCache(ctx, deleteAttendancePeriodKeysPattern)
	if err != nil {
		a.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return attendancePeriod, nil
}

func (a *attendancePeriod) Update(ctx context.Context, updateParam entity.AttendancePeriodUpdateParam, selectParam entity.AttendancePeriodParam) error {
	err := a.updateSQL(ctx, updateParam, selectParam)
	if err != nil {
		return err
	}

	err = a.deleteCache(ctx, deleteAttendancePeriodKeysPattern)
	if err != nil {
		a.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return nil
}
