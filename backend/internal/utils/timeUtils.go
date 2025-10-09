package utils

import (
	"time"
)

const DATETIME_FORMAT = "2006-01-02 15:04:05"

// CleanTime removes the monotonic clock information from a time.Time
// and returns it formatted as a human-readable string.
// Example output: "2025-10-08 09:03:21"
func CleanTime(t time.Time) string {
	// Round(0) strips the monotonic part
	clean := t.Round(0)

	// Return formatted string (consistent across app)
	return clean.Format(DATETIME_FORMAT)
}

// NowClean returns the current time formatted cleanly (shortcut)
func NowClean() string {
	return CleanTime(time.Now())
}
