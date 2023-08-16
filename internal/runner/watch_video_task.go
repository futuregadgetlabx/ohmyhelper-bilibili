package runner

import (
	"fmt"
	"math/rand"
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

func runWatchVideoTask() {
	rewardStatus, err := d.GetExpRewardStatus()
	if err != nil {
		log.WithError(err).Error("获取每日经验奖励状态失败")
		return
	}

	if rewardStatus.Watch {
		//log.Info("今日观看视频任务已完成")
		//return
	}

	videos, err := d.GetTrendVideo(regionIds[rand.Intn(len(regionIds))])
	if err != nil {
		log.WithError(err).Error("获取分区视频失败")
		return
	}
	v := videos[0]
	seconds, err := convertToSeconds(v.Duration)
	if err != nil {
		log.WithError(err).Error("转换视频时长失败")
		return
	}
	playtime := 10 + rand.Intn(seconds-10)
	err = d.PlayVideo(v.Bvid, playtime)
	if err != nil {
		log.WithError(err).Error("观看视频失败")
		return
	}
	log.Infof("播放视频[%s]成功,已观看至%d秒", v.Title, playtime)
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
