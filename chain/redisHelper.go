package chain

import (
	json2 "encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"sync"
)

type TimeLogEntry struct {
	Client string `json:"client"`
	Cid    string `json:"cid"`
	Start  int64  `json:"start"`
	End    int64  `json:"end"`
}

var rhelper *RedisHelper

type RedisHelper struct {
	redisClient *redis.Client
	redisDo     bool
	redisMutex  sync.Mutex
}

func init() {
	rhelper = &RedisHelper{
		redisClient: nil,
		redisDo:     true,
	}
}
func GetRedisHelper() *RedisHelper {
	return rhelper
}

func (rh *RedisHelper) redisInitClient() {
	rh.redisMutex.Lock()
	if rh.redisClient == nil {
		addr, exist := os.LookupEnv("LOTUS_REDIS_ADDR")
		if !exist {
			rh.redisDo = false
		}
		pw, exist := os.LookupEnv("LOTUS_REDIS_PW")
		if !exist {
			pw = ""
		}

		rh.redisClient = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pw,
			DB:       0,
		})
	}
	rh.redisMutex.Unlock()
}

func (rh *RedisHelper) RedisSaveStartTime(cid string, starttime int64) {
	rh.redisInitClient()
	if !rh.redisDo {
		return
	}
	hostname, _ := os.Hostname()
	json, err := json2.Marshal(TimeLogEntry{Client: hostname, Cid: cid, Start: starttime})
	if err != nil {
		fmt.Printf("RedisSaveStarttime.json.Marshal: %v\n", err)
	}
	err = rh.redisClient.Set(cid+"-"+hostname, json, 0).Err()
	if err != nil {
		fmt.Printf("RedisSaveStarttime.redisClient.set: %v\n", err)
	}
}
func (rh *RedisHelper) RedisSaveEndTime(cid string, endtime int64) {
	rh.redisInitClient()
	if !rh.redisDo {
		return
	}
	hostname, _ := os.Hostname()

	var stored TimeLogEntry

	val, err := rh.redisClient.Get(cid + "-" + hostname).Result()
	if err != nil {
		if err != redis.Nil {
			fmt.Printf("RedisSaveEndTime.redisClient.get: %v\n", err)
		}
		stored = TimeLogEntry{Client: hostname, Cid: cid, Start: 0, End: endtime}

	} else {

		err = json2.Unmarshal([]byte(val), &stored)
		if err != nil {
			fmt.Printf("RedisSaveEndTime.json.Unmarshal %v gives %v\n", val, err)
		}

	}

	json, err := json2.Marshal(TimeLogEntry{Client: hostname, Cid: cid, Start: stored.Start, End: endtime})
	if err != nil {
		fmt.Printf("RedisSaveEndTime.json.Marshal %v\n", err)
	}
	err = rh.redisClient.Set(cid+"-"+hostname, json, 0).Err()
	if err != nil {
		fmt.Printf("RedisSaveEndTime.redisClient.set: %v\n", err)
	}

}
