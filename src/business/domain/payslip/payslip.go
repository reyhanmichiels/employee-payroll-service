package payslip

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
	Get(ctx context.Context, param entity.PayslipParam) (entity.Payslip, error)
	GetList(ctx context.Context, param entity.PayslipParam) ([]entity.Payslip, *entity.Pagination, error)
	Create(ctx context.Context, param entity.PayslipInputParam) (entity.Payslip, error)
	CreateMany(ctx context.Context, inputParams []entity.PayslipInputParam) error
	Update(ctx context.Context, updateParam entity.PayslipUpdateParam, selectParam entity.PayslipParam) error
}

type payslip struct {
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
	return &payslip{
		db:    param.Db,
		log:   param.Log,
		redis: param.Redis,
		json:  param.Json,
	}
}

func (p *payslip) Get(ctx context.Context, param entity.PayslipParam) (entity.Payslip, error) {
	payslip := entity.Payslip{}

	marshalledParam, err := p.json.Marshal(param)
	if err != nil {
		return payslip, err
	}

	if !param.BypassCache {
		payslip, err = p.getCache(ctx, fmt.Sprintf(getPayslipByKey, string(marshalledParam)))
		switch {
		case errors.Is(err, redis.Nil):
			p.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedisNil, err.Error()))
		case err != nil:
			p.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
		default:
			return payslip, nil
		}
	}

	payslip, err = p.getSQL(ctx, param)
	if err != nil {
		return payslip, err
	}

	err = p.upsertCache(ctx, fmt.Sprintf(getPayslipByKey, string(marshalledParam)), payslip, p.redis.GetDefaultTTL(ctx))
	if err != nil {
		p.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return payslip, nil
}

func (p *payslip) GetList(ctx context.Context, param entity.PayslipParam) ([]entity.Payslip, *entity.Pagination, error) {
	if !param.BypassCache {
		payslipList, pg, err := p.getCacheList(ctx, param)
		switch {
		case errors.Is(err, redis.Nil):
			p.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedisNil, err.Error()))
		case err != nil:
			p.log.Warn(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
		default:
			return payslipList, &pg, nil
		}
	}

	payslipList, pg, err := p.getListSQL(ctx, param)
	if err != nil {
		return payslipList, pg, err
	}

	err = p.upsertCacheList(ctx, param, payslipList, *pg, p.redis.GetDefaultTTL(ctx))
	if err != nil {
		p.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return payslipList, pg, nil
}

func (p *payslip) Create(ctx context.Context, param entity.PayslipInputParam) (entity.Payslip, error) {
	payslip, err := p.createSQL(ctx, param)
	if err != nil {
		return payslip, err
	}

	err = p.deleteCache(ctx, deletePayslipKeysPattern)
	if err != nil {
		p.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return payslip, nil
}

func (p *payslip) CreateMany(ctx context.Context, inputParams []entity.PayslipInputParam) error {
	err := p.createManySQL(ctx, inputParams)
	if err != nil {
		return err
	}

	err = p.deleteCache(ctx, deletePayslipKeysPattern)
	if err != nil {
		p.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return nil
}

func (p *payslip) Update(ctx context.Context, updateParam entity.PayslipUpdateParam, selectParam entity.PayslipParam) error {
	err := p.updateSQL(ctx, updateParam, selectParam)
	if err != nil {
		return err
	}

	err = p.deleteCache(ctx, deletePayslipKeysPattern)
	if err != nil {
		p.log.Error(ctx, fmt.Sprintf(entity.ErrorRedis, err.Error()))
	}

	return nil
}
