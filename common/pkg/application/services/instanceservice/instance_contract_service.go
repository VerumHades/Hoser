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

type ContractPaymentService interface {
	CreateContractPaymentLedgerTransaction(ctx context.Context, context *contract.InstanceRentalContract) (events.DomainEventEnvelope[accounting.LedgerTransactionCreatedEvent], error)
	CreateContractRefundLedgerTransaction(ctx context.Context, context *contract.InstanceRentalContract) (events.DomainEventEnvelope[accounting.LedgerTransactionCreatedEvent], error)
}

type InstanceContractService struct {
	paymentService         ContractPaymentService
	transactionalPublisher *unitofwork.TransactionalEventPublisher

	listingQueryRepository    repositories.ListingQueryRepository
	contractQueryRepository   repositories.InstanceRentalContractQueryRepository
	contractCommandRepository repositories.InstanceRentalContractCommandRepository
	userQueryRepository       repositories.UserQueryRepository
}

func NewInstanceContractService(
	paymentService ContractPaymentService,
	transactionalPublisher *unitofwork.TransactionalEventPublisher,
	listingQueryRepository repositories.ListingQueryRepository,
	contractQueryRepository repositories.InstanceRentalContractQueryRepository,
	contractCommandRepository repositories.InstanceRentalContractCommandRepository,
	userQueryRepository repositories.UserQueryRepository,
) *InstanceContractService {
	return &InstanceContractService{
		paymentService:            paymentService,
		transactionalPublisher:    transactionalPublisher,
		listingQueryRepository:    listingQueryRepository,
		contractQueryRepository:   contractQueryRepository,
		contractCommandRepository: contractCommandRepository,
		userQueryRepository:       userQueryRepository,
	}
}

func (service *InstanceContractService) withContract(
	ctx context.Context,
	contractID shared.InstanceRentalContractID,
	processFunction func(ctx context.Context, contract *contract.InstanceRentalContract) error,
) error {
	contract, err := service.contractQueryRepository.GetByID(ctx, contractID)
	if err != nil {
		return err
	}

	if contract == nil {
		return fmt.Errorf("contract does not exist")
	}

	return processFunction(ctx, contract)
}

func (service *InstanceContractService) RentInstanceOfListing(
	ctx context.Context,
	userID shared.UserID,
	listingID shared.ListingID,
	contract *contract.InstanceRentalContract,
) error {
	if err := util.EnsureEntityExists(ctx, userID, fmt.Errorf("user %s does not exist", userID), service.userQueryRepository); err != nil {
		return err
	}

	if err := util.EnsureEntityExists(ctx, listingID, fmt.Errorf("listing %s does not exist", listingID), service.listingQueryRepository); err != nil {
		return err
	}
	return service.transactionalPublisher.PublishWithTransaction(ctx, func(txContext context.Context) ([]events.DomainEvent, error) {
		event, err := service.paymentService.CreateContractPaymentLedgerTransaction(ctx, contract)
		return []events.DomainEvent{event}, err
	})
}

func (service *InstanceContractService) applyToContract(
	ctx context.Context,
	contractID shared.InstanceRentalContractID,
	updater func(contract *contract.InstanceRentalContract),
) error {
	return service.withContract(ctx, contractID, func(ctx context.Context, contract *contract.InstanceRentalContract) error {
		updater(contract)
		return service.contractCommandRepository.Update(ctx, contract)
	})
}

func (service *InstanceContractService) EnableContractRenewal(
	ctx context.Context,
	contractID shared.InstanceRentalContractID,
	renewDuration time.Duration,
) error {
	return service.applyToContract(ctx, contractID, func(contract *contract.InstanceRentalContract) {
		contract.EnableRenewal(renewDuration)
	})
}

func (service *InstanceContractService) DisableContractRenewal(
	ctx context.Context,
	contractID shared.InstanceRentalContractID,
) error {
	return service.applyToContract(ctx, contractID, func(contract *contract.InstanceRentalContract) {
		contract.DisableRenewal()
	})
}

func (service *InstanceContractService) ChangeContractHardwareSpecification(
	ctx context.Context,
	contractID shared.InstanceRentalContractID,
	hardwareSpecification *shared.HardwareSpecification,
) error {
	return service.withContract(ctx, contractID, func(ctx context.Context, contract *contract.InstanceRentalContract) error {
		return service.transactionalPublisher.PublishWithTransaction(ctx, func(txContext context.Context) ([]events.DomainEvent, error) {
			event, err := service.paymentService.CreateContractRefundLedgerTransaction(txContext, contract)

			if err != nil {
				return []events.DomainEvent{}, err
			}

			newContract, err := contract.WithNewHardwareSpecification(hardwareSpecification)
			if err != nil {
				return []events.DomainEvent{}, err
			}

			if err := service.contractCommandRepository.Create(txContext, newContract); err != nil {
				return []events.DomainEvent{}, err
			}

			event2, err := service.paymentService.CreateContractPaymentLedgerTransaction(txContext, newContract)
			return []events.DomainEvent{event, event2}, err
		})
	})
}
