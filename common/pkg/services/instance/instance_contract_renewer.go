package instance

import (
	"common/pkg/domain/instance"
	"common/pkg/domain/ledger"
	"common/pkg/shared"
	"common/pkg/util"
	"context"
	"fmt"
	"time"
)

type ContractPaymentCreationService interface {
	CreateContractPaymentLedgerTransaction(ctx context.Context, context *instance.InstanceRentalContract, transaction shared.Transaction) error
}

type InstanceContractRenewer struct {
	transactionProvider            shared.TransactionProvider
	contractPaymentCreationService ContractPaymentCreationService

	contractQueryRepository   instance.InstanceRentalContractQueryRepository
	contractCommandRepository instance.InstanceRentalContractCommandRepository

	ledgerRepository ledger.LedgerTransactionCommandRepository
	tickerWorker     *util.TickerWorker
}

func NewInstanceSubscriptionRenewalReconciler(
	instanceRepository instance.InstanceRentalContractQueryRepository,
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
		request shared.BatchRequest[instance.InstanceRentalContractCursor],
	) ([]*instance.InstanceRentalContract, instance.InstanceRentalContractCursor, error) {
		return r.contractQueryRepository.FetchNextBatchPendingRenewal(ctx, time.Now(), request)
	}

	for contract := range util.GenerateInBatches(ctx, 100, fetchNextBatch) {
		if err := shared.WithTransaction(ctx, r.transactionProvider, func(ctx context.Context, transaction shared.Transaction) error {
			contract.Cancel(time.Now())

			if err := r.contractCommandRepository.Update(ctx, transaction, contract); err != nil {
				return err
			}

			newContract, err := contract.Renew(time.Now())
			if err != nil {
				return err
			}

			if err := r.contractCommandRepository.Create(ctx, transaction, newContract); err != nil {
				return err
			}

			return r.contractPaymentCreationService.CreateContractPaymentLedgerTransaction(ctx, newContract, transaction)
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
