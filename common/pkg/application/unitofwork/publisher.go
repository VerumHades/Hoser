package unitofwork

import (
	"common/pkg/application/outbox"
	"common/pkg/domain/entities/events"
	"context"
)

type TransactionalEventPublisher struct {
	transactionProvider TransactionProvider
	eventPublisher      outbox.EventPublisher
}

func NewTransactionalEventPublisher(
	transactionProvider TransactionProvider,
	eventPublisher outbox.EventPublisher,
) *TransactionalEventPublisher {
	return &TransactionalEventPublisher{
		transactionProvider: transactionProvider,
		eventPublisher:      eventPublisher,
	}
}

type TransactionalEventFunction func(txContext context.Context) ([]events.DomainEvent, error)

/*
WithTransactionalEvents runs the given function within a transaction.
All events returned by the function will be persisted automatically on success.
Rollback occurs automatically on error or panic.
*/
func (executor *TransactionalEventPublisher) PublishWithTransaction(
	ctx context.Context,
	fn TransactionalEventFunction,
) error {
	return WithTransaction(ctx, executor.transactionProvider, func(txCtx context.Context) error {
		events, err := fn(txCtx)
		if err != nil {
			return err
		}
		if len(events) > 0 {
			return executor.eventPublisher.PublishEnvelopes(txCtx, events...)
		}
		return nil
	})
}
