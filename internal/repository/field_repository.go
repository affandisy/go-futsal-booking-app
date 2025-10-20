package repository

import (
	"database/sql"
	"fmt"
	"futsal-booking-app/internal/domain"
)

type FieldRepository interface {
	Create(field *domain.Field) error
	FindByID(id int) (*domain.Field, error)
	FindByOwnerID(ownerID int) ([]*domain.Field, error)
	FindAll() ([]*domain.Field, error)
	Update(field *domain.Field) error
	Delete(id int) error

	CreateSchedule(schedule *domain.Schedule) error
	FindScheduleByFieldID(fieldID int) ([]*domain.Schedule, error)
	UpdateSchedule(schedule *domain.Schedule) error
	DeleteScheduleByFieldID(fieldID int) error
}

type fieldRepository struct {
	db *sql.DB
}

func NewFieldRepository(db *sql.DB) FieldRepository {
	return &fieldRepository{db: db}
}

func (r *fieldRepository) Create(field *domain.Field) error {
	query := `INSERT INTO fields (owner_id, name, address, description, price_per_hour, image_url, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := r.db.QueryRow(
		query,
		field.OwnerID,
		field.Name,
		field.Address,
		field.Description,
		field.PricePerHour,
		field.ImageURL,
		field.CreatedAt,
	).Scan(&field.ID)

	if err != nil {
		return fmt.Errorf("error creating field: %w", err)
	}

	return nil
}

func (r *fieldRepository) FindByID(id int) (*domain.Field, error) {
	query := `SELECT id, owner_id, address, description, price_per_hour, image_url, created_at FROM fields WHERE id=$1`

	field := &domain.Field{}

	err := r.db.QueryRow(query, id).Scan(
		&field.ID,
		&field.OwnerID,
		&field.Name,
		&field.Address,
		&field.Description,
		&field.PricePerHour,
		&field.ImageURL,
		&field.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("field not found")
		}
		return nil, fmt.Errorf("error finding field: %w", err)
	}

	return field, nil
}

func (r *fieldRepository) FindByOwnerID(ownerID int) ([]*domain.Field, error) {
	query := `SELECT id, owner_id, name, address, description, price_per_hour, image_url, created_at FROM fields WHERE owner_id=$1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("error finding fields by owner: %w", err)
	}

	defer rows.Close()

	fields := []*domain.Field{}

	for rows.Next() {
		field := &domain.Field{}
		err := rows.Scan(
			&field.ID,
			&field.OwnerID,
			&field.Name,
			&field.Address,
			&field.Description,
			&field.PricePerHour,
			&field.ImageURL,
			&field.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning field: %w", err)
		}

		fields = append(fields, field)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating fields: %w", err)
	}

	return fields, nil
}

func (r *fieldRepository) FindAll() ([]*domain.Field, error) {
	query := `SELECT id, owner_id, name, address, description, price_per_hour, image_url, created_at FROM fields ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error finding all fields: %w", err)
	}

	defer rows.Close()

	fields := []*domain.Field{}

	for rows.Next() {
		field := &domain.Field{}
		err := rows.Scan(
			&field.ID,
			&field.OwnerID,
			&field.Name,
			&field.Address,
			&field.Description,
			&field.PricePerHour,
			&field.ImageURL,
			&field.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning field: %w", err)
		}

		fields = append(fields, field)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating fields: %w", err)
	}

	return fields, nil
}

func (r *fieldRepository) Update(field *domain.Field) error {
	query := `UPDATE fields SET name=$1, address=$2, description=$3, price_per_hour=$4, image_url=$5 WHERE id=$6`

	result, err := r.db.Exec(
		query,
		field.Name,
		field.Address,
		field.Description,
		field.PricePerHour,
		field.ImageURL,
		field.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating field: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("field not found")
	}

	return nil
}

func (r *fieldRepository) Delete(id int) error {
	query := `DELETE FROM fields WHERE id=$1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting field: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("field not found")
	}

	return nil
}

func (r *fieldRepository) CreateSchedule(schedule *domain.Schedule) error {
	query := `INSERT INTO schedules (field_id, day_of_week, open_time, close_time) VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.db.QueryRow(
		query,
		schedule.FieldID,
		schedule.DayOfWeek,
		schedule.OpenTime,
		schedule.CloseTime,
	).Scan(&schedule.ID)

	if err != nil {
		return fmt.Errorf("error creating schedule: %w", err)
	}

	return nil
}

func (r *fieldRepository) FindScheduleByFieldID(fieldID int) ([]*domain.Schedule, error) {
	query := `SELECT id, field_id, day_of_week, open_time, close_time FROM schedules WHERE field_id=$1 ORDER BY day_of_week`

	rows, err := r.db.Query(query, fieldID)
	if err != nil {
		return nil, fmt.Errorf("error finding schedules: %w", err)
	}
	defer rows.Close()

	schedules := []*domain.Schedule{}

	for rows.Next() {
		schedule := &domain.Schedule{}
		err := rows.Scan(
			&schedule.ID,
			&schedule.FieldID,
			&schedule.DayOfWeek,
			&schedule.OpenTime,
			&schedule.CloseTime,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning schedule: %w", err)
		}
		schedules = append(schedules, schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating schedules: %w", err)
	}

	return schedules, nil
}

func (r *fieldRepository) UpdateSchedule(schedule *domain.Schedule) error {
	query := `UPDATE schedules SET day_of_week=$1, open_time=$2, close_time=$3 WHERE id=$4`

	result, err := r.db.Exec(
		query,
		schedule.DayOfWeek,
		schedule.OpenTime,
		schedule.CloseTime,
		schedule.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating schedule: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("schedule not found")
	}

	return nil
}

func (r *fieldRepository) DeleteScheduleByFieldID(fieldID int) error {
	query := `DELETE FROM schedules WHERE field_id=$1`

	_, err := r.db.Query(query, fieldID)
	if err != nil {
		return fmt.Errorf("error deleting schedules: %w", err)
	}

	return nil
}
