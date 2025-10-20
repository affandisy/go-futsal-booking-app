package domain

import "time"

type Role string

const (
	RoleCustomer Role = "CUSTOMER"
	RoleOwner    Role = "OWNER"
)

type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
}

func (u *User) IsOwner() bool {
	return u.Role == RoleOwner
}

func (u *User) IsCustomer() bool {
	return u.Role == RoleCustomer
}
