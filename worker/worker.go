package worker

import (
	"cube/task"
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Worker struct {
	Name      string
	Queue     queue.Queue
	Db        map[uuid.UUID]*task.Task
	TaskCount int
}

func (w *Worker) CollectStats() {
	fmt.Println("CollectStats")
}

func (w *Worker) RunTask() {
	fmt.Println("RunTask")
}

func (w *Worker) StartTask() {
	fmt.Println("StartTask")
}

func (w *Worker) StopTask() {
	fmt.Println("StopTask")
}
