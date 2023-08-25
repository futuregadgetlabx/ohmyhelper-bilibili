package task

import (
	"context"
	"ohmyhelper-bilibili/internal/delegate"
	"os"
	"strconv"
	"time"
)

type DonateGiftTask struct {
	ctx context.Context
	d   *delegate.BiliDelegate
}

func NewDonateGiftTask(ctx context.Context, d *delegate.BiliDelegate) *DonateGiftTask {
	return &DonateGiftTask{
		ctx: ctx,
		d:   d,
	}
}
func (d *DonateGiftTask) Run() {
	config := d.ctx.Value("taskConfig").(*delegate.BiliTaskConfig)
	if !config.DonateGift {
		log.Info("未启用赠送即将过期礼物")
		return
	}

	userID, roomID, err := d.getRoomID(config)
	if err != nil {
		log.WithError(err).Error("获取目标直播间ID失败")
		return
	}

	gifts, err := d.d.ListGifts()
	if err != nil {
		log.WithError(err).Error("获取背包礼物失败")
		return
	}
	for _, gift := range gifts.List {
		ddl := time.Now().Add(3 * time.Hour * 24).Unix()
		if ddl > gift.ExpireAt {
			err = d.d.DonateGift(userID, roomID, gift.BagID, gift.GiftID, gift.GiftNum)
			if err != nil {
				log.WithError(err).Errorf("为直播间%d赠送礼物失败", roomID)
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

func (d *DonateGiftTask) getRoomID(config *delegate.BiliTaskConfig) (userID, roomID int, err error) {
	userID = config.DonateGiftTarget
	if userID == 0 {
		userID, err = strconv.Atoi(os.Getenv("AUTHOR_ID"))
		if err != nil {
			return 0, 0, err
		}
	}

	liveRoomInfo, err := d.d.GetLiveRoomInfo(userID)
	if err != nil {
		return 0, 0, err
	}

	return userID, liveRoomInfo.RoomID, nil
}
