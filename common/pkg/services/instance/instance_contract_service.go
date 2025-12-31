package instance

import (
	"common/pkg/domain/instance"
	"common/pkg/domain/listing"
	"common/pkg/domain/user"
	"common/pkg/shared"
	"common/pkg/util"
	"context"
	"fmt"
	"time"
)

type ContractPaymentService interface {
	CreateContractPaymentLedgerTransaction(ctx context.Context, context *instance.InstanceRentalContract, transaction shared.Transaction) error
	CreateContractRefundLedgerTransaction(ctx context.Context, context *instance.InstanceRentalContract, transaction shared.Transaction) error
}

type InstanceContractService struct {
	paymentService      ContractPaymentService
	transactionProvider shared.TransactionProvider

	listingQueryRepository    listing.ListingQueryRepository
	contractQueryRepository   instance.InstanceRentalContractQueryRepository
	contractCommandRepository instance.InstanceRentalContractCommandRepository
	userQueryRepository       user.UserQueryRepository
}

func (service *InstanceContractService) withContract(
	ctx context.Context,
	contractID shared.InstanceRentalContractID,
	processFunction func(ctx context.Context, contract *instance.InstanceRentalContract) error,
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
	contract *instance.InstanceRentalContract,
) error {
	if err := util.EnsureEntityExists(ctx, userID, fmt.Errorf("user %s does not exist", userID), service.userQueryRepository); err != nil {
		return err
	}

	if err := util.EnsureEntityExists(ctx, listingID, fmt.Errorf("listing %s does not exist", listingID), service.listingQueryRepository); err != nil {
		return err
	}

	return service.paymentService.CreateContractPaymentLedgerTransaction(ctx, contract, nil)
}

func (service *InstanceContractService) applyToContract(
	ctx context.Context,
	contractID shared.InstanceRentalContractID,
	updater func(contract *instance.InstanceRentalContract),
) error {
	return service.withContract(ctx, contractID, func(ctx context.Context, contract *instance.InstanceRentalContract) error {
		updater(contract)
		return service.contractCommandRepository.Update(ctx, nil, contract)
	})
}

func (service *InstanceContractService) EnableContractRenewal(
	ctx context.Context,
	contractID shared.InstanceRentalContractID,
	renewDuration time.Duration,
) error {
	return service.applyToContract(ctx, contractID, func(contract *instance.InstanceRentalContract) {
		contract.EnableRenewal(renewDuration)
	})
}

func (service *InstanceContractService) DisableContractRenewal(
	ctx context.Context,
	contractID shared.InstanceRentalContractID,
) error {
	return service.applyToContract(ctx, contractID, func(contract *instance.InstanceRentalContract) {
		contract.DisableRenewal()
	})
}

func (service *InstanceContractService) ChangeContractHardwareSpecification(
	ctx context.Context,
	contractID shared.InstanceRentalContractID,
	hardwareSpecification *shared.HardwareSpecification,
) error {
	return service.withContract(ctx, contractID, func(ctx context.Context, contract *instance.InstanceRentalContract) error {
		return shared.WithTransaction(ctx, service.transactionProvider, func(ctx context.Context, transaction shared.Transaction) error {
			if err := service.paymentService.CreateContractRefundLedgerTransaction(ctx, contract, transaction); err != nil {
				return err
			}

			newContract, err := contract.WithNewHardwareSpecification(hardwareSpecification)
			if err != nil {
				return err
			}

			if err := service.contractCommandRepository.Create(ctx, transaction, newContract); err != nil {
				return err
			}

			return service.paymentService.CreateContractPaymentLedgerTransaction(ctx, newContract, nil)
		})
	})
}
