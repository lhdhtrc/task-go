package core

import "time"

func (s *TaskCoreEntity) WithAddTaskSuccess(handle func(id string)) {
	s.withAddTaskSuccess = handle
}
func (s *TaskCoreEntity) WithAddTaskError(handle func(err error)) {
	s.withAddTaskError = handle
}
func (s *TaskCoreEntity) WithRunTask(handle func(id string, et time.Duration, err error)) {
	s.withRunTask = handle
}

func (s *TaskCoreEntity) WithAddRoutine(handle func()) {
	s.withAddRoutine = handle
}
func (s *TaskCoreEntity) WithRemoveRoutine(handle func()) {
	s.withRemoveRoutine = handle
}
