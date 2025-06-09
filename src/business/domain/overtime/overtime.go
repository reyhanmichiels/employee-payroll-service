package overtime

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
	Get(ctx context.Context, param entity.OvertimeParam) (entity.Overtime, error)
	GetList(ctx context.Context, param entity.OvertimeParam) ([]entity.Overtime, *entity.Pagination, error)
	Create(ctx context.Context, param entity.OvertimeInputParam) (entity.Overtime, error)
	Update(ctx context.Context, updateParam entity.OvertimeUpdateParam, selectParam entity.OvertimeParam) error
}

type overtime struct {
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
	return &overtime{
		db:    param.Db,
		log:   param.Log,
		redis: param.Redis,
		json:  param.Json,
	}
}

func (o *overtime) Get(ctx context.Context, param entity.OvertimeParam) (entity.Overtime, error) {
	overtime := entity.Overtime{}

	marshalledParam, err := o.json.Marshal(param)
	if err != nil {
		return overtime, err
	}

	if !param.BypassCache {
		overtime, err = o.getCache(ctx, fmt.Sprintf(getOvertimeByKey, string(marshalledParam)))
		switch {
		case errors.Is(err, redis.Nil):
			o.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedisNil, err.Error()))
		case err != nil:
			o.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
		default:
			return overtime, nil
		}
	}

	overtime, err = o.getSQL(ctx, param)
	if err != nil {
		return overtime, err
	}

	err = o.upsertCache(ctx, fmt.Sprintf(getOvertimeByKey, string(marshalledParam)), overtime, o.redis.GetDefaultTTL(ctx))
	if err != nil {
		o.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return overtime, nil
}

func (o *overtime) GetList(ctx context.Context, param entity.OvertimeParam) ([]entity.Overtime, *entity.Pagination, error) {
	if !param.BypassCache {
		overtimeList, pg, err := o.getCacheList(ctx, param)
		switch {
		case errors.Is(err, redis.Nil):
			o.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedisNil, err.Error()))
		case err != nil:
			o.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
		default:
			return overtimeList, &pg, nil
		}
	}

	overtimeList, pg, err := o.getListSQL(ctx, param)
	if err != nil {
		return overtimeList, pg, err
	}

	err = o.upsertCacheList(ctx, param, overtimeList, *pg, o.redis.GetDefaultTTL(ctx))
	if err != nil {
		o.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return overtimeList, pg, nil
}

func (o *overtime) Create(ctx context.Context, param entity.OvertimeInputParam) (entity.Overtime, error) {
	overtime, err := o.createSQL(ctx, param)
	if err != nil {
		return overtime, err
	}

	err = o.deleteCache(ctx, deleteOvertimeKeysPattern)
	if err != nil {
		o.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return overtime, nil
}

func (o *overtime) Update(ctx context.Context, updateParam entity.OvertimeUpdateParam, selectParam entity.OvertimeParam) error {
	err := o.updateSQL(ctx, updateParam, selectParam)
	if err != nil {
		return err
	}

	err = o.deleteCache(ctx, deleteOvertimeKeysPattern)
	if err != nil {
		o.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return nil
}
