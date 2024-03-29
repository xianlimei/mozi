package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-done/mozi/jober"
	"github.com/go-done/mozi/jober/queue"
	"github.com/go-done/mozi/jober/structs"
)

// Payload just for test
type Payload struct {
	Element []int `json:"element"`
}

func main() {
	q := queue.NewRedisQueue("localhost:6379", "", 0)
	jer := jober.NewJober(filepath.Join(".", "tasks"), q)
	go jer.Start()

	// time.Sleep(2 * time.Second)

	for index := 0; index < 10; index++ {

		// time.Sleep(1 * time.Second)
		fmt.Println("******* send math default job *******")
		n := []int{1, 2, 3, 4, 5}
		pld := &Payload{
			Element: n,
		}
		plgb, err := json.Marshal(pld)
		if err != nil {
			return
		}
		args2 := &structs.JobArgs{
			Name: "Math",
			Args: plgb,
		}
		jobBody, _ := json.Marshal(args2)
		err = jer.AddJob(jobBody)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("******* send math sum job *******")
		args2 = &structs.JobArgs{
			Name: "Math.Sum",
			Args: plgb,
		}
		jobBody, _ = json.Marshal(args2)
		err = jer.AddJob(jobBody)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("******* send math Multiply job *******")
		args2 = &structs.JobArgs{
			Name: "Math.Multiply",
			Args: plgb,
		}
		jobBody, _ = json.Marshal(args2)
		err = jer.AddJob(jobBody)
		if err != nil {
			fmt.Println(err)
		}

	}

	time.Sleep(30 * time.Second)
}
