package payslip_detail

import (
	"context"
	"fmt"
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
)

const (
	getPayslipDetailByKey           = "employeePayroll:payslipDetail:get:%s"
	getPayslipDetailByQueryKey      = "employeePayroll:payslipDetail:get:q:%s"
	getPayslipDetailByPaginationKey = "employeePayroll:payslipDetail:get:p:%s"
	deletePayslipDetailKeysPattern  = "employeePayroll:payslipDetail*"
)

func (p *payslipDetail) upsertCache(ctx context.Context, key string, payslipDetail entity.PayslipDetail, ttl time.Duration) error {
	marshalledPayslipDetail, err := p.json.Marshal(payslipDetail)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	err = p.redis.SetEX(ctx, key, string(marshalledPayslipDetail), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheSetSimpleKey, err.Error())
	}

	return nil
}

func (p *payslipDetail) getCache(ctx context.Context, key string) (entity.PayslipDetail, error) {
	payslipDetail := entity.PayslipDetail{}

	marshalledPayslipDetail, err := p.redis.Get(ctx, key)
	if err != nil {
		return payslipDetail, err
	}

	err = p.json.Unmarshal([]byte(marshalledPayslipDetail), &payslipDetail)
	if err != nil {
		return payslipDetail, errors.NewWithCode(codes.CodeCacheUnmarshal, err.Error())
	}

	return payslipDetail, nil
}

func (p *payslipDetail) upsertCacheList(ctx context.Context, param entity.PayslipDetailParam, payslipDetailList []entity.PayslipDetail, pg entity.Pagination, ttl time.Duration) error {
	keyValue, err := p.json.Marshal(param)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	// Set payslipDetail list to cache
	marshalledPayslipDetailList, err := p.json.Marshal(payslipDetailList)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}
	err = p.redis.SetEX(ctx, fmt.Sprintf(getPayslipDetailByQueryKey, string(keyValue)), string(marshalledPayslipDetailList), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheSetSimpleKey, err.Error())
	}

	// Set pagination to cache
	marshalledPagination, err := p.json.Marshal(pg)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	err = p.redis.SetEX(ctx, fmt.Sprintf(getPayslipDetailByPaginationKey, string(keyValue)), string(marshalledPagination), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheSetSimpleKey, err.Error())
	}

	return nil
}

func (p *payslipDetail) getCacheList(ctx context.Context, param entity.PayslipDetailParam) ([]entity.PayslipDetail, entity.Pagination, error) {
	var (
		payslipDetailList = []entity.PayslipDetail{}
		pg                = entity.Pagination{}
	)

	keyValue, err := p.json.Marshal(param)
	if err != nil {
		return payslipDetailList, pg, errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	// Get payslipDetail list from redis
	marshalledPayslipDetailList, err := p.redis.Get(ctx, fmt.Sprintf(getPayslipDetailByQueryKey, string(keyValue)))
	if err != nil {
		return payslipDetailList, pg, err
	}

	err = p.json.Unmarshal([]byte(marshalledPayslipDetailList), &payslipDetailList)
	if err != nil {
		return payslipDetailList, pg, errors.NewWithCode(codes.CodeCacheUnmarshal, err.Error())
	}

	// Get pagination from redis
	marshalledPagination, err := p.redis.Get(ctx, fmt.Sprintf(getPayslipDetailByPaginationKey, string(keyValue)))
	if err != nil {
		return payslipDetailList, pg, err
	}

	err = p.json.Unmarshal([]byte(marshalledPagination), &pg)
	if err != nil {
		return payslipDetailList, pg, errors.NewWithCode(codes.CodeCacheUnmarshal, err.Error())
	}

	return payslipDetailList, pg, nil
}

func (p *payslipDetail) deleteCache(ctx context.Context, key string) error {
	err := p.redis.Del(ctx, key)
	if err != nil {
		return err
	}

	return nil
}
