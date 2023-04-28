package time

import "time"

type TimeUnit int

const (
	Second TimeUnit = iota
	Minute
	Hour
	Day
	Month
	Year
)

func Now() time.Time {
	return time.Now()
}

// MsToTime 毫秒转Time
func MsToTime(epochMilli int64) time.Time {
	return time.UnixMicro(epochMilli)
}

// ParseDateTime 解析日期时间字符串为时间,支持时间格式:yyyy-MM-dd HH:mm:ss
func ParseDateTime(timeStr string) (time.Time, error) {
	return time.Parse(patternToFormat(NormalDatetime), timeStr)
}

// ParseDate 解析日期时间字符串为时间,支持时间格式:yyyy-MM-dd
func ParseDate(timeStr string) (time.Time, error) {
	return time.Parse(patternToFormat(NormalDate), timeStr)
}

// ParseTime 解析日期时间字符串为时间,自定义时间格式
func ParseTime(timeStr string, pattern string) (time.Time, error) {
	return time.Parse(patternToFormat(pattern), timeStr)
}

// FormatNormal 格式化日期时间为yyyy-MM-dd HH:mm:ss格式
func FormatNormal(t0 time.Time) string {
	return t0.Format(patternToFormat(NormalDatetime))
}

// FormatNormalDate 格式化日期时间为yyyy-MM-dd格式
func FormatNormalDate(t0 time.Time) string {
	return t0.Format(patternToFormat(NormalDate))
}

// FormatNormalTime 格式化日期时间为HH:mm:ss格式
func FormatNormalTime(t0 time.Time) string {
	return t0.Format(time.TimeOnly)
}

// FormatWithPattern 格式化日期时间为自定义格式
func FormatWithPattern(t0 time.Time, pattern string) string {
	return t0.Format(patternToFormat(pattern))
}

// Between 获取两个时间差
func Between(start, end time.Time) time.Duration {
	return end.Sub(start)
}

// PlusSeconds 增加时间:单位:秒,seconds < 0,表示时间相减
func PlusSeconds(t0 time.Time, seconds int) time.Time {
	return t0.Add(time.Second * time.Duration(seconds))
}

// PlusMinutes 增加时间:单位:分,minutes < 0,表示时间相减
func PlusMinutes(t0 time.Time, minutes int) time.Time {
	return t0.Add(time.Minute * time.Duration(minutes))
}

// PlusHours 增加时间:单位:小时,hours < 0,表示时间相减
func PlusHours(t0 time.Time, hours int) time.Time {
	return t0.Add(time.Hour * time.Duration(hours))
}

// PlusDays 增加时间:单位:天,days < 0,表示时间相减
func PlusDays(t0 time.Time, days int) time.Time {
	return t0.AddDate(0, 0, days)
}

// PlusMonths 增加时间:单位:月,months < 0,表示时间相减
func PlusMonths(t0 time.Time, months int) time.Time {
	return t0.AddDate(0, months, 0)
}

// PlusYears 增加时间:单位:年,years < 0,表示时间相减
func PlusYears(t0 time.Time, years int) time.Time {
	return t0.AddDate(years, 0, 0)
}

func PlusTime(t0 time.Time, sub int, unit TimeUnit) time.Time {
	switch unit {
	case Year:
		return PlusYears(t0, sub)
	case Month:
		return PlusMonths(t0, sub)
	case Day:
		return PlusDays(t0, sub)
	case Hour:
		return PlusHours(t0, sub)
	case Minute:
		return PlusMinutes(t0, sub)
	case Second:
		return PlusSeconds(t0, sub)
	default:
		return t0
	}
}
