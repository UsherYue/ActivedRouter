package tools

import (
	"fmt"
	"math"
	"time"
)

var months = map[string]int{
	"January":   1,
	"February":  2,
	"March":     3,
	"April":     4,
	"May":       5,
	"June":      6,
	"July":      7,
	"August":    8,
	"September": 9,
	"October":   10,
	"November":  11,
	"December":  12,
}

var days = map[string]int{
	"Sunday":    0,
	"Monday":    1,
	"Tuesday":   2,
	"Wednesday": 3,
	"Thursday":  4,
	"Friday":    5,
	"Saturday":  6,
}

//时间工具
type DateTool struct {
}

func (this *DateTool) CurrentUnixTimestamp() int64 {
	return time.Now().Unix()
}

func (this *DateTool) FormatUnixTimestamp(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func (this *DateTool) ParseDatetimeString(datetime string) (time.Time, error) {
	timeLayout := "2006-01-02 15:04:05"
	tm, err := time.ParseInLocation(timeLayout, datetime, time.Local)
	return tm, err
}

//fmt.Println(tm.Format("2006-01-0203:04:05PM"))
//fmt.Println(tm.Format("02/01/200615:04:05PM"))
func (this *DateTool) Format(timestamp int64, layoutString string) string {
	tm := this.FormatUnixTimestamp(timestamp)
	return tm.Format(layoutString)
}

//获取今天是星期几
func (this *DateTool) TodayWeekday() (string, int) {
	tmNow := time.Now()
	tmWeekdayStr := tmNow.Weekday().String()
	tmWeekday := days[tmWeekdayStr]
	return tmWeekdayStr, tmWeekday
}

//today y m d
func (this *DateTool) TodayYMD() (int, int, int) {
	tmNow := time.Now()
	tmMonthStr := tmNow.Month().String()
	tmMonth := months[tmMonthStr]
	return tmNow.Year(), tmMonth, tmNow.Day()
}

//get  hour  minutes  second
func (this *DateTool) NowHMS() (int, int, int) {
	tmNow := time.Now()
	return tmNow.Hour(), tmNow.Minute(), tmNow.Second()
}

//convert month string to month
func (this *DateTool) MonthStrToInt(month string) int {
	return months[month]
}

func (this *DateTool) UnixToYMD(unixTimestamp int64) (int, int, int) {
	tm := this.FormatUnixTimestamp(unixTimestamp)
	tmMonthStr := tm.Month().String()
	tmMonth := months[tmMonthStr]
	//tmWeekday:=tm.Weekday()
	return tm.Year(), tmMonth, tm.Day()
}

// 获取指定年月份的天数
func (this *DateTool) GetMonthDays(month, year int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	}
	if year%4 == 0 {
		return 29
	} else {
		return 28
	}
}

//获取当前时间戳
func (this *DateTool) GetNowUnixTimestamp() int64 {
	return time.Now().Local().Unix()
}

//获取某一天的凌晨开始 向前一个周的该天的凌晨 00:00:00
//返回begin  end  timestamp
func (this *DateTool) GetLastWeekInterval(weekday int) (int64, int64) {
	//tmNow := time.Now()
	tmWeekday := weekday //days[tmNow.Weekday().String()]
	tmYear, tmMonth, tmDay := this.TodayYMD()
	//begin date -> end  date
	var endDay, endMonth, endYear = tmDay, tmMonth, tmYear
	//saturday 如果是周六那么 结束日期就是本月
	switch tmWeekday {
	case 6:
		{
			endDay = tmDay
			endMonth = tmMonth
			endYear = tmYear
		}
	default:
		{
			//不是周六计算出 差量days
			reduceDay := tmWeekday + 1
			if tmDay > 6 {
				endDay = tmDay - reduceDay
			} else {
				//上一个月 以及天数如果是1月份那么要进入上一年
				if tmMonth == 1 {
					endMonth = 12
					endYear -= 1
				} else {
					endMonth -= 1
				}
				endMonthDays := this.GetMonthDays(endMonth, endYear)
				endDay = endMonthDays - int(math.Abs(float64(endDay-reduceDay)))
			}
		}
	}
	endDateStr := fmt.Sprintf("%d-%02d-%02d 00:00:00", endYear, endMonth, endDay)
	endTime, _ := this.ParseDatetimeString(endDateStr)
	endTimestamp := endTime.Unix()
	beginTimestamp := endTimestamp - 24*7*3600
	return beginTimestamp, endTimestamp - 1
}

//获取多个偏移 n个 23:59:59 － 00:00:00
func (this *DateTool) GetLastWeekIntervalOffset(weekday int, offset int64) (int64, int64) {
	var begin, end int64 = this.GetLastWeekInterval(weekday)
	begin -= offset * 24 * 3600 * 7
	end -= offset * 24 * 3600 * 7
	return begin, end
}

//获取多个偏移 n个 23:59:59 － 00:00:00
func (this *DateTool) GetFutureWeekIntervalOffset(weekday int, offset int64) (int64, int64) {
	var begin, end int64 = this.GetLastWeekInterval(weekday)
	begin += offset * 24 * 3600 * 7
	end += offset * 24 * 3600 * 7
	return begin, end
}
