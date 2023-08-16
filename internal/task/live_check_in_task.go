package task

func doLiveCheckIn() {
	checkInResp, err := d.LiveCheckIn()
	if err != nil {
		log.WithError(err).Error("直播签到失败")
		return
	}
	log.Infof("直播签到成功，本次获得%s,%s", checkInResp.Text, checkInResp.SpecialText)
}
