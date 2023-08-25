package task

import (
	"context"
	"ohmyhelper-bilibili/internal/delegate"
)

type LiveCheckInTask struct {
	ctx context.Context
	d   *delegate.BiliDelegate
}

func NewLiveCheckInTask(ctx context.Context, d *delegate.BiliDelegate) *LiveCheckInTask {
	return &LiveCheckInTask{
		ctx: ctx,
		d:   d,
	}
}

func (l LiveCheckInTask) Run() {
	checkInResp, err := l.d.LiveCheckIn()
	if err != nil {
		log.WithError(err).Error("直播签到失败")
		return
	}
	log.Infof("直播签到成功，本次获得%s,%s", checkInResp.Text, checkInResp.SpecialText)
}

func (l LiveCheckInTask) Name() string {
	return "直播签到"
}
