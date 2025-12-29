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
	return util.ProcessInBatches(
		ctx,
		100,
		func(ctx context.Context, request shared.BatchRequest) (items []*instance.InstanceRentalContract, nextCursor shared.Cursor, err error) {
			return r.contractQueryRepository.FetchNextBatchPendingRenewal(ctx, time.Now(), request)
		},
		func(ctx context.Context, contract *instance.InstanceRentalContract) error {
			return shared.WithTransaction(ctx, r.transactionProvider, func(transaction shared.Transaction) error {
				contract.Cancel(time.Now())

				err := r.contractCommandRepository.Update(ctx, transaction, contract)
				if err != nil {
					return err
				}

				newContract, err := contract.Renew(time.Now())
				if err != nil {
					return err
				}

				err = r.contractCommandRepository.Create(ctx, transaction, newContract)
				if err != nil {
					return err
				}

				return r.contractPaymentCreationService.CreateContractPaymentLedgerTransaction(ctx, newContract, transaction)
			})
		},
	)
}

func (r *InstanceContractRenewer) Start(ctx context.Context) {
	r.tickerWorker.Start(ctx)
}

func (r *InstanceContractRenewer) Stop() {
	r.tickerWorker.Stop()
}
