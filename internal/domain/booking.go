package domain

import "time"

type BookingStatus string

const (
	BookingPending   BookingStatus = "PENDING"
	BookingConfirmed BookingStatus = "CONFIRMED"
	BookingCancelled BookingStatus = "CANCELLED"
	BookingCompleted BookingStatus = "COMPLETED"
)

type Booking struct {
	ID         int
	UserID     int
	FieldID    int
	StartTime  time.Time
	EndTime    time.Time
	TotalPrice int
	Status     BookingStatus
	PaymentID  *int
	CreatedAt  time.Time
}

func (b *Booking) GetDuration() float64 {
	duration := b.EndTime.Sub(b.StartTime)

	return duration.Hours()
}

func (b *Booking) GetDurationHours() int {
	duration := b.GetDuration()

	hours := int(duration)
	if duration > float64(hours) {
		hours++
	}

	return hours
}

func (b *Booking) IsPending() bool {
	return b.Status == BookingPending
}

func (b *Booking) IsConfirmed() bool {
	return b.Status == BookingConfirmed
}

func (b *Booking) IsCancelled() bool {
	return b.Status == BookingCancelled
}

func (b *Booking) IsCompleted() bool {
	return b.Status == BookingCompleted
}

func (b *Booking) CanBeCancelled(now time.Time) bool {
	if b.Status != BookingPending && b.Status != BookingConfirmed {
		return false
	}

	timeUntilStart := b.StartTime.Sub(now)
	minCancellationTime := 2 * time.Hour

	return timeUntilStart > minCancellationTime
}

func (b *Booking) IsActive(now time.Time) bool {
	return b.Status == BookingConfirmed && now.After(b.StartTime) && now.Before(b.EndTime)
}
