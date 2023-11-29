package task

import (
	"context"
	"ohmyhelper-bilibili/internal/delegate"
	"time"
)

type DonateGiftTask struct {
	ctx      context.Context
	delegate *delegate.BiliDelegate
}

func NewDonateGiftTask(ctx context.Context, d *delegate.BiliDelegate) *DonateGiftTask {
	return &DonateGiftTask{
		ctx:      ctx,
		delegate: d,
	}
}
func (d *DonateGiftTask) Run() {
	config := d.delegate.Config
	if !config.DonateGift {
		log.Info("未启用赠送即将过期礼物")
		return
	}

	userID, roomID, err := d.getRoomID(config.DonateGiftTarget)
	if err != nil {
		log.WithError(err).Error("获取目标直播间ID失败")
		return
	}

	gifts, err := d.delegate.ListGifts()
	if err != nil {
		log.WithError(err).Error("获取背包礼物失败")
		return
	}
	for _, gift := range gifts.List {
		ddl := time.Now().Add(3 * time.Hour * 24).Unix()
		if ddl > gift.ExpireAt {
			err = d.delegate.DonateGift(userID, roomID, gift.BagID, gift.GiftID, gift.GiftNum)
			if err != nil {
				log.WithError(err).Errorf("为直播间%s赠送礼物失败", roomID)
				return
			}
			log.Infof("为直播间%d赠送了%d个%s", roomID, gift.GiftNum, gift.GiftName)
		}
	}
	log.Info("赠送礼物任务完成")
}

func (d *DonateGiftTask) Name() string {
	return "赠送过期礼物"
}

func (d *DonateGiftTask) getRoomID(target string) (userId, roomId string, err error) {
	if target == "" {
		userId = "287969457"
	}

	liveRoomInfo, err := d.delegate.GetLiveRoomInfo(target)
	if err != nil {
		return "287969457", "11526309", nil
	}

	return userId, liveRoomInfo.RoomID, nil
}
