package payslip

import (
	"context"
	"fmt"
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichies/employee-payroll-service/src/business/entity"
)

const (
	getPayslipByKey           = "employeePayroll:payslip:get:%s"
	getPayslipByQueryKey      = "employeePayroll:payslip:get:q:%s"
	getPayslipByPaginationKey = "employeePayroll:payslip:get:p:%s"
	deletePayslipKeysPattern  = "employeePayroll:payslip*"
)

func (p *payslip) upsertCache(ctx context.Context, key string, payslip entity.Payslip, ttl time.Duration) error {
	marshalledPayslip, err := p.json.Marshal(payslip)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	err = p.redis.SetEX(ctx, key, string(marshalledPayslip), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheSetSimpleKey, err.Error())
	}

	return nil
}

func (p *payslip) getCache(ctx context.Context, key string) (entity.Payslip, error) {
	payslip := entity.Payslip{}

	marshalledPayslip, err := p.redis.Get(ctx, key)
	if err != nil {
		return payslip, err
	}

	err = p.json.Unmarshal([]byte(marshalledPayslip), &payslip)
	if err != nil {
		return payslip, errors.NewWithCode(codes.CodeCacheUnmarshal, err.Error())
	}

	return payslip, nil
}

func (p *payslip) upsertCacheList(ctx context.Context, param entity.PayslipParam, payslipList []entity.Payslip, pg entity.Pagination, ttl time.Duration) error {
	keyValue, err := p.json.Marshal(param)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	// Set payslip list to cache
	marshalledPayslipList, err := p.json.Marshal(payslipList)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}
	err = p.redis.SetEX(ctx, fmt.Sprintf(getPayslipByQueryKey, string(keyValue)), string(marshalledPayslipList), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheSetSimpleKey, err.Error())
	}

	// Set pagination to cache
	marshalledPagination, err := p.json.Marshal(pg)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	err = p.redis.SetEX(ctx, fmt.Sprintf(getPayslipByPaginationKey, string(keyValue)), string(marshalledPagination), ttl)
	if err != nil {
		return errors.NewWithCode(codes.CodeCacheSetSimpleKey, err.Error())
	}

	return nil
}

func (p *payslip) getCacheList(ctx context.Context, param entity.PayslipParam) ([]entity.Payslip, entity.Pagination, error) {
	var (
		payslipList = []entity.Payslip{}
		pg          = entity.Pagination{}
	)

	keyValue, err := p.json.Marshal(param)
	if err != nil {
		return payslipList, pg, errors.NewWithCode(codes.CodeCacheMarshal, err.Error())
	}

	// Get payslip list from redis
	marshalledPayslipList, err := p.redis.Get(ctx, fmt.Sprintf(getPayslipByQueryKey, string(keyValue)))
	if err != nil {
		return payslipList, pg, err
	}

	err = p.json.Unmarshal([]byte(marshalledPayslipList), &payslipList)
	if err != nil {
		return payslipList, pg, errors.NewWithCode(codes.CodeCacheUnmarshal, err.Error())
	}

	// Get pagination from redis
	marshalledPagination, err := p.redis.Get(ctx, fmt.Sprintf(getPayslipByPaginationKey, string(keyValue)))
	if err != nil {
		return payslipList, pg, err
	}

	err = p.json.Unmarshal([]byte(marshalledPagination), &pg)
	if err != nil {
		return payslipList, pg, errors.NewWithCode(codes.CodeCacheUnmarshal, err.Error())
	}

	return payslipList, pg, nil
}

func (p *payslip) deleteCache(ctx context.Context, key string) error {
	err := p.redis.Del(ctx, key)
	if err != nil {
		return err
	}

	return nil
}
