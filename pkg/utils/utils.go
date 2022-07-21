package utils

import (
	"errors"
	"time"
)

func SetTime(str string, now time.Time) (time.Time, error) {
	layOut := "15:04:05"
	timeStamp, err := time.Parse(layOut, str)
	if err != nil {
		return time.Time{}, errors.New("the annotation time.start or time.end is malformed")
	}

	hour, min, sec := timeStamp.Clock()
	dateTime := time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, time.UTC)

	return dateTime, nil
}

func InRangeTime(dateStart time.Time, dateEnd time.Time, now time.Time) bool {
	if dateEnd.Before(dateStart) {
		dateEnd = dateEnd.AddDate(0, 0, 1)
	}

	if now.Equal(dateStart) || now.Equal(dateEnd) {
		return true
	}

	return now.After(dateStart) && now.Before(dateEnd)
}

func DaysBetweenDates(todayDate time.Time, eventDate time.Time) int {
	if todayDate.After(eventDate) {
		days := todayDate.Sub(eventDate).Hours() / 24
		return int(days)
	}
	return 0
}
