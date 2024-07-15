package main

import (
	"fmt"
	task "github.com/lhdhtrc/task-go/pkg"
	"time"
)

func TaskHandle() {
	time.Sleep(1 * time.Second)
	fmt.Println("Time-consuming asynchronous tasks!")
}

func main() {
	instance := task.New(&task.ConfigEntity{
		MaxCache:       100000,
		MaxConcurrency: 50,
		MinConcurrency: 1,
	})
	instance.WithRunTask(func(id string, et time.Duration) {
		fmt.Printf("%s success, run time %s\n", id, et)
	})

	// How to add a task to a Task queue (asynchronous)?
	for i := 0; i < 1000; i++ {
		instance.Add(&task.RawEntity{
			Id:     fmt.Sprintf("%s_%d", "task", i+1),
			Handle: TaskHandle,
		})
	}

	// How do I wait for an asynchronous task to finish
	instance.Await()
}
