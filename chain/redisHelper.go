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
type BlockLog struct {
	Client     string `json:"client"`
	Cid        string `json:"cid"`
	MsgCount   uint64 `json:"amount"`
	FirstKnown int64  `json:"firstKnown"`
	Accepted   int64  `json:"accepted"`
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
func (rh *RedisHelper) RedisBlockMsgCount(cid string, count uint64) {
	rh.redisInitClient()
	if !rh.redisDo {
		return
	}
	hostname, _ := os.Hostname()

	var stored BlockLog

	val, err := rh.redisClient.Get(cid + "-b-" + hostname).Result()
	if err != nil {
		if err != redis.Nil {
			fmt.Printf("RedisBlockMsgCount.redisClient.get: %v\n", err)
		}
		stored = BlockLog{Client: hostname, Cid: cid, MsgCount: count}

	} else {

		err = json2.Unmarshal([]byte(val), &stored)
		if err != nil {
			fmt.Printf("RedisBlockMsgCount.json.Unmarshal %v gives %v\n", val, err)
		}

	}

	json, err := json2.Marshal(BlockLog{Client: hostname, Cid: cid, FirstKnown: stored.FirstKnown, Accepted: stored.Accepted, MsgCount: count})
	if err != nil {
		fmt.Printf("RedisBlockMsgCount.json.Marshal %v\n", err)
	}
	err = rh.redisClient.Set(cid+"-"+hostname, json, 0).Err()
	if err != nil {
		fmt.Printf("RedisBlockMsgCount.redisClient.set: %v\n", err)
	}
}
func (rh *RedisHelper) RedisBlockFirstKnown(cid string, time int64) {
	rh.redisInitClient()
	if !rh.redisDo {
		return
	}
	hostname, _ := os.Hostname()
	json, err := json2.Marshal(BlockLog{Client: hostname, Cid: cid, FirstKnown: time})
	if err != nil {
		fmt.Printf("RedisBlockFirstKnown.json.Marshal: %v\n", err)
	}
	err = rh.redisClient.Set(cid+"-b-"+hostname, json, 0).Err()
	if err != nil {
		fmt.Printf("RedisBlockFirstKnown.redisClient.set: %v\n", err)
	}
}
func (rh *RedisHelper) RedisBlockApproved(cid string, time int64) {
	rh.redisInitClient()
	if !rh.redisDo {
		return
	}
	hostname, _ := os.Hostname()

	var stored BlockLog

	val, err := rh.redisClient.Get(cid + "-b-" + hostname).Result()
	if err != nil {
		if err != redis.Nil {
			fmt.Printf("RedisBlockApproved.redisClient.get: %v\n", err)
		}
		stored = BlockLog{Client: hostname, Cid: cid, FirstKnown: 0, Accepted: time}

	} else {

		err = json2.Unmarshal([]byte(val), &stored)
		if err != nil {
			fmt.Printf("RedisBlockApproved.json.Unmarshal %v gives %v\n", val, err)
		}

	}

	json, err := json2.Marshal(BlockLog{Client: hostname, Cid: cid, FirstKnown: stored.FirstKnown, Accepted: time})
	if err != nil {
		fmt.Printf("RedisBlockApproved.json.Marshal %v\n", err)
	}
	err = rh.redisClient.Set(cid+"-"+hostname, json, 0).Err()
	if err != nil {
		fmt.Printf("RedisBlockApproved.redisClient.set: %v\n", err)
	}
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
