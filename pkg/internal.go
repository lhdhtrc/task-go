package tpg

import (
	"sync/atomic"
	"time"
)

func (core *CoreEntity) install() {
	// 启动初始数量的routine
	for i := int32(0); i < core.MinConcurrency; i++ {
		core.addRoutine()
	}

	go core.monitorRoutine()
}

// monitorQueue 监控任务队列长度并动态调整worker数量
func (core *CoreEntity) monitorRoutine() {
	ticker := time.NewTicker(time.Duration(core.MonitorTime) * time.Millisecond) // 每秒检查一次任务队列长度并调整worker数量
	defer ticker.Stop()

	for {
		select {
		case <-core.ctx.Done():
			return
		case <-ticker.C:
			l := len(core.queue)
			i := atomic.LoadInt32(&core.routineCount) - int32(l)
			if l > 0 && i < core.MaxConcurrency-core.MinConcurrency {
				core.addRoutine()
			} else if i > 0 {
				core.removeRoutine()
			}
		}
	}
}

func (core *CoreEntity) addRoutine() {
	if atomic.LoadInt32(&core.routineCount) >= core.MaxConcurrency {
		// 如果routine数量已经达到最大值，则不添加
		return
	}

	// 增加routine计数，然后启动routine
	atomic.AddInt32(&core.routineCount, 1)
	go core.routine()

	if core.withAddRoutine != nil {
		core.withAddRoutine()
	}
}

func (core *CoreEntity) removeRoutine() {
	if atomic.LoadInt32(&core.routineCount) <= core.MinConcurrency {
		// 如果routine数量已经达到最小值，则不移除
		return
	}

	// 尝试发送停止信号给routine
	select {
	case core.stop <- 1:
		if core.withRemoveRoutine != nil {
			core.withRemoveRoutine()
		}
	default:
		// 没有routine可移除，routine处于忙碌状态
	}
}

func (core *CoreEntity) routine() {
	defer atomic.AddInt32(&core.routineCount, -1) // 无论任务是否完成，routine退出时都减少计数

	for {
		select {
		case task := <-core.queue:
			// 执行任务
			st := time.Now()
			task.Handle()
			et := time.Now()
			dt := et.Sub(st)
			if core.withRunTask != nil {
				core.withRunTask(task.Id, dt)
			}
			core.twg.Done()
		case <-core.stop:
			// 收到信号，routine终止
			return
		case <-core.ctx.Done():
			// 池被关闭，routine退出
			return
		}
	}
}
