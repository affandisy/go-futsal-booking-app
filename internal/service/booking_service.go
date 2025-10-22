package service

import (
	"fmt"
	"futsal-booking-app/internal/domain"
	"futsal-booking-app/internal/repository"
	"time"
)

type BookingService interface {
	CreateBooking(userID, fieldID int, startTime time.Time, durationHours int) (*domain.Booking, error)
	GetBookingByID(id int) (*domain.Booking, error)
	GetMyBookings(userID int) ([]*domain.Booking, error)
	GetFileBookings(fieldID int) ([]*domain.Booking, error)
	CancelBooking(userID, bookingID int) error

	ConfirmBooking(bookingID int) error
	CompleteBooking(bookingID int) error
}

type bookingService struct {
	bookingRepo repository.BookingRepository
	fieldRepo   repository.FieldRepository
	paymentRepo repository.PaymentRepository
}

// func NewBookingService(bookingRepo repository.BookingRepository, fieldRepo repository.FieldRepository, paymentRepo repository.PaymentRepository) BookingService {
// 	return &bookingService{
// 		bookingRepo: bookingRepo,
// 		fieldRepo:   fieldRepo,
// 		paymentRepo: paymentRepo,
// 	}
// }

func (u *bookingService) CreateBooking(userID, fieldID int, startTime time.Time, durationHours int) (*domain.Booking, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	if fieldID <= 0 {
		return nil, fmt.Errorf("invalid field ID")
	}

	if durationHours <= 0 {
		return nil, fmt.Errorf("duration must be at least 1 hour")
	}

	if startTime.Before(time.Now()) {
		return nil, fmt.Errorf("cannot book in the past")
	}

	endTime := startTime.Add(time.Duration(durationHours) * time.Hour)

	field, err := u.fieldRepo.FindByID(fieldID)
	if err != nil {
		return nil, fmt.Errorf("field not found")
	}

	available, err := u.bookingRepo.CheckAvailability(fieldID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("error checking availability: %w", err)
	}

	if !available {
		return nil, fmt.Errorf("time slot is not available")
	}

	totalPrice := field.CalculatePrice(durationHours)

	booking := &domain.Booking{
		UserID:     userID,
		FieldID:    fieldID,
		StartTime:  startTime,
		EndTime:    endTime,
		TotalPrice: totalPrice,
		Status:     domain.BookingPending,
		CreatedAt:  time.Now(),
	}

	if err := u.bookingRepo.Create(booking); err != nil {
		return nil, fmt.Errorf("error creating booking: %w", err)
	}

	payment := &domain.Payment{
		BookingID:      booking.ID,
		Amount:         totalPrice,
		PaymentGateway: "Midtrans",
		TransactionID:  fmt.Sprintf("TRX-%d-%d", booking.ID, time.Now().Unix()),
		Status:         domain.PaymentPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := u.paymentRepo.Create(payment); err != nil {
		return nil, fmt.Errorf("error creating payment: %w", err)
	}

	booking.PaymentID = &payment.ID

	if err := u.bookingRepo.Update(booking); err != nil {
		return nil, fmt.Errorf("error updating booking: %w", err)
	}

	return booking, nil

}
