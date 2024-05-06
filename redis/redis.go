package redis

import (
	"strings"

	"github.com/redis/go-redis/v9"
)

type ClientRedis struct {
	config *ConfigRedis
}

func NewRedisClient(config *ConfigRedis) *ClientRedis {
	return &ClientRedis{
		config: config,
	}
}

func (r *ClientRedis) splitClusterAddrs(addr string) []string {
	addrs := strings.Split(addr, ",")
	unique := make(map[string]struct{})
	for _, each := range addrs {
		unique[strings.TrimSpace(each)] = struct{}{}
	}
	addrs = addrs[:0]
	for k := range unique {
		addrs = append(addrs, k)
	}
	return addrs
}

func (r *ClientRedis) GetClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     r.config.Host,
		Password: r.config.Pass,
		DB:       r.config.Db,
	})
	return rdb
}

func (r *ClientRedis) GetSentinelClient(masterName string) *redis.Client {
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: r.splitClusterAddrs(r.config.Host),
		Password:      r.config.Pass,
		DB:            r.config.Db,
	})
	return rdb
}

func (r *ClientRedis) GetClusterClient() *redis.ClusterClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    r.splitClusterAddrs(r.config.Host),
		Password: r.config.Pass,
	})
}
