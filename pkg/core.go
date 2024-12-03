package task

import (
	"context"
	"errors"
	"fmt"
)

func New(config *Config) *Instance {
	ctx, cancel := context.WithCancel(context.Background())

	core := &Instance{
		ctx:    ctx,
		cancel: cancel,
		config: config,

		queue: make(chan *Raw, config.MaxCache),
		stop:  make(chan int),
	}
	if core.config.MaxCache == 0 {
		core.config.MaxCache = 1000
	}
	if core.config.MaxConcurrency == 0 {
		core.config.MaxCache = 5
	}
	if core.config.MinConcurrency == 0 {
		core.config.MaxCache = 1
	}
	if core.config.MonitorTime == 0 {
		core.config.MonitorTime = 100
	}

	core.install()

	return core
}

func (core *Instance) Add(task *Raw) {
	select {
	case core.queue <- task:
		core.twg.Add(1)
		// 任务成功提交到队列
		if core.withAddTaskSuccess != nil {
			core.withAddTaskSuccess(task.Id)
		}
	default:
		// 队列已满，可以选择丢弃任务或进行其他处理
		if core.withAddTaskError != nil {
			core.withAddTaskError(errors.New("queue is full, job was not submitted"))
		}
	}
}

func (core *Instance) Await() {
	core.twg.Wait()
}

func (core *Instance) Count() int {
	return len(core.queue)
}

func (core *Instance) RoutineCount() int32 {
	return core.routineCount
}

func (core *Instance) Uninstall() {
	core.twg.Wait() // 等待所有任务执行完成

	core.cancel()     // 取消上下文，通知所有routine停止工作
	close(core.stop)  // 发送信号到所有routine
	close(core.queue) // 关闭任务队列

	fmt.Println("uninstall task success")
}
