package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var loc = time.Now().Local().Location()

func TestSetTime(t *testing.T) {

	faketime := time.Date(2022, time.March, 2, 21, 0, 0, 0, loc)

	testCases := []struct {
		name         string
		dateStr      string
		dateExpected time.Time
	}{
		{
			name:         "OK - Test 1",
			dateStr:      "10:00:00",
			dateExpected: time.Date(2022, time.March, 2, 10, 0, 0, 0, loc),
		},
		{
			name:         "OK - Test 2",
			dateStr:      "00:00:00",
			dateExpected: time.Date(2022, time.March, 2, 0, 0, 0, 0, loc),
		},
	}

	for _, testCase := range testCases {
		e, _ := SetTime(testCase.dateStr, faketime)
		assert.Equal(t, e, testCase.dateExpected, testCase.name)
	}
}

func TestSetTimeError(t *testing.T) {

	faketime := time.Date(2022, time.March, 2, 21, 0, 0, 0, loc)

	testCases := []struct {
		name    string
		dateStr string
	}{
		{
			name:    "KO - minute and second not set",
			dateStr: "10",
		},
		{
			name:    "KO - second not set",
			dateStr: "10:00",
		},
		{
			name:    "KO - hour not set",
			dateStr: ":10:00",
		},
		{
			name:    "KO - values is nil",
			dateStr: "",
		},
		{
			name:    "KO - min, sec is not integer",
			dateStr: "10:aa:aa",
		},
		{
			name:    "KO - sec is not integer",
			dateStr: "10:10:aa",
		},
	}

	for _, testCase := range testCases {
		_, err := SetTime(testCase.dateStr, faketime)
		assert.Error(t, err, testCase.name)
	}
}

func TestInRangeTime(t *testing.T) {

	faketime := time.Date(2022, time.March, 2, 21, 0, 0, 0, loc)

	testCases := []struct {
		expected  bool
		dateStart time.Time
		dateEnd   time.Time
		name      string
	}{
		{
			name:      "OK - is inside the period ",
			dateStart: time.Date(2022, time.March, 2, 20, 0, 0, 0, loc),
			dateEnd:   time.Date(2022, time.March, 2, 22, 0, 0, 0, loc),
			expected:  true,
		},
		{
			name:      "OK - dateStart and time.Now is Equal",
			dateStart: time.Date(2022, time.March, 2, 21, 0, 0, 0, loc),
			dateEnd:   time.Date(2022, time.March, 2, 22, 0, 0, 0, loc),
			expected:  true,
		},
		{
			name:      "OK - dateEnd and time.Now is Equal",
			dateStart: time.Date(2022, time.March, 2, 20, 0, 0, 0, loc),
			dateEnd:   time.Date(2022, time.March, 2, 21, 0, 0, 0, loc),
			expected:  true,
		},
		{
			name:      "OK - dateEnd is after midnight",
			dateStart: time.Date(2022, time.March, 2, 20, 0, 0, 0, loc),
			dateEnd:   time.Date(2022, time.March, 2, 0, 30, 0, 0, loc),
			expected:  true,
		},
		{
			name:      "KO - dateStart and dateEnd is inverted",
			dateStart: time.Date(2022, time.March, 2, 22, 0, 0, 0, loc),
			dateEnd:   time.Date(2022, time.March, 2, 0, 30, 0, 0, loc),
			expected:  false,
		},
		{
			name:      "KO - is outside the period",
			dateStart: time.Date(2022, time.March, 2, 19, 0, 0, 0, loc),
			dateEnd:   time.Date(2022, time.March, 2, 20, 30, 0, 0, loc),
			expected:  false,
		},
	}

	for _, testCase := range testCases {
		err := InRangeTime(testCase.dateStart, testCase.dateEnd, faketime)
		assert.Equal(t, err, testCase.expected, testCase.name)
	}
}

func TestDaysBetweenDates(t *testing.T) {

	testCases := []struct {
		expected  int
		today     time.Time
		eventDate time.Time
		name      string
	}{
		{
			name:      "OK - 10d8h",
			today:     time.Date(2022, time.March, 13, 20, 0, 0, 0, loc),
			eventDate: time.Date(2022, time.March, 2, 22, 0, 0, 0, loc),
			expected:  10,
		},
		{
			name:      "OK - 9d,8h",
			today:     time.Date(2022, time.March, 12, 20, 0, 0, 0, loc),
			eventDate: time.Date(2022, time.March, 2, 22, 0, 0, 0, loc),
			expected:  9,
		},
		{
			name:      "OK - 10d,1h",
			today:     time.Date(2022, time.March, 12, 20, 0, 0, 0, loc),
			eventDate: time.Date(2022, time.March, 2, 19, 0, 0, 0, loc),
			expected:  10,
		},
		{
			name:      "OK - 10d,1h",
			today:     time.Date(2022, time.March, 2, 20, 0, 0, 0, loc),
			eventDate: time.Date(2022, time.March, 12, 19, 0, 0, 0, loc),
			expected:  0,
		},
		{
			name:      "OK - 10d,1h",
			today:     time.Date(2022, time.March, 2, 20, 0, 0, 0, loc),
			eventDate: time.Date(2022, time.March, 2, 19, 0, 0, 0, loc),
			expected:  0,
		},
	}

	for _, testCase := range testCases {
		days := DaysBetweenDates(testCase.today, testCase.eventDate)
		assert.Equal(t, testCase.expected, days, testCase.name)
	}
}
