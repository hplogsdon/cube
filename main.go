package main

import (
	"cube/manager"
	"cube/task"
	"cube/worker"
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	host := os.Getenv("CUBE_HOST")
	if host == "" {
		host = "localhost"
	}
	portStr := os.Getenv("CUBE_PORT")
	if portStr == "" {
		portStr = "5555"
	}
	port, _ := strconv.Atoi(portStr)

	fmt.Printf("Starting Cube Worker\n")

	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	api := worker.Api{Address: host, Port: port, Worker: &w}
	go runtasks(&w)
	go w.CollectStats()
	go api.Start()

	workers := []string{fmt.Sprintf("%s:%d", host, port)}
	m := manager.NewManager(workers)
	for i := 0; i < 3; i++ {
		t := task.Task{
			ID:    uuid.New(),
			Name:  fmt.Sprintf("test-container-%d", i+1),
			State: task.Scheduled,
			Image: "strm/helloworld-http",
		}
		te := task.TaskEvent{
			ID:    uuid.New(),
			Task:  t,
			State: task.Running,
		}
		m.AddTask(te)
		m.SendWork()
	}

	go func() {
		for {
			fmt.Printf("Updating tasks from %d workers\n", len(m.Workers))
			m.UpdateTasks()
			time.Sleep(11 * time.Second)
		}
	}()

	for {
		for _, t := range m.TaskDb {
			fmt.Printf("[Manager] Task: id: %s, state: %d\n", t.ID, t.State)
			time.Sleep(12 * time.Second)
		}
	}
}

func runtasks(w *worker.Worker) {
	for {
		if w.Queue.Len() != 0 {
			result := w.RunTask()
			if result.Error != nil {
				panic(result.Error)
			}
		} else {
			log.Printf("No tasks\n")
		}
		time.Sleep(time.Second * 10)
	}
}
