package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// SetTime handles both RFC3339 format and legacy simple time format for backward compatibility
// RFC3339: "2024-01-15T18:30:00+02:00" (preferred) - preserves original timezone
// Legacy: "18:30:00" (assumes UTC)
func SetTime(timeStr string, now time.Time) (time.Time, error) {
	// Check if it's RFC3339 format (contains 'T' and timezone info)
	if strings.Contains(timeStr, "T") && (strings.Contains(timeStr, "+") || strings.Contains(timeStr, "-") || strings.HasSuffix(timeStr, "Z")) {
		// Parse RFC3339 format
		parsedTime, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			return time.Time{}, errors.New("the time string is malformed - expected RFC3339 format")
		}

		return time.Date(
			now.Year(), now.Month(), now.Day(),
			parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(),
			0, parsedTime.Location()), nil
	}

	// Legacy format: simple time "HH:MM:SS"
	parsedTime, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("the time string is malformed - expected either RFC3339 format or HH:MM:SS format, got: %s", timeStr)
	}

	// Use UTC for consistency when no timezone is specified
	return time.Date(
		now.Year(), now.Month(), now.Day(),
		parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(),
		0, time.UTC), nil
}

func InRangeTime(dateStart time.Time, dateEnd time.Time, now time.Time) bool {
	// Convert all times to UTC for consistent comparison
	dateStart = dateStart.UTC()
	dateEnd = dateEnd.UTC()
	now = now.UTC()

	if dateEnd.Before(dateStart) {
		dateEnd = dateEnd.AddDate(0, 0, 1)
	}

	if now.Equal(dateStart) || now.Equal(dateEnd) {
		return true
	}

	return now.After(dateStart) && now.Before(dateEnd)
}

func DaysBetweenDates(todayDate time.Time, eventDate time.Time) int {
	// Convert both dates to UTC for accurate day calculation
	todayDate = todayDate.UTC()
	eventDate = eventDate.UTC()

	if todayDate.After(eventDate) {
		days := todayDate.Sub(eventDate).Hours() / 24
		return int(days)
	}
	return 0
}
