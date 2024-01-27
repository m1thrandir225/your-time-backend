package util

import "time"

func IsReminderBeforeDue(reminder, due time.Time) bool {
	return reminder.Before(due)
}
