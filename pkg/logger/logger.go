package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"strconv"
)

var Logger *logrus.Entry

func Init(ctx context.Context) {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{
		//以下设置只是为了使输出更美观
		TimestampFormat: "2006-01-02 15:03:04",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			//处理文件名
			fileName := path.Base(frame.File)

			return "", fileName + ":" + strconv.Itoa(frame.Line)
		},
	})
	Logger = logrus.WithField("traceID", ctx.Value("traceID")).
		WithField("biliUserID", ctx.Value("biliUserID"))
}
