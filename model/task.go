package model

type TaskEntity struct {
	Id     string
	Handle func() error
}
