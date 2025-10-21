package repository

import (
	"database/sql"
	"fmt"
	"futsal-booking-app/internal/domain"
)

type PaymentRepository interface {
	Create(payment *domain.Payment) error
	FindByID(id int) (*domain.Payment, error)
	FindByBookingID(bookingID int) (*domain.Payment, error)
	FindByTransactionID(transactionID string) (*domain.Payment, error)
	Update(payment *domain.Payment) error
	Delete(id int) error
}

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *domain.Payment) error {
	query := `INSERT INTO payments (booking_id, amount, payment_gateway, transaction_id, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := r.db.QueryRow(
		query,
		payment.BookingID,
		payment.Amount,
		payment.PaymentGateway,
		payment.TransactionID,
		payment.Status,
		payment.CreatedAt,
		payment.UpdatedAt,
	).Scan(&payment.ID)

	if err != nil {
		return fmt.Errorf("error creating payment: %w", err)
	}

	return nil
}

func (r *paymentRepository) FindByID(id int) (*domain.Payment, error) {
	query := `SELECT id, booking_id, amount, payment_gateway, transaction_id, status, created_at, updated_at FROM payments WHERE id=$1`

	payment := &domain.Payment{}

	err := r.db.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.BookingID,
		&payment.Amount,
		&payment.PaymentGateway,
		&payment.TransactionID,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("error finding payment: %w", err)
	}

	return payment, nil
}

func (r *paymentRepository) FindByBookingID(bookingID int) (*domain.Payment, error) {
	query := `SELECT id, booking_id, amount, payment_gateway, transaction_Id, status, created_at, updated_at FROM payments WHERE booking_id=$1`

	payment := &domain.Payment{}

	err := r.db.QueryRow(query, bookingID).Scan(
		&payment.ID,
		&payment.BookingID,
		&payment.Amount,
		&payment.PaymentGateway,
		&payment.TransactionID,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("error finding payment: %w", err)
	}

	return payment, nil
}

func (r *paymentRepository) FindByTransactionID(transactionID string) (*domain.Payment, error) {
	query := `SELECT id, booking_id, amount, payment_gateway, transaction_id, status, created_at, updated_at FROM payments WHERE transaction_id=$1`

	payment := &domain.Payment{}

	err := r.db.QueryRow(query, transactionID).Scan(
		&payment.ID,
		&payment.BookingID,
		&payment.Amount,
		&payment.PaymentGateway,
		&payment.TransactionID,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("error finding payment: %w", err)
	}

	return payment, nil
}

func (r *paymentRepository) Update(payment *domain.Payment) error {
	query := `UPDATE payments SET booking_id=$1, amount=$2, payment_gateway=$3, transaction_id=$4, status=$5, updated_at=$6 WHERE id=$7`

	result, err := r.db.Exec(
		query,
		payment.BookingID,
		payment.Amount,
		payment.PaymentGateway,
		payment.TransactionID,
		payment.Status,
		payment.UpdatedAt,
		payment.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating payment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

func (r *paymentRepository) Delete(id int) error {
	query := `DELETE FROM payments WHERE id=$1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting payment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}
