package transactor

import (
	"context"

	"github.com/reyhanmichiels/go-pkg/v2/sql"
)

type Interface interface {
	Execute(ctx context.Context, name string, txOpts sql.TxOptions, f func(context.Context) error) error
}

type transactor struct {
	db sql.Interface
}

func Init(db sql.Interface) Interface {
	return &transactor{
		db: db,
	}
}

func (t *transactor) Execute(ctx context.Context, name string, txOpts sql.TxOptions, f func(context.Context) error) error {
	return t.db.Transaction(ctx, name, txOpts, f)
}
