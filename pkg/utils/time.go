package utils

import "time"

var cstSh, _ = time.LoadLocation("Asia/Shanghai")

func GetLocalDate() string {
	return time.Now().In(cstSh).Local().Format("2006-01-02")
}

func GetLocalDateTime() string {
	return time.Now().In(cstSh).Local().Format("2006-01-02 15:04:05")
}
func GetExpireDateTime() string {
	return time.Now().In(cstSh).Add(time.Hour * 24).Local().Format("2006-01-02 15:04:05")
}
func TranslateTime(time1 string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", time1, cstSh)
}
func GetLocalTime() time.Time {
	return time.Now().In(cstSh)
}
