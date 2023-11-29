package delegate

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var taskConfig biliTaskConfig

type biliTaskConfig struct {
	Dedeuserid         string  `json:"dedeUserId"  env:"DEDE_USER_ID" envDefault:""`
	Sessdata           string  `json:"sessdata"  env:"SESSDATA" envDefault:""`
	BiliJct            string  `json:"biliJct" env:"BILI_JCT" envDefault:""`
	DonateCoins        float64 `json:"donateCoins"  env:"DONATE_COINS" envDefault:"0"`
	ReserveCoins       int     `json:"reserveCoins"  env:"RESERVE_COINS" envDefault:"0"`
	AutoCharge         bool    `json:"autoCharge" env:"AUTO_CHARGE" envDefault:"false"`
	DonateGift         bool    `json:"donateGift"  env:"DONATE_GIFT" envDefault:"false"`
	DonateGiftTarget   string  `json:"donateGiftTarget"  env:"DONATE_GIFT_TARGET" envDefault:""`
	AutoChargeTarget   string  `json:"autoChargeTarget" env:"AUTO_CHARGE_TARGET" envDefault:""`
	DonateCoinStrategy string  `json:"donateCoinStrategy" env:"DONATE_COIN_STRATEGY" envDefault:"1"`
	UserAgent          string  `json:"userAgent"  env:"USER_AGENT" envDefault:"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"`
	FollowDeveloper    bool    `json:"followDeveloper"  env:"FOLLOW_DEVELOPER" envDefault:"true"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Error("Can not read env from file system, please check the right this program owned.")
	}

	taskConfig = biliTaskConfig{}

	if err := env.Parse(&taskConfig); err != nil {
		log.Error("Can not parse env from file system, please check the env.")
		panic("Can not parse env from file system, please check the env.")
	}
	return
}
