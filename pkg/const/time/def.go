package time

import (
	"time"
)

const (
	// DaysPerYear ..
	DaysPerYear = 365
	// HoursPerDay ..
	HoursPerDay = 24
	// MinutesPerHour ..
	MinutesPerHour = 60
	// SecondsPerMinute ..
	SecondsPerMinute = 60
	// SecondsPerHour ..
	SecondsPerHour = 60 * SecondsPerMinute
	// SecondsPerDay ..
	SecondsPerDay = 24 * SecondsPerHour
	// SecondsPerMonth ..
	SecondsPerMonth = 30 * SecondsPerDay
	// SecondsPerYear ..
	SecondsPerYear = DaysPerYear * SecondsPerDay
)

func TomorrowStart() time.Time {
	now := time.Now()
	y, m, d := now.Date()
	return time.Date(y, m, d+1, 0, 0, 0, 0, now.Location())
}
