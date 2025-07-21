package main

import (
	"cube/manager"
	"cube/task"
	"cube/worker"
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"os"
	"strconv"
)

func main() {
	whost := os.Getenv("CUBE_WORKER_HOST")
	wport, _ := strconv.Atoi(os.Getenv("CUBE_WORKER_PORT"))
	if whost == "" || wport == 0 {
		whost = "localhost"
		wport = 5678
	}

	mhost := os.Getenv("CUBE_MANAGER_HOST")
	mport, _ := strconv.Atoi(os.Getenv("CUBE_MANAGER_PORT"))
	if mhost == "" || mport == 0 {
		mhost = "localhost"
		mport = 5555
	}

	fmt.Printf("Starting Cube Worker\n")

	w := worker.Worker{
		Name:  "worker-1",
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	api := worker.Api{Address: whost, Port: wport, Worker: &w}

	go w.RunTasks()
	go w.CollectStats()
	go api.Start()

	workers := []string{fmt.Sprintf("%s:%d", whost, wport)}
	m := manager.NewManager(workers)
	mApi := manager.Api{
		Address: mhost,
		Port:    mport,
		Manager: m,
	}
	fmt.Printf("Starting Manager at %s:%d\n", mhost, mport)

	go m.ProcessTasks()
	go m.UpdateTasks()
	mApi.Start()
}
