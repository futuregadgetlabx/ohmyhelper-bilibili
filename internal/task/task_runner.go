package task

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"ohmyhelper-bilibili/internal/delegate"
	"os"
	"time"
)

var log *logrus.Entry

type Runner struct {
	ctx   context.Context
	tasks []Task
	d     *delegate.BiliDelegate
}

func NewRunner(ctx context.Context) (*Runner, error) {
	traceId := ctx.Value("traceId").(string)
	dedeuserid := os.Getenv("DEDE_USER_ID")
	log = logrus.WithField("traceId", traceId).WithField("biliUserID", dedeuserid)

	// 获取B站用户详情
	d := delegate.NewDelegate(false)
	details, err := d.GetUserDetails()
	if err != nil {
		return nil, errors.Wrap(err, "cookie过期")
	}
	ctx = context.WithValue(ctx, "details", details)

	coinsLogTask := NewCoinsLogTask(ctx, d)
	videoTask := NewVideoTask(ctx, d)
	bigVipPrivilegeTask := NewBigVipPrivilegeTask(ctx, d)
	donateCoinTask := NewDonateCoinTask(ctx, d)
	donateGiftTask := NewDonateGiftTask(ctx, d)
	silver2CoinTask := NewSilver2CoinTask(ctx, d)
	liveCheckInTask := NewLiveCheckInTask(ctx, d)
	chargeTask := NewChargeTask(ctx, d)
	tasks := []Task{
		coinsLogTask,
		videoTask,
		bigVipPrivilegeTask,
		donateCoinTask,
		donateGiftTask,
		silver2CoinTask,
		liveCheckInTask,
		chargeTask,
	}
	return &Runner{
		ctx:   ctx,
		tasks: tasks,
		d:     d,
	}, nil
}

func (r *Runner) Run(ctx context.Context) {
	for _, task := range r.tasks {
		log.Infof("====🤖执行任务：%s🤖====", task.Name())
		task.Run()
		time.Sleep(3 * time.Second)
	}
	log.Info("=====🌟任务全部完成🌟====")
}

func (r *Runner) Summary(ctx context.Context) {
	details, err := r.d.GetUserDetails()
	if err != nil {
		log.WithError(err).Error("任务总结异常")
	}
	expToday, err := r.d.GetCoinExpToday()
	if err != nil {
		log.WithError(err).Error("获取当日投币经验失败")
	}
	rewardStatus, err := r.d.GetExpRewardStatus()
	if err != nil {
		log.WithError(err).Error("获取当日奖励状态失败")
	}
	if rewardStatus.Share {
		expToday += 5
	}
	if rewardStatus.Watch {
		expToday += 5
	}
	log.Infof("🍓今日已获得%d点经验", int(expToday))
	levelExp := details.LevelExp
	currentExp := levelExp.CurrentExp
	log.Infof("🍇当前经验：%d", currentExp)
	if levelExp.CurrentLevel < 6 {
		diff := levelExp.NextExp - levelExp.CurrentExp
		upgradeDays := (diff / int(expToday)) + 1
		log.Infof("🍖按照当前进度，升级到Lv%d还需要: %d天", levelExp.CurrentLevel+1, upgradeDays)
	}
}
