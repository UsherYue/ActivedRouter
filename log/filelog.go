//ActivedRouter
//Author:usher.yue
//Amail:usher.yue@gmail.com
//TencentQQ:4223665

package log

import (
	"log"
	"os"
)

type MyLogger struct {
	//log
	loger      *log.Logger
	fileHandle *os.File
}

//正常日志输出
func (this *MyLogger) LogOut(msg ...interface{}) {
	this.loger.Println(msg)

}

//打印错误
func (this *MyLogger) LogErr(msg ...interface{}) {
	this.loger.Fatalln(msg)
}

//创建日志
func NewLogger(logfile string) *MyLogger {
	//初始化log
	file, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		os.Exit(-1)
	}
	loger := log.New(file, "\r\n", log.Ldate|log.Ltime|log.Lshortfile)
	obj := &MyLogger{}
	obj.fileHandle = file
	obj.loger = loger
	return obj
}

//关闭文件
func (this *MyLogger) Close() {
	this.fileHandle.Close()
}

//每天记录一次日志......
func startLogService() {
	log.Printf("正在启动日志记录服务,位于%s目录下.......\n", "logs")
}
