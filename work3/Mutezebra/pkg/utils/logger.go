package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"runtime"
	"time"
)

var LogrusObj *logrus.Logger

// InitLog 初始化logger
func InitLog() {
	// 如果已经初始化过则更新日志输出路径
	if LogrusObj != nil {
		src, err := setOutputFile()
		if err != nil {
			log.Println(err)
		}
		LogrusObj.Out = src
		fmt.Println(src)
		return
	}

	// 日志对象的实例化
	logger := logrus.New()
	src, err := setOutputFile()
	if err != nil {
		log.Println(err)
	}
	logger.Out = src
	logger.SetLevel(logrus.DebugLevel)

	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File)
			return frame.Function, fmt.Sprintf("%s:%d", fileName, frame.Line)
		},
	})
	logger.Infoln("this is a info")
	LogrusObj = logger
}

// setOutFile 设置日志输出路径
func setOutputFile() (*os.File, error) {
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "\\log\\"
	}
	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(logFilePath, 0777)
		if err != nil {
			return nil, err
		}
	}
	logFileName := now.Format("2006-01-02") + ".log"

	fileName := logFilePath + logFileName
	_, err = os.Stat(fileName)
	if err != nil {
		_, err := os.Create(fileName)
		if err != nil {
			return nil, err
		}
	}
	src, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	return src, nil
}
