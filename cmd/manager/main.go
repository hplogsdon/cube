package main

import (
	"cube/manager"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	host := os.Getenv("CUBE_MANAGER_HOST")
	port, _ := strconv.Atoi(os.Getenv("CUBE_MANAGER_PORT"))
	workers := strings.Split(os.Getenv("CUBE_WORKERS"), ",")

	m := manager.NewManager(workers)
	mApi := manager.Api{
		Address: host,
		Port:    port,
		Manager: m,
	}
	fmt.Printf("Starting Manager at %s:%d\n", host, port)

	go m.ProcessTasks()
	go m.UpdateTasks()
	mApi.Start()
}
