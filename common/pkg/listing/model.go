package listing

import (
	"common/pkg/currency"
	"common/pkg/hardware"
)

type PricingType int
type ListingAccessMode int

const (
	OneTime PricingType = iota
	Monthly
	Yearly
)

const (
	Private ListingAccessMode = iota
	Public
)

type Pricing struct {
	id     string
	typ    PricingType
	amount currency.Money
}

type Listing struct {
	id                    string
	authorID              string
	title                 string
	description           string
	accessMode            ListingAccessMode
	pricing               []*Pricing
	hardwareSpecification *hardware.HardwareSpecification
}
