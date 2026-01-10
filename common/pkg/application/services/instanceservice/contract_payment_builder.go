package instanceservice

import (
	"common/pkg/domain/entities/accounting"
	"common/pkg/domain/entities/contract"
	"common/pkg/domain/entities/events"
	"common/pkg/domain/repositories"
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

	transactionCommandRepository repositories.LedgerTransactionCommandRepository
	transactionQueryRepository   repositories.LedgerTransactionQueryRepository
}

func NewContractPaymentBuilder(
	accountQueryService accountQueryService,
	hardwareCostCalculationService hardwareCostCalculationService,
	transactionCommandRepository repositories.LedgerTransactionCommandRepository,
	transactionQueryRepository repositories.LedgerTransactionQueryRepository,
) *ContractPaymentBuilder {
	return &ContractPaymentBuilder{
		accountQueryService:            accountQueryService,
		hardwareCostCalculationService: hardwareCostCalculationService,
		transactionCommandRepository:   transactionCommandRepository,
		transactionQueryRepository:     transactionQueryRepository,
	}
}

func (builder *ContractPaymentBuilder) calculateContractHardwareCost(ctx context.Context, contract *contract.InstanceRentalContract) (int64, error) {
	return builder.hardwareCostCalculationService.CalculateCost(
		ctx,
		contract.HardwareSpecification(),
		contract.PeriodEnd().Sub(contract.PeriodStart()),
		time.Now(),
	)
}

func (builder *ContractPaymentBuilder) CreateContractPaymentLedgerTransaction(
	ctx context.Context,
	contract *contract.InstanceRentalContract,
) (events.DomainEventEnvelope[accounting.LedgerTransactionCreatedEvent], error) {
	if !contract.IsActiveAt(time.Now()) {
		return events.DomainEventEnvelope[accounting.LedgerTransactionCreatedEvent]{}, nil
	}

	purchaseAmount, err := builder.calculateContractHardwareCost(ctx, contract)
	paymentAccountID, err := builder.accountQueryService.GetPlatformHardwareRentAccountID(ctx)
	if err != nil {
		return events.DomainEventEnvelope[accounting.LedgerTransactionCreatedEvent]{}, err
	}
	userAccountID, err := builder.accountQueryService.GetUserAccountID(ctx, contract.OwnerID())
	if err != nil {
		return events.DomainEventEnvelope[accounting.LedgerTransactionCreatedEvent]{}, err
	}

	entries, err := accounting.NewLedgerEntries(
		accounting.LedgerEntryShorthand(userAccountID, -purchaseAmount),
		accounting.LedgerEntryShorthand(paymentAccountID, purchaseAmount),
	)

	if err != nil {
		return events.DomainEventEnvelope[accounting.LedgerTransactionCreatedEvent]{}, err
	}

	purchaseTransaction, event, err := accounting.NewLedgerTransaction(accounting.ReferenceTypePurchase, string(contract.ID()), entries)
	if err != nil {
		return events.DomainEventEnvelope[accounting.LedgerTransactionCreatedEvent]{}, err
	}
	err = builder.transactionCommandRepository.Create(ctx, purchaseTransaction)
	return event, err
}

func (builder *ContractPaymentBuilder) CreateContractRefundLedgerTransaction(
	ctx context.Context,
	contract *contract.InstanceRentalContract,
) (events.DomainEventEnvelope[accounting.LedgerTransactionCreatedEvent], error) {
	if !contract.IsActiveAt(time.Now()) {
		return events.DomainEventEnvelope[accounting.LedgerTransactionCreatedEvent]{}, nil
	}

	remainingHardwareCost, err := builder.calculateContractHardwareCost(ctx, contract)
	refundAccountID, err := builder.accountQueryService.GetPlatformHardwareRefundAccountID(ctx)
	if err != nil {
		return events.DomainEventEnvelope[accounting.LedgerTransactionCreatedEvent]{}, err
	}
	userAccountID, err := builder.accountQueryService.GetUserAccountID(ctx, contract.OwnerID())
	if err != nil {
		return events.DomainEventEnvelope[accounting.LedgerTransactionCreatedEvent]{}, err
	}

	entries, err := accounting.NewLedgerEntries(
		accounting.LedgerEntryShorthand(userAccountID, remainingHardwareCost),
		accounting.LedgerEntryShorthand(refundAccountID, -remainingHardwareCost),
	)

	if err != nil {
		return events.DomainEventEnvelope[accounting.LedgerTransactionCreatedEvent]{}, err
	}

	purchaseTransaction, event, err := accounting.NewLedgerTransaction(accounting.ReferenceTypePurchase, string(contract.ID()), entries)
	if err != nil {
		return events.DomainEventEnvelope[accounting.LedgerTransactionCreatedEvent]{}, err
	}
	err = builder.transactionCommandRepository.Create(ctx, purchaseTransaction)
	return event, err
}
