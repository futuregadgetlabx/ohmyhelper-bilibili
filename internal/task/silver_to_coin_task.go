package task

func doSilver2Coin() {
	wallet, err := d.GetLiveWallet()
	if err != nil {
		log.WithError(err).Errorln("获取直播间钱包信息失败")
		return
	}

	if wallet.Silver < 700 {
		log.Infof("银瓜子余额为[%d]，不足700，无法执行兑换任务\n", wallet.Silver)
		return
	}

	err = d.Silver2Coin()
	if err != nil {
		log.WithError(err).Errorln("银瓜子兑换硬币失败")
		return
	}
	log.Infoln("银瓜子兑换硬币成功")
}
