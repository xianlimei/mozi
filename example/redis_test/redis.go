package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type test struct {
	Name string `json:"name"`
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "139.219.99.164:8379",
		Password: "visual3d", // no password set
		DB:       0,          // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong)
	if err != nil {
		fmt.Println(err)
		fmt.Println("connect failed")
		return
	}
	t := &test{Name: "ck"}
	d, _ := json.Marshal(t)
	ri, err := client.RPush("test", string(d)).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ri)

	r, err := client.BLPop(1*time.Hour, "test").Result()
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	fmt.Println([]byte(r[0]))
}
