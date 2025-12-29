package shared

import "context"

type Transaction interface {
	Rollback() error
	Commit() error
}

type TransactionProvider interface {
	BeginTransaction(ctx context.Context) (Transaction, error)
}

type TransactionFunction func(transaction Transaction) error

func WithTransaction(
	ctx context.Context,
	factory TransactionProvider,
	transactionFunction TransactionFunction,
) error {

	transaction, err := factory.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if recoveredValue := recover(); recoveredValue != nil {
			_ = transaction.Rollback()
			panic(recoveredValue)
		}
	}()

	if err := transactionFunction(transaction); err != nil {
		_ = transaction.Rollback()
		return err
	}

	return transaction.Commit()
}
