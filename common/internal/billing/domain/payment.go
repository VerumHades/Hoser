package domain

import (
	"errors"
	"time"

	"common/internal/shared"
)

// Payment represents a financial transaction associated with a billing account.
type Payment struct {
	id               shared.PaymentID
	billingAccountID shared.BillingAccountID
	amount           shared.Money
	status           PaymentStatus
	kind             PaymentKind
	createdAt        time.Time
	paidAt           *time.Time

	oneTimeMetadata      *OneTimePaymentMetadata
	subscriptionMetadata *SubscriptionPaymentMetadata
}

// OneTimePaymentMetadata represents additional info for one-time payments.
type OneTimePaymentMetadata struct {
	listingID shared.ListingID
	userID    shared.UserID
}

// SubscriptionPaymentMetadata represents additional info for subscription payments.
type SubscriptionPaymentMetadata struct {
	instanceID       shared.InstanceID
	currentPeriodEnd time.Time
	paymentType      SubscriptionPaymentType
}

// NewOneTimePayment creates a new one-time payment with generated ID.
func NewOneTimePayment(billingAccountID shared.BillingAccountID, amount shared.Money, listingID shared.ListingID, userID shared.UserID) (*Payment, error) {
	id := shared.PaymentID(shared.GenerateUUID())
	metadata := &OneTimePaymentMetadata{listingID: listingID, userID: userID}
	payment := &Payment{
		id:               id,
		billingAccountID: billingAccountID,
		amount:           amount,
		status:           Pending,
		kind:             OneTime,
		createdAt:        time.Now().UTC(),
		oneTimeMetadata:  metadata,
	}
	if err := payment.Validate(); err != nil {
		return nil, err
	}
	return payment, nil
}

// NewSubscriptionPayment creates a new subscription payment with generated ID.
func NewSubscriptionPayment(billingAccountID shared.BillingAccountID, amount shared.Money, instanceID shared.InstanceID, periodEnd time.Time, paymentType SubscriptionPaymentType) (*Payment, error) {
	id := shared.PaymentID(shared.GenerateUUID())
	metadata := &SubscriptionPaymentMetadata{
		instanceID:       instanceID,
		currentPeriodEnd: periodEnd,
		paymentType:      paymentType,
	}
	payment := &Payment{
		id:                   id,
		billingAccountID:     billingAccountID,
		amount:               amount,
		status:               Pending,
		kind:                 Subscription,
		createdAt:            time.Now().UTC(),
		subscriptionMetadata: metadata,
	}
	if err := payment.Validate(); err != nil {
		return nil, err
	}
	return payment, nil
}

// NewPaymentWithID creates a payment with an existing ID (for repository hydration).
func NewPaymentWithID(id shared.PaymentID, billingAccountID shared.BillingAccountID, amount shared.Money, status PaymentStatus, kind PaymentKind, createdAt time.Time, paidAt *time.Time, oneTime *OneTimePaymentMetadata, subscription *SubscriptionPaymentMetadata) (*Payment, error) {
	payment := &Payment{
		id:                   id,
		billingAccountID:     billingAccountID,
		amount:               amount,
		status:               status,
		kind:                 kind,
		createdAt:            createdAt,
		paidAt:               paidAt,
		oneTimeMetadata:      oneTime,
		subscriptionMetadata: subscription,
	}
	if err := payment.Validate(); err != nil {
		return nil, err
	}
	return payment, nil
}

// ID returns the payment ID.
func (p *Payment) ID() shared.PaymentID {
	return p.id
}

// shared.BillingAccountID returns the associated billing account ID.
func (p *Payment) BillingAccountID() shared.BillingAccountID {
	return p.billingAccountID
}

// Amount returns the payment amount.
func (p *Payment) Amount() shared.Money {
	return p.amount
}

// Status returns the current status.
func (p *Payment) Status() PaymentStatus {
	return p.status
}

// Kind returns the payment kind.
func (p *Payment) Kind() PaymentKind {
	return p.kind
}

// CreatedAt returns creation timestamp.
func (p *Payment) CreatedAt() time.Time {
	return p.createdAt
}

// PaidAt returns the timestamp when payment was completed (if any).
func (p *Payment) PaidAt() *time.Time {
	return p.paidAt
}

// MarkPaid sets the payment status to Paid and updates PaidAt timestamp.
func (p *Payment) MarkPaid() {
	now := time.Now().UTC()
	p.status = Paid
	p.paidAt = &now
}

// MarkFailed sets the payment status to Failed.
func (p *Payment) MarkFailed() {
	p.status = Failed
}

// Refund sets the payment status to Refunded.
func (p *Payment) Refund() {
	p.status = Refunded
}

// OneTimeMetadata returns the metadata for one-time payments, or nil if not a one-time payment.
func (p *Payment) OneTimeMetadata() *OneTimePaymentMetadata {
	return p.oneTimeMetadata
}

// SubscriptionMetadata returns the metadata for subscription payments, or nil if not a subscription payment.
func (p *Payment) SubscriptionMetadata() *SubscriptionPaymentMetadata {
	return p.subscriptionMetadata
}

// Validate ensures the payment entity is in a valid state.
func (p *Payment) Validate() error {
	if p.id == "" {
		return errors.New("payment ID cannot be empty")
	}
	if p.billingAccountID == "" {
		return errors.New("billing account ID cannot be empty")
	}
	if p.amount.IsZero() {
		return errors.New("payment amount cannot be zero")
	}
	switch p.status {
	case Pending, Paid, Failed, Refunded:
		// valid
	default:
		return errors.New("invalid payment status")
	}
	switch p.kind {
	case OneTime:
		if p.oneTimeMetadata == nil {
			return errors.New("one-time payment metadata is required")
		}
	case Subscription:
		if p.subscriptionMetadata == nil {
			return errors.New("subscription payment metadata is required")
		}
	case UnknownKind:
		// allow unknown
	default:
		return errors.New("invalid payment kind")
	}
	return nil
}

// ListingID returns the listing ID associated with the one-time payment.
func (m *OneTimePaymentMetadata) ListingID() shared.ListingID {
	return m.listingID
}

// UserID returns the user ID associated with the one-time payment.
func (m *OneTimePaymentMetadata) UserID() shared.UserID {
	return m.userID
}

// SubscriptionPaymentMetadata getters

// InstanceID returns the instance ID associated with the subscription payment.
func (m *SubscriptionPaymentMetadata) InstanceID() shared.InstanceID {
	return m.instanceID
}

// CurrentPeriodEnd returns the end of the current subscription period.
func (m *SubscriptionPaymentMetadata) CurrentPeriodEnd() time.Time {
	return m.currentPeriodEnd
}

// PaymentType returns the type of the subscription payment.
func (m *SubscriptionPaymentMetadata) PaymentType() SubscriptionPaymentType {
	return m.paymentType
}
