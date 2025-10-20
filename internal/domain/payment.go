package domain

import "time"

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "PENDING"
	PaymentSuccess PaymentStatus = "SUCCESS"
	PaymentFailed  PaymentStatus = "FAILED"
)

type Payment struct {
	ID             int
	BookingID      int
	Amount         int
	PaymentGateway string
	TransactionID  string
	Status         PaymentStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (p *Payment) IsPending() bool {
	return p.Status == PaymentPending
}

func (p *Payment) IsSuccess() bool {
	return p.Status == PaymentSuccess
}

func (p *Payment) IsFailed() bool {
	return p.Status == PaymentFailed
}

func (p *Payment) MarkAsSuccess() {
	p.Status = PaymentSuccess
	p.UpdatedAt = time.Now()
}

func (p *Payment) MarkAsFailed() {
	p.Status = PaymentFailed
	p.UpdatedAt = time.Now()
}
