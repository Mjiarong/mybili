package utils

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var Logger *logrus.Logger

func Init() {
	scr, err := os.OpenFile(os.Getenv("LOG_FILE_PATH"), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("err", err)
	}

	Logger = logrus.New()
	Logger.SetLevel(logrus.DebugLevel)
	Logger.SetReportCaller(true)
	Logger.Out = scr

	// 设置 file-rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		os.Getenv("LOG_FILE_PATH")+"%Y%m%d.log",
		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		//rotatelogs.WithLinkName(linkName),
	)
	if err != nil {
		fmt.Println("err", err)
	}

	// 为不同级别设置不同的输出目的
	writeMap := lfshook.WriterMap{
		logrus.PanicLevel: logWriter,
		logrus.FatalLevel: logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.InfoLevel:  logWriter,
		logrus.DebugLevel: logWriter,
	}

	hook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	Logger.AddHook(hook)

	//Logger.WithFields(logrus.Fields{}).Infof("Logger Init Success %v", s)
	Logger.Debugln("Logger Init Success")
}
