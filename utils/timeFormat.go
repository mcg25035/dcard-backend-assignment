package utils

import (
	"time"
)

func DateStringToTimestamp(dateStr string) (int64, error) {
	// "xxxZ" -> "000Z"
	layout := time.RFC3339 
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return 0, err 
	}
	return t.UnixMilli(), nil
}

func TimestampToDateString(timestamp int64) string {
	t := time.UnixMilli(timestamp).UTC()
	return t.Format("2006-01-02T15:04:05.000Z")
}

func TimestampToTime(timestamp int64) time.Time {
	return time.UnixMilli(timestamp).UTC()
}

func DateStringToTime(dateStr string) (time.Time, error) {
	var timestamp, err = DateStringToTimestamp(dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return TimestampToTime(timestamp), nil
}