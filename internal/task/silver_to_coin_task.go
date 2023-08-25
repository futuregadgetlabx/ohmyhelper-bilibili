package task

import (
	"context"
	"ohmyhelper-bilibili/internal/delegate"
)

type Silver2CoinTask struct {
	ctx context.Context
	d   *delegate.BiliDelegate
}

func NewSilver2CoinTask(ctx context.Context, d *delegate.BiliDelegate) *Silver2CoinTask {
	return &Silver2CoinTask{
		ctx: ctx,
		d:   d,
	}
}

func (s Silver2CoinTask) Run() {
	wallet, err := s.d.GetLiveWallet()
	if err != nil {
		log.WithError(err).Error("获取直播间钱包信息失败")
		return
	}

	if wallet.Silver < 700 {
		log.Infof("银瓜子余额为[%d]，不足700，无法执行兑换任务", wallet.Silver)
		return
	}

	err = s.d.Silver2Coin()
	if err != nil {
		log.WithError(err).Error("银瓜子兑换硬币失败")
		return
	}
	log.Info("银瓜子兑换硬币成功")
}

func (s Silver2CoinTask) Name() string {
	return "银瓜子兑换硬币"
}
