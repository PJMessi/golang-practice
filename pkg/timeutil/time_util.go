package timeutil

import (
	"strconv"
	"time"
)

func ConvertDurationStrToSec(durationStr string) (seconds int, err error) {
	seconds, err = strconv.Atoi(durationStr)
	if err != nil {
		return 0, err
	}

	return seconds, nil
}

func GetTimestampAfterDurationStr(durationStr string) (timestamp int64, err error) {
	seconds, err := ConvertDurationStrToSec(durationStr)
	if err != nil {
		return 0, err
	}

	return time.Now().Add(time.Second * time.Duration(seconds)).Unix(), nil
}

func GetCurrentDateTimeStr() time.Time {
	return time.Now()
}
