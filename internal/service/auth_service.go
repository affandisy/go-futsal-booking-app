package service

import (
	"fmt"
	"futsal-booking-app/internal/domain"
	"futsal-booking-app/internal/repository"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(name, email, password string, role domain.Role) (*domain.User, error)
	LoginUser(email, password string) (*domain.User, error)
	GetUserByID(id int) (*domain.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// RegisterUser mendaftarkan user baru ke sistem
// Business logic yang diterapkan:
// 1. Validasi input (name, email, password tidak boleh kosong)
// 2. Validasi format email
// 3. Validasi role (harus CUSTOMER atau OWNER)
// 4. Cek duplikasi email
// 5. Hash password menggunakan bcrypt
// 6. Simpan user ke database
// Parameter:
//   - name: nama lengkap user
//   - email: email user (akan dijadikan username)
//   - password: plain password dari user
//   - role: role user (CUSTOMER atau OWNER)
//
// Return:
//   - *entity.User: user yang berhasil dibuat
//   - error: error jika ada validasi yang gagal
func (u *authService) RegisterUser(name, email, password string, role domain.Role) (*domain.User, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	if strings.TrimSpace(email) == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return nil, fmt.Errorf("invalid email format")
	}

	if len(password) < 6 {
		return nil, fmt.Errorf("password must be at least 6 characters")
	}

	if role != domain.RoleCustomer && role != domain.RoleOwner {
		return nil, fmt.Errorf("invalid role, must be CUSTOMER or OWNER")
	}

	existingUser, err := u.userRepo.FindByEmail(email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	user := &domain.User{
		Name:         name,
		Email:        strings.ToLower(strings.TrimSpace(email)),
		PasswordHash: string(hashedPassword),
		Role:         role,
		CreatedAt:    time.Now(),
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	user.PasswordHash = ""
	return user, nil
}

// LoginUser mengautentikasi user dengan email dan password
// Business logic yang diterapkan:
// 1. Validasi input tidak kosong
// 2. Cari user berdasarkan email
// 3. Verifikasi password dengan bcrypt
// 4. Return user jika berhasil
// Parameter:
//   - email: email user
//   - password: plain password dari user
//
// Return:
//   - *entity.User: user jika autentikasi berhasil
//   - error: error jika autentikasi gagal
func (u *authService) LoginUser(email, password string) (*domain.User, error) {
	if strings.TrimSpace(email) == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	if strings.TrimSpace(password) == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	user, err := u.userRepo.FindByEmail(strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	user.PasswordHash = ""
	return user, nil
}

// GetUserByID mengambil data user berdasarkan ID
// Digunakan untuk mendapatkan info user yang sedang login
// Parameter:
//   - id: ID user
//
// Return:
//   - *entity.User: user jika ditemukan
//   - error: error jika user tidak ditemukan
func (u *authService) GetUserByID(id int) (*domain.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	user.PasswordHash = ""
	return user, nil
}
