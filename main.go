package main

import (
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
	api.Start()
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
