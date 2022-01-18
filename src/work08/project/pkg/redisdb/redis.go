package rdb

import (
	"sync"
	"time"

	"github.com/go-redis/redis"
)

var (
	lock sync.RWMutex
	rd   *redis.Client = nil
)

type RedisConf struct {
	Addr     string
	Password string
	DB       int
}

func RedisDial(config *RedisConf) error {
	lock.RLock()
	if rd != nil {
		lock.RUnlock()
		return nil
	}
	lock.RUnlock()

	lock.Lock()
	defer lock.Unlock()
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return err
	}
	rd = client
	return nil
}

func Set(key string, val string, expTime int32) {
	rd.Set(key, val, time.Duration(expTime)*time.Second)
}

func Get(key string) string {
	val, err := rd.Get(key).Result()
	if err != nil {
		return ""
	}
	return val
}

func Del(key string) {
	rd.Del(key)
}

func SetHash(key string, field string, value interface{}) {
	rd.HSet(key, field, value)
}

func GetHash(key string, field string) (s string, e error) {
	s, e = rd.HGet(key, field).Result()
	return
}

func DelHash(key string, field string) {
	rd.HDel(key, field)
}

func DoExpire(Key string, Time int) {
	rd.Do("EXPIRE", Key, Time)
}

//计算redis total.allocated
func DoInfo() (int64, error) {
	var allocated int64
	result, err := rd.Do("MEMORY", "STATS").Result()
	if err != nil {
		panic(err)
	}
	if stats, ok := result.([]interface{}); ok {
		for i, max := 0, len(stats)-1; i < max; i += 2 {
			if stats[i] == "total.allocated" {
				if allocated, ok = stats[i+1].(int64); ok {
					// Use allocated
				}
			}
		}
	}
	return allocated, nil
}
