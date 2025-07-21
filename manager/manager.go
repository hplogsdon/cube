package manager

import (
	"bytes"
	"cube/task"
	"cube/worker"
	"encoding/json"
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type Manager struct {
	Pending       queue.Queue
	TaskDb        map[uuid.UUID]*task.Task
	EventDb       map[uuid.UUID]*task.TaskEvent
	Workers       []string
	WorkerTaskMap map[string][]uuid.UUID
	TaskWorkerMap map[uuid.UUID]string
	LastWorker    int
}

func (m *Manager) SelectWorker() string {
	var newWorker int
	if m.LastWorker+1 < len(m.Workers) {
		newWorker = m.LastWorker + 1
		m.LastWorker++
	} else {
		newWorker = 0
		m.LastWorker = 0
	}

	return m.Workers[newWorker]
}

func (m *Manager) SendWork() {
	if m.Pending.Len() > 0 {
		w := m.SelectWorker()
		e := m.Pending.Dequeue()
		te := e.(task.TaskEvent)
		t := te.Task
		log.Printf("Pulled %v off Pending queue\n", t)

		// update the worker task map, allowing manager to identify tasks it has
		// assigned to the workers
		m.EventDb[te.ID] = &te
		m.WorkerTaskMap[w] = append(m.WorkerTaskMap[w], te.Task.ID)
		m.TaskWorkerMap[t.ID] = w

		t.State = task.Scheduled
		m.TaskDb[t.ID] = &t

		data, err := json.Marshal(te)
		if err != nil {
			log.Printf("Failed to marshal task: %v\n", t)
			return
		}

		url := fmt.Sprintf("http://%s/tasks", w)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Printf("error connecting to worker %s: %v\n", w, err)
			m.Pending.Enqueue(t)
			return
		}

		d := json.NewDecoder(resp.Body)
		if resp.StatusCode != http.StatusCreated {
			e := worker.ErrResponse{}
			err := d.Decode(&e)
			if err != nil {
				log.Printf("error decoding response from worker %s: %v\n", w, err)
				return
			}
			log.Printf("error response (%d) %s\n", e.HTTPStatusCode, e.Message)
			return
		}
		t = task.Task{}
		err = d.Decode(&t)
		if err != nil {
			log.Printf("error decoding response from worker %s: %v\n", w, err)
			return
		}
		log.Printf("%#v\n", t)
	} else {
		log.Printf("No work in queue\n")
	}
}

func (m *Manager) UpdateTasks() {
	for _, w := range m.Workers {
		log.Printf("Updating w %s\n", w)
		url := fmt.Sprintf("http://%s/tasks", w)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("error connecting to w %s: %v\n", w, err)
			return
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("error sending request (%d) %v\n", resp.StatusCode, err)
			return
		}

		d := json.NewDecoder(resp.Body)
		var tasks []*task.Task
		if err := d.Decode(&tasks); err != nil {
			log.Printf("error deserializing response from w %s: %s\n", w, err.Error())
			return
		}
		for _, task := range tasks {
			log.Printf("Updating task %v\n", task)
			_, ok := m.TaskDb[task.ID]
			if !ok {
				log.Printf("Task with id %v not found\n", task.ID)
				return
			}

			if m.TaskDb[task.ID].State != task.State {
				m.TaskDb[task.ID].State = task.State
			}

			m.TaskDb[task.ID].StartTime = task.StartTime
			m.TaskDb[task.ID].EndTime = task.EndTime
			m.TaskDb[task.ID].ContainerID = task.ContainerID
		}
	}
}

func (m *Manager) AddTask(te task.TaskEvent) {
	m.Pending.Enqueue(te)
}

func NewManager(workers []string) *Manager {
	taskDb := make(map[uuid.UUID]*task.Task)
	eventDb := make(map[uuid.UUID]*task.TaskEvent)
	workerTaskMap := make(map[string][]uuid.UUID)
	taskWorkerMap := make(map[uuid.UUID]string)
	for w := range workers {
		workerTaskMap[workers[w]] = []uuid.UUID{}
	}

	return &Manager{
		Pending:       *queue.New(),
		Workers:       workers,
		TaskDb:        taskDb,
		EventDb:       eventDb,
		WorkerTaskMap: workerTaskMap,
		TaskWorkerMap: taskWorkerMap,
	}
}
