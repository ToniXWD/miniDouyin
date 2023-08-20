package utils

import (
	"time"
)

func Realurl(s_url string) string {
	r_url := "http://" + ServerIP + ":" + PORT + "/data/" + s_url
	return r_url
}

func I64ToTime(num int64) time.Time {
	// 提取时间戳的秒和微秒部分
	seconds := num / 1000
	microseconds := (num % 1000) * 1000
	// 将 Unix 时间戳转换为 time.Time 类型
	timeObj := time.Unix(seconds, microseconds)

	// 格式化为指定的日期时间格式
	//formattedTime := timeObj.Format("2006-01-02 15:04:05.000000")

	return timeObj
}

func TimeToI64(t time.Time) int64 {
	// 将 time.Time 转换为 Unix 时间戳（秒数）
	seconds := t.Unix()

	// 将秒数转换为毫秒级别的时间戳
	milliseconds := seconds*1000 + int64(t.Nanosecond())/int64(time.Millisecond)

	return milliseconds
}
