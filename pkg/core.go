package task

import (
	"context"
	"errors"
	"fmt"
)

func New(config *ConfigEntity) *CoreEntity {
	ctx, cancel := context.WithCancel(context.Background())

	core := &CoreEntity{
		ConfigEntity: *config,

		queue:  make(chan *RawEntity, config.MaxCache),
		stop:   make(chan int),
		ctx:    ctx,
		cancel: cancel,
	}
	if core.ConfigEntity.MaxCache == 0 {
		core.ConfigEntity.MaxCache = 1000
	}
	if core.ConfigEntity.MaxConcurrency == 0 {
		core.ConfigEntity.MaxCache = 5
	}
	if core.ConfigEntity.MinConcurrency == 0 {
		core.ConfigEntity.MaxCache = 1
	}
	if core.ConfigEntity.MonitorTime == 0 {
		core.ConfigEntity.MonitorTime = 100
	}

	core.install()

	return core
}

func (core *CoreEntity) Add(task *RawEntity) {
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

func (core *CoreEntity) Await() {
	core.twg.Wait()
}

func (core *CoreEntity) Count() int {
	return len(core.queue)
}

func (core *CoreEntity) RoutineCount() int32 {
	return core.routineCount
}

func (core *CoreEntity) Uninstall() {
	core.twg.Wait() // 等待所有任务执行完成

	core.cancel()     // 取消上下文，通知所有routine停止工作
	close(core.stop)  // 发送信号到所有routine
	close(core.queue) // 关闭任务队列

	fmt.Println("uninstall task success")
}
