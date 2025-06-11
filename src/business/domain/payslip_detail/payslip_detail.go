package payslip_detail

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
	Get(ctx context.Context, param entity.PayslipDetailParam) (entity.PayslipDetail, error)
	GetList(ctx context.Context, param entity.PayslipDetailParam) ([]entity.PayslipDetail, *entity.Pagination, error)
	Create(ctx context.Context, param entity.PayslipDetailInputParam) (entity.PayslipDetail, error)
	CreateMany(ctx context.Context, inputParams []entity.PayslipDetailInputParam) error
	Update(ctx context.Context, updateParam entity.PayslipDetailUpdateParam, selectParam entity.PayslipDetailParam) error
}

type payslipDetail struct {
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
	return &payslipDetail{
		db:    param.Db,
		log:   param.Log,
		redis: param.Redis,
		json:  param.Json,
	}
}

func (p *payslipDetail) Get(ctx context.Context, param entity.PayslipDetailParam) (entity.PayslipDetail, error) {
	payslipDetail := entity.PayslipDetail{}

	marshalledParam, err := p.json.Marshal(param)
	if err != nil {
		return payslipDetail, err
	}

	if !param.BypassCache {
		payslipDetail, err = p.getCache(ctx, fmt.Sprintf(getPayslipDetailByKey, string(marshalledParam)))
		switch {
		case errors.Is(err, redis.Nil):
			p.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedisNil, err.Error()))
		case err != nil:
			p.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
		default:
			return payslipDetail, nil
		}
	}

	payslipDetail, err = p.getSQL(ctx, param)
	if err != nil {
		return payslipDetail, err
	}

	err = p.upsertCache(ctx, fmt.Sprintf(getPayslipDetailByKey, string(marshalledParam)), payslipDetail, p.redis.GetDefaultTTL(ctx))
	if err != nil {
		p.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return payslipDetail, nil
}

func (p *payslipDetail) GetList(ctx context.Context, param entity.PayslipDetailParam) ([]entity.PayslipDetail, *entity.Pagination, error) {
	if !param.BypassCache {
		payslipDetailList, pg, err := p.getCacheList(ctx, param)
		switch {
		case errors.Is(err, redis.Nil):
			p.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedisNil, err.Error()))
		case err != nil:
			p.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
		default:
			return payslipDetailList, &pg, nil
		}
	}

	payslipDetailList, pg, err := p.getListSQL(ctx, param)
	if err != nil {
		return payslipDetailList, pg, err
	}

	err = p.upsertCacheList(ctx, param, payslipDetailList, *pg, p.redis.GetDefaultTTL(ctx))
	if err != nil {
		p.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return payslipDetailList, pg, nil
}

func (p *payslipDetail) Create(ctx context.Context, param entity.PayslipDetailInputParam) (entity.PayslipDetail, error) {
	payslipDetail, err := p.createSQL(ctx, param)
	if err != nil {
		return payslipDetail, err
	}

	err = p.deleteCache(ctx, deletePayslipDetailKeysPattern)
	if err != nil {
		p.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return payslipDetail, nil
}

func (p *payslipDetail) CreateMany(ctx context.Context, inputParams []entity.PayslipDetailInputParam) error {
	err := p.createManySQL(ctx, inputParams)
	if err != nil {
		return err
	}

	err = p.deleteCache(ctx, deletePayslipDetailKeysPattern)
	if err != nil {
		p.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return nil
}

func (p *payslipDetail) Update(ctx context.Context, updateParam entity.PayslipDetailUpdateParam, selectParam entity.PayslipDetailParam) error {
	err := p.updateSQL(ctx, updateParam, selectParam)
	if err != nil {
		return err
	}

	err = p.deleteCache(ctx, deletePayslipDetailKeysPattern)
	if err != nil {
		p.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return nil
}
