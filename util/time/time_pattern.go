package time

import (
	"strings"
	"time"
)

// 标准年月格式
const (
	//NormalYear 年格式：yyyy
	NormalYear = "yyyy"

	//NormalMonth 年月格式：yyyy-MM
	NormalMonth = "yyyy-MM"

	//NormalDate 标准日期格式：yyyy-MM-dd
	NormalDate = time.DateOnly

	//NormalHour 标准日期格式：yyyy-MM-dd HH
	NormalHour = "yyyy-MM-dd HH"

	//NormalMinute 标准日期格式：yyyy-MM-dd HH:mm
	NormalMinute = "yyyy-MM-dd HH:mm"

	//NormalDatetime 标准日期时间格式，精确到秒yyyy-MM-dd HH:mm:ss
	NormalDatetime = time.DateTime

	//NormalDatetimeMS 标准日期时间格式，精确到毫秒yyyy-MM-dd HH:mm:ss.SSS
	NormalDatetimeMS = "yyyy-MM-dd HH:mm:ss.SSS"
)

// 简单年月格式
const (
	//SimpleMonth 简单年月格式：yyyyMM
	SimpleMonth = "yyyyMM"

	//SimpleDate 简单年月日格式：yyyyMMdd
	SimpleDate = "yyyyMMdd"

	//SimpleHour 标准日期格式：yyyyMMddHH
	SimpleHour = "yyyyMMddHH"

	//SimpleMinute 标准日期格式：yyyyMMddHHmm
	SimpleMinute = "yyyyMMddHHmm"

	//SimpleDatetime 标准日期时间格式，精确到秒yyyyMMddHHmmss
	SimpleDatetime = "yyyyMMddHHmmss"

	//SimpleDatetimeMS 标准日期时间格式，精确到毫秒yyyyMMdd HH:mm:ssSSS
	SimpleDatetimeMS = "yyyyMMddHHmmssSSS"
)

// 中文年月格式
const (
	//ChineseYear 简单年月格式：yyyy年
	ChineseYear = "yyyy年"

	//ChineseMonth 简单年月格式：yyyy年MM月
	ChineseMonth = "yyyy年MM月"

	//ChineseDate 简单年月日格式：yyyy年MM月dd日
	ChineseDate = "yyyy年MM月dd日"

	//ChineseHour 标准日期格式：yyyy年MM月dd日HH
	ChineseHour = "yyyy年MM月dd日HH时"

	//ChineseMinute 标准日期格式：yyyy年MM月dd日HHmm
	ChineseMinute = "yyyy年MM月dd日HH时mm分"

	//ChineseDatetime 标准日期时间格式，精确到秒yyyy年MM月dd日HHmmss
	ChineseDatetime = "yyyy年MM月dd日HH时mm分ss秒"
)

// UTC时间 yyyy-MM-dd'T'HH:mm:ss.SSS'Z'
const (
	UTC_MS_PATTERN = "yyyy-MM-dd'T'HH:mm:ss.SSS'Z'"

	UTC_MS_WITH_ZONE_OFFSET_PATTERN = "yyyy-MM-dd'T'HH:mm:ss.SSSZ"

	UTC_MS_WITH_XXX_OFFSET_PATTERN = "yyyy-MM-dd'T'HH:mm:ss.SSSXXX"

	UTC_MS_WITH_TIME_PATTERN = time.RFC3339
)

func patternToFormat(pattern string) string {
	pattern = strings.ReplaceAll(pattern, "yyyy", "2006")
	pattern = strings.ReplaceAll(pattern, "MM", "01")
	pattern = strings.ReplaceAll(pattern, "dd", "02")
	pattern = strings.ReplaceAll(pattern, "HH", "15")
	pattern = strings.ReplaceAll(pattern, "mm", "04")
	pattern = strings.ReplaceAll(pattern, "ss", "05")
	pattern = strings.ReplaceAll(pattern, "SSS", "000")
	pattern = strings.ReplaceAll(pattern, "XXX", "000")
	return pattern
}
