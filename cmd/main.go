package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"ohmyhelper-bilibili/internal/task"
	"os"
	"path"
	"runtime"
	"strconv"
)

var log *logrus.Entry

func init() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{
		//以下设置只是为了使输出更美观
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			//处理文件名
			fileName := path.Base(frame.File)

			return "", fileName + ":" + strconv.Itoa(frame.Line)
		},
	})
	f, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.WithError(err).Panic("open log file failed")
	}
	writers := []io.Writer{
		f,
		os.Stderr,
	}
	multiWriter := io.MultiWriter(writers...)
	logrus.SetOutput(multiWriter)
}
func main() {
	ctx := context.Background()
	// 生成UUID
	traceId := uuid.New().String()

	ctx = context.WithValue(ctx, "traceId", traceId)
	dedeuserid := os.Getenv("DEDE_USER_ID")
	ctx = context.WithValue(ctx, "biliUserID", dedeuserid)
	log = logrus.WithFields(logrus.Fields{
		"traceId":    traceId,
		"biliUserId": dedeuserid,
	})
	ctx = context.WithValue(ctx, "taskConfig", dedeuserid)

	// run task
	runner, err := task.NewRunner(ctx)
	if err != nil {
		log.WithError(err).Fatal("初始化任务失败")
	}
	runner.Run(ctx)
	runner.Summary(ctx)
}
