package task

import "time"

func (core *Instance) WithAddTaskSuccess(handle func(id string)) {
	core.withAddTaskSuccess = handle
}
func (core *Instance) WithAddTaskError(handle func(err error)) {
	core.withAddTaskError = handle
}
func (core *Instance) WithRunTask(handle func(id string, et time.Duration)) {
	core.withRunTask = handle
}

func (core *Instance) WithAddRoutine(handle func()) {
	core.withAddRoutine = handle
}
func (core *Instance) WithRemoveRoutine(handle func()) {
	core.withRemoveRoutine = handle
}
