package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"ohmyhelper-bilibili/internal/delegate"
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

	taskConfig, err := parseConfig()
	if err != nil {
		log.WithError(err).Error("初始化任务配置失败")
		return
	}
	traceID := uuid.New().String()
	ctx = context.WithValue(ctx, "traceID", traceID)
	ctx = context.WithValue(ctx, "biliUserID", taskConfig.Dedeuserid)
	log = logrus.WithFields(logrus.Fields{
		"traceId":    traceID,
		"biliUserID": taskConfig.Dedeuserid,
	})
	ctx = context.WithValue(ctx, "taskConfig", taskConfig)

	// run task
	runner, err := task.NewRunner(ctx)
	if err != nil {
		log.WithError(err).Fatal("初始化任务失败")
	}
	runner.Run(ctx)
	runner.Summery(ctx)
}

func parseConfig() (*delegate.BiliTaskConfig, error) {
	dedeuseridStr := os.Getenv("dedeuserid")
	dedeuserid, err := strconv.Atoi(dedeuseridStr)
	if err != nil {
		log.WithError(err).Errorf("string to int error:  %s", dedeuseridStr)
		return nil, err
	}

	donateCoinsStr := os.Getenv("donateCoins")
	donateCoins, err := strconv.ParseFloat(donateCoinsStr, 10)
	if err != nil {
		log.WithError(err).Errorf("string to float64 error:  %s", donateCoinsStr)
		return nil, err
	}

	reserveCoinsStr := os.Getenv("reserveCoins")
	reserveCoins, err := strconv.Atoi(reserveCoinsStr)
	if err != nil {
		log.WithError(err).Errorf("string to int error:  %s", reserveCoinsStr)
		return nil, err
	}

	autoChargeStr := os.Getenv("autoCharge")
	autoCharge, err := strconv.ParseBool(autoChargeStr)
	if err != nil {
		log.WithError(err).Errorf("string to bool error:  %s", autoChargeStr)
		return nil, err
	}

	donateGiftStr := os.Getenv("donateGift")
	donateGift, err := strconv.ParseBool(donateGiftStr)
	if err != nil {
		log.WithError(err).Errorf("string to bool error:  %s", autoChargeStr)
		return nil, err
	}

	donateGiftTargetStr := os.Getenv("donateGiftTarget")
	donateGiftTarget, err := strconv.Atoi(donateGiftTargetStr)
	if err != nil {
		log.WithError(err).Errorf("string to int error:  %s", donateGiftTargetStr)
		return nil, err
	}

	autoChargeTargetStr := os.Getenv("autoChargeTarget")
	autoChargeTarget, err := strconv.Atoi(autoChargeTargetStr)
	if err != nil {
		log.WithError(err).Errorf("string to int error:  %s", autoChargeTargetStr)
		return nil, err
	}

	followDeveloperStr := os.Getenv("followDeveloper")
	followDeveloper, err := strconv.ParseBool(followDeveloperStr)
	if err != nil {
		log.WithError(err).Errorf("string to bool error:  %s", followDeveloperStr)
		return nil, err
	}
	taskConfig := &delegate.BiliTaskConfig{
		Dedeuserid:         dedeuserid,
		Sessdata:           os.Getenv("sessdata"),
		BiliJct:            os.Getenv("biliJct"),
		DonateCoins:        donateCoins,
		ReserveCoins:       reserveCoins,
		AutoCharge:         autoCharge,
		DonateGift:         donateGift,
		DonateGiftTarget:   donateGiftTarget,
		AutoChargeTarget:   autoChargeTarget,
		DevicePlatform:     os.Getenv("devicePlatform"),
		DonateCoinStrategy: os.Getenv("donateCoinStrategy"),
		UserAgent:          os.Getenv("userAgent"),
		FollowDeveloper:    followDeveloper,
	}
	return taskConfig, nil
}
