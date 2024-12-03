package task

import (
	"context"
	"sync"
	"time"
)

type Config struct {
	MaxCache       int   `json:"max_cache" yaml:"max_cache" bson:"max_cache" mapstructure:"max_cache"`
	MaxConcurrency int32 `json:"max_concurrency" yaml:"max_concurrency" bson:"max_concurrency" mapstructure:"max_concurrency"`
	MinConcurrency int32 `json:"min_concurrency" yaml:"min_concurrency" bson:"min_concurrency" mapstructure:"min_concurrency"`
	MonitorTime    int32 `json:"monitor_time" yaml:"monitor_time" bson:"monitor_time" mapstructure:"monitor_time"`
}

type Instance struct {
	config *Config

	ctx    context.Context
	cancel context.CancelFunc

	queue chan *Raw // 任务队列
	stop  chan int  // 用于通知routine停止的信号chan

	routineCount int32 // 使用原子操作来更新worker数量

	twg sync.WaitGroup // task任务异步函数

	withAddTaskSuccess func(id string)
	withAddTaskError   func(err error)
	withRunTask        func(id string, et time.Duration)
	withAddRoutine     func()
	withRemoveRoutine  func()
}

type Raw struct {
	Id     string
	Handle func()
}
