package purchase

import (
	"common/pkg/currency"
	"common/pkg/util"
)

// PurchaseService handles business logic related to purchases.
type PurchaseService struct {
	repo PurchaseRepository
}

// NewPurchaseService creates a new purchase service instance.
func NewPurchaseService(repo PurchaseRepository) *PurchaseService {
	return &PurchaseService{repo: repo}
}

// CreatePurchase creates and persists a new purchase.
func (s *PurchaseService) CreatePurchase(userID, listingID string, price currency.Money) (*Purchase, error) {
	p := &Purchase{
		id:        util.GenerateUUID(),
		userID:    userID,
		listingID: listingID,
		pricePaid: price,
	}
	err := s.repo.Save(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PurchaseService) HasPurchased(userID, listingID string) (bool, error) {
	return s.repo.ExistsByUserAndListing(userID, listingID)
}

// GetPurchaseView retrieves a purchase as a read-only view.
func (s *PurchaseService) GetPurchaseView(purchaseID string) (*PurchaseView, error) {
	p, err := s.repo.GetByID(purchaseID)
	if err != nil {
		return nil, err
	}
	return p.ToView(), nil
}
