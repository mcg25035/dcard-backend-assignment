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