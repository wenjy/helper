package helper

import "time"

const TimeLayout = "2006-01-02 15:04:05"
const timezone = "Asia/Shanghai"

var cstZone *time.Location

func init() {
	var err error
	cstZone, err = time.LoadLocation(timezone)
	if err != nil {
		cstZone = time.FixedZone("CST", 8*3600)
	}
}

// CustomFormatTime 格式时间为 `Y-m-d H:i:s`
func CustomFormatTime(t time.Time) string {
	return t.In(cstZone).Format(TimeLayout)
}

func Strtotime(date string) (time.Time, error) {
	return time.ParseInLocation(TimeLayout, date, cstZone)
}

func NowTime() time.Time {
	return time.Now().In(cstZone)
}

func DateTime() string {
	return NowTime().Format(TimeLayout)
}
