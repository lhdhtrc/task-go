package task

import (
	"context"
	"sync"
	"time"
)

type ConfigEntity struct {
	MaxCache       int   `json:"max_cache" yaml:"max_cache" bson:"max_cache" mapstructure:"max_cache"`
	MaxConcurrency int32 `json:"max_concurrency" yaml:"max_concurrency" bson:"max_concurrency" mapstructure:"max_concurrency"`
	MinConcurrency int32 `json:"min_concurrency" yaml:"min_concurrency" bson:"min_concurrency" mapstructure:"min_concurrency"`
}

type CoreEntity struct {
	ConfigEntity

	ctx    context.Context
	cancel context.CancelFunc

	queue chan RawEntity // 任务队列
	stop  chan int       // 用于通知routine停止的信号chan

	routineCount int32 // 使用原子操作来更新worker数量

	wg sync.WaitGroup
	mu sync.Mutex // 用于保护routines数量检查时的竞态关系

	withAddTaskSuccess func(id string)
	withAddTaskError   func(err error)
	withRunTask        func(id string, et time.Duration, err error)
	withAddRoutine     func()
	withRemoveRoutine  func()
}

type RawEntity struct {
	Id     string
	Handle func() error
}
