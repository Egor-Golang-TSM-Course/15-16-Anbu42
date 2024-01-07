package chiserver

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetTasks(dataStore *DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		default:
			tasks := make([]Task, 0)
			dataStore.Tasks.Range(func(key, value interface{}) bool {
				tasks = append(tasks, value.(Task))
				return true
			})

			respond(w, http.StatusOK, tasks)
		case <-r.Context().Done():
			log.Println("GetTasks: Request timeout")
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
		}
	}
}

func GetTask(dataStore *DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		default:
			taskIDStr := chi.URLParam(r, "id")

			if taskIDStr == "" {
				respond(w, http.StatusBadRequest, "Task ID is required")
				return
			}

			taskID, err := strconv.Atoi(taskIDStr)
			if err != nil {
				respond(w, http.StatusBadRequest, "Invalid Task ID")
				return
			}

			if task, ok := dataStore.Tasks.Load(taskID); ok {
				respond(w, http.StatusOK, task)
			} else {
				respond(w, http.StatusNotFound, nil)
			}
		case <-r.Context().Done():
			log.Println("GetTask: Request timeout")
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
		}
	}
}

func CreateTask(dataStore *DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// time.Sleep(2 * time.Second)
		select {
		default:
			var task Task
			if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
				log.Println("CreateTask: Invalid JSON payload")
				respond(w, http.StatusBadRequest, "Invalid JSON payload")
				return
			}

			dataStore.Tasks.Store(task.ID, task)

			log.Printf("CreateTask: Creating new task - ID: %d, Name: %s\n", task.ID, task.Name)
			respond(w, http.StatusCreated, nil)
		case <-r.Context().Done():
			log.Println("CreateTask: Request timeout")
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
		}
	}
}

func UpdateTask(dataStore *DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		default:
			taskIDStr := chi.URLParam(r, "id")

			if taskIDStr == "" {
				respond(w, http.StatusBadRequest, "Task ID is required")
				return
			}

			taskID, err := strconv.Atoi(taskIDStr)
			if err != nil {
				respond(w, http.StatusBadRequest, "Invalid Task ID")
				return
			}

			if _, ok := dataStore.Tasks.Load(taskID); ok {
				var updatedTask Task
				if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
					log.Println("UpdateTask: Invalid JSON payload")
					respond(w, http.StatusBadRequest, "Invalid JSON payload")
					return
				}
				dataStore.Tasks.Store(taskID, updatedTask)

				log.Printf("UpdateTask: Updating task with ID %d\n", taskID)
				respond(w, http.StatusOK, nil)
			} else {
				respond(w, http.StatusNotFound, nil)
			}
		case <-r.Context().Done():
			log.Println("UpdateTask: Request timeout")
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
		}
	}
}

func DeleteTask(dataStore *DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		default:
			taskIDStr := chi.URLParam(r, "id")

			if taskIDStr == "" {
				respond(w, http.StatusBadRequest, "Task ID is required")
				return
			}

			taskID, err := strconv.Atoi(taskIDStr)
			if err != nil {
				respond(w, http.StatusBadRequest, "Invalid Task ID")
				return
			}

			if _, ok := dataStore.Tasks.Load(taskID); ok {
				dataStore.Tasks.Delete(taskID)

				log.Printf("DeleteTask: Deleting task with ID %d\n", taskID)
				respond(w, http.StatusOK, nil)
			} else {
				respond(w, http.StatusNotFound, nil)
			}
		case <-r.Context().Done():
			log.Println("DeleteTask: Request timeout")
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
		}
	}
}
