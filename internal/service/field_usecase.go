package service

import (
	"fmt"
	"futsal-booking-app/internal/domain"
	"futsal-booking-app/internal/repository"
	"strings"
	"time"
)

type FieldService interface {
	CreateField(ownerID int, name, address, description, imageURL string, pricePerHour int) (*domain.Field, error)
	GetFieldByID(id int) (*domain.Field, error)
	GetAllFields() ([]*domain.Field, error)
	GetFieldsByOwnerID(ownerID int) ([]*domain.Field, error)
	UpdateField(fieldID, ownerID int, name, address, description, imageURL string, pricePerHour int) (*domain.Field, error)
	DeleteField(fieldID, ownerID int) error

	SetupSchedules(fieldID, ownerID int, schedules []ScheduleInput) error
	GetScheduleByFieldID(fieldID int) ([]*domain.Schedule, error)

	FindAvailableSlots(fieldID int, date time.Time) ([]TimeSlot, error)
}

type ScheduleInput struct {
	DayOfWeek int
	OpenTime  string
	CloseTime string
}

type TimeSlot struct {
	StartTime time.Time
	EndTime   time.Time
	Available bool
}

type fieldService struct {
	fieldRepo   repository.FieldRepository
	bookingRepo repository.BookingRepository
}

func NewFieldService(fieldRepo repository.FieldRepository, bookingRepo repository.BookingRepository) FieldService {
	return &fieldService{fieldRepo: fieldRepo, bookingRepo: bookingRepo}
}

// CreateField membuat lapangan baru
// Business logic:
// 1. Validasi input (name, address tidak boleh kosong, price harus positif)
// 2. Simpan field ke database
func (u *fieldService) CreateField(ownerID int, name, address, description, imageURL string, pricePerHour int) (*domain.Field, error) {
	if ownerID <= 0 {
		return nil, fmt.Errorf("invalid owner id")
	}

	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("field name cannot be empty")
	}

	if strings.TrimSpace(address) == "" {
		return nil, fmt.Errorf("field address cannot be empty")
	}

	if pricePerHour <= 0 {
		return nil, fmt.Errorf("price per hour must be positive")
	}

	field := &domain.Field{
		OwnerID:      ownerID,
		Name:         name,
		Address:      address,
		Description:  description,
		PricePerHour: pricePerHour,
		ImageURL:     imageURL,
		CreatedAt:    time.Now(),
	}

	if err := u.fieldRepo.Create(field); err != nil {
		return nil, fmt.Errorf("error creating field: %w", err)
	}

	return field, nil
}

func (u *fieldService) GetFieldByID(id int) (*domain.Field, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid field ID")
	}

	field, err := u.fieldRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching fields: %w", err)
	}

	return field, nil
}

func (u *fieldService) GetAllFields() ([]*domain.Field, error) {
	fields, err := u.fieldRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("error fetching fields: %w", err)
	}

	return fields, nil
}

func (u *fieldService) GetFieldsByOwnerID(ownerID int) ([]*domain.Field, error) {
	if ownerID == 0 {
		return nil, fmt.Errorf("invalid owner ID")
	}

	fields, err := u.fieldRepo.FindByOwnerID(ownerID)
	if err != nil {
		return nil, fmt.Errorf("error fetching fields: %w", err)
	}

	return fields, nil
}

func (u *fieldService) UpdateField(fieldID, ownerID int, name, address, description, imagerURL string, pricePerHour int) (*domain.Field, error) {
	if fieldID <= 0 {
		return nil, fmt.Errorf("invalid field ID")
	}

	field, err := u.fieldRepo.FindByID(fieldID)
	if err != nil {
		return nil, fmt.Errorf("field not found")
	}

	if !field.IsOwnedBy(ownerID) {
		return nil, fmt.Errorf("unauthorized: you are not the owner of this field")
	}

	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("field name cannot be empty")
	}

	if strings.TrimSpace(address) == "" {
		return nil, fmt.Errorf("field address cannot be empty")
	}

	if pricePerHour <= 0 {
		return nil, fmt.Errorf("price per hour must be positive")
	}

	field.Name = name
	field.Address = address
	field.Description = description
	field.PricePerHour = pricePerHour
	field.ImageURL = imagerURL

	if err := u.fieldRepo.Update(field); err != nil {
		return nil, fmt.Errorf("error updating field: %w", err)
	}

	return field, nil
}

func (u *fieldService) DeleteField(fieldID, ownerID int) error {
	if fieldID <= 0 {
		return fmt.Errorf("invalid field ID")
	}

	field, err := u.fieldRepo.FindByID(fieldID)
	if err != nil {
		return fmt.Errorf("field not found")
	}

	if !field.IsOwnedBy(ownerID) {
		return fmt.Errorf("unauthorized: you are not the owner of this field")
	}

	if err := u.fieldRepo.Delete(fieldID); err != nil {
		return fmt.Errorf("error deleting field: %w", err)
	}

	return nil
}

func (u *fieldService) SetupSchedules(fieldID, ownerID int, schedules []ScheduleInput) error {
	if fieldID <= 0 {
		return fmt.Errorf("invalid field ID")
	}

	field, err := u.fieldRepo.FindByID(fieldID)
	if err != nil {
		return fmt.Errorf("field not found")
	}

	if !field.IsOwnedBy(ownerID) {
		return fmt.Errorf("unauthorized: you are not the owner of this field")
	}

	if len(schedules) == 0 {
		return fmt.Errorf("at least one schedule is required")
	}

	if err := u.fieldRepo.DeleteScheduleByFieldID(fieldID); err != nil {
		return fmt.Errorf("error deleting old schedules: %w", err)
	}

	for _, input := range schedules {
		if input.DayOfWeek < 0 || input.DayOfWeek > 6 {
			return fmt.Errorf("invalid day of week: %d", input.DayOfWeek)
		}

		openTime, err := time.Parse("15:04", input.OpenTime)
		if err != nil {
			return fmt.Errorf("invalid open time format: %s", input.OpenTime)
		}

		closeTime, err := time.Parse("15:04", input.CloseTime)
		if err != nil {
			return fmt.Errorf("invalid close time format: %s", input.CloseTime)
		}

		if closeTime.Before(openTime) || closeTime.Equal(openTime) {
			return fmt.Errorf("close time must be after open time")
		}

		schedule := &domain.Schedule{
			FieldID:   fieldID,
			DayOfWeek: domain.DayOfWeek(input.DayOfWeek),
			OpenTime:  openTime,
			CloseTime: closeTime,
		}

		if err := u.fieldRepo.CreateSchedule(schedule); err != nil {
			return fmt.Errorf("error creating schedule: %w", err)
		}
	}

	return nil
}

func (u *fieldService) GetScheduleByFieldID(fieldID int) ([]*domain.Schedule, error) {
	if fieldID <= 0 {
		return nil, fmt.Errorf("invalid field ID")
	}

	schedules, err := u.fieldRepo.FindScheduleByFieldID(fieldID)
	if err != nil {
		return nil, fmt.Errorf("error fetching schedules: %w", err)
	}

	return schedules, nil
}

func (u *fieldService) FindAvailableSlots(fieldID int, date time.Time) ([]TimeSlot, error) {
	if fieldID <= 0 {
		return nil, fmt.Errorf("invalid field ID")
	}

	schedules, err := u.fieldRepo.FindScheduleByFieldID(fieldID)
	if err != nil {
		return nil, fmt.Errorf("error fecthing schedules: %w", err)
	}

	dayOfWeek := domain.DayOfWeek(date.Weekday())
	var relevantSchedule *domain.Schedule
	for _, schedule := range schedules {
		if schedule.DayOfWeek == dayOfWeek {
			relevantSchedule = schedule
			break
		}
	}

	if relevantSchedule == nil {
		return []TimeSlot{}, nil
	}

	slots := []TimeSlot{}

	openHour, openMin, _ := relevantSchedule.OpenTime.Clock()
	closeHour, closeMin, _ := relevantSchedule.CloseTime.Clock()

	currentSlot := time.Date(date.Year(), date.Month(), date.Day(), openHour, openMin, 0, 0, date.Location())
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), closeHour, closeMin, 0, 0, date.Location())

	for currentSlot.Before(endOfDay) {
		slotEnd := currentSlot.Add(1 * time.Hour)

		available, err := u.bookingRepo.CheckAvailability(fieldID, currentSlot, slotEnd)
		if err != nil {
			return nil, fmt.Errorf("error checking availability: %w", err)
		}

		slots = append(slots, TimeSlot{
			StartTime: currentSlot,
			EndTime:   slotEnd,
			Available: available,
		})

		currentSlot = slotEnd
	}

	return slots, nil
}
