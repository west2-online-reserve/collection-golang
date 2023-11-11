package log

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

var workDir string

func InitLog() {
	if LogrusObj != nil {
		src, err := setOutPutFile()
		if err != nil {
			log.Println(err)
		}
		LogrusObj.Out = src
		return
	}
	workdir, _ := os.Getwd()
	workDir = workdir
	logger := logrus.New()
	src, err := setOutPutFile()
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
	logger.Infoln("logs init success!")
	LogrusObj = logger
}

func setOutPutFile() (*os.File, error) {
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) {
		err = os.Mkdir(logFilePath, 0777)
		if err != nil {
			return nil, err
		}
	}
	logFileName := now.Format("2006-01-02") + ".log"
	fileName := logFilePath + "/" + logFileName
	_, err = os.Stat(fileName)
	if err != nil {
		_, err = os.Create(fileName)
		if err != nil {
			return nil, err
		}
	}
	src, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	return src, err
}
