package core

import (
	"context"
	"errors"
	"github.com/lhdhtrc/task-go/model"
	"sync"
	"time"
)

type TaskCoreEntity struct {
	model.ConfigEntity

	queue chan model.TaskEntity // 任务队列
	stop  chan int              // 用于通知routine停止的信号chan

	ctx    context.Context
	cancel context.CancelFunc

	routineCount int32 // 使用原子操作来更新worker数量

	wg sync.WaitGroup
	mu sync.Mutex // 用于保护routines数量检查时的竞态关系

	withAddTaskSuccess func(id string)
	withAddTaskError   func(err error)
	withRunTask        func(id string, et time.Duration, err error)
	withAddRoutine     func()
	withRemoveRoutine  func()
}

func New(config model.ConfigEntity) *TaskCoreEntity {
	ctx, cancel := context.WithCancel(context.Background())

	return &TaskCoreEntity{
		ConfigEntity: config,

		queue:  make(chan model.TaskEntity, config.MaxCache),
		stop:   make(chan int),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *TaskCoreEntity) Setup() {
	// 启动初始数量的routine
	for i := int32(0); i < s.MinConcurrency; i++ {
		s.addRoutine()
	}

	go s.monitorTask()
}

func (s *TaskCoreEntity) Add(task model.TaskEntity) {
	select {
	case s.queue <- task:
		// 任务成功提交到队列
		if s.withAddTaskSuccess != nil {
			s.withAddTaskSuccess(task.Id)
		}
	default:
		// 队列已满，可以选择丢弃任务或进行其他处理
		if s.withAddTaskError != nil {
			s.withAddTaskError(errors.New("queue is full, job was not submitted"))
		}
	}
}

func (s *TaskCoreEntity) Count() int {
	return len(s.queue)
}

func (s *TaskCoreEntity) Stop() {
	s.cancel()     // 取消上下文，通知所有routine停止工作
	close(s.stop)  // 发送信号到所有routine
	s.wg.Wait()    // 等待所有routine退出
	close(s.queue) // 关闭任务队列
}

func (s *TaskCoreEntity) RoutineCount() int32 {
	return s.routineCount
}
