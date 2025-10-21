package repository

import (
	"database/sql"
	"fmt"
	"futsal-booking-app/internal/domain"
	"time"
)

type BookingRepository interface {
	Create(booking *domain.Booking) error
	FindByID(id int) (*domain.Booking, error)
	FindByUserID(userID int) ([]*domain.Booking, error)
	FindByFieldID(fieldID int) ([]*domain.Booking, error)
	Update(booking *domain.Booking) error
	Delete(id int) error

	CheckAvailability(fieldID int, startTime, endTime time.Time) (bool, error)
	FindConflictingBookings(fieldID int, startTime, endTime time.Time) ([]*domain.Booking, error)
}

type bookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(booking *domain.Booking) error {
	query := `INSERT INTO bookings (user_id, field_id, start_time, end_time, total_price, status, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := r.db.QueryRow(
		query,
		booking.UserID,
		booking.FieldID,
		booking.StartTime,
		booking.EndTime,
		booking.TotalPrice,
		booking.Status,
		booking.CreatedAt,
	).Scan(&booking.ID)

	if err != nil {
		return fmt.Errorf("error creating booking: %w", err)
	}

	return nil
}

func (r *bookingRepository) FindByID(id int) (*domain.Booking, error) {
	query := `SELECT id, user_id, field_id, start_time, end_time, total_price, status, created_at FROM bookings WHERE id=$1`

	booking := &domain.Booking{}
	err := r.db.QueryRow(query, id).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.FieldID,
		&booking.StartTime,
		&booking.EndTime,
		&booking.TotalPrice,
		&booking.Status,
		&booking.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, fmt.Errorf("error finding booking: %w", err)
	}

	return booking, nil
}

func (r *bookingRepository) FindByUserID(userID int) ([]*domain.Booking, error) {
	query := `SELECT id, user_id, field_id, start_time, end_time, total_price, status, created_at FROM bookings WHERE user_id=$1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error finding bookings by user: %w", err)
	}
	defer rows.Close()

	bookings := []*domain.Booking{}

	for rows.Next() {
		booking := &domain.Booking{}
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.FieldID,
			&booking.StartTime,
			&booking.EndTime,
			&booking.TotalPrice,
			&booking.Status,
			&booking.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning booking: %w", err)
		}
		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating bookings: %w", err)
	}

	return bookings, nil
}

func (r *bookingRepository) FindByFieldID(fieldID int) ([]*domain.Booking, error) {
	query := `SELECT id, user_id, field_id, start_time, end_time, total_price, status, created_at FROM bookings WHERE field_id=$1 ORDER BY start_time DESC`

	rows, err := r.db.Query(query, fieldID)
	if err != nil {
		return nil, fmt.Errorf("error finding bookings by field: %w", err)
	}
	defer rows.Close()

	bookings := []*domain.Booking{}

	for rows.Next() {
		booking := &domain.Booking{}
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.FieldID,
			&booking.StartTime,
			&booking.EndTime,
			&booking.TotalPrice,
			&booking.Status,
			&booking.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning booking: %w", err)
		}
		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating bookings: %w", err)
	}

	return bookings, nil
}

func (r *bookingRepository) Update(booking *domain.Booking) error {
	query := `UPDATE bookings SET user_id=$1, field_id=$2, start_time=$3, end_time=$4, total_price=$5, status=$6 WHERE id=$7`

	result, err := r.db.Exec(
		query,
		booking.UserID,
		booking.FieldID,
		booking.StartTime,
		booking.EndTime,
		booking.TotalPrice,
		booking.Status,
		booking.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating booking: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("bboking not found")
	}

	return nil
}

func (r *bookingRepository) Delete(id int) error {
	query := `DELETE FROM bookings WHERE id=$1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting booking: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("booking not found")
	}

	return nil
}

func (r *bookingRepository) CheckAvailability(fieldID int, startTime, endTime time.Time) (bool, error) {
	query := `SELECT COUNT(*) FROM bookings WHERE field_id=$1 AND status IN ('CONFIRMED', 'PENDING') AND start_time < $3, AND end_time > $2`

	var count int

	err := r.db.QueryRow(query, fieldID, startTime, endTime).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking availability: %w", err)
	}

	return count == 0, nil
}

func (r *bookingRepository) FindConflictingBookings(fieldID int, startTime, endTime time.Time) ([]*domain.Booking, error) {
	query := `SELECT id, user_id, field_id, start_time, end_time, total_price, status, created_at FROM bookings WHERE field_id=$1 AND status IN ('CONFIRMED','PENDING') AND start_time < $3 AND end_time < $2 ORDER BY start_time`

	rows, err := r.db.Query(query, fieldID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("error finding conflicting bookings: %w", err)
	}
	defer rows.Close()

	bookings := []*domain.Booking{}

	for rows.Next() {
		booking := &domain.Booking{}
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.FieldID,
			&booking.StartTime,
			&booking.EndTime,
			&booking.TotalPrice,
			&booking.Status,
			&booking.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning booking: %w", err)
		}
		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating bookings: %w", err)
	}

	return bookings, nil
}
