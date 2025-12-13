package listing

import (
	"common/pkg/hardware"
	"common/pkg/money"
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
	amount money.Money
}

type Listing struct {
	id                    string
	authorID              string
	title                 string
	description           string
	accessMode            ListingAccessMode
	hardwareSpecification *hardware.HardwareSpecification
}
