package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// task object
type Task struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Result string `json:"result,omitempty"`
}

// memory tasks
var (
	tasks       = make(map[string]*Task)
	Mutex       = &sync.Mutex{}
	taskResults = make(map[string]string)
)

// функция для тестирования горутин
func longRunningTask(taskID string) {
	log.Printf("Task id :%s starting...", taskID)
	time.Sleep(3 * time.Minute)
	taskResults[taskID] = "Task completed successfully"
	Mutex.Lock()
	tasks[taskID].Status = "completed"
	log.Printf("Task id :%s complete!", taskID)
	tasks[taskID].Result = taskResults[taskID]
	Mutex.Unlock()
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	taskID := fmt.Sprintf("%d", time.Now().UnixNano())
	task := &Task{
		ID:     taskID,
		Status: "in-progress",
	}

	Mutex.Lock()
	tasks[taskID] = task
	Mutex.Unlock()

	go longRunningTask(taskID) // start task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Task ID: %s\n was created..", task.ID)
}

func GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["task_id"]

	Mutex.Lock()
	defer Mutex.Unlock()

	task, exists := tasks[taskID]
	if !exists {
		http.Error(w, "Task not found, please check the correctness of the data you entered", http.StatusNotFound)
		return
	}
	log.Printf("Task id :%s , status: %s!", task.ID, task.Status)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}
