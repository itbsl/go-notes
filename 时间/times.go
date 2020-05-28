package times

import (
	"fmt"
	"time"
)

type Times struct {
	time.Time
}

func Now() Times {
	return Times{time.Now()}
}

//日期
func (this *Times) Date() string {
	return fmt.Sprintf("%d-%02d-%02d",
		this.Year(), this.Month(), this.Day())
}

//日期时间
func (this *Times) DateTime() string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		this.Year(), this.Month(), this.Day(),
		this.Hour(), this.Minute(), this.Second())
}

//时间戳(秒)
func (this *Times) UnixSec() int64 {
	return this.Unix()
}

//时间戳(毫秒)
func (this *Times) UnixMsec() int64 {
	return this.UnixNano() / 1e6
}

//时间戳(纳秒)
func (this *Times) UnixNanosec() int64 {
	return this.UnixNano()
}

//当前日期 格式:YYYY-mm-dd
func Date() string {
	now := time.Now()
	return fmt.Sprintf("%d-%02d-%02d",
		now.Year(), now.Month(), now.Day())
}

//当前日期时间 格式:YYYY-mm-dd HH:ii:ss
func DateTime() string {
	now := time.Now()
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())
}

//当前时间戳(秒)
func UnixSec() int64 {
	return time.Now().Unix()
}

//当前时间戳(毫秒)
func UnixMsec() int64 {
	return time.Now().UnixNano() / 1e6
}

//当前时间戳(纳秒)
func UnixNanosec() int64 {
	return time.Now().UnixNano()
}

//时间戳转日期
func UnixToDate(unix int64) string {
	onTime := time.Unix(unix, 0)
	return fmt.Sprintf("%d-%02d-%02d",
		onTime.Year(), onTime.Month(), onTime.Day())
}

//时间戳转日期时间
func UnixToDateTime(unix int64) string {
	onTime := time.Unix(unix, 0)
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",
		onTime.Year(), onTime.Month(), onTime.Day(),
		onTime.Hour(), onTime.Minute(), onTime.Second())
}

//日期转时间戳
func DateToUnix(date string) int64 {
	oneTime, _ := time.ParseInLocation("2006-01-02", date, time.Local)
	return oneTime.Unix()
}

//日期时间转时间戳
func DateTimeToUnix(dateTime string) int64 {
	oneTime, _ := time.Parse("2006-01-02 15:04:05", dateTime)
	return oneTime.Unix()
}
