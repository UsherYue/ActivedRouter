package tools

import (
	"fmt"
	"testing"
	//	"time"
)

func TestDatetool(t *testing.T) {
	dateTool := DateTool{}
	begin, _ := dateTool.ParseDatetimeString("2016-02-13 00:00:00")
	end, _ := dateTool.ParseDatetimeString("2016-02-20 00:00:00")
	fmt.Printf("BeginTime:%d\n", begin.Unix())
	fmt.Printf("EndTime:%d\n", end.Unix())
}

func TestGetWeekInterval(t *testing.T) {
	datatTool := DateTool{}
	datatTool.GetLastWeekInterval(1)
	y, m, d := datatTool.TodayYMD()
	weekstr, week := datatTool.TodayWeekday()
	h, m, s := datatTool.NowHMS()
	fmt.Printf("YMD:%d,%d,%d\n", y, m, d)
	fmt.Printf("Week:%s,%d\n", weekstr, week)
	fmt.Printf("HMS:  %d:%d:%d\n", h, m, s)

	begin, end := datatTool.GetLastWeekInterval(week)
	fmt.Printf("Begin_End:%d,%d\n", begin, end)
	strBegin := datatTool.Format(begin, "2006-01-02 15:04:05")
	strEnd := datatTool.Format(end, "2006-01-02 15:04:05")
	fmt.Println(strEnd)
	fmt.Println(strBegin)

	//unix := datatTool.GetNowUnixTimestamp()
	//fmt.Println("Unix:", unix)
	//var tmUnix int64 = time.Now().Unix()
	//var ss string = time.Unix(tmUnix, 0).Format("2006-01-02 15:04:05")
	//fmt.Println(ss)

	beginPre, endPre := datatTool.GetFutureWeekIntervalOffset(week, 1)
	fmt.Println("begin,end:", beginPre, ":", endPre)

	str := datatTool.Format(782166400, "2006-01-02 15:04:05")
	fmt.Println(str)
}
