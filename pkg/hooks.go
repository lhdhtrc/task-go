package task

import "time"

func (s *CoreEntity) WithAddTaskSuccess(handle func(id string)) {
	s.withAddTaskSuccess = handle
}
func (s *CoreEntity) WithAddTaskError(handle func(err error)) {
	s.withAddTaskError = handle
}
func (s *CoreEntity) WithRunTask(handle func(id string, et time.Duration, err error)) {
	s.withRunTask = handle
}

func (s *CoreEntity) WithAddRoutine(handle func()) {
	s.withAddRoutine = handle
}
func (s *CoreEntity) WithRemoveRoutine(handle func()) {
	s.withRemoveRoutine = handle
}
