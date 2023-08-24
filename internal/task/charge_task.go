package task

import (
	"context"
	"ohmyhelper-bilibili/internal/delegate"
	"ohmyhelper-bilibili/internal/model"
	"os"
	"strconv"
)

func doChargeTask(ctx context.Context) {
	details := ctx.Value("details").(*model.BiliUserDetails)
	config := ctx.Value("taskConfig").(*delegate.BiliTaskConfig)

	vipType := details.Vip.Type
	if vipType != 2 {
		log.Infoln("账号非年费大会员，停止执行充电任务")
		return
	}

	chargeInfo, err := d.GetChargeInfo()
	if err != nil {
		log.WithError(err).Errorln("获取充电信息失败")
		return
	}

	couponBalance := chargeInfo.BpWallet.CouponBalance
	if couponBalance < 2 {
		log.Infoln("B币券余额不足，停止执行充电任务")
		return
	}

	target, err := getChargeTarget(config)
	if err != nil {
		log.WithError(err).Errorln("获取充电对象信息失败")
		return
	}

	chargeResponse, err := d.DoCharge(couponBalance, target)
	if err != nil {
		log.WithError(err).Errorln("执行充电失败")
		return
	}

	log.Infof("为账号%d充电成功\n", target)
	orderNo := chargeResponse.OrderNo
	_ = d.DoChargeComment(orderNo)
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
