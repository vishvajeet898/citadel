package utils

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"gorm.io/gorm"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
)

func GetGoLangTimeFromGormDeletedAt(t *gorm.DeletedAt) *time.Time {
	if t == nil || t.Time.IsZero() {
		return nil
	}
	return &t.Time
}

func GetGormDeletedAtFromGoLangTime(t *time.Time) *gorm.DeletedAt {
	if t == nil || t.IsZero() {
		return nil
	}
	return &gorm.DeletedAt{Time: *t}
}

func GetCurrentTime() *time.Time {
	t := time.Now()
	return &t
}

func AddMinutesToTime(t *time.Time, minutes int) *time.Time {
	if t == nil || t.IsZero() {
		return t
	}
	newTime := t.Add(time.Duration(minutes) * time.Minute)
	return &newTime
}

func StringToTime(s string, format string) time.Time {
	t, _ := time.Parse(format, s)
	return t
}

func GetTimeInString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.String()
}

func GetTimeFromString(t string) *time.Time {
	if t == "" {
		return nil
	}
	time, _ := time.Parse(commonConstants.DateTimeUTCLayout, t)
	return &time
}

func GetCurrentTimeInMilliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetDeletedAtString(deletedAt gorm.DeletedAt) string {
	if deletedAt.Valid {
		return deletedAt.Time.Format("2006-01-02 15:04:05")
	}
	return ""
}

func GetExponentialBackoff(attemptCount int, retryDelay time.Duration) time.Duration {
	baseDelay := math.Pow(2, float64(attemptCount)) * float64(retryDelay)
	jitter := rand.Float64() * baseDelay
	return time.Duration(baseDelay + jitter)
}

func CalculateLabEta(lisSyncTime *time.Time, tat float32) *time.Time {
	tatTime, _ := time.ParseDuration(fmt.Sprintf("%vh", tat))
	labEta := lisSyncTime.Add(tatTime)

	return &labEta
}

func UtcToIst(utc time.Time) time.Time {
	loc, _ := time.LoadLocation(commonConstants.LocalTimeZoneLocation)
	return utc.In(loc)
}

func GetDaysBetweenTimes(startTime, endTime time.Time) int {
	if startTime.IsZero() || endTime.IsZero() {
		return 0
	}
	duration := endTime.Sub(startTime)
	return int(math.Ceil(duration.Hours() / 24))
}
