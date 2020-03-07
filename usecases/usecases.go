package usecases

import "time"

var (
	durationNone = time.Duration(0)
	durationWeek = time.Minute * 60 * 24 * 7
)

func GetCurrentOffset() time.Duration {
	now := time.Now()

	b := time.Date(now.Year(), now.Month(), now.Day() - int(now.Weekday()), 0, 0, 0, 0, now.Location())
	return time.Since(b)
}

func SleepUntilOffset(offset time.Duration) {
	currentOffset := GetCurrentOffset()

	if offset < currentOffset {
		offset = offset + durationWeek
	}

	timeToSleep := offset - currentOffset
	time.Sleep(timeToSleep)
}
