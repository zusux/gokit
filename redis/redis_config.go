package redis

type ConfigRedis struct {
	Host string `yaml:"host" json:"host"`
	Pass string `yaml:"pass" json:"pass"`
	Db   int    `yaml:"db" json:"db"`
}
