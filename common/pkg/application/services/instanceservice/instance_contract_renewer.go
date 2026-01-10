package instanceservice

import (
	"common/pkg/application/unitofwork"
	"common/pkg/domain/entities/accounting"
	"common/pkg/domain/entities/contract"
	"common/pkg/domain/entities/events"
	"common/pkg/domain/repositories"
	"common/pkg/shared"
	"common/pkg/util"
	"context"
	"fmt"
	"time"
)

type ContractPaymentCreationService interface {
	CreateContractPaymentLedgerTransaction(ctx context.Context, context *contract.InstanceRentalContract) (events.DomainEventEnvelope[accounting.LedgerTransactionCreatedEvent], error)
}

type InstanceContractRenewer struct {
	transactionalPublisher         *unitofwork.TransactionalEventPublisher
	contractPaymentCreationService ContractPaymentCreationService

	contractQueryRepository   repositories.InstanceRentalContractQueryRepository
	contractCommandRepository repositories.InstanceRentalContractCommandRepository

	tickerWorker *util.TickerWorker
}

func NewInstanceSubscriptionRenewalReconciler(
	instanceRepository repositories.InstanceRentalContractQueryRepository,
	interval time.Duration,
) *InstanceContractRenewer {
	reconciler := &InstanceContractRenewer{
		contractQueryRepository: instanceRepository,
	}

	reconciler.tickerWorker = util.NewTickerWorker(
		interval,
		reconciler.tick,
		func(err error) {
			fmt.Println("instance renewal reconciliation error:", err)
		},
	)

	return reconciler
}

func (r *InstanceContractRenewer) tick(ctx context.Context) error {
	fetchNextBatch := func(
		ctx context.Context,
		request shared.BatchRequest[repositories.InstanceRentalContractCursor],
	) ([]*contract.InstanceRentalContract, repositories.InstanceRentalContractCursor, error) {
		return r.contractQueryRepository.FetchNextBatchPendingRenewal(ctx, time.Now(), request)
	}

	for contract := range util.GenerateInBatches(ctx, 100, fetchNextBatch) {
		if err := r.transactionalPublisher.PublishWithTransaction(ctx, func(txContext context.Context) ([]events.DomainEvent, error) {
			contract.Cancel(time.Now())

			if err := r.contractCommandRepository.Update(txContext, contract); err != nil {
				return []events.DomainEvent{}, err
			}

			newContract, err := contract.Renew(time.Now())
			if err != nil {
				return []events.DomainEvent{}, err
			}

			if err := r.contractCommandRepository.Create(txContext, newContract); err != nil {
				return []events.DomainEvent{}, err
			}

			event, err := r.contractPaymentCreationService.CreateContractPaymentLedgerTransaction(txContext, newContract)
			return []events.DomainEvent{event}, err

		}); err != nil {
			return err
		}
	}

	return nil
}

func (r *InstanceContractRenewer) Start(ctx context.Context) {
	r.tickerWorker.Start(ctx)
}

func (r *InstanceContractRenewer) Stop() {
	r.tickerWorker.Stop()
}
