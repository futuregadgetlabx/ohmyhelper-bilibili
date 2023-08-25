package task

import (
	"context"
	"ohmyhelper-bilibili/internal/delegate"
)

type CoinsLogTask struct {
	ctx context.Context
	d   *delegate.BiliDelegate
}

func NewCoinsLogTask(ctx context.Context, d *delegate.BiliDelegate) *CoinsLogTask {
	return &CoinsLogTask{
		ctx: ctx,
		d:   d,
	}
}

func (c *CoinsLogTask) Run() {
	// 查询B站硬币日志
	coinsLog, err := c.d.GetCoinLog()
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
		log.Infof("%s--%s：%.2f", clog.Time, clog.Reason, clog.Delta)
	}
	log.Infof("最近一周收入%.2f个硬币", income)
	log.Infof("最近一周支出%.2f个硬币", expend)
}

func (c *CoinsLogTask) Name() string {
	return "硬币日志"
}
