package task

import (
	"context"
	"ohmyhelper-bilibili/internal/delegate"
	"ohmyhelper-bilibili/internal/model"
	"os"
	"strconv"
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
	config := c.ctx.Value("taskConfig").(*delegate.BiliTaskConfig)

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

	target, err := getChargeTarget(config)
	if err != nil {
		log.WithError(err).Error("获取充电对象信息失败")
		return
	}

	chargeResponse, err := c.d.DoCharge(couponBalance, target)
	if err != nil {
		log.WithError(err).Error("执行充电失败")
		return
	}

	log.Infof("为账号%d充电成功", target)
	orderNo := chargeResponse.OrderNo
	_ = c.d.DoChargeComment(orderNo)
}

func (c *ChargeTask) Name() string {
	return "充电"
}

func getChargeTarget(config *delegate.BiliTaskConfig) (userID int, err error) {
	userID = config.AutoChargeTarget
	if userID == 0 {
		userID, err = strconv.Atoi(os.Getenv("AUTHOR_ID"))
		if err != nil {
			return 0, err
		}
	}

	return userID, nil
}
