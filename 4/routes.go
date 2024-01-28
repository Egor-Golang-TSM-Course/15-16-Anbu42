package chiserver

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes() http.Handler {
	dataStore := &DataStore{
		Tasks: sync.Map{},
		Users: sync.Map{},
	}

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(timeoutMiddleware) // Добавляем middleware для установки таймаута

	// Routes
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", GetTasks(dataStore))
		r.Get("/{id}", GetTask(dataStore))
		r.Post("/", CreateTask(dataStore))
		r.Put("/{id}", UpdateTask(dataStore))
		r.Delete("/{id}", DeleteTask(dataStore))
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/", GetUsers(dataStore))
		r.Get("/{id}", GetUser(dataStore))
		r.Post("/", CreateUser(dataStore))
		r.Put("/{id}", UpdateUser(dataStore))
		r.Delete("/{id}", DeleteUser(dataStore))
	})

	return r
}
