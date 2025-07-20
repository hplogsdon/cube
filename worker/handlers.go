package worker

import (
	"cube/task"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type ErrResponse struct {
	HTTPStatusCode int    `json:"status"`
	Message        string `json:"message"`
}

func (a *Api) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	evt := task.TaskEvent{}
	err := d.Decode(&evt)
	if err != nil {
		msg := fmt.Sprintf("Error deserializing task event: %v", err)
		log.Printf(msg)
		w.WriteHeader(http.StatusBadRequest)
		e := ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			Message:        msg,
		}
		json.NewEncoder(w).Encode(e)
		return
	}

	a.Worker.AddTask(evt.Task)
	log.Printf("Added task %v\n", evt.Task.ID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(evt.Task)
}

func (a *Api) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(a.Worker.GetTasks())
}

func (a *Api) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskId := chi.URLParam(r, "taskID")
	if taskId == "" {
		log.Printf("No taskId found in request\n")
		w.WriteHeader(http.StatusBadRequest)
	}

	tId, _ := uuid.Parse(taskId)
	_, ok := a.Worker.Db[tId]
	if !ok {
		log.Printf("Task %v does not exist\n", tId)
		w.WriteHeader(http.StatusNotFound)
	}
	taskToStop := a.Worker.Db[tId]
	taskCopy := *taskToStop
	taskCopy.State = task.Completed
	a.Worker.StopTask(taskCopy)
	w.WriteHeader(http.StatusNoContent)
}
