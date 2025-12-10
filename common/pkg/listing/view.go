package listing

import (
	"common/pkg/currency"
	"common/pkg/hardware"
)

// PricingView is a read-only representation of a Pricing entry.
type PricingView struct {
	ID     string
	Type   PricingType
	Amount currency.Money
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

// ToView converts a Listing to a read-only ListingView.
func (l *Listing) ToView() *ListingView {
	pricingViews := make([]*PricingView, len(l.pricing))
	for i, p := range l.pricing {
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
