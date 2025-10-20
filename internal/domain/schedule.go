package domain

import "time"

type DayOfWeek int

const (
	Sunday    DayOfWeek = 0
	Monday    DayOfWeek = 1
	Tuesday   DayOfWeek = 2
	Wednesday DayOfWeek = 3
	Thursday  DayOfWeek = 4
	Friday    DayOfWeek = 5
	Saturday  DayOfWeek = 6
)

type Schedule struct {
	ID        int
	FieldID   int
	DayOfWeek DayOfWeek
	OpenTime  time.Time
	CloseTime time.Time
}

func (s *Schedule) IsOpen(t time.Time) bool {
	day := DayOfWeek(t.Weekday())

	if day != s.DayOfWeek {
		return false
	}

	checkHour, checkMin, checkSec := t.Clock()
	checkTime := checkHour*3600 + checkMin*60 + checkSec

	openHour, openMin, openSec := s.OpenTime.Clock()
	openTimeInSec := openHour*3600 + openMin*60 + openSec

	closeHour, closeMin, closeSec := s.CloseTime.Clock()
	closeTimeInSec := closeHour*3600 + closeMin*60 + closeSec

	return checkTime >= openTimeInSec && checkTime <= closeTimeInSec
}

func (s *Schedule) GetDayName() string {
	dayNames := map[DayOfWeek]string{
		Sunday:    "Minggu",
		Monday:    "Senin",
		Tuesday:   "Selasa",
		Wednesday: "Rabu",
		Thursday:  "Kamis",
		Friday:    "Jumat",
		Saturday:  "Sabtu",
	}
	return dayNames[s.DayOfWeek]
}
