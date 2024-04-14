## Task Go
A minimalist task scheduler.

### Example
``` go
package main

import (
	"fmt"
	"github.com/lhdhtrc/task-go/core"
	"github.com/lhdhtrc/task-go/model"
	"go.uber.org/zap"
	"time"
)

func main() {
	logger, _ := zap.NewDevelopment()
	timeFormat := "2006-01-02 15:04:05"

	task := core.New(model.ConfigEntity{
		MaxCache:       1000000,
		MaxConcurrency: 10,
		MinConcurrency: 3,
	})
	task.WithAddTaskSuccess(func(id string) {
		now := time.Now()
		var field []zap.Field
		field = append(field, zap.String("TaskId", id))
		field = append(field, zap.String("AddDate", now.Format(timeFormat)))
		logger.Info("Task Add Success", field...)
	})
	task.WithRunTask(func(id string, et time.Duration, err error) {
		now := time.Now()
		var field []zap.Field
		field = append(field, zap.String("TaskId", id))
		field = append(field, zap.String("RunDate", now.Format(timeFormat)))
		field = append(field, zap.String("TimeConsuming", fmt.Sprintf("%v", et)))
		logger.Info("Task Run Success", field...)
	})
	task.Setup()

	for i := 0; i < 10; i++ {
		jobID := i
		task.Add(core.TaskEntity{
			Id: fmt.Sprintf("TestTask_%d", jobID),
			Handle: func() error {
				fmt.Printf("Starting job %d\n", jobID)
				time.Sleep(time.Second) // 模拟耗时任务
				fmt.Printf("Finished job %d\n", jobID)
				return nil
			},
		})
	}

	time.Sleep(5 * time.Second)
}
```