package task

import (
	"context"
	"github.com/sirupsen/logrus"
	"ohmyhelper-bilibili/internal/delegate"
)

var log *logrus.Entry
var d *delegate.BiliDelegate

func Run(ctx context.Context) {
	traceID := ctx.Value("traceID").(string)
	biliUserID := ctx.Value("biliUserID").(string)
	config := ctx.Value("taskConfig").(*delegate.BiliTaskConfig)
	log = logrus.WithField("traceID", traceID).WithField("biliUserID", biliUserID)

	// 获取B站用户详情
	d = delegate.NewDelegate(config, false)
	details, err := d.GetUserDetails()
	if err != nil {
		log.WithError(err).Error("获取用户详情失败")
		return
	}
	ctx = context.WithValue(ctx, "details", details)

	// 执行任务
	doTask(ctx)
}

func doTask(ctx context.Context) {
	doCoinsLog()
	doWatchVideo()
	doShareVideo()
	doBigVipPrivilegeTask(ctx)
	doDonateCoin(ctx)
	doSilver2Coin()
	doLiveCheckIn()
	doDonateGift(ctx)
	doChargeTask(ctx)
}
