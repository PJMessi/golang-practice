package timeutil

import (
	"fmt"
	"strconv"
	"time"
)

func ConvertDurationStrToSec(durationStr string) (seconds int64, err error) {
	numStr := durationStr[:len(durationStr)-1]
	unitStr := durationStr[len(durationStr)-1:]

	// Parse the numeric part as an integer
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, err
	}

	// Calculate the duration in seconds based on the unit string
	switch unitStr {
	case "s":
		seconds = int64(num)
	case "m":
		seconds = int64(num) * 60 // 1 min = 60 seconds
	case "h":
		seconds = int64(num) * 60 * 60 // 1 hour = 60 minutes = 60 * 60 seconds
	case "d":
		seconds = int64(num) * 24 * 60 * 60 // 1 day = 24 hours = 24 * 60 minutes = 24 * 60 * 60 seconds
	case "M":
		seconds = int64(num) * 30 * 24 * 60 * 60 // 1 month = 30 days = 30 * 24 hours = 30 * 24 * 60 minutes = 30 * 24 * 60 * 60 seconds
	case "y":
		seconds = int64(num) * 365 * 24 * 60 * 60 // 1 year = 365 days = 365 * 24 hours = 365 * 24 * 60 minutes = 365 * 24 * 60 * 60 seconds
	default:
		return 0, fmt.Errorf("timeutil.ConvertDurationStrToSec: unsupported unit '%s'", unitStr)
	}

	return seconds, nil
}

func GetTimestampAfterNSec(seconds int64) int64 {
	return time.Now().Add(time.Second * time.Duration(seconds)).Unix()
}

func GetCurrentTime() time.Time {
	return time.Now()
}
