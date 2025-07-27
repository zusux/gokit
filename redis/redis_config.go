package redis

type ConfigRedis struct {
	Host string `yaml:"host"`
	Pass string `yaml:"pass"`
	Db   int    `yaml:"db"`
}
