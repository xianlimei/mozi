package queue

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-done/mozi/jober/structs"
	"github.com/go-redis/redis"
)

const queuekey = "job_queue"

type redisqueue struct {
	sendBuf  []*structs.JobArgs // buffer for Add func
	popBuf   []*structs.JobArgs // buffer for Pop func
	si       int                // current index for sendBuf
	pi       int                // current index for popBuf
	client   *redis.Client      // redis client
	redisOpt *redis.Options     // redis option
	addlock  *sync.Mutex
	poplock  *sync.Mutex
}

// NewRedisQueue create a new redis queue
func NewRedisQueue(addr, password string, db int) Queue {
	opt := &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	}
	client := redis.NewClient(opt)
	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err)
		panic("redis connect failed. addr: " + addr + "password: " + password + "db" + string(db))
	}

	return &redisqueue{
		sendBuf:  make([]*structs.JobArgs, 20),
		popBuf:   make([]*structs.JobArgs, 20),
		si:       -1,
		pi:       -1,
		client:   client,
		redisOpt: opt,
		addlock:  new(sync.Mutex),
		poplock:  new(sync.Mutex),
	}
}

func (r *redisqueue) Add(job *structs.JobArgs) error {
	r.addlock.Lock()
	defer r.addlock.Unlock()
	jobName := job.Name
	bt, err := json.Marshal(job)
	if err != nil {
		return err
	}

	_, err = r.client.RPush(getQueueKey(jobName), base64.StdEncoding.EncodeToString(bt)).Result()
	return err
}

func getQueueKey(jobName string) string {
	// return fmt.Sprintf("queue::%s", jobName)
	return queuekey
}

func (r *redisqueue) Pop() (*structs.JobArgs, error) {
	jobArgs := &structs.JobArgs{}
	res, err := r.client.BLPop(0, queuekey).Result()
	if err != nil {
		return nil, err
	}
	body, err := base64.StdEncoding.DecodeString(res[1])
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, jobArgs)
	if err != nil {
		return nil, err
	}

	return jobArgs, nil
}
