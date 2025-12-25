package domain

// BillingAccountStatus defines the possible states of a billing account.
type BillingAccountStatus string

const (
	Active    BillingAccountStatus = "active"
	Suspended BillingAccountStatus = "suspended"
	Closed    BillingAccountStatus = "closed"
)

// PaymentKind represents the type of payment.
type PaymentKind string

const (
	OneTime      PaymentKind = "one_time"
	Subscription PaymentKind = "subscription"
	UnknownKind  PaymentKind = "unknown"
)

// PaymentStatus represents the status of a payment.
type PaymentStatus string

const (
	Pending  PaymentStatus = "pending"
	Paid     PaymentStatus = "paid"
	Failed   PaymentStatus = "failed"
	Refunded PaymentStatus = "refunded"
)

// SubscriptionPaymentType defines types of subscription payments.
type SubscriptionPaymentType string

const (
	SubscriptionPaymentInstanceHosting SubscriptionPaymentType = "instance_hosting"
	SubscriptionPaymentUnknown         SubscriptionPaymentType = "unknown"
)
