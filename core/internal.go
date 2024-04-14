package core

import (
	"sync/atomic"
	"time"
)

// monitorQueue 监控任务队列长度并动态调整worker数量
func (s *TaskCoreEntity) monitorTask() {
	ticker := time.NewTicker(time.Second) // 每秒检查一次任务队列长度并调整worker数量
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			l := len(s.queue)
			i := atomic.LoadInt32(&s.routineCount) - int32(l)
			if l > 0 && i < s.MaxConcurrency-s.MinConcurrency {
				s.addRoutine()
			} else if i > 0 {
				s.removeRoutine()
			}
		}
	}
}

func (s *TaskCoreEntity) addRoutine() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if atomic.LoadInt32(&s.routineCount) >= s.MaxConcurrency {
		// 如果routine数量已经达到最大值，则不添加
		return
	}

	// 增加routine计数，然后启动routine
	atomic.AddInt32(&s.routineCount, 1)
	s.wg.Add(1)
	go s.routine()

	if s.withAddRoutine != nil {
		s.withAddRoutine()
	}
}

func (s *TaskCoreEntity) removeRoutine() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if atomic.LoadInt32(&s.routineCount) <= s.MinConcurrency {
		// 如果routine数量已经达到最小值，则不移除
		return
	}

	// 尝试发送停止信号给routine
	select {
	case s.stop <- 1:
		if s.withRemoveRoutine != nil {
			s.withRemoveRoutine()
		}
	default:
		// 没有routine可移除，routine处于忙碌状态
	}
}

func (s *TaskCoreEntity) routine() {
	defer atomic.AddInt32(&s.routineCount, -1) // 无论任务是否完成，routine退出时都减少计数

	for {
		select {
		case task := <-s.queue:
			// 执行任务
			st := time.Now()
			te := task.Handle()
			et := time.Now()
			dt := et.Sub(st)
			if s.withRunTask != nil {
				s.withRunTask(task.Id, dt, te)
			}
		case <-s.stop:
			// 收到信号，routine终止
			return
		case <-s.ctx.Done():
			// 池被关闭，routine退出
			return
		}
	}
}
