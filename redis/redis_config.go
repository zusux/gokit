package redis

type ConfigRedis struct {
	Host string `yaml:"Host"`
	Pass string `yaml:"Pass"`
	Db   int    `yaml:"Db"`
}
