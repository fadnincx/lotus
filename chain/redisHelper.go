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

var redisClient *redis.Client = nil
var redisDo = true
var redisMutex sync.Mutex

func redisInitClient() {
	redisMutex.Lock()
	if redisClient == nil {
		addr, exist := os.LookupEnv("LOTUS_REDIS_ADDR")
		if !exist {
			redisDo = false
		}
		pw, exist := os.LookupEnv("LOTUS_REDIS_PW")
		if !exist {
			pw = ""
		}

		redisClient = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pw,
			DB:       0,
		})
	}
	redisMutex.Unlock()
}

func RedisSaveStartTime(cid string, starttime int64) {
	redisInitClient()
	if !redisDo {
		return
	}
	hostname, _ := os.Hostname()
	json, err := json2.Marshal(TimeLogEntry{Client: hostname, Cid: cid, Start: starttime})
	if err != nil {
		fmt.Println(err)
	}
	err = redisClient.Set(cid+"-"+hostname, json, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}
func RedisSaveEndTime(cid string, endtime int64) {
	redisInitClient()
	if !redisDo {
		return
	}
	hostname, _ := os.Hostname()

	val, err := redisClient.Get(cid + "-" + hostname).Result()
	if err != nil {
		fmt.Println(err)
	}
	var stored TimeLogEntry
	err = json2.Unmarshal([]byte(val), &stored)
	if err != nil {
		fmt.Println(err)
	}

	json, err := json2.Marshal(TimeLogEntry{Client: hostname, Cid: cid, Start: stored.Start, End: endtime})
	if err != nil {
		fmt.Println(err)
	}
	err = redisClient.Set(cid+"-"+hostname, json, 0).Err()
	if err != nil {
		fmt.Println(err)
	}

}
