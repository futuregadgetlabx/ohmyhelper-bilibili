package task

import (
	"context"
	"ohmyhelper-bilibili/internal/model"
)

// BP B币券
const BP = 1

func doBigVipPrivilegeTask(ctx context.Context) {
	details := ctx.Value("details").(*model.BiliUserDetails)
	if details.Vip.Type == 0 || details.Vip.Status != 1 {
		log.Infoln("该账号非大会员，取消执行领取大会员B币券")
	}

	err := d.GetVipReward(BP)
	if err != nil {
		log.WithError(err).Errorln("领取大会员B币券失败")
		return
	}

	log.Infoln("领取大会员B币券成功")
}
