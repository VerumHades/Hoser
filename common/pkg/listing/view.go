package listing

import (
	"common/pkg/hardware"
	"common/pkg/money"
	"time"
)

// PricingView is a read-only representation of a Pricing entry.
type PricingView struct {
	ID     string
	Type   PricingType
	Amount money.Money
}

// ListingView is a read-only representation of a Listing.
type ListingView struct {
	ID                    string
	AuthorID              string
	Title                 string
	Description           string
	AccessMode            ListingAccessMode
	Pricing               []*PricingView
	HardwareSpecification *hardware.HardwareSpecificationView
}

// ToView converts a Listing and its pricings to a read-only ListingView.
func (l *Listing) ToView(pricings []*Pricing) *ListingView {
	pricingViews := make([]*PricingView, len(pricings))
	for i, p := range pricings {
		pricingViews[i] = &PricingView{
			ID:     p.id,
			Type:   p.typ,
			Amount: p.amount,
		}
	}

	var hwView *hardware.HardwareSpecificationView
	if l.hardwareSpecification != nil {
		hwView = l.hardwareSpecification.ToView()
	}

	return &ListingView{
		ID:                    l.id,
		AuthorID:              l.authorID,
		Title:                 l.title,
		Description:           l.description,
		AccessMode:            l.accessMode,
		Pricing:               pricingViews,
		HardwareSpecification: hwView,
	}
}

type OneTimePaymentMetadataView struct {
	ListingID string
	UserID    string
}

// SubscriptionPaymentMetadataView is a read-only representation of SubscriptionPaymentMetadata.
type SubscriptionPaymentMetadataView struct {
	InstanceID       string
	CurrentPeriodEnd time.Time
}
