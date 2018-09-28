package redis

import (
	"strconv"
	"time"

	"github.com/splitio/split-synchronizer/conf"
	redis "gopkg.in/redis.v5"
)

// Client is a redis client with a connection pool
var Client *redis.ClusterClient

// BaseStorageAdapter basic redis storage adapter
type BaseStorageAdapter struct {
	*prefixAdapter
	client *redis.ClusterClient
}

// Initialize Redis module with a pool connection
func Initialize(redisOptions conf.RedisSection) error {
	var err error
	Client, err = NewInstance(redisOptions)
	return err
}

// NewInstance returns an instance of Redis Client
func NewInstance(opt conf.RedisSection) (*redis.ClusterClient, error) {
//	if !opt.SentinelReplication {
		return redis.NewClusterClient(
			&redis.ClusterOptions{
			//	Network:      opt.Network,
				Addrs:        []string{opt.Host, strconv.FormatInt(int64(opt.Port), 10)},
				Password:     opt.Pass,
//				DB:           opt.Db,
			//	MaxRetries:   opt.MaxRetries,
				PoolSize:     opt.PoolSize,
				DialTimeout:  time.Duration(opt.DialTimeout) * time.Second,
				ReadTimeout:  time.Duration(opt.ReadTimeout) * time.Second,
				WriteTimeout: time.Duration(opt.WriteTimeout) * time.Second,
			}), nil
/*	}

	if opt.SentinelMaster == "" {
		return nil, errors.New("Missing redis sentinel master name")
	}

	if opt.SentinelAddresses == "" {
		return nil, errors.New("Missing redis sentinels addresses")
	}

	addresses := strings.Split(opt.SentinelAddresses, ",")

	return redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    opt.SentinelMaster,
		SentinelAddrs: addresses,
		Password:      opt.Pass,
		DB:            opt.Db,
		MaxRetries:    opt.MaxRetries,
		PoolSize:      opt.PoolSize,
		DialTimeout:   time.Duration(opt.DialTimeout) * time.Second,
		ReadTimeout:   time.Duration(opt.ReadTimeout) * time.Second,
		WriteTimeout:  time.Duration(opt.WriteTimeout) * time.Second,
	}), nil*/
}
