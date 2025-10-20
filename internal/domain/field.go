package domain

import "time"

type Field struct {
	ID           int
	OwnerID      int
	Name         string
	Address      string
	Description  string
	PricePerHour int
	ImageURL     string
	CreatedAt    time.Time
}

func (f *Field) CalculatePrice(hours int) int {
	return f.PricePerHour * hours
}

func (f *Field) IsOwnedBy(userID int) bool {
	return f.OwnerID == userID
}
