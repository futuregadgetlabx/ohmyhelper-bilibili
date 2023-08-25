package task

import (
	"context"
	"ohmyhelper-bilibili/internal/delegate"
	"ohmyhelper-bilibili/internal/model"
)

// BP B币券
const BP = 1

type BigVipPrivilegeTask struct {
	ctx context.Context
	d   *delegate.BiliDelegate
}

func NewBigVipPrivilegeTask(ctx context.Context, d *delegate.BiliDelegate) *BigVipPrivilegeTask {
	return &BigVipPrivilegeTask{
		ctx: ctx,
		d:   d,
	}
}

func (b *BigVipPrivilegeTask) Run() {
	details := b.ctx.Value("details").(*model.BiliUserDetails)
	if details.Vip.Type == 0 || details.Vip.Status != 1 {
		log.Info("该账号非大会员，取消执行领取大会员B币券")
	}

	err := b.d.GetVipReward(BP)
	if err != nil {
		log.WithError(err).Errorln("领取大会员B币券失败")
		return
	}

	log.Info("领取大会员B币券成功")
}

func (b *BigVipPrivilegeTask) Name() string {
	return "大会员特权"
}
