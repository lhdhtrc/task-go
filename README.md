## Task Go
A minimalist task scheduler.

### How to use it?
`go get github.com/lhdhtrc/task-go`
```go
package main

import (
	"fmt"
	task "github.com/lhdhtrc/task-go/pkg"
	"time"
)

func main() {
	config := &task.Config{
		MaxCache:       1000000, // Set a buffer large enough for your business needs, because if you Add more data at once, the task will be discarded
		MaxConcurrency: 50,
		MinConcurrency: 1,
	}
	
    instance := task.New(config)
    instance.WithRunTask(func(id string, et time.Duration) {
        fmt.Printf("%s success, run time %s\n", id, et)
        fmt.Println(instance.RoutineCount())
    })
    instance.WithAddTaskError(func(err error) {
        fmt.Println(err.Error())
    })
    
    // How to add a task to a Task queue (asynchronous)?
    for i := 0; i < config.MaxCache; i++ {
        instance.Add(&task.RawEntity{
            Id:     fmt.Sprintf("%s_%d", "task", i+1),
            Handle: TaskHandle,
        })
    }
    
    // How do I wait for an asynchronous task to finish
    instance.Await()
    
    // Note that at the end of the process, please reclaim your lease!
    instance.Uninstall()
}
```

### Finally
- If you feel good, click on star.
- If you have a good suggestion, please ask the issue.