package reimbursement

import (
	"context"
	"fmt"

	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/log"
	"github.com/reyhanmichiels/go-pkg/v2/parser"
	"github.com/reyhanmichiels/go-pkg/v2/redis"
	"github.com/reyhanmichiels/go-pkg/v2/sql"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/business/entity"
)

type Interface interface {
	Get(ctx context.Context, param entity.ReimbursementParam) (entity.Reimbursement, error)
	GetList(ctx context.Context, param entity.ReimbursementParam) ([]entity.Reimbursement, *entity.Pagination, error)
	Create(ctx context.Context, param entity.ReimbursementInputParam) (entity.Reimbursement, error)
	Update(ctx context.Context, updateParam entity.ReimbursementUpdateParam, selectParam entity.ReimbursementParam) error
}

type reimbursement struct {
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
	return &reimbursement{
		db:    param.Db,
		log:   param.Log,
		redis: param.Redis,
		json:  param.Json,
	}
}

func (r *reimbursement) Get(ctx context.Context, param entity.ReimbursementParam) (entity.Reimbursement, error) {
	reimbursement := entity.Reimbursement{}

	marshalledParam, err := r.json.Marshal(param)
	if err != nil {
		return reimbursement, err
	}

	if !param.BypassCache {
		reimbursement, err = r.getCache(ctx, fmt.Sprintf(getReimbursementByKey, string(marshalledParam)))
		switch {
		case errors.Is(err, redis.Nil):
			r.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedisNil, err.Error()))
		case err != nil:
			r.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
		default:
			return reimbursement, nil
		}
	}

	reimbursement, err = r.getSQL(ctx, param)
	if err != nil {
		return reimbursement, err
	}

	err = r.upsertCache(ctx, fmt.Sprintf(getReimbursementByKey, string(marshalledParam)), reimbursement, r.redis.GetDefaultTTL(ctx))
	if err != nil {
		r.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return reimbursement, nil
}

func (r *reimbursement) GetList(ctx context.Context, param entity.ReimbursementParam) ([]entity.Reimbursement, *entity.Pagination, error) {
	if !param.BypassCache {
		reimbursementList, pg, err := r.getCacheList(ctx, param)
		switch {
		case errors.Is(err, redis.Nil):
			r.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedisNil, err.Error()))
		case err != nil:
			r.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
		default:
			return reimbursementList, &pg, nil
		}
	}

	reimbursementList, pg, err := r.getListSQL(ctx, param)
	if err != nil {
		return reimbursementList, pg, err
	}

	err = r.upsertCacheList(ctx, param, reimbursementList, *pg, r.redis.GetDefaultTTL(ctx))
	if err != nil {
		r.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return reimbursementList, pg, nil
}

func (r *reimbursement) Create(ctx context.Context, param entity.ReimbursementInputParam) (entity.Reimbursement, error) {
	reimbursement, err := r.createSQL(ctx, param)
	if err != nil {
		return reimbursement, err
	}

	err = r.deleteCache(ctx, deleteReimbursementKeysPattern)
	if err != nil {
		r.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return reimbursement, nil
}

func (r *reimbursement) Update(ctx context.Context, updateParam entity.ReimbursementUpdateParam, selectParam entity.ReimbursementParam) error {
	err := r.updateSQL(ctx, updateParam, selectParam)
	if err != nil {
		return err
	}

	err = r.deleteCache(ctx, deleteReimbursementKeysPattern)
	if err != nil {
		r.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return nil
}
