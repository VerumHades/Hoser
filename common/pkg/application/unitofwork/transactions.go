package unitofwork

import (
	"context"
)

type TransactionProvider interface {
	BeginTransaction(ctx context.Context) (txContext context.Context, err error)
}

type TransactionFunction func(txContext context.Context) error

func WithTransaction(
	ctx context.Context,
	factory TransactionProvider,
	transactionFunction TransactionFunction,
) error {
	transactionContext, err := factory.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	return transactionFunction(transactionContext)
}
