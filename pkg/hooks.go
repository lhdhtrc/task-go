package task

import "time"

func (core *CoreEntity) WithAddTaskSuccess(handle func(id string)) {
	core.withAddTaskSuccess = handle
}
func (core *CoreEntity) WithAddTaskError(handle func(err error)) {
	core.withAddTaskError = handle
}
func (core *CoreEntity) WithRunTask(handle func(id string, et time.Duration)) {
	core.withRunTask = handle
}

func (core *CoreEntity) WithAddRoutine(handle func()) {
	core.withAddRoutine = handle
}
func (core *CoreEntity) WithRemoveRoutine(handle func()) {
	core.withRemoveRoutine = handle
}
