package conf

type Config struct {
	DB struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"db"`

	Bilibili struct {
		UserId string `yaml:"userid"`
	} `yaml:"bilibili"`
}
