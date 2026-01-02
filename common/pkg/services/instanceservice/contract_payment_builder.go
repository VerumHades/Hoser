package instanceservice

import (
	"common/pkg/domain/instance"
	"common/pkg/domain/ledger"
	"common/pkg/shared"
	"context"
	"time"
)

type hardwareCostCalculationService interface {
	CalculateCost(
		ctx context.Context,
		spec *shared.HardwareSpecification,
		duration time.Duration,
		at time.Time,
	) (int64, error)
}
type accountQueryService interface {
	GetPlatformHardwareRentAccountID(ctx context.Context) (shared.AccountID, error)
	GetPlatformHardwareRefundAccountID(ctx context.Context) (shared.AccountID, error)
	GetUserAccountID(ctx context.Context, userID shared.UserID) (shared.AccountID, error)
}

type ContractPaymentBuilder struct {
	accountQueryService            accountQueryService
	hardwareCostCalculationService hardwareCostCalculationService

	transactionCommandRepository ledger.LedgerTransactionCommandRepository
	transactionQueryRepository   ledger.LedgerTransactionQueryRepository
}

func NewContractPaymentBuilder(
	accountQueryService accountQueryService,
	hardwareCostCalculationService hardwareCostCalculationService,
	transactionCommandRepository ledger.LedgerTransactionCommandRepository,
	transactionQueryRepository ledger.LedgerTransactionQueryRepository,
) *ContractPaymentBuilder {
	return &ContractPaymentBuilder{
		accountQueryService:            accountQueryService,
		hardwareCostCalculationService: hardwareCostCalculationService,
		transactionCommandRepository:   transactionCommandRepository,
		transactionQueryRepository:     transactionQueryRepository,
	}
}

func (builder *ContractPaymentBuilder) calculateContractHardwareCost(ctx context.Context, contract *instance.InstanceRentalContract) (int64, error) {
	return builder.hardwareCostCalculationService.CalculateCost(
		ctx,
		contract.HardwareSpecification(),
		contract.PeriodEnd().Sub(contract.PeriodStart()),
		time.Now(),
	)
}

func (builder *ContractPaymentBuilder) CreateContractPaymentLedgerTransaction(
	ctx context.Context,
	contract *instance.InstanceRentalContract,
	transaction shared.Transaction,
) error {
	if !contract.IsActiveAt(time.Now()) {
		return nil
	}

	purchaseAmount, err := builder.calculateContractHardwareCost(ctx, contract)
	paymentAccountID, err := builder.accountQueryService.GetPlatformHardwareRentAccountID(ctx)
	if err != nil {
		return err
	}
	userAccountID, err := builder.accountQueryService.GetUserAccountID(ctx, contract.OwnerID())
	if err != nil {
		return err
	}

	entries, err := ledger.NewLedgerEntries(
		ledger.LedgerEntryShorthand(userAccountID, -purchaseAmount),
		ledger.LedgerEntryShorthand(paymentAccountID, purchaseAmount),
	)

	if err != nil {
		return err
	}

	purchaseTransaction, err := ledger.NewLedgerTransaction(ledger.ReferenceTypePurchase, string(contract.ID()), entries)
	if err != nil {
		return err
	}
	return builder.transactionCommandRepository.Create(ctx, transaction, purchaseTransaction)
}

func (builder *ContractPaymentBuilder) CreateContractRefundLedgerTransaction(
	ctx context.Context,
	contract *instance.InstanceRentalContract,
	transaction shared.Transaction,
) error {
	if !contract.IsActiveAt(time.Now()) {
		return nil
	}

	remainingHardwareCost, err := builder.calculateContractHardwareCost(ctx, contract)
	refundAccountID, err := builder.accountQueryService.GetPlatformHardwareRefundAccountID(ctx)
	if err != nil {
		return err
	}
	userAccountID, err := builder.accountQueryService.GetUserAccountID(ctx, contract.OwnerID())
	if err != nil {
		return err
	}

	entries, err := ledger.NewLedgerEntries(
		ledger.LedgerEntryShorthand(userAccountID, remainingHardwareCost),
		ledger.LedgerEntryShorthand(refundAccountID, -remainingHardwareCost),
	)

	if err != nil {
		return err
	}

	purchaseTransaction, err := ledger.NewLedgerTransaction(ledger.ReferenceTypePurchase, string(contract.ID()), entries)
	if err != nil {
		return err
	}
	return builder.transactionCommandRepository.Create(ctx, transaction, purchaseTransaction)
}
