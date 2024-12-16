package helpers

import "time"

// GetNextMonday calculates the date for the next Monday (if today is Monday, the next Monday will be 7 days later).
func GetNextMonday(currentDate time.Time) time.Time {
	// Calculate the number of days to add to get to the next Monday
	daysUntilMonday := (time.Monday - currentDate.Weekday() + 7) % 7
	if daysUntilMonday == 0 {
		daysUntilMonday = 7 // If today is Monday, we want the next Monday, so we add 7 days.
	}
	// Return the next Monday
	return currentDate.AddDate(0, 0, int(daysUntilMonday))
}

// GenerateLastBillDate calculates the last bill date for a loan.
func GenerateLastBillDate(startDate time.Time, numWeeks int) time.Time {
	// Get the next Monday
	nextMonday := GetNextMonday(startDate)
	// The last bill date will be numWeeks-1 weeks after the next Monday
	lastBillDate := nextMonday.AddDate(0, 0, (numWeeks-1)*7) // Subtract 1 to get the last billing week
	return lastBillDate
}
