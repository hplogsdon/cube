package main

import (
	"cube/task"
	"cube/worker"
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"time"
)

func main() {
	db := make(map[uuid.UUID]*task.Task)
	w := worker.Worker{
		Queue: *queue.New(),
		Db:    db,
	}

	t := task.Task{
		ID:    uuid.New(),
		Name:  "Task-1",
		State: task.Scheduled,
		Image: "strm/helloworld-http",
	}

	w.AddTask(t)
	result := w.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}

	t.ContainerID = result.ContainerId
	fmt.Printf("task %s is running in container: %s\n", t.ID, t.ContainerID)
	fmt.Printf("Sleeping\n")
	time.Sleep(time.Second * 30)

	fmt.Printf("stopping task %s\n", t.ID)
	result = w.StopTask(t)
	if result.Error != nil {
		panic(result.Error)
	}
}
