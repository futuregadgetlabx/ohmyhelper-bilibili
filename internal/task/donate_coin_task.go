package task

import (
	"context"
	"math/rand"
	"ohmyhelper-bilibili/internal/delegate"
	"ohmyhelper-bilibili/internal/model"
)

type DonateCoinTask struct {
	ctx context.Context
	d   *delegate.BiliDelegate
}

func NewDonateCoinTask(ctx context.Context, d *delegate.BiliDelegate) *DonateCoinTask {
	return &DonateCoinTask{
		ctx: ctx,
		d:   d,
	}
}

func (d *DonateCoinTask) Run() {
	details := d.ctx.Value("details").(*model.BiliUserDetails)
	coins := details.Coins
	config := d.d.Config
	if coins <= float64(config.ReserveCoins) {
		log.Errorf("账户余额不足%delegate，无法进行投币任务", config.ReserveCoins)
		return
	}

	if details.Level >= 6 {
		log.Info("账号已到达6级，取消执行投币任务")
		return
	}

	expToday, err := d.d.GetCoinExpToday()
	if err != nil {
		log.WithError(err).Error("获取今日投币经验失败")
		return
	}

	needCoins := config.DonateCoins - expToday/10
	if needCoins <= 0 {
		log.Info("今日投币任务已完成")
		return
	}

	trendVideo, err := d.d.GetTrendVideo(regionIds[rand.Intn(len(regionIds))])
	if err != nil {
		log.WithError(err).Error("获取分区视频排行失败")
		return
	}

	for _, v := range trendVideo {
		if needCoins <= 0 {
			log.Info("今日投币任务已完成")
			return
		}
		err = d.d.DonateCoin(v.Bvid, 1, 1)
		if err != nil {
			log.WithError(err).Errorf("为视频[%s]投币失败", v.Title)
			return
		}
		log.Infof("为视频[%s]投币成功", v.Title)
		needCoins--
	}
}

func (d *DonateCoinTask) Name() string {
	return "投币"
}
