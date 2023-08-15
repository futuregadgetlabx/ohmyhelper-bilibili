package runner

import (
	"context"
	"github.com/sirupsen/logrus"
	"ohmyhelper-bilibili/internal/delegate"
)

func Run(ctx context.Context) error {
	traceID := ctx.Value("X-Trace-ID").(string)
	biliUserID := ctx.Value("biliUserID").(string)
	config := ctx.Value("taskConfig").(*delegate.BiliTaskConfig)
	log := logrus.WithField("traceID", traceID).WithField("biliUserID", biliUserID)
	// 获取B站用户详情
	d := delegate.NewDelegate(config, false)
	details, err := d.GetUserDetails()
	if err != nil {
		log.Errorf("获取用户详情失败: %s", err)
		return err
	}
	ctx = context.WithValue(ctx, "details", details)

	// 查询B站硬币日志
	coinsLog, err := d.GetCoinLog()
	if err != nil {
		log.Errorf("获取B站硬币日志失败: %s", err)
		return err
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
	}
	log.Infof("最近一周收入%f个硬币", income)
	log.Infof("最近一周支出%f个硬币", expend)
	return nil
}
