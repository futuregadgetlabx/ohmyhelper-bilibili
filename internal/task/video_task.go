package task

import (
	"context"
	"fmt"
	"math/rand"
	"ohmyhelper-bilibili/internal/delegate"
	"ohmyhelper-bilibili/internal/model"
	"strconv"
	"strings"
)

// regionIds 分类id
var regionIds = []int{
	1,   // 动画
	129, // 舞蹈
	4,   // 游戏
	36,  // 知识
	188, // 科技
	234, // 运动
	223, // 汽车
	160, // 生活
	211, // 美食
	217, // 动物圈
	119, // 鬼畜
	181, // 影视
}

type VideoTask struct {
	ctx context.Context
	d   *delegate.BiliDelegate
}

func NewVideoTask(ctx context.Context, d *delegate.BiliDelegate) *VideoTask {
	return &VideoTask{
		ctx: ctx,
		d:   d,
	}
}

func (v *VideoTask) Run() {
	v.doWatchVideo()
	v.doShareVideo()
}

func (v *VideoTask) Name() string {
	return "视频任务"
}

func (v *VideoTask) getVideo() *model.RegionRank {
	videos, err := v.d.GetTrendVideo(regionIds[rand.Intn(len(regionIds))])
	if err != nil {
		log.WithError(err).Error("获取分区视频失败")
		return nil
	}
	return &videos[0]
}

func (v *VideoTask) doWatchVideo() {
	target := v.getVideo()
	rewardStatus, err := v.d.GetExpRewardStatus()
	if err != nil {
		log.WithError(err).Error("获取每日经验奖励状态失败")
		return
	}

	if rewardStatus.Watch {
		log.Info("今日观看视频任务已完成")
		return
	}

	seconds, err := convertToSeconds(target.Duration)
	if err != nil {
		log.WithError(err).Error("转换视频时长失败")
		return
	}
	playtime := 10 + rand.Intn(seconds-10)
	err = v.d.PlayVideo(target.Bvid, playtime)
	if err != nil {
		log.WithError(err).Error("观看视频失败")
		return
	}
	log.Infof("播放视频[%s]成功,已观看至%d秒", target.Title, playtime)
}

func (v *VideoTask) doShareVideo() {
	target := v.getVideo()
	rewardStatus, err := v.d.GetExpRewardStatus()
	if err != nil {
		log.WithError(err).Error("获取每日经验奖励状态失败")
		return
	}

	if rewardStatus.Share {
		log.Info("今日分享视频任务已完成")
		return
	}

	err = v.d.ShareVideo(target.Bvid)
	if err != nil {
		log.WithError(err).Error("分享视频失败")
		return
	}
	log.Infof("分享视频[%s]成功", target.Title)
}

func convertToSeconds(timeStr string) (int, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid time format")
	}

	minutes, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	seconds, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	return minutes*60 + seconds, nil
}
