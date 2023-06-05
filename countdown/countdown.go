package countdown

import (
	"time"
)

func GetRemainingTime(targetDate time.Time) (remainingHours int, remainingMinutes int) {
	remainingTime := time.Until(targetDate)
	remainingHours = int(remainingTime.Hours()) - 1
	remainingMinutes = int(remainingTime.Minutes()) % 60
	return
}
