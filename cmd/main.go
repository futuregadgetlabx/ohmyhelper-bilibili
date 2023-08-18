package main

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ohmyhelper-bilibili/internal/delegate"
	"ohmyhelper-bilibili/internal/task"
	"ohmyhelper-bilibili/pkg/util"
	"os"
	"path"
	"runtime"
	"strconv"
)

var log *logrus.Entry

func main() {
	ctx := context.Background()
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{
		//以下设置只是为了使输出更美观
		TimestampFormat: "2006-01-02 15:03:04",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			//处理文件名
			fileName := path.Base(frame.File)

			return "", fileName + ":" + strconv.Itoa(frame.Line)
		},
	})
	// 生成UUID
	traceID := uuid.New().String()
	ctx = context.WithValue(ctx, "traceID", traceID)
	biliUserID := os.Getenv("BILIBILI_USER_ID")
	ctx = context.WithValue(ctx, "biliUserID", biliUserID)
	log = logrus.WithFields(logrus.Fields{
		"traceId":    traceID,
		"biliUserID": biliUserID,
	})
	// 使用viper解析yaml配置
	viper.SetConfigFile("conf/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 监听配置文件变更
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("配置文件被修改:", e.Name)
	})
	// 获取配置数据
	dbConfig := viper.Sub("db")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		dbConfig.Get("user"),
		dbConfig.Get("password"),
		dbConfig.Get("host"),
		dbConfig.Get("port"),
		dbConfig.Get("database"))

	// 创建数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	userid := biliUserID
	taskConfig := &delegate.BiliTaskConfig{}
	err = db.Table("task_config").Where("dedeuserid = ?", userid).Find(&taskConfig).Error
	decodeSessdata, _ := hex.DecodeString(taskConfig.Sessdata)
	decodeBiliJct, _ := hex.DecodeString(taskConfig.BiliJct)
	secret := []byte(os.Getenv("SECRET_KEY"))
	taskConfig.Sessdata = string(util.AesDecrypt(decodeSessdata, secret))
	taskConfig.BiliJct = string(util.AesDecrypt(decodeBiliJct, secret))
	ctx = context.WithValue(ctx, "taskConfig", taskConfig)
	if err != nil {
		log.Fatal(err)
	}

	// run task
	task.Run(ctx)

	d, err := db.DB()
	if err != nil {
		return
	}
	defer func(d *sql.DB) {
		err := d.Close()
		if err != nil {
			return
		}
	}(d)
}
