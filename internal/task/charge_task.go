package task

import (
	"context"
	"ohmyhelper-bilibili/internal/delegate"
	"ohmyhelper-bilibili/internal/model"
)

type ChargeTask struct {
	ctx context.Context
	d   *delegate.BiliDelegate
}

func NewChargeTask(ctx context.Context, d *delegate.BiliDelegate) *ChargeTask {
	return &ChargeTask{
		ctx: ctx,
		d:   d,
	}
}

func (c *ChargeTask) Run() {
	details := c.ctx.Value("details").(*model.BiliUserDetails)
	config := c.d.Config

	vipType := details.Vip.Type
	if vipType != 2 {
		log.Info("账号非年费大会员，停止执行充电任务")
		return
	}

	chargeInfo, err := c.d.GetChargeInfo()
	if err != nil {
		log.WithError(err).Error("获取充电信息失败")
		return
	}

	couponBalance := chargeInfo.BpWallet.CouponBalance
	if couponBalance < 2 {
		log.Info("B币券余额不足，停止执行充电任务")
		return
	}

	chargeTarget := config.AutoChargeTarget
	if chargeTarget == "" {
		chargeTarget = "287969457"
	}

	chargeResponse, err := c.d.DoCharge(couponBalance, chargeTarget)
	if err != nil {
		log.WithError(err).Error("执行充电失败")
		return
	}

	log.Infof("为账号%d充电成功", chargeTarget)
	orderNo := chargeResponse.OrderNo
	_ = c.d.DoChargeComment(orderNo)
}

func (c *ChargeTask) Name() string {
	return "充电"
}
