package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-done/mozi/jober"
)

type payload struct {
	Element []int
}

func main() {
	jer := jober.NewJober(filepath.Join(".", "tasks"), filepath.Join(".", "sos"))
	jer.Start()

	time.Sleep(2 * time.Second)
	// jober.
	for index := 0; index < 100; index++ {

		time.Sleep(1 * time.Second)
		fmt.Println("******* send math default job *******")
		n := []int{1, 2, 3, 4, 5}
		pld := &payload{
			Element: n,
		}
		plgb, err := json.Marshal(pld)
		if err != nil {
			return
		}
		args2 := &jober.JobArgs{
			Name: "Math",
			Args: plgb,
		}
		err = jer.AddJob(args2)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("******* send math sum job *******")
		args2 = &jober.JobArgs{
			Name: "Math.Sum",
			Args: plgb,
		}
		err = jer.AddJob(args2)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("******* send math Multiply job *******")
		args2 = &jober.JobArgs{
			Name: "Math.Multiply",
			Args: plgb,
		}
		err = jer.AddJob(args2)
		if err != nil {
			fmt.Println(err)
		}

	}
	time.Sleep(30 * time.Second)

	jer.Clear()
}
