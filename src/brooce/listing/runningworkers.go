package listing

import (
	"encoding/json"

	"brooce/heartbeat"

	redis "gopkg.in/redis.v3"
)

func RunningWorkers() (workers []*heartbeat.HeartbeatType, err error) {
	var keys []string
	keys, err = redisClient.Keys(redisHeader + ":workerprocs:*").Result()
	if err != nil {
		return
	}

	var heartbeatStrs []*redis.StringCmd
	_, err = redisClient.Pipelined(func(pipe *redis.Pipeline) error {
		for _, key := range keys {
			result := pipe.Get(key)
			heartbeatStrs = append(heartbeatStrs, result)
		}
		return nil
	})
	if err != nil {
		return
	}

	for _, str := range heartbeatStrs {
		worker := &heartbeat.HeartbeatType{}
		err = json.Unmarshal([]byte(str.Val()), worker)
		if err != nil {
			return
		}

		workers = append(workers, worker)
	}

	return
}
