package task

func doCoinsLog() {
	// 查询B站硬币日志
	coinsLog, err := d.GetCoinLog()
	if err != nil {
		log.WithError(err).Error("获取B站硬币日志失败")
	}
	// 逐条打印硬币日志
	income := 0.0
	expend := 0.0
	log.Infof("最近一周共产生%d条变更日志", len(coinsLog.List))
	for _, clog := range coinsLog.List {
		if clog.Delta > 0 {
			income += clog.Delta
		} else {
			expend += clog.Delta
		}
		log.Infof("%s--%s：%f", clog.Time, clog.Reason, clog.Delta)
	}
	log.Infof("最近一周收入%f个硬币", income)
	log.Infof("最近一周支出%f个硬币", expend)
}
