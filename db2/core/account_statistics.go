package core

import (
	"time"
)

// helper struct selected only with account, field names reflect that
type Statistics struct {
	DailyOutcome   int64 `db:"st_daily_out"`
	WeeklyOutcome  int64 `db:"st_weekly_out"`
	MonthlyOutcome int64 `db:"st_monthly_out"`
	AnnualOutcome  int64 `db:"st_annual_out"`
	UpdatedAt      int64 `db:"st_updated_at"`
}

func getDaysPassed(updatedAt, currentTime time.Time) int {
	if updatedAt.Year() == currentTime.Year() {
		return currentTime.YearDay() - updatedAt.YearDay()
	}

	lastDayOfUpdateAtYear := time.Date(updatedAt.Year(), time.December, 31, 0, 0, 0, 0, time.UTC).YearDay()
	return lastDayOfUpdateAtYear - updatedAt.YearDay() + currentTime.YearDay()
}

func isWeekPassed(updateAt, currentTime time.Time) bool {
	if getDaysPassed(updateAt, currentTime) >= 7 {
		return true
	}

	return currentTime.Weekday() < updateAt.Weekday()
}

func (s *Statistics) ClearObsolete(currentTime time.Time) {
	if s == nil {
		return
	}
	updatedAt := time.Unix(s.UpdatedAt, 0).UTC()

	isYear := updatedAt.Year() < currentTime.Year()
	if isYear {
		s.AnnualOutcome = 0
	}

	isMonth := isYear || updatedAt.Month() < currentTime.Month()
	if isMonth {
		s.MonthlyOutcome = 0
	}

	isWeek := isWeekPassed(updatedAt, currentTime)
	if isWeek {
		s.WeeklyOutcome = 0
	}

	isDay := isYear || updatedAt.YearDay() < currentTime.YearDay()
	if isDay {
		s.DailyOutcome = 0
	}
}
