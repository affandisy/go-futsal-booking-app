package repository

import (
	"database/sql"
	"fmt"
	"futsal-booking-app/internal/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id int) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id int) error
	FindByRole(role domain.Role) ([]*domain.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `INSERT INTO users (name, email, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRow(
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.CreatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (r *userRepository) FindByID(id int) (*domain.User, error) {
	query := `SELECT id, name, email, password_hash, role, created_at FROM users WHERE id = $1`

	user := &domain.User{}

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}

		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return user, nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	query := `SELECT id, name, email, password_hash, role, created_at FROM users WHERE email = $1`

	user := &domain.User{}

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}

		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	query := `UPDATE users SET name=$1, email=$2, password_hash=$3, role=$4 WHERE id=$5`

	result, err := r.db.Exec(
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id=$1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) FindByRole(role domain.Role) ([]*domain.User, error) {
	query := `SELECT id, name, email, password_hash, role, created_at FROM users WHERE role=$1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, role)
	if err != nil {
		return nil, fmt.Errorf("error finding users by role: %w", err)
	}

	defer rows.Close()

	users := []*domain.User{}

	for rows.Next() {
		user := &domain.User{}

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.PasswordHash,
			&user.Role,
			&user.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error literating users: %w", err)
	}

	return users, nil
}
