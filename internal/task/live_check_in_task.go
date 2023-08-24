package task

func doLiveCheckIn() {
	checkInResp, err := d.LiveCheckIn()
	if err != nil {
		log.WithError(err).Errorln("直播签到失败")
		return
	}
	log.Infof("直播签到成功，本次获得%s,%s\n", checkInResp.Text, checkInResp.SpecialText)
}
