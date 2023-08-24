package task

import (
	"context"
	"ohmyhelper-bilibili/internal/delegate"
	"os"
	"strconv"
	"time"
)

func doDonateGift(ctx context.Context) {
	config := ctx.Value("taskConfig").(*delegate.BiliTaskConfig)
	if !config.DonateGift {
		log.Infoln("未启用赠送即将过期礼物")
		return
	}

	userID, roomID, err := getRoomID(config)
	if err != nil {
		log.WithError(err).Errorln("获取目标直播间ID失败")
		return
	}

	gifts, err := d.ListGifts()
	if err != nil {
		log.WithError(err).Errorln("获取背包礼物失败")
		return
	}
	for _, gift := range gifts.List {
		ddl := time.Now().Add(3 * time.Hour * 24).Unix()
		if ddl > gift.ExpireAt {
			err = d.DonateGift(userID, roomID, gift.BagID, gift.GiftID, gift.GiftNum)
			if err != nil {
				log.WithError(err).Errorf("为直播间%d赠送礼物失败\n", roomID)
				return
			}
			log.Infof("为直播间%d赠送了%d个%s\n", roomID, gift.GiftNum, gift.GiftName)
		}
	}
	log.Infoln("赠送礼物任务完成")
}

func getRoomID(config *delegate.BiliTaskConfig) (userID, roomID int, err error) {
	userID = config.DonateGiftTarget
	if userID == 0 {
		userID, err = strconv.Atoi(os.Getenv("AUTHOR_ID"))
		if err != nil {
			return 0, 0, err
		}
	}

	liveRoomInfo, err := d.GetLiveRoomInfo(userID)
	if err != nil {
		return 0, 0, err
	}

	return userID, liveRoomInfo.RoomID, nil
}
